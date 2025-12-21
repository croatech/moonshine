<template>
  <div class="shop">
    <div class="item-categories">
      <div class="col-md-3">
        <div class="category" v-for="(category, i) in categories">
          <a @click="showCategory(i)" :class="'category-btn btn btn-default ' + resolveCategoryClass(i)">{{category.name}}</a>
        </div>
      </div>

      <div class="col-md-9">
        <div class="items" v-if="currentCategoryIndex != null">
          <h3 v-if="items.length == 0">Items will be soon</h3>
          <div class="item row" v-for="item in items">
            <tool-item :item="item" :playerSkill="playerSkill"></tool-item>

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

        <div v-if="currentCategoryIndex == null">
          <img :src="'../assets/locations/cities/moon_light/craft_shop/bg.jpg'" alt="equipment" class="center background">
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import config from '../../config.js'
  import axios from 'axios'
  import ToolItem from '../tool/item.vue'

  export default {
    data: function () {
      return {
        categories: [],
        items: [],
        currentCategoryIndex: null,
        playerSkill: null,
        boughtItemId: null,
        resource_name: null,
        dismissSecs: 1,

        // flashes
        errorMessage: '',
        successMessage: ''
      }
    },
    components: {
      toolItem: ToolItem
    },
    computed: {
      player() {
        return this.$store.state.player
      }
    },
    methods: {
      getCategoriesList: function() {
        axios.get('/stuff/categories?category_type=tool')
        .then(response => {
          this.categories = response.data
        })
        .catch(e => {
          console.log(e)
        })
      },
      showCategory: function(index) {
        this.currentCategoryIndex = index

        var categoryName = this.categories[index].name.toLowerCase()
        this.playerSkill = this.player[categoryName + '_skill']

        this.items = this.categories[index].items
      },
      buyItem: function(item_id) {
        this.boughtItemId = item_id
        this.errorMessage = ''
        this.successMessage = ''

        var link = '/stuff/items/' + item_id + '/buy'
        axios.post(link, {
          item_type: 'tool'
        })
        .then(response => {
          this.successMessage = 'Congrats! You have bought ' + response.data.name
          this.$store.commit('gold-decrease', reponse.data.price)
        })
        .catch(e => {
          this.errorMessage = e.response.data
        })
      },
      resolveCategoryClass: function(index) {
        if(this.currentCategoryIndex == index) {
          return 'active'
        }
      }
    },
    created: function() {
      this.getCategoriesList()
    }
  }
</script>
