/* eslint no-console:0 */
// This file is automatically compiled by Webpack, along with any other files
// present in this directory. You're encouraged to place your actual application logic in
// a relevant structure within app/javascript and only use these pack files to reference
// that code so it'll be compiled.
//
// To reference this file, add <%= javascript_pack_tag 'application' %> to the appropriate
// layout file, like app/views/layouts/application.html.erb

import Vue from 'vue/dist/vue.esm'
import Vuex from 'vuex'
import VueResource from 'vue-resource'
import BootstrapVue from 'bootstrap-vue'
import axios from 'axios'
import config from './config'
import { ws_subscriptions } from './mixins/ws_subscriptions'
import WeaponShop from './components/locations/weapon_shop.vue'
import CraftShop from './components/locations/craft_shop.vue'
import WaywardPines from './components/locations/wayward_pines.vue'
import Gold from './components/player/gold.vue'
import Stats from './components/player/stats.vue'
import Inventory from './components/player/inventory.vue'
import Fight from './components/fight.vue'
import Framing from './components/framing.vue'

import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

Vue.use(VueResource);
Vue.use(BootstrapVue);
Vue.use(Vuex);

axios.defaults.baseURL = config.apiUrl;

export const store = new Vuex.Store({
  state: {
    player: {},
    messages: {},
    recipient: ''
  },
  mutations: {
    increment_gold(state, n) {
      state.player.gold += n
    },
    decrement_gold(state, n) {
      state.player.gold -= n
    },
    updatePlayer(state, data) {
      state.player = data
    },
    setRecipient(state, data) {
      state.recipient = data
    },
    addMessage(state, data) {
      state.messages.unshift(data)
    },
    addEvent(state, data) {
      state.player.events.unshift(data)
    }
  },
  actions: {
    loadPlayer() {
      axios.get('/players/current').then(response => {
        this.state.player = response.data
      }).catch(e => {
          console.log(e)
      })
    },
    loadMessages() {
      var link = '/messages'
      axios.get(link).then(response => {
        this.state.messages = response.data
      }).catch(e => {
          console.log(e.response.data)
      })
    }
  }
});

document.addEventListener('DOMContentLoaded', () => {
  const app = new Vue({
    el: '#app',
    store,
    mixins: [ws_subscriptions],
    components: {
      WeaponShop,
      CraftShop,
      Gold,
      Stats,
      Inventory,
      Fight,
      WaywardPines,
      Framing
    },
    created: function() {
      this.$store.dispatch('loadPlayer')
      this.$store.dispatch('loadMessages')
      this.subscribeToChatStream()
      this.subscribeToEventsStream()
    }
  })
})
