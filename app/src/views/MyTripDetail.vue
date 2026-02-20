<template>
  <div class="my-trip-detail min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部导航 -->
    <div class="page-header safe-area-top sticky top-0 z-10">
      <div class="page-header-bg"></div>
      <div class="relative flex items-center h-12 px-4">
        <button @click="goBack" class="w-10 h-10 -ml-2 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-white">我的行程</h1>
        <!-- 分享按钮 -->
        <button @click="handleShare" class="w-10 h-10 -mr-2 flex items-center justify-center">
          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
        </button>
      </div>
    </div>

    <!-- 加载状态 -->
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="loading-spinner"></div>
    </div>

    <!-- 错误状态 -->
    <div v-else-if="error" class="text-center py-20">
      <p class="text-gray-500">{{ error }}</p>
      <button @click="loadTrip" class="btn btn-primary mt-4">重试</button>
    </div>

    <!-- 行程详情 -->
    <div v-else-if="trip" class="pb-24">
      <!-- 图片区域 -->
      <div class="relative h-48">
        <div class="absolute inset-0 bg-gradient-to-br" :style="{ background: `linear-gradient(135deg, var(--theme-primary), #F97316)` }"></div>
        <template v-if="tripImages.length > 0">
          <img :src="tripImages[currentImageIndex]" class="w-full h-full object-cover" />
          <div v-if="tripImages.length > 1" class="absolute bottom-3 left-1/2 -translate-x-1/2 flex gap-1.5">
            <span
              v-for="(_, idx) in tripImages"
              :key="idx"
              @click="currentImageIndex = idx"
              class="w-2 h-2 rounded-full cursor-pointer transition-all"
              :class="idx === currentImageIndex ? 'bg-white w-4' : 'bg-white/50'"
            ></span>
          </div>
        </template>
        <div v-else class="absolute inset-0 flex items-center justify-center">
          <svg class="w-16 h-16 text-white/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01m-7.08-7.071c3.904-3.905 10.236-3.905 14.14 0M1.394 9.393c5.857-5.857 15.355-5.857 21.213 0" />
          </svg>
        </div>
        <!-- 状态标签 -->
        <div class="absolute top-3 left-3">
          <span
            class="badge"
            :class="getStatusClass(trip.status)"
          >
            {{ getStatusText(trip.status) }}
          </span>
        </div>
        <!-- 修改图片按钮 -->
        <label class="absolute top-3 right-3 px-3 py-1.5 bg-black/30 rounded-full flex items-center gap-1 cursor-pointer">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <span class="text-xs text-white">修改</span>
          <input type="file" accept="image/*" class="hidden" @change="handleImageUpload" :disabled="uploading" />
        </label>
      </div>

      <!-- 数据统计 -->
      <div class="px-4 -mt-6 relative z-10">
        <div class="card p-4">
          <div class="flex items-center justify-around text-center">
            <div>
              <div class="text-2xl font-bold" :style="{ color: 'var(--theme-primary)' }">{{ trip.view_count || 0 }}</div>
              <div class="text-xs text-gray-500">浏览次数</div>
            </div>
            <div class="w-px h-10" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"></div>
            <div>
              <div class="text-2xl font-bold text-green-500">{{ trip.grabbers?.length || 0 }}</div>
              <div class="text-xs text-gray-500">抢单人数</div>
            </div>
            <div class="w-px h-10" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"></div>
            <div>
              <div class="text-2xl font-bold">{{ trip.seats }}</div>
              <div class="text-xs text-gray-500">座位数</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 行程信息 -->
      <div class="px-4 mt-4">
        <div class="card p-4">
          <div class="flex items-center gap-2 mb-3">
            <span
              class="badge"
              :class="trip.trip_type === 1 ? 'trip-type-driver' : 'trip-type-passenger'"
            >
              {{ trip.trip_type === 1 ? '车找人' : '人找车' }}
            </span>
            <span class="text-xs text-gray-400">发布于 {{ formatTime(trip.created_at) }}</span>
          </div>

          <!-- 起终点 -->
          <div class="space-y-3">
            <div class="flex items-start gap-3">
              <div class="w-6 flex flex-col items-center pt-1">
                <span class="w-2.5 h-2.5 rounded-full" style="background: var(--theme-primary);"></span>
                <span class="w-px h-8" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"></span>
              </div>
              <div class="flex-1">
                <div class="text-sm text-gray-500">出发地</div>
                <div class="font-medium">{{ trip.departure_city }}</div>
                <div class="text-xs text-gray-400">{{ trip.departure_address }}</div>
              </div>
            </div>
            <div class="flex items-start gap-3">
              <div class="w-6 flex flex-col items-center pt-1">
                <span class="w-2.5 h-2.5 rounded-full bg-orange-500"></span>
              </div>
              <div class="flex-1">
                <div class="text-sm text-gray-500">目的地</div>
                <div class="font-medium">{{ trip.destination_city }}</div>
                <div class="text-xs text-gray-400">{{ trip.destination_address }}</div>
              </div>
            </div>
          </div>

          <!-- 出发时间 -->
          <div class="mt-4 pt-4 border-t" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-500">出发时间</span>
              <span class="font-medium">{{ formatTime(trip.departure_time) }}</span>
            </div>
            <div v-if="trip.price > 0" class="flex items-center justify-between mt-2">
              <span class="text-sm text-gray-500">费用</span>
              <span class="font-medium" :style="{ color: 'var(--theme-primary)' }">¥{{ trip.price }}/人</span>
            </div>
          </div>

          <!-- 备注 -->
          <div v-if="trip.remark" class="mt-4 pt-4 border-t" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
            <div class="text-sm text-gray-500 mb-1">备注</div>
            <div class="text-sm">{{ trip.remark }}</div>
          </div>
        </div>
      </div>

      <!-- 抢单用户列表 -->
      <div class="px-4 mt-4">
        <div class="section-header px-0">
          <h3 class="section-title">
            <span>👥</span>
            抢单用户
          </h3>
          <span class="text-xs text-gray-400">{{ trip.grabbers?.length || 0 }}人</span>
        </div>
        
        <div v-if="!trip.grabbers || trip.grabbers.length === 0" class="card p-6 text-center">
          <div class="empty-state-icon mx-auto mb-3" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
            <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <p class="text-gray-400 text-sm">暂无抢单用户</p>
          <p class="text-gray-400 text-xs mt-1">等待其他用户发现您的行程</p>
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="grab in trip.grabbers"
            :key="grab.id"
            class="card p-4"
          >
            <div class="flex items-center gap-3">
              <div class="w-12 h-12 rounded-full flex items-center justify-center overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-100'">
                <img v-if="grab.user?.avatar" :src="grab.user.avatar" class="w-full h-full object-cover" />
                <span v-else class="text-lg text-gray-400">{{ grab.user?.nickname?.charAt(0) || '?' }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <div class="font-medium text-sm truncate">{{ grab.user?.nickname || '用户' }}</div>
                <div v-if="grab.message" class="text-xs text-gray-500 truncate">留言: {{ grab.message }}</div>
                <div class="text-xs text-gray-400">{{ formatShortTime(grab.created_at) }}</div>
              </div>
              <button
                @click="goToChat(grab)"
                class="btn btn-outline px-4 py-2 text-sm"
              >
                联系TA
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- 修改行程 -->
      <div class="px-4 mt-4">
        <button
          @click="showEditModal = true"
          class="btn btn-secondary w-full"
        >
          修改行程信息
        </button>
      </div>
    </div>

    <!-- 底部操作栏 -->
    <div v-if="trip && trip.status === 1" class="fixed bottom-0 left-0 right-0 p-4 safe-area-bottom border-t" :class="appStore.theme === 'dark' ? 'bg-gray-900 border-gray-800' : 'bg-white border-gray-100'">
      <div class="flex gap-3">
        <button
          @click="handleComplete"
          class="flex-1 btn text-white bg-green-500 hover:bg-green-600"
        >
          标记为已成行
        </button>
        <button
          @click="handleCancel"
          class="flex-1 btn btn-secondary"
        >
          取消行程
        </button>
      </div>
    </div>

    <!-- 编辑弹窗 -->
    <div v-if="showEditModal" class="action-sheet" @click.self="showEditModal = false">
      <div class="action-sheet-overlay" @click="showEditModal = false"></div>
      <div class="action-sheet-content max-h-[80vh] overflow-y-auto">
        <div class="sticky top-0 px-4 py-3 flex items-center justify-between border-b" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
          <span class="font-medium">修改行程</span>
          <button @click="showEditModal = false" class="p-2">
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-4 space-y-4">
          <!-- 提示 -->
          <div class="p-3 rounded-lg text-xs" :class="appStore.theme === 'dark' ? 'bg-yellow-900/30 text-yellow-400' : 'bg-yellow-50 text-yellow-700'">
            <p class="font-medium mb-1">⚠️ 修改说明</p>
            <p>• 图片、备注、座位、费用可直接修改</p>
            <p>• 起终点、出发时间修改需系统审核</p>
          </div>

          <!-- 图片 -->
          <div>
            <label class="block text-sm text-gray-600 mb-2">行程图片</label>
            <div class="flex flex-wrap gap-2">
              <div v-for="(img, idx) in editForm.images" :key="idx" class="relative w-16 h-16 rounded-lg overflow-hidden">
                <img :src="img" class="w-full h-full object-cover" />
                <button @click="removeEditImage(idx)" class="absolute top-0.5 right-0.5 w-5 h-5 bg-black/50 rounded-full flex items-center justify-center">
                  <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              <label v-if="editForm.images.length < 3" class="w-16 h-16 border-2 border-dashed rounded-lg flex items-center justify-center cursor-pointer" :class="appStore.theme === 'dark' ? 'border-gray-600' : 'border-gray-300'">
                <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                <input type="file" accept="image/*" class="hidden" @change="handleEditImageUpload" :disabled="uploading" />
              </label>
            </div>
          </div>

          <!-- 备注 -->
          <div>
            <label class="block text-sm text-gray-600 mb-2">备注</label>
            <textarea v-model="editForm.remark" rows="2" class="input resize-none" placeholder="其他说明"></textarea>
          </div>

          <!-- 座位数 -->
          <div>
            <label class="block text-sm text-gray-600 mb-2">座位数</label>
            <input v-model.number="editForm.seats" type="number" min="1" max="7" class="input" />
          </div>

          <!-- 费用 -->
          <div v-if="trip?.trip_type === 1">
            <label class="block text-sm text-gray-600 mb-2">费用（元/人）</label>
            <input v-model.number="editForm.price" type="number" min="0" class="input" />
          </div>

          <div class="border-t pt-4" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
            <p class="text-xs text-gray-500 mb-3">以下修改需要审核：</p>
            
            <!-- 出发城市 -->
            <div class="mb-3">
              <label class="block text-sm text-gray-600 mb-2">出发城市</label>
              <input v-model="editForm.departure_city" type="text" class="input" placeholder="出发城市" />
            </div>

            <!-- 目的城市 -->
            <div class="mb-3">
              <label class="block text-sm text-gray-600 mb-2">目的城市</label>
              <input v-model="editForm.destination_city" type="text" class="input" placeholder="目的城市" />
            </div>

            <!-- 出发时间 -->
            <div>
              <label class="block text-sm text-gray-600 mb-2">出发时间</label>
              <input v-model="editForm.departure_time" type="datetime-local" class="input" />
            </div>
          </div>

          <!-- 提交按钮 -->
          <button
            @click="handleSubmitEdit"
            :disabled="saving"
            class="btn btn-primary w-full"
          >
            {{ saving ? '保存中...' : '保存修改' }}
          </button>
        </div>
      </div>
    </div>

    <!-- 分享弹窗 -->
    <div v-if="showShareModal" class="action-sheet" @click.self="showShareModal = false">
      <div class="action-sheet-overlay" @click="showShareModal = false"></div>
      <div class="action-sheet-content">
        <div class="px-4 py-4 border-b" :class="appStore.theme === 'dark' ? 'border-gray-700' : 'border-gray-100'">
          <h3 class="text-lg font-semibold text-center">分享行程</h3>
        </div>
        
        <div class="p-4">
          <!-- 分享预览 -->
          <div class="card p-4 mb-4" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-50'">
            <div class="text-sm font-medium mb-2">
              {{ trip?.trip_type === 1 ? '🚗 车找人' : '🙋 人找车' }}
            </div>
            <div class="text-sm">
              📍 {{ trip?.departure_city }} → {{ trip?.destination_city }}
            </div>
            <div class="text-sm mt-1">
              🕐 {{ formatTime(trip?.departure_time) }}
            </div>
            <div class="text-sm mt-1">
              💺 {{ trip?.seats }}座{{ trip?.price > 0 ? ` · ¥${trip?.price}/人` : '' }}
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
import { useRouter, useRoute } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useAppStore } from '@/stores/app'
import { uploadImage } from '@/utils/api'

