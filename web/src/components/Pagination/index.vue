<template>
  <el-pagination
    :current-page="page"
    :page-size="pageSize"
    :page-sizes="[10, 30, 50, 100]"
    :total="total"
    layout="total, sizes, prev, pager, next, jumper"
    @current-change="handleCurrentChange"
    @size-change="handleSizeChange"
  />
</template>

<script setup>
import { ref, defineProps, defineEmits } from 'vue';

const props = defineProps({
  page: {
    type: Number,
    default: 1
  },
  pageSize: {
    type: Number,
    default: 10
  },
  total: {
    type: Number,
    required: true
  }
});

const emits = defineEmits(['size-change', 'current-change']);

const handleSizeChange = (val) => {
  emits('size-change', val);
};

const handleCurrentChange = (val) => {
  emits('current-change', val);
};
</script>