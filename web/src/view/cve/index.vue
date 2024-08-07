<template>
    <div>
      <div class="gva-table-box">
        <div class="gva-search-box">
         <el-form
           ref="searchForm"
           :inline="true"
           :model="searchInfo"
         >
           <el-form-item label="名称">
             <el-input
               v-model="searchInfo.cveName"
               placeholder="请输入漏洞名称"
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
          :columns="tableColumns"
          :tableData="tableData"
          :listQuery="listQuery"
          :statusData="statusData"
          :pagination="handleCurrentChange"
          :index="true"
        >
        <template v-slot:customSeverity="slotProps">
        <!-- 自定义的字段 -->
        <span>
          <el-tag effect="dark" >{{ slotProps.row.severity }}</el-tag>
        </span>
      </template>
        </advance-table>
  
      </div>
    </div>
  </template>
  
  <script setup>
  import { getCveList } from "@/api/cve.js"
  import { ref, watch, reactive } from 'vue'
  
  defineOptions({
    name: 'CveData',
  })

  const page = ref(1)
  const tableData = ref([])
  const listQuery = reactive({
     page : 1,
     total: 0,
     pageSize: 10,
  })
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }
  
  // 操作按钮
  const statusData = reactive([])
  const searchInfo = ref({})
  //查询
  const onSubmit = () => {
     listQuery.page = 1
     getTableData()
   }
  //重置
   const onReset = () => {
     searchInfo.value = {}   
   }
  
  // 查询
  const getTableData = async() => {
    const table = await getCveList({
        page: listQuery.page,
        pageSize: listQuery.pageSize,
        ...searchInfo,
      });
      if (table.code === 0) {
        tableData.value = table.data.list;
        listQuery.total = table.data.total;
        listQuery.page = table.data.page;
        listQuery.pageSize = table.data.pageSize;
      }
  }
  getTableData()
  
  const tableColumns = reactive([
      { label:'名称', prop:'name'},
      { label:'等级', prop:'severity', slot: 'customSeverity' },
      { label:'关联模板', prop:'templateId'},
  ])  
  </script>
  
  <style lang="scss">
  </style>
  