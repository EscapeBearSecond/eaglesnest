<script  setup>
import { reactive, ref } from 'vue'
const labelPosition = ref('right')
const itemLabelPosition = ref('top')
const form = reactive({
    policyName: '',
    policyDesc: '',
    policyConfig: [{
        "name": "",
        "kind": "",
        "timeout": "",
        "count": 0,
        "format": "",
        "rateLimit": 0,
        "concurrency": 0
    }],
    "onlineConfig": {
        "use": true,
        "timeout": "5s",
        "count": 1,
        "format": "csv",
        "rateLimit": 1000,
        "concurrency": 1000
    },
    "portScanConfig": {
        "use": true,
        "timeout": "5s",
        "count": 1,
        "format": "csv",
        "ports": "http",
        "rateLimit": 1000,
        "concurrency": 1000
    }
})
const editableTabsValue = ref('3')
const handleTabsEdit = () => {

}

</script>
<template>
    <el-form
      :label-position="labelPosition"
      label-width="auto"
      :model="form"
      :rules="rules"
      style="max-width: 600px"
    >
      <el-form-item label="策略名称" :label-position="itemLabelPosition">
        <el-input v-model="form.policyName" />
      </el-form-item>
      <el-form-item label="策略描述" :label-position="itemLabelPosition">
        <el-input v-model="form.policyDesc" />
      </el-form-item>
      <label  class="el-form-item__label">策略选项</label>
      <el-tabs type="border-card" class="demo-tabs">
        <el-tab-pane>
          <template #label>
            <span class="custom-tabs-label">
              <el-icon><calendar /></el-icon>
              在线检测
            </span>
          </template>
          <el-form-item :label-position="itemLabelPosition" class="one-lab">
                <el-checkbox v-model="form.onlineConfig.use" label="在线检测"  size="large" />
          </el-form-item>
          <div v-if="form.onlineConfig.use">
            <el-form-item  label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
            <el-input v-model="form.onlineConfig.concurrency" />
          </el-form-item>
          <el-form-item label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
            <el-input v-model="form.onlineConfig.timeout" />
          </el-form-item>
          <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
            <el-input v-model="form.onlineConfig.count" />
          </el-form-item>
          <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
            <el-input v-model="form.onlineConfig.rateLimit" />
          </el-form-item>
          </div> 
    </el-tab-pane>
    <el-tab-pane label="Config">
      <template #label>
        <span class="custom-tabs-label">
          <el-icon><calendar /></el-icon>
          端口检测
        </span>
      </template>
        <el-form-item :label-position="itemLabelPosition" class="one-lab">
            <el-checkbox v-model="form.portScanConfig.use" label="端口检测"  size="large" />
          </el-form-item>
        <div v-if="form.portScanConfig.use">
        <el-form-item label="端口范围" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.ports" />
        </el-form-item>
        <el-form-item  label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.rateLimit" />
        </el-form-item>
        <el-form-item label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.timeout" />
        </el-form-item>
        <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.count" />
        </el-form-item>
        <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.rateLimit" />
        </el-form-item>
        <el-form-item label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
          <el-input v-model="form.portScanConfig.rateLimit" />
        </el-form-item>
      </div>
    </el-tab-pane>
      </el-tabs>
      <label  class="el-form-item__label">模板设置</label> 
      
    </el-form>

  </template>
  

  
<style scoped>
  .one-lab {
    margin-left: 30px;
  }

  .sec-lab {
    margin-left: 50px;
  }

  .demo-tabs > .el-tabs__content {
  padding: 32px;
  color: #6b778c;
  font-size: 32px;
  font-weight: 600;
}
</style>
