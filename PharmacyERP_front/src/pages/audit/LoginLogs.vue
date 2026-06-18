<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="用户名"><a-input v-model:value="searchParams.username" allow-clear /></a-form-item>
        <a-form-item label="结果">
          <a-select v-model:value="searchParams.success" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option :value="true">成功</a-select-option>
            <a-select-option :value="false">失败</a-select-option>
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
          <template v-if="column.key === 'success'">
            <a-tag :color="record.success ? 'green' : 'red'">{{ record.success ? '成功' : '失败' }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getLoginLogs } from '@/api/system'
import type { LoginLog } from '@/types/system'

const loading = ref(false)
const list = ref<LoginLog[]>([])
const searchParams = reactive({ username: '', success: undefined as boolean | undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: 'IP地址', dataIndex: 'ip', width: 140 },
  { title: '结果', key: 'success', width: 70 },
  { title: '失败原因', dataIndex: 'message', ellipsis: true },
  { title: '设备信息', dataIndex: 'user_agent', ellipsis: true },
  { title: '登录时间', dataIndex: 'created_at', width: 160 },
]

async function fetchList() {
  loading.value = true
  try {
    const res = await getLoginLogs({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { username: '', success: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
