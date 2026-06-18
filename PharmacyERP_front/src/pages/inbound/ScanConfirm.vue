<template>
  <div class="page-content">
    <!-- 顶部信息栏 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="24">
        <a-card size="small">
          <a-descriptions :column="5" size="small">
            <a-descriptions-item label="入库单号">{{ order?.order_no }}</a-descriptions-item>
            <a-descriptions-item label="供应商">{{ order?.supplier_name }}</a-descriptions-item>
            <a-descriptions-item label="状态">
              <StatusTag type="inbound" :value="order?.status ?? ''" />
            </a-descriptions-item>
            <a-descriptions-item label="扫码进度">
              <a-progress
                :percent="progress"
                size="small"
                :format="progressFormat"
                style="width: 120px"
              />
            </a-descriptions-item>
            <a-descriptions-item label="操作">
              <a-space>
                <a-button
                  type="primary"
                  size="small"
                  :disabled="!canComplete"
                  :loading="completeLoading"
                  @click="handleComplete"
                >
                  完成入库
                </a-button>
                <a-button
                  size="small"
                  danger
                  @click="handleCancel"
                >
                  取消入库单
                </a-button>
              </a-space>
            </a-descriptions-item>
          </a-descriptions>
        </a-card>
      </a-col>
    </a-row>

    <div class="workbench-layout">
      <!-- 左侧：待确认明细 -->
      <div class="detail-panel">
        <a-card title="待确认明细" size="small" style="height: 100%; overflow: auto">
          <a-list :data-source="details" size="small">
            <template #renderItem="{ item }">
              <a-list-item
                :class="{ 'active-detail': selectedDetailId === item.id }"
                @click="selectedDetailId = item.id"
                style="cursor: pointer"
              >
                <a-list-item-meta>
                  <template #title>
                    <span>{{ item.drug_name }}</span>
                    <a-tag
                      :color="item.confirmed_qty >= item.planned_qty ? 'green' : 'blue'"
                      style="margin-left: 6px"
                      size="small"
                    >
                      {{ item.confirmed_qty }}/{{ item.planned_qty }}
                    </a-tag>
                  </template>
                  <template #description>
                    <span style="font-size: 11px">
                      批号:{{ item.batch_number }} 效期:{{ item.expire_date }}
                    </span>
                  </template>
                </a-list-item-meta>
                <template #extra>
                  <CheckCircleFilled
                    v-if="item.confirmed_qty >= item.planned_qty"
                    style="color: #52c41a; font-size: 18px"
                  />
                </template>
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </div>

      <!-- 右侧：扫码区 -->
      <div class="scan-panel">
        <a-card title="扫码确认" size="small" style="height: 100%">
          <!-- 当前选中明细 -->
          <div v-if="selectedDetail" class="selected-detail-info">
            <a-descriptions :column="2" size="small" bordered>
              <a-descriptions-item label="药品">{{ selectedDetail.drug_name }}</a-descriptions-item>
              <a-descriptions-item label="规格">{{ selectedDetail.specification }}</a-descriptions-item>
              <a-descriptions-item label="批号">{{ selectedDetail.batch_number }}</a-descriptions-item>
              <a-descriptions-item label="有效期">{{ selectedDetail.expire_date }}</a-descriptions-item>
              <a-descriptions-item label="计划数量">{{ selectedDetail.planned_qty }}</a-descriptions-item>
              <a-descriptions-item label="已确认">
                <span :style="{ color: selectedDetail.confirmed_qty >= selectedDetail.planned_qty ? '#52c41a' : '#1677ff' }">
                  {{ selectedDetail.confirmed_qty }}
                </span>
              </a-descriptions-item>
            </a-descriptions>
          </div>

          <a-alert
            v-else
            type="info"
            message="请先在左侧选择一个入库明细，再进行扫码确认"
            style="margin-bottom: 16px"
          />

          <!-- 扫码输入 -->
          <div style="margin-top: 16px">
            <ScanInput
              ref="scanInputRef"
              :disabled="!selectedDetailId || order?.status !== 'PENDING_CONFIRM'"
              placeholder="扫描追溯码，回车确认"
              @submit="handleScan"
              :auto-focus="true"
            />
          </div>

          <!-- 扫码结果提示 -->
          <div style="margin-top: 12px; max-height: 300px; overflow-y: auto">
            <div
              v-for="(result, idx) in scanResults"
              :key="idx"
              :class="['scan-result-' + result.type]"
            >
              <div style="font-weight: 500; font-size: 13px">{{ result.trace_code }}</div>
              <div style="font-size: 12px">{{ result.message }}</div>
            </div>
          </div>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { CheckCircleFilled } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import ScanInput from '@/components/common/ScanInput.vue'
