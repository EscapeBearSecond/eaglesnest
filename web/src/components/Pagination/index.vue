<template>
  <el-pagination
    :class="{ 'single-page': singlePage, 'show': isShow, 'fix-foot': fixFoot }"
    :total="total"
    :current-page="currentPage"
    :page-size="pageSize"
    :page-sizes="pagesizes"
    :layout="layout"
    @size-change="handleSizeChange"
    @current-change="handleCurrentChange"
  />
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue';

const props = defineProps({
  page: {
    type: Number,
    default: 1,
  },
  pagesize: {
    type: Number,
    default: 1,
  },
  pagesizes: {
    type: Array,
    default: () => [10, 20, 30, 40],
  },
  layout: {
    type: String,
    default: 'total, sizes, prev, pager, next, jumper',
  },
  total: {
    type: Number,
    required: true,
  },
  autoScroll: {
    type: Boolean,
    default: true,
  },
  fixFoot: {
    type: Boolean,
    default: true,
  },
});

const emit = defineEmits(['pagination']);

const isShow = ref(false);
const currentPage = ref(props.page);
const pageSize = ref(props.pagesize);

watch(
  () => props.page,
  (newVal) => {
    currentPage.value = newVal;
  }
);

watch(
  () => props.pagesize,
  (newVal) => {
    pageSize.value = newVal;
  }
);

const singlePage = computed(() => {
  return pageSize.value >= props.total;
});

const handleSizeChange = (val) => {
  emit('pagination', { page: currentPage.value, limit: val });
  if (props.autoScroll) {
    window.scrollTo(0, 800);
  }
};

const handleCurrentChange = (val = currentPage.value) => {
  emit('pagination', { page: val, limit: pageSize.value });
  if (props.autoScroll) {
    window.scrollTo(0, 800);
  }
};

onMounted(() => {
  setTimeout(() => {
    isShow.value = true;
  }, 500);
});

onUnmounted(() => {
  isShow.value = false;
});
</script>

<style scoped>
/* 添加你的样式 */
</style>