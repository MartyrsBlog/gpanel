<template>
  <div class="site-create-container">
    <el-card>
      <template #header>
        <span>创建网站</span>
      </template>
      
      <el-form :model="siteForm" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="域名" prop="domain">
          <el-input v-model="siteForm.domain" placeholder="请输入域名" />
        </el-form-item>
        
        <el-form-item label="网站目录" prop="path">
          <el-input v-model="siteForm.path" placeholder="请输入网站目录路径" />
        </el-form-item>
        
        <el-form-item label="PHP版本" prop="php_version">
          <el-select v-model="siteForm.php_version" placeholder="选择PHP版本">
            <el-option label="PHP 7.4" value="7.4" />
            <el-option label="PHP 8.0" value="8.0" />
            <el-option label="PHP 8.1" value="8.1" />
            <el-option label="PHP 8.2" value="8.2" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="启用SSL">
          <el-switch v-model="siteForm.ssl" />
        </el-form-item>
        
        <el-form-item label="SSL邮箱" prop="email" v-if="siteForm.ssl">
          <el-input v-model="siteForm.email" placeholder="用于SSL证书申请的邮箱" />
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="createSite" :loading="loading">
            创建网站
          </el-button>
          <el-button @click="$router.go(-1)">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import siteApi from '@/api/site'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

const siteForm = reactive({
  domain: '',
  path: '',
  php_version: '8.0',
  ssl: false,
  email: ''
})

const rules = {
  domain: [
    { required: true, message: '请输入域名', trigger: 'blur' }
  ],
  path: [
    { required: true, message: '请输入网站目录', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱地址', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ]
}

const createSite = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await siteApi.createSite(siteForm)
    ElMessage.success('网站创建成功')
    router.push('/sites/list')
  } catch (error) {
    if (error !== false) {
      ElMessage.error('创建网站失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.site-create-container {
  padding: 20px;
  max-width: 600px;
}
</style>