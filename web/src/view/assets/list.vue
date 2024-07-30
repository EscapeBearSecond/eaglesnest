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
  
  </div>
</template>

<script setup>
import { getListApi } from "@/api/assets"

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

const initPage = async() => {
  getTableData()
}

initPage()

const tableColumns = reactive([
    { label:'名称', prop:'assetName'},
    { label:'IP', prop:'assetIP'},
    { label:'归属', prop:'assetArea'},
    { label:'类型', prop:'assetType'},
    { label:'备注', prop:'desc'},
])

</script>

<style lang="scss">
  .header-img-box {
    @apply w-52 h-52 border border-solid border-gray-300 rounded-xl flex justify-center items-center cursor-pointer;
 }
</style>
