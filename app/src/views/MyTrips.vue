<template>
  <div class="my-trips-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部导航 -->
    <div class="page-header safe-area-top">
      <div class="page-header-bg"></div>
      <div class="relative flex items-center h-12 px-4">
        <button @click="goBack" class="w-10 h-10 -ml-2 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-white">我的行程</h1>
        <router-link to="/publish" class="w-10 h-10 -mr-2 flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
        </router-link>
      </div>
    </div>

    <!-- 筛选标签 -->
    <div class="px-4 py-3 sticky top-0 z-10" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
      <div class="flex gap-2 overflow-x-auto pb-1">
        <button
          v-for="tab in statusTabs"
          :key="tab.value"
          @click="currentStatus = tab.value"
          class="px-4 py-2 rounded-full text-sm font-medium whitespace-nowrap transition-all"
          :class="currentStatus === tab.value 
            ? 'text-white' 
            : appStore.theme === 'dark' ? 'bg-gray-800 text-gray-400' : 'bg-white text-gray-500'"
          :style="currentStatus === tab.value ? { backgroundColor: 'var(--theme-primary)' } : {}"
        >
          {{ tab.label }}
          <span 
            v-if="tab.count > 0" 
            class="ml-1"
            :class="currentStatus === tab.value ? 'text-white/70' : 'text-gray-400'"
          >
            ({{ tab.count }})
          </span>
        </button>
      </div>
    </div>

    <!-- 列表 -->
    <div class="px-4 pb-4">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="filteredTrips.length === 0" class="empty-state py-16">
        <div class="empty-state-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <p class="empty-state-text">{{ currentStatus === null ? '暂无行程' : '暂无该状态的行程' }}</p>
        <router-link to="/publish" class="text-sm" :style="{ color: 'var(--theme-primary)' }">
          发布行程
        </router-link>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="trip in filteredTrips"
          :key="trip.id"
          @click="goToDetail(trip.id)"
          class="card p-4 cursor-pointer active:scale-98 transition-transform"
        >
          <!-- 第一行：类型 + 状态 + 操作 -->
          <div class="flex items-center justify-between mb-3">
            <div class="flex items-center gap-2">
              <span
                class="badge"
                :class="trip.trip_type === 1 ? 'trip-type-driver' : 'trip-type-passenger'"
              >
                {{ trip.trip_type === 1 ? '车找人' : '人找车' }}
              </span>
              <span
                class="badge"
                :class="getStatusClass(trip.status)"
              >
                {{ getStatusText(trip.status) }}
              </span>
            </div>
            
            <!-- 数据统计 -->
            <div class="flex items-center gap-3 text-xs text-gray-400">
              <span class="flex items-center gap-1">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                {{ trip.view_count || 0 }}
              </span>
              <span class="flex items-center gap-1">
                <svg class="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                {{ trip.grabbers?.length || 0 }}
              </span>
            </div>
          </div>

          <!-- 路线信息 -->
          <div class="flex items-center gap-2 mb-2">
            <span class="w-2 h-2 rounded-full flex-shrink-0" style="background: var(--theme-primary);"></span>
            <span class="text-sm font-medium truncate flex-1">{{ trip.departure_city }}</span>
            <svg class="w-4 h-4 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
            <span class="w-2 h-2 rounded-full bg-orange-500 flex-shrink-0"></span>
            <span class="text-sm font-medium truncate flex-1">{{ trip.destination_city }}</span>
          </div>

          <!-- 时间和价格 -->
          <div class="flex items-center justify-between text-sm">
            <span class="text-gray-500">🕐 {{ formatTime(trip.departure_time) }}</span>
            <div class="flex items-center gap-2">
              <span class="text-gray-500">{{ trip.seats }}座</span>
              <span v-if="trip.price > 0" class="font-semibold" :style="{ color: 'var(--theme-primary)' }">
                ¥{{ trip.price }}/人
              </span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div v-if="trip.status === 1" class="flex gap-2 mt-3 pt-3 border-t border-gray-100" :class="appStore.theme === 'dark' ? 'border-gray-700' : ''">
            <button
              @click.stop="handleShare(trip)"
              class="flex-1 py-2 text-xs font-medium border rounded-lg transition-colors"
              :class="appStore.theme === 'dark' ? 'border-gray-600 text-gray-300 hover:bg-gray-700' : 'border-gray-200 text-gray-600 hover:bg-gray-50'"
            >
              分享行程
            </button>
            <button
              @click.stop="handleComplete(trip.id)"
              class="flex-1 py-2 text-xs font-medium text-white rounded-lg transition-colors bg-green-500 hover:bg-green-600"
            >
              已成行
            </button>
            <button
              @click.stop="handleCancel(trip.id)"
              class="flex-1 py-2 text-xs font-medium border rounded-lg transition-colors"
              :class="appStore.theme === 'dark' ? 'border-gray-600 text-gray-300 hover:bg-gray-700' : 'border-gray-200 text-gray-600 hover:bg-gray-50'"
            >
              取消
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 分享弹窗 -->
    <div v-if="showShareModal" class="action-sheet" @click.self="showShareModal = false">
      <div class="action-sheet-overlay" @click="showShareModal = false"></div>
      <div class="action-sheet-content">
        <div class="px-4 py-4 border-b border-gray-100" :class="appStore.theme === 'dark' ? 'border-gray-700' : ''">
          <h3 class="text-lg font-semibold text-center">分享行程</h3>
        </div>
        
        <div class="p-4">
          <!-- 分享预览 -->
          <div class="card p-4 mb-4" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'">
            <div class="text-sm font-medium mb-2">
              {{ shareTrip?.trip_type === 1 ? '🚗 车找人' : '🙋 人找车' }}
            </div>
            <div class="text-sm">
              📍 {{ shareTrip?.departure_city }} → {{ shareTrip?.destination_city }}
            </div>
            <div class="text-sm mt-1">
              🕐 {{ formatTime(shareTrip?.departure_time) }}
            </div>
            <div class="text-sm mt-1">
              💺 {{ shareTrip?.seats }}座{{ shareTrip?.price > 0 ? ` · ¥${shareTrip?.price}/人` : '' }}
            </div>
          </div>

          <!-- 分享方式 -->
          <div class="grid grid-cols-4 gap-4 mb-4">
            <button 
              @click="shareToWechat"
              class="flex flex-col items-center gap-2 active:scale-95 transition-transform"
            >
              <div class="w-12 h-12 rounded-full bg-green-500 flex items-center justify-center">
                <svg class="w-7 h-7 text-white" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 01.213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 00.167-.054l1.903-1.114a.864.864 0 01.717-.098 10.16 10.16 0 002.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348zM5.785 5.991c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178A1.17 1.17 0 014.623 7.17c0-.651.52-1.18 1.162-1.18zm5.813 0c.642 0 1.162.529 1.162 1.18a1.17 1.17 0 01-1.162 1.178 1.17 1.17 0 01-1.162-1.178c0-.651.52-1.18 1.162-1.18z"/>
                  <path d="M23.27 14.643c0-3.171-3.022-5.794-6.66-5.794-3.715 0-6.738 2.548-6.738 5.794 0 3.245 3.022 5.871 6.738 5.871.708 0 1.387-.097 2.04-.276a.637.637 0 01.51.068l1.379.81a.24.24 0 00.12.04.213.213 0 00.213-.21c0-.053-.02-.103-.036-.152l-.281-1.067a.431.431 0 01.155-.48c1.408-1.05 2.56-2.727 2.56-4.604zm-8.894-1.022c-.474 0-.857-.39-.857-.87 0-.48.383-.87.857-.87s.857.39.857.87c0 .48-.383.87-.857.87zm4.468 0c-.474 0-.857-.39-.857-.87 0-.48.383-.87.857-.87s.858.39.858.87c0 .48-.384.87-.858.87z"/>
                </svg>
              </div>
              <span class="text-xs">微信好友</span>
            </button>

            <button 
              @click="shareToMoments"
              class="flex flex-col items-center gap-2 active:scale-95 transition-transform"
            >
              <div class="w-12 h-12 rounded-full bg-green-600 flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
              </div>
              <span class="text-xs">朋友圈</span>
            </button>

            <button 
              @click="copyShareLink"
              class="flex flex-col items-center gap-2 active:scale-95 transition-transform"
            >
              <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-200'">
                <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                </svg>
              </div>
              <span class="text-xs">复制链接</span>
            </button>

            <button 
              @click="copyShareText"
              class="flex flex-col items-center gap-2 active:scale-95 transition-transform"
            >
              <div class="w-12 h-12 rounded-full flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-200'">
                <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5H6a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2v-1M8 5a2 2 0 002 2h2a2 2 0 002-2M8 5a2 2 0 012-2h2a2 2 0 012 2m0 0h2a2 2 0 012 2v3m2 4H10m0 0l3-3m-3 3l3 3" />
                </svg>
              </div>
              <span class="text-xs">复制文字</span>
            </button>
          </div>

          <button 
            @click="showShareModal = false"
            class="btn btn-secondary w-full"
          >
            取消
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useAppStore } from '@/stores/app'

