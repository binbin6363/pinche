import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useTripStore = defineStore('trip', () => {
  const trips = ref([])
  const myTrips = ref([])
  const total = ref(0)
  const loading = ref(false)

  async function fetchTrips(params = {}) {
    loading.value = true
    try {
      const result = await api.get('/trips', { params })
      trips.value = result.list || []
      total.value = result.total
      return result
    } finally {
      loading.value = false
    }
  }

  async function fetchMyTrips() {
    loading.value = true
    try {
      const result = await api.get('/trips/my')
      myTrips.value = result || []
      return result
    } finally {
      loading.value = false
    }
  }

  async function getTripById(id) {
    return await api.get(`/trips/${id}`)
  }

  // get my trip detail with grabbers list
  async function getMyTripDetail(id) {
    return await api.get(`/trips/my/${id}`)
  }

  async function createTrip(data) {
    const result = await api.post('/trips', data)
    return result
  }

  // update trip (images/remark can be updated directly, location/time requires review)
  async function updateTrip(id, data) {
    return await api.put(`/trips/${id}`, data)
  }

  async function cancelTrip(id) {
    return await api.put(`/trips/${id}/cancel`)
  }

  async function completeTrip(id) {
    return await api.put(`/trips/${id}/complete`)
  }

  async function deleteTrip(id) {
    return await api.delete(`/trips/${id}`)
  }

  async function grabTrip(id, message = '') {
    return await api.post(`/trips/${id}/grab`, { message })
  }

  return {
    trips,
    myTrips,
    total,
    loading,
    fetchTrips,
    fetchMyTrips,
    getTripById,
    getMyTripDetail,
    createTrip,
    updateTrip,
    cancelTrip,
    completeTrip,
    deleteTrip,
    grabTrip
  }
})
