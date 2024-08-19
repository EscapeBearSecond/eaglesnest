<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column type="index" width="80" label="序号" />
      <el-table-column prop="host" label="IP" />
      <el-table-column prop="count" label="总漏洞" width="100" />
      <el-table-column prop="critical" label="超危" width="100" />
      <el-table-column prop="high" label="高危" width="100" />
      <el-table-column prop="medium" label="中危" width="100" />
    </el-table>
  </div>
</template>

<script setup>
defineProps({
  tableData: {
    type: Array,
    default() {
      return [];
    },
  },
})
</script>

<style scoped lang="scss">

</style>
