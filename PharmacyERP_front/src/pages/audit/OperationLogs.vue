<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="操作人"><a-input v-model:value="searchParams.operator_name" allow-clear /></a-form-item>
        <a-form-item label="操作模块"><a-input v-model:value="searchParams.module" allow-clear /></a-form-item>
        <a-form-item label="操作类型"><a-input v-model:value="searchParams.action" allow-clear /></a-form-item>
        <a-form-item label="时间">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
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
          <template v-if="column.key === 'detail'">
            <a-button type="link" size="small" @click="showDetail(record)">查看</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <a-modal v-model:open="detailVisible" title="操作详情" :footer="null" width="600">
      <a-descriptions bordered :column="1" size="small" v-if="currentLog">
        <a-descriptions-item label="操作人">{{ currentLog.operator_name }}</a-descriptions-item>
        <a-descriptions-item label="操作模块">{{ currentLog.module }}</a-descriptions-item>
        <a-descriptions-item label="操作类型">{{ currentLog.action }}</a-descriptions-item>
        <a-descriptions-item label="资源类型">{{ currentLog.resource_type }}</a-descriptions-item>
        <a-descriptions-item label="资源ID">{{ currentLog.resource_id }}</a-descriptions-item>
        <a-descriptions-item label="IP">{{ currentLog.ip }}</a-descriptions-item>
        <a-descriptions-item v-if="currentLog.before_data" label="变更前">
          <pre style="font-size: 12px; max-height: 200px; overflow: auto">{{ JSON.stringify(currentLog.before_data, null, 2) }}</pre>
        </a-descriptions-item>
        <a-descriptions-item v-if="currentLog.after_data" label="变更后">
          <pre style="font-size: 12px; max-height: 200px; overflow: auto">{{ JSON.stringify(currentLog.after_data, null, 2) }}</pre>
        </a-descriptions-item>
      </a-descriptions>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getOperationLogs } from '@/api/system'
import type { OperationLog } from '@/types/system'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const detailVisible = ref(false)
const list = ref<OperationLog[]>([])
const currentLog = ref<OperationLog | null>(null)
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ operator_name: '', module: '', action: '', start_time: '', end_time: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '操作人', dataIndex: 'operator_name', width: 100 },
  { title: '模块', dataIndex: 'module', width: 120 },
  { title: '操作', dataIndex: 'action', width: 100 },
  { title: '资源类型', dataIndex: 'resource_type', width: 120 },
  { title: 'IP', dataIndex: 'ip', width: 130 },
  { title: '时间', dataIndex: 'created_at', width: 160 },
  { title: '详情', key: 'detail', width: 70 },
]

function showDetail(log: OperationLog) {
  currentLog.value = log
  detailVisible.value = true
}

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.start_time = dates[0].format('YYYY-MM-DD')
    searchParams.end_time = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.start_time = ''
    searchParams.end_time = ''
  }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getOperationLogs({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { operator_name: '', module: '', action: '', start_time: '', end_time: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
