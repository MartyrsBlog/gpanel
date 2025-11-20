<template>
  <div class="sites-list-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>网站列表</span>
          <el-button type="primary" @click="$router.push('/sites/create')">
            <el-icon><Plus /></el-icon>
            创建网站
          </el-button>
        </div>
      </template>
      
      <el-table :data="sites" v-loading="loading" style="width: 100%">
        <el-table-column prop="domain" label="域名" width="200" />
        <el-table-column prop="path" label="网站目录" />
        <el-table-column prop="php_version" label="PHP版本" width="100" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.status === '运行中' ? 'success' : 'danger'">
              {{ scope.row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ssl" label="SSL" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.ssl ? 'success' : 'info'">
              {{ scope.row.ssl ? '已启用' : '未启用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button size="small" @click="manageSite(scope.row)">
              管理
            </el-button>
            <el-button size="small" type="danger" @click="deleteSite(scope.row)">
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
import siteApi from '@/api/site'

const sites = ref([])
const loading = ref(false)

const loadSites = async () => {
  loading.value = true
  try {
    const response = await siteApi.getSiteList()
    sites.value = response.data
  } catch (error) {
    ElMessage.error('获取网站列表失败')
  } finally {
    loading.value = false
  }
}

const manageSite = (site) => {
  // 跳转到网站管理页面
  ElMessage.info(`管理网站: ${site.domain}`)
}

const deleteSite = async (site) => {
  try {
    await ElMessageBox.confirm(`确定要删除网站 ${site.domain} 吗？`, '确认操作', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await siteApi.deleteSite(site.id)
    ElMessage.success('网站已删除')
    loadSites()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('删除网站失败')
    }
  }
}

onMounted(() => {
  loadSites()
})
</script>

<style scoped>
.sites-list-container {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>