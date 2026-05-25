# Soybean Admin前端开发规范

## 一、核心信息

- **技术栈**: Vue 3 / TypeScript / Vite / Element Plus / Pinia / axios
- **项目路径**: `app/web`
- **AI角色**: 资深前端开发，严格遵循本规范生成代码

## 二、目录结构（核心）

```
src/
├── components/     # 全局组件 (advanced/common/custom)
├── hooks/          # 组合式函数 (business/common)
├── layouts/        # 布局组件
├── locales/        # 国际化
├── router/         # 路由 (routes/guard)
├── service/        # API服务 (api/request)
├── store/          # Pinia状态 (modules)
├── typings/        # TS类型声明
├── utils/          # 工具函数
└── views/          # 页面
```

## 三、命名规范

| 类型 | 规范 | 示例 |
|:---|:---|:---|
| 文件/文件夹 | kebab-case | `user-list.vue`, `home/` |
| Vue组件 | PascalCase | `<AppProvider />` |
| Icon组件 | kebab-case | `<icon-mdi-emoticon />` |
| 类/接口/类型 | PascalCase | `Person`, `UserProps` |
| 变量/函数 | camelCase | `userName`, `getUser()` |
| 常量 | UPPER_SNAKE | `MAX_COUNT` |
| CSS类名 | kebab-case | `.container-item` |
| 请求函数 | fetch + 资源名 | `fetchUser()`, `fetchUserList()` |

## 四、Vue组件规范

### SFC代码顺序

```vue
<script setup lang="ts">
// 1. import (顺序: vue → router → pinia → @vueuse → UI库 → 项目依赖 → 别名 → 相对路径)
import { ref } from 'vue'
import type { Ref } from 'vue'  // 类型单独导入

// 2. defineOptions

// 3. Props
interface Props { name: string }
const props = defineProps<Props>()

// 4. Emits
interface Emits { change: [value: string] }
const emit = defineEmits<Emits>()

// 5. Hooks
const router = useRouter()

// 6. 响应式状态
const count = ref(0)

// 7. 初始化函数
async function init() { }

// 8. watch
watch(count, (val) => { })

// 9. 生命周期
init()  // 等同于 created
onMounted(() => { })

// 10. defineExpose
defineExpose({ count })
</script>

<template>
  <!-- 模板 -->
</template>

<style scoped>
/* 样式 */
</style>
```

## 五、关键约束

1. **组件名**: Vue组件用PascalCase，Icon组件用kebab-case
2. **类型优先**: 所有Props/Emits必须定义TypeScript类型
3. **初始化逻辑**: 统一放在`init`函数中
4. **样式隔离**: 使用`scoped`避免样式污染
5. **请求函数**: 统一以`fetch`开头