<template>
  <div class="container">
    <div class="winner center" v-if="winner != null">
      <h1 class="text-center">
        Winner: {{ winner.name }}
      </h1>
      <div v-if="fight['winner_type'] == 'Player'">
        <div class="progress">
          <div aria-valuemax="100" aria-valuemin="0" aria-valuenow="80" class="progress-bar progress-bar-danger" role="progressbar" :style="{ width: playerHpPercent }">
            {{ round.player_hp }} / {{ player.hp }} HP
          </div>
        </div>

        <div class="progress">
          <div aria-valuemax="100" aria-valuemin="0" aria-valuenow="80" class="progress-bar progress-bar-warning" role="progressbar" :style="{ width: playerExpPercent }">
            {{ player.exp }} / {{ player.exp_next }} EXP
          </div>
        </div>
      </div>

      <img :src="winner.avatar.url" class="image center avatar" alt="avatar">

      <div v-if="fight.dropped_gold != null" class="center drop">
        <b>Dropped gold:</b>
        <img src="../assets/gold.png" alt="gold">
        {{ fight.dropped_gold }}
      </div>

      <div v-if="fight.dropped_item != null" class="center drop">
        <b>Dropped item:</b>
        <img :src="fight.dropped_item.image.url" alt="item">
      </div>
    </div>

    <div class="fight" v-if="winner == null">
      <div class="characters row">
        <div class="col-md-6">
          <div class="character player">
            <h2>{{ player.name }}[{{ player.level }}]</h2>
            <div class="progress">
              <div aria-valuemax="100" aria-valuemin="0" aria-valuenow="80" class="progress-bar progress-bar-danger" role="progressbar" :style="{ width: playerHpPercent }">
                {{ round.player_hp }} HP
              </div>
            </div>

            <img :src="player.avatar.url" class="image avatar" alt="avatar">

            <div class="points">
              <h4 class="block-title">Defense: {{ attackPoint }}</h4>
              <ul v-for="point in points">
                <li><button @click=setAttackPoint(point) :class="'btn btn-default attack ' + resolveAttackPointClass(point)">{{ point }}</button></li>
              </ul>
            </div>
          </div>
        </div>

        <div class="col-md-6">
          <div class="character bot">
            <h2>{{ bot.name }}[{{ bot.level }}]</h2>
            <div class="progress">
              <div aria-valuemax="100" aria-valuemin="0" aria-valuenow="80" class="progress-bar progress-bar-danger" role="progressbar" :style="{ width: botHpPercent }">
              {{ round.bot_hp }} HP
              </div>
            </div>

            <img :src="bot.avatar.url" class="image avatar" alt="avatar">

            <div class="points">
              <h4 class="block-title">Attack: {{ defensePoint }}</h4>
              <ul v-for="point in points">
                <li><button @click=setDefensePoint(point) :class="'btn btn-default attack ' + resolveDefensePointClass(point)">{{ point }}</button></li>
              </ul>
            </div>
          </div>
        </div>
      </div>
      <div class="row">
        <button @click="attack" class="btn btn-danger btn-attack">Attack</button>
      </div>
    </div>

    <div class="rounds center">
      <div v-for="round in rounds">
        <div v-if="round.player_damage != null">
          <hr/>
          <div class="damage">
            <p>
              {{ player.name }} has dealt
              <span>{{ round.player_damage }}</span>
              damage to {{ bot.name }} in the
              <b>{{ round.player_attack_point }}</b>.
            </p>
          </div>

          <div class="damage">
            <p>
              {{ bot.name }} has dealt
              <span>{{ round.bot_damage }}</span>
              damage to {{ player.name }} in the
              <b>{{ round.bot_attack_point }}</b>.
            </p>
          </div>
          <hr/>
        </div>
      </div>
    </div>
  </div>

</template>

<script>
  import axios from 'axios'
  import { calculate } from '../mixins/calculate'

  export default {
    data: function() {
      return {
        rounds: [],
        fight: {},
        round: {},
        player: {},
        bot: {},
        winner: {},
        points: [],
        attackPoint: '',
        defensePoint: ''
      }
    },
    mixins: [calculate],
    methods: {
      getFight: function() {
        axios.get('/fight')
        .then(response => {
          this.setDataFromResponse(response.data)
        })
        .catch(e => {
          console.log(e)
        })
      },
      attack: function() {
        axios.patch('/fight', {
          defensePoint: this.defensePoint,
          attackPoint: this.attackPoint
        })
        .then(response => {
          this.setDataFromResponse(response.data)
          this.commit_if_gold_dropped()
        })
        .catch(e => {
          console.log(e)
        })
      },
      resolveAttackPointClass: function(point) {
        if(this.attackPoint == point) {
          return 'active'
        }
      },
      resolveDefensePointClass: function(point) {
        if(this.defensePoint == point) {
          return 'active'
        }
      },
      setAttackPoint: function(point) {
        this.attackPoint = point
      },
      setDefensePoint: function(point) {
        this.defensePoint = point
      },
      setDataFromResponse: function(data) {
        this.fight = data['fight']
        this.rounds = this.fight['rounds'].reverse()
        this.round = this.fight['current_round']
        this.player = this.fight['player']
        this.bot = this.fight['bot']
        this.winner = this.fight['winner']
        this.points = data['points']
        this.attackPoint = this.points[this.points.length - 1]
        this.defensePoint = this.points[this.points.length - 1]
      },
      commit_if_gold_dropped: function() {
        if(this.fight.dropped_gold != null) {
          this.$store.commit('increment_gold', this.fight.dropped_gold)
        }
      }
    },
    created: function() {
      this.getFight()
    },
    computed: {
      playerHpPercent: function() {
        return this.percentProgressBar(this.round['player_hp'], this.player['hp'])
      },
      botHpPercent: function() {
        return this.percentProgressBar(this.round['bot_hp'], this.bot['hp'])
      },
      playerExpPercent: function() {
        return this.percentProgressBar(this.player['exp'], this.player['exp_next'])
      },
    }
  }
</script>
