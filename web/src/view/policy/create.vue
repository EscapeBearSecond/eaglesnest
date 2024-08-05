
<template>
<div>
    <el-row :gutter="10" style="margin-bottom:18px">
        <el-col :span="18" :offset="2">
           
            <el-form
                ref="formRef"
                :model="form"
                :rules="rules"
                style="padding:10px 30px;"
                label-width="100px"
                :label-position="labelPosition"
            >
            <el-form-item label="策略名称"  prop="policyName">
                <el-input v-model="form.policyName" placeholder="请输入策略名称" />
            </el-form-item>
            <el-form-item label=" 其他描述" >
                    <el-input type="textarea" :rows="2" v-model="form.policyDesc" />
                </el-form-item>
            <el-form-item label="策略配置" ></el-form-item>
            <el-tabs type="border-card" style="margin-left: 100px">
                <el-tab-pane label="在线检测">
                    <el-form-item label="策略状态"  class="one-lab">
                            <el-checkbox v-model="form.onlineConfig.use" label="开启"  size="large" />
                    </el-form-item>
                    <div v-if="form.onlineConfig.use">
                        <el-form-item  label="并发数量"  class="sec-lab">
                            <el-input v-model="form.onlineConfig.concurrency" />
                        </el-form-item>
                        <el-form-item label="超时设置"  class="sec-lab">
                            <el-input v-model="form.onlineConfig.timeout" />
                        </el-form-item>
                        <el-form-item label="探活轮次"  class="sec-lab">
                            <el-input v-model="form.onlineConfig.count" />
                        </el-form-item>
                        <el-form-item label="探活频率"  class="sec-lab">
                            <el-input v-model="form.onlineConfig.rateLimit" />
                        </el-form-item>
                    </div>
                </el-tab-pane>
                <el-tab-pane label="端口检测">
                    <el-form-item label="策略状态"  class="one-lab">
                        <el-checkbox v-model="form.portScanConfig.use" label="开启"  size="large" />
                    </el-form-item>
                    <div v-if="form.portScanConfig.use">
                        <el-form-item label="端口范围"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.ports" />
                        </el-form-item>
                        <el-form-item  label="并发数量"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.rateLimit" />
                        </el-form-item>
                        <el-form-item label="超时设置"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.timeout" />
                        </el-form-item>
                        <el-form-item label="探活轮次"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.count" />
                        </el-form-item>
                        <el-form-item label="探活频率"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.rateLimit" />
                        </el-form-item>
                        <el-form-item label="探活频率"  class="sec-lab">
                            <el-input v-model="form.portScanConfig.rateLimit" />
                        </el-form-item>
                    </div>
                </el-tab-pane>
            </el-tabs>
            <el-form-item label="模板配置"  style="margin-top: 15px">
                <el-button type="primary" @click="addTemplate">新增</el-button>
            </el-form-item>
            <div style="margin-left:100px;margin-button:20px">
                <el-table :data="form.policyConfig" style="width: 100%;">
                    <el-table-column type="index" width="60" label="序号" />
                    <el-table-column prop="kind" label="模板类型">
                        <template #default="scope">
                            {{  getKind(scope.row.kind) }}
                        </template>
                    </el-table-column>
                    <el-table-column prop="concurrency" label="最大并发" />
                    <el-table-column prop="timeout" label="超时时间" />
                    <el-table-column prop="rateLimit" label="限流速度"/>
                    <el-table-column prop="count" label="探活轮次" />
                    <el-table-column prop="tag" label="操作" width="80" >
                        <template #default="scope">
                            <el-button type="primary" @click="deleteTemplateConfig">删除</el-button>
                        </template>
                    </el-table-column>
                </el-table>
            </div>
            <div class="btn-save">
                <el-button type="info" @click="goStep">返回</el-button>
                <el-button type="primary" @click="savePolicy">确定</el-button>
            </div>
            </el-form>
        </el-col>
    </el-row>
    <el-drawer
      v-model="templateDialog"
      size="45%"
      :show-close="false"
      :close-on-press-escape="false"
      :close-on-click-modal="false"
    >
        <template #header>
            <div class="flex justify-between items-center">
            <span class="text-lg">配置</span>
            <div>
                <el-button @click="closeDialog">取 消</el-button>
                <el-button
                type="primary"
                @click="enterDialog"
                >保存</el-button>
            </div>
            </div>
        </template>
        <el-form ref="searchRef" :rules="searchRules" :model="searchInfo" label-width="80px">
                <el-form-item label="模板类型"  class="sec-lab" prop="kind">
                    <el-select v-model="searchInfo.kind" placeholder="请选择模板类型" @change="changeScanType(searchInfo, 'kind')">
                        <el-option
                            v-for="type in typeNameList"
                            :key="type.value"
                            :label="type.label"
                            :value="type.value"
                            :disabled="type.disabled"
                        />
                    </el-select>
                </el-form-item>
                <el-form-item label="设备类型"  class="sec-lab"> 
                    <el-select v-model="searchInfo.tagOne" placeholder="请选择执行方式" @change="changeScanType(searchInfo, 'tagOne')">
                        <el-option label="全部" value=""></el-option>
                        <el-option v-for="(tagOne, key) in tagList.tag1" :label="tagOne" :value="tagOne" :key="key" />
                    </el-select>
                </el-form-item>
                <el-form-item label="系统类型"  class="sec-lab">
                    <el-select v-model="searchInfo.tagTwo" placeholder="请选择执行方式" @change="changeScanType(searchInfo, 'tagTwo')">
                        <el-option label="全部" value=""></el-option>
                        <el-option v-for="(tagTwo, key) in tagList.tag2" :label="tagTwo" :value="tagTwo" :key="key" />
                    </el-select>
                </el-form-item>
                <el-form-item label="厂商名称"  class="sec-lab" >
                    <el-select v-model="searchInfo.tagThree" placeholder="请选择执行方式" @change="changeScanType(searchInfo, 'tagThree')">
                        <el-option label="全部" value=""></el-option>
                        <el-option v-for="(tagThree, key) in tagList.tag3" :label="tagThree" :value="tagThree" :key="key" />
                    </el-select>
                </el-form-item>
                <el-form-item label="产品型号"  class="sec-lab">
                    <el-select v-model="searchInfo.tagFour" placeholder="请选择执行方式" @change="changeScanType(searchInfo, 'tagFour')">
                        <el-option label="全部" value=""></el-option>
                        <el-option v-for="(tagFour, key) in tagList.tag4" :label="tagFour" :value="tagFour" :key="key" />
                    </el-select>
                </el-form-item>
                <el-form-item label="模板选择" :label-position="itemLabelPosition"  class="sec-lab" v-if="searchInfo.kind != ''" prop="templates">
                    <el-select 
                        v-model="searchInfo.templates" 
                        placeholder="请选择模板，可多选"   
                        multiple
                        collapse-tags
                        collapse-tags-tooltip
                    >
                    <template #header>
                        <el-checkbox
                        v-model="checkAll"
                        :indeterminate="indeterminate"
                        @change="handleCheckAll(checkAll, searchInfo.kind)"
                        >
                        全选
                        </el-checkbox>
                    </template>
                        <el-option
                        v-for="item in tmpOption[searchInfo.kind - 1]"
                        :key="item.value"
                        :label="item.label"
                        :value="item.value"
                        :disabled="item.disabled"
                        />
                    </el-select>
                </el-form-item>  
                <el-form-item label="最大并发" :label-position="itemLabelPosition" class="sec-lab" prop="concurrency">
                    <el-input v-model="searchInfo.concurrency" placeholder="请输入最大并发"  />
                </el-form-item>
                <el-form-item label="超时时间" :label-position="itemLabelPosition" class="sec-lab"  prop="timeout">
                    <el-input v-model="searchInfo.timeout" placeholder="请输入超时时间"  />
                </el-form-item>
                <el-form-item label="限流速度" :label-position="itemLabelPosition" class="sec-lab"  prop="rateLimit">
                    <el-input v-model="searchInfo.rateLimit" placeholder="请输入限流速度"  />
                </el-form-item>
                <el-form-item label="探活轮次" :label-position="itemLabelPosition" class="sec-lab"  prop="count">
                    <el-input v-model="searchInfo.count" placeholder="请输入探活轮次"  />
                </el-form-item> 
        </el-form>
    </el-drawer>     
