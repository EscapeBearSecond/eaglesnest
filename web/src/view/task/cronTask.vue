<template>
  <div>
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
        <el-form-item label="扫描状态：" :label-position="itemLabelPosition" prop="status">
          <el-select v-model="taskForm.status" placeholder="请选择扫描状态">
            <el-option label="开启" value="1" />
            <el-option label="关闭" value="0" />
          </el-select>
        </el-form-item>
        <p style="margin-left:100px"><warning-bar title="注：多个地址段请用逗号分隔" /></p>
        <el-form-item label="扫描I P：" :label-position="itemLabelPosition" prop="targetIp">
          <el-input type="textarea" :rows="4" v-model="taskForm.targetIpStr" placeholder="请输入扫描IP, 例：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 "></el-input>
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
import WarningBar from '@/components/warningBar/warningBar.vue'
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
defineOptions({
  name: 'Task',
})

const page = ref(1)
const tableData = ref([])
const listQuery = reactive({
   page : 1,
   total: 0,
   pageSize: 10,
})

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
      name: "停止",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleStop(scope.row),
  },
  {
      name: "删除",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleDel(scope.row),
  },
  {
      name: "生成报告",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleReport(scope.row),
  }
])

const searchInfo = reactive({
  taskName:''
})

// 查询
const getTableData = async() => {
  const table = await getTaskList({
      page: listQuery.page,
      pageSize: listQuery.pageSize,
      isAll:true,
      taskPlan: [3],
      ...searchInfo,
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
const setPolicyOption = async() => {
    const data = await getPolicyList({ page: 1, pageSize: 99999 })
    
    policyOption.value = data.data.list.map((item)=> {
      return {label: item.policyName, value: item.ID}
    })
    console.log(data.data.list, policyOption.value)
}

const getPolicyName = (id) => {
   let item = policyOption.value.find((item) => item.value == id);   
   return item.label
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
  status:"",
  targetIp:"",
  targetIpStr:"",
  policyId:"",
  taskPlan:[3],
  date:"",
  frequency:"",
})

const tableColumns = reactive([
  { label:'名称', prop:'taskName'},
  { label:'目标', prop:'targetIp'},
  { label:'执行方式', prop:'taskPlan', slot: 'customTaskPlan'},
  { label:'状态', prop:'status', formatter(row, column) {
      let res = ['创建中','执行中','已完成', '执行失败']
      return res[row.status]
  }},
  { label:'计划配置', prop:'planConfig'},
  { label:'描述', prop:'taskDesc'},
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
      // req.areaIp = getIpArr(req.areaIpStr)
      if (dialogFlag.value === 'add') {
        req.targetIp = getIpArr(req.targetIpStr)
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

</script>

<style lang="scss">
</style>
