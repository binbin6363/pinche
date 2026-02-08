<template>
  <div class="profile-edit min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-50'">
    <!-- 顶部导航 -->
    <div class="page-header safe-area-top">
      <div class="page-header-bg"></div>
      <div class="relative flex items-center h-12 px-4">
        <button @click="goBack" class="w-10 h-10 -ml-2 flex items-center justify-center">
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold text-white">编辑资料</h1>
        <button 
          @click="handleSave" 
          :disabled="saving"
          class="px-3 py-1.5 text-sm font-medium text-white bg-white/20 rounded-full hover:bg-white/30 transition-all disabled:opacity-50"
        >
          {{ saving ? '保存中...' : '保存' }}
        </button>
      </div>
    </div>

    <!-- 表单内容 -->
    <div class="px-4 py-4 space-y-4">
      <!-- 头像 -->
      <div class="card p-4">
        <div class="flex items-center justify-between">
          <span class="text-sm font-medium">头像</span>
          <div class="flex items-center gap-3">
            <div 
              @click="triggerAvatarUpload"
              class="w-16 h-16 rounded-full bg-gray-100 flex items-center justify-center overflow-hidden cursor-pointer hover:ring-2 hover:ring-offset-2 transition-all"
              :class="appStore.theme === 'dark' ? 'bg-gray-700 hover:ring-gray-500' : 'hover:ring-gray-300'"
            >
              <img v-if="form.avatar" :src="form.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-2xl text-gray-400">{{ form.nickname?.charAt(0) || '?' }}</span>
            </div>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </div>
        <input
          ref="avatarInput"
          type="file"
          accept="image/*"
          class="hidden"
          @change="handleAvatarSelect"
        />
      </div>

      <!-- 基本信息 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">基本信息</span>
        </div>
        
        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">昵称</span>
          <input
            v-model="form.nickname"
            type="text"
            placeholder="请输入昵称"
            class="flex-1 text-sm text-right outline-none bg-transparent"
            maxlength="20"
          />
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">性别</span>
          <div class="flex-1 flex justify-end gap-2">
            <button
              @click="form.gender = 1"
              class="px-4 py-1.5 text-sm rounded-full transition-all"
              :class="form.gender === 1 
                ? 'bg-blue-500 text-white' 
                : appStore.theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-600'"
            >
              男
            </button>
            <button
              @click="form.gender = 2"
              class="px-4 py-1.5 text-sm rounded-full transition-all"
              :class="form.gender === 2 
                ? 'bg-pink-500 text-white' 
                : appStore.theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-600'"
            >
              女
            </button>
          </div>
        </div>

        <div class="list-item">
          <span class="text-sm text-gray-500 w-24">所在地</span>
          <input
            v-model="form.city"
            type="text"
            placeholder="请输入城市"
            class="flex-1 text-sm text-right outline-none bg-transparent"
          />
        </div>
      </div>

      <!-- 联系方式 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">联系方式</span>
          <p class="text-xs text-gray-400 mt-1">仅向已匹配用户展示</p>
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">手机号</span>
          <input
            v-model="form.contact_phone"
            type="tel"
            placeholder="请输入手机号"
            class="flex-1 text-sm text-right outline-none bg-transparent"
            maxlength="11"
          />
        </div>

        <div class="list-item">
          <span class="text-sm text-gray-500 w-24">微信号</span>
          <input
            v-model="form.contact_wechat"
            type="text"
            placeholder="请输入微信号"
            class="flex-1 text-sm text-right outline-none bg-transparent"
          />
        </div>
      </div>

      <!-- 紧急联系人 -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">紧急联系人</span>
          <p class="text-xs text-gray-400 mt-1">行程出现异常时通知的联系人</p>
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">联系人姓名</span>
          <input
            v-model="form.emergency_contact_name"
            type="text"
            placeholder="请输入姓名"
            class="flex-1 text-sm text-right outline-none bg-transparent"
          />
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">联系人电话</span>
          <input
            v-model="form.emergency_contact_phone"
            type="tel"
            placeholder="请输入电话"
            class="flex-1 text-sm text-right outline-none bg-transparent"
            maxlength="11"
          />
        </div>

        <div class="list-item">
          <span class="text-sm text-gray-500 w-24">与我关系</span>
          <div class="flex-1 flex justify-end gap-2 flex-wrap">
            <button
              v-for="rel in relationOptions"
              :key="rel"
              @click="form.emergency_contact_relation = rel"
              class="px-3 py-1 text-xs rounded-full transition-all"
              :class="form.emergency_contact_relation === rel
                ? 'bg-blue-500 text-white'
                : appStore.theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-600'"
            >
              {{ rel }}
            </button>
          </div>
        </div>
      </div>

      <!-- 车辆信息（司机） -->
      <div class="card">
        <div class="px-4 py-3 border-b border-gray-100">
          <span class="text-sm font-semibold">车辆信息</span>
          <p class="text-xs text-gray-400 mt-1">司机发布行程时展示</p>
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">车牌号</span>
          <input
            v-model="form.car_number"
            type="text"
            placeholder="如: 粤B12345"
            class="flex-1 text-sm text-right outline-none bg-transparent uppercase"
            maxlength="8"
          />
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">车辆品牌</span>
          <input
            v-model="form.car_brand"
            type="text"
            placeholder="如: 大众"
            class="flex-1 text-sm text-right outline-none bg-transparent"
          />
        </div>

        <div class="list-item border-b border-gray-100">
          <span class="text-sm text-gray-500 w-24">车辆型号</span>
          <input
            v-model="form.car_model"
            type="text"
            placeholder="如: 帕萨特"
            class="flex-1 text-sm text-right outline-none bg-transparent"
          />
        </div>

        <div class="list-item">
          <span class="text-sm text-gray-500 w-24">车辆颜色</span>
          <div class="flex-1 flex justify-end gap-2 flex-wrap">
            <button
              v-for="color in carColorOptions"
              :key="color.value"
              @click="form.car_color = color.value"
              class="px-3 py-1 text-xs rounded-full transition-all flex items-center gap-1"
              :class="form.car_color === color.value
                ? 'bg-blue-500 text-white'
                : appStore.theme === 'dark' ? 'bg-gray-700 text-gray-300' : 'bg-gray-100 text-gray-600'"
            >
              <span 
                class="w-3 h-3 rounded-full border border-white/30"
                :style="{ backgroundColor: color.hex }"
              ></span>
              {{ color.label }}
            </button>
          </div>
        </div>
      </div>

      <!-- 资料完整度 -->
      <div class="card p-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium">资料完整度</span>
          <span 
            class="text-sm font-semibold"
            :class="completeness >= 80 ? 'text-green-500' : completeness >= 50 ? 'text-yellow-500' : 'text-red-500'"
          >
            {{ completeness }}%
          </span>
        </div>
        <div class="h-2 bg-gray-200 rounded-full overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
          <div 
            class="h-full rounded-full transition-all duration-500"
            :class="completeness >= 80 ? 'bg-green-500' : completeness >= 50 ? 'bg-yellow-500' : 'bg-red-500'"
            :style="{ width: `${completeness}%` }"
          ></div>
        </div>
        <p class="text-xs text-gray-400 mt-2">
          完善资料可提高匹配成功率
        </p>
      </div>
    </div>

    <!-- Image Cropper Modal -->
    <ImageCropper
      v-model:visible="showCropper"
      :image-src="cropperImageSrc"
      @confirm="handleCropConfirm"
      @cancel="handleCropCancel"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'
