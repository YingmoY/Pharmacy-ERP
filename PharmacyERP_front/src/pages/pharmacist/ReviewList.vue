<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="销售单号">
          <a-input v-model:value="searchParams.order_no" placeholder="请输入" allow-clear />
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="PENDING">待审核</a-select-option>
            <a-select-option value="APPROVED">已通过</a-select-option>
            <a-select-option value="REJECTED">已驳回</a-select-option>
            <a-select-option value="CANCELLED">已取消</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="提交时间">
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
          <template v-if="column.key === 'submitter'">
            <span>{{ record.submitter_name || (record.submitter_id ? `用户#${record.submitter_id}` : '-') }}</span>
          </template>
          <template v-if="column.key === 'pharmacist'">
            <span>{{ record.pharmacist_name || (record.pharmacist_id ? `用户#${record.pharmacist_id}` : '-') }}</span>
          </template>
          <template v-if="column.key === 'status'">
            <StatusTag type="review" :value="record.status" />
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="router.push(`/pharmacist/reviews/${record.id}`)">
              {{ record.status === 'PENDING' ? '去审核' : '查看' }}
            </a-button>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import type { Dayjs } from 'dayjs'
import StatusTag from '@/components/common/StatusTag.vue'
import { getReviewList } from '@/api/pharmacist'
import type { PharmacistReview } from '@/types/sales'

const router = useRouter()
const loading = ref(false)
const list = ref<PharmacistReview[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)

const searchParams = reactive({
  order_no: '',
  status: undefined as string | undefined,
  submitted_start: undefined as string | undefined,
  submitted_end: undefined as string | undefined,
})

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '审核单号', dataIndex: 'review_no', ellipsis: true },
  { title: '销售单号', dataIndex: 'order_no', ellipsis: true },
  { title: '提交人', key: 'submitter', width: 100 },
  { title: '提交时间', dataIndex: 'submitted_at', width: 160 },
  { title: '状态', key: 'status', width: 100 },
  { title: '审核人', key: 'pharmacist', width: 100 },
  { title: '审核时间', dataIndex: 'reviewed_at', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  searchParams.submitted_start = dates?.[0]?.format('YYYY-MM-DD')
  searchParams.submitted_end = dates?.[1]?.format('YYYY-MM-DD')
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getReviewList({
      page: pagination.current,
      page_size: pagination.pageSize,
      status: searchParams.status,
      submitted_start: searchParams.submitted_start,
      submitted_end: searchParams.submitted_end,
    })
    list.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

function handleReset() {
  searchParams.order_no = ''
  searchParams.status = undefined
  searchParams.submitted_start = undefined
  searchParams.submitted_end = undefined
  dateRange.value = null
  pagination.current = 1
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
