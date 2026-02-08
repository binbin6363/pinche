<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <div class="grid grid-cols-2 gap-3 mb-6">
      <div class="stat-card cursor-pointer active:scale-[0.98] transition-transform" @click="goUsers()">
        <div class="flex items-center justify-between">
          <div>
            <div class="stat-value">{{ stats.total_users }}</div>
            <div class="stat-label">总用户</div>
          </div>
          <div class="w-10 h-10 bg-blue-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
          </div>
        </div>
        <div class="stat-change positive">今日 +{{ stats.today_users }}</div>
      </div>

      <div class="stat-card cursor-pointer active:scale-[0.98] transition-transform" @click="goTrips()">
        <div class="flex items-center justify-between">
          <div>
            <div class="stat-value">{{ stats.total_trips }}</div>
            <div class="stat-label">总行程</div>
          </div>
          <div class="w-10 h-10 bg-green-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM19 17a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
        </div>
        <div class="stat-change positive">今日 +{{ stats.today_trips }}</div>
      </div>

      <div class="stat-card cursor-pointer active:scale-[0.98] transition-transform" @click="goTrips('1')">
        <div class="flex items-center justify-between">
          <div>
            <div class="stat-value">{{ stats.active_trips }}</div>
            <div class="stat-label">进行中</div>
          </div>
          <div class="w-10 h-10 bg-yellow-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div class="text-xs text-gray-400 mt-1">待匹配行程</div>
      </div>

      <div class="stat-card cursor-pointer active:scale-[0.98] transition-transform" @click="goBannedUsers()">
        <div class="flex items-center justify-between">
          <div>
            <div class="stat-value">{{ stats.banned_users }}</div>
            <div class="stat-label">封禁用户</div>
          </div>
          <div class="w-10 h-10 bg-red-100 rounded-xl flex items-center justify-center">
            <svg class="w-5 h-5 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728A9 9 0 015.636 5.636m12.728 12.728L5.636 5.636" />
            </svg>
          </div>
        </div>
        <div class="text-xs text-gray-400 mt-1">累计封禁</div>
      </div>
    </div>

    <!-- 最近用户 -->
    <div class="card mb-4">
      <div class="p-4 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-semibold text-gray-800">最近注册</h3>
        <router-link to="/users" class="text-sm text-blue-500">查看全部</router-link>
      </div>
      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>
      <div v-else-if="recentUsers.length === 0" class="empty-state py-6">
        <p>暂无数据</p>
      </div>
      <div v-else>
        <div
          v-for="user in recentUsers"
          :key="user.id"
          class="list-item cursor-pointer"
          @click="router.push('/users')"
        >
          <div class="w-10 h-10 bg-gray-100 rounded-full flex items-center justify-center mr-3">
            <span class="text-sm font-medium text-gray-500">{{ user.nickname?.charAt(0) || 'U' }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="font-medium text-gray-800 truncate">{{ user.nickname }}</div>
            <div class="text-xs text-gray-400">{{ user.phone }}</div>
          </div>
          <div class="text-xs text-gray-400">{{ formatTime(user.created_at) }}</div>
        </div>
      </div>
    </div>

    <!-- 最近行程 -->
    <div class="card">
      <div class="p-4 border-b border-gray-100 flex items-center justify-between">
        <h3 class="font-semibold text-gray-800">最近行程</h3>
        <router-link to="/trips" class="text-sm text-blue-500">查看全部</router-link>
      </div>
      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>
      <div v-else-if="recentTrips.length === 0" class="empty-state py-6">
        <p>暂无数据</p>
      </div>
      <div v-else>
        <div
          v-for="trip in recentTrips"
          :key="trip.id"
          class="list-item cursor-pointer"
          @click="router.push('/trips')"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="font-medium text-gray-800">{{ trip.departure_city }}</span>
              <svg class="w-4 h-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
              <span class="font-medium text-gray-800">{{ trip.destination_city }}</span>
            </div>
            <div class="text-xs text-gray-400">
              {{ trip.trip_type === 1 ? '司机' : '乘客' }} · {{ trip.user?.nickname || '未知' }}
            </div>
          </div>
          <span :class="getStatusBadgeClass(trip.status)" class="badge">
            {{ getStatusText(trip.status) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/utils/api'

const router = useRouter()
const loading = ref(true)
const stats = ref({
  total_users: 0,
  total_trips: 0,
  active_trips: 0,
  banned_users: 0,
  today_users: 0,
  today_trips: 0
})
const recentUsers = ref([])
const recentTrips = ref([])

onMounted(async () => {
  await fetchData()
})

async function fetchData() {
  loading.value = true
  try {
    const [statsData, usersData, tripsData] = await Promise.all([
      api.get('/admin/stats'),
      api.get('/admin/users', { params: { page: 1, page_size: 5 } }),
      api.get('/admin/trips', { params: { page: 1, page_size: 5 } })
    ])
    stats.value = statsData || {}
    recentUsers.value = usersData.list || []
    recentTrips.value = tripsData.list || []
  } catch (e) {
    console.error('Failed to fetch dashboard data:', e)
  } finally {
    loading.value = false
  }
}

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}/${day} ${hour}:${minute}`
}

function getStatusBadgeClass(status) {
  switch (status) {
    case 1: return 'badge-warning'
    case 2: return 'badge-success'
    case 3: return 'badge-gray'
    case 4: return 'badge-danger'
    case 5: return 'badge-danger'
    default: return 'badge-gray'
  }
}

function getStatusText(status) {
  switch (status) {
    case 1: return '待匹配'
    case 2: return '已匹配'
    case 3: return '已完成'
    case 4: return '已取消'
    case 5: return '已封禁'
    default: return '未知'
  }
}

function goUsers() {
  router.push('/users')
}

function goTrips(status) {
  if (status) {
    router.push({ path: '/trips', query: { status } })
  } else {
    router.push('/trips')
  }
}

function goBannedUsers() {
  router.push({ path: '/users', query: { status: '1' } })
}
</script>
