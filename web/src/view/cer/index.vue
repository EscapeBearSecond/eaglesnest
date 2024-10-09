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
import { ref } from "vue";
import { uploadLicense } from "@/api/api";
import router from "@/router/index";
import { ElMessage } from "element-plus";

const fileList = ref([]);
const handleCustomUpload = async (file) => {
  const formData = new FormData();
  formData.append("license", file);

  // 阻止默认的上传行为
  let data = await uploadLicense(formData);
  if (data.code === 0) {
    ElMessage({ type: "success", message: data.msg });
    router.push({ name: "Login", replace: true });
  } else {
    ElMessage({ type: "error", message: data.msg });
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
