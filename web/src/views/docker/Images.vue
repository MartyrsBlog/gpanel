<template>
  <div class="docker-images-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>镜像管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="pullImage">
              <el-icon><Download /></el-icon>
              拉取镜像
            </el-button>
            <el-button @click="refreshImages">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="images" v-loading="loading" style="width: 100%">
        <el-table-column prop="repository" label="镜像名称" width="300" />
        <el-table-column prop="tag" label="标签" width="100" />
        <el-table-column prop="size" label="大小" width="100" />
        <el-table-column prop="created" label="创建时间" width="180" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" type="danger" @click="removeImage(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 拉取镜像对话框 -->
    <el-dialog v-model="pullDialogVisible" title="拉取镜像" width="500px">
      <el-form :model="pullForm" label-width="100px">
        <el-form-item label="镜像名称" required>
          <el-input 
            v-model="pullForm.image" 
            placeholder="例如: nginx:latest 或 ubuntu:20.04"
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="pullDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="pullImageAction" :loading="pulling">
          拉取
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Refresh } from '@element-plus/icons-vue'
import dockerApi from '@/api/docker'

const images = ref([])
const loading = ref(false)
const pullDialogVisible = ref(false)
const pulling = ref(false)

const pullForm = reactive({
  image: ''
})

const refreshImages = async () => {
  loading.value = true
  try {
    const response = await dockerApi.getImages()
    images.value = response.data
  } catch (error) {
    ElMessage.error('获取镜像列表失败')
  } finally {
    loading.value = false
  }
}

const pullImage = () => {
  pullForm.image = ''
  pullDialogVisible.value = true
}

const pullImageAction = async () => {
  if (!pullForm.image.trim()) {
    ElMessage.warning('请输入镜像名称')
    return
  }
  
  pulling.value = true
  try {
    await dockerApi.pullImage({ image: pullForm.image })
    ElMessage.success('镜像拉取成功')
    pullDialogVisible.value = false
    refreshImages()
  } catch (error) {
    ElMessage.error('镜像拉取失败')
  } finally {
    pulling.value = false
  }
}

const removeImage = async (image) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除镜像 ${image.repository}:${image.tag} 吗？此操作不可恢复！`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await dockerApi.removeImage(image.id)
    ElMessage.success('镜像删除成功')
    refreshImages()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除镜像失败')
    }
  }
}

onMounted(() => {
  refreshImages()
})
</script>

<style scoped>
.docker-images-container {
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