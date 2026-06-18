<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <!-- 顶部状态卡 -->
      <a-card class="status-card" style="margin-bottom: 16px" v-if="order">
        <a-row :gutter="24" align="middle">
          <a-col :flex="1">
            <a-space size="large">
              <div>
                <div style="font-size: 12px; color: #999">入库单号</div>
                <div style="font-size: 18px; font-weight: 600">{{ order.order_no }}</div>
              </div>
              <div>
                <div style="font-size: 12px; color: #999">状态</div>
                <StatusTag type="inbound" :value="order.status" />
              </div>
              <div>
                <div style="font-size: 12px; color: #999">供应商</div>
                <div>{{ order.supplier_name }}</div>
              </div>
              <div>
                <div style="font-size: 12px; color: #999">总金额</div>
                <div style="color: #e6550d; font-weight: 600">¥{{ parseFloat(order.total_amount).toFixed(2) }}</div>
              </div>
              <div>
                <div style="font-size: 12px; color: #999">扫码进度</div>
                <a-progress
                  :percent="progress"
                  size="small"
                  :format="() => `${confirmedTotal}/${plannedTotal}`"
                  style="width: 100px"
                />
              </div>
            </a-space>
          </a-col>
          <a-col>
            <a-space>
              <a-button
                v-if="order.status === 'DRAFT'"
                @click="router.push(`/inbound/orders/${order.id}/edit`)"
              >
                编辑
              </a-button>
              <a-button
                v-if="order.status === 'DRAFT'"
                type="primary"
                @click="handleSubmit"
              >
                提交入库
              </a-button>
              <a-button
                v-if="order.status === 'PENDING_CONFIRM'"
                type="primary"
                @click="router.push(`/inbound/orders/${order.id}/confirm`)"
              >
                进入扫码确认
              </a-button>
              <a-button
                v-if="['DRAFT', 'PENDING_CONFIRM'].includes(order.status)"
                danger
                @click="handleCancel"
              >
                取消入库单
              </a-button>
            </a-space>
          </a-col>
        </a-row>
      </a-card>

      <!-- Tabs -->
      <a-card>
        <a-tabs v-model:activeKey="activeTab">
          <!-- 基础信息 -->
          <a-tab-pane key="basic" tab="基础信息">
            <a-descriptions :column="3" bordered size="small" v-if="order">
              <a-descriptions-item label="入库单号">{{ order.order_no }}</a-descriptions-item>
              <a-descriptions-item label="供应商">{{ order.supplier_name }}</a-descriptions-item>
              <a-descriptions-item label="发票号">{{ order.invoice_no || '-' }}</a-descriptions-item>
              <a-descriptions-item label="创建人">{{ order.creator_name }}</a-descriptions-item>
              <a-descriptions-item label="创建时间">{{ order.created_at }}</a-descriptions-item>
              <a-descriptions-item label="提交时间">{{ order.submitted_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="完成时间">{{ order.completed_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="取消时间">{{ order.cancelled_at || '-' }}</a-descriptions-item>
              <a-descriptions-item label="备注" :span="3">{{ order.remark || '-' }}</a-descriptions-item>
            </a-descriptions>
          </a-tab-pane>

          <!-- 入库明细 -->
          <a-tab-pane key="details" tab="入库明细">
            <a-table
              :columns="detailColumns"
              :data-source="details"
              :loading="detailLoading"
              row-key="id"
              :pagination="false"
              size="small"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'progress'">
                  <a-progress
                    :percent="record.planned_qty ? Math.round(record.confirmed_qty / record.planned_qty * 100) : 0"
                    size="small"
                    :format="() => `${record.confirmed_qty}/${record.planned_qty}`"
                  />
                </template>
                <template v-if="column.key === 'amount'">
                  ¥{{ parseFloat(record.amount).toFixed(2) }}
                </template>
              </template>
            </a-table>
          </a-tab-pane>
        </a-tabs>
      </a-card>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInboundOrderDetail, getInboundDetailList, submitInboundOrder, cancelInboundOrder } from '@/api/inbound'
import type { InboundOrder, InboundOrderDetail } from '@/types/inbound'

const route = useRoute()
const router = useRouter()
const orderId = Number(route.params.id)

const order = ref<InboundOrder | null>(null)
const details = ref<InboundOrderDetail[]>([])
const loading = ref(false)
const detailLoading = ref(false)
const activeTab = ref('basic')

const plannedTotal = computed(() => details.value.reduce((s, d) => s + d.planned_qty, 0))
const confirmedTotal = computed(() => details.value.reduce((s, d) => s + d.confirmed_qty, 0))
const progress = computed(() =>
  plannedTotal.value ? Math.round((confirmedTotal.value / plannedTotal.value) * 100) : 0,
)

const detailColumns = [
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 120 },
  { title: '厂家', dataIndex: 'manufacturer', ellipsis: true },
  { title: '批号', dataIndex: 'batch_number', width: 120 },
  { title: '有效期', dataIndex: 'expire_date', width: 100 },
  { title: '单价', dataIndex: 'unit_price', width: 80 },
  { title: '金额', key: 'amount', width: 100 },
  { title: '扫码进度', key: 'progress', width: 140 },
]

async function loadData() {
  loading.value = true
  try {
    order.value = await getInboundOrderDetail(orderId)
  } finally {
    loading.value = false
  }

  detailLoading.value = true
  try {
    details.value = (await getInboundDetailList(orderId)) ?? []
  } finally {
    detailLoading.value = false
  }
}

async function handleSubmit() {
  Modal.confirm({
    title: '确认提交入库单',
    content: '提交后明细不可再修改，确认提交吗？',
    okText: '确认提交',
    onOk: async () => {
      await submitInboundOrder(orderId)
      message.success('入库单已提交')
      loadData()
    },
  })
}

async function handleCancel() {
  Modal.confirm({
    title: '确认取消入库单',
    content: '取消后将删除所有已扫码追溯码，此操作不可撤销！',
    okText: '确认取消',
    okType: 'danger',
    onOk: async () => {
      await cancelInboundOrder(orderId)
      message.success('入库单已取消')
      router.push('/inbound/orders')
    },
  })
}

onMounted(loadData)
</script>
