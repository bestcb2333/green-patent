// @ts-check
import axios from 'axios'
import { ElMessage } from 'element-plus'
import usePersistedStore from './stores/persisted'

export const request = axios.create({
  timeout: 5000,
})

request.interceptors.request.use((config) => {
  const persisted = usePersistedStore()
  config.baseURL = persisted.setting.apiAddr
  return config
})

request.interceptors.response.use(
  (res) => {
    return res.data.data
  },
  (err) => {
    if (err.response) {
      const message = 'Error'
      ElMessage({type: 'error', message: message})
      return Promise.reject(new Error(message))
    } else {
      return Promise.reject(err)
    }
  },
)
