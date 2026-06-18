<template>
  <div class="page-content">
    <a-row :gutter="16">
      <!-- 左侧操作区 -->
      <a-col :span="10">
        <a-card title="追溯码移位" style="height: 100%">
          <a-form layout="vertical">
            <a-form-item label="① 扫描追溯码">
              <ScanInput ref="traceRef" v-model="traceCode" placeholder="扫描追溯码，回车确认" :auto-focus="true" :clear-after-submit="false" @submit="handleTraceScan" />
            </a-form-item>

            <a-form-item label="② 新目标货位" style="margin-top: 12px">
              <ScanInput v-model="newLocationCode" placeholder="扫描或输入新货位码" :auto-focus="false" :clear-after-submit="false" :disabled="!traceInfo" @submit="(v) => (newLocationCode = v)" />
            </a-form-item>

            <a-form-item label="备注">
              <a-textarea v-model:value="remark" :rows="2" placeholder="移位原因（可选）" />
            </a-form-item>

            <a-button type="primary" block size="large" :disabled="!canRelocate" :loading="loading" @click="handleRelocate">
              确认移位
            </a-button>
          </a-form>

          <a-divider />

          <!-- 结果提示 -->
          <div v-if="lastResult">
            <a-result :status="lastResult.success ? 'success' : 'error'" :title="lastResult.message" style="padding: 12px 0" />
          </div>
        </a-card>
      </a-col>

      <!-- 右侧信息区 -->
      <a-col :span="14">
        <a-card title="追溯码信息" style="margin-bottom: 16px">
          <a-empty v-if="!traceInfo" description="请先扫描追溯码" />
          <a-descriptions v-else :column="2" bordered size="small">
            <a-descriptions-item label="药品名称">{{ traceInfo.drug_name }}</a-descriptions-item>
            <a-descriptions-item label="规格">{{ traceInfo.specification }}</a-descriptions-item>
            <a-descriptions-item label="批号">{{ traceInfo.batch_number }}</a-descriptions-item>
            <a-descriptions-item label="有效期">{{ traceInfo.expire_date }}</a-descriptions-item>
            <a-descriptions-item label="当前货位">
              <a-tag color="blue">{{ traceInfo.location_code || '未上架' }}</a-tag>
            </a-descriptions-item>
            <a-descriptions-item label="状态">
              <StatusTag type="trace" :value="traceInfo.status" />
            </a-descriptions-item>
          </a-descriptions>
        </a-card>

        <a-card title="移位记录">
          <a-timeline v-if="relocateHistory.length > 0">
            <a-timeline-item v-for="(item, idx) in relocateHistory" :key="idx" :color="item.success ? 'green' : 'red'">
              <div>{{ item.trace_code }}</div>
              <div style="font-size: 12px; color: #666">{{ item.from_location }} → {{ item.to_location }}</div>
              <div style="font-size: 12px; color: #999">{{ item.message }}</div>
            </a-timeline-item>
          </a-timeline>
          <a-empty v-else description="暂无移位记录" />
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import ScanInput from '@/components/common/ScanInput.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getTraceCodeDetail, relocateTraceCode } from '@/api/inventory'
import type { DrugTraceInventory } from '@/types/inventory'

const traceCode = ref('')
const newLocationCode = ref('')
const remark = ref('')
const loading = ref(false)
const traceInfo = ref<DrugTraceInventory | null>(null)
const lastResult = ref<{ success: boolean; message: string } | null>(null)
const relocateHistory = ref<Array<{ trace_code: string; from_location: string; to_location: string; success: boolean; message: string }>>([])

const canRelocate = computed(() => traceCode.value && newLocationCode.value && traceInfo.value?.status === 'IN_STOCK')

async function handleTraceScan(code: string) {
  try {
    traceInfo.value = await getTraceCodeDetail(code)
    traceCode.value = code
    if (traceInfo.value.status !== 'IN_STOCK') {
      message.warning(`该追溯码状态为 ${traceInfo.value.status}，无法移位`)
    }
  } catch {
    traceInfo.value = null
    message.error('未找到该追溯码')
  }
}

async function handleRelocate() {
  if (!canRelocate.value) return
  loading.value = true
  try {
    await relocateTraceCode({ trace_code: traceCode.value, location_code: newLocationCode.value, remark: remark.value })
    const histEntry = {
      trace_code: traceCode.value,
      from_location: traceInfo.value?.location_code ?? '',
      to_location: newLocationCode.value,
      success: true,
      message: '移位成功',
    }
    lastResult.value = { success: true, message: '移位成功' }
    relocateHistory.value.unshift(histEntry)

    traceCode.value = ''
    newLocationCode.value = ''
    remark.value = ''
    traceInfo.value = null
  } catch (e: unknown) {
    const err = e as Error
    lastResult.value = { success: false, message: err.message || '移位失败' }
  } finally {
    loading.value = false
  }
}
</script>
