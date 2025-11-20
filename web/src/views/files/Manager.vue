<template>
  <div class="files-manager-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>文件管理</span>
          <div class="header-actions">
            <el-button @click="uploadFile">
              <el-icon><Upload /></el-icon>
              上传文件
            </el-button>
            <el-button @click="createFolder">
              <el-icon><FolderAdd /></el-icon>
              新建文件夹
            </el-button>
          </div>
        </div>
      </template>
      
      <!-- 面包屑导航 -->
      <el-breadcrumb separator="/" style="margin-bottom: 20px;">
        <el-breadcrumb-item 
          v-for="(item, index) in breadcrumbs" 
          :key="index"
          @click="navigateTo(item.path)"
          style="cursor: pointer;"
        >
          {{ item.name }}
        </el-breadcrumb-item>
      </el-breadcrumb>
      
      <!-- 文件列表 -->
      <el-table :data="files" v-loading="loading" style="width: 100%">
        <el-table-column width="50">
          <template #default="scope">
            <el-icon v-if="scope.row.type === 'directory'" color="#409EFF"><Folder /></el-icon>
            <el-icon v-else color="#67C23A"><Document /></el-icon>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="名称">
          <template #default="scope">
            <span 
              @click="handleFileClick(scope.row)"
              style="cursor: pointer; color: #409EFF;"
            >
              {{ scope.row.name }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="size" label="大小" width="100" />
        <el-table-column prop="modified" label="修改时间" width="180" />
        <el-table-column prop="permissions" label="权限" width="120" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="downloadFile(scope.row)">
              下载
            </el-button>
            <el-button size="small" @click="renameFile(scope.row)">
              重命名
            </el-button>
            <el-button size="small" type="danger" @click="deleteFile(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 上传文件对话框 -->
    <el-dialog v-model="uploadDialogVisible" title="上传文件" width="500px">
      <el-upload
        drag
        :action="uploadUrl"
        :data="{ path: currentPath }"
        :on-success="handleUploadSuccess"
        :on-error="handleUploadError"
        multiple
      >
        <el-icon class="el-icon--upload"><upload-filled /></el-icon>
        <div class="el-upload__text">
          将文件拖到此处，或<em>点击上传</em>
        </div>
      </el-upload>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Upload, FolderAdd, Folder, Document, UploadFilled } from '@element-plus/icons-vue'
import fileApi from '@/api/file'

const files = ref([])
const loading = ref(false)
const currentPath = ref('/')
const uploadDialogVisible = ref(false)

const breadcrumbs = computed(() => {
  const parts = currentPath.value.split('/').filter(Boolean)
  const crumbs = [{ name: '根目录', path: '/' }]
  
  let path = ''
  parts.forEach(part => {
    path += '/' + part
    crumbs.push({ name: part, path })
  })
  
  return crumbs
})

const uploadUrl = computed(() => {
  return '/api/file/upload'
})

const loadFiles = async () => {
  loading.value = true
  try {
    const response = await fileApi.getFileList({ path: currentPath.value })
    files.value = response.data
  } catch (error) {
    ElMessage.error('获取文件列表失败')
  } finally {
    loading.value = false
  }
}

const handleFileClick = (file) => {
  if (file.type === 'directory') {
    currentPath.value = currentPath.value === '/' ? `/${file.name}` : `${currentPath.value}/${file.name}`
    loadFiles()
  } else {
    ElMessage.info(`打开文件: ${file.name}`)
  }
}

const navigateTo = (path) => {
  currentPath.value = path
  loadFiles()
}

const uploadFile = () => {
  uploadDialogVisible.value = true
}

const createFolder = async () => {
  try {
    const { value: folderName } = await ElMessageBox.prompt('请输入文件夹名称', '新建文件夹', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPattern: /^[^/\\:*?"<>|]+$/,
      inputErrorMessage: '文件夹名称不能包含特殊字符'
    })
    
    if (folderName) {
      await fileApi.createFolder({
        path: currentPath.value,
        name: folderName
      })
      ElMessage.success('文件夹创建成功')
      loadFiles()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('创建文件夹失败')
    }
  }
}

const downloadFile = (file) => {
  if (file.type === 'directory') {
    ElMessage.warning('无法下载文件夹')
    return
  }
  
  const downloadUrl = `/api/file/download?path=${currentPath.value}/${file.name}`
  window.open(downloadUrl)
}

const renameFile = async (file) => {
  try {
    const { value: newName } = await ElMessageBox.prompt('请输入新名称', '重命名', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValue: file.name,
      inputPattern: /^[^/\\:*?"<>|]+$/,
      inputErrorMessage: '文件名不能包含特殊字符'
    })
    
    if (newName && newName !== file.name) {
      await fileApi.renameFile({
        path: currentPath.value,
        oldName: file.name,
        newName: newName
      })
      ElMessage.success('重命名成功')
      loadFiles()
    }
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('重命名失败')
    }
  }
}

const deleteFile = async (file) => {
  try {
    await ElMessageBox.confirm(`确定要删除 ${file.type === 'directory' ? '文件夹' : '文件'} ${file.name} 吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await fileApi.deleteFile({
      path: currentPath.value,
      name: file.name
    })
    ElMessage.success('删除成功')
    loadFiles()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除失败')
    }
  }
}

const handleUploadSuccess = () => {
  ElMessage.success('文件上传成功')
  uploadDialogVisible.value = false
  loadFiles()
}

const handleUploadError = () => {
  ElMessage.error('文件上传失败')
}

onMounted(() => {
  loadFiles()
})
</script>

<style scoped>
.files-manager-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 10px;
}
</style>