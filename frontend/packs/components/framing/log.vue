<template>
  <div class="log col-md-8">
    <button @click="switchComponent" class="btn btn-default switch">Switch to {{ inactiveComponent }}</button>
    <events v-if="inactiveComponent != 'events'"></events>
    <chat v-if="inactiveComponent != 'chat'"></chat>
  </div>
</template>

<script>
  import Events from './log/events.vue'
  import Chat from './log/chat.vue'

  export default {
    data: function() {
      return {
        inactiveComponent: 'events'
      }
    },
    methods: {
      switchComponent() {
        if(this.inactiveComponent == 'events') {
          this.inactiveComponent = 'chat'
        } else {
          this.inactiveComponent = 'events'
        }
      }
    },
    components: {
      Events: Events,
      Chat: Chat
    }
  }
</script>