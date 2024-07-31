<template>
  <div>
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

const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 操作
const statusData = reactive([])

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
