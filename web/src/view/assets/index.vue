<template>
  <div>
    <warning-bar title="注：右上角头像下拉可切换角色" />
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="addUser"
        >新增区域</el-button>
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
      v-model="addUserDialog"
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
        <el-form-item
          label="区域名称"
          prop="areaName"
        >
          <el-input v-model="districtInfo.areaName" />
        </el-form-item>
        <el-form-item
          label="IP范围"
          prop="areaIpStr"
        >
        <el-input type="textarea" :rows="6"  v-model="districtInfo.areaIpStr" placeholder="参考：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 多个地址段请用逗号分隔"  />
        </el-form-item>
        <el-form-item
          label="备注"
        >
        <el-input v-model="districtInfo.areaDesc" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>

import {
  getUserList,
  setUserAuthorities,
  register,
  deleteUser
} from '@/api/user'
import { getAreaList, createArea, editArea, delArea } from "@/api/area"

import { getAuthorityList } from '@/api/authority'
import WarningBar from '@/components/warningBar/warningBar.vue'
import { setUserInfo, resetPassword } from '@/api/user.js'

import { nextTick, ref, watch, reactive } from 'vue'
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
const total = ref(0)
const pageSize = ref(10)
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
   {
      name: "修改",
      type: "primary",
      icon: "Edit",
      handleClick: (scope) => openEdit(scope.row), 
  },
  {
      name: "删除",
      type: "danger",
      icon: "Delete",
      handleClick: (scope) => deleteUserFunc(scope.row), 
  }
])

const searchInfo = reactive({
    areaName:''
})

// 查询
const getTableData = async() => {
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
}

watch(() => tableData.value, () => {
  setAuthorityIds()
})

const initPage = async() => {
  getTableData()
  const res = await getAuthorityList({ page: 1, pageSize: 999 })
  setOptions(res.data.list)
}

initPage()

const resetPasswordFunc = (row) => {
  ElMessageBox.confirm(
    '是否将此用户密码重置为123456?',
    '警告',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(async() => {
    const res = await resetPassword({
      ID: row.ID,
    })
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: res.msg,
      })
    } else {
      ElMessage({
        type: 'error',
        message: res.msg,
      })
    }
  })
}
const setAuthorityIds = () => {
  tableData.value && tableData.value.forEach((user) => {
    user.authorityIds = user.authorities && user.authorities.map(i => {
      return i.authorityId
    })
  })
}

const chooseImg = ref(null)
const openHeaderChange = () => {
  chooseImg.value.open()
}

const authOptions = ref([])
const setOptions = (authData) => {
  authOptions.value = []
  setAuthorityOptions(authData, authOptions.value)
}

const deleteUserFunc = async(row) => {
  ElMessageBox.confirm('确定要删除吗?', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async() => {
    const res = await deleteUser({ id: row.ID })
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
    { label:'区域名称', prop:'areaName'},
    { label:'IP范围', prop:'areaIP'},
    { label:'备注', prop:'areaDesc'},
])

const rules = ref({
  areaName: [
    { required: true, message: '请输入区域名称', trigger: 'blur' },
  ],
  areaIpStr: [
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
        const res = await createArea(req)
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

const addUserDialog = ref(false)
const closeAddUserDialog = () => {
  form.value.resetFields()
  districtInfo.value.headerImg = ''
  districtInfo.value.authorityIds = []
  addUserDialog.value = false
}

const dialogFlag = ref('add')

const addUser = () => {
  dialogFlag.value = 'add'
  addUserDialog.value = true
}

const tempAuth = {}
const changeAuthority = async(row, flag, removeAuth) => {
  if (flag) {
    if (!removeAuth) {
      tempAuth[row.ID] = [...row.authorityIds]
    }
    return
  }
  await nextTick()
  const res = await setUserAuthorities({
    ID: row.ID,
    authorityIds: row.authorityIds
  })
  if (res.code === 0) {
    ElMessage({ type: 'success', message: '角色设置成功' })
  } else {
    if (!removeAuth) {
      row.authorityIds = [...tempAuth[row.ID]]
      delete tempAuth[row.ID]
    } else {
      row.authorityIds = [removeAuth, ...row.authorityIds]
    }
  }
}

const openEdit = (row) => {
  row.areaIpStr = row.areaIP.join(',');
  dialogFlag.value = 'edit'
  districtInfo.value = JSON.parse(JSON.stringify(row))
  addUserDialog.value = true
}

const switchEnable = async(row) => {
  districtInfo.value = JSON.parse(JSON.stringify(row))
  await nextTick()
  const req = {
    ...districtInfo.value
  }
  const res = await setUserInfo(req)
  if (res.code === 0) {
    ElMessage({ type: 'success', message: `${req.enable === 2 ? '禁用' : '启用'}成功` })
    await getTableData()
    districtInfo.value.headerImg = ''
    districtInfo.value.authorityIds = []
  }
}

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
</style>
