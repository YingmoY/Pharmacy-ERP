<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="供应商编码">
          <a-input v-model:value="searchParams.supplier_code" allow-clear />
        </a-form-item>
        <a-form-item label="供应商名称">
          <a-input v-model:value="searchParams.name" allow-clear />
        </a-form-item>
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
      <a-button type="primary" @click="showDrawer()"><PlusOutlined /> 新增供应商</a-button>
    </div>

    <a-card>
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="enable" :value="String(record.status)" />
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="showDrawer(record)">编辑</a-button>
              <a-popconfirm
                :title="record.status === 1 ? '确认停用？' : '确认启用？'"
                @confirm="toggleStatus(record)"
              >
                <a-button type="link" size="small" :danger="record.status === 1">
                  {{ record.status === 1 ? '停用' : '启用' }}
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:open="drawerVisible" :title="editing ? '编辑供应商' : '新增供应商'" width="480">
      <a-form :model="form" layout="vertical" ref="formRef" :rules="rules">
        <a-form-item label="供应商编码" name="supplier_code">
          <a-input v-model:value="form.supplier_code" :disabled="!!editing" />
        </a-form-item>
        <a-form-item label="供应商名称" name="name">
          <a-input v-model:value="form.name" />
        </a-form-item>
        <a-form-item label="联系人" name="contact_name">
          <a-input v-model:value="form.contact_name" />
        </a-form-item>
        <a-form-item label="电话" name="contact_phone">
          <a-input v-model:value="form.contact_phone" />
        </a-form-item>
        <a-form-item label="地址" name="address">
          <a-textarea v-model:value="form.address" :rows="2" />
        </a-form-item>
      </a-form>
      <template #footer>
        <a-space>
          <a-button @click="drawerVisible = false">取消</a-button>
          <a-button type="primary" :loading="saveLoading" @click="handleSave">保存</a-button>
        </a-space>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getSupplierList, createSupplier, updateSupplier, toggleSupplierStatus } from '@/api/suppliers'
import type { Supplier } from '@/types/supplier'

const loading = ref(false)
const saveLoading = ref(false)
const list = ref<Supplier[]>([])
const drawerVisible = ref(false)
const editing = ref<Supplier | null>(null)
const formRef = ref()

const searchParams = reactive({ supplier_code: '', name: '', status: undefined as number | undefined })
const form = reactive({ supplier_code: '', name: '', contact_name: '', contact_phone: '', address: '' })
const rules = {
  supplier_code: [{ required: true, message: '供应商编码必填' }],
  name: [{ required: true, message: '供应商名称必填' }],
}

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })
const columns = [
  { title: '供应商编码', dataIndex: 'supplier_code', width: 140 },
  { title: '名称', dataIndex: 'name', ellipsis: true },
  { title: '联系人', dataIndex: 'contact_name', width: 90 },
  { title: '电话', dataIndex: 'contact_phone', width: 130 },
  { title: '状态', key: 'status', width: 70 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 120 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getSupplierList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { supplier_code: '', name: '', status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

function showDrawer(supplier?: Supplier) {
  editing.value = supplier ?? null
  Object.assign(form, supplier ?? { supplier_code: '', name: '', contact_name: '', contact_phone: '', address: '' })
  drawerVisible.value = true
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (editing.value) {
      await updateSupplier(editing.value.id, form)
    } else {
      await createSupplier(form)
    }
    message.success('保存成功')
    drawerVisible.value = false
    fetchList()
  } finally { saveLoading.value = false }
}

async function toggleStatus(supplier: Supplier) {
  await toggleSupplierStatus(supplier.id, supplier.status === 1 ? 0 : 1)
  message.success('状态已更新')
  fetchList()
}

onMounted(fetchList)
</script>
