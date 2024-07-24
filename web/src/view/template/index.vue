<script setup>
import { ref, reactive } from 'vue' 
import districtForm from "./components/assetForm.vue"


const searchInfo = reactive({
    areaName:''
})
const tableColumns = reactive([
    { label:'名称', prop:'templateName'},
    { label:'类型', prop:'templateType'},
    { label:'备注', prop:'templateDesc'},
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
const dialogTitle = ref('新增模板')
const formData = reactive({
  templateName:"",
  templateContent:"",
  templateType:"",
})
const labelPosition = ref('left')
const itemLabelPosition = ref('top')

const formRules = reactive({
  templateName: [
    { required: true, message: '请输入模板名称', trigger: 'blur' }
  ],
  templateType: [
    { required: true, message: '请选择模板类型', trigger: 'blur' }
  ]
});
const onCancel = () => {
  addDialogFlag.value = false
}
const onSubmit = (searchInfo) => {}

const onReset = () => {
  searchInfo.templateName.value = ""
}

const  createAsset = ()=> { 
  addDialogFlag.value = true;
}
const  handleDel = (e) => { console.log(e);}
const  handleEdit = (e) => { console.log(e);}

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
               v-model="searchInfo.templateName"
               placeholder="模板名称"
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
        style="padding: 10px 50px;"
        @update:modelValue="val => addDialogFlag = val"
      >
    <div>
      <district-form 
        :form="formData"
        :rules="formRules"
        :label-position="labelPosition"
        :item-label-position="itemLabelPosition"
        @submit="onSubmit"
        @cancel="onCancel"
      />
    </div>
  </el-dialog>
  </div>

</template>
<style lang='scss' scoped>
  
</style>
