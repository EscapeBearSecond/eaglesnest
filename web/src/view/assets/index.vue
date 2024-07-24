<script setup>
import { ref, reactive } from 'vue' 
import districtForm from "./components/districtForm.vue"

const searchInfo = reactive({
    areaName:''
})
const tableColumns = reactive([
       { label:'名称', prop:'areaName'},
       { label:'IP范围', prop:'areaIP'},
       { label:'备注', prop:'areaDesc'},
    ])
const tableData = ref([])
const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
const statusData = reactive([
   {
      name: "修改",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleEdit(scope.row), 
  },
  {
      name: "删除",
      type: "primary",
      icon: "delete",
      handleClick: (scope) => handleDel(scope.row), 
  }
])

const addDialogFlag = ref(false)
const dialogTitle = ref('新增区域')
const formData = reactive({
  areaName:"",
  areaIP:"",
  areaDesc:"",
})
const labelPosition = ref('left')
const itemLabelPosition = ref('top')

const rules = reactive({
  areaName: [
    { required: true, message: '请输入区域名称', trigger: 'blur' }
  ],
  areaIP: [
    { required: true, message: '请输入区域IP范围', trigger: 'blur' }
  ]
});
const onCancel = () => {
  addDialogFlag.value = false
}
const onSubmit = (searchInfo) => {}

const onReset = () => {
  searchInfo.areaName.value = ""
}

const createAsset = ()=> { 
  addDialogFlag.value = true;
}
const handleDel = (e) => { console.log(e);}
const handleEdit = (e) => { console.log(e);}
const onSubmitDialog = (formValues) => {
  console.log('表单数据：', formValues);
  // 执行保存逻辑
};
const pagination = (listQuery)=> {}
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
               v-model="searchInfo.areaName"
               placeholder="区域名称"
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
          <el-button type="primary" icon="plus" @click="createAsset">新增区域</el-button>
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
        style="padding:50px;"
        width="35%"
        @update:modelValue="val => addDialogFlag = val"
      >
    <div>
      <district-form 
        :form="formData"
        :rules="rules"
        :label-position="labelPosition"
        :item-label-position="itemLabelPosition"
        @submit="onSubmitDialog"
        @cancel="onCancel"
      />
    </div>
  </el-dialog>
  </div>

</template>
<style lang='scss' scoped>
  
</style>
