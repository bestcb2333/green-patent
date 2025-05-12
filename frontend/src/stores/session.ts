import { defineStore } from 'pinia'
import type { User } from '@/tables'
import { ref } from 'vue'
import { request } from '@/axios'

export const useSessionStore = defineStore('session', () => {

  const user = ref<User | null>(null)
  async function loadUser() {
    user.value = await request.get<any, User>('/myinfo')
  }

  return { user, loadUser }
})

export default useSessionStore
