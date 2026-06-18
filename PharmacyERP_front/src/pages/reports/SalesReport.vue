<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="时间范围">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
        </a-form-item>
        <a-form-item label="药品名称"><a-input v-model:value="searchParams.drug_name" allow-clear /></a-form-item>
        <a-form-item label="收银员"><a-input v-model:value="searchParams.cashier_name" allow-clear /></a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button style="margin-left: 8px" @click="handleExport" :loading="exportLoading">导出</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <!-- 汇总统计 -->
    <a-row :gutter="16" style="margin-bottom: 16px" v-if="summary">
      <a-col :span="6">
        <a-card size="small"><a-statistic title="销售总额" :value="summary.total_amount" :precision="2" prefix="¥" :value-style="{ color: '#1677ff' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="订单总数" :value="summary.total_orders" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="销售品种数" :value="summary.total_drugs" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="退款金额" :value="summary.total_refund" :precision="2" prefix="¥" :value-style="{ color: '#ff4d4f' }" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'amount'">
            <span style="color: #1677ff; font-weight: 500">¥{{ record.amount }}</span>
          </template>
        </template>
        <template #summary>
          <a-table-summary fixed>
            <a-table-summary-row>
              <a-table-summary-cell :index="0" :col-span="4">本页合计</a-table-summary-cell>
              <a-table-summary-cell :index="4">
                <span style="color: #1677ff; font-weight: 500">¥{{ pageTotal }}</span>
              </a-table-summary-cell>
            </a-table-summary-row>
          </a-table-summary>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { getSalesReport, exportSalesReport } from '@/api/report'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const exportLoading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const summary = ref<Record<string, number> | null>(null)
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ drug_name: '', cashier_name: '', start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const pageTotal = computed(() => list.value.reduce((s, r) => s + Number(r.amount ?? 0), 0).toFixed(2))

const columns = [
  { title: '订单号', dataIndex: 'order_no', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 130 },
  { title: '数量', dataIndex: 'qty', width: 70 },
  { title: '金额', key: 'amount', width: 110 },
  { title: '收银员', dataIndex: 'cashier_name', width: 90 },
  { title: '销售时间', dataIndex: 'created_at', width: 160 },
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
    const res = await getSalesReport(searchParams)
    list.value = res.items
    pagination.total = res.items.length
    summary.value = { total_amount: res.total_amount, total_orders: res.total_orders }
  } finally { loading.value = false }
}

async function handleExport() {
  exportLoading.value = true
  try {
    await exportSalesReport(searchParams as Record<string, unknown>)
    message.success('导出任务已创建，请前往「导出任务」页面下载')
  } finally { exportLoading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { drug_name: '', cashier_name: '', start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
