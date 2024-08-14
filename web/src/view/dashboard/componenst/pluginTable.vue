<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="ranking" label="序号" width="80" align="center" />
      <el-table-column prop="title" label="类型" show-overflow-tooltip>
      </el-table-column>
      <el-table-column prop="click_num" label="数量" width="100" />
      <el-table-column prop="hot" label="占比" width="100" />
    </el-table>
  </div>
</template>

<script setup>
const tableData = [
  {
    ranking: 1,
    title : "CVE-2021-1102",
    click_num : 523,
    hot : 263,
  },
  {
    ranking: 2,
    title : "CVE-2021-33044",
    click_num : 416,
    hot : 223
  },
  {
    ranking: 3,
    title : "CVE-2013-6117",
    click_num : 337,
    hot : 176,
  },
  {
    ranking: 4,
    title : "CVE-2020-12062",
    click_num : 292,
    hot : 145,
  },
  {
    ranking: 5,
    title : "CVE-2014-0160",
    click_num : 173,
    hot : 110,
  }
]
</script>

<style scoped lang="scss">

</style>
