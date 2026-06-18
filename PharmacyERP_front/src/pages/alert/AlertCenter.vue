<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="预警类型">
          <a-select v-model:value="searchParams.alert_type" placeholder="全部" style="width: 130px" allow-clear>
            <a-select-option value="NEAR_EXPIRE">近效期</a-select-option>
            <a-select-option value="LOW_STOCK">库存不足</a-select-option>
            <a-select-option value="MISPLACED">错架</a-select-option>
            <a-select-option value="LOSS">盘亏</a-select-option>
            <a-select-option value="SYSTEM">系统异常</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="级别">
          <a-select v-model:value="searchParams.priority" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option value="HIGH">紧急</a-select-option>
            <a-select-option value="MEDIUM">一般</a-select-option>
            <a-select-option value="LOW">低</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option value="ACTIVE">未处理</a-select-option>
            <a-select-option value="RESOLVED">已解决</a-select-option>
            <a-select-option value="IGNORED">已忽略</a-select-option>
          </a-select>
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
          <template v-if="column.key === 'alert_type'">
            <StatusTag type="alertType" :value="record.alert_type" />
          </template>
          <template v-if="column.key === 'priority'">
            <StatusTag type="alertLevel" :value="record.priority" />
          </template>
          <template v-if="column.key === 'status'">
            <StatusTag type="alert" :value="record.status" />
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions" v-if="record.status === 'ACTIVE'">
              <a-button type="link" size="small" @click="showResolve(record)">标记解决</a-button>
              <a-popconfirm title="确认忽略此预警？" @confirm="handleIgnore(record)">
                <a-button type="link" size="small" danger>忽略</a-button>
              </a-popconfirm>
            </div>
            <span v-else style="color: #999">-</span>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 标记解决弹窗 -->
    <a-modal v-model:open="resolveVisible" title="标记预警已解决" @ok="handleResolve" :confirm-loading="resolveLoading">
      <a-form :model="resolveForm" layout="vertical">
        <a-form-item label="处理备注" required>
          <a-textarea v-model:value="resolveForm.remark" :rows="4" placeholder="请填写处理备注" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getAlertList, resolveAlert, ignoreAlert } from '@/api/alert'
import type { AlertInfo } from '@/types/alert'

const loading = ref(false)
const resolveLoading = ref(false)
const resolveVisible = ref(false)
const currentAlert = ref<AlertInfo | null>(null)
const list = ref<AlertInfo[]>([])
const searchParams = reactive({ alert_type: undefined as string | undefined, priority: undefined as string | undefined, status: 'ACTIVE' as string | undefined })
const resolveForm = reactive({ remark: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '预警类型', key: 'alert_type', width: 100 },
  { title: '级别', key: 'priority', width: 80 },
  { title: '预警标题', dataIndex: 'title', ellipsis: true },
  { title: '内容', dataIndex: 'content', ellipsis: true },
  { title: '状态', key: 'status', width: 80 },
  { title: '触发时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 150 },
]

function showResolve(alert: AlertInfo) {
  currentAlert.value = alert
  resolveForm.remark = ''
  resolveVisible.value = true
}

async function handleResolve() {
  if (!resolveForm.remark.trim()) {
    message.warning('请填写处理备注')
    return
  }
  resolveLoading.value = true
  try {
    await resolveAlert(currentAlert.value!.id, { remark: resolveForm.remark })
    message.success('预警已标记为解决')
    resolveVisible.value = false
    fetchList()
  } finally { resolveLoading.value = false }
}

async function handleIgnore(alert: AlertInfo) {
  await ignoreAlert(alert.id)
  message.success('已忽略该预警')
  fetchList()
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getAlertList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { alert_type: undefined, priority: undefined, status: 'ACTIVE' })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

onMounted(fetchList)
</script>
