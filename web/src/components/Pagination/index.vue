<template>
    <el-pagination
      :class="{ 'single-page': singlePage, 'show': isShow, 'fix-foot': fixFoot }"
      :total="total"
      :current-page.sync="currentPage"
      :page-size.sync="pageSize"
      :page-sizes="pagesizes"
      :layout="layout"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </template>
  
  <script>
  import { defineComponent, ref, computed, watch, onActivated, onDeactivated } from 'vue';
  import { ElPagination } from 'element-plus';
  
  export default defineComponent({
    name: 'PageNation',
    components: {
      ElPagination,
    },
    props: {
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
        default: () => [1, 2, 3, 4],
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
    },
    setup(props, { emit }) {
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
  
      onActivated(() => {
        setTimeout(() => {
          isShow.value = true;
        }, 500);
      });
  
      onDeactivated(() => {
        isShow.value = false;
      });
  
      return {
        isShow,
        currentPage,
        pageSize,
        singlePage,
        handleSizeChange,
        handleCurrentChange,
      };
    },
  });
  </script>
  
  <style scoped>
  /* 添加你的样式 */
  </style>
  