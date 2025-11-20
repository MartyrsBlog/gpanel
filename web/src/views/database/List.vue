<template>
  <div class="database-list-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>数据库列表</span>
          <el-button type="primary" @click="$router.push('/database/create')">
            <el-icon><Plus /></el-icon>
            创建数据库
          </el-button>
        </div>
      </template>
      
      <el-table :data="databases" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="数据库名" width="200" />
        <el-table-column prop="username" label="用户名" width="150" />
        <el-table-column prop="charset" label="字符集" width="100" />
        <el-table-column prop="collation" label="排序规则" width="150" />
        <el-table-column prop="size" label="大小" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250">
          <template #default="scope">
            <el-button size="small" @click="backupDatabase(scope.row)">
              备份
            </el-button>
            <el-button size="small" @click="manageDatabase(scope.row)">
              管理
            </el-button>
            <el-button size="small" type="danger" @click="deleteDatabase(scope.row)">
              删除
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
import { Plus } from '@element-plus/icons-vue'
import databaseApi from '@/api/database'

const databases = ref([])
const loading = ref(false)

const loadDatabases = async () => {
  loading.value = true
  try {
    const response = await databaseApi.getDatabaseList()
    databases.value = response.data
  } catch (error) {
    ElMessage.error('获取数据库列表失败')
  } finally {
    loading.value = false
  }
}

const backupDatabase = async (database) => {
  try {
    await databaseApi.backupDatabase(database.name)
    ElMessage.success('数据库备份已开始')
  } catch (error) {
    ElMessage.error('数据库备份失败')
  }
}

const manageDatabase = (database) => {
  ElMessage.info(`管理数据库: ${database.name}`)
}

const deleteDatabase = async (database) => {
  try {
    await ElMessageBox.confirm(`确定要删除数据库 ${database.name} 吗？此操作不可恢复！`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await databaseApi.deleteDatabase(database.name)
    ElMessage.success('数据库已删除')
    loadDatabases()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除数据库失败')
    }
  }
}

onMounted(() => {
  loadDatabases()
})
</script>

<style scoped>
.database-list-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>