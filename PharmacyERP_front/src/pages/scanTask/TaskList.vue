<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="任务号"><a-input v-model:value="searchParams.task_no" allow-clear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 110px" allow-clear>
            <a-select-option value="PENDING">待分配</a-select-option>
            <a-select-option value="ASSIGNED">已分配</a-select-option>
            <a-select-option value="IN_PROGRESS">进行中</a-select-option>
            <a-select-option value="COMPLETED">已完成</a-select-option>
            <a-select-option value="CANCELLED">已取消</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="创建时间">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="scanTask" :value="record.status" />
          </template>
          <template v-if="column.key === 'task_no'">
            <router-link :to="`/scan-tasks/${record.id}`">{{ record.task_no }}</router-link>
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="$router.push(`/scan-tasks/${record.id}`)">详情</a-button>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getScanTaskList } from '@/api/scanTask'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ task_no: '', status: undefined as string | undefined, start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '任务号', key: 'task_no', width: 160 },
  { title: '任务类型', dataIndex: 'task_type', width: 100 },
  { title: '状态', key: 'status', width: 90 },
  { title: '指派员工', dataIndex: 'assigned_to_name', width: 100 },
  { title: '关联单号', dataIndex: 'related_order_no', width: 160 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.start_date = dates[0].format('YYYY-MM-DD')
    searchParams.end_date = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.start_date = ''
    searchParams.end_date = ''
  }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getScanTaskList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { task_no: '', status: undefined, start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
