<template>
  <div class="gold">
    <img src="../../assets/gold.png" alt="gold">
    {{ gold }}
  </div>
</template>

<script>
  export default {
    props: ['player'],
    computed: {
      gold() {
        return this.$store.state.player.gold
      }
    }
  }
</script>
