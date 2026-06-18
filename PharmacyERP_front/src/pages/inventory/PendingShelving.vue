<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="追溯码"><a-input v-model:value="searchParams.trace_code" allow-clear /></a-form-item>
        <a-form-item label="药品名称"><a-input v-model:value="searchParams.drug_name" allow-clear /></a-form-item>
        <a-form-item label="批号"><a-input v-model:value="searchParams.batch_number" allow-clear /></a-form-item>
        <a-form-item label="入库单"><a-input v-model:value="searchParams.inbound_order_no" allow-clear /></a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <template #title>
        待上架追溯码
        <a-badge :count="pagination.total" :overflow-count="9999" style="margin-left: 8px" />
      </template>
      <template #extra>
        <a-button type="primary" @click="$router.push('/shelving')">前往上架工作台</a-button>
      </template>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="trace_code" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'drug_name'">
            <div>{{ record.drug_name }}</div>
            <div style="font-size: 12px; color: #666">{{ record.specification }}</div>
          </template>
          <template v-if="column.key === 'expire_date'">
            <span :style="{ color: isNearExpire(record.expire_date) ? '#ff4d4f' : 'inherit' }">{{ record.expire_date }}</span>
          </template>
          <template v-if="column.key === 'trace_code'">
            <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="$router.push('/shelving')">去上架</a-button>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getInventoryList } from '@/api/inventory'
import type { DrugTraceInventory } from '@/types/inventory'
import dayjs from 'dayjs'

const loading = ref(false)
const list = ref<DrugTraceInventory[]>([])
const searchParams = reactive({ trace_code: '', drug_name: '', batch_number: '', inbound_order_no: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', key: 'drug_name', ellipsis: true },
  { title: '批号', dataIndex: 'batch_number', width: 140 },
  { title: '有效期', key: 'expire_date', width: 120 },
  { title: '入库时间', dataIndex: 'created_at', width: 160 },
  { title: '入库单', dataIndex: 'inbound_order_no', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

function isNearExpire(date: string): boolean {
  return dayjs(date).diff(dayjs(), 'day') <= 30
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInventoryList({ page: pagination.current, page_size: pagination.pageSize, status: 'PENDING', ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { trace_code: '', drug_name: '', batch_number: '', inbound_order_no: '' })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
