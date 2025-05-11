// @ts-check
import axios from 'axios'
import { ElMessage } from 'element-plus'

export const request = axios.create({
  baseURL: 'http://axogc.net:8702',
  timeout: 5000,
})

request.interceptors.request.use((config) => {
  config.baseURL = 'http://axogc.net:8702'
  return config
})

request.interceptors.response.use(
  (res) => {
    ElMessage({type: 'success', message: 'Error'})
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