import {
  getInboundOrderDetail,
  getInboundDetailList,
  confirmTraceCode,
  completeInboundOrder,
  cancelInboundOrder,
} from '@/api/inbound'
import type { InboundOrder, InboundOrderDetail } from '@/types/inbound'

const route = useRoute()
const router = useRouter()
const orderId = Number(route.params.id)

const order = ref<InboundOrder | null>(null)
const details = ref<InboundOrderDetail[]>([])
const selectedDetailId = ref<number | null>(null)
const completeLoading = ref(false)

const scanResults = ref<Array<{ trace_code: string; message: string; type: 'success' | 'error' | 'warning' }>>([])

const selectedDetail = computed(() =>
  details.value.find((d) => d.id === selectedDetailId.value) ?? null,
)

const plannedTotal = computed(() => details.value.reduce((s, d) => s + d.planned_qty, 0))
const confirmedTotal = computed(() => details.value.reduce((s, d) => s + d.confirmed_qty, 0))
const progress = computed(() =>
  plannedTotal.value ? Math.round((confirmedTotal.value / plannedTotal.value) * 100) : 0,
)

function progressFormat() {
  return `${confirmedTotal.value}/${plannedTotal.value}`
}

const canComplete = computed(() =>
  order.value?.status === 'PENDING_CONFIRM' &&
  confirmedTotal.value > 0 &&
  confirmedTotal.value === plannedTotal.value,
)

async function loadData() {
  order.value = await getInboundOrderDetail(orderId)
  details.value = (await getInboundDetailList(orderId)) ?? []
  if (details.value.length > 0 && !selectedDetailId.value) {
    const pending = details.value.find((d) => d.confirmed_qty < d.planned_qty)
    selectedDetailId.value = pending?.id ?? details.value[0].id
  }
}

async function refreshDetails() {
  details.value = (await getInboundDetailList(orderId)) ?? []
}

async function handleScan(traceCode: string) {
  if (!selectedDetailId.value) return

  try {
    await confirmTraceCode(orderId, {
      detail_id: selectedDetailId.value,
      trace_code: traceCode,
    })

    scanResults.value.unshift({ trace_code: traceCode, message: '扫码成功', type: 'success' })

    // 刷新明细列表以获取最新已确认数量
    await refreshDetails()

    // 自动跳转到下一个未完成明细
    const detail = details.value.find((d) => d.id === selectedDetailId.value)
    if (detail && detail.confirmed_qty >= detail.planned_qty) {
      const next = details.value.find((d) => d.id !== detail.id && d.confirmed_qty < d.planned_qty)
      if (next) selectedDetailId.value = next.id
    }
  } catch (err: unknown) {
    const msg = (err as Error).message || '扫码失败'
    scanResults.value.unshift({ trace_code: traceCode, message: msg, type: 'error' })
  }

  if (scanResults.value.length > 20) scanResults.value.length = 20
}

async function handleComplete() {
  Modal.confirm({
    title: '确认完成入库',
    content: '完成入库后将无法再修改，确认所有追溯码已扫描完毕吗？',
    okText: '确认完成',
    cancelText: '取消',
    onOk: async () => {
      completeLoading.value = true
      try {
        await completeInboundOrder(orderId)
        message.success('入库已完成！')
        router.push(`/inbound/orders/${orderId}`)
      } finally {
        completeLoading.value = false
      }
    },
  })
}

async function handleCancel() {
  Modal.confirm({
    title: '确认取消入库单',
    content: '取消入库单将删除所有已扫码的追溯码库存记录，此操作不可撤销！',
    okText: '确认取消',
    okType: 'danger',
    cancelText: '保留',
    onOk: async () => {
      await cancelInboundOrder(orderId)
      message.success('入库单已取消')
      router.push('/inbound/orders')
    },
  })
}

onMounted(loadData)
</script>

<style scoped>
.workbench-layout {
  display: grid;
  grid-template-columns: 320px 1fr;
  gap: 16px;
  height: calc(100vh - 220px);
}

.detail-panel,
.scan-panel {
  height: 100%;
  overflow: hidden;
}

.active-detail {
  background: #e6f4ff;
}

.selected-detail-info {
  margin-bottom: 16px;
}
</style>
