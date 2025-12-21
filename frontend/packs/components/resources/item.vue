<template>
  <div class="row top">
    <div class="col-md-3">
      <img v-bind:src="item.image.url"/>
    </div>

    <div class="col-md-9">
      <h3>{{ item.name }}</h3>
    </div>
  </div>
</template>

<script>
  export default {
    props: ['item']
  }
</script>
