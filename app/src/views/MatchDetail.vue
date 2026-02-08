<template>
  <div class="match-detail-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部 -->
    <div class="border-b safe-area-top" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-center h-12 px-4">
        <button @click="goBack" class="p-2 -ml-2">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : ''">匹配详情</h1>
        <div class="w-10"></div>
      </div>
    </div>

    <div v-if="loading" class="flex justify-center py-12">
      <div class="loading-spinner"></div>
    </div>

    <div v-else-if="!match" class="text-center py-12 text-gray-400">
      匹配记录不存在
    </div>

    <div v-else class="p-4">
      <!-- 状态卡片 -->
      <div class="card p-4 mb-4 text-center">
        <div
          class="inline-flex items-center justify-center w-16 h-16 rounded-full mb-3"
          :class="getStatusIconClass()"
        >
          <svg v-if="match.status === 1" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
          <svg v-else-if="match.status === 2" class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          <svg v-else class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <h2 class="text-lg font-semibold" :class="getStatusTextColor()">
          {{ getStatusTitle() }}
        </h2>
        <p class="text-sm text-gray-500 mt-1">{{ getStatusDesc() }}</p>
      </div>

      <!-- 双方用户信息 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-semibold text-gray-800 mb-3">双方信息</h3>
        <div class="flex items-center justify-around">
          <!-- 司机 -->
          <div class="text-center" @click="showUserProfile(match.driver)">
            <div class="w-14 h-14 mx-auto rounded-full flex items-center justify-center mb-2 bg-blue-100 overflow-hidden cursor-pointer">
              <img v-if="match.driver?.avatar" :src="match.driver.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-lg font-medium text-blue-500">{{ match.driver?.nickname?.charAt(0) || '司' }}</span>
            </div>
            <div class="text-sm font-medium">{{ getMaskedNickname(match.driver) }}</div>
            <div class="text-xs text-gray-400">司机</div>
          </div>
          
          <!-- 连接线 -->
          <div class="flex flex-col items-center">
            <svg class="w-8 h-8 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </div>
          
          <!-- 乘客 -->
          <div class="text-center" @click="showUserProfile(match.passenger)">
            <div class="w-14 h-14 mx-auto rounded-full flex items-center justify-center mb-2 bg-green-100 overflow-hidden cursor-pointer">
              <img v-if="match.passenger?.avatar" :src="match.passenger.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-lg font-medium text-green-500">{{ match.passenger?.nickname?.charAt(0) || '乘' }}</span>
            </div>
            <div class="text-sm font-medium">{{ getMaskedNickname(match.passenger) }}</div>
            <div class="text-xs text-gray-400">乘客</div>
          </div>
        </div>
      </div>

      <!-- 行程信息 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-semibold text-gray-800 mb-3">行程信息</h3>
        <div class="flex items-center gap-3">
          <div class="flex-1">
            <div class="text-xs text-gray-400">出发</div>
            <div class="text-sm font-medium">{{ tripInfo?.departure_city || '-' }}</div>
          </div>
          <svg class="w-6 h-6 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
          <div class="flex-1 text-right">
            <div class="text-xs text-gray-400">到达</div>
            <div class="text-sm font-medium">{{ tripInfo?.destination_city || '-' }}</div>
          </div>
        </div>
        <div class="mt-3 pt-3 border-t border-gray-100 text-sm text-gray-500 space-y-1">
          <div>出发时间：{{ formatTime(tripInfo?.departure_time) }}</div>
          <div>座位数：{{ tripInfo?.seats || '-' }}</div>
          <div v-if="tripInfo?.price > 0">费用：¥{{ tripInfo?.price }}/人</div>
        </div>
      </div>

      <!-- 双方确认状态 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-semibold text-gray-800 mb-3">确认状态</h3>
        <div class="flex items-center justify-around">
          <div class="text-center">
            <div class="w-12 h-12 mx-auto rounded-full flex items-center justify-center mb-2"
                 :class="match.driver_status === 1 ? 'bg-green-100' : 'bg-gray-100'">
              <svg v-if="match.driver_status === 1" class="w-6 h-6 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else-if="match.driver_status === 2" class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              <span v-else class="text-gray-400">?</span>
            </div>
            <div class="text-sm font-medium">司机</div>
            <div class="text-xs text-gray-500">{{ getConfirmText(match.driver_status) }}</div>
          </div>
          <div class="text-center">
            <div class="w-12 h-12 mx-auto rounded-full flex items-center justify-center mb-2"
                 :class="match.passenger_status === 1 ? 'bg-green-100' : 'bg-gray-100'">
              <svg v-if="match.passenger_status === 1" class="w-6 h-6 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else-if="match.passenger_status === 2" class="w-6 h-6 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              <span v-else class="text-gray-400">?</span>
            </div>
            <div class="text-sm font-medium">乘客</div>
            <div class="text-xs text-gray-500">{{ getConfirmText(match.passenger_status) }}</div>
          </div>
        </div>
      </div>

      <!-- 联系信息（仅匹配成功后显示） -->
      <div v-if="match.status === 1 && contactInfo" class="card p-4 mb-4">
        <h3 class="text-sm font-semibold text-gray-800 mb-3">联系方式</h3>
        <div class="space-y-3">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <div class="text-sm font-medium">司机：{{ contactInfo.driver_nickname }}</div>
              <div class="text-primary-500">{{ contactInfo.driver_phone }}</div>
            </div>
            <a :href="'tel:' + contactInfo.driver_phone" class="btn btn-primary btn-sm">
              拨打
            </a>
          </div>
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div>
              <div class="text-sm font-medium">乘客：{{ contactInfo.passenger_nickname }}</div>
              <div class="text-primary-500">{{ contactInfo.passenger_phone }}</div>
            </div>
            <a :href="'tel:' + contactInfo.passenger_phone" class="btn btn-primary btn-sm">
              拨打
            </a>
          </div>
        </div>
        
        <!-- 免责声明 -->
        <div class="mt-4 p-3 bg-yellow-50 rounded-lg">
          <p class="text-xs text-yellow-700">
            <strong>免责声明：</strong>本平台仅提供信息匹配服务，不参与任何实际交易。双方请自行协商出行细节，注意人身和财产安全。如发生任何纠纷，请自行协商解决或寻求法律途径，本平台不承担任何责任。
          </p>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div v-if="match.status === 0 && needConfirm" class="flex gap-3">
        <button
          @click="handleConfirm(false)"
          :disabled="confirming"
          class="btn btn-secondary flex-1"
        >
          拒绝
        </button>
        <button
          @click="handleConfirm(true)"
          :disabled="confirming"
          class="btn btn-primary flex-1"
        >
          <span v-if="confirming" class="loading-spinner mr-2"></span>
          接受匹配
        </button>
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
import { useMatchStore } from '@/stores/match'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const route = useRoute()
const router = useRouter()
const matchStore = useMatchStore()
const userStore = useUserStore()
const appStore = useAppStore()

