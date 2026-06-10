# Admin Dashboard 全面规划

## 当前状态分析

### 已有 12 个模块

| 模块 | 内容 | 形式 |
|------|------|------|
| header-banner.vue | 问候语 + 分类/标签/角色快数 | 巨大横条 |
| card-data.vue | 6 项指标（藏书/章节/阅读/新增/完成率/连胜） | 数字卡片 |
| library-trend.vue | 12 月书库增长（月度新增柱+累计面积线） | ECharts 双轴 |
| pie-chart.vue | 5 分类占比 | ECharts 饼图 |
| reading-hours.vue | 24 小时阅读时段 | ECharts 柱图 |
| weekly-pattern.vue | 周一~周日阅读模式 | ECharts 柱图 |
| calendar-heatmap.vue | 30 天阅读热力 + 连胜标注 | HTML 网格 |
| radar-chart.vue | 6 分类阅读偏好 | ECharts 雷达图 |
| top-authors.vue | TOP 6 作者阅读时长 | ECharts 横向柱图 |
| tag-distribution.vue | TOP 8 标签数量 | ECharts 横向柱图 |
| line-chart.vue | 14 天阅读时长+次数趋势 | ECharts 面积线 |
| project-news.vue | 书库操作动态时间线 | 列表 |

### 路由
- 看板: `/admin/dashboard`
- 小说列表: `/admin/library/book`
- 分类管理: `/admin/library/book-category`
- 标签管理: `/admin/library/book-tag`
- 小说编辑模态框: 内嵌在 `book/index.vue` 的 `BookOperateModal`

### 技术栈
- Naive UI (NCard, NGrid, NGi, NTag, NButton, NImage, NList, etc.)
- ECharts 6 (useEcharts hook)
- Vue 3 + TypeScript + UnoCSS
- i18n: vue-i18n

---

## 规划方案：四大区域 + 10 行布局

### 区域一：书库概览（Library Overview）— 3 行

#### 第 1 行：顶部横幅
- **HeaderBanner**（保留现有，加强统计）
  - 问候语 + 头像
  - 快数：分类数 / 标签数 / 角色数 / 作者数

#### 第 2 行：关键指标卡片（扩展为 8 卡片）
- **CardData**（修改，8 项指标，2 行 × 4 列）
  - 📚 藏书总数 | 📖 总章节数 | ⏱ 累计阅读 | ✨ 本周新增
  - 🎯 完成率 | 🔥 阅读连胜 | 👤 作者总数 | 📝 角色总数

#### 第 3 行：书库增长 + 分类分布（2 列, 16:8）
- **LibraryTrend**（保留）— 12 月双轴柱+面积
- **PieChart**（保留）— 5 分类占比

---

### 区域二：书库内容（Library Content）— 3 行

#### 第 4 行：近期入库小说列表 + 字数分布（2 列, 14:10）
- **新建 `recent-imports.vue`** — 近期入库小说列表
  - 模拟 6~8 条数据
  - 列：封面（NImage 缩略图）、书名、作者、分类、入库时间
  - 封面用 CSS 占位色块模拟（无需真实图片）
- **新建 `word-count-dist.vue`** — 字数分布统计
  - ECharts 饼图 / 环形图
  - 短篇(<10w) / 中篇(10w~50w) / 长篇(50w~100w) / 超长篇(>100w)

#### 第 5 行：最新编辑小说 + 收藏评分分布（2 列, 14:10）
- **新建 `recently-edited.vue`** — 最新完善/编辑的小说
  - 模拟 6 条数据
  - 封面 + 书名 + 作者 + 最后编辑时间 + 编辑类型（新增/修改章节/修改信息）
  - 封面缩略图 + 书名链接到 `/admin/library/book`
- **新建 `rating-dist.vue`** — 收藏/评分分布
  - 双图表：
    - 收藏数分布（横向柱状图，TOP 5）
    - 评分分布（饼图，1~5 星占比）

#### 第 6 行：待完善小说 + 随机推荐（2 列, 14:10）
- **新建 `pending-books.vue`** — 待完善小说列表
  - 模拟 6 条数据
  - 封面 + 书名 + 缺少信息标注（如"缺简介"、"缺分类"、"缺标签"）
  - 每行右侧有「去完善」按钮，点击跳转到 `/admin/library/book`
  - 使用 NButton tertiary + router.push
- **新建 `random-recommend.vue`** — 随机推荐
  - 模拟 6 条数据
  - 封面 + 书名 + 作者 + 分类 + 推荐理由（如"同类型阅读较多"）
  - 刷新按钮换一批

---

### 区域三：阅读习惯（Reading Habits）— 2 行

#### 第 7 行：时段 + 每周 + 热力（3 列, 8:8:8）
- **ReadingHours**（保留）— 24 小时时段
- **WeeklyPattern**（保留）— 周一~周日
- **CalendarHeatmap**（保留）— 30 天热力 + 连胜

#### 第 8 行：阅读趋势 + 分类阅读时长 + 标签阅读时长（3 列, 8:8:8）
- **LineChart**（保留，改标题为"每日阅读趋势"）— 14 天双面积线
- **新建 `reading-by-category.vue`** — 各分类阅读时长
  - ECharts 横向柱状图
  - X 轴：阅读时长(小时)，Y 轴：分类名