const router = useRouter()
const tripStore = useTripStore()
const appStore = useAppStore()

const trips = ref([])
const loading = ref(true)
const currentStatus = ref(null)
const showShareModal = ref(false)
const shareTrip = ref(null)

const statusTabs = computed(() => [
  { value: null, label: '全部', count: trips.value.length },
  { value: 1, label: '待匹配', count: trips.value.filter(t => t.status === 1).length },
  { value: 2, label: '已匹配', count: trips.value.filter(t => t.status === 2).length },
  { value: 3, label: '已完成', count: trips.value.filter(t => t.status === 3).length },
  { value: 4, label: '已取消', count: trips.value.filter(t => t.status === 4).length }
])

const filteredTrips = computed(() => {
  if (currentStatus.value === null) return trips.value
  return trips.value.filter(t => t.status === currentStatus.value)
})

onMounted(async () => {
  try {
    trips.value = await tripStore.fetchMyTrips() || []
  } finally {
    loading.value = false
  }
})

function goBack() {
  router.back()
}

function goToDetail(id) {
  router.push(`/my-trip/${id}`)
}

function getStatusClass(status) {
  switch (status) {
    case 1: return 'status-pending'
    case 2: return 'status-matched'
    case 3: return 'status-completed'
    case 4: return 'status-cancelled'
    default: return 'badge-secondary'
  }
}

