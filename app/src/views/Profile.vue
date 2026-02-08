<template>
  <div class="profile-page min-h-screen" :class="appStore.theme === 'dark' ? 'bg-gray-900' : 'bg-gray-100'">
    <div v-if="!userStore.isLoggedIn" class="text-center py-12">
      <div class="w-20 h-20 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
        <svg class="w-10 h-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>
      </div>
      <p class="text-gray-500 mb-4">登录后查看个人信息</p>
      <router-link to="/login" class="btn btn-primary">
        立即登录
      </router-link>
    </div>

    <div v-else>
      <!-- 用户信息卡片 - 仿微信风格 -->
      <div class="bg-white px-4 py-5" :class="appStore.theme === 'dark' ? 'bg-gray-800' : ''">
        <div class="flex items-center gap-4">
          <!-- Avatar -->
          <router-link to="/profile/edit" class="relative flex-shrink-0">
            <div
              class="w-16 h-16 bg-gray-200 rounded-xl flex items-center justify-center overflow-hidden"
            >
              <img v-if="userStore.user?.avatar" :src="userStore.user.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-2xl font-bold text-gray-400">{{ userStore.user?.nickname?.charAt(0) || '?' }}</span>
            </div>
          </router-link>
          
          <div class="flex-1 min-w-0">
            <div class="text-lg font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ userStore.user?.nickname }}</div>
            <div class="text-sm text-gray-500 mt-0.5">手机号：{{ maskPhone(userStore.user?.phone) }}</div>
            <!-- 当前身份标签 -->
            <div class="mt-1.5 flex items-center gap-2">
              <span 
                class="px-2 py-0.5 text-xs rounded-full"
                :class="userStore.identity === 1 ? 'bg-blue-100 text-blue-600' : 'bg-green-100 text-green-600'"
              >
                {{ userStore.identity === 1 ? '🚗 车主' : '🙋 乘客' }}
              </span>
            </div>
          </div>
          
          <!-- 右侧箭头和二维码图标 -->
          <div class="flex items-center gap-3">
          <button @click="showQRCode = true" class="p-1">
            <svg class="w-5 h-5" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
              </svg>
            </button>
            <router-link to="/profile/edit">
              <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </router-link>
          </div>
        </div>
      </div>

      <!-- 功能列表 - 仿微信风格，统一蓝灰色调 -->
      <div class="mt-2">
        <!-- 第一组 - 主要功能（主题色） -->
        <div class="bg-white" :class="appStore.theme === 'dark' ? 'bg-gray-800' : ''">
          <router-link 
            to="/my-trips" 
            class="list-item-wechat"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: var(--theme-primary);">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">我的行程</span>
            </div>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
          
          <router-link 
            to="/publish" 
            class="list-item-wechat"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: var(--theme-primary);">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">发布行程</span>
            </div>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>

          <router-link 
            to="/notifications" 
            class="list-item-wechat"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" style="background: var(--theme-primary);">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">系统通知</span>
            </div>
            <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>

        <!-- 第二组 - 服务与帮助（灰色系） -->
        <div class="mt-2 bg-white" :class="appStore.theme === 'dark' ? 'bg-gray-800' : ''">
          <div 
            class="list-item-wechat cursor-pointer"
            @click="showHelp('safety')"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-400'">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">安全须知</span>
            </div>
            <svg class="w-5 h-5" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          
          <div 
            class="list-item-wechat cursor-pointer"
            @click="showHelp('faq')"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-400'">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">常见问题</span>
            </div>
            <svg class="w-5 h-5" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          
          <div 
            class="list-item-wechat cursor-pointer"
            @click="showHelp('feedback')"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-400'">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">意见反馈</span>
            </div>
            <svg class="w-5 h-5" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>

          <router-link 
            to="/settings" 
            class="list-item-wechat"
          >
            <div class="flex items-center gap-3">
              <div class="w-8 h-8 rounded-lg flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-600' : 'bg-gray-400'">
                <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <span class="text-base" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">设置</span>
            </div>
            <svg class="w-5 h-5" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </div>
      </div>
    </div>

    <!-- 帮助弹窗 -->
    <div v-if="helpModal" class="action-sheet" @click.self="helpModal = null">
      <div class="action-sheet-overlay" @click="helpModal = null"></div>
      <div class="action-sheet-content max-h-[80vh]">
        <div class="px-4 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="text-lg font-bold" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ helpModalTitle }}</h3>
          <button @click="helpModal = null" class="w-8 h-8 rounded-full flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-100'">
            <svg class="w-4 h-4" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-4 overflow-y-auto max-h-[60vh]">
          <!-- 安全须知 -->
          <template v-if="helpModal === 'safety'">
            <div class="space-y-4 text-sm">
              <div class="bg-red-50 p-4 rounded-xl" :class="appStore.theme === 'dark' ? 'bg-red-900/20' : ''">
                <h4 class="font-semibold text-red-600 mb-2">⚠️ 重要提醒</h4>
                <p class="text-red-600/80">本平台仅提供信息发布服务，不对行程安全负责。请务必注意以下事项：</p>
              </div>
              
              <div>
                <h4 class="font-semibold mb-2">出行前</h4>
                <ul class="list-disc list-inside space-y-1 text-gray-600" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">
                  <li>核实对方真实身份和联系方式</li>
                  <li>选择公共场所见面交接</li>
                  <li>告知家人或朋友行程信息</li>
                  <li>查看对方历史评价</li>
                </ul>
              </div>
              
              <div>
                <h4 class="font-semibold mb-2">出行中</h4>
                <ul class="list-disc list-inside space-y-1 text-gray-600" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">
                  <li>使用行程分享功能</li>
                  <li>保持手机电量充足</li>
                  <li>注意保管贵重物品</li>
                  <li>遇到异常及时报警</li>
                </ul>
              </div>
            </div>
          </template>
          
          <!-- 常见问题 -->
          <template v-else-if="helpModal === 'faq'">
            <div class="space-y-4">
              <div v-for="(faq, idx) in faqs" :key="idx" class="border-b border-gray-100 pb-4" :class="appStore.theme === 'dark' ? 'border-gray-700' : ''">
                <h4 class="font-semibold text-sm mb-2">Q: {{ faq.q }}</h4>
                <p class="text-sm text-gray-600" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">A: {{ faq.a }}</p>
              </div>
            </div>
          </template>
          
          <!-- 意见反馈 -->
          <template v-else-if="helpModal === 'feedback'">
            <div class="space-y-4">
              <p class="text-sm text-gray-500">您的反馈对我们非常重要，请描述您遇到的问题或建议：</p>
              <textarea
                v-model="feedbackContent"
                rows="5"
                class="input resize-none"
                placeholder="请输入您的反馈内容..."
              ></textarea>
              <button 
                @click="submitFeedback" 
                class="btn btn-primary w-full"
                :disabled="!feedbackContent.trim()"
              >
                提交反馈
              </button>
            </div>
          </template>
        </div>
      </div>
    </div>

    <!-- 二维码弹窗 -->
    <div v-if="showQRCode" class="action-sheet" @click.self="showQRCode = false">
      <div class="action-sheet-overlay" @click="showQRCode = false"></div>
      <div class="action-sheet-content">
        <div class="px-4 py-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="text-lg font-bold" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">我的二维码</h3>
          <button @click="showQRCode = false" class="w-8 h-8 rounded-full flex items-center justify-center" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-100'">
            <svg class="w-4 h-4" :class="appStore.theme === 'dark' ? 'text-gray-400' : 'text-gray-500'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-6 flex flex-col items-center">
          <!-- 用户信息 -->
          <div class="flex items-center gap-3 mb-6">
            <div class="w-12 h-12 rounded-xl flex items-center justify-center overflow-hidden" :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'">
              <img v-if="userStore.user?.avatar" :src="userStore.user.avatar" class="w-full h-full object-cover" />
              <span v-else class="text-lg font-bold" :class="appStore.theme === 'dark' ? 'text-gray-500' : 'text-gray-400'">{{ userStore.user?.nickname?.charAt(0) || '?' }}</span>
            </div>
            <div>
              <div class="font-semibold" :class="appStore.theme === 'dark' ? 'text-white' : 'text-gray-900'">{{ userStore.user?.nickname }}</div>
              <div class="text-xs text-gray-500">{{ userStore.identity === 1 ? '车主' : '乘客' }}</div>
            </div>
          </div>
          
          <!-- 二维码区域 -->
          <div class="qrcode-container p-4 rounded-2xl border shadow-sm" :class="appStore.theme === 'dark' ? 'bg-gray-700 border-gray-600' : 'bg-white border-gray-200'">
            <canvas ref="qrcodeCanvas" class="w-48 h-48"></canvas>
          </div>
          
          <p class="text-sm text-gray-500 mt-4">扫一扫上面的二维码，加我为好友</p>
          <p class="text-xs text-gray-400 mt-1">（好友功能即将上线）</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useUserStore } from '@/stores/user'
