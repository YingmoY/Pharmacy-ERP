<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="追溯码">
          <a-input v-model:value="searchParams.trace_code" allow-clear style="width: 200px" />
        </a-form-item>
        <a-form-item label="药品">
          <a-input v-model:value="searchParams.drug_name" allow-clear />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="PENDING">待上架</a-select-option>
            <a-select-option value="IN_STOCK">在库</a-select-option>
            <a-select-option value="SOLD">已售</a-select-option>
            <a-select-option value="MISPLACED">错架</a-select-option>
            <a-select-option value="LOSS_CANDIDATE">盘亏候选</a-select-option>
            <a-select-option value="LOST">已盘亏</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="货位">
          <a-input v-model:value="searchParams.location_code" allow-clear style="width: 120px" />
        </a-form-item>
        <a-form-item label="批号">
          <a-input v-model:value="searchParams.batch_number" allow-clear style="width: 120px" />
        </a-form-item>
        <a-form-item>
          <a-checkbox v-model:checked="searchParams.near_expire">近效期</a-checkbox>
          <a-checkbox v-model:checked="searchParams.is_reserved" style="margin-left: 8px">已预占</a-checkbox>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

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
          <template v-if="column.key === 'trace_code'">
            <a-button type="link" size="small" @click="router.push(`/trace/${record.trace_code}`)">
              {{ record.trace_code }}
            </a-button>
          </template>
          <template v-if="column.key === 'status'">
            <StatusTag type="trace" :value="record.status" />
            <a-tag v-if="record.is_reserved" color="orange" size="small" style="margin-left: 2px">预占</a-tag>
          </template>
          <template v-if="column.key === 'expire_date'">
            <span :style="{ color: isExpiringSoon(record.expire_date) ? '#ff4d4f' : 'inherit' }">
              {{ record.expire_date }}
            </span>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInventoryList } from '@/api/inventory'
import type { DrugTraceInventory } from '@/types/inventory'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const list = ref<DrugTraceInventory[]>([])

const searchParams = reactive({
  trace_code: '',
  drug_name: '',
  status: undefined as string | undefined,
  location_code: '',
  batch_number: '',
  near_expire: false,
  is_reserved: false,
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '追溯码', key: 'trace_code', width: 220 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 120 },
  { title: '批号', dataIndex: 'batch_number', width: 120 },
  { title: '有效期', key: 'expire_date', width: 100 },
  { title: '状态', key: 'status', width: 120 },
  { title: '货位', dataIndex: 'location_code', width: 90 },
  { title: '入库单', dataIndex: 'inbound_order_no', width: 160, ellipsis: true },
  { title: '更新时间', dataIndex: 'updated_at', width: 160 },
]

function isExpiringSoon(date: string): boolean {
  return dayjs(date).diff(dayjs(), 'day') <= 30
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInventoryList({
      page: pagination.current,
      page_size: pagination.pageSize,
      trace_code: searchParams.trace_code || undefined,
      status: searchParams.status as any,
      batch_number: searchParams.batch_number || undefined,
      near_expire: searchParams.near_expire || undefined,
      is_reserved: searchParams.is_reserved || undefined,
    })
    list.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

function handleReset() {
  Object.assign(searchParams, { trace_code: '', drug_name: '', status: undefined, location_code: '', batch_number: '', near_expire: false, is_reserved: false })
  pagination.current = 1
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
