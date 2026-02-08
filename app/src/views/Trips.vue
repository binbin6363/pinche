<template>
  <div class="trips-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 喜庆主题顶栏 -->
    <div class="page-header safe-area-top">
      <div class="page-header-bg"></div>
      
      <!-- 装饰元素 - 仅新春主题 -->
      <template v-if="appStore.theme === 'spring'">
        <div class="absolute top-1 right-4 text-yellow-300 opacity-30 text-3xl">🧧</div>
        <div class="absolute top-2 left-4 text-yellow-200 opacity-20 text-2xl">✨</div>
        <div class="absolute -top-1 right-16 flex gap-2">
          <div class="spring-lantern"></div>
          <div class="spring-lantern w-4 h-6"></div>
        </div>
      </template>
      
      <!-- 内容 -->
      <div class="relative px-4 py-4 text-white">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-2 text-sm">
            <span class="text-lg">{{ userStore.identity === 1 ? '🚗' : '🙋' }}</span>
            <span>
              <span class="font-medium">{{ userStore.identity === 1 ? '司机' : '乘客' }}</span>
              · 查看{{ userStore.identity === 1 ? '乘客' : '司机' }}行程
            </span>
          </div>
          <router-link to="/settings" class="text-xs bg-white/20 px-3 py-1.5 rounded-full hover:bg-white/30 transition-colors">
            切换身份
          </router-link>
        </div>
      </div>
    </div>

    <!-- 安全提示 - 可关闭 -->
    <div v-if="showSafetyTip" class="px-4 py-2" :class="appStore.theme === 'dark' ? 'bg-amber-900/30' : 'bg-amber-50'">
      <div class="flex items-center justify-between">
        <div class="flex items-center gap-2 flex-1 min-w-0">
          <span class="text-sm">⚠️</span>
          <p class="text-xs truncate" :class="appStore.theme === 'dark' ? 'text-amber-400' : 'text-amber-700'">
            自助拼车，勿提供个人敏感信息 · 公共场所见面 · 注意安全
          </p>
        </div>
        <button 
          @click="closeSafetyTip"
          class="ml-2 p-2 flex-shrink-0"
          :class="appStore.theme === 'dark' ? 'text-amber-400 hover:text-amber-300' : 'text-amber-500 hover:text-amber-700'"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 筛选栏 -->
    <div class="px-4 py-3 border-b sticky top-0 z-10" 
         :class="appStore.theme === 'dark' ? 'bg-gray-900 border-gray-800' : 'bg-white border-gray-100'">
      <div class="flex gap-2">
        <input
          v-model="filter.departure_city"
          type="text"
          placeholder="出发城市"
          class="input flex-1 text-sm"
        />
        <input
          v-model="filter.destination_city"
          type="text"
          placeholder="目的城市"
          class="input flex-1 text-sm"
        />
        <button @click="doSearch" class="btn btn-primary px-4">
          搜索
        </button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="pb-4">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="trips.length === 0" class="empty-state py-16">
        <div class="empty-state-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <p class="empty-state-text">暂无相关行程</p>
        <p class="text-xs text-gray-400">可以修改筛选条件试试</p>
      </div>

      <div v-else>
        <!-- 行程卡片列表 -->
        <div
          v-for="trip in trips"
          :key="trip.id"
          @click="goTripDetail(trip.id)"
          class="px-4 py-4 border-b cursor-pointer active:bg-opacity-50 transition-colors"
          :class="appStore.theme === 'dark' 
            ? 'bg-gray-900 border-gray-800 active:bg-gray-800' 
            : 'bg-white border-gray-100 active:bg-gray-50'"
        >
          <!-- 顶部：时间 + 类型标签 -->
          <div class="flex items-center justify-between mb-3">
            <div class="text-base font-semibold">
              {{ formatTripDate(trip.departure_time) }}
            </div>
            <span
              class="badge"
              :class="trip.trip_type === 1 ? 'trip-type-driver' : 'trip-type-passenger'"
            >
              {{ trip.trip_type === 1 ? '车找人' : '人找车' }}
            </span>
          </div>

          <!-- 中间：起点→终点 一行展示 -->
          <div class="flex items-center gap-3">
            <div class="flex-1 min-w-0 flex items-center gap-2">
              <span class="w-2 h-2 rounded-full flex-shrink-0" style="background: var(--theme-primary);"></span>
              <span class="text-sm truncate max-w-[100px]">{{ trip.departure_address || trip.departure_city }}</span>
              <svg class="w-4 h-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
              </svg>
              <span class="w-2 h-2 rounded-full bg-orange-500 flex-shrink-0"></span>
              <span class="text-sm truncate max-w-[100px]">{{ trip.destination_address || trip.destination_city }}</span>
            </div>
            
            <!-- 右侧：价格 -->
            <div class="flex flex-col items-end flex-shrink-0">
              <div v-if="trip.price > 0" class="text-xl font-bold">
                {{ trip.price }}<span class="text-sm font-normal text-gray-500">元</span>
              </div>
              <div v-else class="text-sm text-gray-400">面议</div>
            </div>
          </div>

          <!-- 底部：座位数、标签 -->
          <div class="flex items-center gap-2 mt-3">
            <span class="badge badge-secondary">
              {{ trip.seats }}人
            </span>
            <span v-if="trip.remark" class="text-xs text-gray-400 truncate max-w-[180px]">
              {{ trip.remark }}
            </span>
          </div>
        </div>
      </div>

      <!-- 分页器 -->
      <div v-if="trips.length > 0" class="flex items-center justify-center gap-3 mt-4 pt-2 px-4">
        <button
          @click="prevPage"
          :disabled="page === 1"
          class="btn btn-secondary px-4 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          上一页
        </button>
        <span class="text-sm" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'">
          {{ page }} / {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="page >= totalPages"
          class="btn btn-secondary px-4 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          下一页
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const tripStore = useTripStore()
const userStore = useUserStore()
const appStore = useAppStore()

const trips = ref([])
const loading = ref(true)
const page = ref(1)
const pageSize = 15
const total = ref(0)
const showSafetyTip = ref(true)

const filter = reactive({
  departure_city: '',
  destination_city: ''
})

onMounted(() => {
  const closed = localStorage.getItem('safety_tip_closed')
  if (closed) {
    showSafetyTip.value = false
  }
  fetchTrips()
})

// show opposite identity trips
const targetTripType = computed(() => {
  return userStore.identity === 1 ? 2 : 1
})

const totalPages = computed(() => {
  return Math.max(1, Math.ceil(total.value / pageSize))
})

async function fetchTrips() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize,
      trip_type: targetTripType.value
    }
    if (filter.departure_city) {
      params.departure_city = filter.departure_city
    }
    if (filter.destination_city) {
      params.destination_city = filter.destination_city
    }
    if (userStore.user?.city) {
      params.user_city = userStore.user.city
    }
    if (userStore.user?.province) {
      params.user_province = userStore.user.province
    }
    const result = await tripStore.fetchTrips(params)
    trips.value = result.list || []
    total.value = result.total || 0
  } finally {
    loading.value = false
  }
}

function doSearch() {
  page.value = 1
  fetchTrips()
}

function prevPage() {
  if (page.value > 1) {
    page.value--
    fetchTrips()
    window.scrollTo(0, 0)
  }
}

function nextPage() {
  if (page.value < totalPages.value) {
    page.value++
    fetchTrips()
    window.scrollTo(0, 0)
  }
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
  
  if (diffDays === 0) {
    return `今天 ${timeStr}`
  } else if (diffDays === 1) {
    return `明天 ${timeStr}`
  } else if (diffDays === -1) {
    return `昨天 ${timeStr}`
  } else {
    const month = date.getMonth() + 1
    const day = date.getDate()
    return `${month}月${day}日 ${timeStr}`
  }
}

function closeSafetyTip() {
  showSafetyTip.value = false
  localStorage.setItem('safety_tip_closed', '1')
}
</script>
