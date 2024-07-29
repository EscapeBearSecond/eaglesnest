<template>
  <div>
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <!-- <el-button
          type="primary"
          icon="plus"
          @click="addForm"
        >新增资产</el-button> -->
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :pagination="handleCurrentChange"
        :index="true"
      >
      </advance-table>

    </div>
    <el-drawer
      v-model="districtInfoDialog"
      size="45%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">区域</span>
          <div>
            <el-button @click="closeAddUserDialog">取 消</el-button>
            <el-button
              type="primary"
              @click="enterAddUserDialog"
            >确 定</el-button>
          </div>
        </div>
      </template>

      <el-form
        ref="form"
        :rules="rules"
        :model="districtInfo"
        label-width="80px"
      >
      <el-form-item label="资产名称" :label-position="itemLabelPosition" prop="assetName">
          <el-input v-model="form.assetName" />
        </el-form-item>
        <el-form-item label="资产IP/范围" :label-position="itemLabelPosition" prop="assetIP">
          <el-input v-model="form.assetIP" />
        </el-form-item>
        <el-form-item label="所属区域" :label-position="itemLabelPosition" prop="assetArea">
          <el-input v-model="form.assetArea" />
        </el-form-item>
        <el-form-item label="资产类型" :label-position="itemLabelPosition" prop="assetType">
          <el-input v-model="form.assetType" />
        </el-form-item>
        <el-form-item label="资产价值" :label-position="itemLabelPosition">
          <el-input v-model="form.manufacturer" />
        </el-form-item>
        <el-form-item label="备注" :label-position="itemLabelPosition">
          <el-input v-model="form.desc" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getListApi, createApi, delApi } from "@/api/assets"
import { getAreaList } from "@/api/area"

import { getAuthorityList } from '@/api/authority'


import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'Area',
})
// 初始化相关
const setAuthorityOptions = (AuthorityData, optionsData) => {
  AuthorityData &&
        AuthorityData.forEach(item => {
          if (item.children && item.children.length) {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName,
              children: []
            }
            setAuthorityOptions(item.children, option.children)
            optionsData.push(option)
          } else {
            const option = {
              authorityId: item.authorityId,
              authorityName: item.authorityName
            }
            optionsData.push(option)
          }
        })
}

const page = ref(1)
const tableData = ref([])
const listQuery = reactive({
   page : 1,
   total: 0,
   pageSize: 10,
})

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 操作
const statusData = reactive([
  // {
  //     name: "删除",
  //     type: "danger",
  //     icon: "Delete",
  //     handleClick: (scope) => handleClickDelete(scope.row), 
  // }
])

const searchInfo = reactive({
    areaName:''
})

// 查询
const getTableData = async() => {
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
}

const initPage = async() => {
  getTableData()
}

initPage()

const handleClickDelete = async(row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await delApi({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

// 弹窗相关
const districtInfo = ref({
  areaName:"",
  areaIpStr:"",
  areaDesc:"",
})

const tableColumns = reactive([
    { label:'名称', prop:'assetName'},
    { label:'IP', prop:'assetIP'},
    { label:'归属', prop:'assetArea'},
    { label:'类型', prop:'assetType'},
    { label:'价值', prop:'manufacturer'},
    { label:'备注', prop:'desc'},
])

const rules = ref({
  assetName: [
    { required: true, message: '请输入区域名称', trigger: 'blur' },
  ],
  assetIp: [
    { required: true, message: '请输入IP范围', trigger: 'blur' },
  ],
  
})
const form = ref(null)
const enterAddUserDialog = async() => {
  
  form.value.validate(async valid => {
    if (valid) {
      const req = {
        ...districtInfo.value
      }
      req.areaIp = getIpArr(req.areaIpStr)
      if (dialogFlag.value === 'add') {
        const res = await createApi(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '创建成功' })
          await getTableData()
          closeAddUserDialog()
        }
      }
      if (dialogFlag.value === 'edit') {
        const res = await editArea(req)
        if (res.code === 0) {
          ElMessage({ type: 'success', message: '编辑成功' })
          await getTableData()
          closeAddUserDialog()
        }
      }
    }
  })
}

function getIpArr(e) {
    if(e.includes(',')) {
        return e.split(',')
    }else {
      return [e]
    }
}

const districtInfoDialog = ref(false)
const closeAddUserDialog = () => {
  form.value.resetFields()
  districtInfo.value.headerImg = ''
  districtInfo.value.authorityIds = []
  districtInfoDialog.value = false
}

const dialogFlag = ref('add')

const addForm = () => {
  dialogFlag.value = 'add'
  districtInfoDialog.value = true
}

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
</style>
