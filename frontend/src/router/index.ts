import { createRouter, createWebHistory } from 'vue-router';
import type { RouteRecordRaw } from 'vue-router';

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/tasks',
  },
  {
    path: '/tasks',
    name: 'TaskList',
    component: () => import('@/views/task/TaskList.vue'),
    meta: {
      title: '任务列表',
    },
  },
  {
    path: '/tasks/:id',
    name: 'TaskDetail',
    component: () => import('@/views/task/TaskDetail.vue'),
    meta: {
      title: '任务详情',
    },
  },
  {
    path: '/rules',
    name: 'RuleManagement',
    component: () => import('@/views/task/RuleManagement.vue'),
    meta: {
      title: '规则管理',
    },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  document.title = `${to.meta.title} - AI翻译系统`;

  // 检查是否需要登录
  const token = localStorage.getItem('token');
  if (!token && to.path !== '/login') {
    next('/login');
  } else {
    next();
  }
});

export default router; 