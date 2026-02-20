<template>
  <div class="chat-page" :class="appStore.theme === 'dark' ? 'chat-page--dark' : ''">
    <!-- 顶部导航 - 固定 -->
    <div class="chat-header border-b safe-area-top" :class="appStore.theme === 'dark' ? 'bg-gray-800 border-gray-700' : 'bg-white border-gray-100'">
      <div class="flex items-center h-12 px-4">
        <button @click="goBack" class="p-2 -ml-2">
          <svg class="w-6 h-6" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <h1 class="flex-1 text-center text-lg font-semibold truncate cursor-pointer" :class="appStore.theme === 'dark' ? 'text-white' : ''" @click="goToUserProfile">{{ peerNickname }}</h1>
        <div class="w-10"></div>
      </div>
    </div>

    <!-- 消息列表 - 中间可滚动区域 -->
    <div ref="messageListRef" class="chat-messages" @scroll="handleScroll">
      <!-- 加载更多 -->
      <div v-if="hasMore" class="flex justify-center mb-4">
        <button
          @click="loadMoreMessages"
          :disabled="loadingMore"
          class="px-4 py-1 text-xs rounded-full"
          :class="appStore.theme === 'dark' ? 'text-gray-400 bg-gray-700' : 'text-gray-500 bg-gray-100'"
        >
          {{ loadingMore ? '加载中...' : '加载更多' }}
        </button>
      </div>

      <div v-if="loading" class="flex justify-center py-8">
        <div class="loading-spinner"></div>
      </div>

      <div v-else-if="messages.length === 0" class="text-center py-8 text-gray-400 text-sm">
        暂无消息，开始聊天吧
      </div>

      <!-- 消息列表（倒序显示，最新在底部） -->
      <div v-else class="space-y-4">
        <div
          v-for="(msg, index) in reversedMessages"
          :key="msg.id"
        >
          <!-- 时间分隔 -->
          <div
            v-if="shouldShowTime(msg, index)"
            class="flex justify-center mb-3"
          >
            <span class="px-2 py-1 text-xs rounded" :class="appStore.theme === 'dark' ? 'text-gray-500 bg-gray-700' : 'text-gray-400 bg-gray-100'">
              {{ formatMessageTime(msg.created_at) }}
            </span>
          </div>

          <!-- 消息气泡 -->
          <div
            class="flex items-end gap-2"
            :class="msg.sender_id === currentUserOpenId ? 'flex-row-reverse' : ''"
          >
            <!-- 头像 -->
            <div 
              class="w-8 h-8 rounded-full flex-shrink-0 flex items-center justify-center overflow-hidden cursor-pointer" 
              :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"
              @click="msg.sender_id !== currentUserOpenId && goToUserProfile()"
            >
              <img 
                v-if="msg.sender_id === currentUserOpenId && userStore.user?.avatar" 
                :src="userStore.user.avatar" 
                class="w-full h-full object-cover" 
              />
              <img 
                v-else-if="msg.sender_id !== currentUserOpenId && peerAvatar" 
                :src="peerAvatar" 
                class="w-full h-full object-cover" 
              />
              <span v-else class="text-xs font-medium text-gray-500">
                {{ msg.sender_id === currentUserOpenId ? myInitial : peerInitial }}
              </span>
            </div>

            <!-- 消息内容 -->
            <div
              class="max-w-[70%]"
              :class="(msg.msg_type === 1 || msg.msg_type === 3)
                ? (msg.sender_id === currentUserOpenId
                  ? 'px-3 py-2 rounded-2xl bg-primary-500 text-white rounded-br-sm'
                  : appStore.theme === 'dark' ? 'px-3 py-2 rounded-2xl bg-gray-700 text-gray-200 rounded-bl-sm' : 'px-3 py-2 rounded-2xl bg-white text-gray-800 rounded-bl-sm shadow-sm')
                : ''"
            >
              <!-- 文字消息 -->
              <p v-if="msg.msg_type === 1" class="text-sm whitespace-pre-wrap break-words">
                {{ msg.content }}
              </p>
              <!-- 图片消息 - 显示缩略图，点击加载原图 -->
              <div
                v-else-if="msg.msg_type === 2"
                class="image-message relative cursor-pointer"
                @click="previewImage(msg.content)"
              >
                <img
                  :src="thumbnailUrlCache.get(msg.content) || imageUrlCache.get(msg.content) || ''"
                  @load="onImageLoad"
                  class="max-w-full rounded-xl"
                  style="max-height: 200px; min-width: 80px; min-height: 80px"
                  :data-key="msg.content"
                />
                <!-- 加载中状态 -->
                <div
                  v-if="!thumbnailUrlCache.get(msg.content) && !imageUrlCache.get(msg.content)"
                  class="w-32 h-32 rounded-xl flex items-center justify-center absolute inset-0"
                  :class="appStore.theme === 'dark' ? 'bg-gray-700' : 'bg-gray-200'"
                >
                  <span class="text-xs text-gray-400">加载中...</span>
                </div>
              </div>
              <!-- 语音消息 -->
              <div
                v-else-if="msg.msg_type === 3"
                @click="playVoice(msg)"
                class="flex items-center gap-2 cursor-pointer min-w-20"
                :class="msg.sender_id === currentUserOpenId ? 'flex-row-reverse' : ''"
              >
                <!-- 语音播放动画图标 -->
                <div class="voice-icon" :class="{ 'voice-icon--playing': playingVoiceId === msg.id, 'voice-icon--self': msg.sender_id === currentUserOpenId }">
                  <span></span>
                  <span></span>
                  <span></span>
                </div>
                <span class="text-sm">{{ msg.duration || 0 }}"</span>
              </div>
              <!-- 表情消息 -->
              <div v-else-if="msg.msg_type === 4" class="text-4xl">
                {{ msg.content }}
              </div>
              <!-- 通话记录消息 -->
              <div
                v-else-if="msg.msg_type === 5"
                class="call-record px-3 py-2 rounded-2xl flex items-center gap-2"
                :class="msg.sender_id === currentUserOpenId
                  ? 'bg-primary-500 text-white rounded-br-sm'
                  : appStore.theme === 'dark' ? 'bg-gray-700 text-gray-200 rounded-bl-sm' : 'bg-white text-gray-800 rounded-bl-sm shadow-sm'"
                @click="initiateCall(parseCallRecord(msg.content).call_type)"
              >
                <!-- 通话图标 -->
                <div class="call-record__icon">
                  <svg v-if="parseCallRecord(msg.content).call_type === 'video'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 10l4.553-2.276A1 1 0 0121 8.618v6.764a1 1 0 01-1.447.894L15 14M5 18h8a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                  <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                  </svg>
                </div>
                <!-- 通话信息 -->
                <div class="call-record__info text-sm">
                  <span>{{ formatCallRecord(msg) }}</span>
                </div>
              </div>
              <!-- 视频消息 - 显示缩略图，点击播放 -->
              <div
                v-else-if="msg.msg_type === 6"
                class="video-message relative cursor-pointer"
                @click="playVideo(msg)"
              >
                <img
                  :src="getVideoThumbnail(msg)"
                  class="max-w-full rounded-xl"
                  style="max-height: 200px; min-width: 120px; min-height: 80px; object-fit: cover"
                />
                <!-- 播放按钮 -->
                <div class="absolute inset-0 flex items-center justify-center">
                  <div class="w-12 h-12 rounded-full bg-black/50 flex items-center justify-center">
                    <svg class="w-6 h-6 text-white ml-1" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M8 5v14l11-7z" />
                    </svg>
                  </div>
                </div>
                <!-- 视频时长 -->
                <div class="absolute bottom-2 right-2 px-1.5 py-0.5 rounded bg-black/60 text-white text-xs">
                  {{ formatVideoDuration(parseVideoContent(msg.content).duration) }}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 底部输入区域 - 微信风格 -->
    <div class="chat-input safe-area-bottom" :class="appStore.theme === 'dark' ? 'chat-input--dark' : ''">
      <div class="chat-input__main">
        <!-- 语音/键盘切换按钮 -->
        <button
          @click="toggleVoiceMode"
          class="chat-input__icon-btn"
          :class="appStore.theme === 'dark' ? 'chat-input__icon-btn--dark' : ''"
        >
          <!-- 语音图标 -->
          <svg v-if="!voiceMode" class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 18.75a6 6 0 006-6v-1.5m-6 7.5a6 6 0 01-6-6v-1.5m6 7.5v3.75m-3.75 0h7.5M12 15.75a3 3 0 01-3-3V4.5a3 3 0 116 0v8.25a3 3 0 01-3 3z" />
          </svg>
          <!-- 键盘图标 -->
          <svg v-else class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.75 6.75h16.5M3.75 12h16.5m-16.5 5.25H12" />
          </svg>
        </button>

        <!-- 输入区域 -->
        <div class="chat-input__field-wrap">
          <!-- 语音录制按钮 -->
          <div
            v-if="voiceMode"
            @touchstart.prevent="startRecording"
            @touchend.prevent="stopRecording"
            @mousedown.prevent="startRecording"
            @mouseup.prevent="stopRecording"
            @mouseleave="cancelRecording"
            class="chat-input__voice-btn"
            :class="[
              recording ? 'chat-input__voice-btn--recording' : '',
              appStore.theme === 'dark' ? 'chat-input__voice-btn--dark' : ''
            ]"
          >
            {{ recording ? `松开 结束` : '按住 说话' }}
          </div>
          <!-- 文字输入框 -->
          <textarea
            v-else
            v-model="inputText"
            ref="textInputRef"
            rows="1"
            placeholder="输入消息"
            class="chat-input__textarea"
            :class="appStore.theme === 'dark' ? 'chat-input__textarea--dark' : ''"
            @keydown.enter.exact.prevent="sendMessage"
            @input="autoResize"
            @focus="handleInputFocus"
            @blur="handleInputBlur"
          ></textarea>
        </div>

        <!-- 表情按钮 -->
        <button
          @click="toggleEmojiPicker"
          class="chat-input__icon-btn"
          :class="[
            showEmojiPicker ? 'chat-input__icon-btn--active' : '',
            appStore.theme === 'dark' ? 'chat-input__icon-btn--dark' : ''
          ]"
        >
          <svg class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M15.182 15.182a4.5 4.5 0 01-6.364 0M21 12a9 9 0 11-18 0 9 9 0 0118 0zM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75zm-.375 0h.008v.015h-.008V9.75zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75zm-.375 0h.008v.015h-.008V9.75z" />
          </svg>
        </button>

        <!-- 更多/发送按钮 -->
        <button
          v-if="inputText.trim()"
          @click="sendMessage"
          :disabled="sending"
          class="chat-input__send-btn"
        >
          发送
        </button>
        <button
          v-else
          @click="toggleMorePanel"
          class="chat-input__icon-btn"
          :class="[
            showMorePanel ? 'chat-input__icon-btn--active' : '',
            appStore.theme === 'dark' ? 'chat-input__icon-btn--dark' : ''
          ]"
        >
          <svg class="w-7 h-7" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
          </svg>
        </button>
      </div>

      <!-- 表情选择器面板 -->
      <div v-if="showEmojiPicker" class="chat-input__panel" :class="appStore.theme === 'dark' ? 'chat-input__panel--dark' : ''">
        <div class="grid grid-cols-8 gap-1 p-3 max-h-52 overflow-y-auto">
          <button
            v-for="emoji in emojiList"
            :key="emoji"
            @click="sendEmoji(emoji)"
            class="chat-input__emoji-item"
            :class="appStore.theme === 'dark' ? 'hover:bg-gray-700 active:bg-gray-600' : ''"
          >
            {{ emoji }}
          </button>
        </div>
      </div>

      <!-- 更多功能面板 -->
      <div v-if="showMorePanel" class="chat-input__panel" :class="appStore.theme === 'dark' ? 'chat-input__panel--dark' : ''">
        <div class="grid grid-cols-4 gap-4 p-4">
          <!-- 图片 -->
          <div class="chat-input__more-item" @click="selectImage">
            <div class="chat-input__more-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
              <svg class="w-7 h-7" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
              </svg>
            </div>
            <span class="chat-input__more-text" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">图片</span>
          </div>
          <!-- 拍摄 -->
          <div class="chat-input__more-item" @click="openVideoRecorder">
            <div class="chat-input__more-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
              <svg class="w-7 h-7" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6.827 6.175A2.31 2.31 0 015.186 7.23c-.38.054-.757.112-1.134.175C2.999 7.58 2.25 8.507 2.25 9.574V18a2.25 2.25 0 002.25 2.25h15A2.25 2.25 0 0021.75 18V9.574c0-1.067-.75-1.994-1.802-2.169a47.865 47.865 0 00-1.134-.175 2.31 2.31 0 01-1.64-1.055l-.822-1.316a2.192 2.192 0 00-1.736-1.039 48.774 48.774 0 00-5.232 0 2.192 2.192 0 00-1.736 1.039l-.821 1.316z" />
                <path stroke-linecap="round" stroke-linejoin="round" d="M16.5 12.75a4.5 4.5 0 11-9 0 4.5 4.5 0 019 0zM18.75 10.5h.008v.008h-.008V10.5z" />
              </svg>
            </div>
            <span class="chat-input__more-text" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">拍摄</span>
          </div>
          <!-- 语音通话 -->
          <div class="chat-input__more-item" @click="startVoiceCall">
            <div class="chat-input__more-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
              <svg class="w-7 h-7" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 6.75c0 8.284 6.716 15 15 15h2.25a2.25 2.25 0 002.25-2.25v-1.372c0-.516-.351-.966-.852-1.091l-4.423-1.106c-.44-.11-.902.055-1.173.417l-.97 1.293c-.282.376-.769.542-1.21.38a12.035 12.035 0 01-7.143-7.143c-.162-.441.004-.928.38-1.21l1.293-.97c.363-.271.527-.734.417-1.173L6.963 3.102a1.125 1.125 0 00-1.091-.852H4.5A2.25 2.25 0 002.25 4.5v2.25z" />
              </svg>
            </div>
            <span class="chat-input__more-text" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">语音通话</span>
          </div>
          <!-- 视频通话 -->
          <div class="chat-input__more-item" @click="startVideoCall">
            <div class="chat-input__more-icon" :class="appStore.theme === 'dark' ? 'bg-gray-700' : ''">
              <svg class="w-7 h-7" :class="appStore.theme === 'dark' ? 'text-gray-300' : 'text-gray-600'" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M15.75 10.5l4.72-4.72a.75.75 0 011.28.53v11.38a.75.75 0 01-1.28.53l-4.72-4.72M4.5 18.75h9a2.25 2.25 0 002.25-2.25v-9a2.25 2.25 0 00-2.25-2.25h-9A2.25 2.25 0 002.25 7.5v9a2.25 2.25 0 002.25 2.25z" />
              </svg>
            </div>
            <span class="chat-input__more-text" :class="appStore.theme === 'dark' ? 'text-gray-400' : ''">视频通话</span>
          </div>
        </div>
      </div>

      <!-- 录音状态提示 -->
      <div v-if="recording" class="chat-input__recording-overlay">
        <div class="chat-input__recording-box">
          <div class="chat-input__recording-wave">
            <span></span><span></span><span></span><span></span><span></span>
          </div>
          <div class="chat-input__recording-time">{{ recordingDuration }}"</div>
          <div class="chat-input__recording-tip">松开发送，上滑取消</div>
        </div>
      </div>

      <input
        ref="imageInputRef"
        type="file"
        accept="image/*"
        class="hidden"
        @change="handleImageSelected"
      />
    </div>

    <!-- 图片预览遮罩 -->
    <div
      v-if="previewImageUrl"
      @click="previewImageUrl = ''"
      class="fixed inset-0 bg-black/90 z-50 flex items-center justify-center"
    >
      <img :src="previewImageUrl" class="max-w-full max-h-full" />
    </div>

    <!-- 视频录制器 -->
    <div v-if="showVideoRecorder" class="fixed inset-0 bg-black z-50 flex flex-col">
      <!-- 预览视频流 -->
      <div class="flex-1 relative">
        <video
          ref="videoPreviewRef"
          class="w-full h-full object-cover"
          autoplay
          playsinline
          muted
        ></video>
        <!-- 录制状态 -->
        <div v-if="videoRecording" class="absolute top-4 left-1/2 -translate-x-1/2 px-3 py-1 rounded-full bg-red-500 text-white text-sm flex items-center gap-2">
          <span class="w-2 h-2 rounded-full bg-white animate-pulse"></span>
          {{ videoRecordingDuration }}s
        </div>
        <!-- 关闭按钮 -->
        <button
          @click="closeVideoRecorder"
          class="absolute top-4 left-4 w-10 h-10 rounded-full bg-black/50 flex items-center justify-center"
        >
          <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      <!-- 控制栏 -->
      <div class="p-6 flex justify-center items-center gap-8 bg-black">
        <!-- 录制按钮 -->
        <button
          @touchstart.prevent="startVideoRecording"
          @touchend.prevent="stopVideoRecording"
          @mousedown.prevent="startVideoRecording"
          @mouseup.prevent="stopVideoRecording"
          class="w-20 h-20 rounded-full border-4 border-white flex items-center justify-center"
          :class="videoRecording ? 'bg-red-500' : 'bg-transparent'"
        >
          <div
            class="rounded-full transition-all"
            :class="videoRecording ? 'w-8 h-8 bg-white' : 'w-16 h-16 bg-red-500'"
          ></div>
        </button>
      </div>
    </div>

    <!-- 视频预览确认 -->
    <div v-if="recordedVideoUrl" class="fixed inset-0 bg-black z-50 flex flex-col">
      <div class="flex-1 relative">
        <video
          :src="recordedVideoUrl"
          class="w-full h-full object-contain"
          controls
          playsinline
        ></video>
      </div>
      <div class="p-4 flex justify-center gap-8 bg-black">
        <button
          @click="cancelRecordedVideo"
          class="px-6 py-3 rounded-full bg-gray-700 text-white"
        >
          重拍
        </button>
        <button
          @click="sendRecordedVideo"
          class="px-6 py-3 rounded-full bg-primary-500 text-white"
          :disabled="sending"
        >
          {{ sending ? '发送中...' : '发送' }}
        </button>
      </div>
    </div>

    <!-- 视频播放器 -->
    <div
      v-if="showVideoPlayer"
      @click="closeVideoPlayer"
      class="fixed inset-0 bg-black z-50 flex items-center justify-center"
    >
      <video
        :src="videoPlayerUrl"
        class="max-w-full max-h-full"
        controls
        autoplay
        playsinline
        @click.stop
      ></video>
      <button
        class="absolute top-4 right-4 w-10 h-10 rounded-full bg-black/50 flex items-center justify-center"
        @click="closeVideoPlayer"
      >
        <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useMessageStore } from '@/stores/message'
