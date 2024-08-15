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
         <el-form-item label="扫描类型">
          <el-select v-model="searchInfo.scanType" placeholder="请选择扫描类型" multiple collapse-tags>
            <el-option label="资产发现" value="1" />
            <el-option label="漏洞扫描" value="2" />
            <el-option label="弱口令" value="3" />
          </el-select>
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
           @click="handleAdd()"
         >新增策略</el-button>
       </div>
       <advance-table
         :columns="tableColumns"
         :tableData="tableData"
         :listQuery="listQuery"
         :statusData="statusData"
         :changePageSize="changeSize"
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
     <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">策略</span>
          <div>
            <el-button @click="closeDialog">取 消</el-button>
            <el-button
              type="primary"
              @click="enterDialog"
            >确 定</el-button>
          </div>
        </div>
      </template>
       <el-form
         ref="formRef"
         :model="form"
         :rules="rules"
         style="padding:10px 20px;"
         label-width="100px"
       >
         <el-form-item label="策略名称" :label-position="itemLabelPosition" prop="policyName">
            <el-input v-model="form.policyName" placeholder="请输入策略名称" />
         </el-form-item>
         <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
            <el-input type="textarea" :rows="3" v-model="form.policyDesc" />
        </el-form-item>
         <div style="margin-left: 8px;">
            <label  class="el-form-item__label">策略配置</label>
            <el-collapse v-model="activeNames" style="padding-left: 40px;" accordion>
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
         <div style="margin: 10px 0 10px 10px;padding-top: 5px;">
            <label  class="el-form-item__label">模板配置：</label>
            <!-- <el-button type="primary" icon="Plus" @click="addTmpData" v-if="form.policyConfig.length < 3">新增</el-button>
            <el-button type="primary" icon="Delete" @click="delTmpData" v-if="form.policyConfig.length >= 2">删除</el-button> -->
            <el-button size="small" @click="addTab(editableTabsValue)">
              add tab
            </el-button>
            <el-button size="small" @click="removeTab(editableTabsValue)">
              del tab
            </el-button>
            <el-tabs
              v-model="editableTabsValue"
              type="card"
              class="demo-tabs"
              @tab-remove="removeTab"
            >
              <el-tab-pane
                v-for="item in editableTabs"
                :key="item.name"
                :label="item.title"
                :name="item.name"
              >
                {{ item.content }}
              </el-tab-pane>
            </el-tabs>
         </div>
         </el-form>
     </el-drawer>
   </div>
 </template>
 
 <script setup>
 import { getPolicyList, createPolicy, deletePolicy, updatePolicy, getPolicyId } from '@/api/policy.js';
 import { getTemplateList } from '@/api/template.js';
 import router from '@/router/index'
 import { ref,reactive } from 'vue'
 import { ElMessage, ElMessageBox } from 'element-plus'
 
 defineOptions({
   name: 'Policy'
 })

 /* test */
 let tabIndex = 1
 const editableTabsValue = ref('1')
 const removeTab = (targetName) => {
    const tabs = editableTabs.value
    let activeName = editableTabsValue.value
    if (activeName === targetName) {
      tabs.forEach((tab, index) => {
        if (tab.name === targetName) {
          const nextTab = tabs[index + 1] || tabs[index - 1]
          if (nextTab) {
            activeName = nextTab.name
          }
        }
      })
    }

    editableTabsValue.value = activeName
    editableTabs.value = tabs.filter((tab) => tab.name !== targetName)
 }
 const addTab  = (targetName) => {
  const newTabName = `${++tabIndex}`
  editableTabs.value.push({
    title: '配置' + tabIndex,
    name: newTabName,
    content: 'New Tab content',
  })
  editableTabsValue.value = newTabName
 }
 const editableTabs = ref([
  {
    title: '配置',
    name: 'setting',
    content: 'Tab 1 content',
  }
])


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
          "timeout": "5s",
          "count": 1,
          "format": "",
          "rateLimit": 150,
          "concurrency": 150
         }],
         "onlineConfig": {
          "use": true,
          "timeout": "5s",
          "count": 1,
          "format": "csv",
          "rateLimit": 150,
          "concurrency": 150
         },
         "portScanConfig": {
          "use": true,
          "timeout": "5s",
          "count": 1,
          "format": "csv",
          "ports": "http",
          "rateLimit": 150,
          "concurrency": 150
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
      { label:'在线检测', prop:'onlineCheck', slot: 'custOnline',width: '90'},
      { label:'端口检测', prop:'portScan' , slot: 'custPortScan',width: '90'},
      { label:'扫描类型', prop:'scanType' , slot: 'custScanType'},
      { label:'描述', prop:'policyDesc'},
    ])
 const rules = ref({
  policyName: [
     { required: true, message: '请输入策略名称', trigger: 'blur' }
  ],
 })

 const listQuery = ref({
     page: 1,
     total: 0,
     pageSize: 10
 })
 const changeSize = (e) => {
    listQuery.page = 1
    listQuery.pageSize = e
    getTableData()
  }
 const statusData = ref([
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
      handleClick: (scope) => handleDelete(scope.row),
   },
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
 }
 
 getTableData()

 // 初始化表单
 const formRef = ref(null)
 const initForm = () => {
   if (formRef.value) {
    formRef.value.resetFields()
   }
   form.value = {
    policyName: '',
    policyDesc: '',
    policyConfig: [{
      "name": "",
      "kind": "",
      "timeout": "5s",
      "count": 1,
      "format": "",
      "rateLimit": 150,
      "concurrency": 150
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
        "rateLimit": 150,
        "concurrency": 150
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
  formRef.value.validate(async valid => {
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
           break;
          case 'edit':
            {
              const res = await updatePolicy(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '修改成功!'
                })
                getTableData()
                closeDialog()
              }
            }
       }
 
       initForm()
       dialogFormVisible.value = false
     }
   })
 } 
 // 新增策略
 const handleAdd = () => {
   router.push({ name: 'create'})
  //  initForm()
  //  dialogTitle.value = '新增策略'
  //  dialogType.value = 'add'
  //  dialogFormVisible.value = true
 }

 const handleEdit = (row) => {
  let id = row.ID
  router.push({ name: 'create', query: { id:id } })

  //  initForm()
  //  dialogTitle.value = '修改策略'
  //  dialogType.value = 'edit'
  //  getPolicyById(row.ID)
   
  //  dialogFormVisible.value = true
 }

 //获取单个策略修改内容
 const getPolicyById = async (id) => {
     const data = await getPolicyId({id: id})     
     form.value = data.data
 }
 
 
 const pagination = (val) => {
  page.value = val
  getTableData()
 }
 
 const handleDelete = (row) => {
   ElMessageBox.confirm('此操作将删除策略, 是否继续?', '提示', {
     confirmButtonText: '确定',
     cancelButtonText: '取消',
     type: 'warning'
   })
     .then(async() => {
       const res = await deletePolicy({ id: row.ID })
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
   getTableData()
 }
 
 const addTmpData = () => {
  if(form.value.policyConfig.length <= 3) {
    form.value.policyConfig.push({
      "name": "",
      "kind": "",
      "timeout": "5s",
      "count": 1,
      "format": "",
      "rateLimit": 150,
      "concurrency": 150
    })
  }
 }

 const delTmpData = () => {
    form.value.policyConfig.pop()
 }

 const getTypeName = (type) => {
    if(type) {
      const typeNameList = ['未知', '资产发现', '漏洞扫描', '弱口令']
      return typeNameList[type]
    }
 }

 // 配置选中扫描类型时返回模板
 const tmpOption = [[],[],[]]
 const getTemplateData = async () => {
    const table = await getTemplateList({
        page: 1,
        pageSize: 99999,
        isAll: false,
    });
    table.data.list.forEach(e => {
        if(e.templateType == 1) {
          tmpOption[0].push({label:e.templateName, value: e.ID})
        }
        if(e.templateType == 2) {
          tmpOption[1].push({label:e.templateName, value: e.ID})
        }
        if(e.templateType == 3) {
          tmpOption[2].push({label:e.templateName, value: e.ID})
        }
    });    
 }

 const changeScanType = (e)=> {
      e.templates = []
 }

 // 判断全选
const checkAll = ref(false)
const indeterminate = ref(false)
const handleCheckAll = (e, w) => {
    indeterminate.value = false
    console.log(e, w)

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
 