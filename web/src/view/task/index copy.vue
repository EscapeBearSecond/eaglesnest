<script setup>
import { ref, reactive } from 'vue' 
import { getAreaList, createArea, editArea, delArea } from "@/api/area"
import { ElMessage, ElMessageBox  } from 'element-plus'


const searchInfo = reactive({
  taskName:''
})

const tableColumns = reactive([
       { label:'名称', prop:'taskName'},
       { label:'描述', prop:'taskDesc'},
       { label:'目标', prop:'targetIp'},
       { label:'任务计划', prop:'taskPlan'},
       { label:'状态', prop:'status'},
       { label:'策略', prop:'planConfig'},
    ])
const tableData = ref([])
const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
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

const addDialogFlag = ref(false)
const date = new Date()

const formData = reactive({
  taskName:"",
  taskDesc:"",
  status:"",
  targetIp:"",
  policyId:"",
  taskPlan:"1",
  date:"",
  frequency:"",
  planConfig:{
    date: date.setTime(date.getTime() - 3600 * 1000 * 24),
    time:"",
    frequency:"",
  },
})

const taskForm = ref(null)
const checkFlag = ref(false)
const initForm = () => {
  checkFlag.value = false
  taskForm.value.resetFields()
  form.value = {
    ID: 0,
    path: '',
    name: '',
    hidden: false,
    parentId: 0,
    component: '',
    meta: {
      title: '',
      icon: '',
      defaultMenu: false,
      closeTab: false,
      keepAlive: false,
    },
  }
}

//策略Option
let policyOptions = reactive([
  { label: "测试", value: 1}
])
const editDialogFlag = ref(false);
let editData = reactive({});
const labelPosition = ref('left')
const itemLabelPosition = ref('left')

const rules = reactive({
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
});
const onCancel = () => {
  addDialogFlag.value = false;
  initForm();
}

const onSubmit = async () => {
  try {
    listQuery.page = 1
    await getTableData();
    ElMessage({
      type: 'success',
      message: '查询成功！',
      showClose: true,
    });
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '查询失败，请重试。',
      showClose: true,
    });
  }
}

const onReset = () => {
  searchInfo.areaName = "";
  onSubmit();
}

const createAsset = ()=> { 
  addDialogFlag.value = true;
  initForm();
}

function getIpArr(e) {
    if(e.includes(',')) {
        return e.split(',')
    }else {
      return [e]
    }
}

const handleDel = (row) => { 
  ElMessageBox.confirm(
    '是否删除该条数据?',
    '提示：',
    {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning',
    }
  )
    .then(() => {
      delArea({id: row.id}).then(res=> {
          if(res.code == 0) {
            ElMessage({
              type: 'success',
              message: '删除成功！',
            })
          }
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '已取消删除.',
      })
    });
};

const handleEdit = (row) => {
  try {
    if (!row || !row.id) {
      throw new Error('无效的行数据');
    }
    editData.id = row.id;
    editData.areaName = row.areaName;
    editData.areaIpStr = row.areaIP.join(',');
    editData.areaDesc = row.areaDesc;

    editDialogFlag.value = true;
  } catch (error) {
    console.error(error);
    ElMessage({
      type: 'error',
      message: '编辑失败，请重试。',
      showClose: true,
    });
  }
};

const onSubmitDialog =  async() => {
  taskForm.value.validate(async(valid) => {
    if (valid) {
      let res
      if (isEdit.value) {
        res = await updateBaseMenu(form.value)
      } else {
        res = await addBaseMenu(form.value)
      }
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: isEdit.value ? '编辑成功' : '添加成功!',
        })
        getTableData()
      }
      initForm()
      dialogFormVisible.value = false
    }
  })
}

const pagination = () => {
  getTableData();
}

const getTableData = async() => {
  try {
    const table = await getAreaList({
      page: listQuery.page,
      pageSize: listQuery.pageSize,
      ...searchInfo,
    });
    if (table.code === 0) {
      tableData.value = table.data.list;
      listQuery.total = table.data.total;
      listQuery.page = table.data.page;
      listQuery.pageSize = table.data.pageSize;
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '获取数据失败，请重试。',
      showClose: true,
    });
  }
}
getTableData();

