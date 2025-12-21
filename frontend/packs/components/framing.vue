<template>
  <div class="frame">
    <log></log>
    <online-players></online-players>
  </div>
</template>

<script>
  import Log from './framing/log.vue'
  import OnlinePlayers from './framing/online_players.vue'

  export default {
    components: {
      Log: Log,
      OnlinePlayers: OnlinePlayers
    }
  }
</script>