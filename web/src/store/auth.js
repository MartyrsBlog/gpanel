import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login, logout, getUserInfo } from '@/api/auth'
import Cookies from 'js-cookie'

const TOKEN_KEY = 'gpanel_token'
const USER_KEY = 'gpanel_user'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(Cookies.get(TOKEN_KEY) || '')
  const user = ref(JSON.parse(localStorage.getItem(USER_KEY) || '{}'))
  const loading = ref(false)

  const isAuthenticated = computed(() => !!token.value)

  // 登录
  const doLogin = async (credentials) => {
    loading.value = true
    try {
      const response = await login(credentials)
      const { token: newToken, user: userInfo } = response.data
      
      // 保存token和用户信息
      token.value = newToken
      user.value = userInfo
      
      // 持久化存储
      Cookies.set(TOKEN_KEY, newToken, { expires: 7 }) // 7天过期
      localStorage.setItem(USER_KEY, JSON.stringify(userInfo))
      
      return response
    } catch (error) {
      throw error
    } finally {
      loading.value = false
    }
  }

  // 登出
  const doLogout = async () => {
    loading.value = true
    try {
      await logout()
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // 清除本地存储
      token.value = ''
      user.value = {}
      Cookies.remove(TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
      loading.value = false
    }
  }

  // 检查认证状态
  const checkAuth = async () => {
    if (!token.value) {
      return false
    }

    try {
      const response = await getUserInfo()
      user.value = response.data
      localStorage.setItem(USER_KEY, JSON.stringify(response.data))
      return true
    } catch (error) {
      // token失效，清除本地存储
      token.value = ''
      user.value = {}
      Cookies.remove(TOKEN_KEY)
      localStorage.removeItem(USER_KEY)
      return false
    }
  }

  // 更新用户信息
  const updateUser = (newUserInfo) => {
    user.value = { ...user.value, ...newUserInfo }
    localStorage.setItem(USER_KEY, JSON.stringify(user.value))
  }

  return {
    token,
    user,
    loading,
    isAuthenticated,
    login: doLogin,
    logout: doLogout,
    checkAuth,
    updateUser
  }
})