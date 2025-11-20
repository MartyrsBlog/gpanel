import request from './index'

// 用户登录
export const login = (data) => {
  return request({
    url: '/auth/login',
    method: 'post',
    data
  })
}

// 用户登出
export const logout = () => {
  return request({
    url: '/auth/logout',
    method: 'post'
  })
}

// 获取用户信息
export const getUserInfo = () => {
  return request({
    url: '/auth/info',
    method: 'get'
  })
}

// 创建用户
export const createUser = (data) => {
  return request({
    url: '/auth/create',
    method: 'post',
    data
  })
}

// 修改密码
export const updatePassword = (data) => {
  return request({
    url: '/auth/password',
    method: 'put',
    data
  })
}