import { uploadImage } from '@/utils/api'
import ImageCropper from '@/components/ImageCropper.vue'

const router = useRouter()
const userStore = useUserStore()
const appStore = useAppStore()

const saving = ref(false)
const avatarInput = ref(null)
const showCropper = ref(false)
const cropperImageSrc = ref('')

const relationOptions = ['父母', '配偶', '子女', '兄弟姐妹', '朋友', '其他']

const carColorOptions = [
  { value: 'white', label: '白色', hex: '#FFFFFF' },
  { value: 'black', label: '黑色', hex: '#1F2937' },
  { value: 'silver', label: '银色', hex: '#9CA3AF' },
  { value: 'gray', label: '灰色', hex: '#6B7280' },
  { value: 'red', label: '红色', hex: '#EF4444' },
  { value: 'blue', label: '蓝色', hex: '#3B82F6' },
  { value: 'green', label: '绿色', hex: '#10B981' },
  { value: 'brown', label: '棕色', hex: '#92400E' }
]

const form = reactive({
  avatar: '',
  nickname: '',
  gender: 0,
  city: '',
  contact_phone: '',
  contact_wechat: '',
  emergency_contact_name: '',
  emergency_contact_phone: '',
  emergency_contact_relation: '',
  car_number: '',
  car_brand: '',
  car_model: '',
  car_color: ''
})

