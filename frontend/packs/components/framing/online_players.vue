<template>
  <div class="players-online col-md-4">
    <h3>Players online:</h3>
    <div class="player" v-for="player in players">
      <p class="name">
        <b @click="setRecipient(player)">{{ player.name }}</b>
        [{{ player.level }}]
      </p>
    </div>
  </div>
</template>

<script>
  import axios from 'axios'

  export default {
    data: function() {
      return {
        players: {}
      }
    },
    methods: {
      getOnlinePlayers() {
        var link = '/players/online'
        axios.get(link)
        .then(response => {
          console.log(response)
          this.players = response.data
        })
        .catch(e => {
          console.log(e.response.data)
        })
      },
      setRecipient(player) {
        if(player.id == this.player.id) {
          return
        }
        this.$store.commit('setRecipient', player)
      }
    },
    created: function() {
      this.getOnlinePlayers()
    },
    computed: {
      player() {
        return this.$store.state.player
      }
    }
  }
</script>