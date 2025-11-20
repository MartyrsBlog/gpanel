<template>
  <div class="settings">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>系统设置</span>
        </div>
      </template>
      
      <el-tabs v-model="activeTab" type="border-card">
        <!-- 基本设置 -->
        <el-tab-pane label="基本设置" name="basic">
          <el-form
            ref="basicFormRef"
            :model="basicForm"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="系统名称">
              <el-input v-model="basicForm.systemName" placeholder="GPanel" />
            </el-form-item>
            
            <el-form-item label="系统描述">
              <el-input
                v-model="basicForm.systemDescription"
                type="textarea"
                :rows="3"
                placeholder="基于Go + Gin + Vue3的服务器管理面板"
              />
            </el-form-item>
            
            <el-form-item label="管理员邮箱">
              <el-input v-model="basicForm.adminEmail" placeholder="admin@example.com" />
            </el-form-item>
            
            <el-form-item label="时区">
              <el-select v-model="basicForm.timezone" placeholder="选择时区">
                <el-option
                  v-for="tz in timezones"
                  :key="tz.value"
                  :label="tz.label"
                  :value="tz.value"
                />
              </el-select>
            </el-form-item>
            
            <el-form-item label="语言">
              <el-select v-model="basicForm.language" placeholder="选择语言">
                <el-option label="简体中文" value="zh-CN" />
                <el-option label="English" value="en-US" />
              </el-select>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveBasicSettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 安全设置 -->
        <el-tab-pane label="安全设置" name="security">
          <el-form
            ref="securityFormRef"
            :model="securityForm"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="JWT密钥">
              <el-input
                v-model="securityForm.jwtSecret"
                type="password"
                show-password
                placeholder="JWT签名密钥"
              />
            </el-form-item>
            
            <el-form-item label="Token过期时间">
              <el-input-number
                v-model="securityForm.tokenExpire"
                :min="3600"
                :max="86400 * 30"
                placeholder="秒"
              />
              <span class="ml-2">秒</span>
            </el-form-item>
            
            <el-form-item label="密码强度">
              <el-radio-group v-model="securityForm.passwordStrength">
                <el-radio label="low">低 (6位以上)</el-radio>
                <el-radio label="medium">中 (包含数字和字母)</el-radio>
                <el-radio label="high">高 (包含数字、字母和特殊字符)</el-radio>
              </el-radio-group>
            </el-form-item>
            
            <el-form-item label="登录失败锁定">
              <el-switch
                v-model="securityForm.loginLockEnabled"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
            
            <el-form-item
              v-if="securityForm.loginLockEnabled"
              label="锁定阈值"
            >
              <el-input-number
                v-model="securityForm.loginLockThreshold"
                :min="3"
                :max="10"
              />
              <span class="ml-2">次失败后锁定</span>
            </el-form-item>
            
            <el-form-item
              v-if="securityForm.loginLockEnabled"
              label="锁定时间"
            >
              <el-input-number
                v-model="securityForm.loginLockTime"
                :min="300"
                :max="3600"
              />
              <span class="ml-2">秒</span>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveSecuritySettings">保存设置</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 通知设置 -->
        <el-tab-pane label="通知设置" name="notification">
          <el-form
            ref="notificationFormRef"
            :model="notificationForm"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="邮件通知">
              <el-switch
                v-model="notificationForm.emailEnabled"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
            
            <template v-if="notificationForm.emailEnabled">
              <el-form-item label="SMTP服务器">
                <el-input v-model="notificationForm.smtpHost" placeholder="smtp.example.com" />
              </el-form-item>
              
              <el-form-item label="SMTP端口">
                <el-input-number v-model="notificationForm.smtpPort" :min="1" :max="65535" />
              </el-form-item>
              
              <el-form-item label="发送邮箱">
                <el-input v-model="notificationForm.smtpUser" placeholder="noreply@example.com" />
              </el-form-item>
              
              <el-form-item label="邮箱密码">
                <el-input
                  v-model="notificationForm.smtpPassword"
                  type="password"
                  show-password
                />
              </el-form-item>
              
              <el-form-item label="启用SSL">
                <el-switch
                  v-model="notificationForm.smtpSSL"
                  active-text="是"
                  inactive-text="否"
                />
              </el-form-item>
            </template>
            
            <el-form-item label="通知事件">
              <el-checkbox-group v-model="notificationForm.events">
                <el-checkbox label="login">登录异常</el-checkbox>
                <el-checkbox label="system">系统警告</el-checkbox>
                <el-checkbox label="disk">磁盘空间不足</el-checkbox>
                <el-checkbox label="backup">备份完成</el-checkbox>
                <el-checkbox label="update">系统更新</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            
            <el-form-item>
              <el-button type="primary" @click="saveNotificationSettings">保存设置</el-button>
              <el-button @click="testEmail">测试邮件</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        
        <!-- 备份设置 -->
        <el-tab-pane label="备份设置" name="backup">
          <el-form
            ref="backupFormRef"
            :model="backupForm"
            label-width="120px"
            class="settings-form"
          >
            <el-form-item label="自动备份">
              <el-switch
                v-model="backupForm.autoEnabled"
                active-text="启用"
                inactive-text="禁用"
              />
            </el-form-item>
            
            <template v-if="backupForm.autoEnabled">
              <el-form-item label="备份频率">
                <el-select v-model="backupForm.frequency">
                  <el-option label="每天" value="daily" />
                  <el-option label="每周" value="weekly" />
                  <el-option label="每月" value="monthly" />
                </el-select>
              </el-form-item>
              
              <el-form-item label="备份时间">
                <el-time-picker
                  v-model="backupForm.time"
                  format="HH:mm"
                  placeholder="选择时间"
                />
              </el-form-item>
              
              <el-form-item label="保留天数">
                <el-input-number
                  v-model="backupForm.retentionDays"
                  :min="1"
                  :max="365"
                />
                <span class="ml-2">天</span>
              </el-form-item>
              
              <el-form-item label="备份内容">
                <el-checkbox-group v-model="backupForm.content">
                  <el-checkbox label="config">配置文件</el-checkbox>
                  <el-checkbox label="database">数据库</el-checkbox>
                  <el-checkbox label="websites">网站文件</el-checkbox>
                  <el-checkbox label="logs">日志文件</el-checkbox>
                </el-checkbox-group>
              </el-form-item>
              
              <el-form-item label="备份路径">
                <el-input v-model="backupForm.path" placeholder="/backup" />
              </el-form-item>
            </template>
            
            <el-form-item>
              <el-button type="primary" @click="saveBackupSettings">保存设置</el-button>
              <el-button @click="createBackup">立即备份</el-button>
            </el-form-item>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'

