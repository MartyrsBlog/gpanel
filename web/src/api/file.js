import request from './index'

export default {
  // 获取文件列表
  getFileList(params) {
    return request({
      url: '/file/list',
      method: 'get',
      params
    })
  },

  // 上传文件
  uploadFile(data) {
    return request({
      url: '/file/upload',
      method: 'post',
      data,
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // 下载文件
  downloadFile(params) {
    return request({
      url: '/file/download',
      method: 'get',
      params,
      responseType: 'blob'
    })
  },

  // 创建文件夹
  createFolder(data) {
    return request({
      url: '/file/create',
      method: 'post',
      data
    })
  },

  // 重命名文件
  renameFile(data) {
    return request({
      url: '/file/rename',
      method: 'post',
      data
    })
  },

  // 删除文件
  deleteFile(data) {
    return request({
      url: '/file/delete',
      method: 'post',
      data
    })
  },

  // 解压文件
  extractFile(data) {
    return request({
      url: '/file/extract',
      method: 'post',
      data
    })
  },

  // 修改文件权限
  chmodFile(data) {
    return request({
      url: '/file/chmod',
      method: 'post',
      data
    })
  }
}