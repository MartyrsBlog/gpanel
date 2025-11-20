import request from './index'

// 获取系统监控信息
export const getSystemMonitor = () => {
  return request({
    url: '/system/monitor',
    method: 'get'
  })
}

// 获取进程列表
export const getProcessList = () => {
  return request({
    url: '/system/processes',
    method: 'get'
  })
}

// 获取磁盘信息
export const getDiskInfo = () => {
  return request({
    url: '/system/disk',
    method: 'get'
  })
}