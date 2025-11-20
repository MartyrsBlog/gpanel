<template>
  <div class="processes-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>进程管理</span>
          <el-button type="primary" @click="refreshProcesses">
            <el-icon><Refresh /></el-icon>
            刷新
          </el-button>
        </div>
      </template>
      
      <el-table :data="processes" v-loading="loading" style="width: 100%">
        <el-table-column prop="pid" label="PID" width="80" />
        <el-table-column prop="name" label="进程名" width="150" />
        <el-table-column prop="user" label="用户" width="100" />
        <el-table-column prop="cpu" label="CPU%" width="80" />
        <el-table-column prop="memory" label="内存%" width="80" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="getStatusType(scope.row.status)">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="command" label="命令" />
        <el-table-column label="操作" width="150">
          <template #default="scope">
            <el-button size="small" @click="killProcess(scope.row.pid)">
              结束进程
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import systemApi from '@/api/system'

const processes = ref([])
const loading = ref(false)

const refreshProcesses = async () => {
  loading.value = true
  try {
    const response = await systemApi.getProcesses()
    processes.value = response.data
  } catch (error) {
    ElMessage.error('获取进程列表失败')
  } finally {
    loading.value = false
  }
}

const killProcess = async (pid) => {
  try {
    await ElMessageBox.confirm(`确定要结束进程 ${pid} 吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await systemApi.killProcess(pid)
    ElMessage.success('进程已结束')
    refreshProcesses()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('结束进程失败')
    }
  }
}

const getStatusType = (status) => {
  const statusMap = {
    '运行中': 'success',
    '睡眠': 'info',
    '停止': 'danger',
    '僵尸': 'warning'
  }
  return statusMap[status] || 'info'
}

onMounted(() => {
  refreshProcesses()
})
</script>

<style scoped>
.processes-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>