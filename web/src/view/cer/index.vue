<template>
  <div class="cer-box">
    <hr style="height: 2px; border: none; border-top: 2px dotted #185598" />
    <h1>证书已过期</h1>
    <div class="tip-text">
      <span class="tipw">您的使用证书已过期,系统功能已被限制</span>
      <span class="tipm">请您及时更新证书以继续使用完整功能</span>
      <span class="tipn">请联系我们获取新的有效证书:</span>
    </div>
    <hr style="height: 2px; border: none; border-top: 2px dotted #185598" />
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
      </el-upload>
    </div>
  </div>
</template>
<script setup>
<<<<<<< HEAD
import { ref, reactive } from 'vue' 
import { uploadLicense } from '@/api/api'
import router from '@/router/index'
import { ElMessage } from 'element-plus'
=======
import { ref } from "vue";
import { uploadLicense } from "@/api/api";
import router from "@/router/index";
import { ElMessage } from "element-plus";
>>>>>>> 6815ed01d3ddf10c23b0e1537960aa0c3fafa8e0

const fileList = ref([]);
const handleCustomUpload = async (file) => {
  const formData = new FormData();
<<<<<<< HEAD
  formData.append('license', file);
  
=======
  formData.append("license", file);
  console.log(file);
>>>>>>> 6815ed01d3ddf10c23b0e1537960aa0c3fafa8e0
  // 阻止默认的上传行为
  let data = await uploadLicense(formData);
  if (data.code === 0) {
<<<<<<< HEAD
      ElMessage({ type: 'success', message: '更新成功，请重新登陆！' })
      router.push({ name: 'Login', replace: true })
  }else {
      ElMessage({ type: 'error', message: '更新失败，请联系相关人员操作！' })
=======
    ElMessage({ type: "success", message: data.msg });
    router.push({ name: "Login", replace: true });
  } else {
    ElMessage({ type: "error", message: data.msg });
>>>>>>> 6815ed01d3ddf10c23b0e1537960aa0c3fafa8e0
  }
  return false;
};
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

.tipw,
.tipm {
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
