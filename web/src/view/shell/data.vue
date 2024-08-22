<template>
    <div>
      <iframe ref="iframe" :src="iframeSrc" frameborder="0" width="100%" height="500"></iframe>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const iframeSrc = ref('http://106.15.237.118:10200/view/shell/dataAnalysis?username=admin&password=yigan2024')

  </script>