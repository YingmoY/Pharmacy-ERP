<template>
  <div class="page-content">
    <div style="margin-bottom: 12px">
      <a-button type="primary" @click="showDrawer()"><PlusOutlined /> 新增角色</a-button>
    </div>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="false" row-key="id">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="showDrawer(record)">编辑</a-button>
              <a-button type="link" size="small" @click="showPermissions(record)">权限配置</a-button>
              <a-popconfirm title="确认删除此角色？" @confirm="handleDelete(record.id)">
                <a-button type="link" size="small" danger>删除</a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:open="drawerVisible" :title="editing ? '编辑角色' : '新增角色'" width="420">
      <a-form :model="form" layout="vertical" ref="formRef" :rules="rules">
        <a-form-item label="角色名称" name="name"><a-input v-model:value="form.name" /></a-form-item>
        <a-form-item label="角色编码" name="code"><a-input v-model:value="form.code" :disabled="!!editing" /></a-form-item>
        <a-form-item label="说明"><a-textarea v-model:value="form.description" :rows="3" /></a-form-item>
      </a-form>
      <template #footer>
        <a-space>
          <a-button @click="drawerVisible = false">取消</a-button>
          <a-button type="primary" :loading="saveLoading" @click="handleSave">保存</a-button>
        </a-space>
      </template>
    </a-drawer>

    <!-- 权限配置弹窗 -->
    <a-modal v-model:open="permVisible" :title="`${currentRole?.name} - 权限配置`" width="640" @ok="handleSavePerms" :confirm-loading="permLoading">
      <a-spin :spinning="permLoading">
        <a-tree
          v-model:checkedKeys="checkedPermIds"
          :tree-data="permTree"
          checkable
          :field-names="{ title: 'name', key: 'code', children: 'children' }"
          check-strictly
        />
      </a-spin>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getRoleList, createRole, updateRole, deleteRole, getRolePermissions, assignRolePermissions, getPermissionList } from '@/api/system'
import type { SysRole, SysPermission } from '@/types/system'

const loading = ref(false)
const saveLoading = ref(false)
const permLoading = ref(false)
const drawerVisible = ref(false)
const permVisible = ref(false)
const editing = ref<SysRole | null>(null)
const currentRole = ref<SysRole | null>(null)
const list = ref<SysRole[]>([])
const permTree = ref<SysPermission[]>([])
const checkedPermIds = ref<{ checked: string[]; halfChecked: string[] }>({ checked: [], halfChecked: [] })
const formRef = ref()

const form = reactive({ name: '', code: '', description: '' })
const rules = {
  name: [{ required: true, message: '角色名称必填' }],
  code: [{ required: true, message: '角色编码必填' }],
}

const columns = [
  { title: '角色名称', dataIndex: 'name', width: 150 },
  { title: '角色编码', dataIndex: 'code', width: 150 },
  { title: '说明', dataIndex: 'description', ellipsis: true },
  { title: '用户数', dataIndex: 'user_count', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 200 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getRoleList({ page: 1, page_size: 100 })
    list.value = res.list
  } finally { loading.value = false }
}

function showDrawer(role?: SysRole) {
  editing.value = role ?? null
  Object.assign(form, role ?? { name: '', code: '', description: '' })
  drawerVisible.value = true
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (editing.value) {
      await updateRole(editing.value.id, form)
    } else {
      await createRole(form)
    }
    message.success('保存成功')
    drawerVisible.value = false
    fetchList()
  } finally { saveLoading.value = false }
}

async function handleDelete(id: number) {
  await deleteRole(id)
  message.success('角色已删除')
  fetchList()
}

async function showPermissions(role: SysRole) {
  currentRole.value = role
  permLoading.value = true
  permVisible.value = true
  try {
    const [tree, assigned] = await Promise.all([getPermissionList(), getRolePermissions(role.id)])
    permTree.value = tree
    checkedPermIds.value = { checked: assigned.map((p: SysPermission) => p.code), halfChecked: [] }
  } finally { permLoading.value = false }
}

async function handleSavePerms() {
  permLoading.value = true
  try {
    await assignRolePermissions(currentRole.value!.id, checkedPermIds.value.checked)
    message.success('权限已保存')
    permVisible.value = false
  } finally { permLoading.value = false }
}

onMounted(fetchList)
</script>
