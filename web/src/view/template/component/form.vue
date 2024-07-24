<script setup>
  import { ref } from 'vue';
  
  const props = defineProps({
    form: Object,
    rules: Object,
    labelPosition: String,
    itemLabelPosition: String,
  });

  const options = ref([
    {label: "漏洞扫描", value: '1'},
    {label: "资产发现", value: '2'},
    {label: "弱口令", value: '3'},
  ])
  
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
        ref="formRef"
      >
        <!-- Form items -->
        <el-form-item label="模板名称" :label-position="itemLabelPosition">
          <el-input v-model="localForm.templateName" />
        </el-form-item>
        <el-form-item label="模板类型"  :label-position="itemLabelPosition">
            <el-select
                v-model="localForm.templateType"
                placeholder="选择模板类型"
                size="large"
                >
                <el-option
                    v-for="item in options"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                />
                </el-select>
        </el-form-item>
        <el-form-item label="模板内容" :label-position="itemLabelPosition">
          <el-input type="textarea" :rows="18" v-model="localForm.templateContent" />
        </el-form-item>  
        <!-- Buttons -->
        <el-form-item>
          <el-button type="primary" @click="submitForm">保存</el-button>
          <el-button @click="$emit('cancel')">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </template>