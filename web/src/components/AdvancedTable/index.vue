<template>
  <el-table 
    :data="tableData"
    :columns="columns"
    :row-key="row => row.id"
    header-align="center"
    @select="onTableSelect"
    @select-all="onTableSelectAll" 
    :fit="true"
    :style="{'thead-tr-height': '20px', 'width': '100%', 'overflow-x': 'hidden'}">
    <!-- <template #empty>
      <el-empty description="未查询到相关数据" />
    </template> -->
    <el-table-column
      v-if="index"
      type="index"
      width="60"
      label="序号"
      align="center"
    >
      <template #default="scope">
        {{ scope.$index + 1 + (listQuery.page - 1) * listQuery.pageSize }}
      </template>
    </el-table-column>
    <el-table-column
      v-if="selection"
      reserve-selection
      type="selection"
      width="45"
    />
    <el-table-column
      v-for="(column, index)  in columns"
      :key="index"
      :label="column.label"
      :prop="column.prop"
      :formatter="column.formatter"
      align="center"
    >
    <!-- 自定义单元格内容 -->
    <template v-if="column.slot" #default="scope">
        <slot :name="column.slot" :row="scope.row" />
      </template>
    </el-table-column>
    <el-table-column
      prop="status"
      label="操作"
      v-if="statusData"
      :width="statusWidth"
      align="center"
    >
      <template #default="scope">
          <template v-for="(btn, index) in statusData">
              <template v-if="btn.visible ? btn.visible(scope, btn) : true">
                  <el-button link :icon="btn.icon ? btn.icon : ''" :type="btn.type" @click="btn.handleClick && btn.handleClick(scope)" :key="index">{{  btn.name }}</el-button>
              </template>
          </template>
      </template>
    </el-table-column>
  </el-table>
  <pagination
    v-if="listQuery && pagination"
      :page="listQuery.page"
      :page-size="listQuery.pageSize"
      :total="listQuery.total"
      @size-change="handleSizeChange"
      @current-change="handlePagination"
  />
</template>

<script>
import { defineComponent } from 'vue';

export default defineComponent({
  name: 'CommonTable',
  props: {
    columns: {
      default: [],
      type: Array,
      required: true
    },
    tableData: {
      default: [],
      type: Array,
      required: true
    },
    listQuery:{
      default: null,
      type:Object,
      required: true
    },
    selection: { default: false, type: Boolean },
    statusData: { default: null, type: Array },
    statusWidth: { default: null, type: String },
    index: { default: false, type: Boolean },
    pagination: { default: null, type: Function },
  },
  emits: ['update:tableData', 'update:listQuery', 'pagination'],
  setup(props, { emit }) {
    const handlePagination = (val) => {
      props.listQuery.page = val;
      emit('update:listQuery', { ...props.listQuery, page: val });
      if (props.pagination) {
        props.pagination(val);
      }
    };

    const onTableSelect = (selection, row) => {
      const selectedStatus = selection.length && selection.includes(row);
      props.selectionRow && props.selectionRow(selection, row, selectedStatus);
      };
    const onTableSelectAll = (selection) => {
      props.selectionAll && props.selectionAll(selection);
    };

    const handleSizeChange = (val) => {
      props.listQuery.value.pageSize = val;
    }

    return {
      handleSizeChange,
      handlePagination,
      onTableSelectAll,
      onTableSelect
    };
  }

  
});
</script>

<style scoped>
 .hide-scrollbar .el-table__body-wrapper {
  overflow-x: hidden !important;
}
</style>
