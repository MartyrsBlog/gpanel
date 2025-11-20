import request from './index'

export default {
  // 获取系统监控信息
  getMonitor() {
    return request({
      url: '/system/monitor',
      method: 'get'
    })
  },

  // 获取进程列表
  getProcesses() {
    return request({
      url: '/system/processes',
      method: 'get'
    })
  },

  // 结束进程
  killProcess(pid) {
    return request({
      url: '/system/kill',
      method: 'post',
      data: { pid }
    })
  },

  // 获取磁盘信息
  getDiskInfo() {
    return request({
      url: '/system/disk',
      method: 'get'
    })
  }
}