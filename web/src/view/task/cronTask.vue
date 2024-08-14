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
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="handleClickAdd"
        >定时任务</el-button>
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :pagination="handleCurrentChange"
        :changePageSize="changeSize"
        :index="true"
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
          <span class="text-lg">计划任务</span>
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
        <el-form-item label="计划配置：" :label-position="itemLabelPosition" prop="planConfig">
          <el-input v-model="taskForm.planConfig" placeholder="请输入Cron表达式，例每天中午12点执行：0 0 12 * * ? " />
        </el-form-item>
        <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
          <el-input type="textarea" :rows="3" v-model="taskForm.taskDesc" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import {
  getTaskList,
  createTask,
  stopTask,
  delTask,
  reportTask,
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

const page = ref(1)
const tableData = ref([])
const itemLabelPosition = ref('right')
const listQuery = reactive({
   page : 1,
   total: 0,
   pageSize: 10,
})

const changeSize = (e) => {
  listQuery.pageSize = e
  getTableData()
}

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
      name: "生成报告",
      type: "primary",
      icon: "Position",
      handleClick: (scope) => handleReport(scope.row),
      visible : (scope) => visibleReport(scope.row)
  }
])

const searchInfo = ref({
  taskName: '',
  status: null,
})
const onSubmit = () => {
  listQuery.page = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}


// 查询
const getTableData = async() => {
  const table = await getTaskList({
      page: listQuery.page,
      pageSize: listQuery.pageSize,
      isAll:true,
      taskPlan: [3],
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
}

const initPage = async() => {
  setPolicyOption()
  getTableData()
}

initPage()

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

const handleDel = (row) => {
  ElMessageBox.confirm('此操作将永久删除该任务, 是否继续?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async() => {
      const res = await delTask({ id: row.ID })
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

const handleReport =  async(row) =>{
    await reportTask({ id: row.ID })
}

const getTypeTagName = (e) => {
    let status = ['其他', '立即执行', '稍后执行','定时执行']
    return status[e]
}

// 弹窗相关
const taskForm = ref({
  taskName:"",
  taskDesc:"",
  targetIp:"",
  targetIpStr:"",
  policyId:"",
  taskPlan:"3",
  scanIpType: '1',
})

const tableColumns = reactive([
  { label:'名称', prop:'taskName'},
  { label:'目标', prop:'targetIp'},
  { label:'执行方式', prop:'taskPlan', slot: 'customTaskPlan'},
  { label:'状态', prop:'status', formatter(row, column) {
      let opt = statusOption.value.find(item => item.value == row.status)
      if (!opt) {
          return ''
      }
      return opt.label
  }},
  { label:'计划配置', prop:'planConfig'},
])

const rules = reactive({
  taskName: [
    { required: true, message: '请输入扫描名称', trigger: 'blur' }
  ],
  targetIpStr: [
    { required: true, message: '请输入扫描IP', trigger: 'blur' }
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
const form = ref(null)
const enterAddDialog = async() => {
  
  form.value.validate(async valid => {
    if (valid) {
      const req = {
        ...taskForm.value
      }
      req.scanIpType != 1 ? (req.targetIp = getIpArr(req.targetIpStr)): req.targetIp = req.areaIp;
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

// 是否显示启动按钮
const visibleStart = (e) => {
    return e.status == 6
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
</style>
