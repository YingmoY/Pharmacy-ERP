<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="追溯码"><a-input v-model:value="searchParams.trace_code" allow-clear /></a-form-item>
        <a-form-item label="操作类型">
          <a-select v-model:value="searchParams.adjustment_type" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="LOSS">盘亏</a-select-option>
            <a-select-option value="RETURN">退库</a-select-option>
            <a-select-option value="RELOCATE">移位</a-select-option>
            <a-select-option value="MANUAL">手动调整</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="操作时间">
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
          <template v-if="column.key === 'adjustment_type'">
            <a-tag :color="getTypeColor(record.adjust_type)">{{ getTypeLabel(record.adjust_type) }}</a-tag>
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
import { getInventoryAdjustments } from '@/api/inventory'
import type { InventoryAdjustment } from '@/types/inventory'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const list = ref<InventoryAdjustment[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ trace_code: '', adjustment_type: undefined as string | undefined, start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const typeMap: Record<string, { label: string; color: string }> = {
  LOSS: { label: '盘亏', color: 'red' },
  RETURN: { label: '退库', color: 'orange' },
  RELOCATE: { label: '移位', color: 'blue' },
  MANUAL: { label: '手动调整', color: 'purple' },
}

function getTypeLabel(type: string) { return typeMap[type]?.label ?? type }
function getTypeColor(type: string) { return typeMap[type]?.color ?? 'default' }

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.start_date = dates[0].format('YYYY-MM-DD')
    searchParams.end_date = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.start_date = ''
    searchParams.end_date = ''
  }
}

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '操作类型', key: 'adjustment_type', width: 100 },
  { title: '原货位', dataIndex: 'from_location_code', width: 100 },
  { title: '目标货位', dataIndex: 'to_location_code', width: 100 },
  { title: '备注', dataIndex: 'reason', ellipsis: true },
  { title: '操作人', dataIndex: 'operator_name', width: 90 },
  { title: '操作时间', dataIndex: 'created_at', width: 160 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getInventoryAdjustments({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { trace_code: '', adjustment_type: undefined, start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
