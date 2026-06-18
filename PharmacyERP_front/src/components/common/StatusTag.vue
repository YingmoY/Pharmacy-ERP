<template>
  <a-tag :color="tagColor">{{ tagLabel }}</a-tag>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  type: string   // 状态类型标识，如 'inbound', 'sales', 'trace', 'review', 'alert'
  value: string  // 状态值
}>()

// 各类状态标签配置
const statusConfig: Record<string, Record<string, { label: string; color: string }>> = {
  // 入库单状态
  inbound: {
    DRAFT: { label: '草稿', color: 'default' },
    PENDING_CONFIRM: { label: '待扫码确认', color: 'blue' },
    COMPLETED: { label: '已完成', color: 'green' },
    CANCELLED: { label: '已取消', color: 'red' },
  },
  // 销售单状态
  sales: {
    PENDING: { label: '待结算', color: 'orange' },
    PENDING_REVIEW: { label: '待药师审核', color: 'blue' },
    APPROVED: { label: '审核通过', color: 'cyan' },
    COMPLETED: { label: '已完成', color: 'green' },
    PARTIALLY_REFUNDED: { label: '部分退货', color: 'purple' },
    REFUNDED: { label: '已退货', color: 'magenta' },
    CANCELLED: { label: '已取消', color: 'red' },
  },
  // 追溯码库存状态
  trace: {
    PENDING: { label: '待上架', color: 'orange' },
    IN_STOCK: { label: '在库', color: 'green' },
    SOLD: { label: '已售', color: 'blue' },
    MISPLACED: { label: '错架', color: 'volcano' },
    LOSS_CANDIDATE: { label: '盘亏候选', color: 'red' },
    LOST: { label: '已盘亏', color: 'gray' },
  },
  // 药师审核状态
  review: {
    PENDING: { label: '待审核', color: 'orange' },
    APPROVED: { label: '审核通过', color: 'green' },
    REJECTED: { label: '审核驳回', color: 'red' },
    CANCELLED: { label: '已取消', color: 'default' },
  },
  // 预警状态
  alert: {
    ACTIVE: { label: '待处理', color: 'red' },
    RESOLVED: { label: '已处理', color: 'green' },
    IGNORED: { label: '已忽略', color: 'default' },
  },
  // 预警类型
  alertType: {
    NEAR_EXPIRE: { label: '近效期', color: 'orange' },
    LOW_STOCK: { label: '低库存', color: 'red' },
    LOSS_CANDIDATE: { label: '盘亏候选', color: 'volcano' },
    MISPLACED: { label: '错架', color: 'purple' },
  },
  // 预警级别
  alertLevel: {
    HIGH: { label: '高', color: 'red' },
    MEDIUM: { label: '中', color: 'orange' },
    LOW: { label: '低', color: 'blue' },
  },
  // 盘库任务状态
  inventoryTask: {
    PENDING: { label: '待开始', color: 'default' },
    IN_PROGRESS: { label: '进行中', color: 'blue' },
    COMPLETED: { label: '已完成', color: 'green' },
    CANCELLED: { label: '已取消', color: 'red' },
  },
  // 扫码任务状态
  scanTask: {
    PENDING: { label: '待开始', color: 'default' },
    IN_PROGRESS: { label: '进行中', color: 'blue' },
    COMPLETED: { label: '已完成', color: 'green' },
    CANCELLED: { label: '已取消', color: 'red' },
  },
  // AI 发票识别状态
  invoice: {
    PENDING: { label: '待处理', color: 'default' },
    PROCESSING: { label: '识别中', color: 'blue' },
    COMPLETED: { label: '识别完成', color: 'green' },
    FAILED: { label: '识别失败', color: 'red' },
  },
  // 导出任务状态
  export: {
    PENDING: { label: '待执行', color: 'default' },
    RUNNING: { label: '执行中', color: 'blue' },
    SUCCESS: { label: '成功', color: 'green' },
    FAILED: { label: '失败', color: 'red' },
  },
  // 退货状态
  refund: {
    NONE: { label: '未退货', color: 'default' },
    REFUNDED: { label: '已退货', color: 'red' },
  },
  // 启用状态
  enable: {
    '1': { label: '启用', color: 'green' },
    '0': { label: '停用', color: 'default' },
  },
  // 扫码结果
  scanResult: {
    NORMAL: { label: '正常', color: 'green' },
    MISPLACED_FOUND: { label: '错架', color: 'volcano' },
    UNEXPECTED: { label: '异常', color: 'red' },
    DUPLICATE: { label: '重复', color: 'orange' },
    SUCCESS: { label: '成功', color: 'green' },
    INVALID: { label: '无效', color: 'red' },
    STATUS_ERROR: { label: '状态错误', color: 'orange' },
  },
}

const config = computed(() => {
  const typeConfig = statusConfig[props.type]
  if (!typeConfig) return { label: props.value, color: 'default' }
  return typeConfig[props.value] ?? { label: props.value, color: 'default' }
})

const tagLabel = computed(() => config.value.label)
const tagColor = computed(() => config.value.color)
</script>