const completeness = computed(() => {
  let score = 0
  if (form.nickname) score += 10
  if (form.avatar) score += 10
  if (form.gender) score += 5
  if (form.city) score += 5
  if (form.contact_phone) score += 15
  if (form.contact_wechat) score += 10
  if (form.emergency_contact_name && form.emergency_contact_phone) score += 15
  if (form.car_number) score += 15
  if (form.car_brand || form.car_model) score += 10
  if (form.car_color) score += 5
  return Math.min(100, score)
})

onMounted(() => {
  if (userStore.user) {
    Object.assign(form, {
      avatar: userStore.user.avatar || '',
      nickname: userStore.user.nickname || '',
      gender: userStore.user.gender || 0,
      city: userStore.user.city || '',
      contact_phone: userStore.user.contact_phone || '',
      contact_wechat: userStore.user.contact_wechat || '',
      emergency_contact_name: userStore.user.emergency_contact_name || '',
      emergency_contact_phone: userStore.user.emergency_contact_phone || '',
      emergency_contact_relation: userStore.user.emergency_contact_relation || '',
      car_number: userStore.user.car_number || '',
      car_brand: userStore.user.car_brand || '',
      car_model: userStore.user.car_model || '',
      car_color: userStore.user.car_color || ''
    })
  }
})

function goBack() {
  router.back()
}

function triggerAvatarUpload() {
  avatarInput.value?.click()
}

function handleAvatarSelect(event) {
  const file = event.target.files?.[0]
  if (!file) return

  if (file.size > 10 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过10MB', 'error')
    return
  }

  cropperImageSrc.value = URL.createObjectURL(file)
  showCropper.value = true
  event.target.value = ''
}

async function handleCropConfirm(croppedFile) {
  try {
    const avatarUrl = await uploadImage(croppedFile, 'avatar')
    form.avatar = avatarUrl
    appStore.showToast('头像上传成功', 'success')
  } catch (err) {
    appStore.showToast(err.message || '头像上传失败', 'error')
  } finally {
    if (cropperImageSrc.value) {
      URL.revokeObjectURL(cropperImageSrc.value)
      cropperImageSrc.value = ''
    }
  }
}

function handleCropCancel() {
  if (cropperImageSrc.value) {
    URL.revokeObjectURL(cropperImageSrc.value)
    cropperImageSrc.value = ''
  }
}

async function handleSave() {
  if (!form.nickname?.trim()) {
    appStore.showToast('请输入昵称', 'error')
    return
  }

  saving.value = true
  try {
    await userStore.updateProfile({
      avatar: form.avatar,
      nickname: form.nickname.trim(),
      gender: form.gender,
      city: form.city,
      contact_phone: form.contact_phone,
      contact_wechat: form.contact_wechat,
      emergency_contact_name: form.emergency_contact_name,
      emergency_contact_phone: form.emergency_contact_phone,
      emergency_contact_relation: form.emergency_contact_relation,
      car_number: form.car_number?.toUpperCase(),
      car_brand: form.car_brand,
      car_model: form.car_model,
      car_color: form.car_color
    })
    appStore.showToast('保存成功', 'success')
    router.back()
  } catch (err) {
    appStore.showToast(err.message || '保存失败', 'error')
  } finally {
    saving.value = false
  }
}
</script>
