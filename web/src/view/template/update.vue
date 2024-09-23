
<template>
    <div style="width: 80%;margin: 40px auto;">
        <el-upload
        ref="uploadRef"
        class="updateTemp"
        drag
        :action="''"
        :auto-upload="false"
        multiple
        :limit="1"
        :on-exceed="handleExceed"
        accept=".zip,.enc" 
        @change="handleFileChange"
      >
          <el-icon class="el-icon--upload"><upload-filled /></el-icon>
          <div class="el-upload__text">
            拖动文件到这里 或<em> 点击上传</em>
          </div>
          <!-- <template #tip>
            <div class="el-upload__tip">
             注：文件数量不可超过1个;
            </div>
          </template> -->
        </el-upload>
        <div style="text-align: center;"><el-button type="primary" @click="saveFile">更新</el-button></div>
    </div>
</template>
<script setup>
import { ref } from 'vue' 
import { updateTemplateContent } from '@/api/template'
import { ElMessage } from 'element-plus'

const uploadRef = ref(null);
const selectedFiles = ref([]);
const handleFileChange = (file, fileList) => {
  selectedFiles.value = fileList.map(item => item.raw);
};

const saveFile =  async() => {
    console.log(selectedFiles.value.length)
    if (selectedFiles.value.length === 0) {
        ElMessage({ type: 'error', message: '未选中文件' })
        return;
    }
    if (selectedFiles.value.length > 1) {
        ElMessage({ type: 'error', message: '当前只支持单文件' })
        return;
    }

    const formData = new FormData();
    for(let i = 0; i < selectedFiles.value.length; i++) {
        formData.append('file', selectedFiles.value[i]);
    }
    
    let data = await updateTemplateContent(formData)
    if (data.code === 0) {
        ElMessage({ type: 'success', message: '提交成功' })
    }
}

const  handleExceed = (files, fileList) => {
    this.$message.warning('最多只能上传1个文件');
  }


</script>
<style lang='scss' scoped>
  
</style>
