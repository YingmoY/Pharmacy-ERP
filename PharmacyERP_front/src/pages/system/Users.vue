<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="用户名"><a-input v-model:value="searchParams.username" allow-clear /></a-form-item>
        <a-form-item label="姓名"><a-input v-model:value="searchParams.real_name" allow-clear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">停用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <div style="margin-bottom: 12px">
      <a-button type="primary" @click="showDrawer()"><PlusOutlined /> 新增用户</a-button>
    </div>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'roles'">
            <a-tag v-for="role in record.roles" :key="role.id" color="blue">{{ role.name }}</a-tag>
          </template>
          <template v-if="column.key === 'status'">
            <StatusTag type="enable" :value="String(record.status)" />
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="showDrawer(record)">编辑</a-button>
              <a-button type="link" size="small" @click="showResetPwd(record)">重置密码</a-button>
              <a-popconfirm :title="record.status === 1 ? '确认停用？' : '确认启用？'" @confirm="toggleStatus(record)">
                <a-button type="link" size="small" :danger="record.status === 1">{{ record.status === 1 ? '停用' : '启用' }}</a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:open="drawerVisible" :title="editing ? '编辑用户' : '新增用户'" width="480">
      <a-form :model="form" layout="vertical" ref="formRef" :rules="rules">
        <a-form-item label="用户名" name="username"><a-input v-model:value="form.username" :disabled="!!editing" /></a-form-item>
        <a-form-item label="姓名" name="real_name"><a-input v-model:value="form.real_name" /></a-form-item>
        <a-form-item label="手机号" name="phone"><a-input v-model:value="form.phone" /></a-form-item>
        <a-form-item label="角色" name="role_codes">
          <a-select v-model:value="form.role_codes" mode="multiple" placeholder="选择角色" :options="roleOptions" />
        </a-form-item>
        <a-form-item v-if="!editing" label="初始密码" name="password">
          <a-input-password v-model:value="form.password" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-space>
          <a-button @click="drawerVisible = false">取消</a-button>
          <a-button type="primary" :loading="saveLoading" @click="handleSave">保存</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 重置密码弹窗 -->
    <a-modal v-model:open="pwdVisible" title="重置密码" @ok="handleResetPwd" :confirm-loading="pwdLoading">
      <a-form :model="pwdForm" layout="vertical">
        <a-form-item label="新密码" required>
          <a-input-password v-model:value="pwdForm.password" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getUserList, createUser, updateUser, assignUserRoles, toggleUserStatus, resetUserPassword, getRoleList } from '@/api/system'
import type { SysUser } from '@/types/system'

const loading = ref(false)
const saveLoading = ref(false)
const pwdLoading = ref(false)
const drawerVisible = ref(false)
const pwdVisible = ref(false)
const editing = ref<SysUser | null>(null)
const currentUser = ref<SysUser | null>(null)
const list = ref<SysUser[]>([])
const roleOptions = ref<{ label: string; value: string }[]>([])
const formRef = ref()

const searchParams = reactive({ username: '', real_name: '', status: undefined as number | undefined })
const form = reactive({ username: '', real_name: '', role_codes: [] as string[], password: '' })
const pwdForm = reactive({ password: '' })
const rules = {
  username: [{ required: true, message: '用户名必填' }],
  real_name: [{ required: true, message: '姓名必填' }],
  password: [{ required: true, message: '初始密码必填' }],
}

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })
const columns = [
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '姓名', dataIndex: 'real_name', width: 100 },
  { title: '手机号', dataIndex: 'phone', width: 130 },
  { title: '角色', key: 'roles', ellipsis: true },
  { title: '状态', key: 'status', width: 70 },
  { title: '最后登录', dataIndex: 'last_login_at', width: 160 },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getUserList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function showDrawer(user?: SysUser) {
  editing.value = user ?? null
  Object.assign(form, user
    ? { username: user.username, real_name: user.real_name, role_codes: user.roles?.map(r => r.code) ?? [], password: '' }
    : { username: '', real_name: '', role_codes: [], password: '' })
  drawerVisible.value = true
}

function showResetPwd(user: SysUser) {
  currentUser.value = user
  pwdForm.password = ''
  pwdVisible.value = true
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (editing.value) {
      await updateUser(editing.value.id, { real_name: form.real_name })
      await assignUserRoles(editing.value.id, form.role_codes)
    } else {
      await createUser({ username: form.username, password: form.password, real_name: form.real_name, role_codes: form.role_codes })
    }
    message.success('保存成功')
    drawerVisible.value = false
    fetchList()
  } finally { saveLoading.value = false }
}

async function handleResetPwd() {
  if (!pwdForm.password) { message.warning('请输入新密码'); return }
  pwdLoading.value = true
  try {
    await resetUserPassword(currentUser.value!.id, pwdForm.password)
    message.success('密码已重置')
    pwdVisible.value = false
  } finally { pwdLoading.value = false }
}

async function toggleStatus(user: SysUser) {
  await toggleUserStatus(user.id, user.status === 1 ? 0 : 1)
  message.success('状态已更新')
  fetchList()
}

function handleReset() {
  Object.assign(searchParams, { username: '', real_name: '', status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(async () => {
  fetchList()
  const res = await getRoleList({ page: 1, page_size: 100 })
  roleOptions.value = res.list.map(r => ({ label: r.name, value: r.code }))
})
</script>
