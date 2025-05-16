import {ChatLineSquare, DataLine, Document, Opportunity, Setting, UserFilled} from '@element-plus/icons-vue'
import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'users',
      component: () => import('@/views/UserManagement.vue'),
      meta: {
        path: '/',
        icon: UserFilled,
        name: '用户管理',
      },
    },
    {
      path: '/patents/:id?',
      name: 'patents',
      component: () => import('@/views/GreenPatent.vue'),
      meta: {
        path: '/patents',
        icon: Document,
        name: '绿色专利',
      },
    },
    {
      path: '/reports/:id?',
      name: 'reports',
      component: () => import('@/views/EnterpriseReport.vue'),
      meta: {
        path: '/reports',
        icon: DataLine,
        name: '企业年报',
      },
    },
    {
      path: '/suggestions/:id?',
      name: 'suggestions',
      component: () => import('@/views/InnovativeSuggestion.vue'),
      meta: {
        path: '/suggestions',
        icon: Opportunity,
        name: '创新建议管理',
      },
    },
    {
      path: '/keywords',
      name: 'keywords',
      component: () => import('@/views/KeyWords.vue'),
      meta: {
        path: '/keywords',
        icon: ChatLineSquare,
        name: '关键词管理',
      },
    },
    {
      path: '/setting',
      name: 'setting',
      component: () => import('@/views/AppSetting.vue'),
      meta: {
        path: '/setting',
        icon: Setting,
        name: '系统设置',
      },
    },
  ],
})

export default router
