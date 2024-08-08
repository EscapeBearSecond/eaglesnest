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
         <el-form-item label="执行方式">
          <el-select v-model="searchInfo.taskPlan" placeholder="请选择执行方式" >
            <el-option label="立即执行" :value="1" />
            <el-option label="稍后执行" :value="2" />
          </el-select>
         </el-form-item>
         <el-form-item label="状态">
          <el-select v-model="searchInfo.status" placeholder="请选择状态">
            <el-option label="创建中" :value="0" />
            <el-option label="执行中" :value="1" />
            <el-option label="已完成" :value="2" />
            <el-option label="执行失败" :value="3" />
            <el-option label="已终止" :value="4" />
            <el-option label="运行中" :value="5" />
            <el-option label="已停止" :value="6" />
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
          <el-select v-model="taskForm.taskPlan" placeholder="请选择执行方式">
            <el-option label="立即执行" value="1" />
            <el-option label="稍后执行" value="2" />
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
  </div>
</template>

<script setup>
import {
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
defineOptions({
  name: 'Task',
})

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
const statusWidth = ref('220')
const templateOptions = reactive([
    {label: "漏洞扫描", value: '1'},
    {label: "资产发现", value: '2'},
    {label: "弱口令", value: '3'},
])

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
      name: "停止",
      type: "primary",
      icon: "edit",
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
      isAll:true,
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
const setPolicyOption = async() => {
    const data = await getPolicyList({ page: 1, pageSize: 99999 })
    policyOption.value = data.data.list.map((item)=> {
      return {label: item.policyName, value: item.ID}
    })

    const areaData = await getAreaList({ page: 1, pageSize: 99999 })
    areaOption.value = areaData.data.list.map((item)=> {
        return { label: item.areaName, value: item.areaIp.join(',') }
    })

}

// 获取策略名称
const getPolicyName = (id) => {
   let item = policyOption.value.find((item) => item.value == id);   
   return item.label
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
        getTableData()
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
    let data = reportTask({...reportData.value });
    if(data == 7) {
      ElMessage({ type: 'error', message: data.data.msg })
    }else {
      const url = window.URL.createObjectURL(new Blob([(await data).data]))
      const link = document.createElement("a");
      link.href = url;
      link.setAttribute(
        "download",
        `report_${reportData.value.entryId}.${reportData.value.format}`
      )
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      window.URL.revokeObjectURL(url);
    }
  }else {
    reportTaskDoc({entryId: reportData.value.entryId}).then(res => {
      const url = window.URL.createObjectURL(new Blob([(res).data]))
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
    })
    
  }
  reportFlag.value = false
}

const handleClose = () => {
  reportFlag.value = false
}

// 获取执行方式
const getTypeTagName = (e) => {
    let status = ['其他', '立即执行', '稍后执行','定时执行']
    return status[e]
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
  { label:'目标', prop:'targetIp'},
  { label:'执行方式', prop:'taskPlan', slot: 'customTaskPlan'},
  { label:'策略', prop:'policyName'},
  { label:'状态', prop:'status', formatter(row, column) {
      let res = ['创建中','执行中','已完成', '执行失败', '已终止', '运行中', '已停止']
      return res[row.status]
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
      req.targetIp = getIpArr(req.targetIpStr)
      
      if (dialogFlag.value === 'add') {  
        const res = await createTask(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '创建成功' })
          await getTableData()
          closeAddDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await updateTemplate(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '编辑成功' })
          await getTableData()
          closeAddDialog()
        }
      }
    }
  })
}

const templateDialog = ref(false)
const closeAddDialog = () => {
  console.log(
  '%c 🍱 CONSOLE_INFO: ',
  'font-size:20px;background-color: #ED9EC7;color:#fff;',
  form.value
  );
  form.value.resetFields()
  templateDialog.value = false
}

const dialogFlag = ref('add')

const handleClickAdd = () => {
  dialogFlag.value = 'add'
  templateDialog.value = true
}

const handleClickUpdate = (row) => {
  console.log(row)
  dialogFlag.value = 'edit'
  taskForm.value = JSON.parse(JSON.stringify(row))
  templateDialog.value = true
}

function getIpArr(e) {
    if(e.includes(',')) {
        return e.split(',')
    }else {
      return [e]
    }
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
  ElMessageBox.confirm('此操作将启动该任务, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await startTask({ id: row.ID })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '任务启动成功!'
        })
        getTableData()
      }
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消启动任务'
      })
    })
}

</script>

<style lang="scss">
.report {
  padding: 2% 5%;
}
</style>
