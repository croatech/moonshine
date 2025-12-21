<template>
  <div class="show inventory">
    <div class="col-md-5 grid col-md-offset-1">
      <profile :player="player"></profile>
    </div>

    <div class="col-md-5 col-md-offset-1">
      <items :player="player"></items>
    </div>
  </div>
</template>

<script>
  import Profile from './profile.vue'
  import Items from './items.vue'

  export default {
    components: {
      profile: Profile,
      items: Items
    },
    computed: {
      player() {
        return this.$store.state.player
      }
    }
  }
</script>
