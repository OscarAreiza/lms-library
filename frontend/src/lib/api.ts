import axios from 'axios'

import { clearToken, getToken } from './auth'

// Base URL points at library-api directly in dev, or through nginx in production
// (library-docs/05-architecture/decisions/records/ADR-003-nginx-reverse-proxy.md).
export const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api/v1',
})

api.interceptors.request.use((config) => {
  const token = getToken()
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

api.interceptors.response.use(
  (response) => response,
  (error) => {
    // A 401 always means "log in again" — there is no refresh-token flow in v1
    // (library-docs/07-api/authentication.md).
    if (error.response?.status === 401) {
      clearToken()
      if (window.location.pathname !== '/login') {
        window.location.assign('/login')
      }
    }
    return Promise.reject(error)
  },
)
