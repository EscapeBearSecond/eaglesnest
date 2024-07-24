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
        // 处理新增 or 修改的数据
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
        <el-form-item label="资产名称" :label-position="itemLabelPosition">
          <el-input v-model="localForm.assetName" />
        </el-form-item>
        <el-form-item label="资产IP/范围" :label-position="itemLabelPosition">
          <el-input v-model="localForm.assetIP" />
        </el-form-item>
        <el-form-item label="所属区域" :label-position="itemLabelPosition">
          <el-input v-model="localForm.assetArea" />
        </el-form-item>
        <el-form-item label="资产类型" :label-position="itemLabelPosition">
          <el-input v-model="localForm.assetType" />
        </el-form-item>
        <el-form-item label="资产价值" :label-position="itemLabelPosition">
          <el-input v-model="localForm.manufacturer" />
        </el-form-item>
        <el-form-item label="备注" :label-position="itemLabelPosition">
          <el-input v-model="localForm.desc" />
        </el-form-item>
  
        <!-- Buttons -->
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="$emit('cancel')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>