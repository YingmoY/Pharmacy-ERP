<template>
  <div class="trace-timeline">
    <a-timeline>
      <a-timeline-item
        v-for="log in logs"
        :key="log.id"
        :color="getTimelineColor(log.action_type)"
      >
        <template #dot>
          <component :is="getActionIcon(log.action_type)" style="font-size: 16px" />
        </template>
        <div class="timeline-content">
          <div class="action-row">
            <span class="action-label">{{ getActionLabel(log.action_type) }}</span>
            <span class="action-time">{{ formatTime(log.created_at) }}</span>
          </div>
          <div class="action-detail">
            <template v-if="log.from_status">
              <StatusTag type="trace" :value="log.from_status" />
              <span style="margin: 0 6px; color: #999">→</span>
            </template>
            <StatusTag type="trace" :value="log.to_status" />
          </div>
          <div v-if="log.from_location_code || log.to_location_code" class="action-location">
            <template v-if="log.from_location_code">
              货位：{{ log.from_location_code }}
              <span style="margin: 0 6px; color: #999">→</span>
            </template>
            <template v-if="log.to_location_code">
              {{ log.to_location_code }}
            </template>
          </div>
          <div class="action-operator">
            操作人：{{ log.operator_name }}
            <template v-if="log.related_no">
              &nbsp;·&nbsp; 关联单号：{{ log.related_no }}
            </template>
          </div>
          <div v-if="log.remark" class="action-remark">备注：{{ log.remark }}</div>
        </div>
      </a-timeline-item>
    </a-timeline>

    <a-empty v-if="logs.length === 0" description="暂无追溯记录" />
  </div>
</template>

<script setup lang="ts">
import {
  InboxOutlined,
  AppstoreOutlined,
  ShoppingCartOutlined,
  RollbackOutlined,
  FileSearchOutlined,
  SwapOutlined,
  MinusCircleOutlined,
} from '@ant-design/icons-vue'
import dayjs from 'dayjs'
import StatusTag from './StatusTag.vue'
import type { DrugTraceLog } from '@/types/inventory'

defineProps<{
  logs: DrugTraceLog[]
}>()

function getActionLabel(actionType: string): string {
  const map: Record<string, string> = {
    INBOUND: '入库扫码',
    SHELVING: '上架',
    SALE: '销售出库',
    RETURN: '退货入库',
    INVENTORY: '盘库操作',
    RELOCATION: '货位调拨',
    LOSS: '确认盘亏',
  }
  return map[actionType] ?? actionType
}

function getTimelineColor(actionType: string): string {
  const map: Record<string, string> = {
    INBOUND: 'blue',
    SHELVING: 'green',
    SALE: 'orange',
    RETURN: 'purple',
    INVENTORY: 'cyan',
    RELOCATION: '#722ed1',
    LOSS: 'red',
  }
  return map[actionType] ?? 'gray'
}

function getActionIcon(actionType: string) {
  const map: Record<string, unknown> = {
    INBOUND: InboxOutlined,
    SHELVING: AppstoreOutlined,
    SALE: ShoppingCartOutlined,
    RETURN: RollbackOutlined,
    INVENTORY: FileSearchOutlined,
    RELOCATION: SwapOutlined,
    LOSS: MinusCircleOutlined,
  }
  return map[actionType] ?? InboxOutlined
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}
</script>

<style scoped>
.trace-timeline {
  padding: 8px 0;
}

.timeline-content {
  padding: 0 0 8px;
}

.action-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 4px;
}

.action-label {
  font-weight: 600;
  font-size: 14px;
}

.action-time {
  font-size: 12px;
  color: #999;
}

.action-detail {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}

.action-location {
  font-size: 12px;
  color: #666;
  margin-bottom: 2px;
}

.action-operator {
  font-size: 12px;
  color: #666;
}

.action-remark {
  font-size: 12px;
  color: #999;
  margin-top: 2px;
  font-style: italic;
}
</style>
