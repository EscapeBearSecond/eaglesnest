<script setup>
 import { ref, reactive } from 'vue'
 import { getPolicyList, createPolicy } from '@/api/policy'
 import AddPolicy from "./component/addPolicy.vue"
   // 搜索组件初始数据

   // 定义列表头部
   const columns = ref([
       { label:'名称', prop:'policyName'},
        { label:'描述', prop:'policyDesc'},
        { label:'在线检测', prop:'onlineCheck'},
        { label:'端口检测', prop:'portScan' },
        { label:'扫描类型', prop:'scanType' , formatter: (row, col, cellValue) => {
         let data = cellValue
        }},
   ])
   // 定义数据
   const tableData = ref([])
   const searchInfo = reactive({
      name: '',
      desc: '',
   })
   const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
   console.log(listQuery)
   const onReset = () => {
      searchInfo.value = {}
   }

   const statusData = ref([
      {
         name: "修改",
         type: "primary",
         icon: "edit",
         handleClick: (scope) => handleEdit(scope.row),
         // visible: (scope, btn) => scope.row.name == 'aaa' 测试按钮根据情况显示 
      }
   ])

   // 获取表格数据
   const getTableData = async() => {
      const table = await getPolicyList({page: listQuery.page, pageSize: listQuery.pageSize, ...searchInfo.value })
      if (table.code === 0) {
         tableData.value = table.data.list
         total.value = table.data.total
         page.value = table.data.page
         pageSize.value = table.data.pageSize
      }
   }
   // 调用
   getTableData()

   // 编辑按钮
   function handleEdit() {
      console.log('edit')
   }

   // 翻页按钮
   function pagination() {
      console.log('pagenation')
      getTableData()
   }

   // 新增

   const dialogFormVisible = ref(false)
   const dialogTitle = ref('创建策略')
   const addForm = reactive({
      policyName: '',
      policyDesc: '',
      headlessFlg: '',
      scanType: '',
      scanRate: '',
      policyConfig: [{
         "name": "",
         "kind": "",
         "timeout": "",
         "count": 0,
         "format": "",
         "rateLimit": 0,
         "concurrency": 0
      }],
      "onlineConfig": {
         "use": true,
         "timeout": "5s",
         "count": 1,
         "format": "csv",
         "rateLimit": 1000,
         "concurrency": 1000
      },
      "portScanConfig": {
         "use": true,
         "timeout": "5s",
         "count": 1,
         "format": "csv",
         "ports": "http",
         "rateLimit": 1000,
         "concurrency": 1000
      }
   })
   const policyFormRef = ref(null);
   const createPolicyFun =  ()=> {
      dialogFormVisible.value = true
   }
   const closeDialog = () => {
      dialogFormVisible.value = false
   }
   const enterDialog = (addForm) => {
      
   }

   const resetPlicyData = ()=> {

   }
  
</script>
<template>
   <div>
      <div class="gva-search-box">
         <el-form
         ref="searchForm"
         :inline="true"
         :model="searchInfo"
         >
         <el-form-item label="名称">
            <el-input
               v-model="searchInfo.policyName"
               placeholder="策略名称"
            />
         </el-form-item>
         <el-form-item label="描述">
            <el-input
               v-model="searchInfo.desc"
               placeholder="描述"
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
      <div class="gva-table-box">
      <div class="gva-btn-list">
         <el-button
            type="primary"
            icon="plus"
            @click="createPolicyFun"
         >新增策略</el-button>
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
      <el-dialog
         v-model="dialogFormVisible"
         :title="dialogTitle"
         style="padding-left:50px"
      >
      <AddPolicy></AddPolicy>
      </el-dialog>      
   </div>
</template>
