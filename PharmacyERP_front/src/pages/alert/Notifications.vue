<template>
  <div class="page-content">
    <a-card>
      <template #extra>
        <a-button @click="markAllRead" :loading="markLoading">全部标为已读</a-button>
      </template>
      <a-list :data-source="list" :loading="loading" item-layout="horizontal">
        <template #renderItem="{ item }">
          <a-list-item>
            <a-list-item-meta>
              <template #avatar>
                <a-badge :dot="!item.is_read">
                  <a-avatar :style="{ backgroundColor: getIconBg(item.notification_type) }">
                    <template #icon>
                      <BellOutlined />
                    </template>
                  </a-avatar>
                </a-badge>
              </template>
              <template #title>
                <span :style="{ fontWeight: item.is_read ? 'normal' : '600' }">{{ item.title }}</span>
              </template>
              <template #description>
                <div>{{ item.content }}</div>
                <div style="font-size: 12px; color: #999; margin-top: 4px">{{ item.created_at }}</div>
              </template>
            </a-list-item-meta>
            <template #actions>
              <a v-if="!item.is_read" @click="markRead(item.id)">标为已读</a>
            </template>
          </a-list-item>
        </template>
      </a-list>
      <div style="text-align: center; margin-top: 12px">
        <a-pagination v-model:current="pagination.current" :total="pagination.total" :page-size="pagination.pageSize" @change="handlePageChange" show-total />
      </div>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { BellOutlined } from '@ant-design/icons-vue'
import { getNotificationList, markNotificationRead, markAllNotificationsRead } from '@/api/alert'
import type { Notification } from '@/types/alert'

const loading = ref(false)
const markLoading = ref(false)
const list = ref<Notification[]>([])
const pagination = reactive({ current: 1, pageSize: 20, total: 0 })

const typeColorMap: Record<string, string> = {
  ALERT: '#ff4d4f',
  SYSTEM: '#1677ff',
  ORDER: '#52c41a',
  TASK: '#fa8c16',
}

function getIconBg(type: string): string {
  return typeColorMap[type] ?? '#1677ff'
}

async function markRead(id: number) {
  await markNotificationRead(id)
  const item = list.value.find(n => n.id === id)
  if (item) item.is_read = true
}

async function markAllRead() {
  markLoading.value = true
  try {
    await markAllNotificationsRead()
    message.success('全部已读')
    fetchList()
  } finally { markLoading.value = false }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getNotificationList({ page: pagination.current, page_size: pagination.pageSize })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handlePageChange(page: number) {
  pagination.current = page
  fetchList()
}

onMounted(fetchList)
</script>
