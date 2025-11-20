<template>
  <div class="docker-containers-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>容器管理</span>
          <el-button type="primary" @click="refreshContainers">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table :data="containers" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="容器名称" width="200" />
        <el-table-column prop="image" label="镜像" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ports" label="端口" width="150" />
        <el-table-column prop="created" label="创建时间" width="180" />
        <el-table-column label="操作" width="250">
          <template #default="scope">
            <el-button 
              size="small" 
              :type="scope.row.status === 'running' ? 'warning' : 'success'"
              @click="toggleContainer(scope.row)"
            >
              {{ scope.row.status === 'running' ? '停止' : '启动' }}
            </el-button>
            <el-button size="small" @click="viewLogs(scope.row)">
              日志
            </el-button>
            <el-button size="small" type="danger" @click="removeContainer(scope.row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 日志对话框 -->
    <el-dialog v-model="logsDialogVisible" title="容器日志" width="80%">
      <div class="logs-container">
        <pre>{{ logs }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import dockerApi from '@/api/docker'

const containers = ref([])
const loading = ref(false)
const logsDialogVisible = ref(false)
const logs = ref('')

const refreshContainers = async () => {
  loading.value = true
  try {
    const response = await dockerApi.getContainers()
    containers.value = response.data
  } catch (error) {
    ElMessage.error('获取容器列表失败')
  } finally {
    loading.value = false
  }
}

const getStatusType = (status) => {
  const statusMap = {
    'running': 'success',
    'stopped': 'danger',
    'paused': 'warning',
    'restarting': 'info'
  }
  return statusMap[status] || 'info'
}

const toggleContainer = async (container) => {
  const action = container.status === 'running' ? 'stop' : 'start'
  
  try {
    await ElMessageBox.confirm(
      `确定要${action === 'stop' ? '停止' : '启动'}容器 ${container.name} 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    if (action === 'stop') {
      await dockerApi.stopContainer(container.id)
    } else {
      await dockerApi.startContainer(container.id)
    }
    
    ElMessage.success(`容器${action === 'stop' ? '停止' : '启动'}成功`)
    refreshContainers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action === 'stop' ? '停止' : '启动'}容器失败`)
    }
  }
}

const viewLogs = async (container) => {
  try {
    const response = await dockerApi.getContainerLogs(container.id)
    logs.value = response.data
    logsDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取容器日志失败')
  }
}

const removeContainer = async (container) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除容器 ${container.name} 吗？此操作不可恢复！`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await dockerApi.removeContainer(container.id)
    ElMessage.success('容器删除成功')
    refreshContainers()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除容器失败')
    }
  }
}

onMounted(() => {
  refreshContainers()
})
</script>

<style scoped>
.docker-containers-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.logs-container {
  max-height: 400px;
  overflow-y: auto;
  background-color: #f5f5f5;
  padding: 10px;
  border-radius: 4px;
}

.logs-container pre {
  margin: 0;
  font-family: 'Courier New', monospace;
  font-size: 12px;
  white-space: pre-wrap;
  word-wrap: break-word;
}
</style>