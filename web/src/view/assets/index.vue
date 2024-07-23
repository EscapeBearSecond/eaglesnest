<script setup>
import { ref, reactive } from 'vue' 
const searchInfo = reactive({
    assetName:''
})
const tableColumns = reactive([
       { label:'名称', prop:'assetName'},
       { label:'IP', prop:'assetIP'},
       { label:'归属', prop:'assetArea'},
       { label:'类型', prop:'assetType'},
       { label:'价值', prop:'manufacturer'},
       { label:'备注', prop:'desc'},
    ])
const tableData = ref([])
const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
const statusData = reactive([
  {
      name: "删除",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleDel(scope.row), 
  }
])

const addDialogFlag = ref(false)
const dialogTitle = ref('新增资产')

const onSubmit = (searchInfo) => {}

const onReset = () => {
  searchInfo.assetName.value = ""
}

const  createAsset = ()=> { 
  addDialogFlag.value = true;
}
const  handleDel = (e) => { console.log(e);}

const pagination = (listQuery)=> {

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
               v-model="searchInfo.assetName"
               placeholder="资产名称"
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
          <el-button type="primary" icon="plus" @click="createAsset">新增资产</el-button>
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
      <el-dialog
    v-model="addDialogFlag"
    :title="dialogTitle"
    style="padding-left:50px"
    @update:modelValue="val => addDialogFlag = val"
  >
    <div>
      新增资产
    </div>
  </el-dialog>
  </div>

</template>
<style lang='scss' scoped>
  
</style>
