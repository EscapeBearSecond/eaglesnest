
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
                <el-button type="primary">上传证书</el-button>
            </div>
    </div>
</template>
<script setup>
import { ref, reactive } from 'vue' 

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