import { useAppStore } from '@/stores/app'

const userStore = useUserStore()
const appStore = useAppStore()

const helpModal = ref(null)
const feedbackContent = ref('')
const showQRCode = ref(false)
const qrcodeCanvas = ref(null)

// 生成二维码
function generateQRCode() {
  if (!qrcodeCanvas.value || !userStore.user) return
  
  const canvas = qrcodeCanvas.value
  const ctx = canvas.getContext('2d')
  const size = 192
  canvas.width = size
  canvas.height = size
  
  // 简单生成一个二维码样式（实际项目可使用qrcode库）
  // 这里先用占位图案
  ctx.fillStyle = '#fff'
  ctx.fillRect(0, 0, size, size)
  
  // 绘制二维码模拟图案
  const moduleSize = 6
  const modules = Math.floor(size / moduleSize)
  const userId = userStore.user?.id || 0
  
  ctx.fillStyle = '#000'
  
  // 位置探测图形（三个角）
  drawFinderPattern(ctx, 0, 0, moduleSize)
  drawFinderPattern(ctx, (modules - 7) * moduleSize, 0, moduleSize)
  drawFinderPattern(ctx, 0, (modules - 7) * moduleSize, moduleSize)
  
  // 基于用户ID生成伪随机数据模块
  const seed = userId * 12345
  for (let i = 8; i < modules - 8; i++) {
    for (let j = 8; j < modules - 8; j++) {
      const hash = ((i * 31 + j * 17 + seed) % 100)
      if (hash < 45) {
        ctx.fillRect(i * moduleSize, j * moduleSize, moduleSize - 1, moduleSize - 1)
      }
    }
  }
}

