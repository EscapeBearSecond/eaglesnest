<template>
   <div class="authority">
     <div class="gva-search-box">
       <el-form
         ref="searchForm"
         :inline="true"
         :model="searchInfo"
       >
         <el-form-item label="名称">
           <el-input
             v-model="searchInfo.policyName"
             placeholder="请输入策略名称"
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
           @click="addPolicy()"
         >新增策略</el-button>
       </div>
       <advance-table
         :columns="tableColumns"
         :tableData="tableData"
         :listQuery="listQuery"
         :statusData="statusData"
         :pagination="pagination"
         :index="true"
       >
       <template v-slot:custOnline="slotProps">
        <!-- 自定义的字段 -->
        <span v-if="slotProps.row.onlineCheck"> 
          <el-icon :size="20" color="#67C23A" >
            <CircleCheck />
          </el-icon>
        </span>
        <span v-else>
          <el-icon :size="20" color="#F56C6C" >
            <CircleClose />
          </el-icon>
        </span>
      </template>
      <template v-slot:custPortScan="slotProps">
        <!-- 自定义的字段 -->
        <span v-if="slotProps.row.portScan"> 
          <el-icon :size="20" color="#67C23A" >
            <CircleCheck />
          </el-icon>
        </span>
        <span v-else>
          <el-icon :size="20" color="#F56C6C" >
            <CircleClose />
          </el-icon>
        </span>
      </template>
      <template v-slot:custScanType="slotProps">
        <!-- 自定义的字段 -->
        <span v-for="(item, key) in slotProps.row.scanType" :key="key" style="margin-left: 5px;"> 
          <el-tag
            type="primary"
            effect="dark"
          >
          {{ getTypeName(item) }}</el-tag>
        </span>
      </template>
      </advance-table>
      
     </div>
     <!-- 新增角色弹窗 -->
     <el-drawer
       v-model="dialogFormVisible"
       :title="dialogTitle"
       :before-close="handleClose"
       size="45%"
     >
       <el-form
         ref="authorityForm"
         :model="form"
         :rules="rules"
         style="padding:10px 20px;"
         label-width="100px"
       >
         <el-form-item label="策略名称" :label-position="itemLabelPosition" prop="taskName">
            <el-input v-model="form.policyName" placeholder="请输入扫描名称" />
         </el-form-item>
         <div style="margin-left: 8px;">
            <label  class="el-form-item__label">策略选项</label>
            <el-collapse v-model="activeNames" style="padding-left: 20px;" accordion>
               <el-collapse-item title="在线检测" name="1">
                  <el-form-item label="策略状态" :label-position="itemLabelPosition" class="one-lab">
                        <el-checkbox v-model="form.onlineConfig.use" label="开启"  size="large" />
                  </el-form-item>
                  <div v-if="form.onlineConfig.use">
                     <el-form-item  label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.onlineConfig.concurrency" />
                     </el-form-item>
                     <el-form-item label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.onlineConfig.timeout" />
                     </el-form-item>
                     <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.onlineConfig.count" />
                     </el-form-item>
                     <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.onlineConfig.rateLimit" />
                     </el-form-item>
                  </div>
               </el-collapse-item>
               <el-collapse-item title="端口检测" name="2">
                  <el-form-item label="策略状态" :label-position="itemLabelPosition" class="one-lab">
                     <el-checkbox v-model="form.portScanConfig.use" label="开启"  size="large" />
                  </el-form-item>
                  <div v-if="form.portScanConfig.use">
                     <el-form-item label="端口范围" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.ports" />
                     </el-form-item>
                     <el-form-item  label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.rateLimit" />
                     </el-form-item>
                     <el-form-item label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.timeout" />
                     </el-form-item>
                     <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.count" />
                     </el-form-item>
                     <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.rateLimit" />
                     </el-form-item>
                     <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
                        <el-input v-model="form.portScanConfig.rateLimit" />
                     </el-form-item>
                  </div>
               </el-collapse-item>
            </el-collapse>
         </div>
         <div style="margin-left: 8px;padding-top: 5px;">
            <label  class="el-form-item__label">其他选项：</label>
             <template v-for="(item, index) in form.policyConfig" :key="index" >
                <label style="display: block;margin-left: 20px;">配置{{ index+1 }}</label>
                <div style="margin-left: 40px;">
                  <el-form-item label="类型" :label-position="itemLabelPosition" class="sec-lab">
                    <el-select v-model="item.kind" placeholder="请选择扫描类型">
                        <el-option
                          v-for="item in typeNameList"
                          :key="item.value"
                          :label="item.label"
                          :value="item.value"
                          :disabled="item.disabled"
                        />
                      </el-select>
                  </el-form-item>
                  <el-form-item label="模板" :label-position="itemLabelPosition" class="sec-lab">
                     <el-input v-model="item.ports" />
                  </el-form-item>
                  <el-form-item label="并发数" :label-position="itemLabelPosition" class="sec-lab">
                     <el-input v-model="item.concurrency" />
                  </el-form-item>
                  <el-form-item label="超时" :label-position="itemLabelPosition" class="sec-lab">
                     <el-input v-model="item.timeout" />
                  </el-form-item>
                  <el-form-item label="限流速度" :label-position="itemLabelPosition" class="sec-lab">
                     <el-input v-model="item.rateLimit" />
                  </el-form-item>
                  <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
                     <el-input v-model="item.count" />
                  </el-form-item>
                </div>
             </template>
            <el-button type="primary" icon="el-icon-plus" @click="addTmpData">新增模板</el-button>
         </div>
         <div style="margin-left: 15px;margin-top: 10px;">
            <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
               <el-input type="textarea" :rows="3" v-model="form.policyDesc" />
            </el-form-item>
         </div>
         </el-form>
       <template #footer>
         <div class="dialog-footer">
           <el-button @click="closeDialog">取 消</el-button>
           <el-button
             type="primary"
             @click="enterDialog"
           >确 定</el-button>
         </div>
       </template>
     </el-drawer>
   </div>
 </template>
 
 <script setup>
 import {
   getTaskList,
   createTask,
   stopTask,
   delTask,
 } from '@/api/task.js'
