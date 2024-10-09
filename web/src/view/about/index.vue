
<template>
    <div class="info">
        <el-descriptions
            title="信息介绍"
            direction="horizontal"
            :column="1"
            :size="size"
            border
            
        >
            <el-descriptions-item label-align="center" label="系统版本">
                <el-tag size="default">{{ info.systemVersion }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label-align="center" label="漏洞库版本">
                <el-tag size="default">{{ info.vulnVersion  }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label-align="center" label="证书过期时间" :span="2">
                <el-tag size="default">{{ info.licenseExpiration }}</el-tag> 
                
            </el-descriptions-item>
            <el-descriptions-item label-align="center" label="最近更新时间">
                <el-tag size="default">{{ info.lastUpdateDate }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label-align="center" label="证书导入">
                <el-upload
                    v-model:file-list="fileList"
                    class="upload-demo"
                    action="''"
                    multiple
                    :limit="1"
                    :before-upload="handleCustomUpload"
                >
                <el-button size="small" type="primary">更新</el-button>
                <template #tip>
                    <div class="el-upload__tip">
                        文件格式：json 类型, 大小不超过 200KB.
                    </div>
                    </template>
                </el-upload>
            </el-descriptions-item>
        </el-descriptions>
    </div>
</template>
<script setup>
import { ref } from 'vue' 
import { getSystemInfo, uploadLicense } from '@/api/api'
import { ElMessage } from 'element-plus'

const size = ref('default')
const info = ref({})
const getData = async() => {
  const data = await getSystemInfo({})
   console.log(data)
    if(data.code === 0) {
        info.value = data.data
        ElMessage({ type: 'success', message: '更新成功！' })
    }else {
        ElMessage({ type: 'error', message: '更新失败！' })
    }
}

const handleCustomUpload = async(file) => {

    const formData = new FormData();
    formData.append('license', file);
    console.log(file)
    // 阻止默认的上传行为
    let data = await uploadLicense(formData)
    if (data.code === 0) {
        getData()
        ElMessage({ type: 'success', message: '更新成功！' })
    }else {
        ElMessage({ type: 'error', message: '更新失败！' })
    }
    return false;
}
getData()
</script>
<style lang='scss' scoped>
  .info {
    width: 50%;
    margin: 5% auto;
  }
</style>
