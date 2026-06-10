# Admin Dashboard 看板内容规划

## 设计目标

多角度、多维度展示书库概况与个人阅读数据，帮助用户：
- **概览书库**：藏书规模、增长趋势、分类/标签分布
- **理解习惯**：何时读、读多久、读多快、阅读节奏
- **分析偏好**：喜欢什么类型/标签/作者/小说
- **跟踪动态**：最近做了什么、阅读趋势变化

## 当前看板分析

### 现有 10 个模块

| 模块 | 类型 | 维度 | 评价 |
|------|------|------|------|
| HeaderBanner | 问候+3个快读数 | 书库(分类/标签/角色) | ✅ 保留 |
| CardData | 5指标卡片 | 书库+阅读 | ✅ 保留，可扩展 |
| CalendarHeatmap | 30天热力图(HTML) | 阅读活跃度 | ✅ 保留 |
| RadarChart | 雷达图 | 6分类偏好 | ✅ 保留 |
| PieChart | 饼图 | 分类占比 | ✅ 保留 |
| ReadingHours | 24h柱状图 | 阅读时段 | ✅ 保留 |
| TagDistribution | 横向柱状图 | 标签数量 | ✅ 保留 |
| LineChart | 面积趋势线 | 7天阅读时长+次数 | ✅ 保留 |
| ProjectNews | 时间线列表 | 最近动态 | ✅ 保留 |

### 缺失维度（待补充）

| 缺失维度 | 重要性 | 说明 |
|----------|--------|------|
| 书库增长趋势 | ⭐⭐⭐ | 月度/周度新增藏书趋势 |
| 作者 TOP | ⭐⭐⭐ | 阅读最多/收藏最多的作者 |
| 小说 TOP | ⭐⭐⭐ | 阅读时长最多的书籍排名 |
| 阅读连胜 | ⭐⭐⭐ | 连续阅读天数，激励性指标 |
| 每周阅读分布 | ⭐⭐ | 周一~周日的阅读模式 |
| 阅读节奏 | ⭐⭐ | 平均每次阅读时长/每日阅读段数 |
| 连载/完结偏好 | ⭐⭐ | 正在阅读的小说完结状态 |
| 字数偏好 | ⭐⭐ | 短/中/长篇分布 |
| 本周 vs 上周对比 | ⭐⭐ | 阅读时长变化趋势 |
| 阅读目标 | ⭐⭐ | 今日/本周目标进度 |

## 规划布局（6行结构）

### 第1行：顶部横幅
- **HeaderBanner**（保留）
  - 问候语 + 头像
  - 书库快数：分类数 / 标签数 / 角色数

### 第2行：关键指标（扩展为6卡片）
- **CardData**（保留，扩展1项）
  - 📚 藏书总数
  - 📖 总章节数
  - ⏱ 累计阅读时长
  - ✨ 本周新增
  - 🎯 完成率
  - 🔥 **阅读连胜（新增）** — 如「连续阅读 12 天」

### 第3行：书库总览（2列, 16:8）
- **左列 (16/24)**：书库增长趋势（**新增模块 `library-trend.vue`**）
  - ECharts 面积图，双轴
  - X轴：最近12个月
  - 双线：月度新增藏书量（柱）+ 累计藏书总量（折线面积）
  - 直观看到书库成长轨迹
- **右列 (8/24)**：分类分布（保留 PieChart）
  - 标注饼图，展示主要分类占比

### 第4行：阅读习惯（3列, 8:8:8）
- **左列**：阅读时段分布（保留 ReadingHours）
  - 24小时阅读时长分布柱状图
- **中列**：每周阅读模式（**新增模块 `weekly-pattern.vue`**）
  - 周一~周日阅读时长柱状图
  - 顶部标注总时长和日均
- **右列**：阅读连胜日历（保留 CalendarHeatmap → 改名为简洁版）
  - 30天热力格 + 连胜天数标注

### 第5行：阅读偏好（3列, 8:8:8）
- **左列**：阅读偏好雷达（保留 RadarChart）
  - 6分类偏好蜘蛛网图