import { useAppStore } from '@/stores/app'
import { useCallStore } from '@/stores/call'
import { addMessageListener, removeMessageListener, sendCallInvite } from '@/utils/websocket'
import { uploadImage, uploadVoice, uploadVideo, uploadThumbnail, getResourceUrl } from '@/utils/api'

// image/video URL cache: key -> signedUrl
// For images: stores both thumbnail URL and original URL
// For videos: stores thumbnail URL and video URL
const imageUrlCache = ref(new Map())
// thumbnail URL cache: key -> thumbnailUrl (for images, COS thumbnail processing)
const thumbnailUrlCache = ref(new Map())

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const messageStore = useMessageStore()
const appStore = useAppStore()
const callStore = useCallStore()

const messageListRef = ref(null)
const textInputRef = ref(null)
const imageInputRef = ref(null)

const peerOpenId = ref('')
const peerNickname = ref('')
const peerAvatar = ref('')
const peerInitial = ref('')
const messages = ref([])
const loading = ref(true)
const loadingMore = ref(false)
const sending = ref(false)
const inputText = ref('')
const page = ref(1)
const total = ref(0)
const hasMore = ref(false)
const previewImageUrl = ref('')

// voice recording
const voiceMode = ref(false)
const recording = ref(false)
const recordingDuration = ref(0)
let mediaRecorder = null
let audioChunks = []
let recordingTimer = null
let recordingStartTime = 0

