# 重构 Admin Dashboard：小说管理系统看板

## 概述

将现有的通用管理后台仪表盘（Soybean Admin 模板）重构为适合**小说阅读管理系统**的看板，所有数据均使用前端硬编码假数据，无需修改后端。

## 当前状态分析

### 项目技术栈
- **框架**: Vue 3 + TypeScript + Vite 7
- **UI**: Naive UI + UnoCSS
- **图表**: ECharts 6（通过自定义 `useEcharts` hook 使用，位于 `src/hooks/common/echarts.ts`）
- **状态管理**: Pinia
- **国际化**: vue-i18n（中/英双语）
- **目录**: `src/views/admin/dashboard/`

### 当前仪表盘模块（6个）
| 文件 | 功能 | 处理方式 |
|------|------|----------|
| `index.vue` | 主布局容器 | **修改** - 重构布局 |
| `header-banner.vue` | 问候语 + 项目数/待办/消息统计 | **修改** - 适配小说系统 |
| `card-data.vue` | 4个指标卡片（访问量/成交额/下载量/成交量） | **修改** - 替换为小说指标 |
| `line-chart.vue` | 折线图（下载量/注册量趋势） | **修改** - 重用于入库趋势 |
| `pie-chart.vue` | 饼图（作息安排分布） | **修改** - 重用于分类分布 |
| `project-news.vue` | 项目动态列表 | **修改** - 重用于最近动态 |
| `creativity-banner.vue` | 创意装饰横幅 | **删除** - 不再需要 |

### 新增模块（3个）
| 文件 | 功能 |
|------|------|
| `bar-chart.vue` | 柱状图 - 近期入库小说数量（近14天） |
| `reading-time.vue` | 阅读时长统计 - Top 5 小说阅读时长柱状图 |
| `tag-distribution.vue` | 标签分布图表 |

## 详细修改方案

### 1. 删除 `creativity-banner.vue`
不再需要此装饰性模块。

### 2. 修改 `index.vue` — 主布局
**变更**：重构网格布局，移除 `CreativityBanner` 组件引用，引入3个新模块组件。

**新布局结构（从上到下）**：
```
┌─────────────────────────────────────────────┐
│ HeaderBanner（问候语 + 书库概览统计）          │
├─────────────────────────────────────────────┤
│ CardData（4个关键指标卡片）                    │
├───────────────────┬─────────────────────────┤
│ BarChart           │ LineChart               │
│ （柱状图-近14天入库） │ （折线图-入库趋势）        │
├───────────────────┼─────────────────────────┤
│ PieChart           │ ReadingTime             │
│ （饼图-分类分布）    │ （柱状图-Top5阅读时长）    │
├───────────────────┼─────────────────────────┤
│ ProjectNews        │ TagDistribution         │
│ （列表-最近动态）    │ （饼图/柱状图-标签分布）   │
└───────────────────┴─────────────────────────┘
```

### 3. 修改 `header-banner.vue` — 问候横幅
**变更**：
- 保留问候语（修改文字为小说场景）
- 将 `statisticData` 从 `[项目数, 待办, 消息]` 改为 `[总藏书数, 总分类数, 总标签数, 总角色数]`
- 修改头像展示

### 4. 修改 `card-data.vue` — 指标卡片
**变更**：
- 将4个卡片指标从 `[访问量, 成交额, 下载量, 成交量]` 改为：
  - 📚 **藏书总数** (totalBooks) - 如 12,586 本
  - 📖 **总章节数** (totalChapters) - 如 368,420 章
  - ⏱ **累计阅读** (totalReadingHours) - 如 2,580 小时
  - 📥 **本周新增** (weeklyNew) - 如 47 本

### 5. 修改 `line-chart.vue` — 折线图 → 入库趋势
**变更**：
- 修改 `legend` 数据从 `[下载量, 注册量]` 为 `[新增入库, 累计入库]`
- 修改 X 轴数据为最近7天日期标签
- 修改 `series` 数据为模拟的入库趋势数据
- 修改图表标题/国际化键

