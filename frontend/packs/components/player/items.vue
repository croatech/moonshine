<template>
  <div>
    <nav class="navbar">
      <ul>
        <li><a @click="setCategory('all')" class="btn btn-default">All</a></li>
        <li><a @click="setCategory('equipment_items')" class="btn btn-default">Equipment</a></li>
        <li><a @click="setCategory('tool_items')" class="btn btn-default">Tools</a></li>
        <li><a @click="setCategory('resources')" class="btn btn-default">Resources</a></li>
      </ul>
    </nav>

    <div class="items">
      <div class="equipment_items" v-if="currentCategory == 'all' || currentCategory == 'equipment_items'">
        <div class="item row" v-for="item in player.equipment_items" v-if="player.equipment_items.length > 0">
          <equipment-item :item="item"></equipment-item>

          <a @click="putOnItem(item, 'equipment')" class="btn btn-info">Put on</a>
          <a @click="sellItem(item, 'equipment')" class="btn btn-info">Sell for {{ item.sell_price }} gold</a>
          <b-alert variant="danger"
                   dismissible
                   :show="errorMessage != '' && item.id == sellingItemId">
            {{ errorMessage }}
          </b-alert>
        </div>
        <h2 v-if="player.equipment_items.length == 0">No equipment items</h2>
      </div>

      <div class="tool_items" v-if="currentCategory == 'all' || currentCategory == 'tool_items'">
        <div class="item row" v-for="item in player.tool_items" v-if="player.tool_items.length > 0">
          <tool-item :item="item"></tool-item>

          <a @click="putOnItem(item, 'tool')" class="btn btn-info">Put on</a>
          <a @click="sellItem(item, 'tool')" class="btn btn-info">Sell for {{ item.sell_price }} gold</a>
          <b-alert variant="danger"
                   dismissible
                   :show="errorMessage != '' && item.id == sellingItemId">
            {{ errorMessage }}
          </b-alert>
        </div>
        <h2 v-if="player.tool_items.length == 0">No tool items</h2>
      </div>

      <div class="resources" v-if="currentCategory == 'all' || currentCategory == 'resources'">
        <div class="item row" v-for="item in player.resources" v-if="player.resources.length > 0">
          <resource :item="item"></resource>

          <a @click="sellItem(item, 'resource')" class="btn btn-info">Sell for {{ item.price }} gold</a>
          <b-alert variant="danger"
                   dismissible
                   :show="errorMessage != '' && item.id == sellingResourceId">
            {{ errorMessage }}
          </b-alert>
        </div>
        <h2 v-if="player.resources.length == 0">No resources</h2>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from 'axios'
  import EquipmentItem from '../equipment/item.vue'
  import ToolItem from '../tool/item.vue'
  import Resource from '../resources/item.vue'

  export default {
    props: ['player'],
    data: function() {
      return {
        currentCategory: 'all',
        errorMessage: '',
        sellingItemId: 0,
        sellingResourceId: 0
      }
    },
    components: {
      equipmentItem: EquipmentItem,
      toolItem: ToolItem,
      resource: Resource
    },
    methods: {
      setCategory: function(name) {
        this.currentCategory = name
      },
      putOnItem(item, type) {
        var link = '/stuff/items/' + item.id + '/put_on'
        axios.post(link, {
          item_type: type
        })
        .then(response => {
          this.$store.commit('updatePlayer', response.data)
        })
        .catch(e => {
          console.log(e.response)
        })
      },
      sellItem(item, type) {
        this.cleanErrors()

        var link = '/stuff/items/' + item.id + '/sell'
        axios.post(link, {
          item_type: type
        })
        .then(response => {
          this.$store.commit('updatePlayer', response.data)
        })
        .catch(e => {
          this.sellingItemId = item.id
          this.errorMessage = e.response.data
        })
      },
      cleanErrors() {
        this.errorMessage = ''
        this.sellingItemId = 0
        this.sellingResourceId = 0
      }
    }
  }
</script>