const match = ref(null)
const contactInfo = ref(null)
const loading = ref(true)
const confirming = ref(false)
const showUserModal = ref(false)
const selectedUser = ref(null)

const isDriver = computed(() => match.value?.driver_open_id === userStore.user?.open_id)
const myStatus = computed(() => isDriver.value ? match.value?.driver_status : match.value?.passenger_status)
const needConfirm = computed(() => myStatus.value === 0)

// get trip info from driver_trip or passenger_trip
const tripInfo = computed(() => {
  if (!match.value) return null
  return match.value.driver_trip || match.value.passenger_trip || null
})

onMounted(async () => {
  try {
    const id = route.params.id
    match.value = await matchStore.getMatchById(id)
    if (match.value?.status === 1) {
      contactInfo.value = await matchStore.getContactInfo(id)
    }
  } finally {
    loading.value = false
  }
})

function goBack() {
  router.back()
}

function getStatusIconClass() {
  switch (match.value?.status) {
    case 1: return 'bg-green-100 text-green-500'
    case 2: return 'bg-gray-100 text-gray-400'
    default: return 'bg-yellow-100 text-yellow-500'
  }
}

function getStatusTextColor() {
  switch (match.value?.status) {
    case 1: return 'text-green-600'
    case 2: return 'text-gray-500'
    default: return 'text-yellow-600'
  }
}

function getStatusTitle() {
  if (match.value?.status === 1) return '拼车成功'
  if (match.value?.status === 2) return '匹配失败'
  return '等待确认'
}

function getStatusDesc() {
  if (match.value?.status === 1) return '双方已确认，可以联系对方确认出行细节'
  if (match.value?.status === 2) return '一方已拒绝，此匹配已结束'
  if (needConfirm.value) return '请确认是否接受此次匹配'
  return '等待对方确认中...'
}

function getConfirmText(status) {
  switch (status) {
    case 0: return '待确认'
    case 1: return '已接受'
    case 2: return '已拒绝'
    default: return '未知'
  }
}

function formatTime(time) {
  if (!time) return '-'
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hour}:${minute}`
}

function getMaskedNickname(user) {
  if (user && user.nickname) {
    const nickname = user.nickname
    if (nickname.length <= 2) {
      return nickname[0] + '**'
    }
    return nickname.slice(0, -2) + '**'
  }
  return '用户**'
}

function showUserProfile(user) {
  if (user) {
    selectedUser.value = user
    showUserModal.value = true
  }
}

async function handleConfirm(accept) {
  confirming.value = true
  try {
    await matchStore.confirmMatch(match.value.id, accept)
    appStore.showToast(accept ? '已接受匹配' : '已拒绝匹配', 'success')
    // 刷新数据
    match.value = await matchStore.getMatchById(match.value.id)
    if (match.value?.status === 1) {
      contactInfo.value = await matchStore.getContactInfo(match.value.id)
    }
  } finally {
    confirming.value = false
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