// voice playback
const playingVoiceId = ref(null)
let currentAudio = null

// video recording
const showVideoRecorder = ref(false)
const videoRecorderRef = ref(null)
const videoPreviewRef = ref(null)
const recordedVideoUrl = ref('')
const recordedVideoBlob = ref(null)
const videoRecording = ref(false)
const videoRecordingDuration = ref(0)
let videoMediaRecorder = null
let videoChunks = []
let videoRecordingTimer = null
let videoStream = null

// video playback
const playingVideoId = ref(null)
const videoPlayerUrl = ref('')
const showVideoPlayer = ref(false)

// emoji picker
const showEmojiPicker = ref(false)
const showMorePanel = ref(false)
const emojiList = [
  '😀', '😃', '😄', '😁', '😅', '😂', '🤣', '😊',
  '😇', '🙂', '🙃', '😉', '😌', '😍', '🥰', '😘',
  '😗', '😙', '😚', '😋', '😛', '😜', '🤪', '😝',
  '🤑', '🤗', '🤭', '🤫', '🤔', '🤐', '🤨', '😐',
  '😑', '😶', '😏', '😒', '🙄', '😬', '🤥', '😌',
  '😔', '😪', '🤤', '😴', '😷', '🤒', '🤕', '🤢',
  '🤮', '🤧', '🥵', '🥶', '🥴', '😵', '🤯', '🤠',
  '🥳', '😎', '🤓', '🧐', '😕', '😟', '🙁', '☹️',
  '😮', '😯', '😲', '😳', '🥺', '😦', '😧', '😨',
  '😰', '😥', '😢', '😭', '😱', '😖', '😣', '😞',
  '😓', '😩', '😫', '🥱', '😤', '😡', '😠', '🤬',
  '👍', '👎', '👏', '🙌', '👋', '🤝', '💪', '❤️'
]

