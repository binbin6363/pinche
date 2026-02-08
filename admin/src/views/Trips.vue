<template>
  <div class="trips-page">
    <!-- 用户筛选提示 -->
    <div v-if="filterUserName" class="card mb-4">
      <div class="p-3 flex items-center justify-between">
        <div class="flex items-center gap-2">
          <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          <span class="text-sm text-gray-600">筛选用户: <strong class="text-gray-800">{{ filterUserName }}</strong></span>
        </div>
        <button @click="clearUserFilter" class="text-xs text-blue-500 px-2 py-1">清除筛选</button>
      </div>
    </div>

    <!-- 搜索栏 -->
    <div class="card mb-4">
      <div class="p-3">
        <div class="search-bar mb-3">
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="search"
            type="text"
            placeholder="搜索出发地/目的地"
            @keyup.enter="fetchTrips"
          />
        </div>
        <!-- 类型筛选 -->
        <div class="flex gap-2 mb-2 flex-wrap">
          <button
            v-for="opt in typeOptions"
            :key="opt.value"
            @click="setTypeFilter(opt.value)"
            class="filter-tag"
            :class="{ active: typeFilter === opt.value }"
          >
            {{ opt.label }}
          </button>
        </div>
        <!-- 状态筛选 -->
        <div class="flex gap-2 flex-wrap">
          <button
            v-for="opt in statusOptions"
            :key="opt.value"
            @click="setStatusFilter(opt.value)"
            class="filter-tag"
            :class="{ active: statusFilter === opt.value }"
          >
            {{ opt.label }}
          </button>
        </div>
      </div>
    </div>

    <!-- 行程列表 -->
    <div class="card">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="trips.length === 0" class="empty-state">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17a2 2 0 11-4 0 2 2 0 014 0zM19 17a2 2 0 11-4 0 2 2 0 014 0z" />
        </svg>
        <p>暂无行程数据</p>
      </div>

      <div v-else>
        <div
          v-for="trip in trips"
          :key="trip.id"
          class="list-item"
          @click="showTripDetail(trip)"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span :class="trip.trip_type === 1 ? 'badge-info' : 'badge-success'" class="badge text-[11px]">
                {{ trip.trip_type === 1 ? '司机' : '乘客' }}
              </span>
              <span :class="getStatusBadgeClass(trip.status)" class="badge text-[11px]">
                {{ getStatusText(trip.status) }}
              </span>
            </div>
            <div class="flex items-center gap-1 mb-1">
              <span class="font-medium text-gray-800 text-sm">{{ trip.departure_city }}</span>
              <svg class="w-3 h-3 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
              <span class="font-medium text-gray-800 text-sm">{{ trip.destination_city }}</span>
            </div>
            <div class="text-xs text-gray-400">
              {{ trip.user?.nickname || '用户' + trip.user_id }} · {{ formatTime(trip.departure_time) }}
            </div>
          </div>
          <svg class="w-5 h-5 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>

      <!-- 分页 -->
      <div v-if="total > pageSize" class="p-4 border-t border-gray-100 flex items-center justify-between">
        <div class="text-xs text-gray-400">共 {{ total }} 条</div>
        <div class="pagination">
          <button :disabled="page === 1" @click="changePage(page - 1)">上一页</button>
          <span class="text-xs text-gray-500">{{ page }}/{{ totalPages }}</span>
          <button :disabled="page >= totalPages" @click="changePage(page + 1)">下一页</button>
        </div>
      </div>
    </div>

    <!-- 行程详情弹窗 -->
    <div v-if="showDetail" class="modal-overlay" @click.self="closeDetail">
      <div class="modal-content">
        <div class="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-800">行程详情</h3>
          <button @click="closeDetail" class="p-2 -mr-2 text-gray-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4" v-if="currentTrip">
          <!-- 路线信息 -->
          <div class="bg-gray-50 rounded-xl p-4 mb-4">
            <div class="flex items-center justify-center gap-3">
              <div class="text-center">
                <div class="text-lg font-bold text-gray-800">{{ currentTrip.departure_city }}</div>
                <div class="text-xs text-gray-400">出发</div>
              </div>
              <svg class="w-6 h-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3" />
              </svg>
              <div class="text-center">
                <div class="text-lg font-bold text-gray-800">{{ currentTrip.destination_city }}</div>
                <div class="text-xs text-gray-400">到达</div>
              </div>
            </div>
          </div>

          <!-- 详细信息 -->
          <div class="space-y-0">
            <div class="detail-row">
              <span class="detail-label">行程ID</span>
              <span class="detail-value">{{ currentTrip.id }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">发布者</span>
              <span class="detail-value">
                <button @click="viewPublisher(currentTrip)" class="text-blue-500 underline">
                  {{ currentTrip.user?.nickname || '用户' + currentTrip.user_id }}
                </button>
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">类型</span>
              <span class="detail-value">
                <span :class="currentTrip.trip_type === 1 ? 'badge-info' : 'badge-success'" class="badge">
                  {{ currentTrip.trip_type === 1 ? '司机' : '乘客' }}
                </span>
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">状态</span>
              <span class="detail-value">
                <span :class="getStatusBadgeClass(currentTrip.status)" class="badge">
                  {{ getStatusText(currentTrip.status) }}
                </span>
              </span>
            </div>
            <div class="detail-row">
              <span class="detail-label">出发时间</span>
              <span class="detail-value">{{ formatTime(currentTrip.departure_time) }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">发布时间</span>
              <span class="detail-value">{{ formatTime(currentTrip.created_at) }}</span>
            </div>
            <div v-if="currentTrip.remark" class="detail-row">
              <span class="detail-label">备注</span>
              <span class="detail-value">{{ currentTrip.remark }}</span>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="mt-6">
            <button
              v-if="currentTrip.status !== 5"
              @click="handleBan(currentTrip)"
              class="w-full btn btn-danger"
            >
              封禁行程
            </button>
            <button
              v-else
              @click="handleUnban(currentTrip)"
              class="w-full btn btn-success"
            >
              解除封禁
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()

const trips = ref([])
const loading = ref(true)
const search = ref('')
const typeFilter = ref('')
const statusFilter = ref('')
const filterUserId = ref('')
const filterUserName = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetail = ref(false)
const currentTrip = ref(null)

const typeOptions = [
  { label: '全部', value: '' },
  { label: '司机', value: '1' },
  { label: '乘客', value: '2' }
]

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待匹配', value: '1' },
  { label: '已匹配', value: '2' },
  { label: '已完成', value: '3' },
  { label: '已取消', value: '4' },
  { label: '已封禁', value: '5' }
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

onMounted(() => {
  if (route.query.user_id) {
    filterUserId.value = route.query.user_id
    filterUserName.value = route.query.user_name || route.query.user_id
  }
  if (route.query.status !== undefined) {
    statusFilter.value = route.query.status
  }
  fetchTrips()
})

watch(() => route.query, (newQuery) => {
  if (newQuery.user_id) {
    filterUserId.value = newQuery.user_id
    filterUserName.value = newQuery.user_name || newQuery.user_id
  } else {
    filterUserId.value = ''
    filterUserName.value = ''
  }
  if (newQuery.status !== undefined) {
    statusFilter.value = newQuery.status
  }
  page.value = 1
  fetchTrips()
})

function setTypeFilter(value) {
  typeFilter.value = value
  page.value = 1
  fetchTrips()
}

function setStatusFilter(value) {
  statusFilter.value = value
  page.value = 1
  fetchTrips()
}

async function fetchTrips() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (search.value) params.search = search.value
    if (typeFilter.value) params.trip_type = parseInt(typeFilter.value)
    if (statusFilter.value) params.status = parseInt(statusFilter.value)
    if (filterUserId.value) params.user_id = filterUserId.value

    const data = await api.get('/admin/trips', { params })
    trips.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.error('Failed to fetch trips:', e)
  } finally {
    loading.value = false
  }
}

function changePage(newPage) {
  page.value = newPage
  fetchTrips()
}

function showTripDetail(trip) {
  currentTrip.value = trip
  showDetail.value = true
}

function closeDetail() {
  showDetail.value = false
  currentTrip.value = null
}

function clearUserFilter() {
  router.push({ path: '/trips' })
}

function viewPublisher(trip) {
  closeDetail()
  router.push({ path: '/users', query: { highlight: trip.user_id } })
}

async function handleBan(trip) {
  if (!confirm('确定要封禁这个行程吗？')) return

  try {
    await api.post(`/admin/trips/${trip.id}/ban`)
    trip.status = 5
    alert('封禁成功')
  } catch (e) {
    alert('封禁失败: ' + (e.response?.data?.message || e.message))
  }
}

async function handleUnban(trip) {
  if (!confirm('确定要解封这个行程吗？')) return

  try {
    await api.post(`/admin/trips/${trip.id}/unban`)
    trip.status = 1
    alert('解封成功')
  } catch (e) {
    alert('解封失败: ' + (e.response?.data?.message || e.message))
  }
}

function formatTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}-${day} ${hour}:${minute}`
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
</script>