### 6. 修改 `pie-chart.vue` — 饼图 → 分类分布
**变更**：
- 修改 `series` 数据从 `[学习, 娱乐, 工作, 休息]` 为小说分类，如：
  - 玄幻 (Fantasy) - 35%
  - 言情 (Romance) - 25%
  - 科幻 (Sci-Fi) - 15%
  - 悬疑 (Mystery) - 15%
  - 历史 (History) - 10%
- 修改图表标题/国际化键

### 7. 修改 `project-news.vue` — 动态列表 → 最近动态
**变更**：
- 修改数据结构从通用项目动态改为小说系统的最近操作动态：
  - 新增小说（书名 + 时间）
  - 新增分类
  - 新增标签
  - 新增角色
  - 编辑小说
- 修改列表标题为"最近动态"

### 8. 新建 `bar-chart.vue` — 近期入库柱状图
**实现**：
- 使用 `useEcharts` hook 创建柱状图
- X轴：最近14天的日期
- Y轴：当天入库小说数量
- 使用 `type: 'bar'` 柱状图
- 模拟14天的入库数据（随机 3-15 本/天）
- 标题："近期入库统计"

### 9. 新建 `reading-time.vue` — 阅读时长统计
**实现**：
- 使用 `useEcharts` hook 创建柱状图
- X轴：Top 5 小说名称
- Y轴：累计阅读时长（分钟/小时）
- 模拟数据示例：
  - 《剑来》- 128 小时
  - 《诡秘之主》- 96 小时
  - 《凡人修仙传》- 85 小时
  - 《大奉打更人》- 72 小时
  - 《三体》- 60 小时
- 标题："阅读时长 TOP 5"

### 10. 新建 `tag-distribution.vue` — 标签分布
**实现**：
- 使用 `useEcharts` hook（饼图或横向柱状图）
- 展示 Top 6-8 标签的分布情况
- 模拟数据示例：
  - 穿越 - 120 本
  - 系统流 - 98 本
  - 后宫 - 85 本
  - 重生 - 76 本
  - 无敌流 - 65 本
  - 脑洞 - 55 本
  - 种田 - 42 本
  - 黑暗流 - 38 本
- 标题："标签分布"

## 国际化（i18n）更新

### 中英文 locale 文件修改
在 `page.home` 命名空间下新增/修改以下键：

**zh-cn.ts 修改**：
- 修改 `greeting` 文案为看板欢迎语
- 修改 `weatherDesc` 为看板副标题
- 新增键：
  - `totalBooks: "藏书总数"`
  - `totalChapters: "总章节数"`
  - `totalReadingHours: "累计阅读"`
  - `weeklyNew: "本周新增"`
  - `bookCategories: "分类"`
  - `bookTags: "标签"`
  - `bookCharacters: "角色"`
  - `recentImport: "近期入库统计"`
  - `importTrend: "入库趋势"`
  - `categoryDist: "分类分布"`
  - `readingTimeTop: "阅读时长 TOP 5"`
  - `tagDist: "标签分布"`
  - `recentActivity: "最近动态"`
  - 以及 `recentActivity.desc1` ~ `desc8` 动态描述消息
  - 动态标签：`activity.newBook`, `activity.editBook`, `activity.newCategory`, `activity.newTag`, `activity.newCharacter`

**en-us.ts 修改**：
- 同上英文对应翻译

## 数据方案

所有数据使用前端硬编码的 mock 数据，通过 `async function mockData()` + `setTimeout` 模拟延迟（复用现有模式），数据在组件内直接写死：

- 数字指标：直接定义在 `computed` 或 `ref` 中
- 图表数据：通过 `updateOptions` 在 `mockData()` 中注入
- 列表数据：直接定义数组常量

## 不变的部分

- 不需要修改路由（`/admin/dashboard` 保持不变）
- 不需要修改任何后端代码
- 不需要修改 `useEcharts` hook
- 不需要修改 store / Pinia 状态管理
- 所有 Naive UI 组件用法保持一致（NCard, NGrid, NGi, NSpace, NList, NListItem, NThing）
- 页面样式体系不变（UnoCSS + Naive UI 主题）

## 验证步骤
1. 启动开发服务器：`pnpm dev`
2. 访问 `/admin/dashboard`
3. 检查所有模块是否正常渲染
4. 检查图表是否正确加载和显示 mock 数据
5. 检查响应式布局是否正常（不同屏幕尺寸）