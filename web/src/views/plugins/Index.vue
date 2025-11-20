<template>
  <div class="plugins-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>插件管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="scanPlugins">
              <el-icon><Search /></el-icon>
              扫描插件
            </el-button>
            <el-button @click="refreshPlugins">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>
        </div>
      </template>
      
      <el-table :data="plugins" v-loading="loading" style="width: 100%">
        <el-table-column prop="name" label="插件名称" width="200" />
        <el-table-column prop="version" label="版本" width="100" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="author" label="作者" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'info'">
              {{ scope.row.enabled ? '已启用' : '已禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="scope">
            <el-button 
              size="small" 
              :type="scope.row.enabled ? 'warning' : 'success'"
              @click="togglePlugin(scope.row)"
            >
              {{ scope.row.enabled ? '禁用' : '启用' }}
            </el-button>
            <el-button size="small" @click="configurePlugin(scope.row)">
              配置
            </el-button>
            <el-button size="small" type="danger" @click="uninstallPlugin(scope.row)">
              卸载
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
    
    <!-- 配置插件对话框 -->
    <el-dialog v-model="configDialogVisible" :title="`配置插件: ${currentPlugin?.name}`" width="600px">
      <el-form :model="pluginConfig" label-width="120px">
        <el-form-item 
          v-for="(value, key) in pluginConfig" 
          :key="key" 
          :label="key"
        >
          <el-input v-model="pluginConfig[key]" />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="savePluginConfig">
          保存配置
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import pluginApi from '@/api/plugin'

const plugins = ref([])
const loading = ref(false)
const configDialogVisible = ref(false)
const currentPlugin = ref(null)
const pluginConfig = reactive({})

const refreshPlugins = async () => {
  loading.value = true
  try {
    const response = await pluginApi.getPluginList()
    plugins.value = response.data
  } catch (error) {
    ElMessage.error('获取插件列表失败')
  } finally {
    loading.value = false
  }
}

const scanPlugins = async () => {
  try {
    await pluginApi.scanPlugins()
    ElMessage.success('插件扫描完成')
    refreshPlugins()
  } catch (error) {
    ElMessage.error('插件扫描失败')
  }
}

const togglePlugin = async (plugin) => {
  const action = plugin.enabled ? 'disable' : 'enable'
  
  try {
    await ElMessageBox.confirm(
      `确定要${action === 'disable' ? '禁用' : '启用'}插件 ${plugin.name} 吗？`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await pluginApi.togglePlugin(plugin.name, action)
    ElMessage.success(`插件${action === 'disable' ? '禁用' : '启用'}成功`)
    refreshPlugins()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error(`${action === 'disable' ? '禁用' : '启用'}插件失败`)
    }
  }
}

const configurePlugin = async (plugin) => {
  try {
    const response = await pluginApi.getPluginConfig(plugin.name)
    currentPlugin.value = plugin
    
    // 清空并重新填充配置对象
    Object.keys(pluginConfig).forEach(key => delete pluginConfig[key])
    Object.assign(pluginConfig, response.data.config || {})
    
    configDialogVisible.value = true
  } catch (error) {
    ElMessage.error('获取插件配置失败')
  }
}

const savePluginConfig = async () => {
  try {
    await pluginApi.savePluginConfig(currentPlugin.value.name, pluginConfig)
    ElMessage.success('插件配置保存成功')
    configDialogVisible.value = false
  } catch (error) {
    ElMessage.error('保存插件配置失败')
  }
}

const uninstallPlugin = async (plugin) => {
  try {
    await ElMessageBox.confirm(
      `确定要卸载插件 ${plugin.name} 吗？此操作不可恢复！`,
      '确认操作',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )
    
    await pluginApi.uninstallPlugin(plugin.name)
    ElMessage.success('插件卸载成功')
    refreshPlugins()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('卸载插件失败')
    }
  }
}

onMounted(() => {
  refreshPlugins()
})
</script>

<style scoped>
.plugins-container {
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