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
</script>
<template>
    <el-form
      :label-position="labelPosition"
      label-width="auto"
      :model="form"
      :rules="rules"
      style="max-width: 500px"
    >
      <el-form-item label="策略名称" :label-position="itemLabelPosition">
        <el-input v-model="form.policyName" />
      </el-form-item>
      <el-form-item label="策略描述" :label-position="itemLabelPosition">
        <el-input v-model="form.policyDesc" />
      </el-form-item>
      <label  class="el-form-item__label">策略选项</label>
      <el-form-item :label-position="itemLabelPosition" class="one-lab">
        <el-checkbox v-model="form.onlineConfig.use" label="在线检测"  size="large" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.onlineConfig.concurrency" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.onlineConfig.timeout" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.onlineConfig.count" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.onlineConfig.rateLimit" />
      </el-form-item>
      <el-form-item :label-position="itemLabelPosition" class="one-lab">
        <el-checkbox v-model="form.portScanConfig.use" label="端口检测"  size="large" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="端口范围" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.ports" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="并发数量" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.rateLimit" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="超时设置" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.timeout" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="探活轮次" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.count" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.rateLimit" />
      </el-form-item>
      <el-form-item v-if="form.onlineConfig.use" label="探活频率" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.portScanConfig.rateLimit" />
      </el-form-item>
      <el-form-item :label-position="itemLabelPosition" class="one-lab">
        <el-checkbox v-model="form.headlessFlg" label="Headless模式"  size="large" />
      </el-form-item>
      <el-form-item v-if="form.headlessFlg" label="扫描类型" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.scanType" />
      </el-form-item>
      <el-form-item v-if="form.headlessFlg" label="扫描速度" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.scanRate" />
      </el-form-item>
      <el-form-item v-if="form.headlessFlg" label="扫描模板" :label-position="itemLabelPosition" class="sec-lab">
        <el-input v-model="form.scanTmp" />
      </el-form-item>
    </el-form>

  </template>
  

  
<style scoped>
  .one-lab {
    margin-left: 30px;
  }

  .sec-lab {
    margin-left: 50px;
  }
</style>
