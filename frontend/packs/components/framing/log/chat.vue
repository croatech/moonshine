<template>
  <div class="chat">
    <div class="input row">
      <input type="text" v-model="messageData.text" placeholder="your message here" class="message-input form-control" required="true" v-on:keyup.enter="submitMessage">
      <div v-if="recipient != ''" class="recipient">
        <p class="for">Message for</p>
        <b>{{ recipient.name }}</b>
        <p class="remove" @click="removeRecipient()">x</p>
      </div>
    </div>

    <div class="messages row">
      <div class="message" v-for="message in messages" v-if="messageBelongsToRecipient(message)">
        <b class="name" @click="setRecipient(message.player)"> {{ message.player.name }}:</b>
        <i class="recipient" v-if="checkForPrivateMessage(message)" @click="setRecipient(message.recipient)">to {{ message.recipient.name }}</i>
        <p :class="'text ' + setPrivateClass(message)">{{ message.text }}</p>
      </div>
    </div>
  </div>
</template>

<script>
  import axios from 'axios'

  export default {
    data: function() {
      return {
        messageData: {
          text: ''
        },
      }
    },
    methods: {
      submitMessage() {
        var link = '/messages'
        axios.post(link, { message:
          {
            text: this.messageData.text,
            recipient_id: this.$store.state.recipient.id
          }
        })
        .then(response => {
          this.messageData.text = ''
        })
        .catch(e => {
          console.log(e.response.data)
        })
      },
      setRecipient(recipient) {
        if(this.player.id == recipient.id) {
          return
        }
        this.$store.commit('setRecipient', recipient)
      },
      removeRecipient() {
        this.$store.commit('setRecipient', '')
      },
      messageBelongsToRecipient(message) {
        if(message.recipient == null || this.isPrivateMessage(message)) {
          return true
        }
        return false
      },
      setPrivateClass(message) {
        if(this.checkForPrivateMessage(message)) {
          return 'private'
        }
      },
      checkForPrivateMessage(message) {
        if(message.recipient == null) {
          return
        }
        if(this.isPrivateMessage(message)) {
          return 'private'
        }
      },
      isPrivateMessage(message) {
        if(message.recipient.id == this.player.id || message.player.id == this.player.id) {
          return true
        }
      }
    },
    computed: {
      player() {
        return this.$store.state.player
      },
      messages() {
        return this.$store.state.messages
      },
      recipient() {
        return this.$store.state.recipient
      }
    }
  }
</script>