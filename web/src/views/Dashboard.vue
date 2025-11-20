<template>
  <div class="dashboard">
    <!-- 统计卡片 -->
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :md="6" v-for="stat in statsData" :key="stat.title">
        <el-card class="stats-card" :body-style="{ padding: '20px' }">
          <div class="stats-content">
            <div class="stats-info">
              <div class="stats-value">{{ stat.value }}</div>
              <div class="stats-label">{{ stat.title }}</div>
            </div>
            <div class="stats-icon" :style="{ color: stat.color }">
              <el-icon :size="40">
                <component :is="stat.icon" />
              </el-icon>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 图表区域 -->
    <el-row :gutter="20" class="charts-row">
      <!-- CPU 使用率 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>CPU 使用率</span>
              <el-tag type="success">{{ systemInfo.cpu?.usage?.toFixed(1) }}%</el-tag>
            </div>
          </template>
          <v-chart class="chart" :option="cpuChartOption" />
        </el-card>
      </el-col>
      
      <!-- 内存使用率 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>内存使用率</span>
              <el-tag type="warning">{{ systemInfo.memory?.usedPercent?.toFixed(1) }}%</el-tag>
            </div>
          </template>
          <v-chart class="chart" :option="memoryChartOption" />
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 系统信息和磁盘使用 -->
    <el-row :gutter="20" class="info-row">
      <!-- 系统信息 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <span>系统信息</span>
          </template>
          <div class="system-info">
            <div class="info-item">
              <span class="info-label">操作系统:</span>
              <span class="info-value">{{ systemInfo.host?.os }} {{ systemInfo.host?.platformVersion }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">主机名:</span>
              <span class="info-value">{{ systemInfo.host?.hostname }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">平台:</span>
              <span class="info-value">{{ systemInfo.host?.platform }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">运行时间:</span>
              <span class="info-value">{{ formatUptime(systemInfo.host?.uptime) }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">Go版本:</span>
              <span class="info-value">{{ systemInfo.runtime?.goVersion }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">CPU核心:</span>
              <span class="info-value">{{ systemInfo.cpu?.cores }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      
      <!-- 磁盘使用 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <span>磁盘使用</span>
          </template>
          <v-chart class="chart" :option="diskChartOption" />
        </el-card>
      </el-col>
    </el-row>
    
    <!-- 网络流量和系统负载 -->
    <el-row :gutter="20" class="info-row">
      <!-- 网络流量 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <span>网络流量</span>
          </template>
          <v-chart class="chart" :option="networkChartOption" />
        </el-card>
      </el-col>
      
      <!-- 系统负载 -->
      <el-col :xs="24" :lg="12">
        <el-card>
          <template #header>
            <span>系统负载</span>
          </template>
          <div class="load-info">
            <div class="load-item">
              <span class="load-label">1分钟:</span>
              <el-progress 
                :percentage="Math.min((systemInfo.load?.load1 || 0) * 100, 100)"
                :color="getLoadColor(systemInfo.load?.load1 || 0)"
              />
              <span class="load-value">{{ systemInfo.load?.load1?.toFixed(2) }}</span>
            </div>
            <div class="load-item">
              <span class="load-label">5分钟:</span>
              <el-progress 
                :percentage="Math.min((systemInfo.load?.load5 || 0) * 100, 100)"
                :color="getLoadColor(systemInfo.load?.load5 || 0)"
              />
              <span class="load-value">{{ systemInfo.load?.load5?.toFixed(2) }}</span>
            </div>
            <div class="load-item">
              <span class="load-label">15分钟:</span>
              <el-progress 
                :percentage="Math.min((systemInfo.load?.load15 || 0) * 100, 100)"
                :color="getLoadColor(systemInfo.load?.load15 || 0)"
              />
              <span class="load-value">{{ systemInfo.load?.load15?.toFixed(2) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart, BarChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import systemApi from '@/api/system'

// 注册 ECharts 组件
use([
  CanvasRenderer,
  LineChart,
  PieChart,
  BarChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const systemInfo = ref({})
let timer = null

// 统计数据
const statsData = computed(() => [
  {
    title: '网站数量',
    value: '12',
    icon: 'Globe',
    color: '#409eff'
  },
  {
    title: '数据库',
    value: '5',
    icon: 'Coin',
    color: '#67c23a'
  },
  {
    title: 'Docker容器',
    value: '8',
    icon: 'Box',
    color: '#e6a23c'
  },
  {
    title: '运行时间',
    value: formatUptime(systemInfo.value.host?.uptime),
    icon: 'Timer',
    color: '#f56c6c'
  }
])

// CPU 图表配置
const cpuChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis',
    formatter: '{b}: {c}%'
  },
  xAxis: {
    type: 'category',
    data: systemInfo.value.cpu?.usagePerCore?.map((_, index) => `Core ${index + 1}`) || []
  },
  yAxis: {
    type: 'value',
    max: 100,
    axisLabel: {
      formatter: '{value}%'
    }
  },
  series: [
    {
      data: systemInfo.value.cpu?.usagePerCore || [],
      type: 'bar',
      itemStyle: {
        color: '#409eff'
      }
    }
  ]
}))

// 内存图表配置
const memoryChartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{a} <br/>{b}: {c} GB ({d}%)'
  },
  legend: {
    orient: 'vertical',
    left: 'left'
  },
  series: [
    {
      name: '内存使用',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      label: {
        show: false,
        position: 'center'
      },
      emphasis: {
        label: {
          show: true,
          fontSize: '18',
          fontWeight: 'bold'
        }
      },
      labelLine: {
        show: false
      },
      data: [
        {
          value: formatBytes(systemInfo.value.memory?.used || 0),
          name: '已使用',
          itemStyle: { color: '#e6a23c' }
        },
        {
          value: formatBytes(systemInfo.value.memory?.available || 0),
          name: '可用',
          itemStyle: { color: '#67c23a' }
        }
      ]
    }
  ]
}))

// 磁盘图表配置
const diskChartOption = computed(() => ({
  tooltip: {
    trigger: 'item',
    formatter: '{a} <br/>{b}: {c} GB ({d}%)'
  },
  series: [
    {
      name: '磁盘使用',
      type: 'pie',
      radius: '50%',
      data: systemInfo.value.disk?.map(disk => ({
        value: formatBytes(disk.used),
        name: disk.mountpoint,
        itemStyle: {
          color: disk.usedPercent > 80 ? '#f56c6c' : disk.usedPercent > 60 ? '#e6a23c' : '#67c23a'
        }
      })) || []
    }
  ]
}))

// 网络图表配置
const networkChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis'
  },
  xAxis: {
    type: 'category',
    data: systemInfo.value.network?.ioStats?.map(net => net.name) || []
  },
  yAxis: {
    type: 'value',
    axisLabel: {
      formatter: function(value) {
        return formatBytes(value) + '/s'
      }
    }
  },
  series: [
    {
      name: '发送',
      type: 'bar',
      data: systemInfo.value.network?.ioStats?.map(net => net.bytesSent) || [],
      itemStyle: { color: '#409eff' }
    },
    {
      name: '接收',
      type: 'bar',
      data: systemInfo.value.network?.ioStats?.map(net => net.bytesRecv) || [],
      itemStyle: { color: '#67c23a' }
    }
  ]
}))