// use open_id for comparison
const currentUserOpenId = computed(() => userStore.user?.open_id || '')
const myInitial = computed(() => userStore.user?.nickname?.charAt(0) || 'U')

// messages are fetched DESC, reverse for display (oldest first)
const reversedMessages = computed(() => [...messages.value].reverse())

// keyboard height handling for mobile
let initialViewportHeight = 0

onMounted(async () => {
  peerOpenId.value = route.params.peerId
  peerNickname.value = route.query.nickname || '用户'
  peerAvatar.value = route.query.avatar || ''
  peerInitial.value = peerNickname.value.charAt(0)

  await fetchMessages()
  // 等待 DOM 渲染完成后再滚动到底部
  await nextTick()
  scrollToBottom()
  // 确保滚动生效，延迟再执行一次
  setTimeout(() => scrollToBottom(), 100)

  // mark messages as read
  messageStore.markAsRead(peerOpenId.value)

  // listen for new messages from websocket
  addMessageListener(handleNewMessage)

  // listen for keyboard events via visualViewport API
  if (window.visualViewport) {
    initialViewportHeight = window.visualViewport.height
    window.visualViewport.addEventListener('resize', handleViewportResize)
  }
})

onUnmounted(() => {
  removeMessageListener(handleNewMessage)
  messageStore.clearCurrentMessages()
  // stop audio playback when leaving chat
  if (currentAudio) {
    currentAudio.pause()
    currentAudio = null
    playingVoiceId.value = null
  }
  // remove viewport listeners
  if (window.visualViewport) {
    window.visualViewport.removeEventListener('resize', handleViewportResize)
  }
})

// 发起语音通话
function startVoiceCall() {
  initiateCall('audio')
}

// 发起视频通话
function startVideoCall() {
  initiateCall('video')
}

// 通用发起通话逻辑
function initiateCall(type) {
  // 检查是否已在通话中
  if (callStore.isInCall) {
    appStore.showToast('当前正在通话中', 'error')
    return
  }

  // 发起通话
  const callId = callStore.initiateCall(
    peerOpenId.value,
    peerNickname.value,
    peerAvatar.value,
    type
  )

  // 发送通话邀请信令
  const callerInfo = {
    nickname: userStore.user?.nickname || '用户',
    avatar: userStore.user?.avatar || ''
  }
  sendCallInvite(peerOpenId.value, callId, type, callerInfo)

  // 跳转到通话页面
  router.push('/call')
}

// handle viewport resize (keyboard show/hide)
function handleViewportResize() {
  if (!window.visualViewport) return
  const currentHeight = window.visualViewport.height
  const diff = initialViewportHeight - currentHeight
  // keyboard is showing if viewport height decreased significantly
  if (diff > 100) {
    nextTick(() => scrollToBottom())
  }
}

function handleNewMessage(msg) {
  // only add if it's from current conversation
  // msg.sender_id and msg.receiver_id are open_ids from backend
  // case 1: peer sends to me -> sender_id = peer, receiver_id = me
  // case 2: I send to peer (won't receive via websocket, already added via API)
  const isFromPeer = msg.sender_id === peerOpenId.value && msg.receiver_id === currentUserOpenId.value
  const isToPeer = msg.sender_id === currentUserOpenId.value && msg.receiver_id === peerOpenId.value
  
  if (isFromPeer || isToPeer) {
    // avoid duplicate: check if message already exists
    const exists = messages.value.some(m => m.id === msg.id)
    if (exists) return
    
    messages.value.unshift(msg)
    // load image/voice URL if needed
    if (msg.msg_type === 2 || msg.msg_type === 3) {
      loadMediaUrls([msg])
    }
    nextTick(() => scrollToBottom())
    // mark as read immediately if from peer
    if (isFromPeer) {
      messageStore.markAsRead(peerOpenId.value)
    }
  }
}

async function fetchMessages() {
  loading.value = true
  try {
    const data = await messageStore.fetchMessages(peerOpenId.value, page.value, 20)
    messages.value = data.list || []
    total.value = data.total
    hasMore.value = messages.value.length < total.value
    // load media URLs for image/voice messages
    await loadMediaUrls(messages.value)
  } finally {
    loading.value = false
  }
}