const activeTab = ref('basic')
const basicFormRef = ref()
const securityFormRef = ref()
const notificationFormRef = ref()
const backupFormRef = ref()

// 基本设置
const basicForm = reactive({
  systemName: 'GPanel',
  systemDescription: '基于Go + Gin + Vue3的服务器管理面板',
  adminEmail: 'admin@example.com',
  timezone: 'Asia/Shanghai',
  language: 'zh-CN'
})

// 安全设置
const securityForm = reactive({
  jwtSecret: '',
  tokenExpire: 86400,
  passwordStrength: 'medium',
  loginLockEnabled: true,
  loginLockThreshold: 5,
  loginLockTime: 1800
})

// 通知设置
const notificationForm = reactive({
  emailEnabled: false,
  smtpHost: '',
  smtpPort: 587,
  smtpUser: '',
  smtpPassword: '',
  smtpSSL: true,
  events: ['login', 'system']
})

// 备份设置
const backupForm = reactive({
  autoEnabled: false,
  frequency: 'daily',
  time: null,
  retentionDays: 7,
  content: ['config', 'database'],
  path: '/backup'
})

// 时区选项
const timezones = [
  { label: '北京时间 (GMT+8)', value: 'Asia/Shanghai' },
  { label: '东京时间 (GMT+9)', value: 'Asia/Tokyo' },
  { label: '纽约时间 (GMT-5)', value: 'America/New_York' },
  { label: '伦敦时间 (GMT+0)', value: 'Europe/London' },
  { label: '巴黎时间 (GMT+1)', value: 'Europe/Paris' }
]

// 保存基本设置
const saveBasicSettings = async () => {
  try {
    // TODO: 调用API保存设置
    ElMessage.success('基本设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 保存安全设置
const saveSecuritySettings = async () => {
  try {
    // TODO: 调用API保存设置
    ElMessage.success('安全设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 保存通知设置
const saveNotificationSettings = async () => {
  try {
    // TODO: 调用API保存设置
    ElMessage.success('通知设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 测试邮件
const testEmail = async () => {
  try {
    // TODO: 调用API测试邮件发送
    ElMessage.success('测试邮件发送成功')
  } catch (error) {
    ElMessage.error('邮件发送失败')
  }
}

// 保存备份设置
const saveBackupSettings = async () => {
  try {
    // TODO: 调用API保存设置
    ElMessage.success('备份设置保存成功')
  } catch (error) {
    ElMessage.error('保存失败')
  }
}

// 立即备份
const createBackup = async () => {
  try {
    // TODO: 调用API创建备份
    ElMessage.success('备份任务已启动')
  } catch (error) {
    ElMessage.error('备份失败')
  }
}
</script>

<style lang="scss" scoped>
.settings {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  
  .settings-form {
    max-width: 600px;
    
    .el-form-item {
      margin-bottom: 24px;
    }
  }
  
  :deep(.el-tabs__content) {
    padding: 20px;
  }
  
  .ml-2 {
    margin-left: 8px;
  }
}
</style>