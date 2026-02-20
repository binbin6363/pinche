<template>
  <div class="user-profile-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-[#EDEDED]'">
    <!-- 顶部导航（微信风格） -->
    <div class="sticky top-0 z-10 px-4 py-3 flex items-center" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-[#EDEDED]'">
      <button @click="goBack" class="w-10 h-10 flex items-center justify-center -ml-2">
        <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex items-center justify-center py-20">
      <div class="w-8 h-8 border-2 border-green-500 border-t-transparent rounded-full animate-spin"></div>
    </div>

    <!-- 用户资料 -->
    <div v-else-if="profile">
      <!-- 用户基本信息卡片（微信风格） -->
      <div class="px-4 py-5" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div class="flex items-start">
          <!-- 头像 -->
          <div class="w-16 h-16 rounded flex items-center justify-center overflow-hidden flex-shrink-0" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
            <img v-if="profile.avatar" :src="profile.avatar" class="w-full h-full object-cover" />
            <span v-else class="text-2xl font-medium" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'">
              {{ profile.nickname?.charAt(0) || '?' }}
            </span>
          </div>
          <!-- 用户信息 -->
          <div class="flex-1 ml-4">
            <div class="flex items-center gap-2">
              <span class="text-lg font-medium" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
                {{ profile.nickname }}
              </span>
              <span v-if="profile.gender === 1" class="w-5 h-5 rounded-full bg-blue-500 flex items-center justify-center">
                <svg class="w-3 h-3 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                </svg>
              </span>
              <span v-else-if="profile.gender === 2" class="w-5 h-5 rounded-full bg-pink-500 flex items-center justify-center">
                <svg class="w-3 h-3 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
                </svg>
              </span>
            </div>
            <!-- 好友状态标签 -->
            <div v-if="profile.friendship_status === 0" class="mt-1.5 text-sm text-gray-500">
              已发送好友申请
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮区（微信风格） -->
      <div class="mt-2" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <!-- 已是好友：显示发消息和更多按钮 -->
        <template v-if="profile.is_friend">
          <div 
            class="flex items-center px-4 py-3 cursor-pointer action-item"
            :class="appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100'"
            @click="startChat"
          >
            <div class="w-6 h-6 flex items-center justify-center mr-4">
              <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M20 2H4c-1.1 0-1.99.9-1.99 2L2 22l4-4h14c1.1 0 2-.9 2-2V4c0-1.1-.9-2-2-2zm-2 12H6v-2h12v2zm0-3H6V9h12v2zm0-3H6V6h12v2z"/>
              </svg>
            </div>
            <span class="flex-1" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">发消息</span>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div 
            class="flex items-center px-4 py-3 cursor-pointer action-item"
            @click="showMoreActions = true"
          >
            <div class="w-6 h-6 flex items-center justify-center mr-4">
              <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z" />
              </svg>
            </div>
            <span class="flex-1" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">更多</span>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </template>

        <!-- 非好友：显示添加好友按钮 -->
        <template v-else-if="profile.friendship_status !== 0">
          <div 
            class="flex items-center px-4 py-3 cursor-pointer action-item"
            @click="showAddFriendModal = true"
          >
            <div class="w-6 h-6 flex items-center justify-center mr-4">
              <svg class="w-5 h-5 text-green-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M15 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm-9-2V7H4v3H1v2h3v3h2v-3h3v-2H6zm9 4c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
              </svg>
            </div>
            <span class="flex-1" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">添加到通讯录</span>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </template>

        <!-- 等待对方同意 -->
        <template v-else>
          <div class="flex items-center px-4 py-3">
            <div class="w-6 h-6 flex items-center justify-center mr-4">
              <svg class="w-5 h-5 text-yellow-500" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
              </svg>
            </div>
            <span class="flex-1 text-gray-500">等待对方通过好友申请</span>
          </div>
        </template>
      </div>

      <!-- 车辆信息（好友可见） -->
      <div v-if="profile.is_friend && (profile.car_brand || profile.car_model || profile.car_color)" class="mt-2" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div class="px-4 py-3 text-sm text-gray-500" :class="appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100'">车辆信息</div>
        <div class="px-4 py-3 flex items-center justify-between" v-if="profile.car_brand">
          <span class="text-gray-500">品牌</span>
          <span :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ profile.car_brand }}</span>
        </div>
        <div class="px-4 py-3 flex items-center justify-between" :class="appStore.theme === 'dark' ? 'border-t border-gray-700' : 'border-t border-gray-100'" v-if="profile.car_model">
          <span class="text-gray-500">型号</span>
          <span :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ profile.car_model }}</span>
        </div>
        <div class="px-4 py-3 flex items-center justify-between" :class="appStore.theme === 'dark' ? 'border-t border-gray-700' : 'border-t border-gray-100'" v-if="profile.car_color">
          <span class="text-gray-500">颜色</span>
          <span :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ profile.car_color }}</span>
        </div>
      </div>

      <!-- 近期行程 -->
      <div class="mt-2" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div class="px-4 py-3 text-sm text-gray-500" :class="appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100'">近期行程</div>
        
        <div v-if="profile.recent_trips && profile.recent_trips.length > 0">
          <div 
            v-for="(trip, index) in profile.recent_trips" 
            :key="trip.id"
            @click="goToTrip(trip.id)"
            class="px-4 py-3 cursor-pointer action-item"
            :class="index < profile.recent_trips.length - 1 ? (appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-100') : ''"
          >
            <div class="flex items-center justify-between">
              <div class="flex items-center gap-2 flex-wrap">
                <span class="px-1.5 py-0.5 text-xs rounded" :class="trip.type === 1 ? 'bg-blue-100 text-blue-600' : 'bg-green-100 text-green-600'">
                  {{ trip.type === 1 ? '车主' : '乘客' }}
                </span>
                <span v-if="trip.status === 3" class="px-1.5 py-0.5 text-xs rounded bg-emerald-100 text-emerald-600">
                  已成行
                </span>
                <span v-else-if="trip.status === 2" class="px-1.5 py-0.5 text-xs rounded bg-gray-100 text-gray-500">
                  已取消
                </span>
                <span v-else class="px-1.5 py-0.5 text-xs rounded bg-orange-100 text-orange-600">
                  待出行
                </span>
              </div>
              <svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </div>
            <div class="mt-2 text-sm font-medium" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
              {{ trip.departure_city || trip.departure_address }} → {{ trip.destination_city || trip.destination_address }}
            </div>
            <div class="mt-1 text-xs text-gray-500">{{ formatTime(trip.departure_time) }}</div>
          </div>
        </div>

        <div v-else class="text-center py-8 text-gray-500 text-sm">
          暂无近期行程
        </div>
      </div>
    </div>

    <!-- 错误状态 -->
    <div v-else class="text-center py-20 text-gray-500">
      <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <p>用户不存在或加载失败</p>
    </div>

    <!-- 添加好友弹窗（微信风格） -->
    <div v-if="showAddFriendModal" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50" @click="showAddFriendModal = false"></div>
      <div class="relative w-72 rounded-lg overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div class="text-center py-4 font-medium" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
          申请添加好友
        </div>
        <div class="px-4 pb-4">
          <textarea 
            v-model="friendMessage"
            rows="3"
            class="w-full px-3 py-2 rounded-lg border resize-none text-sm"
            :class="appStore.theme === 'dark' ? 'bg-gray-700 border-gray-600 text-white' : 'bg-gray-50 border-gray-200 text-gray-900'"
            placeholder="向对方打个招呼吧..."
          ></textarea>
        </div>
        <div class="flex" :class="appStore.theme === 'dark' ? 'border-t border-gray-700' : 'border-t border-gray-200'">
          <button 
            @click="showAddFriendModal = false"
            class="flex-1 py-3 text-center"
            :class="appStore.theme === 'dark' ? 'text-gray-400 border-r border-gray-700' : 'text-gray-600 border-r border-gray-200'"
          >
            取消
          </button>
          <button 
            @click="sendRequest"
            :disabled="sending"
            class="flex-1 py-3 text-center text-green-500 font-medium disabled:opacity-50"
          >
            {{ sending ? '发送中...' : '发送' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 更多操作弹窗 -->
    <div v-if="showMoreActions" class="fixed inset-0 z-50 flex items-end justify-center">
      <div class="absolute inset-0 bg-black/50" @click="showMoreActions = false"></div>
      <div class="relative w-full" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div 
          class="py-4 text-center cursor-pointer text-red-500 font-medium"
          :class="appStore.theme === 'dark' ? 'border-b border-gray-700' : 'border-b border-gray-200'"
          @click="confirmDelete"
        >
          删除好友
        </div>
        <div 
          class="py-4 text-center cursor-pointer"
          :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'"
          @click="showMoreActions = false"
        >
          取消
        </div>
      </div>
    </div>

    <!-- 删除确认弹窗 -->
    <div v-if="showDeleteConfirm" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-black/50" @click="showDeleteConfirm = false"></div>
      <div class="relative w-72 rounded-lg overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <div class="text-center py-5 px-4">
          <p :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">
            将联系人"{{ profile?.nickname }}"删除，将同时删除与该联系人的聊天记录
          </p>
        </div>
        <div class="flex" :class="appStore.theme === 'dark' ? 'border-t border-gray-700' : 'border-t border-gray-200'">
          <button 
            @click="showDeleteConfirm = false"
            class="flex-1 py-3 text-center"
            :class="appStore.theme === 'dark' ? 'text-gray-400 border-r border-gray-700' : 'text-gray-600 border-r border-gray-200'"
          >
            取消
          </button>
          <button 
            @click="doDelete"
            :disabled="deleting"
            class="flex-1 py-3 text-center text-red-500 font-medium disabled:opacity-50"
          >
            {{ deleting ? '删除中...' : '删除' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAppStore } from '@/stores/app'
import { useFriendStore } from '@/stores/friend'

const route = useRoute()
const router = useRouter()
const appStore = useAppStore()
const friendStore = useFriendStore()

const loading = ref(true)
const profile = ref(null)
const showAddFriendModal = ref(false)
const friendMessage = ref('')
const sending = ref(false)
const showMoreActions = ref(false)
const showDeleteConfirm = ref(false)
const deleting = ref(false)

onMounted(async () => {
  await loadProfile()
})

async function loadProfile() {
  loading.value = true
  const userId = route.params.id
  const result = await friendStore.getUserProfile(userId)
  profile.value = result
  loading.value = false
}

function goBack() {
  router.back()
}

function goToTrip(tripId) {
  router.push(`/trip/${tripId}`)
}

function startChat() {
  if (profile.value?.open_id) {
    router.push({
      path: `/chat/${profile.value.open_id}`,
      query: {
        nickname: profile.value.nickname || '用户',
        avatar: profile.value.avatar || ''
      }
    })
  }
}

function formatTime(timeStr) {
  if (!timeStr) return ''
  const date = new Date(timeStr)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hours}:${minutes}`
}

async function sendRequest() {
  if (sending.value) return
  sending.value = true
  await friendStore.sendFriendRequest(profile.value.open_id, friendMessage.value)
  appStore.showToast('好友申请已发送', 'success')
  showAddFriendModal.value = false
  friendMessage.value = ''
  // refresh profile to update friendship status
  await loadProfile()
  sending.value = false
}

function confirmDelete() {
  showMoreActions.value = false
  showDeleteConfirm.value = true
}

async function doDelete() {
  if (deleting.value) return
  deleting.value = true
  await friendStore.deleteFriend(profile.value.open_id)
  appStore.showToast('已删除好友', 'success')
  showDeleteConfirm.value = false
  deleting.value = false
  router.back()
}
</script>

<style scoped>
.action-item:active {
  background-color: #ECECEC;
}

[data-theme="dark"] .action-item:active {
  background-color: #374151;
}
</style>