function drawFinderPattern(ctx, x, y, moduleSize) {
  // 外框
  ctx.fillRect(x, y, 7 * moduleSize, moduleSize)
  ctx.fillRect(x, y + 6 * moduleSize, 7 * moduleSize, moduleSize)
  ctx.fillRect(x, y, moduleSize, 7 * moduleSize)
  ctx.fillRect(x + 6 * moduleSize, y, moduleSize, 7 * moduleSize)
  // 内框
  ctx.fillRect(x + 2 * moduleSize, y + 2 * moduleSize, 3 * moduleSize, 3 * moduleSize)
}

// 监听弹窗显示，生成二维码
watch(showQRCode, (val) => {
  if (val) {
    nextTick(() => {
      generateQRCode()
    })
  }
})

const helpModalTitle = computed(() => {
  switch (helpModal.value) {
    case 'safety': return '🛡️ 安全须知'
    case 'faq': return '❓ 常见问题'
    case 'feedback': return '💬 意见反馈'
    default: return ''
  }
})

const faqs = [
  { q: '如何发布行程？', a: '点击"发布行程"按钮，填写出发地、目的地、出发时间等信息即可发布。' },
  { q: '如何联系对方？', a: '在行程详情页点击"联系TA"按钮，即可通过平台私信与对方沟通。' },
  { q: '如何取消行程？', a: '在"我的行程"中找到对应行程，点击"取消行程"即可。' },
  { q: '费用如何支付？', a: '平台不参与费用收取，请与同行者自行协商支付方式。' },
  { q: '遇到问题怎么办？', a: '如遇紧急情况请立即报警；非紧急问题可通过"意见反馈"联系我们。' }
]

function maskPhone(phone) {
  if (!phone || phone.length !== 11) return phone
  return phone.slice(0, 3) + '****' + phone.slice(7)
}

function showHelp(type) {
  helpModal.value = type
}

function submitFeedback() {
  if (!feedbackContent.value.trim()) return
  appStore.showToast('感谢您的反馈！', 'success')
  feedbackContent.value = ''
  helpModal.value = null
}
</script>

<style scoped>
.list-item-wechat {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  border-bottom: 1px solid #f3f4f6;
}

.list-item-wechat:last-child {
  border-bottom: none;
}

.list-item-wechat:active {
  background-color: #f3f4f6;
}

[data-theme="dark"] .list-item-wechat {
  border-bottom-color: #374151;
}

[data-theme="dark"] .list-item-wechat:last-child {
  border-bottom-color: transparent;
}

[data-theme="dark"] .list-item-wechat:active {
  background-color: #374151;
}
</style>
