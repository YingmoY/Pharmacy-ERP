<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <!-- 原销售单信息 -->
      <a-card title="原销售单信息" style="margin-bottom: 16px" v-if="order">
        <a-descriptions :column="4" size="small">
          <a-descriptions-item label="销售单号">{{ order.order_no }}</a-descriptions-item>
          <a-descriptions-item label="状态"><StatusTag type="sales" :value="order.status" /></a-descriptions-item>
          <a-descriptions-item label="实收金额">¥{{ order.actual_amount }}</a-descriptions-item>
          <a-descriptions-item label="已退款">¥{{ order.refund_amount }}</a-descriptions-item>
        </a-descriptions>
      </a-card>

      <a-row :gutter="16">
        <!-- 退货明细 -->
        <a-col :span="16">
          <a-card title="选择退货明细">
            <a-alert
              v-if="refundMode === 'FULL'"
              type="info"
              message="全单退货：将退还所有未退货明细"
              style="margin-bottom: 12px"
            />

            <a-table
              :columns="itemColumns"
              :data-source="refundableItems"
              :pagination="false"
              size="small"
              row-key="id"
              :row-selection="refundMode === 'PARTIAL' ? rowSelection : undefined"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'drug_name'">
                  <div>{{ record.drug_name }}</div>
                  <div style="font-size: 12px; color: #666">{{ record.specification }}</div>
                </template>
                <template v-if="column.key === 'price'">¥{{ record.unit_price }}</template>
                <template v-if="column.key === 'refund_status'">
                  <a-tag v-if="record.refund_status === 'REFUNDED'" color="red">已退货</a-tag>
                  <a-tag v-else color="green">可退货</a-tag>
                </template>
              </template>
            </a-table>
          </a-card>
        </a-col>

        <!-- 退货操作面板 -->
        <a-col :span="8">
          <a-card title="退货操作">
            <a-form layout="vertical">
              <a-form-item label="退货模式">
                <a-radio-group v-model:value="refundMode" button-style="solid">
                  <a-radio-button value="FULL">全单退货</a-radio-button>
                  <a-radio-button value="PARTIAL">部分退货</a-radio-button>
                </a-radio-group>
              </a-form-item>

              <a-form-item label="退款金额">
                <div style="font-size: 24px; font-weight: 700; color: #ff4d4f">
                  ¥{{ refundAmount }}
                </div>
                <!-- 医保全单退款：显示分摊明细 -->
                <template v-if="isMedicareFullRefund">
                  <a-descriptions size="small" :column="1" style="margin-top: 8px">
                    <a-descriptions-item label="退还客户（现金）">
                      <span style="color: #ff4d4f">¥{{ order!.personal_amount }}</span>
                    </a-descriptions-item>
                    <a-descriptions-item label="退还医保基金">
                      <span style="color: #fa8c16">¥{{ order!.medicare_amount }}</span>
                    </a-descriptions-item>
                  </a-descriptions>
                  <a-alert
                    type="warning"
                    message="医保结算订单退款：系统将自动发起医保撤销结算（2103）"
                    banner
                    style="margin-top: 8px"
                  />
                </template>
              </a-form-item>

              <a-form-item label="退货原因" required>
                <a-textarea
                  v-model:value="refundReason"
                  :rows="3"
                  placeholder="请输入退货原因"
                />
              </a-form-item>

              <a-alert
                type="warning"
                message="退货后追溯码将恢复在库状态，回到原货位"
                style="margin-bottom: 16px"
              />

              <a-button
                type="primary"
                danger
                block
                size="large"
                :disabled="!canRefund"
                :loading="refundLoading"
                @click="handleRefund"
              >
                确认退货 ¥{{ refundAmount }}
              </a-button>

              <a-button block style="margin-top: 8px" @click="router.back()">
                取消
              </a-button>
            </a-form>
          </a-card>
        </a-col>
      </a-row>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getSalesOrderDetail, refundSalesOrder } from '@/api/sales'
import type { SalesOrder, SalesOrderItem, RefundMode } from '@/types/sales'

const route = useRoute()
const router = useRouter()
const orderId = Number(route.params.id)

const loading = ref(false)
const refundLoading = ref(false)
const order = ref<SalesOrder | null>(null)
const refundMode = ref<RefundMode>('FULL')
const refundReason = ref('')
const selectedItemIds = ref<number[]>([])

// 可退货的明细（未退货的）
const refundableItems = computed(() =>
  (order.value?.items ?? []).filter((i) => i.refund_status === 'NONE'),
)

// 退款金额计算
const refundAmount = computed(() => {
  if (refundMode.value === 'FULL') {
    return refundableItems.value.reduce((s, i) => s + parseFloat(i.unit_price), 0).toFixed(2)
  }
  return order.value?.items
    ?.filter((i) => selectedItemIds.value.includes(i.id))
    .reduce((s, i) => s + parseFloat(i.unit_price), 0)
    .toFixed(2) ?? '0.00'
})

const canRefund = computed(() => {
  if (!refundReason.value.trim()) return false
  if (refundMode.value === 'PARTIAL' && selectedItemIds.value.length === 0) return false
  return true
})

// 是否为医保订单的全单退款（有分摊数据可展示）
const isMedicareFullRefund = computed(() =>
  refundMode.value === 'FULL' &&
  order.value?.payment_method === 'MEDICARE' &&
  parseFloat(order.value?.medicare_amount ?? '0') > 0,
)

const rowSelection = computed(() => ({
  selectedRowKeys: selectedItemIds.value,
  onChange: (keys: number[]) => {
    selectedItemIds.value = keys
  },
  getCheckboxProps: (record: SalesOrderItem) => ({
    disabled: record.refund_status !== 'NONE',
  }),
}))

const itemColumns = [
  { title: '药品', key: 'drug_name', ellipsis: true },
  { title: '追溯码', dataIndex: 'trace_code', width: 200 },
  { title: '批号', dataIndex: 'batch_number', width: 120 },
  { title: '单价', key: 'price', width: 80 },
  { title: '退货状态', key: 'refund_status', width: 80 },
]

async function handleRefund() {
  const confirmContent = isMedicareFullRefund.value
    ? `退还客户（现金）¥${order.value!.personal_amount}，退还医保基金 ¥${order.value!.medicare_amount}。\n此操作将同步撤销医保结算，不可撤销！`
    : `确认退款 ¥${refundAmount.value} 给顾客吗？此操作不可撤销！`

  Modal.confirm({
    title: '确认退货',
    content: confirmContent,
    okText: '确认退货',
    okType: 'danger',
    onOk: async () => {
      refundLoading.value = true
      try {
        await refundSalesOrder(orderId, {
          refund_mode: refundMode.value,
          refund_reason: refundReason.value,
          detail_ids: refundMode.value === 'PARTIAL' ? selectedItemIds.value : undefined,
        })
        message.success('退货成功')
        router.push(`/sales/orders/${orderId}`)
      } finally {
        refundLoading.value = false
      }
    },
  })
}

onMounted(async () => {
  loading.value = true
  try {
    order.value = await getSalesOrderDetail(orderId)
  } finally {
    loading.value = false
  }
})
</script>
