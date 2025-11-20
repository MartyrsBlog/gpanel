import request from './index'

export default {
  // 获取数据库列表
  getDatabaseList() {
    return request({
      url: '/database/list',
      method: 'get'
    })
  },

  // 创建数据库
  createDatabase(data) {
    return request({
      url: '/database/create',
      method: 'post',
      data
    })
  },

  // 删除数据库
  deleteDatabase(name) {
    return request({
      url: '/database/delete',
      method: 'post',
      data: { name }
    })
  },

  // 备份数据库
  backupDatabase(name) {
    return request({
      url: '/database/backup',
      method: 'get',
      params: { name }
    })
  }
}