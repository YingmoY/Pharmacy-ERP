<template>
  <div class="page-content">
    <a-card class="search-card" size="small" v-if="task">
      <a-descriptions :column="4" size="small">
        <a-descriptions-item label="任务号">{{ task.task_no }}</a-descriptions-item>
        <a-descriptions-item label="盘点范围">{{ task.scope_type }} / {{ task.scope_value }}</a-descriptions-item>
        <a-descriptions-item label="盘亏候选数">
          <span style="color: #ff4d4f; font-weight: 600">{{ pagination.total }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="操作">
          <a-button size="small" @click="$router.back()">返回任务详情</a-button>
        </a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-alert type="warning" style="margin-bottom: 12px"
      message="以下追溯码在本次盘库范围内未被扫描到，已标记为盘亏候选。确认盘亏后将不可恢复，请谨慎操作。" show-icon banner />

    <a-card>
      <template #extra>
        <a-button type="primary" danger :disabled="selectedIds.length === 0" :loading="confirmLoading" @click="handleConfirmLoss">
          确认盘亏（{{ selectedIds.length }}条）
        </a-button>
      </template>
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="pagination"
        row-key="trace_code"
        :row-selection="{ selectedRowKeys: selectedIds, onChange: (keys: string[]) => (selectedIds = keys) }"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'trace_code'">
            <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
          </template>
          <template v-if="column.key === 'expire_date'">
            <span :style="{ color: isNearExpire(record.expire_date) ? '#ff4d4f' : 'inherit' }">{{ record.expire_date }}</span>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Modal, message } from 'ant-design-vue'
import { getInventoryTaskDetail, getLossCandidates, confirmLossCandidate } from '@/api/inventoryTask'
import type { InventoryTask, LossCandidateItem } from '@/types/inventoryTask'
import dayjs from 'dayjs'

const route = useRoute()
const taskId = Number(route.params.id)

const loading = ref(false)
const confirmLoading = ref(false)
const task = ref<InventoryTask | null>(null)
const list = ref<LossCandidateItem[]>([])
const selectedIds = ref<string[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 140 },
  { title: '批号', dataIndex: 'batch_number', width: 140 },
  { title: '有效期', key: 'expire_date', width: 120 },
  { title: '系统货位', dataIndex: 'system_location_code', width: 100 },
]

function isNearExpire(date: string): boolean {
  return dayjs(date).diff(dayjs(), 'day') <= 30
}

async function handleConfirmLoss() {
  Modal.confirm({
    title: '确认盘亏',
    content: `即将确认 ${selectedIds.value.length} 条追溯码为盘亏，此操作不可撤销！请确认无误后再操作。`,
    okText: '确认盘亏',
    okType: 'danger',
    onOk: async () => {
      confirmLoading.value = true
      try {
        // 按追溯码逐条调用确认接口
        await Promise.all(
          selectedIds.value.map((traceCode) => confirmLossCandidate(taskId, traceCode, '盘库确认盘亏'))
        )
        message.success('已确认盘亏')
        selectedIds.value = []
        fetchList()
      } finally { confirmLoading.value = false }
    },
  })
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getLossCandidates(taskId, { page: pagination.current, page_size: pagination.pageSize })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(async () => {
  task.value = await getInventoryTaskDetail(taskId)
  fetchList()
})
</script>
