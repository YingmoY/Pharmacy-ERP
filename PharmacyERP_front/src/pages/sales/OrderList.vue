<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="销售单号">
          <a-input v-model:value="searchParams.order_no" allow-clear />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 130px" allow-clear>
            <a-select-option value="PENDING">待结算</a-select-option>
            <a-select-option value="PENDING_REVIEW">待药师审核</a-select-option>
            <a-select-option value="APPROVED">审核通过</a-select-option>
            <a-select-option value="COMPLETED">已完成</a-select-option>
            <a-select-option value="PARTIALLY_REFUNDED">部分退货</a-select-option>
            <a-select-option value="REFUNDED">已退货</a-select-option>
            <a-select-option value="CANCELLED">已取消</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="支付方式">
          <a-select v-model:value="searchParams.payment_method" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="CASH">现金</a-select-option>
            <a-select-option value="ALIPAY">支付宝</a-select-option>
            <a-select-option value="WECHAT">微信</a-select-option>
            <a-select-option value="BANK_CARD">银行卡</a-select-option>
            <a-select-option value="MEDICARE">医保</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="创建时间">
          <a-range-picker v-model:value="dateRange" @change="handleDateChange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="sales" :value="record.status" />
          </template>
          <template v-if="column.key === 'need_audit'">
            <a-tag v-if="record.need_audit" color="orange">需审核</a-tag>
            <a-tag v-else color="green">无需审核</a-tag>
          </template>
          <template v-if="column.key === 'amount'">
            <div>
              <div>应收：¥{{ record.total_amount }}</div>
              <div style="color: #e6550d">实收：¥{{ record.actual_amount }}</div>
              <div v-if="parseFloat(record.refund_amount) > 0" style="color: #ff4d4f">
                退款：¥{{ record.refund_amount }}
              </div>
            </div>
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="router.push(`/sales/orders/${record.id}`)">
                查看
              </a-button>
              <a-button
                v-if="['PENDING', 'APPROVED'].includes(record.status)"
                type="link"
                size="small"
                @click="router.push(`/pos?orderId=${record.id}`)"
              >
                结算
              </a-button>
              <a-button
                v-if="['COMPLETED', 'PARTIALLY_REFUNDED'].includes(record.status)"
                type="link"
                size="small"
                @click="router.push(`/sales/orders/${record.id}/refund`)"
              >
                退货
              </a-button>
              <a-button
                v-if="['PENDING', 'PENDING_REVIEW', 'APPROVED'].includes(record.status)"
                type="link"
                size="small"
                danger
                @click="handleCancel(record.id)"
              >
                取消
              </a-button>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import type { Dayjs } from 'dayjs'
import StatusTag from '@/components/common/StatusTag.vue'
import { getSalesOrderList, cancelSalesOrder } from '@/api/sales'
import type { SalesOrder } from '@/types/sales'

const router = useRouter()
const loading = ref(false)
const list = ref<SalesOrder[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)

const searchParams = reactive({
  order_no: '',
  status: undefined as string | undefined,
  payment_method: undefined as string | undefined,
  created_start: undefined as string | undefined,
  created_end: undefined as string | undefined,
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '销售单号', dataIndex: 'order_no', ellipsis: true },
  { title: '状态', key: 'status', width: 110 },
  { title: '金额', key: 'amount', width: 140 },
  { title: '支付方式', dataIndex: 'payment_method', width: 90 },
  { title: '是否审核', key: 'need_audit', width: 90 },
  { title: '收银员', dataIndex: 'cashier_name', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 180 },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  searchParams.created_start = dates?.[0]?.format('YYYY-MM-DD')
  searchParams.created_end = dates?.[1]?.format('YYYY-MM-DD')
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getSalesOrderList({
      page: pagination.current,
      page_size: pagination.pageSize,
      order_no: searchParams.order_no || undefined,
      status: searchParams.status as any,
      payment_method: searchParams.payment_method as any,
      created_start: searchParams.created_start,
      created_end: searchParams.created_end,
    })
    list.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

function handleReset() {
  Object.assign(searchParams, { order_no: '', status: undefined, payment_method: undefined, created_start: undefined, created_end: undefined })
  dateRange.value = null
  pagination.current = 1
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

async function handleCancel(id: number) {
  Modal.confirm({
    title: '确认取消销售单',
    content: '取消后预占的追溯码将被释放，确认取消吗？',
    okText: '确认取消',
    okType: 'danger',
    onOk: async () => {
      await cancelSalesOrder(id)
      message.success('销售单已取消')
      fetchList()
    },
  })
}

onMounted(fetchList)
</script>
