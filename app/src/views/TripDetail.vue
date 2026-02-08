<template>
  <div class="trip-detail-page min-h-screen flex flex-col" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 头图区域 -->
    <div class="relative bg-gray-200" style="height: 45vh;">
      <!-- 返回按钮 -->
      <button 
        @click="goBack" 
        class="absolute top-4 left-4 z-10 w-8 h-8 bg-white/80 backdrop-blur rounded-full shadow flex items-center justify-center"
      >
        <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- 图片展示区域 -->
      <div class="w-full h-full">
        <!-- 有上传图片时展示图片轮播 -->
        <div v-if="tripImages.length > 0" class="w-full h-full relative">
          <img 
            :src="tripImages[currentImageIndex]" 
            class="w-full h-full object-cover"
            @error="handleImageError"
          />
          <!-- 图片指示器 -->
          <div v-if="tripImages.length > 1" class="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-1.5">
            <span 
              v-for="(_, index) in tripImages" 
              :key="index"
              @click="currentImageIndex = index"
              class="w-2 h-2 rounded-full cursor-pointer transition-all"
              :class="index === currentImageIndex ? 'bg-white w-4' : 'bg-white/50'"
            ></span>
          </div>
          <!-- 左右切换按钮 -->
          <button 
            v-if="tripImages.length > 1 && currentImageIndex > 0"
            @click="currentImageIndex--"
            class="absolute left-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/30 rounded-full flex items-center justify-center"
          >
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <button 
            v-if="tripImages.length > 1 && currentImageIndex < tripImages.length - 1"
            @click="currentImageIndex++"
            class="absolute right-2 top-1/2 -translate-y-1/2 w-8 h-8 bg-black/30 rounded-full flex items-center justify-center"
          >
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>
        <!-- 无图片时展示默认背景 -->
        <div v-else class="w-full h-full bg-gradient-to-br from-blue-400 via-blue-500 to-indigo-600 flex items-center justify-center">
          <div class="text-center text-white">
            <svg class="w-16 h-16 mx-auto mb-3 opacity-80" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7" />
            </svg>
            <p class="text-sm opacity-80">行程详情</p>
          </div>
        </div>
      </div>

    </div>

    <!-- 加载中 -->
    <div v-if="loading" class="flex-1 flex justify-center items-center py-12">
      <div class="loading-spinner"></div>
    </div>

    <!-- 行程不存在 -->
    <div v-else-if="!trip" class="flex-1 flex justify-center items-center text-gray-400">
      行程不存在
    </div>

    <!-- 行程详情 -->
    <div v-else class="flex-1 flex flex-col">
      <!-- 信息卡片 -->
      <div class="flex-1 -mt-4 rounded-t-2xl relative z-10" :class="appStore.theme === 'dark' ? 'bg-gray-800' : 'bg-white'">
        <!-- 提示条 -->
        <div v-if="trip.remark" class="px-4 py-2 bg-blue-50 flex items-center gap-2 rounded-t-2xl">
          <svg class="w-4 h-4 text-blue-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-xs text-blue-700 truncate">{{ trip.remark }}</p>
        </div>

        <!-- 发布者信息 -->
        <div class="px-4 py-3 flex items-center justify-between border-b border-gray-100">
          <div class="flex items-center gap-3" @click="showUserProfile(trip.user)">
            <div class="w-10 h-10 bg-gray-200 rounded-full flex items-center justify-center overflow-hidden cursor-pointer">
              <img v-if="trip.user?.avatar" :src="trip.user.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-sm font-medium text-gray-500">{{ getUserInitial(trip) }}</span>
            </div>
            <div>
              <div class="text-sm font-medium text-gray-800">{{ getMaskedNickname(trip) }}</div>
            </div>
          </div>
          <!-- 聊天按钮 -->
          <button
            v-if="canChat"
            @click="goChat"
            class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center hover:bg-gray-200 transition-colors"
          >
            <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
            </svg>
          </button>
        </div>

        <!-- 行程信息 -->
        <div class="px-4 py-3">
          <!-- 时间和标签 -->
          <div class="flex items-center gap-2 mb-3 text-sm text-gray-600">
            <span>{{ formatDetailTime(trip.departure_time) }}</span>
            <span class="px-1.5 py-0.5 text-xs border border-gray-200 rounded">{{ trip.seats }}人</span>
            <span 
              class="px-1.5 py-0.5 text-xs rounded"
              :class="trip.trip_type === 1 ? 'bg-primary-100 text-primary-600' : 'bg-green-100 text-green-600'"
            >
              {{ trip.trip_type === 1 ? '车找人' : '人找车' }}
            </span>
          </div>

          <!-- 起点 -->
          <div class="flex items-start gap-2 mb-2">
            <span class="w-2 h-2 mt-1.5 rounded-full flex-shrink-0" style="background: var(--theme-primary);"></span>
            <div class="flex-1">
              <div class="text-sm text-gray-800">{{ trip.departure_address || trip.departure_city }}</div>
            </div>
          </div>
          <!-- 终点 -->
          <div class="flex items-start gap-2">
            <span class="w-2 h-2 mt-1.5 rounded-full bg-orange-500 flex-shrink-0"></span>
            <div class="flex-1">
              <div class="text-sm text-gray-800">{{ trip.destination_address || trip.destination_city }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 底部操作栏 -->
      <div class="border-t px-4 py-3 safe-area-bottom" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
        <div class="flex items-center gap-3">
          <!-- 价格 -->
          <div class="flex-shrink-0">
            <div v-if="trip.price > 0" class="text-xl font-bold text-gray-900">
              {{ trip.price }}<span class="text-sm font-normal text-gray-500">元</span>
            </div>
            <div v-else class="text-sm text-gray-500">价格面议</div>
          </div>
          <!-- 抢单按钮 -->
          <button
            v-if="canGrab"
            @click="handleGrab"
            :disabled="grabbing"
            class="flex-1 py-3 bg-gradient-to-r from-orange-500 to-orange-400 text-white text-base font-semibold rounded-full hover:from-orange-600 hover:to-orange-500 active:scale-98 transition-all disabled:opacity-50"
          >
            {{ grabbing ? '处理中...' : grabButtonText }}
          </button>
          <!-- 已是自己的行程 -->
          <div v-else-if="isOwner" class="flex-1 py-3 bg-gray-100 text-gray-500 text-base text-center rounded-full">
            这是您发布的行程
          </div>
          <!-- 未登录 -->
          <button
            v-else
            @click="goLogin"
            class="flex-1 py-3 bg-primary-500 text-white text-base font-semibold rounded-full"
          >
            登录后联系
          </button>
        </div>
      </div>
    </div>

    <!-- 用户资料弹窗 -->
    <div v-if="showUserModal && selectedUser" class="user-profile-modal" @click.self="showUserModal = false">
      <div class="user-profile-modal-overlay" @click="showUserModal = false"></div>
      <div class="user-profile-modal-content">
        <div class="flex items-center gap-4 mb-4">
          <div class="w-16 h-16 bg-gray-200 rounded-full flex items-center justify-center overflow-hidden">
            <img v-if="selectedUser.avatar" :src="selectedUser.avatar" class="w-full h-full object-cover" />
            <span v-else class="text-xl font-bold text-gray-500">{{ selectedUser.nickname?.charAt(0) || '?' }}</span>
          </div>
          <div>
            <div class="text-lg font-semibold text-gray-800">{{ selectedUser.nickname || '用户' }}</div>
            <div class="text-sm text-gray-500">{{ selectedUser.gender === 1 ? '男' : selectedUser.gender === 2 ? '女' : '未设置' }}</div>
          </div>
        </div>
        
        <div class="space-y-3 text-sm">
          <div v-if="selectedUser.location" class="flex items-center gap-2 text-gray-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span>{{ selectedUser.location }}</span>
          </div>
          <div v-if="selectedUser.car_brand" class="flex items-center gap-2 text-gray-600">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            <span>{{ selectedUser.car_brand }} {{ selectedUser.car_model || '' }}</span>
          </div>
        </div>
        
        <button 
          @click="showUserModal = false" 
          class="w-full mt-6 py-3 rounded-xl text-sm font-medium"
          :class="appStore.theme === 'dark' ? 'bg-gray-600 text-gray-200' : 'bg-gray-100 text-gray-600'"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const route = useRoute()
const router = useRouter()
const tripStore = useTripStore()
const userStore = useUserStore()
const appStore = useAppStore()

const trip = ref(null)
const loading = ref(true)
const grabbing = ref(false)
const currentImageIndex = ref(0)
const showUserModal = ref(false)
const selectedUser = ref(null)

// parse trip images from JSON string
const tripImages = computed(() => {
  if (!trip.value || !trip.value.images) return []
  try {
    const images = JSON.parse(trip.value.images)
    return Array.isArray(images) ? images.filter(url => url) : []
  } catch {
    return []
  }
})

// estimated distance based on coordinates
const estimatedDistance = computed(() => {
  if (!trip.value) return 0
  const lat1 = trip.value.departure_lat || 0
  const lng1 = trip.value.departure_lng || 0
  const lat2 = trip.value.destination_lat || 0
  const lng2 = trip.value.destination_lng || 0
  
  if (!lat1 || !lat2) return 30 // default estimate
  
  const R = 6371
  const dLat = (lat2 - lat1) * Math.PI / 180
  const dLng = (lng2 - lng1) * Math.PI / 180
  const a = Math.sin(dLat/2) * Math.sin(dLat/2) +
            Math.cos(lat1 * Math.PI / 180) * Math.cos(lat2 * Math.PI / 180) *
            Math.sin(dLng/2) * Math.sin(dLng/2)
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
  const distance = R * c
  return Math.round(distance * 1.2)
})

const estimatedTime = computed(() => {
  return Math.round(estimatedDistance.value / 40 * 60)
})

// check if current user can chat with the trip publisher
const canChat = computed(() => {
  if (!trip.value) return false
  if (!userStore.isLoggedIn) return false
  return trip.value.user_id !== userStore.user?.open_id
})

// check if current user is the owner
const isOwner = computed(() => {
  if (!trip.value) return false
  if (!userStore.isLoggedIn) return false
  return trip.value.user_id === userStore.user?.open_id
})

// check if can grab this trip
const canGrab = computed(() => {
  if (!trip.value) return false
  if (!userStore.isLoggedIn) return false
  if (isOwner.value) return false
  return trip.value.status === 1 // only pending trips can be grabbed
})

// button text based on trip type
const grabButtonText = computed(() => {
  if (!trip.value) return '立即联系'
  // trip_type=1: driver looking for passengers, so grabber is passenger -> "立即抢单"
  // trip_type=2: passenger looking for driver, so grabber is driver -> "立即接单"
  return trip.value.trip_type === 1 ? '立即抢单' : '立即接单'
})

onMounted(async () => {
  try {
    const id = route.params.id
    trip.value = await tripStore.getTripById(id)
  } finally {
    loading.value = false
  }
})

function handleImageError(e) {
  // hide broken image
  e.target.style.display = 'none'
}

function goBack() {
  router.back()
}

function goChat() {
  if (!userStore.isLoggedIn) {
    appStore.showToast('请先登录', 'info')
    router.push('/login')
    return
  }
  const peerId = trip.value.user_id
  const nickname = getMaskedNickname(trip.value)
  router.push({ path: `/chat/${peerId}`, query: { nickname } })
}

function goLogin() {
  router.push('/login')
}

async function handleGrab() {
  if (!userStore.isLoggedIn) {
    goLogin()
    return
  }
  
  grabbing.value = true
  try {
    // call grab trip API
    const result = await tripStore.grabTrip(trip.value.id)
    
    if (result && result.success) {
      appStore.showToast(result.message || '抢单成功', 'success')
      // go to chat with trip owner
      setTimeout(() => {
        goChat()
      }, 500)
    }
  } catch (error) {
    appStore.showToast(error.message || '抢单失败', 'error')
  } finally {
    grabbing.value = false
  }
}

function openNavigation() {
  if (!trip.value) return
  
  // open navigation in external app
  const destLat = trip.value.destination_lat
  const destLng = trip.value.destination_lng
  const destName = encodeURIComponent(trip.value.destination_address || trip.value.destination_city)
  
  if (destLat && destLng) {
    // try to open AMap navigation
    window.location.href = `https://uri.amap.com/navigation?to=${destLng},${destLat},${destName}&mode=car`
  } else {
    appStore.showToast('目的地坐标不完整', 'info')
  }
}

function formatDetailTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const tripDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const diffDays = Math.floor((tripDate - today) / (1000 * 60 * 60 * 24))
  
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  const timeStr = `${hour}:${minute}`
  
  // calculate end time (add estimated time)
  const endDate = new Date(date.getTime() + estimatedTime.value * 60 * 1000)
  const endHour = endDate.getHours().toString().padStart(2, '0')
  const endMinute = endDate.getMinutes().toString().padStart(2, '0')
  const endTimeStr = `${endHour}:${endMinute}`
  
  let prefix = ''
  if (diffDays === 0) {
    prefix = '今天'
  } else if (diffDays === 1) {
    prefix = '明天'
  } else if (diffDays === -1) {
    prefix = '昨天'
  } else {
    const month = date.getMonth() + 1
    const day = date.getDate()
    prefix = `${month}月${day}日`
  }
  
  return `${prefix} ${timeStr}~${endTimeStr}`
}

function getMaskedNickname(trip) {
  if (trip.user && trip.user.nickname) {
    const nickname = trip.user.nickname
    if (nickname.length <= 2) {
      return nickname[0] + '**'
    }
    return nickname.slice(0, -2) + '**'
  }
  return '用户**'
}

function getUserInitial(trip) {
  if (trip.user && trip.user.nickname) {
    return trip.user.nickname.charAt(0)
  }
  return 'U'
}

function showUserProfile(user) {
  if (user) {
    selectedUser.value = user
    showUserModal.value = true
  }
}
</script>

<style scoped>
.user-profile-modal {
  position: fixed;
  inset: 0;
  z-index: 100;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.user-profile-modal-overlay {
  position: absolute;
  inset: 0;
  background: rgba(0, 0, 0, .5);
}

.user-profile-modal-content {
  position: relative;
  width: 100%;
  max-width: 400px;
  background: #fff;
  border-radius: 24px 24px 0 0;
  padding: 24px;
  animation: slideUp .3s ease;
}

[data-theme="dark"] .user-profile-modal-content {
  background: #1f2937;
}

@keyframes slideUp {
  from {
    transform: translateY(100%);
  }
  to {
    transform: translateY(0);
  }
}
</style>
