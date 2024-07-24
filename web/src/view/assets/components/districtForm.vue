<script setup>
  import { ref } from 'vue';
  
  const props = defineProps({
    form: Object,
    rules: Object,
    labelPosition: String,
    itemLabelPosition: String,
  });
  
  const emits = defineEmits(['submit', 'cancel']);
  
  const formRef = ref(null);
  const localForm = ref({ ...props.form });
  
  const submitForm = () => {
    formRef.value.validate((valid) => {
      if (valid) {
        console.log("test")
        // 处理新增 or 修改的数据
        emits('submit', localForm.value);
      } else {
        console.log('Error submit!');
      }
    });
  };
</script>
<template>
    <div>
      <el-form
        :label-position="labelPosition"
        label-width="auto"
        :model="localForm"
        :rules="rules"
        style="max-width: 500px"
        ref="formRef"
      >
        <!-- Form items -->
        <el-form-item label="区域名称" :label-position="itemLabelPosition" prop="areaName">
          <el-input v-model="localForm.areaName" />
        </el-form-item>
        <el-form-item label="IP范围" :label-position="itemLabelPosition" prop="areaIpStr">
            <el-input type="textarea" :rows="6" v-model="localForm.areaIpStr" placeholder="参考：10.0.0.1/24, 10.0.0.1 ~ 10.0.0.255 多个地址段请用逗号分隔" />
        </el-form-item>
        <el-form-item label="备注" :label-position="itemLabelPosition">
          <el-input v-model="localForm.areaDesc" />
        </el-form-item>
  
        <!-- Buttons -->
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="$emit('cancel')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>