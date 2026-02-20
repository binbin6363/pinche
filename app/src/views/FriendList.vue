<template>
  <div class="friend-list-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-[#EDEDED]'">
    <!-- 顶部导航（微信风格） -->
    <div class="sticky top-0 z-10 px-4 py-3 flex items-center" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-[#EDEDED]'">
      <button @click="goBack" class="w-10 h-10 flex items-center justify-center -ml-2">
        <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
      <h1 class="flex-1 text-center text-lg font-medium -ml-10" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">通讯录</h1>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-green-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <template v-else>
      <!-- 好友数量统计 -->
      <div class="px-4 py-2 text-xs text-gray-500">
        {{ friendStore.friendCount }} 位好友
      </div>

      <!-- 好友列表 -->
      <div v-if="friendStore.friends.length > 0" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div 
          v-for="(item, index) in friendStore.friends" 
          :key="item.id"
          class="friend-item flex items-center px-4 py-3 cursor-pointer"
          :class="[
            index < friendStore.friends.length - 1 ? (appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100') : ''
          ]"
          @click="goToProfile(item.friend)"
        >
          <!-- 头像 -->
          <div class="w-11 h-11 rounded flex items-center justify-center overflow-hidden flex-shrink-0" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
            <img v-if="item.friend?.avatar" :src="item.friend.avatar" class="w-full h-full object-cover" />
            <span v-else class="text-base font-medium" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'">
              {{ item.friend?.nickname?.charAt(0) || '?' }}
            </span>
          </div>

          <!-- 昵称 -->
          <div class="flex-1 ml-3 min-w-0">
            <div class="font-normal text-base truncate" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
              {{ item.friend?.nickname }}
            </div>
          </div>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-else class="text-center py-20">
        <svg class="w-20 h-20 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <p class="text-gray-500 mb-2">还没有好友</p>
        <p class="text-sm text-gray-400">在行程详情页可以添加感兴趣的用户为好友</p>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useFriendStore } from '@/stores/friend'

const router = useRouter()
const appStore = useAppStore()
const friendStore = useFriendStore()

const loading = ref(true)

onMounted(async () => {
  await friendStore.fetchFriends()
  loading.value = false
})

function goBack() {
  router.back()
}

function goToProfile(friend) {
  if (friend?.open_id) {
    router.push(`/user/${friend.open_id}`)
  }
}
</script>

<style scoped>
.friend-item:active {
  background-color: #ECECEC;
}

[data-theme="dark"] .friend-item:active {
  background-color: #374151;
}
</style>