</div>
</template>
<script setup>
import { ref, reactive } from 'vue' 
import { getPolicyList, createPolicy, updatePolicy, getPolicyId } from '@/api/policy'
import { getTemplateTagList, getTemplateList } from '@/api/template'
import router from '@/router/index'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useRoute } from 'vue-router';


const formRef = ref(null)
const form = ref(
    {
        policyName: '',
         policyDesc: '',
         headlessFlg: '',
         scanType: '',
         scanRate: '',
         policyConfig: [],
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
const labelPosition = ref('left')
const itemLabelPosition = ref('left')

const editableTabsValue = ref('config1')
//模板类型筛选
const typeNameList = reactive([
  {id: '1', label: '资产发现', value:'1', disabled: false},
  {id: '2', label: '漏洞扫描', value:'2', disabled: false},
  {id: '3', label: '弱口令', value:'3', disabled: false}
])
//获取四层筛选
const tagList = ref({})
const getTemplateTagData = async () => {
     const data = await getTemplateTagList()
     tagList.value = data.data
}
getTemplateTagData()


// 获取模板
const checkAll = ref(false)
const indeterminate = ref(false)

// 筛选模板联动
const changeScanType = (e, f) => {
    console.log(
    '%c 🍱 CONSOLE_INFO: ',
    'font-size:20px;background-color: #ED9EC7;color:#fff;',
    e, f
    );
    if(f == 'kind') {
        searchInfo.value.templates = []
    }
    updatetmpOption(e.kind)
}

const updatetmpOption = async (kind) => {
    const table = await getTemplateList({
        page: 1,
        pageSize: 99999,
        isAll: false,
        tag1: searchInfo.value.tagOne,
        tag2: searchInfo.value.tagTwo,
        tag3: searchInfo.value.tagThree,
        tag4: searchInfo.value.tagFour,
    });
    tmpOption[kind - 1] = []
    table.data.list.forEach(e => {
        tmpOption[kind - 1].push({label:e.templateName, value: e.ID})
    })
}

const route = useRoute();
const id = ref(route.query.id);
const initForm = async () => {
    // 修改
    if(id.value  != undefined) {
        let data = await getPolicyId({id: id.value})
        form.value = data.data
    }
}
initForm();

// 全选模板
const handleCheckAll = (e, f) => {
    if(e) {
        searchInfo.value.templates = tmpOption[e - 1].map((_)=> _.value)
    } else {
        searchInfo.value.templates = []
    }
}

 // 配置选中扫描类型时返回模板
 const searchRules = ref({
    kind: [
      { required: true, message: '请选择模板类型', trigger: 'blur' }
    ],
    concurrency: [
      { required: true, message: '最大并发未填写', trigger: 'blur' }
    ],
    timeout: [
      { required: true, message: '超时时间未填写', trigger: 'blur' }
    ],
    rateLimit: [
      { required: true, message: '限流速度未填写', trigger: 'blur' }
    ],
    count: [
      { required: true, message: '探活轮次未填写', trigger: 'blur' }
    ],
    templates:[
        { required: true, message: '请选择模板', trigger: 'blur' }
    ]
 })
 const tmpOption = [[],[],[]]
 const searchRef = ref(null)
 const searchInfo = ref({
    "tag1":"",
    "tagTwo":"",
    "tagThree":"",
    "tagFour":"",
    "name": "",
    "kind": "1",
    "timeout": "5s",
    "count": 1,
    "format": "",
    "rateLimit": 150,
    "concurrency": 150,
    "templates":[]
 })
 const getTemplateData = async () => {
    const table = await getTemplateList({
        page: 1,
        pageSize: 99999,
        isAll: false
    });
    table.data.list.forEach(e => {
        if(e.templateType == 1) {
          tmpOption[0].push({label:e.templateName, value: e.ID})
        }
        if(e.templateType == 2) {
          tmpOption[1].push({label:e.templateName, value: e.ID})
        }
        if(e.templateType == 3) {
          tmpOption[2].push({label:e.templateName, value: e.ID})
        }
    });    
 }
 getTemplateData()

const templateDialog = ref(false)
const addTemplate = () => {
    templateDialog.value = true
}

const closeDialog = ()=> {
    onReset()
    templateDialog.value = false
}

const onReset = () => {
    searchInfo.value = {}
    searchInfo.value.templates = []
}

// 保存模板选择
const enterDialog = () => {
    let pushData = JSON.parse(JSON.stringify(searchInfo.value))
    // 如果已经添加同一类型就提醒不能添加
    let flag = false
    form.value.policyConfig.forEach(item => {
        if(item.kind === pushData.kind) {
            flag =  true
        }
    });
    if(!flag) {
        form.value.policyConfig.push(pushData)
        closeDialog()
    }else {
        ElMessage({
            type: 'warning',
            message: '策略已经存在相同类型模板!'
        })
    }
}

// 删除模板配置
const deleteTemplateConfig = () => {
    form.value.policyConfig.pop()
}

// 获取模板类型 
const getKind = (e)=> {
    const item = typeNameList.find(item => item.id === e);
    return item ? item.label : null;
}

const goStep = () => {
    window.history.go(-1)
}

const savePolicy = async () => {
    formRef.value.validate(async valid => {
        if(form.value.id != '' && form.value.id == undefined) {

            const res = await createPolicy(form.value)
             if (res.code === 0) {
               ElMessage({
                 type: 'success',
                 message: '添加成功!'
               })
               window.history.go(-1)
            }
        }else {
            const res = await updatePolicy(form.value)
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '修改成功!'
                })
                window.history.go(-1)
            }
        }
    })  
}
</script>
<style lang='scss' scoped>
.btn-save {
    display: flex;
    justify-content:center;
    margin: 10px;
}
</style>
