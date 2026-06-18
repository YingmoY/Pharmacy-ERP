<template>
  <div class="page-content">
    <a-alert type="info" style="margin-bottom: 12px" message="权限字典由后端维护，此页面仅供查看。" show-icon banner />

    <a-card>
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="false"
        row-key="id"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'">
            <a-tag :color="typeColor(record.type)">{{ record.type }}</a-tag>
          </template>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getPermissionList } from '@/api/system'
import type { SysPermission } from '@/types/system'

const loading = ref(false)
const list = ref<SysPermission[]>([])

const typeColorMap: Record<string, string> = { MENU: 'blue', BUTTON: 'green', API: 'purple' }
function typeColor(type: string): string { return typeColorMap[type] ?? 'default' }

const columns = [
  { title: '权限名称', dataIndex: 'name', width: 200 },
  { title: '权限编码', dataIndex: 'code', ellipsis: true },
  { title: '类型', key: 'type', width: 80 },
  { title: '排序', dataIndex: 'sort_order', width: 70 },
]

async function fetchList() {
  loading.value = true
  try {
    list.value = await getPermissionList()
  } finally { loading.value = false }
}

onMounted(fetchList)
</script>
