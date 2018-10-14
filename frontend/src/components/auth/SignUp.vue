<template>
  <div class="container sign_up">
    <h1>Sign up</h1>
    <form @submit="sendForm" action="#" method="post">
      <div class="form-group">
        <label for="exampleInputEmail1">Email address</label>
        <input v-model="formData.email" type="email" class="form-control" id="exampleInputEmail1" aria-describedby="emailHelp" placeholder="Enter email">
      </div>
      <div class="form-group">
        <label for="exampleInputPassword1">Password</label>
        <input v-model="formData.password" type="password" class="form-control" id="exampleInputPassword1" placeholder="Password">
      </div>
      <button type="submit" class="btn btn-primary">Submit</button>
    </form>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'SignUp',
  data () {
    return {
      msg: 'WEEEWWCOMMME',
      formData: {
        email: '',
        password: ''
      }
    }
  },
  methods: {
    sendForm: function() {
      var link = '/auth/sign_up'
      axios.post(link, {
        params: this.formData
      })
      .then(response => {
        console.log(e.response)
      })
      .catch(e => {
        console.log(e.response)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h1, h2 {
  font-weight: normal;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.sign_up {
  max-width: 300px;
}
</style>
