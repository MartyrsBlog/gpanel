import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/store/auth'
import NProgress from 'nprogress'
import 'nprogress/nprogress.css'

// 配置进度条
NProgress.configure({ showSpinner: false })

const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { 
      title: '登录',
      requiresAuth: false 
    }
  },
  {
    path: '/',
    name: 'Layout',
    component: () => import('@/layout/index.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/Dashboard.vue'),
        meta: { 
          title: '仪表盘',
          icon: 'Monitor',
          requiresAuth: true 
        }
      },
      {
        path: '/system',
        name: 'System',
        meta: { 
          title: '系统监控',
          icon: 'Setting',
          requiresAuth: true 
        },
        children: [
          {
            path: '/system/monitor',
            name: 'SystemMonitor',
            component: () => import('@/views/system/Monitor.vue'),
            meta: { 
              title: '实时监控',
              requiresAuth: true 
            }
          },
          {
            path: '/system/processes',
            name: 'Processes',
            component: () => import('@/views/system/Processes.vue'),
            meta: { 
              title: '进程管理',
              requiresAuth: true 
            }
          }
        ]
      },
      {
        path: '/sites',
        name: 'Sites',
        meta: { 
          title: '网站管理',
          icon: 'Globe',
          requiresAuth: true 
        },
        children: [
          {
            path: '/sites/list',
            name: 'SiteList',
            component: () => import('@/views/sites/List.vue'),
            meta: { 
              title: '网站列表',
              requiresAuth: true 
            }
          },
          {
            path: '/sites/create',
            name: 'SiteCreate',
            component: () => import('@/views/sites/Create.vue'),
            meta: { 
              title: '创建网站',
              requiresAuth: true 
            }
          }
        ]
      },
      {
        path: '/files',
        name: 'Files',
        component: () => import('@/views/files/Manager.vue'),
        meta: { 
          title: '文件管理',
          icon: 'Folder',
          requiresAuth: true 
        }
      },
      {
        path: '/database',
        name: 'Database',
        meta: { 
          title: '数据库',
          icon: 'Coin',
          requiresAuth: true 
        },
        children: [
          {
            path: '/database/list',
            name: 'DatabaseList',
            component: () => import('@/views/database/List.vue'),
            meta: { 
              title: '数据库列表',
              requiresAuth: true 
            }
          },
          {
            path: '/database/create',
            name: 'DatabaseCreate',
            component: () => import('@/views/database/Create.vue'),
            meta: { 
              title: '创建数据库',
              requiresAuth: true 
            }
          }
        ]
      },
      {
        path: '/docker',
        name: 'Docker',
        meta: { 
          title: 'Docker',
          icon: 'Box',
          requiresAuth: true 
        },
        children: [
          {
            path: '/docker/containers',
            name: 'DockerContainers',
            component: () => import('@/views/docker/Containers.vue'),
            meta: { 
              title: '容器管理',
              requiresAuth: true 
            }
          },
          {
            path: '/docker/images',
            name: 'DockerImages',
            component: () => import('@/views/docker/Images.vue'),
            meta: { 
              title: '镜像管理',
              requiresAuth: true 
            }
          }
        ]
      },
      {
        path: '/plugins',
        name: 'Plugins',
        component: () => import('@/views/plugins/Index.vue'),
        meta: { 
          title: '插件管理',
          icon: 'Grid',
          requiresAuth: true 
        }
      },
      {
        path: '/settings',
        name: 'Settings',
        component: () => import('@/views/Settings.vue'),
        meta: { 
          title: '系统设置',
          icon: 'Tools',
          requiresAuth: true 
        }
      }
    ]
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/404.vue'),
    meta: { title: '页面不存在' }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

// 全局前置守卫
router.beforeEach(async (to, from, next) => {
  NProgress.start()
  
  const authStore = useAuthStore()
  
  // 设置页面标题
  document.title = to.meta.title ? `${to.meta.title} - GPanel` : 'GPanel'
  
  // 检查是否需要认证
  if (to.meta.requiresAuth !== false && !authStore.isAuthenticated) {
    // 尝试从本地存储恢复登录状态
    await authStore.checkAuth()
    
    if (!authStore.isAuthenticated) {
      next('/login')
      return
    }
  }
  
  // 如果已登录且访问登录页，重定向到首页
  if (to.path === '/login' && authStore.isAuthenticated) {
    next('/')
    return
  }
  
  next()
})

// 全局后置钩子
router.afterEach(() => {
  NProgress.done()
})

export default router