async function loadMoreMessages() {
  if (loadingMore.value || !hasMore.value) return
  loadingMore.value = true
  page.value++
  try {
    const data = await messageStore.fetchMessages(peerOpenId.value, page.value, 20)
    const newMessages = data.list || []
    messages.value = [...messages.value, ...newMessages]
    hasMore.value = messages.value.length < data.total
    // load media URLs for new messages
    await loadMediaUrls(newMessages)
  } finally {
    loadingMore.value = false
  }
}

async function sendMessage() {
  const text = inputText.value.trim()
  if (!text || sending.value) return

  sending.value = true
  try {
    const msg = await messageStore.sendTextMessage(peerOpenId.value, text)
    messages.value.unshift(msg)
    inputText.value = ''
    nextTick(() => {
      scrollToBottom()
      autoResize()
    })
  } finally {
    sending.value = false
  }
}

function selectImage() {
  imageInputRef.value?.click()
}

async function handleImageSelected(e) {
  const file = e.target.files?.[0]
  if (!file) return

  // validate file type and size
  if (!file.type.startsWith('image/')) {
    appStore.showToast('请选择图片文件', 'error')
    return
  }
  if (file.size > 100 * 1024 * 1024) {
    appStore.showToast('图片大小不能超过100MB', 'error')
    return
  }

  sending.value = true
  try {
    // upload image to server first, get object key
    const imageKey = await uploadImage(file, 'images')
    // send image message with the object key
    const msg = await messageStore.sendImageMessage(peerOpenId.value, imageKey)
    // get signed URL for display
    const signedUrl = await getResourceUrl(imageKey)
    imageUrlCache.value.set(imageKey, signedUrl)
    messages.value.unshift(msg)
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '图片发送失败', 'error')
  } finally {
    sending.value = false
  }

  // clear input
  e.target.value = ''
}

function scrollToBottom() {
  nextTick(() => {
    if (messageListRef.value) {
      messageListRef.value.scrollTop = messageListRef.value.scrollHeight
    }
  })
}

function handleScroll() {
  // could implement pull-to-load-more here
}

function autoResize() {
  const textarea = textInputRef.value
  if (textarea) {
    textarea.style.height = 'auto'
    textarea.style.height = Math.min(textarea.scrollHeight, 100) + 'px'
  }
}

// handle input focus - just scroll messages to bottom
function handleInputFocus() {
  // simple delay to let keyboard appear, then scroll messages
  setTimeout(() => {
    scrollToBottom()
  }, 100)
}

function handleInputBlur() {
  // no-op, let the layout handle itself
}

function shouldShowTime(msg, index) {
  if (index === 0) return true
  const prev = reversedMessages.value[index - 1]
  if (!prev) return true
  // show time if more than 5 minutes apart
  const prevTime = new Date(prev.created_at).getTime()
  const currTime = new Date(msg.created_at).getTime()
  return currTime - prevTime > 5 * 60 * 1000
}

function formatMessageTime(time) {
  const date = new Date(time)
  const now = new Date()
  const isToday = date.toDateString() === now.toDateString()
  const hour = date.getHours().toString().padStart(2, '0')
  const minute = date.getMinutes().toString().padStart(2, '0')

  if (isToday) {
    return `${hour}:${minute}`
  }
  const month = date.getMonth() + 1
  const day = date.getDate()
  return `${month}月${day}日 ${hour}:${minute}`
}

// 解析通话记录内容
function parseCallRecord(content) {
  try {
    return JSON.parse(content)
  } catch {
    return { call_type: 'audio', duration: 0, status: 'completed' }
  }
}

// 格式化通话记录显示
function formatCallRecord(msg) {
  const record = parseCallRecord(msg.content)
  const isSelf = msg.sender_id === currentUserOpenId.value
  const callTypeName = record.call_type === 'video' ? '视频' : '语音'
  
  switch (record.status) {
    case 'completed':
      return `${callTypeName}通话 ${formatCallDuration(record.duration)}`
    case 'missed':
      return isSelf ? `对方未接听` : `未接听`
    case 'rejected':
      return isSelf ? `对方已拒绝` : `已拒绝`
    case 'cancelled':
      return isSelf ? `已取消` : `对方已取消`
    default:
      return `${callTypeName}通话`
  }
}

// 格式化通话时长
function formatCallDuration(seconds) {
  if (!seconds || seconds <= 0) return ''
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  if (mins > 0) {
    return `${mins}分${secs}秒`
  }
  return `${secs}秒`
}

// 点击图片预览原图
async function previewImage(key) {
  // 如果已有原图URL，直接显示
  if (imageUrlCache.value.get(key)) {
    previewImageUrl.value = imageUrlCache.value.get(key)
    return
  }
  // 否则先加载原图
  try {
    appStore.showToast('加载原图中...', 'info')
    const url = await getResourceUrl(key)
    imageUrlCache.value.set(key, url)
    previewImageUrl.value = url
  } catch (err) {
    appStore.showToast('图片加载失败', 'error')
  }
}

// 解析视频消息内容
function parseVideoContent(content) {
  try {
    return JSON.parse(content)
  } catch {
    return { key: '', thumbnail: '', duration: 0, width: 0, height: 0 }
  }
}

