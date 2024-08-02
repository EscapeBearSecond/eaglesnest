
<template>
<div>
    <el-row :gutter="10">
        <el-col :span="16">
            <el-form
                ref="formRef"
                :model="form"
                :rules="rules"
                style="padding:10px 20px;"
                label-width="100px"
            >
            <el-form-item label="策略名称" :label-position="itemLabelPosition" prop="policyName">
                <el-input v-model="form.policyName" placeholder="请输入策略名称" />
            </el-form-item>
            <el-form-item style="padding-left: 8px" label="策略配置" :label-position="itemLabelPosition"></el-form-item>
            <el-tabs type="border-card" style="margin-left: 15px">
                <el-tab-pane label="在线检测">
                    <el-form-item label="策略状态" :label-position="itemLabelPosition" class="one-lab">
                            <el-checkbox v-model="form.onlineConfig.use" label="开启"  size="large" />
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
                <el-tab-pane label="端口检测">
                    <el-form-item label="策略状态" :label-position="itemLabelPosition" class="one-lab">
                        <el-checkbox v-model="form.portScanConfig.use" label="开启"  size="large" />
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
            <div style="margin: 10px 0 10px 10px;padding-top: 5px;">
                <label  class="el-form-item__label">模板配置：</label>
                <el-tabs
                    v-model="editableTabsValue"
                    type="card"
                    class="demo-tabs"
                    closable
                    @tab-remove="removeTab"
                >
                    <el-tab-pane
                        v-for="item in editableTabs"
                        :key="item.name"
                        :label="item.title"
                        :name="item.name"
                        >
                    {{ item.content }}
                    </el-tab-pane>
                </el-tabs>
            </div>
            <el-form-item label=" 其他描述：" :label-position="itemLabelPosition">
                <el-input type="textarea" :rows="3" v-model="form.policyDesc" />
            </el-form-item>
            </el-form>
        </el-col>
        <el-col :span="8">
        </el-col>
    </el-row>
     
</div>
</template>
<script setup>
import { ref, reactive } from 'vue' 
import { getPolicyList } from '@/api/policy'

const formRef = ref(null)
const form = ref(
    {
        policyName: '',
         policyDesc: '',
         headlessFlg: '',
         scanType: '',
         scanRate: '',
         policyConfig: [{
          "name": "",
          "kind": "",
          "timeout": "5s",
          "count": 1,
          "format": "",
          "rateLimit": 150,
          "concurrency": 150
         }],
         "onlineConfig": {
          "use": true,
          "timeout": "5s",
          "count": 1,
          "format": "csv",
          "rateLimit": 150,
          "concurrency": 150
         },
         "portScanConfig": {
          "use": true,
          "timeout": "5s",
          "count": 1,
          "format": "csv",
          "ports": "http",
          "rateLimit": 150,
          "concurrency": 150
         }
    }
)
const rules = ref({
  policyName: [
     { required: true, message: '请输入策略名称', trigger: 'blur' }
  ],
 })
const labelPosition = ref('right')
const itemLabelPosition = ref('left')

const editableTabs= ref([
  {
    title: '配置',
    name: 'config',
    content: 'Tab 1 content',
  }
])


</script>
<style lang='scss' scoped>
  
</style>