- **中列**：作者 TOP（**新增模块 `top-authors.vue`**）
  - TOP 6 作者阅读时长横向柱状图
- **右列**：标签分布（保留 TagDistribution）
  - TOP 8 标签横向柱状图

### 第6行：动态与趋势（2列, 14:10）
- **左列**：阅读趋势（保留 LineChart，改用 14天数据）
  - 双面积线：每日阅读时长 + 阅读次数
  - 14天更完整反映趋势
- **右列**：最近动态（保留 ProjectNews）
  - 时间线列表，带有色圆点标记类型

### 可选补充（若空间允许，可增加第7行 3列）
- 连载/完结分布（Pie）
- 字数偏好分布（Pie）
- 本周 vs 上周对比指标卡

## 文件变更清单

### 保留不变（5个）
- `header-banner.vue` — 内容不变
- `pie-chart.vue` — 内容不变
- `radar-chart.vue` — 内容不变
- `reading-hours.vue` — 内容不变
- `tag-distribution.vue` — 内容不变
- `project-news.vue` — 内容不变

### 需修改（3个）
- `card-data.vue` — 扩展为6项（增加「阅读连胜」）
- `line-chart.vue` — 扩大数据范围为14天
- `calendar-heatmap.vue` — 增加连胜标注（当前连续天数 + 最长连胜）
- `index.vue` — 重新布局为6行结构，引入新模块

### 新增（3个）
- `library-trend.vue` — 书库增长趋势（月度双线面积图）
- `weekly-pattern.vue` — 每周阅读模式（Mon-Sun柱状图）
- `top-authors.vue` — 作者 TOP 排名（横向柱状图）

### i18n 新增键
```
page.home.libraryTrend: "书库增长趋势" / "Library Growth"
page.home.monthlyNew: "月度新增" / "Monthly New"
page.home.cumulativeTotal: "累计总量" / "Cumulative Total"
page.home.weeklyPattern: "每周阅读模式" / "Weekly Pattern"
page.home.mon: "周一" ~ 周日
page.home.topAuthors: "作者 TOP" / "Top Authors"
page.home.readingStreak: "阅读连胜" / "Reading Streak"
page.home.currentStreak: "当前连胜" / "Current Streak"
page.home.longestStreak: "最长连胜" / "Longest Streak"
page.home.avgDaily: "日均阅读" / "Daily Avg"
page.home.totalSessions: "阅读次数" / "Reading Sessions"
```

## 技术方案

- 所有数据使用前端硬编码假数据，无需后端
- 新图表模块复用 `useEcharts` hook
- `library-trend.vue`：ECharts 双轴图（柱+面积线）
- `weekly-pattern.vue`：ECharts 柱状图
- `top-authors.vue`：ECharts 横向柱状图
- 布局使用 `NGrid` + `NGi`，响应式断点 `s:24 m:8/16/14/10`
- 移动端自动折叠为单列

## 数据维度总览

| 维度 | 模块 | 图表类型 | 视角 |
|------|------|----------|------|
| 藏书规模 | CardData | 数字卡片 | 当前总量 |
| 书库增长 | LibraryTrend | 柱+面积 | 时间趋势 |
| 分类占比 | PieChart | 饼图 | 结构分布 |
| 标签排名 | TagDistribution | 横向柱 | 排序对比 |
| 作者排名 | TopAuthors | 横向柱 | 排序对比 |
| 阅读总时长 | CardData | 数字卡片 | 当前总量 |
| 阅读时段 | ReadingHours | 柱状图 | 时间分布 |
| 阅读偏好 | RadarChart | 雷达图 | 多维度 |
| 每周模式 | WeeklyPattern | 柱状图 | 周期规律 |
| 阅读连胜 | CardData/Heatmap | 数字+热力 | 持续性 |
| 阅读趋势 | LineChart | 面积线 | 时间趋势 |
| 近期动态 | ProjectNews | 列表 | 事件流 |

共 **12 个维度**，覆盖 **书库、习惯、偏好、动态** 四大领域。