// 格式化视频时长
function formatVideoDuration(seconds) {
  if (!seconds || seconds <= 0) return '0:00'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

// 获取视频缩略图URL
function getVideoThumbnail(msg) {
  const content = parseVideoContent(msg.content)
  if (content.thumbnail) {
    return thumbnailUrlCache.value.get(content.thumbnail) || ''
  }
  return ''
}

// 播放视频
async function playVideo(msg) {
  const content = parseVideoContent(msg.content)
  if (!content.key) return
  
  // 获取视频URL
  let videoUrl = imageUrlCache.value.get(content.key)
  if (!videoUrl) {
    try {
      appStore.showToast('加载视频中...', 'info')
      videoUrl = await getResourceUrl(content.key)
      imageUrlCache.value.set(content.key, videoUrl)
    } catch (err) {
      appStore.showToast('视频加载失败', 'error')
      return
    }
  }
  
  videoPlayerUrl.value = videoUrl
  showVideoPlayer.value = true
  playingVideoId.value = msg.id
}

// 关闭视频播放器
function closeVideoPlayer() {
  showVideoPlayer.value = false
  videoPlayerUrl.value = ''
  playingVideoId.value = null
}

// load signed URLs for image/video messages (thumbnails for list view)
async function loadMediaUrls(msgList) {
  for (const msg of msgList) {
    try {
      if (msg.msg_type === 2) {
        // 图片消息 - 加载缩略图（使用COS图片处理）
        if (!thumbnailUrlCache.value.has(msg.content)) {
          const url = await getResourceUrl(msg.content)
          // 使用COS图片处理生成缩略图 (添加参数限制尺寸)
          const thumbUrl = url.includes('?') ? `${url}&imageMogr2/thumbnail/400x` : `${url}?imageMogr2/thumbnail/400x`
          thumbnailUrlCache.value.set(msg.content, thumbUrl)
        }
      } else if (msg.msg_type === 3) {
        // 语音消息
        if (!imageUrlCache.value.has(msg.content)) {
          const url = await getResourceUrl(msg.content)
          imageUrlCache.value.set(msg.content, url)
        }
      } else if (msg.msg_type === 6) {
        // 视频消息 - 加载缩略图
        const content = parseVideoContent(msg.content)
        if (content.thumbnail && !thumbnailUrlCache.value.has(content.thumbnail)) {
          const url = await getResourceUrl(content.thumbnail)
          thumbnailUrlCache.value.set(content.thumbnail, url)
        }
      }
    } catch (err) {
      console.error('Failed to load media URL:', err)
    }
  }
}

function onImageLoad() {
  // trigger reactivity update
  imageUrlCache.value = new Map(imageUrlCache.value)
}

function goBack() {
  router.back()
}

function goToUserProfile() {
  if (peerOpenId.value) {
    router.push(`/user/${peerOpenId.value}`)
  }
}

// Voice recording functions
function toggleVoiceMode() {
  voiceMode.value = !voiceMode.value
  showEmojiPicker.value = false
  showMorePanel.value = false
}

async function startRecording() {
  if (recording.value) return
  
  try {
    // check microphone permission first
    if (navigator.permissions) {
      const permission = await navigator.permissions.query({ name: 'microphone' })
      if (permission.state === 'prompt') {
        // permission not yet granted, request it without starting recording
        appStore.showToast('请先允许麦克风权限', 'info')
        try {
          const testStream = await navigator.mediaDevices.getUserMedia({ audio: true })
          testStream.getTracks().forEach(track => track.stop())
          appStore.showToast('权限已授予，请再次按住录音', 'success')
        } catch {
          appStore.showToast('无法访问麦克风', 'error')
        }
        return
      } else if (permission.state === 'denied') {
        appStore.showToast('麦克风权限被拒绝，请在设置中开启', 'error')
        return
      }
    }
    
    const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
    
    // detect supported audio format for cross-browser compatibility
    let mimeType = 'audio/webm'
    if (MediaRecorder.isTypeSupported('audio/webm')) {
      mimeType = 'audio/webm'
    } else if (MediaRecorder.isTypeSupported('audio/mp4')) {
      mimeType = 'audio/mp4'
    } else if (MediaRecorder.isTypeSupported('audio/ogg')) {
      mimeType = 'audio/ogg'
    } else if (MediaRecorder.isTypeSupported('audio/wav')) {
      mimeType = 'audio/wav'
    }
    
    mediaRecorder = new MediaRecorder(stream, { mimeType })
    audioChunks = []
    
    mediaRecorder.ondataavailable = (e) => {
      audioChunks.push(e.data)
    }
    
    mediaRecorder.onstop = async () => {
      stream.getTracks().forEach(track => track.stop())
      
      const duration = Math.round((Date.now() - recordingStartTime) / 1000)
      if (duration < 1) {
        appStore.showToast('录音时间太短', 'error')
        return
      }
      
      // use the actual mimeType from mediaRecorder
      const actualMimeType = mediaRecorder.mimeType || mimeType
      const ext = actualMimeType.includes('mp4') ? '.m4a' : 
                  actualMimeType.includes('ogg') ? '.ogg' : 
                  actualMimeType.includes('wav') ? '.wav' : '.webm'
      const audioBlob = new Blob(audioChunks, { type: actualMimeType })
      await sendVoiceMessage(audioBlob, duration, ext)
    }
    
    mediaRecorder.start()
    recording.value = true
    recordingStartTime = Date.now()
    recordingDuration.value = 0
    
    recordingTimer = setInterval(() => {
      recordingDuration.value = Math.round((Date.now() - recordingStartTime) / 1000)
      if (recordingDuration.value >= 60) {
        stopRecording()
      }
    }, 100)
  } catch (err) {
    appStore.showToast('无法访问麦克风', 'error')
  }
}

function stopRecording() {
  if (!recording.value || !mediaRecorder) return
  
  clearInterval(recordingTimer)
  recording.value = false
  
  if (mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
  }
}

function cancelRecording() {
  if (!recording.value) return
  
  clearInterval(recordingTimer)
  recording.value = false
  
  if (mediaRecorder && mediaRecorder.state === 'recording') {
    mediaRecorder.stop()
    audioChunks = []
  }
}

async function sendVoiceMessage(audioBlob, duration, ext = '.webm') {
  sending.value = true
  try {
    // upload voice file with correct extension
    const filename = `voice${ext}`
    const file = new File([audioBlob], filename, { type: audioBlob.type })
    const voiceKey = await uploadVoice(file)
    // send voice message
    const msg = await messageStore.sendVoiceMessage(peerOpenId.value, voiceKey, duration)
    // get signed URL for playback, use msg.content as cache key for consistency
    const signedUrl = await getResourceUrl(msg.content)
    imageUrlCache.value.set(msg.content, signedUrl)
    messages.value.unshift(msg)
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '语音发送失败', 'error')
  } finally {
    sending.value = false
  }
}

// Voice playback
async function playVoice(msg) {
  const url = imageUrlCache.value.get(msg.content)
  if (!url) {
    try {
      const signedUrl = await getResourceUrl(msg.content)
      imageUrlCache.value.set(msg.content, signedUrl)
      playVoiceUrl(msg.id, signedUrl)
    } catch (err) {
      appStore.showToast('语音加载失败', 'error')
    }
    return
  }
  playVoiceUrl(msg.id, url)
}

function playVoiceUrl(msgId, url) {
  // stop current playing
  if (currentAudio) {
    currentAudio.pause()
    currentAudio = null
    if (playingVoiceId.value === msgId) {
      playingVoiceId.value = null
      return
    }
  }
  
  currentAudio = new Audio(url)
  playingVoiceId.value = msgId
  
  currentAudio.onended = () => {
    playingVoiceId.value = null
    currentAudio = null
  }
  
  currentAudio.onerror = (e) => {
    console.error('Audio playback error:', e, 'url:', url)
    playingVoiceId.value = null
    currentAudio = null
    appStore.showToast('语音播放失败', 'error')
  }
  
  currentAudio.play().catch((err) => {
    console.error('Audio play() rejected:', err)
    playingVoiceId.value = null
    currentAudio = null
    appStore.showToast('语音播放失败', 'error')
  })
}

