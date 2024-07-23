<script setup>
import { ref, reactive } from 'vue' 
const searchInfo = reactive({
    name:''
})
const tableColumns = reactive([
       { label:'名称', prop:'taskName'},
       { label:'描述', prop:'taskDesc'},
       { label:'目标', prop:'targetIp'},
       { label:'任务计划', prop:'taskPlan'},
       { label:'状态', prop:'status'},
       { label:'策略', prop:'planConfig'},
    ])
const tableData = ref([])
const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
const statusData = reactive([
  {
      name: "停止",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleStop(scope.row),
      // visible: (scope, btn) => scope.row.name == 'aaa' 测试按钮根据情况显示 
  },
  {
      name: "删除",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleDel(scope.row),
      // visible: (scope, btn) => scope.row.name == 'aaa' 测试按钮根据情况显示 
  },
  {
      name: "生成报告",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleReport(scope.row),
      // visible: (scope, btn) => scope.row.name == 'aaa' 测试按钮根据情况显示 
  }
])

const onSubmit = (searchInfo) => {}

const onReset = () => {
  searchInfo.name = ""
}

const  createScan = ()=> {}

const  handleStop = (e) => { console.log(e);}
const  handleDel = (e) => { console.log(e);}
const  handleReport = (e) => { console.log(e);}

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
               v-model="searchInfo.name"
               placeholder="任务名称"
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
            @click="createScan"
         >创建扫描</el-button>
      </div>
      <advance-table
         :columns = "tableColumns"
         :tableData="tableData"
         :listQuery="listQuery"
         :statusData="statusData"
         :pagination="pagination"
         :index= "true"
      >
      </advance-table>
    </div>
    </div>
</template>
<style lang='scss' scoped>
  
</style>
