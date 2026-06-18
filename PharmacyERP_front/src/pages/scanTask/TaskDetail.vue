<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <template v-if="task">
        <div class="status-card" style="margin-bottom: 16px">
          <a-row align="middle" justify="space-between">
            <a-col>
              <div style="font-size: 20px; font-weight: 600">{{ task.task_no }}</div>
              <div style="margin-top: 4px"><StatusTag type="scanTask" :value="task.status" /></div>
            </a-col>
            <a-col>
              <a-button @click="$router.push('/scan-tasks')">返回列表</a-button>
            </a-col>
          </a-row>
        </div>

        <a-tabs v-model:activeKey="activeTab">
          <a-tab-pane key="info" tab="任务信息">
            <a-descriptions bordered :column="2" size="small">
              <a-descriptions-item label="任务类型">{{ task.task_type }}</a-descriptions-item>
              <a-descriptions-item label="状态"><StatusTag type="scanTask" :value="task.status" /></a-descriptions-item>
              <a-descriptions-item label="指派员工">{{ task.assigned_to_name || '-' }}</a-descriptions-item>
              <a-descriptions-item label="关联单号">{{ task.related_order_no || '-' }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ task.created_at }}</a-descriptions-item>
              <a-descriptions-item label="完成时间">{{ task.completed_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="备注" :span="2">{{ task.remark || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="items" tab="扫码明细">
            <a-table :columns="itemColumns" :data-source="items" :pagination="itemPagination" row-key="trace_code" @change="handleItemPageChange" size="small">
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'trace_code'">
                  <router-link :to="`/inventory/trace?code=${record.trace_code}`">{{ record.trace_code }}</router-link>
                </template>
                <template v-if="column.key === 'result'">
                  <a-tag :color="record.success ? 'green' : 'red'">{{ record.success ? '成功' : '失败' }}</a-tag>
                </template>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </template>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import StatusTag from '@/components/common/StatusTag.vue'
import { getScanTaskDetail, getScanTaskItems } from '@/api/scanTask'

const route = useRoute()
const taskId = Number(route.params.id)

const loading = ref(false)
const activeTab = ref('info')
const task = ref<Record<string, unknown> | null>(null)
const items = ref<Record<string, unknown>[]>([])
const itemPagination = reactive({ current: 1, pageSize: 20, total: 0 })

const itemColumns = [
  { title: '追溯码', key: 'trace_code', width: 180 },
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '操作结果', key: 'result', width: 90 },
  { title: '扫码时间', dataIndex: 'scanned_at', width: 160 },
  { title: '备注', dataIndex: 'message', ellipsis: true },
]

async function handleItemPageChange(pag: typeof itemPagination) {
  itemPagination.current = pag.current
  const res = await getScanTaskItems(taskId, { page: pag.current, page_size: pag.pageSize })
  items.value = res.list
  itemPagination.total = res.total
}

onMounted(async () => {
  loading.value = true
  try {
    [task.value] = await Promise.all([
      getScanTaskDetail(taskId),
      getScanTaskItems(taskId, { page: 1, page_size: 20 }).then(res => {
        items.value = res.list
        itemPagination.total = res.total
      }),
    ])
  } finally { loading.value = false }
})
</script>
