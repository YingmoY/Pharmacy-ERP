<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="货位编码"><a-input v-model:value="searchParams.location_code" allow-clear /></a-form-item>
        <a-form-item label="货位名称"><a-input v-model:value="searchParams.location_name" allow-clear /></a-form-item>
        <a-form-item label="区域"><a-input v-model:value="searchParams.area" allow-clear /></a-form-item>
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
      <a-button type="primary" @click="showDrawer()"><PlusOutlined /> 新增货位</a-button>
    </div>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="enable" :value="String(record.status)" />
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="showDrawer(record)">编辑</a-button>
              <a-popconfirm :title="record.status === 1 ? '确认停用？' : '确认启用？'" @confirm="toggleStatus(record)">
                <a-button type="link" size="small" :danger="record.status === 1">{{ record.status === 1 ? '停用' : '启用' }}</a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-drawer v-model:open="drawerVisible" :title="editing ? '编辑货位' : '新增货位'" width="480">
      <a-form :model="form" layout="vertical" ref="formRef" :rules="rules">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="货位编码" name="location_code"><a-input v-model:value="form.location_code" :disabled="!!editing" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="货位名称" name="location_name"><a-input v-model:value="form.location_name" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="区域" name="area"><a-input v-model:value="form.area" placeholder="如：A区" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="货架" name="shelf"><a-input v-model:value="form.shelf" placeholder="如：货架1" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="层" name="layer"><a-input-number v-model:value="form.layer" :min="1" style="width:100%" /></a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="位" name="position"><a-input-number v-model:value="form.position" :min="1" style="width:100%" /></a-form-item>
          </a-col>
        </a-row>
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
import { getLocationList, createLocation, updateLocation, toggleLocationStatus } from '@/api/locations'
import type { LocationInfo } from '@/types/location'

const loading = ref(false)
const saveLoading = ref(false)
const list = ref<LocationInfo[]>([])
const drawerVisible = ref(false)
const editing = ref<LocationInfo | null>(null)
const formRef = ref()

const searchParams = reactive({ location_code: '', location_name: '', area: '', status: undefined as number | undefined })
const form = reactive({ location_code: '', location_name: '', area: '', shelf: '', layer: undefined as number | undefined, position: undefined as number | undefined })
const rules = {
  location_code: [{ required: true, message: '货位编码必填' }],
  location_name: [{ required: true, message: '货位名称必填' }],
  area: [{ required: true, message: '区域必填' }],
}

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })
const columns = [
  { title: '货位编码', dataIndex: 'location_code', width: 120 },
  { title: '货位名称', dataIndex: 'location_name', ellipsis: true },
  { title: '区域', dataIndex: 'area', width: 80 },
  { title: '货架', dataIndex: 'shelf', width: 80 },
  { title: '层', dataIndex: 'layer', width: 60 },
  { title: '位', dataIndex: 'position', width: 60 },
  { title: '状态', key: 'status', width: 70 },
  { title: '操作', key: 'action', width: 120 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getLocationList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { location_code: '', location_name: '', area: '', status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

function showDrawer(loc?: LocationInfo) {
  editing.value = loc ?? null
  Object.assign(form, loc ?? { location_code: '', location_name: '', area: '', shelf: '', layer: undefined, position: undefined })
  drawerVisible.value = true
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (editing.value) {
      await updateLocation(editing.value.id, form)
    } else {
      await createLocation(form)
    }
    message.success('保存成功')
    drawerVisible.value = false
    fetchList()
  } finally { saveLoading.value = false }
}

async function toggleStatus(loc: LocationInfo) {
  await toggleLocationStatus(loc.id, loc.status === 1 ? 0 : 1)
  message.success('状态已更新')
  fetchList()
}

onMounted(fetchList)
</script>
