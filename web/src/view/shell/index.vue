<template>
  <div class="view-content" ref="viewContent" v-loading="loading">
    <iframe ref="iframe" :src="iframeSrc" frameborder="0" width="100%" :height="iframeHeight"></iframe>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, ref, watchEffect } from 'vue';

const iframeSrc = ref(import.meta.env.VITE_DKY_URL + '/view/shell/project/list?username=admin&password=yigan2024');
const viewContent = ref(null);
const iframeHeight = ref(0);
const iframe = ref(null)
const loading = ref(true)

// 在页面加载后计算高度
onMounted(() => {
  updateIframeHeight();
  window.addEventListener('resize', updateIframeHeight);
  iframe.value.onload = () => {
    loading.value = false
  }
});

// 在页面卸载前移除事件监听器
onUnmounted(() => {
  window.removeEventListener('resize', updateIframeHeight);
});

// 更新 iframe 的高度
function updateIframeHeight() {
  const windowHeight = window.innerHeight;
  let paddingAndBorder = 0;

  if (viewContent.value) { // 检查元素是否存在
    paddingAndBorder = getPaddingAndBorder(viewContent.value);
  }

  const adjustedHeight = windowHeight - paddingAndBorder;
  iframeHeight.value = adjustedHeight;
}

// 计算 padding 和 border 的总和
function getPaddingAndBorder(element) {
  const style = window.getComputedStyle(element);
  const paddingTop = parseFloat(style.paddingTop);
  const paddingBottom = parseFloat(style.paddingBottom);
  const borderTopWidth = parseFloat(style.borderTopWidth);
  const borderBottomWidth = parseFloat(style.borderBottomWidth);
  return paddingTop + paddingBottom + borderTopWidth + borderBottomWidth;
}

// 监听 iframeSrc 变化时更新高度
watchEffect(() => {
  updateIframeHeight();
});
</script>

<style scoped>
.view-content {
  position: relative;
  width: 100%;
  height: var(--iframe-height); /* 使用 CSS 变量来动态设置高度 */
  overflow: hidden;
  padding: 5px;
}
</style>