// Emoji functions
function toggleEmojiPicker() {
  showEmojiPicker.value = !showEmojiPicker.value
  showMorePanel.value = false
  voiceMode.value = false
}

function toggleMorePanel() {
  showMorePanel.value = !showMorePanel.value
  showEmojiPicker.value = false
  voiceMode.value = false
}

async function sendEmoji(emoji) {
  sending.value = true
  try {
    const msg = await messageStore.sendEmojiMessage(peerOpenId.value, emoji)
    messages.value.unshift(msg)
    showEmojiPicker.value = false
    nextTick(() => scrollToBottom())
  } catch (err) {
    appStore.showToast(err.message || '表情发送失败', 'error')
  } finally {
    sending.value = false
  }
}

// Video recording functions
async function openVideoRecorder() {
  showMorePanel.value = false
  try {
    // 请求摄像头和麦克风权限
    videoStream = await navigator.mediaDevices.getUserMedia({
      video: { facingMode: 'environment' },
      audio: true
    })
    showVideoRecorder.value = true
    nextTick(() => {
      if (videoPreviewRef.value) {
        videoPreviewRef.value.srcObject = videoStream
      }
    })
  } catch (err) {
    console.error('Failed to access camera:', err)
    appStore.showToast('无法访问摄像头', 'error')
  }
}

function closeVideoRecorder() {
  if (videoStream) {
    videoStream.getTracks().forEach(track => track.stop())
    videoStream = null
  }
  showVideoRecorder.value = false
  videoRecording.value = false
  videoRecordingDuration.value = 0
  if (videoRecordingTimer) {
    clearInterval(videoRecordingTimer)
    videoRecordingTimer = null
  }
}

function startVideoRecording() {
  if (videoRecording.value || !videoStream) return
  
  // 检测支持的视频格式
  let mimeType = 'video/webm'
  if (MediaRecorder.isTypeSupported('video/webm;codecs=vp9')) {
    mimeType = 'video/webm;codecs=vp9'
  } else if (MediaRecorder.isTypeSupported('video/webm;codecs=vp8')) {
    mimeType = 'video/webm;codecs=vp8'
  } else if (MediaRecorder.isTypeSupported('video/mp4')) {
    mimeType = 'video/mp4'
  }
  
  videoMediaRecorder = new MediaRecorder(videoStream, { mimeType })
  videoChunks = []
  
  videoMediaRecorder.ondataavailable = (e) => {
    if (e.data.size > 0) {
      videoChunks.push(e.data)
    }
  }
  
  videoMediaRecorder.onstop = () => {
    const blob = new Blob(videoChunks, { type: mimeType })
    recordedVideoBlob.value = blob
    recordedVideoUrl.value = URL.createObjectURL(blob)
    closeVideoRecorder()
  }
  
  videoMediaRecorder.start()
  videoRecording.value = true
  videoRecordingDuration.value = 0
  
  const startTime = Date.now()
  videoRecordingTimer = setInterval(() => {
    videoRecordingDuration.value = Math.round((Date.now() - startTime) / 1000)
    // 最长录制60秒
    if (videoRecordingDuration.value >= 60) {
      stopVideoRecording()
    }
  }, 100)
}

function stopVideoRecording() {
  if (!videoRecording.value || !videoMediaRecorder) return
  
  if (videoRecordingTimer) {
    clearInterval(videoRecordingTimer)
    videoRecordingTimer = null
  }
  
  if (videoMediaRecorder.state === 'recording') {
    videoMediaRecorder.stop()
  }
  videoRecording.value = false
}

function cancelRecordedVideo() {
  if (recordedVideoUrl.value) {
    URL.revokeObjectURL(recordedVideoUrl.value)
  }
  recordedVideoUrl.value = ''
  recordedVideoBlob.value = null
  openVideoRecorder()
}

async function sendRecordedVideo() {
  if (!recordedVideoBlob.value) return
  
  sending.value = true
  try {
    // 生成视频缩略图
    const thumbnailBlob = await generateVideoThumbnail(recordedVideoUrl.value)
    
    // 上传缩略图
    const thumbFile = new File([thumbnailBlob], 'thumb.jpg', { type: 'image/jpeg' })
    const thumbKey = await uploadThumbnail(thumbFile)
    
    // 上传视频
    const ext = recordedVideoBlob.value.type.includes('mp4') ? '.mp4' : '.webm'
    const videoFile = new File([recordedVideoBlob.value], `video${ext}`, { type: recordedVideoBlob.value.type })
    const videoKey = await uploadVideo(videoFile)
    
    // 发送视频消息
    const msg = await messageStore.sendVideoMessage(peerOpenId.value, {
      key: videoKey,
      thumbnail: thumbKey,
      duration: videoRecordingDuration.value || Math.round(recordedVideoBlob.value.size / 50000), // 估算时长
      width: 0,
      height: 0
    })
    
    // 缓存缩略图URL
    const thumbUrl = await getResourceUrl(thumbKey)
    thumbnailUrlCache.value.set(thumbKey, thumbUrl)
    
    messages.value.unshift(msg)
    nextTick(() => scrollToBottom())
    
    // 清理
    if (recordedVideoUrl.value) {
      URL.revokeObjectURL(recordedVideoUrl.value)
    }
    recordedVideoUrl.value = ''
    recordedVideoBlob.value = null
  } catch (err) {
    appStore.showToast(err.message || '视频发送失败', 'error')
  } finally {
    sending.value = false
  }
}

// 从视频生成缩略图
function generateVideoThumbnail(videoUrl) {
  return new Promise((resolve, reject) => {
    const video = document.createElement('video')
    video.src = videoUrl
    video.crossOrigin = 'anonymous'
    video.muted = true
    
    video.onloadeddata = () => {
      video.currentTime = 0.1 // 取第0.1秒的帧
    }
    
    video.onseeked = () => {
      const canvas = document.createElement('canvas')
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      const ctx = canvas.getContext('2d')
      ctx.drawImage(video, 0, 0)
      
      canvas.toBlob((blob) => {
        if (blob) {
          resolve(blob)
        } else {
          reject(new Error('Failed to generate thumbnail'))
        }
      }, 'image/jpeg', .8)
    }
    
    video.onerror = () => {
      reject(new Error('Failed to load video'))
    }
    
    video.load()
  })
}
</script>

<style scoped>
.chat-page {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  background-color: #ededed;
  /* use dvh for better mobile keyboard handling */
  height: 100dvh;
  height: 100vh; /* fallback */
  overflow: hidden;
}

