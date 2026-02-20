import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useFriendStore = defineStore('friend', () => {
  const friends = ref([])
  const friendRequests = ref([])
  const friendCount = ref(0)
  const requestCount = ref(0)
  const loading = ref(false)

  // get friends list
  async function fetchFriends() {
    loading.value = true
    const result = await api.get('/friends')
    friends.value = result.list || []
    friendCount.value = result.total || 0
    loading.value = false
    return result
  }

  // get friend requests
  async function fetchFriendRequests() {
    loading.value = true
    const result = await api.get('/friends/requests')
    friendRequests.value = result.list || []
    requestCount.value = result.total || 0
    loading.value = false
    return result
  }

  // get friend and request counts
  async function fetchFriendCount() {
    const result = await api.get('/friends/count')
    friendCount.value = result.friend_count || 0
    requestCount.value = result.request_count || 0
    return result
  }

  // send friend request
  async function sendFriendRequest(friendOpenId, message = '') {
    await api.post('/friends/request', {
      friend_id: friendOpenId,
      message
    })
    // refresh counts
    await fetchFriendCount()
  }

  // accept friend request
  async function acceptFriendRequest(requestId) {
    await api.post(`/friends/requests/${requestId}/accept`)
    // remove from local list
    friendRequests.value = friendRequests.value.filter(r => r.id !== requestId)
    requestCount.value = Math.max(0, requestCount.value - 1)
    // refresh friends list
    await fetchFriends()
  }

  // reject friend request
  async function rejectFriendRequest(requestId) {
    await api.post(`/friends/requests/${requestId}/reject`)
    // remove from local list
    friendRequests.value = friendRequests.value.filter(r => r.id !== requestId)
    requestCount.value = Math.max(0, requestCount.value - 1)
  }

  // cancel sent friend request
  async function cancelFriendRequest(requestId) {
    await api.delete(`/friends/requests/${requestId}`)
  }

  // delete friend
  async function deleteFriend(friendOpenId) {
    await api.delete(`/friends/${friendOpenId}`)
    // remove from local list
    friends.value = friends.value.filter(f => f.friend?.open_id !== friendOpenId)
    friendCount.value = Math.max(0, friendCount.value - 1)
  }

  // get user public profile
  async function getUserProfile(userOpenId) {
    const result = await api.get(`/users/${userOpenId}/profile`)
    return result
  }

  // clear store data
  function clearData() {
    friends.value = []
    friendRequests.value = []
    friendCount.value = 0
    requestCount.value = 0
  }

  return {
    friends,
    friendRequests,
    friendCount,
    requestCount,
    loading,
    fetchFriends,
    fetchFriendRequests,
    fetchFriendCount,
    sendFriendRequest,
    acceptFriendRequest,
    rejectFriendRequest,
    cancelFriendRequest,
    deleteFriend,
    getUserProfile,
    clearData
  }
})
