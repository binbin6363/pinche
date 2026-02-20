/**
 * WebRTC 封装类
 * 管理 RTCPeerConnection 生命周期、媒体流获取、SDP/ICE 协商
 */

// STUN 服务器配置
const ICE_SERVERS = [
  { urls: 'stun:stun.l.google.com:19302' },
  { urls: 'stun:stun1.l.google.com:19302' },
  { urls: 'stun:stun2.l.google.com:19302' }
]

export class WebRTCManager {
  constructor(options = {}) {
    this.pc = null
    this.localStream = null
    this.remoteStream = null
    this.callType = 'audio' // 'audio' | 'video'
    this.pendingCandidates = [] // 缓存在 remote description 设置前收到的 ICE candidates
    this.hasRemoteDescription = false
    
    // 回调函数
    this.onLocalStream = options.onLocalStream || null
    this.onRemoteStream = options.onRemoteStream || null
    this.onIceCandidate = options.onIceCandidate || null
    this.onConnectionStateChange = options.onConnectionStateChange || null
    this.onError = options.onError || null
  }

  /**
   * 初始化 PeerConnection
   */
  initPeerConnection() {
    if (this.pc) {
      this.pc.close()
    }

    this.pc = new RTCPeerConnection({
      iceServers: ICE_SERVERS
    })

    // ICE candidate 事件
    this.pc.onicecandidate = (event) => {
      if (event.candidate && this.onIceCandidate) {
        this.onIceCandidate(event.candidate)
      }
    }

    // 远端流事件
    this.pc.ontrack = (event) => {
      console.log('WebRTC: Received remote track', event.track.kind, 'enabled:', event.track.enabled)
      
      // 优先使用 event.streams[0]，更可靠
      if (event.streams && event.streams[0]) {
        if (this.remoteStream !== event.streams[0]) {
          this.remoteStream = event.streams[0]
          console.log('WebRTC: Using stream from event, tracks:', this.remoteStream.getTracks().length)
          if (this.onRemoteStream) {
            this.onRemoteStream(this.remoteStream)
          }
        }
      } else {
        // 回退方案：手动创建 MediaStream
        if (!this.remoteStream) {
          this.remoteStream = new MediaStream()
        }
        // 检查是否已有该轨道
        const existingTrack = this.remoteStream.getTracks().find(t => t.id === event.track.id)
        if (!existingTrack) {
          this.remoteStream.addTrack(event.track)
          console.log('WebRTC: Added track to remote stream, total tracks:', this.remoteStream.getTracks().length)
          if (this.onRemoteStream) {
            this.onRemoteStream(this.remoteStream)
          }
        }
      }
    }

    // 连接状态变化
    this.pc.onconnectionstatechange = () => {
      console.log('WebRTC: Connection state:', this.pc.connectionState)
      if (this.onConnectionStateChange) {
        this.onConnectionStateChange(this.pc.connectionState)
      }
    }

    // ICE 连接状态
    this.pc.oniceconnectionstatechange = () => {
      console.log('WebRTC: ICE connection state:', this.pc.iceConnectionState)
    }

    return this.pc
  }

  /**
   * 获取本地媒体流
   * @param {string} type - 'audio' | 'video'
   */
  async getLocalStream(type = 'audio') {
    this.callType = type
    
    const constraints = {
      audio: {
        echoCancellation: true,
        noiseSuppression: true,
        autoGainControl: true
      },
      video: type === 'video' ? {
        width: { ideal: 1280 },
        height: { ideal: 720 },
        facingMode: 'user'
      } : false
    }

    try {
      this.localStream = await navigator.mediaDevices.getUserMedia(constraints)
      console.log('WebRTC: Got local stream', type)
      
      if (this.onLocalStream) {
        this.onLocalStream(this.localStream)
      }
      
      return this.localStream
    } catch (err) {
      console.error('WebRTC: Failed to get local stream:', err)
      if (this.onError) {
        this.onError('media_permission_denied', err.message)
      }
      throw err
    }
  }

  /**
   * 添加本地流到 PeerConnection
   */
  addLocalStreamToPC() {
    if (!this.pc || !this.localStream) {
      console.error('WebRTC: PC or localStream not ready')
      return
    }

    this.localStream.getTracks().forEach(track => {
      console.log('WebRTC: Adding track to PC:', track.kind)
      this.pc.addTrack(track, this.localStream)
    })
  }

  /**
   * 创建 Offer (发起方调用)
   */
  async createOffer() {
    if (!this.pc) {
      throw new Error('PeerConnection not initialized')
    }

    try {
      const offer = await this.pc.createOffer({
        offerToReceiveAudio: true,
        offerToReceiveVideo: this.callType === 'video'
      })
      await this.pc.setLocalDescription(offer)
      console.log('WebRTC: Created offer')
      return offer
    } catch (err) {
      console.error('WebRTC: Failed to create offer:', err)
      if (this.onError) {
        this.onError('create_offer_failed', err.message)
      }
      throw err
    }
  }

  /**
   * 创建 Answer (接收方调用)
   */
  async createAnswer() {
    if (!this.pc) {
      throw new Error('PeerConnection not initialized')
    }

    try {
      const answer = await this.pc.createAnswer()
      await this.pc.setLocalDescription(answer)
      console.log('WebRTC: Created answer')
      return answer
    } catch (err) {
      console.error('WebRTC: Failed to create answer:', err)
      if (this.onError) {
        this.onError('create_answer_failed', err.message)
      }
      throw err
    }
  }

