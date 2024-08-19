<!--
    @auther: bypanghu<bypanghu@163.com>
    @date: 2024/5/8
!-->

<template>
  <div>
    <el-table :data="tableData" stripe style="width: 100%">
      <el-table-column type="index" width="100" label="序号" />
      <el-table-column prop="name" label="类型" show-overflow-tooltip>
      </el-table-column>
      <el-table-column prop="count" label="数量" width="100" />
      <el-table-column prop="severity" label="等级" width="100" >
        <template #default="scope">
          <el-tag
            :type="getTagClass(scope.row.severity)"
            disable-transitions
            >{{ getTag(scope.row.severity) }}</el-tag
          >
        </template>
      </el-table-column>
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

const getTag = (tag) => {
    switch(tag) {
        case 'high':
          return '高危';
        case 'medium':
          return '中危';
        case 'low':
          return '低危';
        default:
          return '高危';
    }
}

const getTagClass= (tag) => {
  switch(tag) {
        case 'high':
          return 'warning';
        case 'medium':
          return 'info';
        case 'low':
          return 'primary';
        default:
          return 'danger';
    }
}
</script>

<style scoped lang="scss">

</style>
