<template>
  <div style="background: #f5f5f5; min-height: 100vh; display: flex; flex-direction: column">
    <!-- 任务信息头 -->
    <div style="background: #1677ff; color: #fff; padding: 12px 16px">
      <div style="font-size: 15px; font-weight: 600">{{ task?.task_no }}</div>
      <div style="font-size: 12px; opacity: 0.85; margin-top: 2px">{{ task?.task_type }}</div>
    </div>

    <!-- 统计栏 -->
    <div style="display: flex; background: #fff; padding: 10px 0; border-bottom: 1px solid #f0f0f0">
      <div style="flex: 1; text-align: center; border-right: 1px solid #f0f0f0">
        <div style="font-size: 20px; font-weight: 700; color: #52c41a">{{ stats.success }}</div>
        <div style="font-size: 12px; color: #999">成功</div>
      </div>
      <div style="flex: 1; text-align: center">
        <div style="font-size: 20px; font-weight: 700; color: #ff4d4f">{{ stats.fail }}</div>
        <div style="font-size: 12px; color: #999">失败</div>
      </div>
    </div>

    <!-- 扫码区 -->
    <div style="flex: 1; padding: 16px">
      <!-- 最近扫码结果 -->
      <div v-if="lastResult" :class="lastResult.success ? 'scan-result-success' : 'scan-result-error'" style="margin-bottom: 12px; border-radius: 8px">
        <div style="font-size: 14px; font-weight: 600">{{ lastResult.trace_code }}</div>
        <div style="font-size: 13px; margin-top: 4px">{{ lastResult.message }}</div>
      </div>

      <!-- 手动输入追溯码 -->
      <div style="background: #fff; border-radius: 8px; padding: 16px">
        <div style="font-size: 14px; font-weight: 500; margin-bottom: 12px">扫描追溯码</div>
        <a-input-search
          v-model:value="inputCode"
          placeholder="输入或扫描追溯码"
          size="large"
          enter-button="确认"
          :loading="scanning"
          @search="handleScan"
          allow-clear
        />
        <a-alert v-if="!task || task.status === 'COMPLETED'" type="info" message="任务已完成，无法扫码" style="margin-top: 12px" />
      </div>

      <!-- 最近记录 -->
      <div style="margin-top: 16px">
        <div style="font-size: 14px; font-weight: 500; margin-bottom: 8px; color: #333">最近记录</div>
        <div v-for="(r, idx) in recentLogs" :key="idx" style="background: #fff; border-radius: 6px; padding: 10px; margin-bottom: 6px; display: flex; justify-content: space-between; align-items: center">
          <div>
            <div style="font-size: 13px">{{ r.trace_code }}</div>
            <div style="font-size: 12px; color: #999">{{ r.message }}</div>
          </div>
          <a-tag :color="r.success ? 'green' : 'red'">{{ r.success ? '✓' : '✗' }}</a-tag>
        </div>
        <a-empty v-if="recentLogs.length === 0" description="暂无记录" :image-style="{ height: '40px' }" />
      </div>
    </div>

    <!-- 底部完成按钮 -->
    <div style="padding: 12px 16px; background: #fff; border-top: 1px solid #f0f0f0">
      <a-button type="primary" block size="large" :disabled="task?.status !== 'IN_PROGRESS'" @click="handleComplete">
        完成任务
      </a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { getScanTaskDetail, submitScanResult, completeScanTask } from '@/api/scanTask'

const route = useRoute()
const router = useRouter()
const taskId = Number(route.params.id)

const task = ref<Record<string, unknown> | null>(null)
const inputCode = ref('')
const scanning = ref(false)
const lastResult = ref<{ trace_code: string; success: boolean; message: string } | null>(null)
const recentLogs = ref<Array<{ trace_code: string; success: boolean; message: string }>>([])
const stats = reactive({ success: 0, fail: 0 })

async function handleScan(code: string) {
  if (!code.trim()) return
  scanning.value = true
  try {
    const result = await submitScanResult(taskId, { trace_code: code })
    const entry = { trace_code: code, success: result.success, message: result.message ?? (result.success ? '扫码成功' : '扫码失败') }
    lastResult.value = entry
    recentLogs.value.unshift(entry)
    if (recentLogs.value.length > 10) recentLogs.value.length = 10
    if (result.success) stats.success++
    else stats.fail++
    inputCode.value = ''
  } finally { scanning.value = false }
}

async function handleComplete() {
  Modal.confirm({
    title: '确认完成任务',
    content: '完成后将无法继续扫码，请确认所有追溯码已扫描。',
    okText: '确认完成',
    onOk: async () => {
      await completeScanTask(taskId)
      message.success('任务已完成')
      router.push('/m/scan-tasks')
    },
  })
}

onMounted(async () => {
  task.value = await getScanTaskDetail(taskId)
})
</script>
