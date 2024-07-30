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
            v-model="searchInfo.taskName"
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
          @click="addAuthority(0)"
        >创建任务</el-button>
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :pagination="pagination"
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
    <!-- 新增角色弹窗 -->
    <el-dialog
      v-model="dialogFormVisible"
      :title="dialogTitle"
    >
      <el-form
        ref="authorityForm"
        :model="form"
        :rules="rules"
        style="padding:10px 20px;"
        label-width="100px"
      >
      <el-form-item label="扫描名称：" :label-position="itemLabelPosition" prop="taskName">
          <el-input v-model="form.taskName" placeholder="请输入扫描名称" />
        </el-form-item>
        
        <el-form-item label="扫描状态：" :label-position="itemLabelPosition" prop="status">
          <el-select v-model="form.status" placeholder="请选择扫描状态">
            <el-option label="开启" value="1" />
            <el-option label="关闭" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="扫描I P：" :label-position="itemLabelPosition" prop="targetIp">
          <el-input type="textarea" :rows="3" v-model="form.targetIp" placeholder="请输入扫描 I P , 参考：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 多个地址段请用逗号分隔"></el-input>
        </el-form-item>
        <el-form-item label="扫描策略：" :label-position="itemLabelPosition" prop="policyId">
          <el-select v-model="form.policyId" placeholder="请选择策略模板">
            <el-option
              v-for="item in policyOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="执行方式：" :label-position="itemLabelPosition" prop="taskPlan">
          <el-select v-model="form.taskPlan" placeholder="请选择执行方式">
            <el-option label="立即执行" value="1" />
            <el-option label="定时执行" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item  v-if="form.taskPlan == 0" label="计划时间" :label-position="itemLabelPosition" prop="date">
                <el-date-picker
                  v-model="form.date"
                  type="datetime"
                  placeholder="选择定时计划时间">
                </el-date-picker>
              </el-form-item>
        <el-form-item v-if="form.taskPlan == 0" label="扫描频率" :label-position="itemLabelPosition">
          <el-input v-model="form.frequency" />
        </el-form-item>
        <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
          <el-input type="textarea" :rows="3" v-model="form.taskDesc" />
        </el-form-item>
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
    </el-dialog>
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
import { ref } from 'vue'
import { ElMessage, ElMessageBox, formatter } from 'element-plus'

defineOptions({
  name: 'Authority'
})

const dialogType = ref('add')

const dialogTitle = ref('新增任务')
const dialogFormVisible = ref(false)
const apiDialogFlag = ref(false)
const form = ref({
  taskName:"",
  taskDesc:"",
  status:"",
  targetIp:"",
  policyId:"",
  taskPlan:"1",
  date:"",
  frequency:"",
  planConfig:{
    date: "",
    time:"",
    frequency:"",
  },
})

const tableColumns = ref([
    { label:'名称', prop:'taskName'},
    { label:'描述', prop:'taskDesc'},
    { label:'目标', prop:'targetIp'},
    { label:'执行方式', prop:'taskPlan', slot: 'customTaskPlan'},
    { label:'状态', prop:'status', formatter(row, column) {
       let res = ['创建中','执行中','已完成', '执行失败']
       return res[row.status]
    }},
    { label:'策略', prop:'planConfig'},
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

const page = ref(1)
const listQuery = ref({
    page: 1,
    total: 0,
    pageSize: 10
})
const statusData = ref([
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

const tableData = ref([])
const searchInfo = ref({})

// 查询
const getTableData = async() => {
  const table = await getTaskList({ page: listQuery.value.page, pageSize: listQuery.value.pageSize, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    listQuery.value.total = table.data.total
    listQuery.value.page= table.data.page
    listQuery.value.pageSize = table.data.pageSize
  }
}
const onSubmit = () => {
  listQuery.value.page = 1
  getTableData()
}
const onReset = () => {
  searchInfo.value = {}
  getTableData
}

getTableData()
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
// 初始化表单
const authorityForm = ref(null)
const initForm = () => {
  if (authorityForm.value) {
    authorityForm.value.resetFields()
  }
  form.value = {
    authorityId: 0,
    authorityName: '',
    parentId: 0
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
      form.value.authorityId = Number(form.value.authorityId)
      switch (dialogType.value) {
        case 'add':
          {
            const res = await createTask(form.value)
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
const addAuthority = (parentId) => {
  initForm()
  dialogTitle.value = '创建任务'
  dialogType.value = 'add'
  form.value.parentId = parentId
  setOptions()
  dialogFormVisible.value = true
}


const pagination = () => {
  getTableData()
}

const handleStop = (row) => {
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

const handleReport =  async(row) =>{
    await reportTask({ id: row.id })
}

const getTypeTagName = (e) => {
    let status = ['其他', '立即执行', '稍后执行','定时执行']
    return status[e]
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
