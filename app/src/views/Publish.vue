<template>
  <div class="publish-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部 -->
    <div class="border-b safe-area-top" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-center h-12 px-4">
        <button @click="goBack" class="p-2 -ml-2">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : ''">发布行程</h1>
        <div class="w-10"></div>
      </div>
    </div>

    <!-- 表单 -->
    <div class="p-4">
      <!-- 当前身份提示 -->
      <div class="card p-4 mb-4">
        <div class="flex items-center justify-between">
          <span class="text-sm text-gray-600">当前身份</span>
          <span
            class="px-3 py-1 text-sm rounded-full font-medium"
            :class="form.trip_type === 1 ? 'bg-primary-100 text-primary-600' : 'bg-green-100 text-green-600'"
          >
            {{ form.trip_type === 1 ? '🚗 司机' : '🙋 乘客' }}
          </span>
        </div>
        <p class="text-xs text-gray-400 mt-2">如需切换身份，请在「我的」-「设置」中修改</p>
      </div>

      <!-- 出发地 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-medium text-gray-800 mb-3">出发地</h3>
        <div class="space-y-3">
          <input
            v-model="form.departure_city"
            type="text"
            placeholder="出发城市"
            class="input"
          />
          <input
            v-model="form.departure_address"
            type="text"
            placeholder="详细地址（可选）"
            class="input"
          />
        </div>
      </div>

      <!-- 目的地 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-medium text-gray-800 mb-3">目的地</h3>
        <div class="space-y-3">
          <input
            v-model="form.destination_city"
            type="text"
            placeholder="目的城市"
            class="input"
          />
          <input
            v-model="form.destination_address"
            type="text"
            placeholder="详细地址（可选）"
            class="input"
          />
        </div>
      </div>

      <!-- 出发时间 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-medium text-gray-800 mb-3">出发时间</h3>
        <input
          v-model="form.departure_time"
          type="datetime-local"
          class="input"
        />
      </div>

      <!-- 其他信息 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-medium text-gray-800 mb-3">其他信息</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs text-gray-500 mb-1">
              {{ form.trip_type === 1 ? '可载人数' : '需要座位数' }}
            </label>
            <input
              v-model.number="form.seats"
              type="number"
              min="1"
              max="7"
              placeholder="座位数"
              class="input"
            />
          </div>
          <div v-if="form.trip_type === 1">
            <label class="block text-xs text-gray-500 mb-1">费用（元/人）</label>
            <input
              v-model.number="form.price"
              type="number"
              min="0"
              placeholder="费用（可选）"
              class="input"
            />
          </div>
          <div>
            <label class="block text-xs text-gray-500 mb-1">备注</label>
            <textarea
              v-model="form.remark"
              placeholder="其他说明（可选）"
              rows="3"
              class="input resize-none"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- 图片上传 -->
      <div class="card p-4 mb-4">
        <h3 class="text-sm font-medium text-gray-800 mb-3">行程图片（可选）</h3>
        <p class="text-xs text-gray-500 mb-3">上传车辆或行程相关图片，增加可信度</p>
        <div class="flex flex-wrap gap-3">
          <!-- 已上传图片预览 -->
          <div 
            v-for="(img, index) in uploadedImages" 
            :key="index"
            class="relative w-20 h-20 rounded-lg overflow-hidden bg-gray-100"
          >
            <img :src="img" class="w-full h-full object-cover" />
            <button
              @click="removeImage(index)"
              class="absolute top-1 right-1 w-5 h-5 bg-black/50 rounded-full flex items-center justify-center"
            >
              <svg class="w-3 h-3 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <!-- 上传按钮 -->
          <label 
            v-if="uploadedImages.length < 3"
            class="w-20 h-20 border-2 border-dashed border-gray-300 rounded-lg flex flex-col items-center justify-center cursor-pointer hover:border-primary-400 transition-colors"
          >
            <svg class="w-6 h-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span class="text-xs text-gray-400 mt-1">上传</span>
            <input
              type="file"
              accept="image/*"
              class="hidden"
              @change="handleImageUpload"
              :disabled="uploading"
            />
          </label>
        </div>
        <p v-if="uploading" class="text-xs text-primary-500 mt-2">图片上传中...</p>
      </div>

      <!-- 提交按钮 -->
      <button
        @click="handleSubmit"
        :disabled="loading"
        class="btn btn-primary w-full py-4 text-base font-semibold min-h-[52px]"
      >
        <span v-if="loading" class="loading-spinner mr-2"></span>
        {{ loading ? '发布中...' : '发布行程' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useTripStore } from '@/stores/trip'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { uploadImage } from '@/utils/api'

const router = useRouter()
const tripStore = useTripStore()
const userStore = useUserStore()
const appStore = useAppStore()

const loading = ref(false)
const uploading = ref(false)
const uploadedImages = ref([])
const form = reactive({
  trip_type: userStore.identity === 1 ? 1 : 2, // 使用用户设置的身份
  departure_city: '',
  departure_address: '',
  destination_city: '',
  destination_address: '',
  departure_time: '',
  seats: 1,
  price: 0,
  remark: ''
})

onMounted(() => {
  // 设置默认时间为明天
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  tomorrow.setHours(8, 0, 0, 0)
  form.departure_time = formatDateTimeLocal(tomorrow)
})

function formatDateTimeLocal(date) {
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day}T${hour}:${minute}`
}

function goBack() {
  router.back()
}

async function handleImageUpload(e) {
  const file = e.target.files[0]
  if (!file) return
  
  // validate file type
  if (!file.type.startsWith('image/')) {
    appStore.showToast('请选择图片文件', 'error')
    return
  }
  
  // validate file size (max 5MB)
  if (file.size > 5 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过5MB', 'error')
    return
  }
  
  uploading.value = true
  try {
    // use 'trip' biz type for trip images, returns public URL
    const url = await uploadImage(file, 'trip')
    if (url) {
      uploadedImages.value.push(url)
    }
  } catch (err) {
    appStore.showToast(err.message || '图片上传失败', 'error')
  } finally {
    uploading.value = false
    // reset input
    e.target.value = ''
  }
}

function removeImage(index) {
  uploadedImages.value.splice(index, 1)
}

async function handleSubmit() {
  if (!form.departure_city) {
    appStore.showToast('请输入出发城市', 'error')
    return
  }
  if (!form.destination_city) {
    appStore.showToast('请输入目的城市', 'error')
    return
  }
  if (!form.departure_time) {
    appStore.showToast('请选择出发时间', 'error')
    return
  }
  if (form.seats < 1) {
    appStore.showToast('座位数至少为1', 'error')
    return
  }

  loading.value = true
  try {
    const submitData = {
      ...form,
      departure_time: form.departure_time.replace('T', ' '),
      departure_address: form.departure_address || form.departure_city,
      destination_address: form.destination_address || form.destination_city,
      departure_lat: 0,
      departure_lng: 0,
      destination_lat: 0,
      destination_lng: 0,
      images: uploadedImages.value.length > 0 ? JSON.stringify(uploadedImages.value) : ''
    }
    await tripStore.createTrip(submitData)
    appStore.showToast('发布成功', 'success')
    router.replace('/trips')
  } catch (e) {
    // error handled in interceptor
  } finally {
    loading.value = false
  }
}
</script>
