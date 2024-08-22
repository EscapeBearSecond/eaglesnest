<template>
    <div class="view-content" ref="viewContent">
      <iframe ref="iframe" :src="iframeSrc" frameborder="0" width="100%" height="600px"></iframe>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const iframeSrc = ref('http://106.15.237.118:10200/view/shell/problemReport?username=admin&password=yigan2024')

  const viewContent = ref(null)

  </script>

<style lang="scss">

</style>