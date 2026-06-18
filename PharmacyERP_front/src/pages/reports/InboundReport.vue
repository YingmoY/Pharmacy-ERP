<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="时间范围">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
        </a-form-item>
        <a-form-item label="供应商"><a-input v-model:value="searchParams.supplier_name" allow-clear /></a-form-item>
        <a-form-item label="药品名称"><a-input v-model:value="searchParams.drug_name" allow-clear /></a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button style="margin-left: 8px" @click="handleExport" :loading="exportLoading">导出</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-row :gutter="16" style="margin-bottom: 16px" v-if="summary">
      <a-col :span="6">
        <a-card size="small"><a-statistic title="入库总金额" :value="summary.total_amount" :precision="2" prefix="¥" :value-style="{ color: '#52c41a' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="入库单数" :value="summary.total_orders" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="入库品种" :value="summary.total_drugs" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="追溯码总数" :value="summary.total_trace_codes" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'amount'">
            <span style="color: #52c41a; font-weight: 500">¥{{ record.amount }}</span>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getInboundReport, exportInboundReport } from '@/api/report'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const exportLoading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const summary = ref<Record<string, number> | null>(null)
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ supplier_name: '', drug_name: '', start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '入库单号', dataIndex: 'order_no', width: 180 },
  { title: '供应商', dataIndex: 'supplier_name', ellipsis: true },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '批号', dataIndex: 'batch_number', width: 130 },
  { title: '数量', dataIndex: 'qty', width: 70 },
  { title: '金额', key: 'amount', width: 110 },
  { title: '入库时间', dataIndex: 'created_at', width: 160 },
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
    const res = await getInboundReport(searchParams)
    list.value = res.items
    pagination.total = res.items.length
    summary.value = { total_amount: res.total_amount, total_orders: res.total_orders }
  } finally { loading.value = false }
}

async function handleExport() {
  exportLoading.value = true
  try {
    await exportInboundReport(searchParams as Record<string, unknown>)
    message.success('导出任务已创建，请前往「导出任务」页面下载')
  } finally { exportLoading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { supplier_name: '', drug_name: '', start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
