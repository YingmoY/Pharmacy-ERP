<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="数据表"><a-input v-model:value="searchParams.table_name" allow-clear /></a-form-item>
        <a-form-item label="记录ID"><a-input v-model:value="searchParams.record_id" allow-clear /></a-form-item>
        <a-form-item label="操作">
          <a-select v-model:value="searchParams.change_type" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option value="INSERT">新增</a-select-option>
            <a-select-option value="UPDATE">修改</a-select-option>
            <a-select-option value="DELETE">删除</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'change_type'">
            <a-tag :color="opColor(record.change_type)">{{ record.change_type }}</a-tag>
          </template>
          <template v-if="column.key === 'diff'">
            <a-button type="link" size="small" @click="showDiff(record)">查看变更</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="diffVisible" title="数据变更详情" :footer="null" width="720">
      <a-row :gutter="16" v-if="currentLog">
        <a-col :span="12">
          <div style="font-weight: 600; margin-bottom: 8px; color: #ff4d4f">变更前</div>
          <pre style="background: #fff1f0; padding: 8px; border-radius: 4px; font-size: 12px; max-height: 400px; overflow: auto">{{ JSON.stringify(currentLog.before_data, null, 2) }}</pre>
        </a-col>
        <a-col :span="12">
          <div style="font-weight: 600; margin-bottom: 8px; color: #52c41a">变更后</div>
          <pre style="background: #f6ffed; padding: 8px; border-radius: 4px; font-size: 12px; max-height: 400px; overflow: auto">{{ JSON.stringify(currentLog.after_data, null, 2) }}</pre>
        </a-col>
      </a-row>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getDataChangeLogs } from '@/api/system'
import type { DataChangeLog } from '@/types/system'

const loading = ref(false)
const diffVisible = ref(false)
const list = ref<DataChangeLog[]>([])
const currentLog = ref<DataChangeLog | null>(null)
const searchParams = reactive({ table_name: '', record_id: '', change_type: undefined as string | undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const opColorMap: Record<string, string> = { INSERT: 'green', UPDATE: 'blue', DELETE: 'red' }
function opColor(op: string): string { return opColorMap[op] ?? 'default' }

const columns = [
  { title: '操作人', dataIndex: 'operator_name', width: 100 },
  { title: '数据表', dataIndex: 'table_name', width: 130 },
  { title: '记录ID', dataIndex: 'record_id', width: 100 },
  { title: '操作类型', key: 'change_type', width: 90 },
  { title: '操作时间', dataIndex: 'created_at', width: 160 },
  { title: '变更详情', key: 'diff', width: 90 },
]

function showDiff(log: DataChangeLog) {
  currentLog.value = log
  diffVisible.value = true
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getDataChangeLogs({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { table_name: '', record_id: '', change_type: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
