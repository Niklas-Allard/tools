import { createRouter, createWebHashHistory } from 'vue-router'
import HomeView from '@/views/HomeView.vue'

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/ytdl',
      name: 'ytdl',
      component: () => import('@/views/YtdlView.vue'),
    },
    {
      path: '/qrcode',
      name: 'qrcode',
      component: () => import('@/views/QrcodeView.vue'),
    },
  ],
})

export default router