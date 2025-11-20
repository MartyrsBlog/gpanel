import request from './index'

export default {
  // 获取容器列表
  getContainers() {
    return request({
      url: '/docker/containers',
      method: 'get'
    })
  },

  // 启动容器
  startContainer(id) {
    return request({
      url: '/docker/container/start',
      method: 'post',
      data: { id }
    })
  },

  // 停止容器
  stopContainer(id) {
    return request({
      url: '/docker/container/stop',
      method: 'post',
      data: { id }
    })
  },

  // 删除容器
  removeContainer(id) {
    return request({
      url: '/docker/container/remove',
      method: 'post',
      data: { id }
    })
  },

  // 获取容器日志
  getContainerLogs(id) {
    return request({
      url: '/docker/container/logs',
      method: 'get',
      params: { id }
    })
  },

  // 获取镜像列表
  getImages() {
    return request({
      url: '/docker/images',
      method: 'get'
    })
  },

  // 拉取镜像
  pullImage(data) {
    return request({
      url: '/docker/image/pull',
      method: 'post',
      data
    })
  },

  // 删除镜像
  removeImage(id) {
    return request({
      url: '/docker/image/remove',
      method: 'post',
      data: { id }
    })
  }
}