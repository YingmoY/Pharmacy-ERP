<template>
  <div class="page-content">
    <div class="workbench-layout">
      <!-- 左侧：扫码操作区 -->
      <div class="shelf-scan-panel">
        <a-card title="上架扫码工作台" size="small" style="height: 100%">
          <!-- 扫描追溯码 -->
          <a-form-item label="① 扫描追溯码">
            <ScanInput
              ref="traceScanRef"
              v-model="traceCode"
              placeholder="扫描或输入追溯码，回车确认"
              :clear-after-submit="false"
              @submit="handleTraceScan"
              :auto-focus="true"
            />
          </a-form-item>

          <!-- 目标货位 -->
          <a-form-item label="② 输入/扫描目标货位" style="margin-top: 16px">
            <ScanInput
              ref="locationScanRef"
              v-model="locationCode"
              placeholder="扫描货位码，回车确认"
              :clear-after-submit="false"
              :auto-focus="false"
              @submit="handleLocationScan"
            />
          </a-form-item>

          <!-- 混放提示 -->
          <a-alert
            v-if="mixCheckResult?.has_mixed_drugs"
            type="warning"
            :message="'当前货位已有其他药品，请注意混放情况'"
            style="margin-bottom: 12px"
            banner
          />

          <!-- 执行上架 -->
          <a-button
            type="primary"
            block
            size="large"
            :disabled="!canShelve"
            :loading="shelveLoading"
            style="margin-top: 16px"
            @click="handleShelve"
          >
            确认上架
          </a-button>

          <!-- 结果提示 -->
          <div v-if="lastResult" style="margin-top: 16px">
            <a-alert
              :type="lastResult.success ? 'success' : 'error'"
              :message="lastResult.message"
              banner
            />
          </div>
        </a-card>
      </div>

      <!-- 右侧：药品和货位信息 -->
      <div class="shelf-info-panel">
        <!-- 追溯码信息 -->
        <a-card title="追溯码信息" size="small" style="margin-bottom: 16px">
          <a-empty v-if="!traceInfo" description="请先扫描追溯码" />
          <a-descriptions v-else :column="1" size="small">
            <a-descriptions-item label="药品名称">{{ traceInfo.drug_name }}</a-descriptions-item>
            <a-descriptions-item label="规格">{{ traceInfo.specification }}</a-descriptions-item>
            <a-descriptions-item label="批号">{{ traceInfo.batch_number }}</a-descriptions-item>
            <a-descriptions-item label="有效期">
              <span :style="{ color: isNearExpire(traceInfo.expire_date) ? '#ff4d4f' : 'inherit' }">
                {{ traceInfo.expire_date }}
              </span>
            </a-descriptions-item>
            <a-descriptions-item label="当前状态">
              <StatusTag type="trace" :value="traceInfo.status" />
            </a-descriptions-item>
            <a-descriptions-item label="入库单">{{ traceInfo.inbound_order_no }}</a-descriptions-item>
          </a-descriptions>
        </a-card>

        <!-- 目标货位信息 -->
        <a-card title="目标货位" size="small">
          <a-empty v-if="!locationCode" description="请输入目标货位" />
          <a-descriptions v-else :column="1" size="small">
            <a-descriptions-item label="货位编码">{{ locationCode }}</a-descriptions-item>
            <template v-if="mixCheckResult">
              <a-descriptions-item label="已存药品">{{ mixCheckResult.drugs?.length ?? 0 }} 种</a-descriptions-item>
              <a-descriptions-item v-for="drug in (mixCheckResult.drugs ?? [])" :key="drug.drug_id" :label="drug.drug_name">
                {{ drug.count }} 盒
              </a-descriptions-item>
            </template>
          </a-descriptions>
        </a-card>
      </div>

      <!-- 最近上架记录 -->
      <div class="shelf-history-panel">
        <a-card title="最近上架记录" size="small" style="height: 100%; overflow: auto">
          <div v-for="(record, idx) in recentResults" :key="idx" :class="record.success ? 'scan-result-success' : 'scan-result-error'">
            <div style="font-weight: 500; font-size: 13px">{{ record.trace_code }}</div>
            <div style="font-size: 12px">{{ record.drug_name || '' }} → {{ record.location_code }}</div>
            <div style="font-size: 11px; color: #666">{{ record.message }}</div>
          </div>
          <a-empty v-if="recentResults.length === 0" description="暂无记录" :image-style="{ height: '40px' }" />
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { message } from 'ant-design-vue'
import ScanInput from '@/components/common/ScanInput.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getTraceCodeDetail, shelveTraceCode } from '@/api/inventory'
import { checkMixPlacement } from '@/api/locations'
import type { DrugTraceInventory, ShelvingResult, MixCheckResult } from '@/types/inventory'
import type { MixCheckResult as MixResult } from '@/types/location'
import dayjs from 'dayjs'

const traceCode = ref('')
const locationCode = ref('')
const traceInfo = ref<DrugTraceInventory | null>(null)
const mixCheckResult = ref<MixResult | null>(null)
const shelveLoading = ref(false)
const lastResult = ref<ShelvingResult | null>(null)
const recentResults = ref<ShelvingResult[]>([])

const canShelve = computed(() =>
  traceCode.value.trim() !== '' &&
  locationCode.value.trim() !== '' &&
  traceInfo.value?.status === 'PENDING',
)

async function handleTraceScan(code: string) {
  try {
    traceInfo.value = await getTraceCodeDetail(code)
    traceCode.value = code
    if (traceInfo.value.status !== 'PENDING') {
      message.warning(`追溯码状态为 ${traceInfo.value.status}，不可上架`)
    }
  } catch {
    traceInfo.value = null
    message.error('未找到该追溯码')
  }
}

async function handleLocationScan(code: string) {
  locationCode.value = code
  // 执行混放检查
  if (traceInfo.value?.drug_id) {
    try {
      mixCheckResult.value = await checkMixPlacement(code)
    } catch {
      mixCheckResult.value = null
    }
  }
}

async function handleShelve() {
  if (!canShelve.value) return
  shelveLoading.value = true
  const currentTraceCode = traceCode.value
  const currentLocationCode = locationCode.value
  const currentDrugName = traceInfo.value?.drug_name
  try {
    await shelveTraceCode({ trace_code: currentTraceCode, location_code: currentLocationCode })
    const result: ShelvingResult = {
      trace_code: currentTraceCode,
      success: true,
      location_code: currentLocationCode,
      drug_name: currentDrugName,
      message: '上架成功',
    }
    lastResult.value = result
    recentResults.value.unshift(result)
    if (recentResults.value.length > 20) recentResults.value.length = 20
    traceCode.value = ''
    locationCode.value = ''
    traceInfo.value = null
    mixCheckResult.value = null
  } catch (e: unknown) {
    const result: ShelvingResult = {
      trace_code: currentTraceCode,
      success: false,
      message: (e as Error).message || '上架失败',
    }
    lastResult.value = result
    recentResults.value.unshift(result)
    if (recentResults.value.length > 20) recentResults.value.length = 20
  } finally {
    shelveLoading.value = false
  }
}

function isNearExpire(date: string): boolean {
  return dayjs(date).diff(dayjs(), 'day') <= 30
}
</script>

<style scoped>
.workbench-layout {
  display: grid;
  grid-template-columns: 360px 280px 1fr;
  gap: 16px;
  height: calc(100vh - 112px);
}

.shelf-scan-panel,
.shelf-info-panel,
.shelf-history-panel {
  height: 100%;
  overflow: hidden;
}
</style>