  /**
   * 设置远端 SDP
   * @param {RTCSessionDescriptionInit} sdp
   */
  async setRemoteDescription(sdp) {
    if (!this.pc) {
      throw new Error('PeerConnection not initialized')
    }

    try {
      await this.pc.setRemoteDescription(new RTCSessionDescription(sdp))
      console.log('WebRTC: Set remote description', sdp.type)
      this.hasRemoteDescription = true
      
      // 处理缓存的 ICE candidates
      if (this.pendingCandidates.length > 0) {
        console.log('WebRTC: Processing', this.pendingCandidates.length, 'pending ICE candidates')
        for (const candidate of this.pendingCandidates) {
          try {
            await this.pc.addIceCandidate(new RTCIceCandidate(candidate))
            console.log('WebRTC: Added pending ICE candidate')
          } catch (err) {
            console.error('WebRTC: Failed to add pending ICE candidate:', err)
          }
        }
        this.pendingCandidates = []
      }
    } catch (err) {
      console.error('WebRTC: Failed to set remote description:', err)
      if (this.onError) {
        this.onError('set_remote_description_failed', err.message)
      }
      throw err
    }
  }

  /**
   * 添加 ICE Candidate
   * @param {RTCIceCandidateInit} candidate
   */
  async addIceCandidate(candidate) {
    if (!this.pc) {
      console.warn('WebRTC: PC not ready, caching ICE candidate')
      this.pendingCandidates.push(candidate)
      return
    }

    // 如果 remote description 还未设置，缓存 candidate
    if (!this.hasRemoteDescription) {
      console.log('WebRTC: Remote description not set, caching ICE candidate')
      this.pendingCandidates.push(candidate)
      return
    }

    try {
      await this.pc.addIceCandidate(new RTCIceCandidate(candidate))
      console.log('WebRTC: Added ICE candidate')
    } catch (err) {
      console.error('WebRTC: Failed to add ICE candidate:', err)
    }
  }

  /**
   * 切换麦克风静音
   */
  toggleMute() {
    if (!this.localStream) return false
    
    const audioTrack = this.localStream.getAudioTracks()[0]
    if (audioTrack) {
      audioTrack.enabled = !audioTrack.enabled
      console.log('WebRTC: Mute toggled:', !audioTrack.enabled)
      return !audioTrack.enabled
    }
    return false
  }

  /**
   * 切换摄像头开关
   */
  toggleVideo() {
    if (!this.localStream) return false
    
    const videoTrack = this.localStream.getVideoTracks()[0]
    if (videoTrack) {
      videoTrack.enabled = !videoTrack.enabled
      console.log('WebRTC: Video toggled:', videoTrack.enabled)
      return videoTrack.enabled
    }
    return false
  }

  /**
   * 切换前后摄像头
   */
  async switchCamera() {
    if (!this.localStream || this.callType !== 'video') return

    const videoTrack = this.localStream.getVideoTracks()[0]
    if (!videoTrack) return

    // 获取当前摄像头朝向
    const settings = videoTrack.getSettings()
    const newFacingMode = settings.facingMode === 'user' ? 'environment' : 'user'

    try {
      // 获取新的视频流
      const newStream = await navigator.mediaDevices.getUserMedia({
        video: { facingMode: newFacingMode }
      })

      const newVideoTrack = newStream.getVideoTracks()[0]
      
      // 替换 PC 中的视频轨道
      const sender = this.pc.getSenders().find(s => s.track?.kind === 'video')
      if (sender) {
        await sender.replaceTrack(newVideoTrack)
      }

      // 停止旧轨道
      videoTrack.stop()
      
      // 更新本地流
      this.localStream.removeTrack(videoTrack)
      this.localStream.addTrack(newVideoTrack)

      if (this.onLocalStream) {
        this.onLocalStream(this.localStream)
      }

      console.log('WebRTC: Camera switched to', newFacingMode)
    } catch (err) {
      console.error('WebRTC: Failed to switch camera:', err)
    }
  }

  /**
   * 检查媒体权限
   * @param {string} type - 'audio' | 'video'
   */
  static async checkPermission(type = 'audio') {
    try {
      const constraints = {
        audio: true,
        video: type === 'video'
      }
      const stream = await navigator.mediaDevices.getUserMedia(constraints)
      stream.getTracks().forEach(track => track.stop())
      return true
    } catch (err) {
      console.error('WebRTC: Permission check failed:', err)
      return false
    }
  }

  /**
   * 关闭连接并释放资源
   */
  close() {
    console.log('WebRTC: Closing connection')
    
    // 停止本地流
    if (this.localStream) {
      this.localStream.getTracks().forEach(track => {
        track.stop()
        console.log('WebRTC: Stopped local track:', track.kind)
      })
      this.localStream = null
    }

    // 停止远端流
    if (this.remoteStream) {
      this.remoteStream.getTracks().forEach(track => track.stop())
      this.remoteStream = null
    }

    // 关闭 PeerConnection
    if (this.pc) {
      this.pc.close()
      this.pc = null
    }
    
    // 重置状态
    this.pendingCandidates = []
    this.hasRemoteDescription = false
  }

  /**
   * 获取连接状态
   */
  getConnectionState() {
    return this.pc ? this.pc.connectionState : 'closed'
  }

  /**
   * 获取 ICE 连接状态
   */
  getIceConnectionState() {
    return this.pc ? this.pc.iceConnectionState : 'closed'
  }
}

export default WebRTCManager
