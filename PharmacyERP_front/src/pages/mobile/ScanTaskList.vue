<template>
  <div style="background: #f5f5f5; min-height: 100vh; padding: 12px">
    <div style="margin-bottom: 12px">
      <a-select v-model:value="statusFilter" style="width: 100%" @change="fetchList">
        <a-select-option value="">全部任务</a-select-option>
        <a-select-option value="ASSIGNED">已分配</a-select-option>
        <a-select-option value="IN_PROGRESS">进行中</a-select-option>
        <a-select-option value="COMPLETED">已完成</a-select-option>
      </a-select>
    </div>

    <a-spin :spinning="loading">
      <div v-if="list.length === 0 && !loading" style="text-align: center; padding: 40px; color: #999">暂无任务</div>

      <div
        v-for="task in list"
        :key="task.id"
        style="background: #fff; border-radius: 8px; padding: 12px; margin-bottom: 10px; cursor: pointer"
        @click="$router.push(`/m/scan-tasks/${task.id}`)"
      >
        <div style="display: flex; justify-content: space-between; align-items: flex-start">
          <div>
            <div style="font-weight: 600; font-size: 15px">{{ task.task_no }}</div>
            <div style="color: #666; font-size: 13px; margin-top: 4px">{{ task.task_type }}</div>
          </div>
          <StatusTag type="scanTask" :value="task.status" />
        </div>
        <div style="margin-top: 8px; font-size: 12px; color: #999">
          {{ task.created_at }}
        </div>
        <div style="margin-top: 8px" v-if="task.status !== 'COMPLETED'">
          <a-button type="primary" size="small" block @click.stop="$router.push(`/m/scan-tasks/${task.id}`)">
            开始扫码
          </a-button>
        </div>
      </div>
    </a-spin>

    <div style="text-align: center; margin-top: 12px; color: #999; font-size: 12px" v-if="hasMore">
      <a @click="loadMore">加载更多</a>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getScanTaskList } from '@/api/scanTask'

const loading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const statusFilter = ref('')
const page = ref(1)
const hasMore = ref(false)

async function fetchList() {
  loading.value = true
  page.value = 1
  try {
    const res = await getScanTaskList({ page: 1, page_size: 20, status: statusFilter.value || undefined })
    list.value = res.list
    hasMore.value = res.total > 20
  } finally { loading.value = false }
}

async function loadMore() {
  page.value++
  const res = await getScanTaskList({ page: page.value, page_size: 20, status: statusFilter.value || undefined })
  list.value.push(...res.list)
  hasMore.value = list.value.length < res.total
}

onMounted(fetchList)
</script>
