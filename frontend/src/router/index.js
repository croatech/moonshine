import Vue from 'vue'
import Router from 'vue-router'
import HelloWorld from '@/components/HelloWorld'
import SignUp from '@/components/auth/SignUp'
import axios from 'axios'
import config from '../config'

Vue.use(Router)

axios.defaults.baseURL = config.apiUrl;

export default new Router({
  routes: [
    { path: '/', component: HelloWorld },
    { path: '/sign_up', component: SignUp }
  ]
})
