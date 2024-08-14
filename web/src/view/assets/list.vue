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
             v-model="searchInfo.assetName"
             placeholder="请输入区域名称"
           />
         </el-form-item>
         <el-form-item label="设备类型"  class="sec-lab"> 
              <el-select v-model="searchInfo.assetType" placeholder="请选择设备类型">
                  <el-option label="全部" value=""></el-option>
                  <el-option v-for="(tagOne, key) in tagList.tag1" :label="tagOne" :value="tagOne" :key="key" />
              </el-select>
          </el-form-item>
          <el-form-item label="系统类型"  class="sec-lab">
              <el-select v-model="searchInfo.systemType" placeholder="请选择系统类型" >
                  <el-option label="全部" value=""></el-option>
                  <el-option v-for="(tagTwo, key) in tagList.tag2" :label="tagTwo" :value="tagTwo" :key="key" />
              </el-select>
          </el-form-item>
          <el-form-item label="厂商名称"  class="sec-lab" >
              <el-select v-model="searchInfo.manufacturer" placeholder="请选择厂商名称" >
                  <el-option label="全部" value=""></el-option>
                  <el-option v-for="(tagThree, key) in tagList.tag3" :label="tagThree" :value="tagThree" :key="key" />
              </el-select>
          </el-form-item>
          <el-form-item label="产品型号"  class="sec-lab">
              <el-select v-model="searchInfo.assetModel" placeholder="请选择产品型号" >
                  <el-option label="全部" value=""></el-option>
                  <el-option v-for="(tagFour, key) in tagList.tag4" :label="tagFour" :value="tagFour" :key="key" />
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
        :changePageSize="changeSize"
        :index="true"
      >
      <template v-slot:customAreaName="slotProps">
        <!-- 自定义的字段 -->
        <span> 
          <el-tag effect="dark" type="warning">{{ getAreaName(slotProps.row.assetArea) }}</el-tag>
        </span>
      </template>
      <template v-slot:customProt="slotProps">
        <!-- 自定义的字段 -->
        <span v-for="(val, key) in slotProps.row.openPorts" :key="key" style="margin-left: 2px;"> 
          <el-tag effect="dark"  type="primary">{{ val }}</el-tag>
        </span>
      </template>
      </advance-table>

    </div>
  
  </div>
</template>

<script setup>
import { getListApi } from "@/api/assets"
import { getAreaList } from "@/api/area"
import { getTemplateTagList } from '@/api/template'
import { ref, reactive } from 'vue'

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
const changeSize = (e) => {
  listQuery.page = 1
  listQuery.pageSize = e
  getTableData()
}

//获取四层筛选
const tagList = ref({})
const getTemplateTagData = async () => {
     const data = await getTemplateTagList()
     tagList.value = data.data
}
getTemplateTagData()

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 操作
const statusData = reactive([])

const searchInfo = ref({})
const onSubmit = () => {
   listQuery.page = 1
   getTableData()
 }

 const onReset = () => {
   searchInfo.value = {}
   getTableData()   
 }

// 查询
const getTableData = async() => {
  const table = await getListApi({
      page: listQuery.page,
      pageSize: listQuery.pageSize,
      // ...searchInfo,
    });
    if (table.code === 0) {
      tableData.value = table.data.list;
      listQuery.total = table.data.total;
      listQuery.page = table.data.page;
      listQuery.pageSize = table.data.pageSize;
    }    
}
const areaList = ref([])
const getAreaData = async() => {
  const table = await getAreaList({
      page: 1,
      pageSize: 100,
  });

  areaList.value = table.data.list
}

const initPage = async() => {
  getTableData()
  getAreaData()
}

initPage()

const tableColumns = reactive([
    { label:'名称', prop:'assetName'},
    { label:'IP', prop:'assetIp'},
    { label:'区域', prop:'assetArea', slot: 'customAreaName'},
    { label:'类型', prop:'assetType'},
    { label:'开放端口', prop:'openPorts', slot: 'customProt'},
])

const getAreaName = (e) => {
  const item = areaList.value.find((item) => item.ID == e);
  return item ? item.areaName : '未知';
};


</script>

<style lang="scss">
</style>
