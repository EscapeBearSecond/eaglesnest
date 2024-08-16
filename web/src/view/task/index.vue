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
             v-model="searchInfo.taskName"
             placeholder="请输入任务名称"
           />
         </el-form-item>
         <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态">
            <el-option
              v-for="item in statusOption"
              :key="item.value"
              :label="item.label"
              :value="parseInt(item.value)"
              :disabled="item.value == 5 || item.value == 6 ? true : false"
          />
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
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="handleClickAdd"
        >新增任务</el-button>
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :pagination="handleCurrentChange"
        :changePageSize="changeSize"
        :index="true"
        :statusWidth="statusWidth"
      >
      <template v-slot:customTaskPlan="slotProps">
        <!-- 自定义的字段 -->
        <span>
          <el-tag effect="dark" >{{ getTypeTagName(slotProps.row.taskPlan) }}</el-tag>
        </span>
      </template>
      </advance-table>

    </div>
    <el-drawer
      v-model="templateDialog"
      size="45%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">普通任务</span>
          <div>
            <el-button @click="closeAddDialog">取 消</el-button>
            <el-button
              type="primary"
              @click="enterAddDialog"
            >确 定</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="form"
        :rules="rules"
        :model="taskForm"
        label-width="100px"
      >
        <el-form-item label="扫描名称：" :label-position="itemLabelPosition" prop="taskName">
          <el-input v-model="taskForm.taskName" placeholder="请输入扫描名称" />
        </el-form-item>
        <el-form-item label="执行方式：" :label-position="itemLabelPosition" prop="taskPlan">
          <el-select v-model="taskForm.taskPlan" placeholder="请选择执行方式" >
            <el-option
                v-for="item in executeTypeOption"
                :key="item.value"
                :label="item.label"
                :value="item.value"
                :disabled="item.value == 3"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="IP类型：" :label-position="itemLabelPosition" prop="targetIp">
          <el-radio-group v-model="taskForm.scanIpType">
            <el-radio-button label="默认" value="1" />
            <el-radio-button label="自定义" value="2" />
          </el-radio-group>
        </el-form-item>
        <p style="margin-left:100px" v-if="taskForm.scanIpType == 2"><warning-bar title="注：多个地址段请用逗号分隔！" /></p>
        <el-form-item label="IP地址：" v-if="taskForm.scanIpType == 2">
          <el-input  type="textarea" :rows="4" v-model="taskForm.targetIpStr" placeholder="请输入扫描IP, 例：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 "></el-input>
        </el-form-item>
        <el-form-item label="扫描区域：" v-if="taskForm.scanIpType == 1">
          <el-select  v-model="taskForm.areaIp" multiple placeholder="请选择扫描任务区域,可多选">
            <el-option  
              v-for="item in areaOption"
              :key="item.value"
              :label="item.label"
              :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="扫描策略：" :label-position="itemLabelPosition" prop="policyId">
          <el-select v-model="taskForm.policyId" placeholder="请选择策略模板">
            <el-option
              v-for="item in policyOption"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
          <el-input type="textarea" :rows="3" v-model="taskForm.taskDesc" />
        </el-form-item>
      </el-form>
    </el-drawer>
    <el-dialog
      v-model="reportFlag"
      title="导出"
      width="500"
      :before-close="handleClose"
    >
      <div class="el-form-item report">
        <span class="el-form-item__label">报告类型</span>
        <el-select v-model="reportData.type" placeholder="请选择导出类型类型">
          <el-option label="默认报告" value="1" />
          <el-option label="任务结果" value="2" />
        </el-select>
      </div>
      <div class="el-form-item report" v-if="reportData.type == 1">
        <span class="el-form-item__label">文件类型</span>
        <el-select v-model="reportData.format" placeholder="请选择导出报告类型">
          <el-option label="Word" value="docx" />
        </el-select>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="getReport">
            确定
          </el-button>
        </div>
      </template>
    </el-dialog>
    <el-dialog v-model="showDialogFlag"  title="任务" :show-close="false" width="50%">

     <el-descriptions
        title=""
        direction="vertical"
        :column="2"
        :size="showSize"
        border
      >
        <el-descriptions-item label="任务名称" align="center">
          {{  showInfo.taskName }}
        </el-descriptions-item>
        <el-descriptions-item label="执行方式" align="center">{{  showInfo.taskPlan == 1 ? '立即执行' : '稍后执行' }}</el-descriptions-item>
        <el-descriptions-item label="关联策略" align="center">{{  showInfo.policyName }}</el-descriptions-item>
        <el-descriptions-item label="当前状态" align="center">
          <el-tag size="small">{{  getStatus(showInfo.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="任务数量" v-if="showInfo.status == 1" align="center">
            {{  stageData.total }}
        </el-descriptions-item>
        <el-descriptions-item label="执行序号" v-if="showInfo.status == 1" align="center">
            {{  stageData.running }}
        </el-descriptions-item>
        <el-descriptions-item label="正在执行" v-if="showInfo.status == 1" align="center">
            {{  stageData.name }}
        </el-descriptions-item>
        <el-descriptions-item label="当前进度" v-if="showInfo.status == 1" align="center">
          <el-progress type="dashboard" :percentage="stageData.percent * 100" :color="colors" />
        </el-descriptions-item>
        <el-descriptions-item label="扫描 IP" :span="2" align="center">
          <div class="ip-content"> 
            {{  IpToStr(showInfo.targetIp) }}
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="任务描述" :span="2" align="center">
          {{  showInfo.taskDesc }}
        </el-descriptions-item>
      </el-descriptions>
      <div class="close-btn"><el-button type="primary" @click="showDialogFlag = false">关闭</el-button></div>
  </el-dialog>
  </div>
</template>

<script setup>
import {
  getTask,
  getTaskStage,
  getTaskList,
  createTask,
  stopTask,
  delTask,
  startTask,
  reportTask,
  reportTaskDoc
} from '@/api/task.js'
import { getPolicyList } from '@/api/policy.js'
import { getAreaList } from '@/api/area.js'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDict } from '@/utils/dictionary'

defineOptions({
  name: 'Task',
})
const itemLabelPosition = ref('right')
const searchInfo = ref({
  taskName: '',
})
const onSubmit = () => {
  listQuery.page = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}


const page = ref(1)
const tableData = ref([])
const listQuery = reactive({
   page : 1,
   total: 0,
   pageSize: 10,
})
const statusWidth = ref('280')
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 操作
const statusData = reactive([
  {
      name: "启动",
      type: "primary",
      icon: "SwitchButton",
      handleClick: (scope) => handleStart(scope.row),
      visible : (scope) => visibleStart(scope.row)
  },
  {
      name: "查看",
      type: "primary",
      icon: "View",
      handleClick: (scope) => handleShow(scope.row),
      // visible : (scope) => visibleStart(scope.row)
  },
  {
      name: "停止",
      type: "primary",
      icon: "VideoPause",
      handleClick: (scope) => handleStop(scope.row),
      visible : (scope) => visibleStop(scope.row)
  },
  {
      name: "删除",
      type: "primary",
      icon: "Delete",
      handleClick: (scope) => handleDel(scope.row),
  },
  {
      name: "导出",
      type: "primary",
      icon: "Download",
      handleClick: (scope) => handleReport(scope.row),
      visible : (scope) => visibleReport(scope.row)
  },
])

// 查询
const getTableData = async() => {  
  const table = await getTaskList({
      page: listQuery.page,
      pageSize: listQuery.pageSize,
      isAll: true,
      ...searchInfo.value,
    });
    if (table.code === 0) {
      tableData.value = table.data.list;
      listQuery.total = table.data.total;
      listQuery.page = table.data.page;
      listQuery.pageSize = table.data.pageSize;
    }
}

// 获取策略模板
const policyOption = ref([])
const areaOption = ref([])
const statusOption = ref([])
const executeTypeOption = ref([])
const setPolicyOption = async() => {
    const data = await getPolicyList({ page: 1, pageSize: 99999 })
    policyOption.value = data.data.list.map((item)=> {
      return {label: item.policyName, value: item.ID}
    })

    const areaData = await getAreaList({ page: 1, pageSize: 99999 })
    areaOption.value = areaData.data.list.map((item)=> {
        return { label: item.areaName, value: item.areaIp.join(',') }
    })
    
    const res = await getDict('taskStatus')
    res && res.forEach(item => {
      statusOption.value.push({label: item.label, value: item.value})
    })

    const executeTypeData = await getDict('executeType')
    executeTypeData && executeTypeData.forEach(item => {
      executeTypeOption.value.push({label: item.label, value: item.value})
    })
}

// 初始化
const initPage = async() => {
  setPolicyOption()
  getTableData()
}

initPage()

// 停止
const handleStop = (row) => {
  ElMessageBox.confirm('此操作将停止该任务, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await stopTask({ id: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '任务停止成功!'
        })
        setTimeout(()=> {
          getTableData()
        }, 1000)
      }else {
        ElMessage({
          type: 'error',
          message: res.msg
        })
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消停止'
      })
    })
}