// 获取系统监控数据
const fetchSystemData = async () => {
  try {
    const response = await systemApi.getMonitor()
    systemInfo.value = response.data
  } catch (error) {
    console.error('获取系统数据失败:', error)
  }
}

// 格式化字节数
const formatBytes = (bytes) => {
  if (bytes === 0) return '0'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// 格式化运行时间
const formatUptime = (seconds) => {
  if (!seconds) return '0天'
  const days = Math.floor(seconds / 86400)
  const hours = Math.floor((seconds % 86400) / 3600)
  const minutes = Math.floor((seconds % 3600) / 60)
  return `${days}天 ${hours}小时 ${minutes}分钟`
}

// 获取负载颜色
const getLoadColor = (load) => {
  if (load < 1) return '#67c23a'
  if (load < 2) return '#e6a23c'
  return '#f56c6c'
}

onMounted(() => {
  fetchSystemData()
  // 每30秒更新一次数据
  timer = setInterval(fetchSystemData, 30000)
})

onUnmounted(() => {
  if (timer) {
    clearInterval(timer)
  }
})
</script>

<style lang="scss" scoped>
.dashboard {
  .stats-row {
    margin-bottom: 20px;
  }
  
  .charts-row {
    margin-bottom: 20px;
  }
  
  .info-row {
    margin-bottom: 20px;
  }
  
  .stats-card {
    .stats-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    
    .stats-info {
      .stats-value {
        font-size: 28px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 8px;
      }
      
      .stats-label {
        font-size: 14px;
        color: #909399;
      }
    }
    
    .stats-icon {
      opacity: 0.8;
    }
  }
  
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .chart {
    height: 300px;
    width: 100%;
  }
  
  .system-info {
    .info-item {
      display: flex;
      justify-content: space-between;
      padding: 12px 0;
      border-bottom: 1px solid #f0f0f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .info-label {
        color: #909399;
        font-weight: 500;
      }
      
      .info-value {
        color: #303133;
      }
    }
  }
  
  .load-info {
    .load-item {
      display: flex;
      align-items: center;
      margin-bottom: 20px;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .load-label {
        width: 60px;
        color: #606266;
        font-weight: 500;
      }
      
      .load-value {
        margin-left: 15px;
        color: #303133;
        font-weight: bold;
        min-width: 50px;
      }
      
      :deep(.el-progress) {
        flex: 1;
        margin: 0 15px;
      }
    }
  }
}

@media (max-width: 768px) {
  .chart {
    height: 250px;
  }
  
  .stats-card .stats-content .stats-info .stats-value {
    font-size: 24px;
  }
}
</style>