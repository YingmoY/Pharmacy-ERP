<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="追溯码"><a-input v-model:value="searchParams.trace_code" allow-clear /></a-form-item>
        <a-form-item label="操作类型">
          <a-select v-model:value="searchParams.action_type" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="INBOUND">入库</a-select-option>
            <a-select-option value="SHELVE">上架</a-select-option>
            <a-select-option value="SELL">销售</a-select-option>
            <a-select-option value="RELOCATE">移位</a-select-option>
            <a-select-option value="LOSS">盘亏</a-select-option>
            <a-select-option value="REFUND">退款</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间范围">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button style="margin-left: 8px" @click="handleExport" :loading="exportLoading">导出</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'action_type'">
            <a-tag :color="actionColor(record.action_type)">{{ record.action_type }}</a-tag>
          </template>
          <template v-if="column.key === 'trace_code'">
            <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getTraceLogReport, exportTraceLogReport } from '@/api/report'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const exportLoading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ trace_code: '', action_type: undefined as string | undefined, start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const actionColorMap: Record<string, string> = {
  INBOUND: 'green',
  SHELVE: 'blue',
  SELL: 'purple',
  RELOCATE: 'cyan',
  LOSS: 'red',
  REFUND: 'orange',
}

function actionColor(type: string): string { return actionColorMap[type] ?? 'default' }

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '操作类型', key: 'action_type', width: 90 },
  { title: '操作前状态', dataIndex: 'from_status', width: 100 },
  { title: '操作后状态', dataIndex: 'to_status', width: 100 },
  { title: '操作人', dataIndex: 'operator_name', width: 90 },
  { title: '操作时间', dataIndex: 'created_at', width: 160 },
  { title: '备注', dataIndex: 'remark', ellipsis: true },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.start_date = dates[0].format('YYYY-MM-DD')
    searchParams.end_date = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.start_date = ''
    searchParams.end_date = ''
  }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getTraceLogReport({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

async function handleExport() {
  exportLoading.value = true
  try {
    await exportTraceLogReport(searchParams as Record<string, unknown>)
    message.success('导出任务已创建，请前往「导出任务」页面下载')
  } finally { exportLoading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { trace_code: '', action_type: undefined, start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