// 删除
const handleDel = (row) => {
  ElMessageBox.confirm('此操作将永久删除该任务, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await delTask({ id: row.id })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功!'
        })
        if (tableData.value.length === 1 && page.value > 1) {
          page.value--
        }
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

// 下载报告
const reportFlag = ref(false)
const reportData = ref({})
const handleReport =  async(row) =>{
  reportFlag.value = true
  reportData.value.entryId = row.entryId
}

const getReport = async() => {
  let date = new Date();
  let timestamp = date.getTime();
  if(reportData.value.type == 1) {
    reportTask({...reportData.value }).then(res => {
        const blob = res.data;
        let resData =  new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = () => {
            try {
              // 如果是个json对象
              const json = JSON.parse(reader.result);
              reject(json);
            } catch (e) {
              // 如果是 blob 
              resolve(blob);
            }
          };
          reader.onerror = () => {
            reject(new Error('未读取到文件对象'));
          };
          reader.readAsText(blob);
      })
      resData.then(blob => {
        // 创建下载链接并触发下载
        const url = window.URL.createObjectURL(new Blob([blob]))
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute(
          "download",
          `report_${timestamp}.docx`
        )
        document.body.appendChild(link);
        link.click();
        setTimeout(() => {
          document.body.removeChild(link);
          window.URL.revokeObjectURL(url); 
        }, 250);
      }).catch(error => {
          if (typeof error === 'object' && error !== null) {
            ElMessage({
              type: 'error',
              message: `${error.msg}`
            })
          } else {
            ElMessage({
              type: 'error',
              message: '下载文档时发生了错误！'
            })
          }
        });
    })
  }else {
    reportTaskDoc({entryId: reportData.value.entryId}).then(res => {
        const blob = res.data;
        let resData =  new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.onload = () => {
          try {
            // 如果是个json对象
            const json = JSON.parse(reader.result);
            reject(json);
          } catch (e) {
            // 如果是 blob 
            resolve(blob);
          }
        };
        reader.onerror = () => {
          reject(new Error('未读取到文件对象'));
        };
        reader.readAsText(blob); 
      })
      
      resData.then(blob => {
        // 创建下载链接并触发下载
        const url = window.URL.createObjectURL(new Blob([blob]))
        const link = document.createElement("a");
        link.href = url;
        link.setAttribute(
          "download",
          `report_${timestamp}.zip`
        )
        document.body.appendChild(link);
        link.click();
        setTimeout(() => {
          document.body.removeChild(link);
          window.URL.revokeObjectURL(url); 
        }, 250);
      }).catch(error => {
          if (typeof error === 'object' && error !== null) {
            ElMessage({
              type: 'error',
              message: `${error.msg}`
            })
          } else {
            ElMessage({
              type: 'error',
              message: '下载文档时发生了错误！'
            })
          }
        });
      }).catch(error => {
        ElMessage({
          type:'error',
          message: '下载文档时发生了错误!'
        })
      });    
  }
  reportFlag.value = false
}

