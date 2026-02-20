<template>
  <div class="home-page min-h-screen" :class="bgClass">
    <!-- 顶部区域 -->
    <div class="page-header">
      <div class="page-header-bg"></div>
      
      <!-- 装饰元素 - 仅新春主题 -->
      <template v-if="appStore.theme === 'spring'">
        <div class="absolute top-0 right-0 text-yellow-300/20 text-6xl">🧧</div>
        <div class="absolute bottom-0 left-4 text-yellow-200/20 text-3xl">✨</div>
        <div class="absolute -top-1 right-16 flex gap-2">
          <div class="spring-lantern"></div>
          <div class="spring-lantern w-4 h-6"></div>
        </div>
      </template>

      <!-- 顶部导航 -->
      <div class="relative flex items-center justify-between px-4 py-3 safe-area-top">
        <h1 class="text-lg font-bold text-white">春节拼车</h1>
        <button 
          @click="showGuide = true" 
          class="w-10 h-10 rounded-full bg-white/20 backdrop-blur-sm flex items-center justify-center active:scale-95 transition-all"
        >
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </button>
      </div>

      <!-- 欢迎语 -->
      <div class="relative px-4 pb-6 text-white">
        <div class="text-2xl font-bold mb-1">
          {{ greeting }}
        </div>
        <p class="text-white/80 text-sm">
          {{ userStore.identity === 1 ? '发布行程，找到同路乘客' : '寻找顺路的司机，安全回家' }}
        </p>
      </div>
    </div>

    <!-- 公告卡片 - 轮播 -->
    <div class="px-4 -mt-4 relative z-10">
      <div v-if="announcements.length > 0" class="card p-4 overflow-hidden">
        <div class="flex items-start gap-3">
          <div class="w-10 h-10 rounded-full flex items-center justify-center flex-shrink-0"
               :class="appStore.theme === 'spring' ? 'bg-red-100' : appStore.theme === 'dark' ? 'bg-blue-900/30' : 'bg-blue-100'">
            <span class="text-lg">📢</span>
          </div>
          <div class="flex-1 min-w-0 relative">
            <transition name="ann-slide" mode="out-in">
              <div :key="currentAnnIndex">
                <h4 class="text-sm font-semibold mb-1">{{ currentAnnouncement.title }}</h4>
                <p class="text-xs text-gray-500 leading-relaxed line-clamp-2">
                  {{ currentAnnouncement.content }}
                </p>
              </div>
            </transition>
            <!-- 轮播指示器 -->
            <div v-if="announcements.length > 1" class="flex items-center gap-1 mt-2">
              <span
                v-for="(_, idx) in announcements"
                :key="idx"
                class="w-1.5 h-1.5 rounded-full transition-all duration-300"
                :class="idx === currentAnnIndex 
                  ? (appStore.theme === 'dark' ? 'bg-blue-400 w-3' : 'bg-blue-500 w-3') 
                  : (appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-300')"
              ></span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 快速操作入口 -->
    <div class="px-4 py-4">
      <div class="grid grid-cols-2 gap-3">
        <button 
          @click="quickPublish(1)"
          class="card p-4 flex items-center gap-3 active:scale-98 transition-transform"
        >
          <div class="w-12 h-12 rounded-xl flex items-center justify-center" :style="{ background: 'var(--theme-primary)' }">
            <span class="text-2xl">🚗</span>
          </div>
          <div class="text-left">
            <div class="font-semibold">我是司机</div>
            <div class="text-xs text-gray-500">发布车找人</div>
          </div>
        </button>
        
        <button 
          @click="quickPublish(2)"
          class="card p-4 flex items-center gap-3 active:scale-98 transition-transform"
        >
          <div class="w-12 h-12 rounded-xl flex items-center justify-center" :style="{ background: 'var(--theme-primary)', opacity: '.85' }">
            <span class="text-2xl">🙋</span>
          </div>
          <div class="text-left">
            <div class="font-semibold">我是乘客</div>
            <div class="text-xs text-gray-500">发布人找车</div>
          </div>
        </button>
      </div>
    </div>

    <!-- 推荐行程 -->
    <div class="px-4 pb-4">
      <div class="section-header px-0">
        <h3 class="section-title">
          <span>🔥</span>
          推荐行程
        </h3>
        <router-link to="/trips" class="text-sm" :style="{ color: 'var(--theme-primary)' }">
          查看更多
        </router-link>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="trips.length === 0" class="empty-state">
        <div class="empty-state-icon">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>
        <p class="empty-state-text">暂无推荐行程</p>
        <router-link to="/publish" class="text-sm" :style="{ color: 'var(--theme-primary)' }">
          发布第一个行程
        </router-link>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="trip in trips"
          :key="trip.id"
          @click="goTripDetail(trip.id)"
          class="card p-4 active:scale-98 transition-transform cursor-pointer"
        >
          <!-- 顶部信息 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span
                class="badge"
                :class="trip.trip_type === 1 ? 'trip-type-driver' : 'trip-type-passenger'"
              >
                {{ trip.trip_type === 1 ? '车找人' : '人找车' }}
              </span>
              <span class="text-xs text-gray-400">{{ formatTripDate(trip.departure_time) }}</span>
            </div>
            <div class="text-lg font-bold" :style="{ color: 'var(--theme-primary)' }">
              <template v-if="trip.price > 0">
                ¥{{ trip.price }}
              </template>
              <span v-else class="text-sm text-gray-400">面议</span>
            </div>
          </div>

          <!-- 路线 -->
          <div class="flex items-center gap-2">
            <span class="w-2 h-2 rounded-full" :style="{ background: 'var(--theme-primary)' }"></span>
            <span class="text-sm font-medium truncate flex-1">{{ trip.departure_city }}</span>
            <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
            <span class="w-2 h-2 rounded-full" :style="{ background: 'var(--theme-primary)', opacity: '.6' }"></span>
            <span class="text-sm font-medium truncate flex-1">{{ trip.destination_city }}</span>
          </div>

          <!-- 底部标签 -->
          <div class="flex items-center gap-2 mt-3 pt-3 border-t" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
            <span class="badge badge-secondary">{{ trip.seats }}座</span>
            <span v-if="trip.remark" class="text-xs text-gray-400 truncate flex-1">{{ trip.remark }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 使用指南弹窗 -->
    <div v-if="showGuide" class="action-sheet" @click.self="showGuide = false">
      <div class="action-sheet-overlay" @click="showGuide = false"></div>
      <div class="action-sheet-content max-h-[80vh]">
        <!-- 头部 -->
        <div class="page-header">
          <div class="page-header-bg"></div>
          <div class="relative px-4 py-4 text-white flex items-center justify-between">
            <h3 class="text-lg font-bold flex items-center gap-2">
              <span>📖</span>
              使用指南
            </h3>
            <button @click="showGuide = false" class="w-8 h-8 rounded-full bg-white/20 flex items-center justify-center">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
        
        <!-- 内容 -->
        <div class="p-4 overflow-y-auto max-h-[60vh]">
          <!-- 司机攻略 -->
          <div class="mb-4">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-full flex items-center justify-center" :style="{ background: 'var(--theme-primary)' }">
                <span class="text-white">🚗</span>
              </div>
              <span class="font-semibold">车找人攻略</span>
            </div>
            <div class="pl-10 space-y-2 text-sm" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-600'">
              <p>① 点击「我是司机」发布行程</p>
              <p>② 在「广场」浏览乘客需求</p>
              <p>③ 私聊沟通或等待系统匹配</p>
              <p>④ 成行后标记行程状态</p>
            </div>
          </div>

          <!-- 乘客攻略 -->
          <div class="mb-4">
            <div class="flex items-center gap-2 mb-2">
              <div class="w-8 h-8 rounded-full flex items-center justify-center" :style="{ background: 'var(--theme-primary)', opacity: '.85' }">
                <span class="text-white">🙋</span>
              </div>
              <span class="font-semibold">人找车攻略</span>
            </div>
            <div class="pl-10 space-y-2 text-sm" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-600'">
              <p>① 点击「我是乘客」发布行程</p>
              <p>② 在「广场」浏览司机行程</p>
              <p>③ 私聊沟通或等待系统匹配</p>
              <p>④ 成行后标记行程状态</p>
            </div>
          </div>

          <!-- 安全提示 -->
          <div class="rounded-xl p-4" :class="appStore.theme === 'dark' ? 'bg-amber-900/20' : 'bg-amber-50'">
            <div class="flex items-center gap-2 mb-2">
              <span>⚠️</span>
              <span class="font-medium" :class="appStore.theme === 'dark' ? 'text-amber-400' : 'text-amber-700'">安全提示</span>
            </div>
            <ul class="space-y-1.5 text-xs" :class="appStore.theme === 'dark' ? 'text-amber-400/80' : 'text-amber-600'">
              <li>• 出行前核实对方身份信息</li>
              <li>• 在公共场所见面交接</li>
              <li>• 告知家人出行信息</li>
              <li>• 保管好个人贵重物品</li>
            </ul>
          </div>
        </div>

        <!-- 底部按钮 -->
        <div class="p-4 border-t" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
          <button 
            @click="showGuide = false" 
            class="btn btn-primary w-full"
          >
            我知道了
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import api from '@/utils/api'

const router = useRouter()
const tripStore = useTripStore()
const userStore = useUserStore()
const appStore = useAppStore()

const trips = ref([])
const announcements = ref([])
const loading = ref(true)
const showGuide = ref(false)
const currentAnnIndex = ref(0)
let annTimer = null

const bgClass = computed(() => {
  return appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'
})

const greeting = computed(() => {
  const hour = new Date().getHours()
  const nickname = userStore.user?.nickname || '旅客'
  if (hour < 6) return `夜深了，${nickname}`
  if (hour < 12) return `早上好，${nickname}`
  if (hour < 18) return `下午好，${nickname}`
  return `晚上好，${nickname}`
})

// show opposite identity trips
const targetTripType = computed(() => {
  return userStore.identity === 1 ? 2 : 1
})

// current announcement for carousel
const currentAnnouncement = computed(() => {
  return announcements.value[currentAnnIndex.value] || {}
})

// start announcement carousel
function startAnnCarousel() {
  if (announcements.value.length <= 1) return
  annTimer = setInterval(() => {
    currentAnnIndex.value = (currentAnnIndex.value + 1) % announcements.value.length
  }, 5000)
}

onMounted(async () => {
  // fetch announcements (default 3 from server)
  try {
    const annResult = await api.get('/announcements')
    announcements.value = annResult || []
    startAnnCarousel()
  } catch (e) {
    // ignore
  }

  // fetch recommended trips (exclude current user's trips)
  try {
    const params = { 
      page: 1, 
      page_size: 5,
      trip_type: targetTripType.value
    }
    // exclude current user's trips if logged in
    if (userStore.user?.open_id) {
      params.exclude_user_id = userStore.user.open_id
    }
    const result = await tripStore.fetchTrips(params)
    trips.value = result.list || []
  } finally {
    loading.value = false
  }
})

onUnmounted(() => {
  if (annTimer) {
    clearInterval(annTimer)
    annTimer = null
  }
})

function quickPublish(identity) {
  userStore.setIdentity(identity)
  router.push('/publish')
}

function goTripDetail(id) {
  router.push(`/trip/${id}`)
}

function formatTripDate(time) {
  const date = new Date(time)
  const now = new Date()
  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const tripDate = new Date(date.getFullYear(), date.getMonth(), date.getDate())
  const diffDays = Math.floor((tripDate - today) / (1000 * 60 * 60 * 24))
  
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  const timeStr = `${hour}:${minute}`
  
  if (diffDays === 0) return `今天 ${timeStr}`
  if (diffDays === 1) return `明天 ${timeStr}`
  if (diffDays === -1) return `昨天 ${timeStr}`
  
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日 ${timeStr}`
}
</script>

<style scoped>
/* announcement carousel transition */
.ann-slide-enter-active,
.ann-slide-leave-active {
  transition: all .3s ease;
}
.ann-slide-enter-from {
  opacity: 0;
  transform: translateY(8px);
}
.ann-slide-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}
</style>
