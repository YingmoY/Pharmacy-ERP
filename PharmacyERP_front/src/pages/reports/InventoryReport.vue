<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="药品名称"><a-input v-model:value="searchParams.drug_name" allow-clear /></a-form-item>
        <a-form-item label="货位"><a-input v-model:value="searchParams.location_code" allow-clear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 110px" allow-clear>
            <a-select-option value="PENDING">待上架</a-select-option>
            <a-select-option value="IN_STOCK">在库</a-select-option>
            <a-select-option value="SOLD">已售出</a-select-option>
            <a-select-option value="LOST">盘亏</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button style="margin-left: 8px" @click="handleExport" :loading="exportLoading">导出</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-row :gutter="16" style="margin-bottom: 16px" v-if="summary">
      <a-col :span="6">
        <a-card size="small"><a-statistic title="在库总数" :value="summary.in_stock_count" :value-style="{ color: '#1677ff' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="待上架" :value="summary.pending_count" :value-style="{ color: '#fa8c16' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="近效期" :value="summary.near_expire_count" :value-style="{ color: '#ff4d4f' }" /></a-card>
      </a-col>
      <a-col :span="6">
        <a-card size="small"><a-statistic title="盘亏数" :value="summary.loss_count" :value-style="{ color: '#d9363e' }" /></a-card>
      </a-col>
    </a-row>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="trace_code" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="trace" :value="record.status" />
          </template>
          <template v-if="column.key === 'expire_date'">
            <span :style="{ color: isNearExpire(record.expire_date) ? '#ff4d4f' : 'inherit' }">{{ record.expire_date }}</span>
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
import StatusTag from '@/components/common/StatusTag.vue'
import { getInventoryReport, exportInventoryReport } from '@/api/report'
import dayjs from 'dayjs'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const exportLoading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const summary = ref<Record<string, number> | null>(null)
const searchParams = reactive({ drug_name: '', location_code: '', status: undefined as string | undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 130 },
  { title: '批号', dataIndex: 'batch_number', width: 130 },
  { title: '有效期', key: 'expire_date', width: 120 },
  { title: '货位', dataIndex: 'location_code', width: 100 },
  { title: '状态', key: 'status', width: 90 },
]

function isNearExpire(date: string): boolean {
  return dayjs(date).diff(dayjs(), 'day') <= 30
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInventoryReport(searchParams)
    list.value = res.items
    pagination.total = res.items.length
    summary.value = { total_count: res.total_count }
  } finally { loading.value = false }
}

async function handleExport() {
  exportLoading.value = true
  try {
    await exportInventoryReport(searchParams as Record<string, unknown>)
    message.success('导出任务已创建，请前往「导出任务」页面下载')
  } finally { exportLoading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { drug_name: '', location_code: '', status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
