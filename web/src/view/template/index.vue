<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button
          type="primary"
          icon="plus"
          @click="addTemplate"
        >新增模板</el-button>
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
      v-model="templateDialog"
      size="45%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">模板</span>
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
        :model="tempFormData"
        label-width="80px"
      >
        <el-form-item
          label="名称"
          prop="templateName"
        >
          <el-input v-model="tempFormData.templateName" placeholder="模板名称" />
        </el-form-item>
        <el-form-item
          label="类型"
          prop="templateType"
        >
        <el-select
          v-model="tempFormData.templateType"
          placeholder="类型"
          size="large"
          >
          <el-option
              v-for="item in templateOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
          />
          </el-select>
        </el-form-item>
        <el-form-item
          label="内容"
          prop="templateContent"
        >
        <el-input type="textarea" :rows="18" v-model="tempFormData.templateContent" placeholder="模板内容" />
        </el-form-item>
      </el-form>
    </el-drawer>
  </div>
</template>

<script setup>
import { getTemplateList,  createTemplate, updateTemplate, delTemplate } from "@/api/template.js"



import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

defineOptions({
  name: 'Template',
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
      name: "修改",
      type: "primary",
      icon: "Edit",
      handleClick: (scope) => handleClickUpdate(scope.row), 
  },
  {
      name: "删除",
      type: "danger",
      icon: "Delete",
      handleClick: (scope) => handleClickDelete(scope.row), 
  }
])

const searchInfo = reactive({
  templateName:''
})

// 查询
const getTableData = async() => {
  const table = await getTemplateList({
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
    const res = await delTemplate({ id: row.ID })
    if (res.code === 0) {
      ElMessage.success('删除成功')
      await getTableData()
    }
  })
}

// 弹窗相关
const tempFormData = ref({
  templateName:"",
  templateType:"",
  templateContent:"",
})

const tableColumns = reactive([
    { label:'名称', prop:'templateName'},
    { label:'类型', prop:'templateType'},
])

const rules = reactive({
  templateName: [
    { required: true, message: '请输入模板名称', trigger: 'blur' },
  ],
  templateType: [
    { required: true, message: '请选择模板类型', trigger: 'blur' },
  ],
  templateContent: [
    { required: true, message: '请输入模板内容', trigger: 'blur' },
  ],
  
})
const form = ref(null)
const enterAddDialog = async() => {
  
  form.value.validate(async valid => {
    if (valid) {
      const req = {
        ...tempFormData.value
      }
      // req.areaIp = getIpArr(req.areaIpStr)
      if (dialogFlag.value === 'add') {
        const res = await createTemplate(req)
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
  form.value.resetFields()
  templateDialog.value = false
}

const dialogFlag = ref('add')

const addTemplate = () => {
  dialogFlag.value = 'add'
  templateDialog.value = true
}

const handleClickUpdate = (row) => {
  console.log(row)
  dialogFlag.value = 'edit'
  row.templateType = row.templateType
  tempFormData.value = JSON.parse(JSON.stringify(row))
  templateDialog.value = true
}

</script>

<style lang="scss">
</style>
