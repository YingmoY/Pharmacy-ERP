<template>
  <div class="page-content">
    <!-- 搜索区域 -->
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="handleSearch">
        <a-form-item label="入库单号">
          <a-input v-model:value="searchParams.order_no" placeholder="请输入" allow-clear />
        </a-form-item>
        <a-form-item label="供应商">
          <a-input v-model:value="searchParams.supplier_name" placeholder="供应商名称" allow-clear />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 140px" allow-clear>
            <a-select-option value="DRAFT">草稿</a-select-option>
            <a-select-option value="PENDING_CONFIRM">待扫码确认</a-select-option>
            <a-select-option value="COMPLETED">已完成</a-select-option>
            <a-select-option value="CANCELLED">已取消</a-select-option>
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

    <!-- 操作按钮 -->
    <div style="margin-bottom: 12px; display: flex; gap: 8px">
      <PermissionButton type="primary" permission="inbound.orders.create" @click="router.push('/inbound/orders/new')">
        <PlusOutlined /> 新建入库单
      </PermissionButton>
    </div>

    <!-- 表格 -->
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
            <StatusTag type="inbound" :value="record.status" />
          </template>
          <template v-if="column.key === 'total_amount'">
            ¥{{ parseFloat(record.total_amount).toFixed(2) }}
          </template>
          <template v-if="column.key === 'progress'">
            <a-progress
              :percent="record.total_planned_qty ? Math.round(record.total_confirmed_qty / record.total_planned_qty * 100) : 0"
              size="small"
              :format="() => String(record.total_confirmed_qty) + '/' + String(record.total_planned_qty)"
            />
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="router.push(`/inbound/orders/${record.id}`)">
                查看
              </a-button>
              <template v-if="record.status === 'DRAFT'">
                <a-button type="link" size="small" @click="router.push(`/inbound/orders/${record.id}/edit`)">
                  编辑
                </a-button>
                <a-button type="link" size="small" @click="handleSubmit(record.id)">提交</a-button>
                <a-button type="link" size="small" danger @click="handleCancel(record.id)">取消</a-button>
              </template>
              <template v-if="record.status === 'PENDING_CONFIRM'">
                <a-button type="link" size="small" @click="router.push(`/inbound/orders/${record.id}/confirm`)">
                  扫码确认
                </a-button>
              </template>
              <template v-if="record.status === 'COMPLETED'">
                <a-button type="link" size="small" @click="router.push('/shelving/workbench')">上架</a-button>
              </template>
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
import { PlusOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import PermissionButton from '@/components/common/PermissionButton.vue'
import { getInboundOrderList, submitInboundOrder, cancelInboundOrder } from '@/api/inbound'
import type { InboundOrder } from '@/types/inbound'

const router = useRouter()
const loading = ref(false)
const list = ref<InboundOrder[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)

const searchParams = reactive({
  order_no: '',
  supplier_name: '',
  status: undefined as string | undefined,
  created_start: undefined as string | undefined,
  created_end: undefined as string | undefined,
})

const pagination = reactive({
  current: 1,
  pageSize: 20,
  total: 0,
  showSizeChanger: true,
  showTotal: (total: number) => `共 ${total} 条`,
})

const columns = [
  { title: '入库单号', dataIndex: 'order_no', key: 'order_no', ellipsis: true },
  { title: '供应商', dataIndex: 'supplier_name', key: 'supplier_name', ellipsis: true },
  { title: '状态', key: 'status', width: 120 },
  { title: '总金额', key: 'total_amount', width: 100 },
  { title: '扫码进度', key: 'progress', width: 160 },
  { title: '创建人', dataIndex: 'creator_name', key: 'creator_name', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 160, ellipsis: true },
  { title: '操作', key: 'action', width: 220 },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.created_start = dates[0].format('YYYY-MM-DD')
    searchParams.created_end = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.created_start = undefined
    searchParams.created_end = undefined
  }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInboundOrderList({
      page: pagination.current,
      page_size: pagination.pageSize,
      order_no: searchParams.order_no || undefined,
      status: searchParams.status as any,
      created_start: searchParams.created_start,
      created_end: searchParams.created_end,
    })
    list.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.current = 1
  fetchList()
}

function handleReset() {
  searchParams.order_no = ''
  searchParams.supplier_name = ''
  searchParams.status = undefined
  searchParams.created_start = undefined
  searchParams.created_end = undefined
  dateRange.value = null
  handleSearch()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

async function handleSubmit(id: number) {
  Modal.confirm({
    title: '确认提交',
    content: '提交后入库单将进入待扫码确认状态，明细不可再修改，确认提交吗？',
    okText: '确认提交',
    onOk: async () => {
      await submitInboundOrder(id)
      message.success('入库单已提交')
      fetchList()
    },
  })
}

async function handleCancel(id: number) {
  Modal.confirm({
    title: '确认取消',
    content: '确认取消此入库单吗？',
    okText: '确认取消',
    okType: 'danger',
    onOk: async () => {
      await cancelInboundOrder(id)
      message.success('入库单已取消')
      fetchList()
    },
  })
}

onMounted(fetchList)
</script>
