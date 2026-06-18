<template>
  <div class="page-content">
    <a-card class="search-card" size="small" v-if="task">
      <a-descriptions :column="4" size="small">
        <a-descriptions-item label="任务号">{{ task.task_no }}</a-descriptions-item>
        <a-descriptions-item label="盘点范围">{{ task.scope_type }} / {{ task.scope_value }}</a-descriptions-item>
        <a-descriptions-item label="错架数量">
          <span style="color: #fa541c; font-weight: 600">{{ pagination.total }}</span>
        </a-descriptions-item>
        <a-descriptions-item label="操作">
          <a-button size="small" @click="$router.back()">返回任务详情</a-button>
        </a-descriptions-item>
      </a-descriptions>
    </a-card>

    <a-alert type="warning" style="margin-bottom: 12px"
      message="以下追溯码扫描位置与系统记录不符（错架），请安排人员核实并移位。" show-icon banner />

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="trace_code" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'trace_code'">
            <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
          </template>
          <template v-if="column.key === 'location_diff'">
            <span>
              <a-tag color="orange">{{ record.system_location_code }}</a-tag>
              <span style="margin: 0 4px">→</span>
              <a-tag color="blue">{{ record.scanned_location_code }}</a-tag>
            </span>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="$router.push('/shelving/relocate')">去移位</a-button>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getInventoryTaskDetail, getMisplacedItems } from '@/api/inventoryTask'
import type { InventoryTask, MisplacedItem } from '@/types/inventoryTask'

const route = useRoute()
const taskId = Number(route.params.id)

const loading = ref(false)
const task = ref<InventoryTask | null>(null)
const list = ref<MisplacedItem[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 130 },
  { title: '货位差异（系统→实际）', key: 'location_diff', width: 220 },
  { title: '操作', key: 'action', width: 80 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getMisplacedItems(taskId, { page: pagination.current, page_size: pagination.pageSize })
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
