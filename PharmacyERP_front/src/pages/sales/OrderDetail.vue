<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <a-card style="margin-bottom: 16px" v-if="order">
        <a-row :gutter="24" align="middle">
          <a-col :flex="1">
            <a-space size="large">
              <div>
                <div style="font-size: 12px; color: #999">销售单号</div>
                <div style="font-size: 18px; font-weight: 600">{{ order.order_no }}</div>
              </div>
              <div>
                <div style="font-size: 12px; color: #999">状态</div>
                <StatusTag type="sales" :value="order.status" />
                <a-tag v-if="order.need_audit" color="orange" style="margin-left: 4px">需药师审核</a-tag>
              </div>
              <div>
                <div style="font-size: 12px; color: #999">应收 / 实收</div>
                <div>¥{{ order.total_amount }} / <span style="color: #e6550d">¥{{ order.actual_amount }}</span></div>
              </div>
              <div v-if="parseFloat(order.refund_amount) > 0">
                <div style="font-size: 12px; color: #999">退款金额</div>
                <div style="color: #ff4d4f">¥{{ order.refund_amount }}</div>
              </div>
            </a-space>
          </a-col>
          <a-col>
            <a-space>
              <a-button
                v-if="['COMPLETED', 'PARTIALLY_REFUNDED'].includes(order.status)"
                type="primary"
                @click="router.push(`/sales/orders/${order.id}/refund`)"
              >
                申请退货
              </a-button>
              <a-button
                v-if="['PENDING', 'PENDING_REVIEW', 'APPROVED'].includes(order.status)"
                danger
                @click="handleCancel"
              >
                取消订单
              </a-button>
            </a-space>
          </a-col>
        </a-row>
      </a-card>

      <a-card>
        <a-tabs>
          <a-tab-pane key="items" tab="销售明细">
            <a-table
              :columns="itemColumns"
              :data-source="order?.items ?? []"
              :pagination="false"
              size="small"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'drug_name'">
                  <div>{{ record.drug_name }}</div>
                  <div style="font-size: 12px; color: #666">{{ record.specification }}</div>
                </template>
                <template v-if="column.key === 'refund_status'">
                  <StatusTag type="refund" :value="record.refund_status" />
                </template>
                <template v-if="column.key === 'price'">¥{{ record.unit_price }}</template>
              </template>
            </a-table>
          </a-tab-pane>

          <a-tab-pane key="review" tab="药师审核">
            <a-empty v-if="!order?.review" description="本销售单无需药师审核" />
            <a-descriptions :column="2" bordered size="small" v-else>
              <a-descriptions-item label="审核单号">{{ order.review.review_no }}</a-descriptions-item>
              <a-descriptions-item label="审核状态">
                <StatusTag type="review" :value="order.review.status" />
              </a-descriptions-item>
              <a-descriptions-item label="提交人">{{ order.review.submitter_name }}</a-descriptions-item>
              <a-descriptions-item label="提交时间">{{ order.review.submitted_at }}</a-descriptions-item>
              <a-descriptions-item label="审核人">{{ order.review.pharmacist_name || '-' }}</a-descriptions-item>
              <a-descriptions-item label="审核时间">{{ order.review.reviewed_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="审核意见" :span="2">{{ order.review.review_opinion || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <a-tab-pane key="info" tab="基础信息">
            <a-descriptions :column="2" bordered size="small" v-if="order">
              <a-descriptions-item label="收银员">{{ order.cashier_name }}</a-descriptions-item>
              <a-descriptions-item label="支付方式">{{ order.payment_method || '-' }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ order.created_at }}</a-descriptions-item>
              <a-descriptions-item label="结算时间">{{ order.paid_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="取消时间">{{ order.cancelled_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="退款时间">{{ order.refunded_at || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getSalesOrderDetail, cancelSalesOrder } from '@/api/sales'
import type { SalesOrder } from '@/types/sales'

const route = useRoute()
const router = useRouter()
const orderId = Number(route.params.id)
const loading = ref(false)
const order = ref<SalesOrder | null>(null)

const itemColumns = [
  { title: '药品', key: 'drug_name', ellipsis: true },
  { title: '追溯码', dataIndex: 'trace_code', width: 200 },
  { title: '批号', dataIndex: 'batch_number', width: 120 },
  { title: '有效期', dataIndex: 'expire_date', width: 100 },
  { title: '货位', dataIndex: 'location_code', width: 80 },
  { title: '单价', key: 'price', width: 80 },
  { title: '退货状态', key: 'refund_status', width: 90 },
]

async function handleCancel() {
  Modal.confirm({
    title: '确认取消销售单',
    content: '取消后预占的追溯码将被释放，确认取消吗？',
    okType: 'danger',
    okText: '确认取消',
    onOk: async () => {
      await cancelSalesOrder(orderId)
      message.success('销售单已取消')
      router.push('/sales/orders')
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
