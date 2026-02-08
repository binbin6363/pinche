<template>
  <div class="users-page">
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
            placeholder="搜索手机号/昵称"
            @keyup.enter="fetchUsers"
          />
        </div>
        <!-- 筛选标签 -->
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

    <!-- 用户列表 -->
    <div class="card">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="users.length === 0" class="empty-state">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
        </svg>
        <p>暂无用户数据</p>
      </div>

      <div v-else>
        <div
          v-for="user in users"
          :key="user.id"
          class="list-item"
          @click="showUserDetail(user)"
        >
          <div class="w-11 h-11 bg-gray-100 rounded-full flex items-center justify-center mr-3 flex-shrink-0">
            <img v-if="user.avatar" :src="user.avatar" class="w-full h-full object-cover rounded-full" />
            <span v-else class="text-sm font-medium text-gray-500">{{ user.nickname?.charAt(0) || 'U' }}</span>
          </div>
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span class="font-medium text-gray-800 truncate">{{ user.nickname }}</span>
              <span :class="user.status === 0 ? 'badge-success' : 'badge-danger'" class="badge text-[11px]">
                {{ user.status === 0 ? '正常' : '封禁' }}
              </span>
            </div>
            <div class="text-xs text-gray-400">{{ user.phone }} · {{ user.city || '未知城市' }}</div>
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

    <!-- 用户详情弹窗 -->
    <div v-if="showDetail" class="modal-overlay" @click.self="closeDetail">
      <div class="modal-content">
        <div class="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-800">用户详情</h3>
          <button @click="closeDetail" class="p-2 -mr-2 text-gray-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4" v-if="currentUser">
          <!-- 用户头像和基本信息 -->
          <div class="flex items-center mb-6">
            <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mr-4">
              <img v-if="currentUser.avatar" :src="currentUser.avatar" class="w-full h-full object-cover rounded-full" />
              <span v-else class="text-xl font-medium text-gray-500">{{ currentUser.nickname?.charAt(0) || 'U' }}</span>
            </div>
            <div>
              <div class="flex items-center gap-2">
                <span class="text-lg font-semibold text-gray-800">{{ currentUser.nickname }}</span>
                <span :class="currentUser.status === 0 ? 'badge-success' : 'badge-danger'" class="badge">
                  {{ currentUser.status === 0 ? '正常' : '已封禁' }}
                </span>
              </div>
              <div class="text-sm text-gray-500 mt-1">ID: {{ currentUser.open_id }}</div>
            </div>
          </div>

          <!-- 详细信息 -->
          <div class="space-y-0">
            <div class="detail-row">
              <span class="detail-label">手机号</span>
              <span class="detail-value">{{ currentUser.phone }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">城市</span>
              <span class="detail-value">{{ currentUser.city || '-' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">省份</span>
              <span class="detail-value">{{ currentUser.province || '-' }}</span>
            </div>
            <div class="detail-row">
              <span class="detail-label">注册时间</span>
              <span class="detail-value">{{ formatTime(currentUser.created_at) }}</span>
            </div>
          </div>

          <!-- 查看行程 -->
          <button
            @click="viewUserTrips(currentUser)"
            class="w-full mt-4 btn btn-primary"
          >
            查看该用户行程
          </button>

          <!-- 操作按钮 -->
          <div class="mt-6">
            <button
              v-if="currentUser.status === 0"
              @click="handleBan(currentUser)"
              class="w-full btn btn-danger"
            >
              封禁用户
            </button>
            <button
              v-else
              @click="handleUnban(currentUser)"
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/utils/api'

const route = useRoute()
const router = useRouter()

const users = ref([])
const loading = ref(true)
const search = ref('')
const statusFilter = ref('')
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetail = ref(false)
const currentUser = ref(null)

const statusOptions = [
  { label: '全部', value: '' },
  { label: '正常', value: '0' },
  { label: '已封禁', value: '1' }
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

onMounted(async () => {
  if (route.query.status !== undefined) {
    statusFilter.value = route.query.status
  }
  await fetchUsers()
  if (route.query.highlight) {
    const target = users.value.find(u => u.open_id === route.query.highlight)
    if (target) showUserDetail(target)
  }
})

function setStatusFilter(value) {
  statusFilter.value = value
  page.value = 1
  fetchUsers()
}

async function fetchUsers() {
  loading.value = true
  try {
    const params = {
      page: page.value,
      page_size: pageSize.value
    }
    if (search.value) params.search = search.value
    if (statusFilter.value !== '') params.status = parseInt(statusFilter.value)

    const data = await api.get('/admin/users', { params })
    users.value = data.list || []
    total.value = data.total || 0
  } catch (e) {
    console.error('Failed to fetch users:', e)
  } finally {
    loading.value = false
  }
}

function changePage(newPage) {
  page.value = newPage
  fetchUsers()
}

function showUserDetail(user) {
  currentUser.value = user
  showDetail.value = true
}

function closeDetail() {
  showDetail.value = false
  currentUser.value = null
}

function viewUserTrips(user) {
  closeDetail()
  router.push({ path: '/trips', query: { user_id: user.open_id, user_name: user.nickname } })
}

async function handleBan(user) {
  if (!confirm(`确定要封禁用户 "${user.nickname}" 吗？`)) return

  try {
    // use open_id for API call
    await api.post(`/admin/users/${user.open_id}/ban`)
    user.status = 1
    alert('封禁成功')
  } catch (e) {
    alert('封禁失败: ' + (e.response?.data?.message || e.message))
  }
}

async function handleUnban(user) {
  if (!confirm(`确定要解封用户 "${user.nickname}" 吗？`)) return

  try {
    // use open_id for API call
    await api.post(`/admin/users/${user.open_id}/unban`)
    user.status = 0
    alert('解封成功')
  } catch (e) {
    alert('解封失败: ' + (e.response?.data?.message || e.message))
  }
}

function formatTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day} ${hour}:${minute}`
}
</script>
