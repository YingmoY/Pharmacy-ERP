<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="事件类型">
          <a-select v-model:value="searchParams.event_type" placeholder="全部" style="width: 140px" allow-clear>
            <a-select-option value="LOGIN_FAIL">登录失败</a-select-option>
            <a-select-option value="PERMISSION_DENY">权限拒绝</a-select-option>
            <a-select-option value="UNUSUAL_ACCESS">异常访问</a-select-option>
            <a-select-option value="ACCOUNT_LOCKED">账号锁定</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="风险等级">
          <a-select v-model:value="searchParams.severity" placeholder="全部" style="width: 110px" allow-clear>
            <a-select-option value="LOW">低</a-select-option>
            <a-select-option value="MEDIUM">中</a-select-option>
            <a-select-option value="HIGH">高</a-select-option>
            <a-select-option value="CRITICAL">严重</a-select-option>
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
          <template v-if="column.key === 'severity'">
            <a-tag :color="severityColor(record.severity)">{{ record.severity }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getSecurityEvents } from '@/api/system'
import type { SecurityEvent } from '@/types/system'

const loading = ref(false)
const list = ref<SecurityEvent[]>([])
const searchParams = reactive({ event_type: undefined as string | undefined, severity: undefined as string | undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '事件类型', dataIndex: 'event_type', width: 130 },
  { title: '风险等级', key: 'severity', width: 90 },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: 'IP地址', dataIndex: 'ip', width: 140 },
  { title: '发生时间', dataIndex: 'created_at', width: 160 },
]

const severityColorMap: Record<string, string> = { CRITICAL: 'red', HIGH: 'volcano', MEDIUM: 'orange', LOW: 'blue' }
function severityColor(level: string): string { return severityColorMap[level] ?? 'default' }

async function fetchList() {
  loading.value = true
  try {
    const res = await getSecurityEvents({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { event_type: undefined, severity: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