function getStatusText(status) {
  switch (status) {
    case 1: return '待匹配'
    case 2: return '已匹配'
    case 3: return '已完成'
    case 4: return '已取消'
    default: return '未知'
  }
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hour}:${minute}`
}

function handleShare(trip) {
  shareTrip.value = trip
  showShareModal.value = true
}

function getShareText() {
  if (!shareTrip.value) return ''
  const trip = shareTrip.value
  return `🚗 ${trip.trip_type === 1 ? '车找人' : '人找车'}\n📍 ${trip.departure_city} → ${trip.destination_city}\n🕐 ${formatTime(trip.departure_time)}\n💺 ${trip.seats}座${trip.price > 0 ? ` · ¥${trip.price}/人` : ''}`
}

function getShareUrl() {
  if (!shareTrip.value) return ''
  return `${window.location.origin}/trip/${shareTrip.value.id}`
}

async function shareToWechat() {
  // native share or copy
  if (navigator.share) {
    try {
      await navigator.share({
        title: '拼车行程分享',
        text: getShareText(),
        url: getShareUrl()
      })
      showShareModal.value = false
      return
    } catch (err) {
      if (err.name === 'AbortError') return
    }
  }
  // fallback
  copyShareText()
}

function shareToMoments() {
  // same as wechat share
  shareToWechat()
}

async function copyShareLink() {
  try {
    await navigator.clipboard.writeText(getShareUrl())
    appStore.showToast('链接已复制到剪贴板', 'success')
    showShareModal.value = false
  } catch (err) {
    appStore.showToast('复制失败', 'error')
  }
}

async function copyShareText() {
  try {
    await navigator.clipboard.writeText(`${getShareText()}\n🔗 ${getShareUrl()}`)
    appStore.showToast('行程信息已复制，可粘贴到微信分享', 'success')
    showShareModal.value = false
  } catch (err) {
    appStore.showToast('复制失败', 'error')
  }
}

async function handleComplete(id) {
  if (!confirm('确定要标记这个行程为已成行吗？')) return
  try {
    await tripStore.completeTrip(id)
    appStore.showToast('行程已标记为已成行', 'success')
    trips.value = await tripStore.fetchMyTrips() || []
  } catch (e) {
    // error handled
  }
}

async function handleCancel(id) {
  if (!confirm('确定要取消这个行程吗？')) return
  try {
    await tripStore.cancelTrip(id)
    appStore.showToast('行程已取消', 'success')
    trips.value = await tripStore.fetchMyTrips() || []
  } catch (e) {
    // error handled
  }
}
</script>
