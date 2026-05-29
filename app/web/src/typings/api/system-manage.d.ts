declare namespace Api {
  /**
   * namespace SystemManage
   *
   * backend api module: "systemManage"
   */
  namespace SystemManage {
    type CommonSearchParams = Pick<
      Common.PaginatingCommonParams,
      "current" | "size"
    >;

    /**
     * data scope
     *
     * - "1": all
     * - "2": custom department
     * - "3": current department
     * - "4": current department and sub-departments
     * - "5": self only
     */
    type DataScope = "1" | "2" | "3" | "4" | "5";

    /** role */
    type Role = Common.CommonRecord<{
      /** role name */
      roleName: string;
      /** role code */
      roleCode: string;
      /** role description */
      roleDesc: string | null;
      /** data scope */
      dataScope: DataScope;
    }>;

    /** role search params */
    type RoleSearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.Role, "roleName" | "roleCode" | "status"> &
        CommonSearchParams
    >;

    /** role list */
    type RoleList = Common.PaginatingQueryRecord<Role>;

    /** all role */
    type AllRole = Pick<Role, "id" | "roleName" | "roleCode">;

    /**
     * user gender
     *
     * - "1": "male"
     * - "2": "female"
     */
    type UserGender = "1" | "2";

    /** user */
    type User = Common.CommonRecord<{
      /** user name */
      userName: string;
      /** user password (only used in create/update request) */
      password: string;
      /** user gender */
      userGender: UserGender | null;
      /** user nick name */
      nickName: string;
      /** user phone */
      userPhone: string;
      /** user email */
      userEmail: string;
      /** user role code collection */
      userRoles: string[];
    }>;

    /** user search params */
    type UserSearchParams = CommonType.RecordNullable<
      Pick<
        Api.SystemManage.User,
        | "userName"
        | "userGender"
        | "nickName"
        | "userPhone"
        | "userEmail"
        | "status"
      > &
        CommonSearchParams
    >;

    /** user list */
    type UserList = Common.PaginatingQueryRecord<User>;

    /** department */
    type Dept = Common.CommonRecord<{
      /** parent dept id */
      parentId: number;
      /** dept name */
      deptName: string;
      /** dept code */
      deptCode: string;
      /** dept leader */
      leader: string | null;
      /** sort order */
      sortOrder: number;
      /** children dept */
      children?: Dept[] | null;
    }>;

    /** dept search params */
    type DeptSearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.Dept, "deptName" | "status"> & CommonSearchParams
    >;

    /** dept list */
    type DeptList = Common.PaginatingQueryRecord<Dept>;

    /**
     * menu type
     *
     * - "1": directory
     * - "2": menu
     */
    type MenuType = "1" | "2";

    type MenuButton = {
      /**
       * button code
       *
       * it can be used to control the button permission
       */
      buttonCode: string;
      /** button description */
      buttonDesc: string;
    };

    /**
     * icon type
     *
     * - "1": iconify icon
     * - "2": local icon
     */
    type IconType = "1" | "2";

    type MenuPropsOfRoute = Pick<
      import("vue-router").RouteMeta,
      | "i18nKey"
      | "keepAlive"
      | "constant"
      | "order"
      | "href"
      | "hideInMenu"
      | "activeMenu"
      | "multiTab"
      | "fixedIndexInTab"
      | "query"
    >;

    type Menu = Common.CommonRecord<{
      /** menu id */
      id: number;
      /** parent menu id */
      parentId: number;
      /** menu type */
      menuType: MenuType;
      /** menu name */
      menuName: string;
      /** route name */
      routeName: string;
      /** route path */
      routePath: string;
      /** component */
      component?: string;
      /** iconify icon name or local icon name */
      icon: string;
      /** icon type */
      iconType: IconType;
      /** buttons */
      buttons?: MenuButton[] | null;
      /** children menu */
      children?: Menu[] | null;
    }> &
      MenuPropsOfRoute;

    /** menu list */
    type MenuList = Common.PaginatingQueryRecord<Menu>;

    /** menu search params */
    type MenuSearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.Menu, "menuName" | "status"> &
        CommonSearchParams
    >;

    type MenuTree = {
      id: number;
      label: string;
      pId: number;
      children?: MenuTree[];
    };

    /** menu tree node (with buttons, from /manage/menu/tree) */
    type MenuTreeNode = {
      id: number;
      parentId: number;
      menuType: string;
      menuName: string;
      routeName: string;
      status: string;
      children?: MenuTreeNode[] | null;
      buttons?: SysMenuButton[] | null;
    };

    /** menu button */
    type SysMenuButton = {
      id: number;
      menuId: number;
      buttonCode: string;
      buttonDesc: string | null;
    };

    /** dict */
    type Dict = Common.CommonRecord<{
      dictName: string;
      dictCode: string;
      dictDesc: string | null;
    }>;

    /** dict item */
    type DictItem = Common.CommonRecord<{
      dictId: number;
      itemLabel: string;
      itemValue: string;
      itemDesc: string;
      sortOrder: number;
    }>;

    /** dict search params */
    type DictSearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.Dict, "dictName" | "dictCode" | "status"> &
        CommonSearchParams
    >;

    /** login log search params */
    type LoginLogSearchParams = CommonType.RecordNullable<
      {
        userName: string;
        loginIp: string;
        loginType: string;
        loginResult: string;
      } &
        CommonSearchParams & {
          startTime?: string;
          endTime?: string;
        }
    >;

    /** operation log search params */
    type OperationLogSearchParams = CommonType.RecordNullable<
      {
        userName: string;
        module: string;
        action: string;
        clientIp: string;
      } &
        CommonSearchParams & {
          startTime?: string;
          endTime?: string;
        }
    >;

    /** book category */
    type BookCategory = Common.CommonRecord<{
      parentId: number;
      categoryName: string;
      categoryCode: string;
      description: string | null;
      sortOrder: number;
      children?: BookCategory[] | null;
    }>;

    /** category search params */
    type CategorySearchParams = CommonType.RecordNullable<
      Pick<Api.SystemManage.BookCategory, "categoryName" | "categoryCode" | "status"> & CommonSearchParams
    >;

    /** book tag */
    type BookTag = {
      id: number;
      tagName: string;
      description: string | null;
      usageCount: number;
    };

    /** tag search params */
    type TagSearchParams = CommonType.RecordNullable<
      { tagName: string } & CommonSearchParams
    >;

    // ============= Book (Novel) =============

    /** serial status: 1-ongoing 2-finished 3-stopped */
    type SerialStatus = '1' | '2' | '3';
    /** visibility: 1-public 2-private 3-dept */
    type Visibility = '1' | '2' | '3';
    /** book listing status: 1-listed 2-unlisted 3-reviewing 4-rejected */
    type BookListingStatus = '1' | '2' | '3' | '4';
    /** aggregate status: 1-single 2-merging 3-done */
    type AggregateStatus = '1' | '2' | '3';

    /** book (novel) */
    type Book = Common.CommonRecord<{
      title: string;
      author: string;
      cover: string | null;
      intro: string | null;
      categoryId: number | null;
      language: string;
      serialStatus: SerialStatus;
      visibility: Visibility;
      primaryFileId: number | null;
      totalChapters: number;
      totalWords: number;
      aggregateStatus: AggregateStatus;
      avgRating: number;
      ratingCount: number;
      ownerId: number;
      ownerType: string;
      deptId: number | null;
      tagIds?: number[];
      tags?: { id: number; tagName: string }[];
    }>;

    /** book search params */
    type BookSearchParams = CommonType.RecordNullable<
      {
        title: string;
        author: string;
        categoryId: number | null;
        status: BookListingStatus;
        visibility: Visibility;
        serialStatus: SerialStatus;
        tagId: number | null;
      } & CommonSearchParams
    >;

    /** book create/update request */
    type BookRequest = {
      title: string;
      author: string;
      cover: string | null;
      intro: string | null;
      categoryId: number | null;
      language: string;
      serialStatus: SerialStatus;
      visibility: Visibility;
      tagIds: number[];
    };

    /** book status update request */
    type BookUpdateStatusRequest = {
      status: BookListingStatus;
    };
  }
}
