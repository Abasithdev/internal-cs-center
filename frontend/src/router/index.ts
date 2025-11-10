import { createRouter, createWebHistory } from 'vue-router'
import DashboardPage from '@/pages/DashboardPage.vue'
import { useAuthStore } from '@/stores/auth'
import LoginPage from '@/pages/LoginPage.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/login',
    },
    {
      path:"/login",
      component: LoginPage
    },
    {
      path: '/dashboard',
      component: DashboardPage,
      meta: {requiresAuth: true},
    },
  ],
})

router.beforeEach((to, _, next) => {
  const auth = useAuthStore();
  if(to.meta.requiresAuth && !auth.token){
    next("/login");
  }else{
    next();
  }
})

export default router
