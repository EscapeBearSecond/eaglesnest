<script setup>
import { ref, reactive } from 'vue' 
import districtForm from "./components/districtForm.vue"
import { getAreaList, createArea, editArea, delArea } from "@/api/area"
import { ElMessage, ElMessageBox  } from 'element-plus'


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
      icon: "Edit",
      handleClick: (scope) => handleEdit(scope.row), 
  },
  {
      name: "删除",
      type: "danger",
      icon: "Delete",
      handleClick: (scope) => handleDel(scope.row), 
  }
])

const addDialogFlag = ref(false)
const dialogTitle = ref('新增区域')
const formData = reactive({
  areaName:"",
  areaIpStr:"",
  areaDesc:"",
})
const editDialogFlag = ref(false);
let editData = reactive({});
const labelPosition = ref('left')
const itemLabelPosition = ref('top')

const rules = reactive({
  areaName: [
    { required: true, message: '请输入区域名称', trigger: 'blur' }
  ],
  areaIpStr: [
    { required: true, message: '请输入区域IP范围', trigger: 'blur' }
  ]
});
const onCancel = () => {
  addDialogFlag.value = false;
  resetFormData();
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
  dialogTitle.value = '新增区域';
  resetFormData();
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

const onSubmitDialog = async (formValues) => {
  let data = {
    areaName: formValues.areaName,
    areaDesc: formValues.areaDesc,
    areaIp: getIpArr(formValues.areaIpStr)
  };
  
  try {
    const res = await createArea(data);
    if(res.code == 0) {
      ElMessage({
        type: 'success',
        message: '新增成功！',
        showClose: true,
      });
      addDialogFlag.value = false;
      getTableData();
    }
  } catch (error) {
    ElMessage({
      type: 'error',
      message: '新增失败，请重试。',
      showClose: true,
    });
  }
};

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

function resetFormData() {
  formData.areaName = "";
  formData.areaIpStr = "";
  formData.areaDesc = "";
}

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
        <el-button type="primary" icon="plus" @click="createAsset">新增区域</el-button>
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
          <span class="text-lg">区域</span>
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
    <el-drawer
      v-model="editDialogFlag"
      size="40%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">区域</span>
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