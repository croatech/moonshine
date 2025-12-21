<template>
 <div class="stats">
   <div v-for="stat in stats" :class="stat + ' stat'">
    <img :src="'../assets/' + stat + '.png'" :alt="stat">
      {{ player[stat] }}
      <a @click="incrementStat(stat)">+</a>
   </div>

   <div class="free-stats stat">
     Free stats:
     {{ player.free_stats }}
   </div>

   <b-alert variant="danger"
            dismissible
            :show="showErrorFlash"
            @dismissed="showErrorFlash=false">
     {{ errorMessage }}
   </b-alert>
 </div>
</template>

<script>
  import config from '../../config'
  import axios from 'axios'

  export default {
    data: function() {
      return {
        stats: config.stats,
        showErrorFlash: false,
        errorMessage: null
      }
    },
    methods: {
      incrementStat: function (stat_name) {
        var link = '/player/stats/' + stat_name + '/increase'
        axios.patch(link)
        .then(response => {
          this.player[stat_name] = response.data[stat_name]
          this.player.free_stats -= 1
        })
        .catch(e => {
          this.showErrorFlash = true
          this.errorMessage = e.response.data
        })
      }
    },
    computed: {
      player() {
        return this.$store.state.player
      }
    }
  }
</script>
