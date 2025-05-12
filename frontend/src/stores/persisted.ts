import { defineStore } from 'pinia'
import {reactive, ref} from 'vue'

const usePersistedStore = defineStore('persisted', () => {

  const setting = reactive({
    apiAddr: 'http://axogc.net:8702',
    darkMode: false,
    themeColor: '#28abce',
    fontSize: 14,
    language: 'zh',
  })

  const token = ref<string|null>(null)

  return {setting, token}
}, {persist: true})

export default usePersistedStore
