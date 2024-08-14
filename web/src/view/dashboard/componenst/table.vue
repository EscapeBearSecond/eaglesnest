<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column prop="ranking" label="排名" width="80" align="center" />
      <el-table-column prop="title" label="IP" show-overflow-tooltip />
      <el-table-column prop="critical" label="超危" width="100" />
      <el-table-column prop="high" label="高危" width="100" />
      <el-table-column prop="medium" label="中危" width="100" />
    </el-table>
  </div>
</template>

<script setup>
const tableData = [
  {
    ranking: 1,
    title : "192.168.1.1",
    critical : 263,
    high : 20,
    medium : 63
  },
  {
    ranking: 2,
    title : "192.168.1.1",
    critical : 193,
    high : 263,
    medium : 33
  },
  {
    ranking: 3,
    title : "192.168.1.1",
    critical : 181,
    high : 30,
    medium : 23
  },
  {
    ranking: 4,
    title : "192.168.1.1",
    critical : 166,
    high : 73,
    medium : 13
  },
  {
    ranking: 5,
    title : "192.168.1.1",
    critical : 106,
    high : 73,
    medium : 163
  },
 
]
</script>

<style scoped lang="scss">

</style>
