<template>
  <div>
    <div class="gva-search-box">
      <el-form
        ref="searchForm"
        :inline="true"
        :model="searchInfo"
      >
        <el-form-item label="模板名称">
          <el-input
            v-model="searchInfo.templateName"
            placeholder="模板名称"
          />
        </el-form-item>
        <el-form-item label="模板类型">
          <el-select v-model="searchInfo.templateType" placeholder="请选择执行方式">
            <el-option label="全部" value=""></el-option>
            <el-option
              v-for="item in templateOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
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
          @click="addTemplate"
        >新增模板</el-button>
        <el-button
          type="primary"
          icon="plus"
          @click="addTemplateFile"
        >批量新增</el-button>
      </div>
      <advance-table
        :columns="tableColumns"
        :tableData="tableData"
        :listQuery="listQuery"
        :statusData="statusData"
        :changePageSize="changeSize"
        :pagination="handleCurrentChange"
        :index="true"
      >
      <template v-slot:custType="slotProps">
        <!-- 自定义的字段 -->
        <span v-for="(item, key) in templateOptions" :key="key" style="margin-left: 5px;"> 
          <el-tag
            type="primary"
            effect="dark"
            v-if="slotProps.row.templateType == item.value"
          >
          {{  item.label }}</el-tag>
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
    <el-drawer v-model="drawerVisible" size="45%" :show-close="false">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="text-lg">批量上传</span>
          <div>
            <el-button @click="closeUpload">取 消</el-button>
            <el-button
              type="primary"
              @click="handleSubmit"
            >确 定</el-button>
          </div>
        </div>
      </template>
      <div>
      <span class="el-form-item__label">模板类型</span> 
      <el-select
        v-model="templateType"
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
      <span class="el-form-item__label">选择文件</span> 
      <el-upload
        ref="uploadRef"
        class="upload-demo"
        drag
        :action="''"
        :auto-upload="false"
        multiple
        @change="handleFileChange"
      >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            拖动文件到这里 或<em> 点击上传</em>
          </div>
          <template #tip>
            <div class="el-upload__tip">
             注：仅可以上传Yaml文件，单个文件小于10M,数量不超过100
            </div>
          </template>
        </el-upload>
      </div>
  </el-drawer>

  </div>
</template>

<script setup>
import { getTemplateList,  createTemplate, updateTemplate, delTemplate, postTemplateImports } from "@/api/template.js"
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { getDict } from '@/utils/dictionary'

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

const changeSize = (e) => {
  listQuery.page = 1
  listQuery.pageSize = e
  getTableData()
}

let templateOptions = ref([])

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

const searchInfo = ref({})
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

const onSubmit = () => {
  listQuery.page = 1
  getTableData()
}

// 查询
const getTableData = async() => {
  const table = await getTemplateList({
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

const initPage = async() => {
  getTableData()
  getOptionDict()
}

const getOptionDict = async() => {
  const res = await getDict('templateType')
  res && res.forEach(item => {
      templateOptions.value.push({label: item.label, value: item.value})
  })
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
  templateType:"",
  templateContent:"",
})

const tableColumns = reactive([
    { label:'名称', prop:'templateName'},
    { label:'I D', prop:'templateId'},
    { label:'类型', prop:'templateType',  slot: 'custType'},
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
  ]
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
            getTableData()
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
    dialogFlag.value = 'edit'
    row.templateType = row.templateType + ''
    tempFormData.value = JSON.parse(JSON.stringify(row))
    templateDialog.value = true
}

// 批量上传
const drawerVisible = ref(false)
const uploadRef = ref(null);
const templateType = ref('1');
const addTemplateFile =  () => {
  drawerVisible.value = true
}
const selectedFiles = ref([]);
const handleFileChange = (file, fileList) => {
  selectedFiles.value = fileList.map(item => item.raw);
};

const handleSubmit = async() => {
  
  if (selectedFiles.value.length === 0) {
    ElMessage({ type: 'error', message: '未选中文件' })
    return;
  }

  const formData = new FormData();
  for(let i = 0; i < selectedFiles.value.length; i++) {
    formData.append('file', selectedFiles.value[i]);
  }
  formData.append('templateType', templateType.value)
  
  let data = await postTemplateImports(formData)
  if (data.code === 0) {
      ElMessage({ type: 'success', message: '提交成功' })
      closeUpload()
    }
    getTableData()
};

const closeUpload = () => {
  drawerVisible.value = false
  selectedFiles.value = []
} 
</script>

<style lang="scss">

.el-upload__tip {
  color: red;
}

v-deep{
  .el-upload-list {
    display: grid;
    grid-template-rows: 1fr 1fr 1fr;
  }
}
</style>
