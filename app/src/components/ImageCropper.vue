<template>
  <teleport to="body">
    <div v-if="visible" class="fixed inset-0 z-50 flex items-center justify-center bg-black/80">
      <div class="w-full max-w-md mx-4 bg-white rounded-2xl overflow-hidden">
        <!-- Header -->
        <div class="flex items-center justify-between px-4 py-3 border-b border-gray-100">
          <button @click="handleCancel" class="text-gray-500 text-sm">取消</button>
          <span class="text-gray-800 font-medium">裁剪头像</span>
          <button @click="handleConfirm" class="text-primary-500 font-medium text-sm">确定</button>
        </div>

        <!-- Cropper Area -->
        <div class="relative w-full aspect-square bg-gray-900 overflow-hidden">
          <div ref="containerRef" class="w-full h-full flex items-center justify-center">
            <canvas ref="canvasRef" class="touch-none" />
          </div>
          <!-- Circular mask overlay -->
          <div class="absolute inset-0 pointer-events-none">
            <svg class="w-full h-full">
              <defs>
                <mask id="circleMask">
                  <rect width="100%" height="100%" fill="white" />
                  <circle cx="50%" cy="50%" r="40%" fill="black" />
                </mask>
              </defs>
              <rect width="100%" height="100%" fill="rgba(0,0,0,0.6)" mask="url(#circleMask)" />
              <circle cx="50%" cy="50%" r="40%" fill="none" stroke="white" stroke-width="2" />
            </svg>
          </div>
        </div>

        <!-- Zoom Controls -->
        <div class="flex items-center justify-center gap-6 py-4">
          <button
            @click="zoom(-0.1)"
            class="w-10 h-10 flex items-center justify-center rounded-full bg-gray-100 text-gray-600 active:bg-gray-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
            </svg>
          </button>
          <button
            @click="rotate(-90)"
            class="w-10 h-10 flex items-center justify-center rounded-full bg-gray-100 text-gray-600 active:bg-gray-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
          </button>
          <button
            @click="zoom(0.1)"
            class="w-10 h-10 flex items-center justify-center rounded-full bg-gray-100 text-gray-600 active:bg-gray-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </teleport>
</template>

<script setup>
import { ref, watch, onUnmounted, nextTick } from 'vue'

const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  imageSrc: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:visible', 'confirm', 'cancel'])

const containerRef = ref(null)
const canvasRef = ref(null)

// cropper state
const state = ref({
  image: null,
  scale: 1,
  rotation: 0,
  offsetX: 0,
  offsetY: 0,
  isDragging: false,
  lastX: 0,
  lastY: 0
})

function loadImage() {
  if (!props.imageSrc) return

  const img = new Image()
  img.crossOrigin = 'anonymous'
  img.onload = () => {
    state.value.image = img
    state.value.scale = 1
    state.value.rotation = 0
    state.value.offsetX = 0
    state.value.offsetY = 0
    nextTick(() => {
      fitImageToCanvas()
      drawImage()
      setupEventListeners()
    })
  }
  img.src = props.imageSrc
}

function fitImageToCanvas() {
  const canvas = canvasRef.value
  const container = containerRef.value
  if (!canvas || !container || !state.value.image) return

  const containerSize = Math.min(container.clientWidth, container.clientHeight)
  canvas.width = containerSize
  canvas.height = containerSize

  const img = state.value.image
  const isRotated = Math.abs(state.value.rotation % 180) === 90
  const imgW = isRotated ? img.height : img.width
  const imgH = isRotated ? img.width : img.height

  // scale image to fit (cover the circular crop area - 80% of container)
  const cropDiameter = containerSize * .8
  const scale = Math.max(cropDiameter / imgW, cropDiameter / imgH)
  state.value.scale = scale
}

function drawImage() {
  const canvas = canvasRef.value
  const ctx = canvas?.getContext('2d')
  if (!ctx || !state.value.image) return

  const { image, scale, rotation, offsetX, offsetY } = state.value

  ctx.clearRect(0, 0, canvas.width, canvas.height)
  ctx.save()

  // move to center
  ctx.translate(canvas.width / 2 + offsetX, canvas.height / 2 + offsetY)
  ctx.rotate((rotation * Math.PI) / 180)
  ctx.scale(scale, scale)

  // draw image centered
  ctx.drawImage(image, -image.width / 2, -image.height / 2)

  ctx.restore()
}

