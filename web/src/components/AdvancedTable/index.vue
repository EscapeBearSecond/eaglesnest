<template>
  <el-table 
    :data="tableData"
    :columns="columns"
    :row-key="row => row.ID"
    header-align="center"
    @select="onTableSelect"
    @select-all="onTableSelectAll" 
    :fit="true"
    border
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
      width="50"
    />
    <el-table-column
      v-for="(column, index)  in columns"
      :key="index"
      :label="column.label"
      :prop="column.prop"
      :width="column.width"
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
      v-if="statusData && statusData.length > 0"
      :width="statusWidth"
      align="left"
    >
      <template #default="scope">
          <template v-for="(btn, index) in statusData">
              <template v-if="btn.visible !== false && (typeof btn.visible !== 'function' || btn.visible(scope, btn))">
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
import { defineComponent, watchEffect } from 'vue';

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
    selectedRows: {
      default: [],
      type: Array,
      required: false
    },
    listQuery:{
      default: null,
      type:Object,
      required: true
    },
    selection: { default: false, type: Boolean },
    statusData: { default: null, type: Array },
    statusWidth: { default: null, type:  [String, Number] },
    index: { default: false, type: Boolean },
    pagination: { default: null, type: Function },
    changePageSize:{default:null, type:Function},
    selectionRow:{ default:null, type:Function},
    selectionAll:{ default:null, type:Function}
  },
  emits: ['update:tableData', 'update:listQuery', 'pagination'],
  setup(props, { emit }) {
    
    const handlePagination = (val) => {
      props.listQuery.page = val;
      emit('update:listQuery', { ...props.listQuery, page: val });
      if (props.pagination) {
        props.pagination(val);
      }

       // 在翻页时，根据 selectedRows 重新设置选中状态
       restoreSelection();
    };

    const onTableSelect = (selection, row) => {
      const selectedStatus = selection.length && selection.includes(row);
      props.selectionRow && props.selectionRow(selection, row, selectedStatus);
    };
    const onTableSelectAll = (selection) => {
      props.selectedRows.splice(0, props.selectedRows.length, ...selection); // 当全选/取消全选时，更新 selectedRows
      props.selectionAll && props.selectionAll(selection);
    };

    const handleSizeChange = (val) => {
      props.listQuery.pageSize = val;
      if (props.changePageSize) {
        props.changePageSize(val);
      }
    }

    // 监听 tableData 变化时恢复选中状态
    watchEffect(() => {
      restoreSelection();
    });

     // 恢复选中状态
     function restoreSelection() {
      
      const newRowSelection = props.selectedRows.map(row => ({
        id: row.ID,
        selected: false,
      }));

      // 将之前选中的行与当前表格数据进行匹配
      props.tableData.forEach(tableRow => {
        newRowSelection.forEach((rowSelection, index) => {
          if (tableRow.id === rowSelection.id) {
            newRowSelection[index].selected = true;
          }
        });
      });

      // 将新的选择状态应用到表格上
      const newSelection = newRowSelection.filter(row => row.selected).map(row => row.id);
      emit('update:tableData', { ...props.tableData, selectedRows: newSelection });
    }

    return {
      handleSizeChange,
      handlePagination,
      onTableSelectAll,
      onTableSelect,
      restoreSelection
    };
  }

  
});
</script>

<style scoped>
 .hide-scrollbar .el-table__body-wrapper {
  overflow-x: hidden !important;
}
</style>
