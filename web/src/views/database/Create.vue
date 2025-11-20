<template>
  <div class="database-create-container">
    <el-card>
      <template #header>
        <span>创建数据库</span>
      </template>
      
      <el-form :model="databaseForm" :rules="rules" ref="formRef" label-width="120px">
        <el-form-item label="数据库名" prop="name">
          <el-input v-model="databaseForm.name" placeholder="请输入数据库名" />
        </el-form-item>
        
        <el-form-item label="用户名" prop="username">
          <el-input v-model="databaseForm.username" placeholder="请输入用户名" />
        </el-form-item>
        
        <el-form-item label="密码" prop="password">
          <el-input v-model="databaseForm.password" type="password" placeholder="请输入密码" show-password />
        </el-form-item>
        
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="databaseForm.confirmPassword" type="password" placeholder="请再次输入密码" show-password />
        </el-form-item>
        
        <el-form-item label="字符集" prop="charset">
          <el-select v-model="databaseForm.charset" placeholder="选择字符集">
            <el-option label="utf8mb4" value="utf8mb4" />
            <el-option label="utf8" value="utf8" />
            <el-option label="latin1" value="latin1" />
          </el-select>
        </el-form-item>
        
        <el-form-item label="排序规则" prop="collation">
          <el-select v-model="databaseForm.collation" placeholder="选择排序规则">
            <el-option label="utf8mb4_general_ci" value="utf8mb4_general_ci" />
            <el-option label="utf8mb4_unicode_ci" value="utf8mb4_unicode_ci" />
            <el-option label="utf8_general_ci" value="utf8_general_ci" />
          </el-select>
        </el-form-item>
        
        <el-form-item>
          <el-button type="primary" @click="createDatabase" :loading="loading">
            创建数据库
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
import databaseApi from '@/api/database'

const router = useRouter()
const formRef = ref()
const loading = ref(false)

const databaseForm = reactive({
  name: '',
  username: '',
  password: '',
  confirmPassword: '',
  charset: 'utf8mb4',
  collation: 'utf8mb4_general_ci'
})

const validateConfirmPassword = (rule, value, callback) => {
  if (value !== databaseForm.password) {
    callback(new Error('两次输入的密码不一致'))
  } else {
    callback()
  }
}

const rules = {
  name: [
    { required: true, message: '请输入数据库名', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '数据库名只能包含字母、数字和下划线，且必须以字母开头', trigger: 'blur' }
  ],
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '用户名只能包含字母、数字和下划线，且必须以字母开头', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码长度不能少于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ],
  charset: [
    { required: true, message: '请选择字符集', trigger: 'change' }
  ],
  collation: [
    { required: true, message: '请选择排序规则', trigger: 'change' }
  ]
}

const createDatabase = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    loading.value = true
    
    await databaseApi.createDatabase({
      name: databaseForm.name,
      username: databaseForm.username,
      password: databaseForm.password,
      charset: databaseForm.charset,
      collation: databaseForm.collation
    })
    
    ElMessage.success('数据库创建成功')
    router.push('/database/list')
  } catch (error) {
    if (error !== false) {
      ElMessage.error('创建数据库失败')
    }
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.database-create-container {
  padding: 20px;
  max-width: 600px;
}
</style>