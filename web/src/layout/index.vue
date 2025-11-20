<template>
  <div class="layout-container">
    <!-- 侧边栏 -->
    <div class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <div class="logo">
          <span v-if="!sidebarCollapsed" class="logo-text">GPanel</span>
          <el-icon v-else class="logo-icon"><Monitor /></el-icon>
        </div>
        <el-button
          v-if="!sidebarCollapsed"
          text
          @click="toggleSidebar"
        >
          <el-icon><Fold /></el-icon>
        </el-button>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        :collapse="sidebarCollapsed"
        :unique-opened="true"
        router
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409eff"
      >
        <template v-for="route in menuRoutes" :key="route.path">
          <!-- 无子菜单 -->
          <el-menu-item
            v-if="!route.children || route.children.length === 0"
            :index="route.path"
          >
            <el-icon v-if="route.meta?.icon">
              <component :is="route.meta.icon" />
            </el-icon>
            <template #title>{{ route.meta?.title }}</template>
          </el-menu-item>
          
          <!-- 有子菜单 -->
          <el-sub-menu
            v-else
            :index="route.path"
          >
            <template #title>
              <el-icon v-if="route.meta?.icon">
                <component :is="route.meta.icon" />
              </el-icon>
              <span>{{ route.meta?.title }}</span>
            </template>
            
            <el-menu-item
              v-for="child in route.children"
              :key="child.path"
              :index="child.path"
            >
              <el-icon v-if="child.meta?.icon">
                <component :is="child.meta.icon" />
              </el-icon>
              <template #title>{{ child.meta?.title }}</template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </el-menu>
    </div>
    
    <!-- 主内容区 -->
    <div class="main-container">
      <!-- 顶部导航 -->
      <div class="navbar">
        <div class="navbar-left">
          <el-button
            v-if="sidebarCollapsed"
            text
            @click="toggleSidebar"
          >
            <el-icon><Expand /></el-icon>
          </el-button>
          
          <el-breadcrumb separator="/">
            <el-breadcrumb-item
              v-for="item in breadcrumbs"
              :key="item.path"
              :to="item.path"
            >
              {{ item.title }}
            </el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        
        <div class="navbar-right">
          <!-- 全屏切换 -->
          <el-tooltip content="全屏" placement="bottom">
            <el-button text @click="toggleFullscreen">
              <el-icon>
                <FullScreen v-if="!isFullscreen" />
                <Aim v-else />
              </el-icon>
            </el-button>
          </el-tooltip>
          
          <!-- 用户下拉菜单 -->
          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" icon="UserFilled" />
              <span class="username">{{ authStore.user.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人资料
                </el-dropdown-item>
                <el-dropdown-item command="password">
                  <el-icon><Lock /></el-icon>
                  修改密码
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </div>
      
      <!-- 页面内容 -->
      <div class="page-content">
        <router-view v-slot="{ Component }">
          <transition name="page" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </div>
    
    <!-- 修改密码对话框 -->
    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="500px"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="100px"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="新密码" prop="newPassword">
          <el-input
            v-model="passwordForm.newPassword"
            type="password"
            show-password
          />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input
            v-model="passwordForm.confirmPassword"
            type="password"
            show-password
          />
        </el-form-item>
      </el-form>
      
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="updatePassword">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, ref, reactive, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import { ElMessage, ElMessageBox } from 'element-plus'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)
const isFullscreen = ref(false)
const passwordDialogVisible = ref(false)
const passwordFormRef = ref()

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordRules = {
  oldPassword: [
    { required: true, message: '请输入旧密码', trigger: 'blur' }
  ],
  newPassword: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.newPassword) {
          callback(new Error('两次输入的密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

// 菜单路由
const menuRoutes = computed(() => {
  return router.options.routes
    .find(route => route.name === 'Layout')
    ?.children.filter(route => route.meta?.title && route.meta?.icon) || []
})

// 当前激活的菜单
const activeMenu = computed(() => {
  const { path } = route
  if (path === '/') return '/dashboard'
  return path
})

// 面包屑导航
const breadcrumbs = computed(() => {
  const matched = route.matched.filter(item => item.meta?.title)
  return matched.map(item => ({
    title: item.meta.title,
    path: item.path
  }))
})

// 切换侧边栏
const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

// 切换全屏
const toggleFullscreen = () => {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}

// 处理下拉菜单命令
const handleCommand = (command) => {
  switch (command) {
    case 'profile':
      // 打开个人资料页面
      break
    case 'password':
      passwordDialogVisible.value = true
      break
    case 'logout':
      handleLogout()
      break
  }
}

// 退出登录
const handleLogout = async () => {
  try {
    await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await authStore.logout()
    ElMessage.success('退出登录成功')
    router.push('/login')
  } catch (error) {
    // 用户取消操作
  }
}

// 修改密码
const updatePassword = async () => {
  if (!passwordFormRef.value) return
  
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // TODO: 调用修改密码API
        ElMessage.success('密码修改成功')
        passwordDialogVisible.value = false
        
        // 重置表单
        Object.assign(passwordForm, {
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        })
      } catch (error) {
        ElMessage.error('密码修改失败')
      }
    }
  })
}

// 监听全屏变化
watch(
  () => document.fullscreenElement,
  (val) => {
    isFullscreen.value = !!val
  }
)
</script>

<style lang="scss" scoped>
.layout-container {
  display: flex;
  height: 100vh;
}

.sidebar {
  width: 250px;
  background-color: #304156;
  transition: width 0.3s ease;
  
  &.collapsed {
    width: 64px;
  }
  
  .sidebar-header {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 20px;
    border-bottom: 1px solid #434a5a;
    
    .logo {
      display: flex;
      align-items: center;
      
      .logo-text {
        font-size: 20px;
        font-weight: bold;
        color: #fff;
      }
      
      .logo-icon {
        font-size: 24px;
        color: #409eff;
      }
    }
  }
  
  :deep(.el-menu) {
    border-right: none;
    
    .el-menu-item,
    .el-sub-menu__title {
      height: 50px;
      line-height: 50px;
      
      &:hover {
        background-color: #434a5a !important;
      }
    }
    
    .el-menu-item.is-active {
      background-color: #409eff !important;
    }
  }
}

.main-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  background-color: #f0f2f5;
}

.navbar {
  height: 60px;
  background: #fff;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  
  .navbar-left {
    display: flex;
    align-items: center;
    gap: 20px;
  }
  
  .navbar-right {
    display: flex;
    align-items: center;
    gap: 15px;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    padding: 5px 10px;
    border-radius: 6px;
    transition: background-color 0.3s;
    
    &:hover {
      background-color: #f5f7fa;
    }
    
    .username {
      font-size: 14px;
      color: #606266;
    }
  }
}

.page-content {
  flex: 1;
  padding: 20px;
  overflow-y: auto;
}

// 页面切换动画
.page-enter-active,
.page-leave-active {
  transition: opacity 0.3s, transform 0.3s;
}

.page-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    left: 0;
    top: 0;
    height: 100vh;
    z-index: 1000;
    transform: translateX(-100%);
    
    &.collapsed {
      transform: translateX(0);
      width: 250px;
    }
  }
  
  .main-container {
    margin-left: 0;
  }
}
</style>