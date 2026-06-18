<template>
  <div class="page-content">
    <!-- 任务信息 -->
    <a-card size="small" style="margin-bottom: 16px" v-if="task">
      <a-descriptions :column="4" size="small">
        <a-descriptions-item label="任务号">{{ task.task_no }}</a-descriptions-item>
        <a-descriptions-item label="盘点范围">{{ task.scope_type }} / {{ task.scope_value }}</a-descriptions-item>
        <a-descriptions-item label="状态"><StatusTag type="inventoryTask" :value="task.status" /></a-descriptions-item>
        <a-descriptions-item label="操作">
          <a-space>
            <a-button
              type="primary"
              size="small"
              :disabled="task.status !== 'IN_PROGRESS'"
              @click="handleComplete"
            >
              完成盘库
            </a-button>
          </a-space>
        </a-descriptions-item>
      </a-descriptions>
    </a-card>

    <div class="workbench-layout">
      <!-- 中间：扫码区 -->
      <div class="scan-area">
        <a-card title="盘库扫码" size="small" style="height: 100%">
          <a-form-item label="当前货位">
            <ScanInput
              v-model="currentLocation"
              placeholder="扫描货位码"
              :clear-after-submit="false"
              :auto-focus="false"
              @submit="(v) => (currentLocation = v)"
            />
          </a-form-item>
          <a-form-item label="追溯码" style="margin-top: 16px">
            <ScanInput
              :disabled="!currentLocation || task?.status !== 'IN_PROGRESS'"
              placeholder="扫描追溯码，回车确认"
              :auto-focus="true"
              @submit="handleScan"
            />
          </a-form-item>

          <a-alert
            v-if="!currentLocation"
            type="info"
            message="请先设置当前货位，再开始扫码"
            style="margin-top: 12px"
          />

          <!-- 统计 -->
          <a-row :gutter="16" style="margin-top: 16px">
            <a-col :span="6">
              <a-statistic title="正常" :value="stats.normal" :value-style="{ color: '#52c41a' }" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="错架" :value="stats.misplaced" :value-style="{ color: '#fa541c' }" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="异常" :value="stats.unexpected" :value-style="{ color: '#ff4d4f' }" />
            </a-col>
            <a-col :span="6">
              <a-statistic title="重复" :value="stats.duplicate" :value-style="{ color: '#faad14' }" />
            </a-col>
          </a-row>
        </a-card>
      </div>

      <!-- 右侧：扫码结果 -->
      <div class="result-area">
        <a-card title="扫码记录" size="small" style="height: 100%; overflow: auto">
          <div v-for="(result, idx) in scanResults" :key="idx" :class="getScanResultClass(result.scan_result)">
            <div style="font-size: 13px; font-weight: 500">{{ result.trace_code }}</div>
            <div style="font-size: 12px">
              <StatusTag type="scanResult" :value="result.scan_result" />
              {{ result.message }}
            </div>
          </div>
          <a-empty v-if="scanResults.length === 0" description="暂无扫码记录" :image-style="{ height: '40px' }" />
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import ScanInput from '@/components/common/ScanInput.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInventoryTaskDetail, scanInventory, completeInventoryTask } from '@/api/inventoryTask'
import type { InventoryTask, InventoryScanResult } from '@/types/inventoryTask'

const route = useRoute()
const router = useRouter()
const taskId = Number(route.params.id)

const task = ref<InventoryTask | null>(null)
const currentLocation = ref('')
const scanResults = ref<Array<InventoryScanResult & { trace_code: string }>>([])
const stats = reactive({ normal: 0, misplaced: 0, unexpected: 0, duplicate: 0 })

function getScanResultClass(result: string): string {
  const map: Record<string, string> = {
    NORMAL: 'scan-result-success',
    MISPLACED_FOUND: 'scan-result-warning',
    UNEXPECTED: 'scan-result-error',
    DUPLICATE: 'scan-result-warning',
  }
  return map[result] ?? 'scan-result-warning'
}

async function handleScan(traceCode: string) {
  try {
    const result = await scanInventory(taskId, {
      trace_code: traceCode,
      scanned_location_code: currentLocation.value,
    })
    const entry = { ...result, trace_code: traceCode }
    scanResults.value.unshift(entry)
    if (scanResults.value.length > 100) scanResults.value.length = 100

    // 更新统计
    const key = result.scan_result.toLowerCase().replace('_found', '') as keyof typeof stats
    if (key in stats) stats[key]++
  } catch {
    scanResults.value.unshift({
      success: false,
      scan_result: 'UNEXPECTED',
      trace_code: traceCode,
      scanned_location_code: currentLocation.value,
      message: '扫码请求失败',
    })
  }
}

async function handleComplete() {
  Modal.confirm({
    title: '确认完成盘库',
    content: '完成后将自动标记范围内未扫到的追溯码为盘亏候选，此操作不可撤销！',
    okText: '确认完成',
    okType: 'danger',
    onOk: async () => {
      await completeInventoryTask(taskId)
      message.success('盘库已完成')
      router.push(`/inventory-tasks/${taskId}`)
    },
  })
}

onMounted(async () => {
  task.value = await getInventoryTaskDetail(taskId)
})
</script>

<style scoped>
.workbench-layout {
  display: grid;
  grid-template-columns: 420px 1fr;
  gap: 16px;
  height: calc(100vh - 220px);
}
</style>
