<template>
  <div class="row top">
    <div class="col-md-3">
      <img v-bind:src="item.image.url"/>
    </div>

    <div class="col-md-9">
      <div :class="'level ' + checkForNotAvailableLevel(item.required_level, this.playerLevel)">[{{ item.required_level }}]</div>

      <h3>{{ item.name }}</h3>

      <div v-for="stat in stats" :class="stat + ' stat'">
        <img :src="'../../assets/' + stat + '.png'" :alt="stat">
        {{ item[stat] }}
      </div>
    </div>
  </div>
</template>

<script>
  import config from '../../config'

  export default {
    props: ['item', 'playerLevel'],
    data: function() {
      return {
        stats: config.stats
      }
    },
    methods: {
      checkForNotAvailableLevel(itemLevel, playerLevel) {
        return itemLevel > playerLevel ? 'red' : ''
      }
    }
  }
</script>