const onEditSubmitDialog = async () => {

  let data = {
    id: editData.id,
    areaName: editData.areaName,
    areaDesc: editData.areaDesc,
    areaIp: getIpArr(editData.areaIpStr)
  };

  try {
    const res = await editArea(data);
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '修改成功！',
        showClose: true,
      });
      editDialogFlag.value = false;
      getTableData();
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '修改失败，请重试。',
      showClose: true,
    });
  }
}

const onEditCancel = ()=> {
  editDialogFlag.value = false;
  resetEditData();
}

const onDialogClose = (val) => {
  if (!val) {
    resetEditData();
  }
};

function resetEditData() {
  editData.id = "";
  editData.areaName = "";
  editData.areaIpStr = "";
  editData.areaDesc = "";
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
        <el-button type="primary" icon="plus" @click="createAsset">创建任务</el-button>
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :pagination="pagination"
        :index="true"
      >
      </advance-table>
    </div>
    <el-drawer
      v-model="addDialogFlag"
      size="40%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">扫描任务</span>
          <div>
            <el-button @click="onCancel">取 消</el-button>
            <el-button
              type="primary"
              @click="onSubmitDialog"
            >确 定</el-button>
          </div>
        </div>
      </template>
        <el-form
        :label-position="labelPosition"
        label-width="auto"
        :model="formData"
        :rules="rules"
        style="max-width: 500px"
        ref="formRef"
      >
        <el-form-item label="扫描名称：" :label-position="itemLabelPosition" prop="taskName">
          <el-input v-model="formData.taskName" placeholder="请输入扫描名称" />
        </el-form-item>
        
        <el-form-item label="扫描状态：" :label-position="itemLabelPosition" prop="status">
          <el-select v-model="formData.status" placeholder="请选择扫描状态">
            <el-option label="开启" value="1" />
            <el-option label="关闭" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item label="扫描I P：" :label-position="itemLabelPosition" prop="targetIp">
          <el-input type="textarea" :rows="6" v-model="formData.targetIp" placeholder="请输入扫描 I P , 参考：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 多个地址段请用逗号分隔"></el-input>
        </el-form-item>
        <el-form-item label="扫描策略：" :label-position="itemLabelPosition" prop="policyId">
          <el-select v-model="formData.policyId" placeholder="请选择策略模板">
            <el-option
              v-for="item in policyOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="执行方式：" :label-position="itemLabelPosition" prop="taskPlan">
          <el-select v-model="formData.taskPlan" placeholder="请选择执行方式">
            <el-option label="立即执行" value="1" />
            <el-option label="定时执行" value="0" />
          </el-select>
        </el-form-item>
        <el-form-item  v-show="formData.taskPlan != 1" label="计划时间" :label-position="itemLabelPosition" prop="date">
                <el-date-picker
                  v-model="formData.date"
                  type="datetime"
                  placeholder="选择定时计划时间">
                </el-date-picker>
              </el-form-item>
        <el-form-item v-show="formData.taskPlan != 1" label="扫描频率" :label-position="itemLabelPosition">
          <el-input v-model="formData.frequency" />
        </el-form-item>
        <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
          <el-input v-model="formData.taskDesc" />
        </el-form-item>
      </el-form>
    </el-drawer>
    <el-drawer
      v-model="editDialogFlag"
      size="40%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">扫描任务</span>
          <div>
            <el-button @click="onEditCancel">取 消</el-button>
            <el-button
              type="primary"
              @click="onEditSubmitDialog"
            >保 存</el-button>
          </div>
        </div>
      </template>
        <el-form
        :label-position="labelPosition"
        label-width="auto"
        :model="editData"
        :rules="rules"
        style="max-width: 500px"
        ref="formRef"
      >
        <!-- Form items -->
        <el-form-item label="区域名称" :label-position="itemLabelPosition" prop="areaName">
          <el-input v-model="editData.areaName" />
        </el-form-item>
        <el-form-item label="IP范围" :label-position="itemLabelPosition" prop="areaIpStr">
            <el-input type="textarea" :rows="6" v-model="editData.areaIpStr" placeholder="参考：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 多个地址段请用逗号分隔" />
        </el-form-item>
        <el-form-item label="备注" :label-position="itemLabelPosition">
          <el-input v-model="editData.areaDesc" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<style lang='scss' scoped>

</style>