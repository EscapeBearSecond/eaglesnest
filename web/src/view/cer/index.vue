
<template>
    <div class="cer-box">
        <hr style=" height:2px;border:none;border-top:2px dotted #185598;" />
        <h1>证书已</h1>
            <div class="tip-text">
                <span class="tipw">过期,当前无法继续使用</span>
                <span class="tipm">请您</span>
                <span class="tipn">上传新证书或联系技术支持以获取更多信息.</span>
            </div>
            <hr style=" height:2px;border:none;border-top:2px dotted #185598;" />
            <div class="btn">
              <el-upload
                    v-model:file-list="fileList"
                    class="upload-demo"
                    action="''"
                    accept=".json, .JSON"
                    multiple
                    :limit="1"
                    :before-upload="handleCustomUpload"
                >
                <el-button type="primary">导入证书</el-button>
                <template #tip>
                    <div class="el-upload__tip">
                        文件类型：JSON文件, 大小不超过 200KB.
                    </div>
                    </template>
                </el-upload>
            </div>
    </div>
</template>
<script setup>
import { ref, reactive } from 'vue' 
import { uploadLicense } from '@/api/api'
import router from '@/router/index'

const fileList = ref([])
const handleCustomUpload = async(file) => {
  const formData = new FormData();
  formData.append('license', file);
  console.log(file)
  // 阻止默认的上传行为
  let data = await uploadLicense(formData)
  if (data.code === 0) {
      getData()
      ElMessage({ type: 'success', message: '更新成功，请重新登陆！' })
      router.push({ name: 'Login', replace: true })
  }else {
      ElMessage({ type: 'error', message: '更新失败，请联系相关人员操作！' })
  }
  return false;
}
</script>
<style lang='scss' scoped>
  .cer-box {
    width: 50%;
    margin: 5% auto;
  }

  .tip-text {
    font-size: 16px;
    margin: 30px 20px;
  }

  .tipw, .tipm {
    display: block;
    font-size: 16px;
    margin: 25px 0px;
  }

  .tipn {
    margin-left: 20px;
  }

  .btn {
    margin: 20px 0;
    align-items: right;
  }
</style>
