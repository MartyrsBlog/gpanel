import request from './index'

export default {
  // 获取插件列表
  getPluginList() {
    return request({
      url: '/plugin/list',
      method: 'get'
    })
  },

  // 扫描插件
  scanPlugins() {
    return request({
      url: '/plugin/scan',
      method: 'post'
    })
  },

  // 启用/禁用插件
  togglePlugin(name, action) {
    return request({
      url: '/plugin/toggle',
      method: 'post',
      data: { name, action }
    })
  },

  // 获取插件配置
  getPluginConfig(name) {
    return request({
      url: '/plugin/config',
      method: 'get',
      params: { name }
    })
  },

  // 保存插件配置
  savePluginConfig(name, config) {
    return request({
      url: '/plugin/config',
      method: 'post',
      data: { name, config }
    })
  },

  // 卸载插件
  uninstallPlugin(name) {
    return request({
      url: '/plugin/uninstall',
      method: 'post',
      data: { name }
    })
  }
}