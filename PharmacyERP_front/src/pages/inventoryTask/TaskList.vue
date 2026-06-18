<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="任务号"><a-input v-model:value="searchParams.task_no" allow-clear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="PENDING">待开始</a-select-option>
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

    <div style="margin-bottom: 12px">
      <a-button type="primary" @click="showCreate"><PlusOutlined /> 新建盘库任务</a-button>
    </div>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="inventoryTask" :value="record.status" />
          </template>
          <template v-if="column.key === 'task_no'">
            <router-link :to="`/inventory-tasks/${record.id}`">{{ record.task_no }}</router-link>
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="$router.push(`/inventory-tasks/${record.id}`)">详情</a-button>
              <a-button v-if="record.status === 'IN_PROGRESS'" type="link" size="small" @click="$router.push(`/inventory-tasks/${record.id}/scan`)">继续扫码</a-button>
              <a-button v-if="record.status === 'PENDING'" type="link" size="small" @click="startTask(record.id)">开始</a-button>
              <a-popconfirm v-if="['PENDING', 'IN_PROGRESS'].includes(record.status)" title="确认取消此盘库任务？" @confirm="cancelTask(record.id)">
                <a-button type="link" size="small" danger>取消</a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 新建任务弹窗 -->
    <a-modal v-model:open="createVisible" title="新建盘库任务" @ok="handleCreate" :confirm-loading="createLoading">
      <a-form :model="createForm" layout="vertical" ref="createFormRef" :rules="createRules">
        <a-form-item label="盘点范围类型" name="scope_type">
          <a-select v-model:value="createForm.scope_type">
            <a-select-option value="LOCATION">按货位</a-select-option>
            <a-select-option value="SHELF">按货架</a-select-option>
            <a-select-option value="AREA">按区域</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item v-if="createForm.scope_type !== 'FULL'" label="范围值" name="scope_value">
          <a-input v-model:value="createForm.scope_value" placeholder="货位码或药品编码" />
        </a-form-item>
        <a-form-item label="备注">
          <a-textarea v-model:value="createForm.remark" :rows="2" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInventoryTaskList, createInventoryTask, startInventoryTask, cancelInventoryTask } from '@/api/inventoryTask'
import type { InventoryTask } from '@/types/inventoryTask'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const createLoading = ref(false)
const createVisible = ref(false)
const list = ref<InventoryTask[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const createFormRef = ref()

const searchParams = reactive({ task_no: '', status: undefined as string | undefined, start_date: '', end_date: '' })
const createForm = reactive({ scope_type: 'LOCATION' as string, scope_value: '', remark: '' })
const createRules = {
  scope_type: [{ required: true, message: '请选择盘点范围类型' }],
  scope_value: [{ required: true, message: '请输入范围值' }],
}

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })
const columns = [
  { title: '任务号', key: 'task_no', width: 160 },
  { title: '范围类型', dataIndex: 'scope_type', width: 90 },
  { title: '范围值', dataIndex: 'scope_value', width: 120 },
  { title: '状态', key: 'status', width: 90 },
  { title: '已扫', dataIndex: 'scanned_count', width: 80, customRender: ({ record }: { record: InventoryTask }) => record.scanned_count ?? 0 },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 180 },
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

function showCreate() {
  Object.assign(createForm, { scope_type: 'LOCATION', scope_value: '', remark: '' })
  createVisible.value = true
}

async function handleCreate() {
  await createFormRef.value?.validate()
  createLoading.value = true
  try {
    await createInventoryTask(createForm)
    message.success('盘库任务创建成功')
    createVisible.value = false
    fetchList()
  } finally { createLoading.value = false }
}

async function startTask(id: number) {
  await startInventoryTask(id)
  message.success('任务已开始')
  fetchList()
}

async function cancelTask(id: number) {
  await cancelInventoryTask(id)
  message.success('任务已取消')
  fetchList()
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInventoryTaskList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
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
