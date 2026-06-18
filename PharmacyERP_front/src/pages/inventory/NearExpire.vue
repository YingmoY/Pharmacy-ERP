<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="药品名称"><a-input v-model:value="searchParams.drug_name" allow-clear /></a-form-item>
        <a-form-item label="批号"><a-input v-model:value="searchParams.batch_number" allow-clear /></a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-alert type="warning" style="margin-bottom: 12px" message="以下药品将在 30 天内到期（法规要求），请及时处理！" banner />

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="batch_number" @change="handleTableChange" :row-class-name="getRowClass">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'drug_name'">
            <div style="font-weight: 500">{{ record.drug_name }}</div>
            <div style="font-size: 12px; color: #666">{{ record.specification }}</div>
          </template>
          <template v-if="column.key === 'expire_date'">
            <span :style="{ color: getDaysColor(record.remaining_days), fontWeight: 500 }">
              {{ record.expire_date }}
              <a-tag :color="getDaysColor(record.remaining_days)" style="margin-left: 4px">
                剩余{{ record.remaining_days }}天
              </a-tag>
            </span>
          </template>
          <template v-if="column.key === 'alert_level'">
            <a-tag :color="levelColor(record.alert_level)">{{ levelLabel(record.alert_level) }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getNearExpireInventory } from '@/api/inventory'
import type { NearExpireInventory } from '@/types/inventory'

const NEAR_EXPIRE_DAYS = 30

const loading = ref(false)
const list = ref<NearExpireInventory[]>([])
const searchParams = reactive({ drug_name: '', batch_number: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '药品名称', key: 'drug_name', ellipsis: true },
  { title: '批号', dataIndex: 'batch_number', width: 140 },
  { title: '效期', key: 'expire_date', width: 220 },
  { title: '货位', dataIndex: 'location_code', width: 100 },
  { title: '库存数量', dataIndex: 'count', width: 90 },
  { title: '预警级别', key: 'alert_level', width: 100 },
]

const levelColorMap: Record<string, string> = { HIGH: 'red', MEDIUM: 'orange', LOW: 'gold' }
const levelLabelMap: Record<string, string> = { HIGH: '高风险', MEDIUM: '中风险', LOW: '低风险' }
function levelColor(level: string): string { return levelColorMap[level] ?? 'default' }
function levelLabel(level: string): string { return levelLabelMap[level] ?? level }

function getDaysColor(days: number): string {
  if (days <= 7) return '#ff4d4f'
  if (days <= 15) return '#fa8c16'
  return '#faad14'
}

function getRowClass(record: NearExpireInventory): string {
  if (record.remaining_days <= 7) return 'near-expire-critical'
  if (record.remaining_days <= 15) return 'near-expire-warning'
  return ''
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getNearExpireInventory({ page: pagination.current, page_size: pagination.pageSize, days: NEAR_EXPIRE_DAYS })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { drug_name: '', batch_number: '' })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>

<style scoped>
:deep(.near-expire-critical) { background: #fff1f0; }
:deep(.near-expire-warning) { background: #fff7e6; }
</style>
