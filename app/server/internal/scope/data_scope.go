package scope

import (
	"context"

	"gorm.io/gorm"

	"boread/internal/model"
)

// DataScopeContext 数据权限上下文 (Service 层组装后传入)
// 设计要点:
//   一次性把"用户能看到的数据范围"算清楚 -> Scope 函数只做无脑过滤
//   消除特殊情况: 任何业务表的 List 都用同一个 ApplyDataScope, 没有 if/else
type DataScopeContext struct {
	UserID       uint64           // 当前用户 id
	IsSuper      bool             // 是否超管 (含 1=全部 范围)
	OwnerOnly    bool             // 是否仅本人 (5)
	DeptIDs      []uint64         // 可访问的部门 id 集合 (2/3/4 计算结果, 空 = 不限制)
}

// BuildDataScope 根据用户角色计算数据权限上下文
// 调用时机: Service 方法收到请求后, 在执行 List 前调一次, 把结果传给 Repository
func BuildDataScope(ctx context.Context, db *gorm.DB, userID uint64) (*DataScopeContext, error) {
	out := &DataScopeContext{UserID: userID}

	// 1. 查用户的角色及其 data_scope, dept_id
	var rows []struct {
		RoleID    uint64           `gorm:"column:role_id"`
		DataScope model.DataScope  `gorm:"column:data_scope"`
		DeptID    *uint64          `gorm:"column:dept_id"`
	}
	err := db.WithContext(ctx).
		Table("sys_user_role AS ur").
		Select("r.id AS role_id, r.data_scope, u.dept_id").
		Joins("JOIN sys_role AS r ON r.id = ur.role_id AND r.deleted_at IS NULL AND r.status = '1'").
		Joins("JOIN sys_user AS u ON u.id = ur.user_id AND u.deleted_at IS NULL").
		Where("ur.user_id = ?", userID).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	deptIDSet := make(map[uint64]struct{})
	hasOwnerOnly := false

	for _, row := range rows {
		switch row.DataScope {
		case model.DataScopeAll:
			out.IsSuper = true
			return out, nil // 全部权限直接短路
		case model.DataScopeSelf:
			hasOwnerOnly = true
		case model.DataScopeDept:
			if row.DeptID != nil {
				deptIDSet[*row.DeptID] = struct{}{}
			}
		case model.DataScopeDeptAndSub:
			if row.DeptID != nil {
				// 查本部门及所有子部门
				subIDs, e := getSubDeptIDs(ctx, db, *row.DeptID)
				if e != nil {
					return nil, e
				}
				for _, id := range subIDs {
					deptIDSet[id] = struct{}{}
				}
			}
		case model.DataScopeCustom:
			// 查 sys_role_dept 自定义部门
			var ids []uint64
			if err := db.WithContext(ctx).
				Model(&model.SysRoleDept{}).
				Where("role_id = ?", row.RoleID).
				Pluck("dept_id", &ids).Error; err != nil {
				return nil, err
			}
			for _, id := range ids {
				deptIDSet[id] = struct{}{}
			}
		}
	}

	// 多角色合并: 取并集
	for id := range deptIDSet {
		out.DeptIDs = append(out.DeptIDs, id)
	}

	// 只有所有角色都是 OwnerOnly 且没有任何部门范围时, 才走 OwnerOnly
	if hasOwnerOnly && len(out.DeptIDs) == 0 {
		out.OwnerOnly = true
	}

	return out, nil
}

// getSubDeptIDs 查本部门及所有子部门 id (含本身)
// 利用 sys_dept.ancestors 前缀索引
func getSubDeptIDs(ctx context.Context, db *gorm.DB, deptID uint64) ([]uint64, error) {
	ids := []uint64{deptID}
	var subs []uint64
	// ancestors 形如 "0,1,3", 子部门的 ancestors 包含 ",deptID," 或以 ",deptID" 结尾或等于 "deptID"
	// 用 LIKE 前缀方式查所有以本部门为祖先的部门
	err := db.WithContext(ctx).
		Model(&model.SysDept{}).
		Where("FIND_IN_SET(?, ancestors) > 0", deptID).
		Pluck("id", &subs).Error
	if err != nil {
		return nil, err
	}
	ids = append(ids, subs...)
	return ids, nil
}

// ApplyDataScope GORM Scope: 把数据权限条件追加到查询
// 用法: db.Scopes(scope.ApplyDataScope(scopeCtx)).Find(&books)
//
// 业务表必须有 owner_id 和 dept_id 两列, 否则不要套这个 Scope
func ApplyDataScope(s *DataScopeContext) func(*gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		if s == nil || s.IsSuper {
			return tx
		}
		if s.OwnerOnly {
			return tx.Where("owner_id = ?", s.UserID)
		}
		if len(s.DeptIDs) > 0 {
			// 部门范围 OR 自己创建的
			return tx.Where("dept_id IN ? OR owner_id = ?", s.DeptIDs, s.UserID)
		}
		// 既不是超管, 也没部门范围, 也不是 OwnerOnly => 默认只看自己
		return tx.Where("owner_id = ?", s.UserID)
	}
}