const router = useRouter()
const route = useRoute()
const tripStore = useTripStore()
const appStore = useAppStore()

const trip = ref(null)
const loading = ref(true)
const error = ref('')
const uploading = ref(false)
const saving = ref(false)
const currentImageIndex = ref(0)
const showEditModal = ref(false)
const showShareModal = ref(false)

const editForm = ref({
  images: [],
  remark: '',
  seats: 1,
  price: 0,
  departure_city: '',
  destination_city: '',
  departure_time: ''
})

const tripImages = computed(() => {
  if (!trip.value?.images) return []
  try {
    return JSON.parse(trip.value.images)
  } catch {
    return []
  }
})

onMounted(() => {
  loadTrip()
})

async function loadTrip() {
  loading.value = true
  error.value = ''
  try {
    const id = route.params.id
    trip.value = await tripStore.getMyTripDetail(id)
    initEditForm()
  } catch (e) {
    error.value = e.message || '加载失败'
  } finally {
    loading.value = false
  }
}

function initEditForm() {
  if (!trip.value) return
  editForm.value = {
    images: tripImages.value,
    remark: trip.value.remark || '',
    seats: trip.value.seats || 1,
    price: trip.value.price || 0,
    departure_city: '',
    destination_city: '',
    departure_time: ''
  }
}

