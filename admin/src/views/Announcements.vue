<template>
  <div class="announcements-page">
    <!-- 顶部操作 -->
    <div class="mb-4">
      <button @click="openCreateModal" class="w-full btn btn-primary">
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        发布公告
      </button>
    </div>

    <!-- 公告列表 -->
    <div class="card">
      <div v-if="loading" class="flex justify-center py-12">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="announcements.length === 0" class="empty-state">
        <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5.882V19.24a1.76 1.76 0 01-3.417.592l-2.147-6.15M18 13a3 3 0 100-6M5.436 13.683A4.001 4.001 0 017 6h1.832c4.1 0 7.625-1.234 9.168-3v14c-1.543-1.766-5.067-3-9.168-3H7a3.988 3.988 0 01-1.564-.317z" />
        </svg>
        <p>暂无公告</p>
      </div>

      <div v-else>
        <div
          v-for="item in announcements"
          :key="item.id"
          class="list-item"
          @click="openEditModal(item)"
        >
          <div class="flex-1 min-w-0">
            <div class="flex items-center gap-2 mb-1">
              <span :class="getTypeBadgeClass(item.type)" class="badge text-[11px]">
                {{ getTypeText(item.type) }}
              </span>
              <span :class="item.is_active ? 'badge-success' : 'badge-gray'" class="badge text-[11px]">
                {{ item.is_active ? '启用' : '禁用' }}
              </span>
            </div>
            <div class="font-medium text-gray-800 text-sm truncate mb-1">{{ item.title }}</div>
            <div class="text-xs text-gray-400">
              {{ formatTime(item.created_at) }}
              <span v-if="item.end_time"> · 有效至 {{ formatTime(item.end_time) }}</span>
            </div>
          </div>
          <svg class="w-5 h-5 text-gray-300 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
          </svg>
        </div>
      </div>
    </div>

    <!-- 创建/编辑弹窗 -->
    <div v-if="showModal" class="modal-overlay" @click.self="closeModal">
      <div class="modal-content">
        <div class="p-4 border-b border-gray-100 flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-800">
            {{ isEdit ? '编辑公告' : '发布公告' }}
          </h3>
          <button @click="closeModal" class="p-2 -mr-2 text-gray-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-4 space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">标题 *</label>
            <input v-model="form.title" type="text" class="input" placeholder="请输入公告标题" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">内容 *</label>
            <textarea
              v-model="form.content"
              rows="4"
              class="input"
              placeholder="请输入公告内容"
            ></textarea>
          </div>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">类型</label>
              <select v-model="form.type" class="input">
                <option value="1">普通</option>
                <option value="2">紧急</option>
                <option value="3">活动</option>
              </select>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">排序</label>
              <input v-model.number="form.sort_order" type="number" class="input" placeholder="越大越靠前" />
            </div>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">开始时间</label>
            <input v-model="form.start_time" type="datetime-local" class="input" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">结束时间</label>
            <input v-model="form.end_time" type="datetime-local" class="input" />
          </div>
          <div class="flex items-center gap-3">
            <input v-model="form.is_active" type="checkbox" id="is_active" class="w-5 h-5 rounded" />
            <label for="is_active" class="text-sm text-gray-700">立即启用</label>
          </div>
        </div>
        <div class="p-4 border-t border-gray-100 space-y-3">
          <button @click="handleSubmit" :disabled="submitting" class="w-full btn btn-primary">
            {{ submitting ? '保存中...' : '保存' }}
          </button>
          <button v-if="isEdit" @click="handleDelete" class="w-full btn btn-danger">
            删除公告
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import api from '@/utils/api'

const announcements = ref([])
const loading = ref(true)
const showModal = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const currentId = ref(null)

const defaultForm = {
  title: '',
  content: '',
  type: 1,
  is_active: true,
  sort_order: 0,
  start_time: '',
  end_time: ''
}

const form = ref({ ...defaultForm })

onMounted(() => {
  fetchAnnouncements()
})

async function fetchAnnouncements() {
  loading.value = true
  try {
    const data = await api.get('/admin/announcements')
    announcements.value = data.list || []
  } catch (e) {
    console.error('Failed to fetch announcements:', e)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEdit.value = false
  currentId.value = null
  form.value = { ...defaultForm }
  showModal.value = true
}

function openEditModal(item) {
  isEdit.value = true
  currentId.value = item.id
  form.value = {
    title: item.title,
    content: item.content,
    type: item.type,
    is_active: item.is_active === 1,
    sort_order: item.sort_order,
    start_time: item.start_time ? formatDateTimeLocal(item.start_time) : '',
    end_time: item.end_time ? formatDateTimeLocal(item.end_time) : ''
  }
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

async function handleSubmit() {
  if (!form.value.title || !form.value.content) {
    alert('请填写标题和内容')
    return
  }

  submitting.value = true
  try {
    const payload = {
      title: form.value.title,
      content: form.value.content,
      type: parseInt(form.value.type),
      is_active: form.value.is_active ? 1 : 0,
      sort_order: form.value.sort_order,
      start_time: form.value.start_time ? new Date(form.value.start_time).toISOString() : null,
      end_time: form.value.end_time ? new Date(form.value.end_time).toISOString() : null
    }

    if (isEdit.value) {
      await api.put(`/admin/announcements/${currentId.value}`, payload)
    } else {
      await api.post('/admin/announcements', payload)
    }

    closeModal()
    fetchAnnouncements()
    alert(isEdit.value ? '更新成功' : '发布成功')
  } catch (e) {
    alert('操作失败: ' + (e.response?.data?.message || e.message))
  } finally {
    submitting.value = false
  }
}

async function handleDelete() {
  if (!confirm('确定要删除这个公告吗？')) return

  try {
    await api.delete(`/admin/announcements/${currentId.value}`)
    closeModal()
    fetchAnnouncements()
    alert('删除成功')
  } catch (e) {
    alert('删除失败: ' + (e.response?.data?.message || e.message))
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

function formatDateTimeLocal(time) {
  const date = new Date(time)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day}T${hour}:${minute}`
}

function getTypeBadgeClass(type) {
  switch (type) {
    case 1: return 'badge-info'
    case 2: return 'badge-danger'
    case 3: return 'badge-success'
    default: return 'badge-gray'
  }
}

function getTypeText(type) {
  switch (type) {
    case 1: return '普通'
    case 2: return '紧急'
    case 3: return '活动'
    default: return '未知'
  }
}
</script>
