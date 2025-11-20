import request from './index'

export default {
  // 获取网站列表
  getSiteList() {
    return request({
      url: '/site/list',
      method: 'get'
    })
  },

  // 创建网站
  createSite(data) {
    return request({
      url: '/site/create',
      method: 'post',
      data
    })
  },

  // 删除网站
  deleteSite(id) {
    return request({
      url: '/site/delete',
      method: 'post',
      data: { id }
    })
  },

  // 申请SSL证书
  applySSL(data) {
    return request({
      url: '/site/ssl',
      method: 'post',
      data
    })
  }
}