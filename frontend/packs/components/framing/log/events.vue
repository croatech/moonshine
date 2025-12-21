<template>
  <div class="events">
    <div class="event" v-for="event in events">
      <div v-html="event.body"></div>
    </div>
  </div>
</template>

<script>
  export default {
    computed: {
      events() {
        return this.$store.state.player.events
      }
    }
  }
</script>