import { getPolicyList, createPolicy } from '@/api/policy.js';

 import { ref,reactive } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 
 defineOptions({
   name: 'Policy'
 })
 const activeNames = ref([1])
 const dialogType = ref('add')
 
 const dialogTitle = ref('新增策略')
 const dialogFormVisible = ref(false)
 const apiDialogFlag = ref(false)
 const labelPosition = ref('right')
const itemLabelPosition = ref('top')
 const form = ref({
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
});

const marks = ref({
  '快速': '串行',
  33: '缓慢',
  66: '适中',
  100: '快速'
})

//'资产发现', '漏洞扫描', '弱口令'
const typeNameList = reactive([
  {id: '1', label: '资产发现', value:'1', disabled: false},
  {id: '2', label: '漏洞扫描', value:'2', disabled: false},
  {id: '3', label: '弱口令', value:'3', disabled: false}
])

 const tableColumns = ref([
       { label:'名称', prop:'policyName'},
        { label:'描述', prop:'policyDesc'},
        { label:'在线检测', prop:'onlineCheck', slot: 'custOnline'},
        { label:'端口检测', prop:'portScan' , slot: 'custPortScan'},
        { label:'扫描类型', prop:'scanType' , slot: 'custScanType'},
   ])
 const rules = ref({
   taskName: [
     { required: true, message: '请输入扫描名称', trigger: 'blur' }
   ],
   targetIp: [
     { required: true, message: '请输入扫描IP', trigger: 'blur' }
   ],
   status: [
     { required: true, message: '请选择扫描状态', trigger: 'blur' }
   ],
   policyId: [
     { required: true, message: '请选择策略模板', trigger: 'blur' }
   ],
   // date: [
   //   { required: true, message: '请选择定时执行时间', trigger: 'blur' }
   // ],
   taskPlan: [
     { required: true, message: '请选择执行方式', trigger: 'blur' }
   ]
 })

 const listQuery = ref({
     page: 1,
     total: 0,
     pageSize: 10
 })
 const statusData = ref([
   {
      name: "修改",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleEdit(scope.row),
   }
 ])
 
 const tableData = ref([])
 const searchInfo = ref({})
 
 // 查询
 const getTableData = async() => {
   const table = await getPolicyList({ page: listQuery.value.page, pageSize: listQuery.value.pageSize, ...searchInfo.value })
   if (table.code === 0) {
     tableData.value = table.data.list
     listQuery.value.total = table.data.total
     listQuery.value.page = table.data.page
     listQuery.value.pageSize = table.data.pageSize
   }

   console.log(table.data.list);
 }
 
 getTableData()

 // 初始化表单
 const authorityForm = ref(null)
 const tmpForm = ref(null)
 const initForm = () => {
   if (authorityForm.value) {
     authorityForm.value.resetFields()
   }
   form.value = {
    policyName: '',
    policyDesc: '',
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
   }
 }
 // 关闭窗口
 const closeDialog = () => {
   initForm()
   dialogFormVisible.value = false
   apiDialogFlag.value = false
 }
 // 确定弹窗
 
 const enterDialog = () => {
   authorityForm.value.validate(async valid => {
     if (valid) {
       switch (dialogType.value) {
         case 'add':
           {
             const res = await createPolicy(form.value)
             if (res.code === 0) {
               ElMessage({
                 type: 'success',
                 message: '添加成功!'
               })
               getTableData()
               closeDialog()
             }
           }
           break
       }
 
       initForm()
       dialogFormVisible.value = false
     }
   })
 }
 const setOptions = () => {
   
 }
 
 // 增加角色
 const addPolicy = (parentId) => {
   initForm()
   dialogTitle.value = '新增策略'
   dialogType.value = 'add'
   form.value.parentId = parentId
   setOptions()
   dialogFormVisible.value = true
 }
 
 
 const pagination = () => {
   getTableData()
 }
 
 const handleEdit = (row) => {
   ElMessageBox.confirm('此操作将停止该任务, 是否继续?', '提示', {
     confirmButtonText: '确定',
     cancelButtonText: '取消',
     type: 'warning'
   })
     .then(async() => {
       const res = await stopTask({ id: row.id })
       if (res.code === 0) {
         ElMessage({
           type: 'success',
           message: '删除成功!'
         })
         getTableData()
       }
     })
     .catch(() => {
       ElMessage({
         type: 'info',
         message: '已取消删除'
       })
     })
 }
 
 const handleClose = () => {

 }

 const onSubmit = () => {
   listQuery.page = 1
   getTableData()
 }

 const onReset = () => {
   searchInfo.value = {}
 }
 
 const addTmpData = () => {
  if(form.value.policyConfig.length <= 3) {
    form.value.policyConfig.push({
      "name": "",
      "kind": "",
      "timeout": "",
      "count": 0,
      "format": "",
      "rateLimit": 0,
      "concurrency": 0
    })
  }
 }

 const getTypeName = (type) => {
      if(type) {
        const typeNameList = ['未知', '资产发现', '漏洞扫描', '弱口令']
        return typeNameList[type]
      }
 }
 
 </script>
 
 <style lang="scss">
 .authority {
   .el-input-number {
     margin-left: 15px;
     span {
       display: none;
     }
   }
 }
 .tree-content{
   margin-top: 10px;
   height: calc(100vh - 158px);
   overflow: auto;
 }
 </style>
 