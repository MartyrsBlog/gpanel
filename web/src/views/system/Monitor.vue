<template>
  <div class="monitor-container">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon cpu">
              <el-icon><Cpu /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ systemInfo.cpu }}%</div>
              <div class="stat-label">CPU使用率</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon memory">
              <el-icon><Monitor /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ systemInfo.memory }}%</div>
              <div class="stat-label">内存使用率</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon disk">
              <el-icon><Folder /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ systemInfo.disk }}%</div>
              <div class="stat-label">磁盘使用率</div>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <div class="stat-icon network">
              <el-icon><Connection /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ systemInfo.network }}</div>
              <div class="stat-label">网络流量</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>CPU使用率趋势</span>
          </template>
          <div ref="cpuChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
      
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>内存使用率趋势</span>
          </template>
          <div ref="memoryChart" style="height: 300px;"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { Cpu, Monitor, Folder, Connection } from '@element-plus/icons-vue'
import systemApi from '@/api/system'

const systemInfo = ref({
  cpu: 0,
  memory: 0,
  disk: 0,
  network: '0 KB/s'
})

const cpuChart = ref(null)
const memoryChart = ref(null)
let cpuChartInstance = null
let memoryChartInstance = null
let refreshTimer = null

const initCharts = () => {
  // CPU图表
  cpuChartInstance = echarts.init(cpuChart.value)
  cpuChartInstance.setOption({
    title: {
      text: 'CPU使用率 (%)'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: []
    },
    yAxis: {
      type: 'value',
      max: 100
    },
    series: [{
      data: [],
      type: 'line',
      smooth: true
    }]
  })
  
  // 内存图表
  memoryChartInstance = echarts.init(memoryChart.value)
  memoryChartInstance.setOption({
    title: {
      text: '内存使用率 (%)'
    },
    tooltip: {
      trigger: 'axis'
    },
    xAxis: {
      type: 'category',
      data: []
    },
    yAxis: {
      type: 'value',
      max: 100
    },
    series: [{
      data: [],
      type: 'line',
      smooth: true
    }]
  })
}

const updateCharts = (cpuData, memoryData, timeLabels) => {
  if (cpuChartInstance && memoryChartInstance) {
    cpuChartInstance.setOption({
      xAxis: {
        data: timeLabels
      },
      series: [{
        data: cpuData
      }]
    })
    
    memoryChartInstance.setOption({
      xAxis: {
        data: timeLabels
      },
      series: [{
        data: memoryData
      }]
    })
  }
}

const refreshMonitor = async () => {
  try {
    const response = await systemApi.getMonitor()
    systemInfo.value = response.data.system
    
    // 更新图表数据
    const now = new Date().toLocaleTimeString()
    // 这里应该从响应中获取历史数据，简化处理
    updateCharts(
      [systemInfo.value.cpu],
      [systemInfo.value.memory],
      [now]
    )
  } catch (error) {
    console.error('获取监控数据失败:', error)
  }
}

onMounted(() => {
  initCharts()
  refreshMonitor()
  
  // 每5秒刷新一次
  refreshTimer = setInterval(refreshMonitor, 5000)
})

onUnmounted(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
  }
  if (cpuChartInstance) {
    cpuChartInstance.dispose()
  }
  if (memoryChartInstance) {
    memoryChartInstance.dispose()
  }
})
</script>

<style scoped>
.monitor-container {
  padding: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 15px;
  font-size: 24px;
  color: white;
}

.stat-icon.cpu {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.memory {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.disk {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.network {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 5px;
}

.stat-label {
  color: #909399;
  font-size: 14px;
}
</style>