const handleClose = () => {
  reportFlag.value = false
}

// 获取执行方式
const getTypeTagName = (e) => {
    let opt = executeTypeOption.value.find(item => item.value == e)
      if(!opt) {
        return ''
      }
      return opt.label
}

// 表单
/**
 * {
    taskName:"",
    taskDesc:"",
    targetIp:[],
    targetIpStr:"",
    policyId:"",
    date:"",
    frequency:"",
    scanIpType:'1',
    areaIp: [],
  }
 */
 const taskForm = ref({scanIpType:'1'})
// 表头
const tableColumns = reactive([
  { label:'名称', prop:'taskName'},
  { label:'执行方式', prop:'taskPlan', slot: 'customTaskPlan'},
  { label:'策略', prop:'policyName'},
  { label:'状态', prop:'status', formatter(row, column) {
      let opt = statusOption.value.find(item => item.value == row.status)
      if(!opt) {
        return ''
      }
      return opt.label
  }},
])

//验证输入
const rules = reactive({
  taskName: [
    { required: true, message: '请输入扫描名称', trigger: 'blur' }
  ],
  targetIpStr: [
    { required: true, message: '请输入扫描IP', trigger: 'blur' }
  ],
  areaIp: [
    { required: true, message: '请选择扫描区域', trigger: 'blur' }
  ],
  status: [
    { required: true, message: '请选择扫描状态', trigger: 'blur' }
  ],
  taskPlan: [
  { required: true, message: '请选择扫执行方式', trigger: 'blur' }
  ],
  policyId: [
    { required: true, message: '请选择策略模板', trigger: 'blur' }
  ]
})
// 提交表单
const form = ref(null)
const enterAddDialog = async() => {
  form.value.validate(async valid => {
    if (valid) {
      const req = {
        ...taskForm.value
      }
      // 这里加了判断 是否是默认执行方式，如果是默认 就是 区域选择 如果是自定义就是输入内容
      req.scanIpType != 1 ? (req.targetIp = getIpArr(req.targetIpStr)): req.targetIp = req.areaIp;
      console.log(req);
      if (dialogFlag.value === 'add') {  
        const res = await createTask(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '创建成功' })
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await updateTemplate(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '编辑成功' })
        }
      }

      await getTableData()
      closeAddDialog()
    }
  })
}

