<template>
    <div>
      <div class="gva-table-box">
        <div class="gva-search-box">
         <el-form
           ref="searchForm"
           :inline="true"
           :model="searchInfo"
         >
           <el-form-item label="名称">
             <el-input
               v-model="searchInfo.name"
               placeholder="请输入名称"
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
        <advance-table
          :columns="tableColumns"
          :tableData="tableData"
          :listQuery="listQuery"
          :statusData="statusData"
          :changePageSize="changeSize"
          :pagination="handleCurrentChange"
          :index="true"
        >
            <template v-slot:customSeverity="slotProps">
            <!-- 自定义的字段 -->
                <span>
                    <el-tag  effect="plain" :color="getColor(slotProps.row.severity)">{{ getSeverityName(slotProps.row.severity) }}</el-tag>
                </span>
            </template>
        </advance-table>
      </div>
      <el-dialog v-model="showFlag" title="漏洞" width="800">
        <el-descriptions
          title=""
          direction="vertical"
          :column="3"
          :size="size"
          border
        >
          <el-descriptions-item label="编号" class="one-third-width" align="center">{{ showData.classification.cve }}</el-descriptions-item>
          <el-descriptions-item label="名称" class="one-third-width" align="center">{{ showData.name }}</el-descriptions-item>
          <el-descriptions-item label="等级" class="one-third-width" align="center">
            <el-tag effect="dark" :color="getColor(showData.severity)">{{ getSeverityName(showData.severity) }}</el-tag>
          </el-descriptions-item>

          <!-- 第二行 -->
          <el-descriptions-item label="作者" :span="1" class="one-third-width">{{ showData.author }}</el-descriptions-item>
          <el-descriptions-item label="描述" :span="2" class="two-thirds-width">
            {{ showData.description }}
          </el-descriptions-item>

          <!-- 第三行 -->
          <el-descriptions-item label="引用信息" class="one-third-width">{{ showData.reference }}</el-descriptions-item>
          <el-descriptions-item label="修复方式" class="two-thirds-width">{{ showData.remediation }}</el-descriptions-item>
        </el-descriptions>
      </el-dialog>
    </div>
  </template>
  
  <script setup>
  import { getCveList } from "@/api/cve.js"
  import { ref, reactive } from 'vue'
  
  defineOptions({
    name: 'CveData',
  })

  const size = ref('large')
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
  
  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }
  
  // 操作按钮
  const statusData = reactive([
  {
        name: "查看",
        type: "primary",
        icon: "View",
        handleClick: (scope) => handleShow(scope.row),
    },  
  ])
  const searchInfo = ref({})
  //查询
  const onSubmit = () => {
     listQuery.page = 1
     getTableData()
   }
  //重置
   const onReset = () => {
     searchInfo.value = {}   
     getTableData()
   }
  
// 查询
const getTableData = async() => {
    const table = await getCveList({
        page: listQuery.page,
        pageSize: listQuery.pageSize,
        ...searchInfo.value,
    });
    if (table.code === 0) {
        tableData.value = table.data.list;
        listQuery.total = table.data.total;
        listQuery.page = table.data.page;
        listQuery.pageSize = table.data.pageSize;
    }
}
getTableData()
  
const tableColumns = reactive([
    { label:'名称', prop:'name'},
    { label:'等级', prop:'severity', slot: 'customSeverity' },
    { label:'关联模板', prop:'templateId'},
])  

const getStyle = (e) => {
  switch (e) {
      case 'critical':
          return 'danger';
          break;
      case 'height':
          return 'warning';
          break
      case 'medium':
          return 'info';
          break;
      default:
          return 'success'
          break;
  }
}
const getColor = (e) => {
  switch (e) {
      case 'critical':
          return '#E6A23C';
          break;
      case 'height':
          return '#b88230';
          break
      case 'medium':
          return '#f3d19e';
          break;
      default:
          return '#b3e19d'
          break;
  }
}
const getSeverityName = (e) => {
switch (e) {
    case 'critical':
        return '严重';
        break;
    case 'height':
        return '高危';
        break
    case 'medium':
        return '中危';
        break;
    default:
        return '低危'
        break;
}
}
const showFlag = ref(false)
const showData = ref({})
const handleShow = (e) => {
    showFlag.value = true
    showData.value = JSON.parse(JSON.stringify(e))
}
</script>

<style lang="scss">

.one-third-width {
  width: calc(33.33% - 8px); /* 减去左右边距的一半 */
}

/* 定义第二列的宽度 */
.two-thirds-width {
  width: calc(66.66% - 8px); /* 减去左右边距的一半 */
}

/* 调整左右边距以对齐 */
.el-descriptions-item__content,
.el-descriptions-item__label {
  padding-left: 4px;
  padding-right: 4px;
}
</style>