@supports (height: 100dvh) {
  .chat-page {
    height: 100dvh;
  }
}

.chat-page--dark {
  background-color: #111827;
}

.chat-header {
  flex-shrink: 0;
  position: relative;
  z-index: 10;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  overflow-x: hidden;
  padding: 1rem;
  -webkit-overflow-scrolling: touch;
  overscroll-behavior: contain;
}

/* 微信风格输入栏 */
.chat-input {
  flex-shrink: 0;
  position: relative;
  z-index: 10;
  background-color: #f7f7f7;
  border-top: 1px solid #e5e5e5;
}

.chat-input--dark {
  background-color: #1f2937;
  border-top-color: #374151;
}

.chat-input__main {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
}

.chat-input__icon-btn {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #1f1f1f;
  transition: opacity .15s;
}

.chat-input__icon-btn:active {
  opacity: .6;
}

.chat-input__icon-btn--dark {
  color: #9ca3af;
}

.chat-input__icon-btn--active {
  color: #07c160;
}

.chat-input__field-wrap {
  flex: 1;
  min-width: 0;
}

.chat-input__textarea {
  display: block;
  width: 100%;
  min-height: 36px;
  max-height: 100px;
  padding: 8px 12px;
  font-size: 16px;
  line-height: 1.4;
  border: none;
  border-radius: 4px;
  background-color: #fff;
  resize: none;
  outline: none;
  -webkit-appearance: none;
  /* prevent iOS zoom on focus */
  touch-action: manipulation;
}

.chat-input__textarea::placeholder {
  color: #b2b2b2;
}

.chat-input__textarea--dark {
  background-color: #374151;
  color: #fff;
}

.chat-input__textarea--dark::placeholder {
  color: #6b7280;
}

.chat-input__voice-btn {
  width: 100%;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 500;
  color: #1f1f1f;
  background-color: #fff;
  border-radius: 4px;
  user-select: none;
  cursor: pointer;
  transition: background-color .15s;
}

.chat-input__voice-btn:active,
.chat-input__voice-btn--recording {
  background-color: #c9c9c9;
}

.chat-input__voice-btn--dark {
  background-color: #374151;
  color: #e5e7eb;
}

.chat-input__voice-btn--dark:active,
.chat-input__voice-btn--dark.chat-input__voice-btn--recording {
  background-color: #4b5563;
}

.chat-input__send-btn {
  flex-shrink: 0;
  height: 36px;
  padding: 0 14px;
  font-size: 15px;
  font-weight: 500;
  color: #fff;
  background-color: #07c160;
  border-radius: 4px;
  transition: opacity .15s;
}

.chat-input__send-btn:active {
  opacity: .8;
}

.chat-input__send-btn:disabled {
  opacity: .5;
}

/* 面板样式 */
.chat-input__panel {
  background-color: #f7f7f7;
  border-top: 1px solid #e5e5e5;
  animation: slideUp .2s ease-out;
}

.chat-input__panel--dark {
  background-color: #1f2937;
  border-top-color: #374151;
}

.chat-input__emoji-item {
  width: 100%;
  aspect-ratio: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 28px;
  border-radius: 8px;
  transition: background-color .15s;
}

.chat-input__emoji-item:active {
  background-color: rgba(0, 0, 0, .05);
}

/* 更多功能面板 */
.chat-input__more-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.chat-input__more-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff;
  border-radius: 12px;
  transition: opacity .15s;
}

.chat-input__more-icon:active {
  opacity: .7;
}

.chat-input__more-text {
  font-size: 12px;
  color: #576b95;
}

/* 录音状态遮罩 */
.chat-input__recording-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: rgba(0, 0, 0, .5);
  z-index: 100;
}

.chat-input__recording-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px 32px;
  background-color: rgba(0, 0, 0, .8);
  border-radius: 12px;
}

.chat-input__recording-wave {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  height: 40px;
}

.chat-input__recording-wave span {
  width: 4px;
  height: 20px;
  background-color: #07c160;
  border-radius: 2px;
  animation: wave 1s ease-in-out infinite;
}

.chat-input__recording-wave span:nth-child(1) { animation-delay: 0s; }
.chat-input__recording-wave span:nth-child(2) { animation-delay: .1s; }
.chat-input__recording-wave span:nth-child(3) { animation-delay: .2s; }
.chat-input__recording-wave span:nth-child(4) { animation-delay: .3s; }
.chat-input__recording-wave span:nth-child(5) { animation-delay: .4s; }

@keyframes wave {
  0%, 100% {
    height: 20px;
  }
  50% {
    height: 40px;
  }
}

.chat-input__recording-time {
  font-size: 48px;
  font-weight: 300;
  color: #fff;
}

.chat-input__recording-tip {
  font-size: 14px;
  color: rgba(255, 255, 255, .7);
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 语音播放图标 - 微信风格声波动画 */
.voice-icon {
  position: relative;
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: flex-start;
  gap: 2px;
  flex-shrink: 0;
}

.voice-icon span {
  display: block;
  width: 3px;
  background-color: currentColor;
  border-radius: 1px;
}

.voice-icon span:nth-child(1) {
  height: 6px;
}

.voice-icon span:nth-child(2) {
  height: 10px;
}

.voice-icon span:nth-child(3) {
  height: 14px;
}

/* 自己发送的语音 - 白色 + 图标靠右 */
.voice-icon--self {
  color: #fff;
  justify-content: flex-end;
}

/* 播放中动画 */
.voice-icon--playing span {
  animation: voiceWave 1s ease-in-out infinite;
}

.voice-icon--playing span:nth-child(1) {
  animation-delay: 0s;
}

.voice-icon--playing span:nth-child(2) {
  animation-delay: .15s;
}

.voice-icon--playing span:nth-child(3) {
  animation-delay: .3s;
}

@keyframes voiceWave {
  0%, 100% {
    opacity: .3;
  }
  50% {
    opacity: 1;
  }
}

/* 通话记录消息 */
.call-record {
  cursor: pointer;
  transition: opacity .15s;
}

.call-record:active {
  opacity: .7;
}

.call-record__icon {
  flex-shrink: 0;
}

.call-record__info {
  white-space: nowrap;
}

/* 图片消息 */
.image-message {
  display: inline-block;
  border-radius: 12px;
  overflow: hidden;
  background-color: rgba(0, 0, 0, .05);
}

.image-message img {
  display: block;
}

/* 视频消息 */
.video-message {
  display: inline-block;
  border-radius: 12px;
  overflow: hidden;
  background-color: #000;
}

.video-message img {
  display: block;
}
</style>
