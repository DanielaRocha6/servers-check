import Vue from 'vue'
import Router from 'vue-router'
import Home from './views/Home.vue'
import AllDomains from './views/AllDomains.vue'
import Domain from './views/Domain.vue'
Vue.use(Router)

export default new Router({
  mode: 'history',
  base: process.env.BASE_URL,
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home
    },
    {
      path: '/domains',
      name: 'domains',
      component: AllDomains
    },
    {
      path: '/domains/:url',
      component: Domain,
      props: true
    }
  ]
})
