
<script setup>
 import { ref } from 'vue'
   // 搜索组件初始数据

   // 定义列表头部
   const columns = ref([
       { label: '名称', prop: 'name' },
        { label: '年龄', prop: 'age' },
   ])
   // 定义列表数据
   const tableData = ref([])
   const statusData = ref([
      {
         name: "修改",
         type: "primary",
         handleClick: (scope) => handleEdit(scope.row),
         // visible: (scope, btn) => scope.row.name == 'aaa' 测试按钮根据情况显示 
      }
   ])
   //翻页
   const listQuery = ref({
      page: 1,
      limit: 10,
      total: 0,
   })

   const searchInfo = ref({})
   const onReset = () => {
      searchInfo.value = {}
   }

   function getTableData() {
      console.log("table")
   }

   function handleEdit() {
      console.log('edit')
   }

   function pagination() {
      console.log('pagenation')
   }
   // 获取页面数据
   getTableData()
</script>
<template>
     <div class="view-container">
      <div class="gva-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="名称">
          <el-input
            v-model="searchInfo.name"
            placeholder="策略名称"
          />
        </el-form-item>
        <el-form-item>
          <el-button
            type="primary"
            icon="search"
            @click="onSubmit"
          >查询</el-button>
          <el-button
            icon="refresh"
            @click="onReset"
          >重置</el-button>
        </el-form-item>
      </el-form>
    </div>
        <advance-table
            :columns = "columns"
            :tableData="tableData"
            :listQuery="listQuery"
            :statusData="statusData"
            :pagination="pagination"
            :index= "true"
         >
         </advance-table>
     </div>
</template>
