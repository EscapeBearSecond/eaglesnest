<template>
    <el-form :model="localForm" :rules="localRules" ref="ruleFormRef" label-width="100px" class="demo-ruleForm">
      <el-form-item v-for="(item, index) in localForm" :key="index" :label="item.label" :prop="item.prop">
        <el-input v-if="item.type === 'input'" :type="item.type" v-model="localForm[item.prop]"></el-input>
        <el-input v-else-if="item.type === 'password'" type="password" v-model="localForm[item.prop]" autocomplete="off"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm">提交</el-button>
        <el-button @click="resetForm">重置</el-button>
      </el-form-item>
    </el-form>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  
  const props = defineProps({
    form: {
      type: Array,
      required: true
    },
    rules: {
      type: Object,
      required: true
    },
    onSubmit: {
      type: Function,
      required: true
    },
    onReset: {
      type: Function,
      required: true
    }
  });
  
  const emits = defineEmits(['update:form']);
  
  const ruleFormRef = ref(null);
  const localForm = ref({});
  const localRules = ref({});
  for (const item of props.form) {
    localForm.value[item.prop] = '';
  }

  Object.assign(localRules.value, props.rules);
  
  const submitForm = () => {
    if (!ruleFormRef.value) return;
    ruleFormRef.value.validate((valid) => {
      if (valid) {
        props.onSubmit(localForm.value);
      } else {
        console.log('error submit!!');
      }
    });
  };
  
  const resetForm = () => {
    if (!ruleFormRef.value) return;
    ruleFormRef.value.resetFields();
    props.onReset();
  };
  </script>