<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="报表类型">
          <a-select v-model:value="searchParams.report_type" placeholder="全部" style="width: 130px" allow-clear>
            <a-select-option value="SALES">销售报表</a-select-option>
            <a-select-option value="INBOUND">入库报表</a-select-option>
            <a-select-option value="INVENTORY">库存报表</a-select-option>
            <a-select-option value="TRACE_LOG">追溯日志</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 100px" allow-clear>
            <a-select-option value="PENDING">等待中</a-select-option>
            <a-select-option value="PROCESSING">处理中</a-select-option>
            <a-select-option value="COMPLETED">已完成</a-select-option>
            <a-select-option value="FAILED">失败</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
          <a-button style="margin-left: 8px" @click="fetchList"><ReloadOutlined /> 刷新</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="export" :value="record.status" />
          </template>
          <template v-if="column.key === 'progress'">
            <a-progress :percent="record.progress ?? 0" size="small" v-if="['PROCESSING'].includes(record.status)" />
            <span v-else>-</span>
          </template>
          <template v-if="column.key === 'action'">
            <a-button
              v-if="record.status === 'COMPLETED' && record.download_url"
              type="link"
              size="small"
              :href="record.download_url"
              target="_blank"
            >
              下载
            </a-button>
            <span v-else style="color: #999">-</span>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ReloadOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import type { ReportExportTask } from '@/types/report'

const loading = ref(false)
const list = ref<ReportExportTask[]>([])
const searchParams = reactive({ report_type: undefined as string | undefined, status: undefined as string | undefined })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '任务ID', dataIndex: 'id', width: 80 },
  { title: '报表类型', dataIndex: 'report_type', width: 110 },
  { title: '状态', key: 'status', width: 90 },
  { title: '进度', key: 'progress', width: 150 },
  { title: '文件名', dataIndex: 'file_name', ellipsis: true },
  { title: '创建时间', dataIndex: 'created_at', width: 160 },
  { title: '完成时间', dataIndex: 'completed_at', width: 160 },
  { title: '操作', key: 'action', width: 80 },
]

async function fetchList() {
  // 后端暂无导出任务列表接口，此页面为预留功能
  list.value = []
}

function handleReset() {
  Object.assign(searchParams, { report_type: undefined, status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

onMounted(fetchList)
</script>
