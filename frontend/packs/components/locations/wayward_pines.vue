<template>
  <div class="wayward-pines">
    <h1 class="text-center">
      {{ location.name }}
    </h1>

    <div v-if="player.active_movement != null" class="center">
      <i>Walking...</i>
    </div>

    <img :src="'../../assets/locations/spinner.gif'" alt="spinner" v-if="!cellsLoaded" class="spinner center">
    <div :class="'map center ' + determineMapClass()">
      <div class="cells">
        <div v-for="cell in cells" :class="'cell ' + checkForInactiveStyle(cell)" @click=changeCell(cell)>
          <img :src="'../../assets/locations/wayward_pines/cells/' + cell.slug + '.png'" :alt="cell.slug">

          <div class="current_position" v-if="checkForCurrentCell(cell)">
            <img :src="'../../assets/locations/point.png'" alt="point">
          </div>

          <div class="arrow" v-if="checkForMovementCell(cell)">
            <img :src="'../../assets/locations/movements/arrow.png'" alt="" :class="'position-' + determineMovementCellClass(cell)">
            <img :src="'../../assets/locations/movements/target.png'" alt="" v-if="checkForLastMovementCell(cell)">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from 'axios'
  import ActionCable from 'actioncable'
  import config from '../../config'

  export default {
    data: function () {
      return {
        location: {},
        cells: {},
        cellsLoaded: false
      }
    },
    methods: {
      getLocation() {
        var link = '/locations/wayward_pines'
        axios.get(link)
        .then(response => {
          this.location = response.data
          this.cells = this.location.children
        })
        .catch(e => {
          console.log(e)
        })
      },
      checkForInactiveStyle(cell) {
        if(cell.inactive == true) {
          return 'inactive'
        }
      },
      checkForCurrentCell(cell) {
        if(cell.id == this.player.location_id) {
          return true
        }
        // Set cellsLoaded = true if cell is the last element in an array
        if(this.cells[this.cells.length - 1] == cell) {
          setTimeout(this.setCellsLoaded, 2000)
        }
      },
      checkForMovementCell: function(cell) {
        if(this.player.active_movement == null) {
          return
        } else {
          var path_array = this.player.active_movement.path
          var index_of_current_cell = path_array.indexOf(String(this.player.location_id))
          var index_of_next_cell = path_array.indexOf(String(cell.id))
          if(this.player.active_movement.path.includes(String(cell.id)) && this.player.location_id != cell.id && index_of_current_cell < index_of_next_cell) {
            return true
          }
        }
      },
      checkForLastMovementCell(cell) {
        if(String(cell.id) == this.player.active_movement.path[this.player.active_movement.path.length - 1]) {
          return true
        }
      },
      determineMapClass() {
        if(!this.cellsLoaded) {
          return 'invisible'
        }
      },
      determineMovementCellClass(cell) {
        var path_array = this.player.active_movement.path
        var index_of_cell = path_array.indexOf(String(cell.id))
        var next_location = path_array[index_of_cell + 1]
        if(next_location == null) {
          return 'none'
        } else {
          return cell.id - next_location
        }
      },
      setCellsLoaded() {
        this.cellsLoaded = true
        // If player has an active movement need to subscribe to stream
        if(this.player.active_movement != null) {
          this.subscribeToStream()
        }
      },
      updatePlayer(data) {
        this.$store.commit('updatePlayer', data)
      },
      changeCell(cell) {
        var link = '/cells/' + cell.id + '/move'
        axios.post(link)
        .then(response => {
          this.updatePlayer(response.data)
          this.subscribeToStream()
        })
        .catch(e => {
          console.log(e)
        })
      },
      subscribeToStream() {
        var that = this
        this.cable.subscriptions.create('MovementChannel', {
          connected() {
            console.log('Connected to MovementChannel')
          },
          disconnected() {
            console.log('Disconnected from MovementChannel')
          },
          received(data) {
            that.updatePlayer(data)
          }
        })
      },
    },
    created: function() {
      this.getLocation()
    },
    computed: {
      player() {
        return this.$store.state.player
      },
      cable() {
        return ActionCable.createConsumer(config.cableUrl)
      }
    }
  }
</script>

<style>
  .invisible {
    display: none;
  }
</style>
