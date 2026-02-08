# 前端代码模板

## 技术栈

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue | 3.4.15 | 前端框架 |
| Vite | 5.0.12 | 构建工具 |
| Pinia | 2.1.7 | 状态管理 |
| vue-router | 4.2.5 | 路由管理 |
| TailwindCSS | 3.4.1 | CSS 框架 |
| Axios | 1.6.5 | HTTP 客户端 |

## Vue 组件模板

### 完整页面组件

```vue
<template>
  <div class="min-h-screen" :class="bgClass">
    <!-- 头部 -->
    <header class="sticky top-0 z-10 backdrop-blur-lg border-b"
            :class="appStore.theme === 'dark' ? 'bg-gray-900/80 border-gray-800' : 'bg-white/80 border-gray-200'">
      <div class="flex items-center justify-between px-4 py-3">
        <h1 class="text-lg font-semibold">页面标题</h1>
      </div>
    </header>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-blue-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="list.length === 0" class="flex flex-col items-center justify-center py-20">
      <span class="text-4xl mb-4">📭</span>
      <p class="text-gray-500">暂无数据</p>
    </div>

    <!-- 列表内容 -->
    <div v-else class="px-4 py-4 space-y-3">
      <div v-for="item in list" :key="item.id" 
           class="card p-4" @click="viewDetail(item)">
        <h3 class="font-medium">{{ item.title }}</h3>
        <p class="text-sm text-gray-500 mt-1">{{ item.description }}</p>
      </div>
    </div>

    <!-- 加载更多 -->
    <div v-if="hasMore" class="px-4 py-4">
      <button @click="loadMore" 
              class="w-full py-3 text-center text-blue-500"
              :disabled="loadingMore">
        {{ loadingMore ? '加载中...' : '加载更多' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useUserStore } from '@/stores/user'
import api from '@/utils/api'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()

// 响应式数据
const list = ref([])
const loading = ref(true)
const loadingMore = ref(false)
const page = ref(1)
const total = ref(0)

// 计算属性
const bgClass = computed(() => {
  return appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'
})

const hasMore = computed(() => {
  return list.value.length < total.value
})

// 生命周期
onMounted(async () => {
  await fetchList()
})

onUnmounted(() => {
  // 清理定时器、监听器等
})

// 方法
async function fetchList() {
  try {
    const result = await api.get('/examples', {
      params: { page: page.value, page_size: 20 }
    })
    list.value = result.list || []
    total.value = result.total || 0
  } finally {
    loading.value = false
  }
}

async function loadMore() {
  if (loadingMore.value) return
  loadingMore.value = true
  
  try {
    page.value++
    const result = await api.get('/examples', {
      params: { page: page.value, page_size: 20 }
    })
    list.value.push(...(result.list || []))
  } finally {
    loadingMore.value = false
  }
}

function viewDetail(item) {
  router.push(`/example/${item.id}`)
}
</script>

<style scoped>
.card {
  @apply rounded-xl transition-all duration-200;
  background: var(--card-bg);
  border: 1px solid var(--card-border);
}

.card:active {
  transform: scale(0.98);
}
</style>
```

## Pinia Store 模板

```javascript
// stores/example.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/utils/api'

export const useExampleStore = defineStore('example', () => {
  // State
  const list = ref([])
  const current = ref(null)
  const loading = ref(false)
  const page = ref(1)
  const total = ref(0)

  // Getters
  const isEmpty = computed(() => list.value.length === 0)
  const hasMore = computed(() => list.value.length < total.value)

  // Actions
  async function fetchList(params = {}) {
    loading.value = true
    try {
      const result = await api.get('/examples', { params })
      list.value = result.list || []
      total.value = result.total || 0
      page.value = params.page || 1
      return result
    } finally {
      loading.value = false
    }
  }

  async function fetchById(id) {
    const result = await api.get(`/examples/${id}`)
    current.value = result
    return result
  }

  async function create(data) {
    const result = await api.post('/examples', data)
    list.value.unshift(result)
    return result
  }

  async function update(id, data) {
    const result = await api.put(`/examples/${id}`, data)
    const index = list.value.findIndex(item => item.id === id)
    if (index > -1) {
      list.value[index] = result
    }
    if (current.value?.id === id) {
      current.value = result
    }
    return result
  }

  async function remove(id) {
    await api.delete(`/examples/${id}`)
    const index = list.value.findIndex(item => item.id === id)
    if (index > -1) {
      list.value.splice(index, 1)
    }
  }

  function reset() {
    list.value = []
    current.value = null
    page.value = 1
    total.value = 0
  }

  return {
    // State
    list,
    current,
    loading,
    page,
    total,
    // Getters
    isEmpty,
    hasMore,
    // Actions
    fetchList,
    fetchById,
    create,
    update,
    remove,
    reset
  }
})
```

## API 调用封装

```javascript
// utils/api.js 已封装，直接使用
import api from '@/utils/api'

// GET 请求
const result = await api.get('/trips', { params: { page: 1 } })

// POST 请求
const data = await api.post('/trips', { title: 'xxx' })

// PUT 请求
const updated = await api.put('/trips/123', { title: 'new' })

// DELETE 请求
await api.delete('/trips/123')
```

## 路由配置

```javascript
// router/index.js
import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    children: [
      { path: '', name: 'Home', component: () => import('@/views/Home.vue') },
      { path: 'list', name: 'List', component: () => import('@/views/List.vue') },
      { path: 'detail/:id', name: 'Detail', component: () => import('@/views/Detail.vue') },
    ]
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
```

## 主题系统使用

```vue
<template>
  <!-- 根据主题动态样式 -->
  <div :class="appStore.theme === 'dark' ? 'bg-gray-900 text-white' : 'bg-white text-gray-900'">
    
    <!-- 使用 CSS 变量 -->
    <div :style="{ color: 'var(--theme-primary)' }">主题色文字</div>
    
    <!-- 春节主题特殊处理 -->
    <div v-if="appStore.theme === 'spring'" class="bg-red-500">
      春节特效
    </div>
  </div>
</template>

<script setup>
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 切换主题
function toggleTheme() {
  const themes = ['light', 'dark', 'spring']
  const currentIndex = themes.indexOf(appStore.theme)
  const nextIndex = (currentIndex + 1) % themes.length
  appStore.setTheme(themes[nextIndex])
}
</script>
```

## 样式规范

### TailwindCSS 常用类

```html
<!-- 布局 -->
<div class="flex items-center justify-between">
<div class="grid grid-cols-2 gap-4">

<!-- 间距 -->
<div class="px-4 py-3 mb-4 space-y-3">

<!-- 圆角 -->
<div class="rounded-lg rounded-xl rounded-full">

<!-- 阴影 -->
<div class="shadow-sm shadow-md shadow-lg">

<!-- 动画 -->
<div class="transition-all duration-200 animate-spin">
```

### 触摸目标规范 (iOS)

```css
/* 最小触摸目标 44px */
.btn {
  min-height: 44px;
  min-width: 44px;
}

/* 点击反馈 */
.card:active {
  transform: scale(0.98);
}
```

### Safe Area 适配

```css
/* 底部安全区 */
.bottom-nav {
  padding-bottom: env(safe-area-inset-bottom);
}

/* 顶部安全区 */
.header {
  padding-top: env(safe-area-inset-top);
}
```
