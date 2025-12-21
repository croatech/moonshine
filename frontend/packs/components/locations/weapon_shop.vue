<template>
  <div class="shop">
    <div class="item-categories">
      <div class="col-md-3">
        <div class="category" v-for="(category, i) in categories">
          <a @click="showCategory(i)" :class="'category-btn btn btn-default ' + resolveCategoryClass(i)">{{category.name}}</a>
        </div>
      </div>

      <div class="col-md-9">
        <div class="items" v-if="currentCategory != null">
          <div class="item row" v-for="item in items">
            <equipment-item :item="item" :playerLevel="player.level"></equipment-item>

            <a @click="buyItem(item.id)" class="buy-button btn btn-success">
              Buy for {{ item.price }} gold
            </a>

            <b-alert variant="danger"
                     dismissible
                     :show="errorMessage != '' && boughtItemId == item.id">
              {{ errorMessage }}
            </b-alert>

            <b-alert variant="success"
                     dismissible
                     :show="successMessage != '' && boughtItemId == item.id">
              {{ successMessage }}
            </b-alert>
          </div>
        </div>

        <div v-if="currentCategory == null">
          <img :src="'../assets/locations/cities/moon_light/' + resource_name + '/bg.jpg'" alt="equipment" class="center background">
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from 'axios'
  import config from '../../config.js'
  import EquipmentItem from '../equipment/item.vue'

  export default {
    data: function () {
      return {
        categories: [],
        items: [],
        stats: config.stats,
        currentCategory: null,
        boughtItemId: null,
        resource_name: null,

        // flashes
        errorMessage: '',
        successMessage: ''
      }
    },
    components: {
      equipmentItem: EquipmentItem
    },
    computed: {
      player() {
        return this.$store.state.player
      }
    },
    methods: {
      getCategoriesList: function() {
        var pathname = window.location.pathname

        // Name of resource like locations/weapon_shop, background image and api path resolving depends on this variable
        this.resource_name = pathname.match(/weapon_shop|shop_of_artifacts/gi)[0]

        // Resolving of path to endpoint api, for example if location is weapon_shop, it must be called
        // api/equipment/categories
        var category_type = this.resource_name == 'weapon_shop' ? 'equipment' : 'artifacts'

        var link = '/stuff/categories?category_type=' + category_type
        axios.get(link)
        .then(response => {
          console.log(response)
          this.categories = response.data
        })
        .catch(e => {
          console.log(e)
        })
      },
      showCategory: function(index) {
        this.currentCategory = index
        this.items = this.categories[index].items
      },
      buyItem: function(item_id) {
        this.boughtItemId = item_id
        this.errorMessage = ''
        this.successMessage = ''

        var link = '/stuff/items/' + item_id + '/buy'
        axios.post(link, {
          item_type: 'equipment'
        })
        .then(response => {
          this.successMessage = 'Congrats! You have bought ' + response.data.name
          this.$store.commit('decrement_gold', response.data.price)
        })
        .catch(e => {
          this.errorMessage = e.response.data
        })
      },
      resolveCategoryClass: function(index) {
        if(this.currentCategory == index) {
          return 'active'
        }
      }
    },
    created: function() {
      this.getCategoriesList()
    }
  }
</script>
