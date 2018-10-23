import Vue from 'vue'
import Router from 'vue-router'

import Login from '../components/Login.vue'
import Dashboard from '../components/Dashboard.vue'

Vue.use(Router)

export const router = new Router({
    mode: 'history',
    routes: [
        { path: '/', component: Dashboard },
        { path: '/login', component: Login },
        { path: '*', redirect: '/' },
    ]
});
  
router.beforeEach((to, from, next) => {
    const publicPages = ['/login', '/register'];
    const authRequired = !publicPages.includes(to.path);
    const loggedIn = localStorage.getItem('user');

    return authRequired && !loggedIn ? next('/login') : next();
});