const templateDialog = ref(false)
const closeAddDialog = () => {
  form.value.resetFields()
  templateDialog.value = false
}

const dialogFlag = ref('add')

const handleClickAdd = () => {
  dialogFlag.value = 'add'
  templateDialog.value = true
}

function getIpArr(e) {
    console.log(e)
    if(e.includes(',')) {
        return e.split(',')
    }else {
      return [e]
    }
}

function IpToStr(e) {   
   if(Array.isArray(e)) {
      if(e.length > 0) {
        let result = [];
        for (let i = 0; i < e.length; i += 2) {
            result.push(e.slice(i, i + 2).join(','));
        }
        return result.join('\n');
      }else {
        return e[0]
      }
   }else {
    return '';
   }
}

function getStatus (e) {
  let opt = statusOption.value.find(item => item.value == e)
  if (!opt) {
      return ''
  }
  return opt.label
}

// 根据状态来判断是否显示停止按钮
const visibleStop = (e) => {
    return e.status == 1
}

// 根据状态来判断是否显示报告按钮
const visibleReport = (e) => {
    return e.status == 2
}

//
const visibleStart = (e) => {
    return e.status != 1
}

const handleStart = (e) => {  
  ElMessageBox.confirm(
    '确定要启动任务吗?',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async() => {
      const res = await startTask({ id: e.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '任务启动成功!'
        })
        getTableData()
      }
    })
}

const changeSize = (e) => {
  listQuery.page = 1
  listQuery.pageSize = e
  getTableData()
}

const showDialogFlag = ref(false)
const showSize = ref('default')
const showInfo = ref({})
let stageData = ref({
  name:'',  
  percent:'',  
  total:'',  
  running:'',  
})
const colors = ref([
  { color: '#f56c6c', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#5cb87a', percentage: 60 },
  { color: '#1989fa', percentage: 80 },
  { color: '#6f7ad3', percentage: 100 },
]) 

const handleShow = async(e)=> {
    let data = await getTask({id: e.ID})
    console.log(data.data.status)
    if(data.data.status == 1) {
        let data = await getTaskStage({id:e.ID})
        stageData = data.data
    }
    showInfo.value = data.data
    showDialogFlag.value = true    
}
</script>

<style lang="scss">
.report {
  padding: 2% 5%;
}

.my-header {
  display: grid;
  grid-template-columns: auto  50px;
  align-items: center;
}

.close-btn {
  display: grid;
  grid-template-rows: 1fr;
  align-items: center;
  justify-items: center;
  margin: 10px 0px;
}

.ip-content {
  text-align: left;
  white-space: normal;
}
::v-deep .el-descriptions__cell , .el-descriptions__label{
    width: 180px;
}
</style>
