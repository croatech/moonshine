<template>
  <div id="app">
    <link rel="stylesheet" type="text/css" href="./static/css/bootstrap-grid.min.css">
    <link rel="stylesheet" type="text/css" href="./static/css/bootstrap.min.css">
    <link rel="stylesheet" type="text/css" href="./static/css/auth.css">
    <router-view/>
  </div>
</template>

<script>
  export default {
    name: 'Sunlight'
  }
</script>

<style>
  #app {
    font-family: 'Avenir', Helvetica, Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    text-align: center;
    color: #2c3e50;
    margin-top: 60px;
  }
</style>
