<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <div class="review-layout">
        <!-- 左侧：销售单详情 -->
        <div class="review-left">
          <a-card title="销售单信息" style="margin-bottom: 16px">
            <a-descriptions :column="2" size="small" v-if="reviewData">
              <a-descriptions-item label="销售单号">{{ reviewData.order?.order_no }}</a-descriptions-item>
              <a-descriptions-item label="提交人">{{ reviewData.submitter_name }}</a-descriptions-item>
              <a-descriptions-item label="提交时间">{{ reviewData.submitted_at }}</a-descriptions-item>
              <a-descriptions-item label="总金额">¥{{ reviewData.order?.total_amount }}</a-descriptions-item>
            </a-descriptions>
          </a-card>

          <a-card title="药品明细">
            <a-table
              :columns="itemColumns"
              :data-source="reviewData?.order?.items ?? []"
              :pagination="false"
              size="small"
              row-key="id"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'drug_name'">
                  <span>{{ record.drug_name }}</span>
                  <a-tag v-if="isPrescription(record)" color="red" size="small" style="margin-left: 4px">处方药</a-tag>
                </template>
                <template v-if="column.key === 'expire_date'">
                  <span :style="{ color: isNearExpire(record.expire_date) ? '#ff4d4f' : 'inherit' }">
                    {{ record.expire_date }}
                  </span>
                </template>
                <template v-if="column.key === 'price'">
                  ¥{{ record.unit_price }}
                </template>
              </template>
            </a-table>
          </a-card>
        </div>

        <!-- 右侧：审核操作 -->
        <div class="review-right">
          <a-card title="审核操作">
            <!-- 当前状态 -->
            <div style="margin-bottom: 16px">
              <span style="color: #666">审核状态：</span>
              <StatusTag type="review" :value="reviewData?.status ?? ''" />
            </div>

            <!-- 已审核时显示结果 -->
            <template v-if="reviewData?.status !== 'PENDING'">
              <a-descriptions :column="1" size="small">
                <a-descriptions-item label="审核人">{{ reviewData?.pharmacist_name }}</a-descriptions-item>
                <a-descriptions-item label="审核时间">{{ reviewData?.reviewed_at }}</a-descriptions-item>
                <a-descriptions-item label="审核意见">{{ reviewData?.review_opinion || '无' }}</a-descriptions-item>
              </a-descriptions>
            </template>

            <!-- 待审核时显示操作 -->
            <template v-else>
              <a-form layout="vertical">
                <a-form-item label="审核意见">
                  <a-textarea
                    v-model:value="opinion"
                    :rows="4"
                    placeholder="请输入审核意见（驳回时必填）"
                  />
                </a-form-item>
                <a-form-item>
                  <a-space direction="vertical" style="width: 100%">
                    <PermissionButton
                      type="primary"
                      block
                      size="large"
                      permission="pharmacist.review.approve"
                      :loading="approveLoading"
                      @click="handleApprove"
                    >
                      <CheckCircleOutlined /> 审核通过
                    </PermissionButton>
                    <PermissionButton
                      danger
                      block
                      size="large"
                      permission="pharmacist.review.reject"
                      :loading="rejectLoading"
                      @click="handleReject"
                    >
                      <CloseCircleOutlined /> 审核驳回
                    </PermissionButton>
                  </a-space>
                </a-form-item>
              </a-form>

              <a-alert
                type="warning"
                message="驳回后销售单将直接取消，所有预占将被释放，此操作不可撤销！"
                banner
              />
            </template>
          </a-card>
        </div>
      </div>
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import PermissionButton from '@/components/common/PermissionButton.vue'
import { getReviewDetail, approveReview, rejectReview } from '@/api/pharmacist'
import type { PharmacistReview, SalesOrder } from '@/types/sales'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const reviewId = Number(route.params.id)

const loading = ref(false)
const approveLoading = ref(false)
const rejectLoading = ref(false)
const opinion = ref('')
const reviewData = ref<(PharmacistReview & { order: SalesOrder }) | null>(null)

const itemColumns = [
  { title: '药品名称', key: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 120 },
  { title: '追溯码', dataIndex: 'trace_code', width: 180 },
  { title: '批号', dataIndex: 'batch_number', width: 120 },
  { title: '有效期', key: 'expire_date', width: 100 },
  { title: '单价', key: 'price', width: 80 },
]

function isPrescription(record: Record<string, unknown>): boolean {
  return !!(record as { is_prescription?: boolean }).is_prescription
}

function isNearExpire(expireDate: string): boolean {
  return dayjs(expireDate).diff(dayjs(), 'day') <= 30
}

async function handleApprove() {
  Modal.confirm({
    title: '确认审核通过',
    content: '确认处方合规，通过审核后销售单可以进行结算。',
    okText: '确认通过',
    onOk: async () => {
      approveLoading.value = true
      try {
        await approveReview(reviewId, opinion.value)
        message.success('审核已通过')
        router.push('/pharmacist/reviews')
      } finally {
        approveLoading.value = false
      }
    },
  })
}

async function handleReject() {
  if (!opinion.value.trim()) {
    message.warning('驳回时必须填写审核意见')
    return
  }
  Modal.confirm({
    title: '确认驳回审核',
    content: '驳回后销售单将直接变为已取消状态，所有预占将被释放。此操作不可撤销！',
    okText: '确认驳回',
    okType: 'danger',
    onOk: async () => {
      rejectLoading.value = true
      try {
        await rejectReview(reviewId, opinion.value)
        message.success('审核已驳回，销售单已取消')
        router.push('/pharmacist/reviews')
      } finally {
        rejectLoading.value = false
      }
    },
  })
}

onMounted(async () => {
  loading.value = true
  try {
    reviewData.value = await getReviewDetail(reviewId)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.review-layout {
  display: grid;
  grid-template-columns: 1fr 360px;
  gap: 16px;
}

.review-left,
.review-right {
  height: fit-content;
}
</style>
