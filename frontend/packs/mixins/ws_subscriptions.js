import ActionCable from 'actioncable'
import config from '../config'

export const ws_subscriptions = {
  methods: {
    subscribeToChatStream() {
      ActionCable.createConsumer(config.cableUrl)
      var that = this
      this.cable.subscriptions.create('ChatChannel', {
        connected() {
          console.log('Connected to ChatChannel')
        },
        disconnected() {
          console.log('Disconnected from ChatChannel')
        },
        received(data) {
          that.$store.commit('addMessage', data)
        }
      })
    },
    subscribeToEventsStream() {
      var that = this
      this.cable.subscriptions.create('EventsChannel', {
        connected() {
          console.log('Connected to EventsChannel')
        },
        disconnected() {
          console.log('Disconnected from EventsChannel')
        },
        received(data) {
          that.$store.commit('addEvent', data)
        }
      })
    }
  },
  computed: {
    cable() {
      return ActionCable.createConsumer(config.cableUrl)
    }
  }
}