- **新建 `reading-by-tag.vue`** — 各标签阅读时长
  - ECharts 横向柱状图
  - X 轴：阅读时长(小时)，Y 轴：标签名

---

### 区域四：阅读偏好 + 动态（Reading Preferences & Activity）— 2 行

#### 第 9 行：雷达 + 作者 + 标签（3 列, 8:8:8）
- **RadarChart**（保留）— 6 分类偏好雷达
- **TopAuthors**（保留）— 作者 TOP 6
- **TagDistribution**（保留）— 标签数量 TOP 8

#### 第 10 行：书库动态（1 列全宽）
- **ProjectNews**（保留，扩展为完整动态）
  - 增加删除操作（删除小说、删除分类、删除标签）
  - 增加作者/字数信息变更
  - 每种操作带不同颜色标记

---

## 文件变更清单

### 保留不变（0 个）
所有模块都需修改或替换，无完全不变的。

### 保留需修改（8 个）
| 文件 | 修改内容 |
|------|----------|
| `card-data.vue` | 指标扩展为 8 项（grid 改为 4×2），新增作者数/角色数 |
| `header-banner.vue` | 快数增加"作者数"（4 项） |
| `library-trend.vue` | 不变 |
| `pie-chart.vue` | 不变 |
| `reading-hours.vue` | 不变 |
| `weekly-pattern.vue` | 不变 |
| `radar-chart.vue` | 不变 |
| `top-authors.vue` | 不变 |
| `tag-distribution.vue` | 不变 |
| `line-chart.vue` | 标题改为"每日阅读趋势" |
| `calendar-heatmap.vue` | 不变 |
| `project-news.vue` | 增加删除/修改的动态类型，丰富内容 |

### 新增（6 个）
| 文件 | 功能 | 形式 |
|------|------|------|
| `recent-imports.vue` | 近期入库小说列表（封面+书名+作者+分类+时间） | NImage + 列表 |
| `recently-edited.vue` | 最新完善/编辑的小说（封面+信息+编辑时间） | 列表 |
| `pending-books.vue` | 待完善小说列表（封面+缺失标注+去完善按钮） | 列表 + NButton |
| `random-recommend.vue` | 随机推荐小说（封面+推荐理由+刷新） | 列表 + 刷新按钮 |
| `word-count-dist.vue` | 字数分布统计 | ECharts 饼图 |
| `rating-dist.vue` | 收藏/评分分布 | ECharts 横向柱图 + 饼图 |
| `reading-by-category.vue` | 各分类阅读时长 | ECharts 横向柱图 |
| `reading-by-tag.vue` | 各标签阅读时长 | ECharts 横向柱图 |

### 删除（0 个）
所有现有模块都保留，只新增。

### i18n 新增键（~20 个）
```
page.home.recentImports: "近期入库"
page.home.authorCount: "作者"
page.home.wordCount: "字数"
page.home.wordCountDist: "字数分布"
page.home.short: "短篇"
page.home.medium: "中篇"
page.home.long: "长篇"
page.home.extraLong: "超长篇"
page.home.recentlyEdited: "最新编辑"
page.home.editType: "编辑类型"
page.home.pendingBooks: "待完善书籍"
page.home.goEdit: "去完善"
page.home.missingIntro: "缺简介"
page.home.missingCategory: "缺分类"
page.home.missingTags: "缺标签"
page.home.randomRecommend: "随机推荐"
page.home.refresh: "换一批"
page.home.recommendReason: "推荐理由"
page.home.ratingDist: "收藏与评分"
page.home.rating: "评分"
page.home.favorites: "收藏"
page.home.readingByCategory: "分类阅读时长"
page.home.readingByTag: "标签阅读时长"
page.home.delete: "删除"
```

### Schema 类型更新
在 `app.d.ts` `page.home` 下新增对应的 string 类型键。

---

## 总布局结构

```
Row 1  ─── HeaderBanner（问候 + 分类/标签/角色/作者 快数）
Row 2  ─── CardData（8 指标卡片，2行×4列）
Row 3  ─── LibraryTrend（书库增长） | PieChart（分类分布）
────── 书库概览 完 ──────

Row 4  ─── RecentImports（入库列表+封面） | WordCountDist（字数饼图）
Row 5  ─── RecentlyEdited（编辑列表+封面） | RatingDist（收藏评分）
Row 6  ─── PendingBooks（待完善+去完善按钮） | RandomRecommend（随机推荐）
────── 书库内容 完 ──────

Row 7  ─── ReadingHours | WeeklyPattern | CalendarHeatmap（习惯三连）
Row 8  ─── LineChart（趋势） | ReadingByCategory | ReadingByTag
────── 阅读习惯 完 ──────

Row 9  ─── RadarChart | TopAuthors | TagDistribution（偏好三连）
Row 10 ─── ProjectNews（书库动态，全宽）
────── 偏好+动态 完 ──────
```

## 验证步骤
1. `pnpm typecheck` 类型检查通过
2. 启动 `pnpm dev` 访问 `/admin/dashboard`
3. 检查 10 行布局是否正确渲染
4. 检查所有 ECharts 图表初始化正常
5. 检查封面占位符显示正常
6. 检查「去完善」按钮点击跳转到 `/admin/library/book`
7. 检查「换一批」按钮刷新推荐列表
8. 检查移动端响应式布局（720px / 375px 断点）