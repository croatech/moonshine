<template>
  <div class="row top">
    <div class="col-md-3">
      <img v-bind:src="item.image.url"/>
    </div>

    <div class="col-md-9">
      <div :class="'level ' + checkForNotAvailableSkill(item.required_skill, this.playerSkill)">[{{ item.required_skill }}]</div>

      <h3>{{ item.name }}</h3>
    </div>
  </div>
</template>

<script>
  import config from '../../config'

  export default {
    props: ['item', 'playerSkill'],
    data: function() {
      return {
        stats: config.stats
      }
    },
    methods: {
      checkForNotAvailableSkill(itemSkill, playerSkill) {
        return itemSkill > playerSkill ? 'red' : ''
      }
    }
  }
</script>