function setupEventListeners() {
  const canvas = canvasRef.value
  if (!canvas) return

  // mouse events
  canvas.addEventListener('mousedown', handleDragStart)
  canvas.addEventListener('mousemove', handleDragMove)
  canvas.addEventListener('mouseup', handleDragEnd)
  canvas.addEventListener('mouseleave', handleDragEnd)

  // touch events
  canvas.addEventListener('touchstart', handleTouchStart, { passive: false })
  canvas.addEventListener('touchmove', handleTouchMove, { passive: false })
  canvas.addEventListener('touchend', handleTouchEnd)

  // wheel zoom
  canvas.addEventListener('wheel', handleWheel, { passive: false })
}

function removeEventListeners() {
  const canvas = canvasRef.value
  if (!canvas) return

  canvas.removeEventListener('mousedown', handleDragStart)
  canvas.removeEventListener('mousemove', handleDragMove)
  canvas.removeEventListener('mouseup', handleDragEnd)
  canvas.removeEventListener('mouseleave', handleDragEnd)
  canvas.removeEventListener('touchstart', handleTouchStart)
  canvas.removeEventListener('touchmove', handleTouchMove)
  canvas.removeEventListener('touchend', handleTouchEnd)
  canvas.removeEventListener('wheel', handleWheel)
}

function handleDragStart(e) {
  state.value.isDragging = true
  state.value.lastX = e.clientX
  state.value.lastY = e.clientY
}

function handleDragMove(e) {
  if (!state.value.isDragging) return
  const dx = e.clientX - state.value.lastX
  const dy = e.clientY - state.value.lastY
  state.value.offsetX += dx
  state.value.offsetY += dy
  state.value.lastX = e.clientX
  state.value.lastY = e.clientY
  drawImage()
}

function handleDragEnd() {
  state.value.isDragging = false
}

function handleTouchStart(e) {
  if (e.touches.length === 1) {
    e.preventDefault()
    state.value.isDragging = true
    state.value.lastX = e.touches[0].clientX
    state.value.lastY = e.touches[0].clientY
  }
}

function handleTouchMove(e) {
  if (!state.value.isDragging || e.touches.length !== 1) return
  e.preventDefault()
  const touch = e.touches[0]
  const dx = touch.clientX - state.value.lastX
  const dy = touch.clientY - state.value.lastY
  state.value.offsetX += dx
  state.value.offsetY += dy
  state.value.lastX = touch.clientX
  state.value.lastY = touch.clientY
  drawImage()
}

function handleTouchEnd() {
  state.value.isDragging = false
}

function handleWheel(e) {
  e.preventDefault()
  const delta = e.deltaY > 0 ? -0.05 : 0.05
  zoom(delta)
}

function zoom(delta) {
  state.value.scale = Math.max(.1, Math.min(5, state.value.scale + delta))
  drawImage()
}

function rotate(degree) {
  state.value.rotation = (state.value.rotation + degree) % 360
  fitImageToCanvas()
  drawImage()
}

function handleCancel() {
  removeEventListeners()
  emit('update:visible', false)
  emit('cancel')
}

function handleConfirm() {
  const canvas = canvasRef.value
  if (!canvas || !state.value.image) return

  // create output canvas for cropped result
  const outputSize = 256
  const outputCanvas = document.createElement('canvas')
  outputCanvas.width = outputSize
  outputCanvas.height = outputSize
  const outputCtx = outputCanvas.getContext('2d')

  // calculate crop area (circular area in the center, 80% of canvas)
  const cropRadius = canvas.width * .4
  const cropCenterX = canvas.width / 2
  const cropCenterY = canvas.height / 2

  // draw circular clip
  outputCtx.beginPath()
  outputCtx.arc(outputSize / 2, outputSize / 2, outputSize / 2, 0, Math.PI * 2)
  outputCtx.closePath()
  outputCtx.clip()

  // calculate source region from main canvas
  const sourceX = cropCenterX - cropRadius
  const sourceY = cropCenterY - cropRadius
  const sourceSize = cropRadius * 2

  outputCtx.drawImage(
    canvas,
    sourceX, sourceY, sourceSize, sourceSize,
    0, 0, outputSize, outputSize
  )

  outputCanvas.toBlob(
    (blob) => {
      if (blob) {
        const file = new File([blob], 'avatar.jpg', { type: 'image/jpeg' })
        emit('confirm', file)
        removeEventListeners()
        emit('update:visible', false)
      }
    },
    'image/jpeg',
    .9
  )
}

watch(
  () => props.visible,
  (val) => {
    if (val) {
      nextTick(() => {
        loadImage()
      })
    } else {
      removeEventListeners()
      state.value.image = null
    }
  }
)

watch(
  () => props.imageSrc,
  () => {
    if (props.visible) {
      loadImage()
    }
  }
)

onUnmounted(() => {
  removeEventListeners()
})
</script>
