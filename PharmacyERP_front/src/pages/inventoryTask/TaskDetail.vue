<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <template v-if="task">
        <!-- 状态卡片 -->
        <div class="status-card" style="margin-bottom: 16px">
          <a-row align="middle" justify="space-between">
            <a-col>
              <div style="font-size: 20px; font-weight: 600">{{ task.task_no }}</div>
              <div style="margin-top: 4px"><StatusTag type="inventoryTask" :value="task.status" /></div>
            </a-col>
            <a-col>
              <a-space>
                <a-button v-if="task.status === 'PENDING'" type="primary" @click="startTask">开始盘库</a-button>
                <a-button v-if="task.status === 'IN_PROGRESS'" type="primary" @click="$router.push(`/inventory-tasks/${task.id}/scan`)">继续扫码</a-button>
                <a-popconfirm v-if="['PENDING', 'IN_PROGRESS'].includes(task.status)" title="确认取消此盘库任务？" @confirm="cancelTask">
                  <a-button danger>取消任务</a-button>
                </a-popconfirm>
              </a-space>
            </a-col>
          </a-row>
          <a-row :gutter="24" style="margin-top: 16px">
            <a-col :span="4">
              <a-statistic title="计划盘点" :value="task.total_count" />
            </a-col>
            <a-col :span="4">
              <a-statistic title="已扫码" :value="task.scanned_count" :value-style="{ color: '#1677ff' }" />
            </a-col>
            <a-col :span="4">
              <a-statistic title="正常" :value="task.normal_count" :value-style="{ color: '#52c41a' }" />
            </a-col>
            <a-col :span="4">
              <a-statistic title="错架" :value="task.misplaced_count" :value-style="{ color: '#fa541c' }" />
            </a-col>
            <a-col :span="4">
              <a-statistic title="盘亏候选" :value="task.loss_candidate_count" :value-style="{ color: '#ff4d4f' }" />
            </a-col>
          </a-row>
        </div>

        <a-tabs v-model:activeKey="activeTab">
          <a-tab-pane key="info" tab="基本信息">
            <a-descriptions bordered :column="2" size="small">
              <a-descriptions-item label="盘点范围类型">{{ task.scope_type }}</a-descriptions-item>
              <a-descriptions-item label="范围值">{{ task.scope_value }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ task.created_at }}</a-descriptions-item>
              <a-descriptions-item label="完成时间">{{ task.completed_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="备注" :span="2">{{ task.remark || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="details" tab="盘点明细">
            <a-table :columns="detailColumns" :data-source="details" :loading="detailLoading" :pagination="detailPagination" row-key="trace_code" @change="handleDetailPageChange">
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'scan_result'">
                  <StatusTag type="scanResult" :value="record.scan_result" />
                </template>
                <template v-if="column.key === 'trace_code'">
                  <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
                </template>
              </template>
            </a-table>
          </a-tab-pane>

          <a-tab-pane key="loss" tab="盘亏候选">
            <a-button type="primary" style="margin-bottom: 12px" @click="$router.push(`/inventory-tasks/${task.id}/loss-candidates`)">查看盘亏候选明细</a-button>
          </a-tab-pane>

          <a-tab-pane key="misplaced" tab="错架明细">
            <a-button type="primary" style="margin-bottom: 12px" @click="$router.push(`/inventory-tasks/${task.id}/misplaced`)">查看错架明细</a-button>
          </a-tab-pane>
        </a-tabs>
      </template>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInventoryTaskDetail, getInventoryTaskDetails, startInventoryTask, cancelInventoryTask } from '@/api/inventoryTask'
import type { InventoryTask, InventoryTaskDetail } from '@/types/inventoryTask'

const route = useRoute()
const router = useRouter()
const taskId = Number(route.params.id)

const loading = ref(false)
const detailLoading = ref(false)
const activeTab = ref('info')
const task = ref<InventoryTask | null>(null)
const details = ref<InventoryTaskDetail[]>([])
const detailPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const detailColumns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '扫描货位', dataIndex: 'scanned_location_code', width: 100 },
  { title: '系统货位', dataIndex: 'system_location_code', width: 100 },
  { title: '扫描结果', key: 'scan_result', width: 100 },
  { title: '扫描时间', dataIndex: 'scanned_at', width: 160 },
]

async function fetchTask() {
  loading.value = true
  try {
    task.value = await getInventoryTaskDetail(taskId)
  } finally { loading.value = false }
}

async function fetchDetails() {
  detailLoading.value = true
  try {
    const res = await getInventoryTaskDetails(taskId, { page: detailPagination.current, page_size: detailPagination.pageSize })
    details.value = res.list
    detailPagination.total = res.total
  } finally { detailLoading.value = false }
}

function handleDetailPageChange(pag: typeof detailPagination) {
  detailPagination.current = pag.current
  fetchDetails()
}

async function startTask() {
  await startInventoryTask(taskId)
  message.success('任务已开始')
  fetchTask()
}

async function cancelTask() {
  await cancelInventoryTask(taskId)
  message.success('任务已取消')
  router.push('/inventory-tasks')
}

onMounted(() => {
  fetchTask()
  fetchDetails()
})
</script>
