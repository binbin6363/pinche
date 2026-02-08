import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/utils/api'

export const useMatchStore = defineStore('match', () => {
  const matches = ref([])
  const loading = ref(false)

  async function fetchMatches() {
    loading.value = true
    try {
      const result = await api.get('/matches')
      matches.value = result || []
      return result
    } finally {
      loading.value = false
    }
  }

  async function getMatchById(id) {
    return await api.get(`/matches/${id}`)
  }

  async function confirmMatch(id, accept) {
    return await api.post(`/matches/${id}/confirm`, { accept })
  }

  async function getContactInfo(id) {
    return await api.get(`/matches/${id}/contact`)
  }

  return {
    matches,
    loading,
    fetchMatches,
    getMatchById,
    confirmMatch,
    getContactInfo
  }
})
