<script setup>
import { ref, reactive } from 'vue' 
import { getListApi, createApi, delApi } from "@/api/assets"
import { ElMessage, ElMessageBox  } from 'element-plus'


const searchInfo = reactive({
  assetName:''
})

const tableColumns = reactive([
      { label:'名称', prop:'assetName'},
       { label:'IP', prop:'assetIP'},
       { label:'归属', prop:'assetArea'},
       { label:'类型', prop:'assetType'},
       { label:'价值', prop:'manufacturer'},
       { label:'备注', prop:'desc'},
    ])
const tableData = ref([])
const listQuery = reactive({
      total: 0,
      page: 1,
      pageSize: 10
   })
const statusData = reactive([
  {
      name: "删除",
      type: "primary",
      icon: "edit",
      handleClick: (scope) => handleDel(scope.row), 
  }
])

const addDialogFlag = ref(false)
const formData = reactive({
  assetName:"",
  assetIP:"",
  assetArea:"",
  assetType:"",
  manufacturer:"",
  desc:"",
})
const labelPosition = ref('left')
const itemLabelPosition = ref('left')

const rules = reactive({
  assetName: [
    { required: true, message: '请输入资产名称', trigger: 'blur' }
  ],
  assetIp: [
    { required: true, message: '请输入资产IP', trigger: 'blur' }
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
      delApi({id: row.id}).then(res=> {
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

const onSubmitDialog = async (formValues) => {
  let data = {
    areaName: formValues.areaName,
    areaDesc: formValues.areaDesc,
    areaIp: getIpArr(formValues.areaIpStr)
  };
  
  try {
    const res = await createApi(data);
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
    const table = await getListApi({
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

function resetFormData() {
  formData.areaName = "";
  formData.areaIpStr = "";
  formData.areaDesc = "";
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
        <el-button type="primary" icon="plus" @click="createAsset">新增资产</el-button>
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
        <el-form-item label="资产名称" :label-position="itemLabelPosition" prop="assetName">
          <el-input v-model="localForm.assetName" />
        </el-form-item>
        <el-form-item label="资产IP/范围" :label-position="itemLabelPosition" prop="assetIP">
          <el-input v-model="localForm.assetIP" />
        </el-form-item>
        <el-form-item label="所属区域" :label-position="itemLabelPosition" prop="assetArea">
          <el-input v-model="localForm.assetArea" />
        </el-form-item>
        <el-form-item label="资产类型" :label-position="itemLabelPosition" prop="assetType">
          <el-input v-model="localForm.assetType" />
        </el-form-item>
        <el-form-item label="资产价值" :label-position="itemLabelPosition">
          <el-input v-model="localForm.manufacturer" />
        </el-form-item>
        <el-form-item label="备注" :label-position="itemLabelPosition">
          <el-input v-model="localForm.desc" />
        </el-form-item>
      </el-form>
    </el-drawer>
    
  </div>
</template>

<style lang='scss' scoped>

</style>