<template>
  <div class="main-layout h-full flex flex-col">
    <!-- 主内容区 -->
    <main class="flex-1 overflow-auto pb-16">
      <router-view />
    </main>

    <!-- 底部导航 - 4个Tab -->
    <nav class="tab-bar">
      <div class="flex items-end">
        <!-- 首页 -->
        <router-link
          to="/"
          class="tab-item flex-1"
          :class="{ active: isActive('/') }"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span>首页</span>
        </router-link>

        <!-- 行程广场 -->
        <router-link
          to="/trips"
          class="tab-item flex-1"
          :class="{ active: isActive('/trips') }"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
          </svg>
          <span>广场</span>
        </router-link>

        <!-- 消息 -->
        <router-link
          to="/matches"
          class="tab-item flex-1 relative"
          :class="{ active: isActive('/matches') }"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <span>消息</span>
          <!-- 未读数 -->
          <span
            v-if="messageStore.unreadCount > 0"
            class="absolute top-0 right-1/4 min-w-[18px] h-[18px] px-1 bg-red-500 text-white text-[11px] font-medium rounded-full flex items-center justify-center"
          >
            {{ messageStore.unreadCount > 99 ? '99+' : messageStore.unreadCount }}
          </span>
        </router-link>

        <!-- 我的 -->
        <router-link
          to="/profile"
          class="tab-item flex-1"
          :class="{ active: isActive('/profile') }"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <span>我的</span>
        </router-link>
      </div>
    </nav>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'

const route = useRoute()
const userStore = useUserStore()
const messageStore = useMessageStore()

onMounted(async () => {
  if (userStore.isLoggedIn) {
    await messageStore.fetchUnreadCount()
  }
})

function isActive(path) {
  if (path === '/') {
    return route.path === '/'
  }
  return route.path.startsWith(path)
}
</script>