function goBack() {
  router.back()
}

function goToChat(grab) {
  if (!grab?.user?.open_id) return
  router.push({
    path: `/chat/${grab.user.open_id}`,
    query: {
      nickname: grab.user.nickname || '用户',
      avatar: grab.user.avatar || ''
    }
  })
}

function getStatusClass(status) {
  switch (status) {
    case 1: return 'status-pending'
    case 2: return 'status-matched'
    case 3: return 'status-completed'
    case 4: return 'status-cancelled'
    case 5: return 'status-cancelled'
    default: return 'badge-secondary'
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

function formatTime(time) {
  if (!time) return ''
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}月${day}日 ${hour}:${minute}`
}

function formatShortTime(time) {
  if (!time) return ''
  const date = new Date(time)
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${month}/${day} ${hour}:${minute}`
}

async function handleImageUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  
  if (file.size > 5 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过5MB', 'error')
    return
  }
  
  uploading.value = true
  try {
    const url = await uploadImage(file, 'trip')
    if (url) {
      const images = [...tripImages.value, url]
      if (images.length > 3) images.shift()
      await tripStore.updateTrip(trip.value.id, { images: JSON.stringify(images) })
      await loadTrip()
      appStore.showToast('图片更新成功', 'success')
    }
  } catch (err) {
    appStore.showToast(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

async function handleEditImageUpload(e) {
  const file = e.target.files?.[0]
  if (!file) return
  
  if (file.size > 5 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过5MB', 'error')
    return
  }
  
  uploading.value = true
  try {
    const url = await uploadImage(file, 'trip')
    if (url) {
      editForm.value.images.push(url)
    }
  } catch (err) {
    appStore.showToast(err.message || '上传失败', 'error')
  } finally {
    uploading.value = false
    e.target.value = ''
  }
}

function removeEditImage(index) {
  editForm.value.images.splice(index, 1)
}

async function handleSubmitEdit() {
  saving.value = true
  try {
    const data = {
      images: JSON.stringify(editForm.value.images),
      remark: editForm.value.remark,
      seats: editForm.value.seats,
      price: editForm.value.price
    }
    
    if (editForm.value.departure_city) {
      data.departure_city = editForm.value.departure_city
    }
    if (editForm.value.destination_city) {
      data.destination_city = editForm.value.destination_city
    }
    if (editForm.value.departure_time) {
      data.departure_time = editForm.value.departure_time.replace('T', ' ')
    }
    
    const result = await tripStore.updateTrip(trip.value.id, data)
    showEditModal.value = false
    
    if (result.needs_review) {
      appStore.showToast(result.message, 'info')
    } else {
      appStore.showToast('修改成功', 'success')
    }
    
    await loadTrip()
  } catch (err) {
    appStore.showToast(err.message || '修改失败', 'error')
  } finally {
    saving.value = false
  }
}

async function handleComplete() {
  if (!confirm('确定要标记这个行程为已成行吗？')) return
  try {
    await tripStore.completeTrip(trip.value.id)
    appStore.showToast('行程已标记为已成行', 'success')
    await loadTrip()
  } catch (e) {
    // error handled
  }
}

async function handleCancel() {
  if (!confirm('确定要取消这个行程吗？')) return
  try {
    await tripStore.cancelTrip(trip.value.id)
    appStore.showToast('行程已取消', 'success')
    await loadTrip()
  } catch (e) {
    // error handled
  }
}

function handleShare() {
  showShareModal.value = true
}

function getShareText() {
  if (!trip.value) return ''
  return `🚗 ${trip.value.trip_type === 1 ? '车找人' : '人找车'}\n📍 ${trip.value.departure_city} → ${trip.value.destination_city}\n🕐 ${formatTime(trip.value.departure_time)}\n💺 ${trip.value.seats}座${trip.value.price > 0 ? ` · ¥${trip.value.price}/人` : ''}`
}

function getShareUrl() {
  if (!trip.value) return ''
  return `${window.location.origin}/trip/${trip.value.id}`
}

async function shareToWechat() {
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
  copyShareText()
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
</script>
