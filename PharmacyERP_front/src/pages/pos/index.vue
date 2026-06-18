<template>
  <div class="page-content">
    <div class="pos-layout">
      <!-- 左侧：药品搜索区 -->
      <div class="pos-search-panel">
        <a-card title="药品搜索" size="small" style="height: 100%; display: flex; flex-direction: column">
          <template #extra>
            <a-radio-group v-model:value="searchMode" size="small" button-style="solid">
              <a-radio-button value="ai">智能搜索</a-radio-button>
              <a-radio-button value="symptom">AI推荐</a-radio-button>
              <a-radio-button value="scan">追溯码</a-radio-button>
            </a-radio-group>
          </template>

          <!-- AI 智能搜索 -->
          <div v-if="searchMode === 'ai'">
            <DrugSearchBox @select="handleDrugSelect" />
          </div>

          <!-- AI 症状推荐 -->
          <div v-else-if="searchMode === 'symptom'">
            <AiDrugRecommend @select="handleDrugSelect" />
          </div>

          <!-- 追溯码扫码 -->
          <div v-else>
            <ScanInput
              placeholder="扫描追溯码"
              :auto-focus="true"
              @submit="handleTraceScan"
            />
          </div>

          <!-- 最近添加提示 -->
          <div v-if="lastAdded" class="last-added-tip">
            <a-alert
              :type="lastAdded.success ? 'success' : 'error'"
              :message="lastAdded.message"
              banner
              closable
            />
          </div>
        </a-card>
      </div>

      <!-- 中间：销售明细 -->
      <div class="pos-cart-panel">
        <a-card
          size="small"
          style="height: 100%; display: flex; flex-direction: column"
        >
          <template #title>
            <span>销售明细</span>
            <a-tag v-if="currentOrder" style="margin-left: 8px" color="blue">
              {{ currentOrder.order_no }}
            </a-tag>
          </template>
          <template #extra>
            <a-space>
              <a-button
                size="small"
                danger
                :disabled="!currentOrder || cartItems.length === 0"
                @click="handleCancelOrder"
              >
                取消订单
              </a-button>
            </a-space>
          </template>

          <!-- 明细表格 -->
          <div style="flex: 1; overflow: auto">
            <a-table
              :columns="cartColumns"
              :data-source="cartItems"
              :pagination="false"
              size="small"
              row-key="id"
              :scroll="{ y: 'calc(100vh - 320px)' }"
            >
              <template #bodyCell="{ column, record }">
                <template v-if="column.key === 'drug_info'">
                  <div style="font-weight: 500">{{ record.drug_name }}</div>
                  <div style="font-size: 12px; color: #666">
                    {{ record.specification }} · {{ record.manufacturer }}
                  </div>
                </template>
                <template v-if="column.key === 'trace_info'">
                  <div style="font-size: 12px">{{ record.trace_code || '未分配' }}</div>
                  <div style="font-size: 11px; color: #999">
                    {{ record.batch_number }} 效期:{{ record.expire_date }}
                  </div>
                </template>
                <template v-if="column.key === 'price'">
                  <span style="color: #e6550d">¥{{ record.unit_price }}</span>
                </template>
                <template v-if="column.key === 'action'">
                  <a-button
                    type="link"
                    danger
                    size="small"
                    :disabled="currentOrder?.status !== 'PENDING'"
                    @click="handleRemoveItem(record.id)"
                  >
                    移除
                  </a-button>
                </template>
              </template>

              <template #summary>
                <a-table-summary fixed>
                  <a-table-summary-row>
                    <a-table-summary-cell :col-span="5" style="text-align: right; font-weight: 600">
                      合计：¥{{ totalAmount }}
                    </a-table-summary-cell>
                  </a-table-summary-row>
                </a-table-summary>
              </template>
            </a-table>
          </div>

          <!-- 处方药审核提示 -->
          <a-alert
            v-if="hasPrescriptionDrug"
            type="warning"
            message="当前订单包含处方药，结算后将自动提交药师审核"
            banner
            style="margin-top: 8px"
          />
        </a-card>
      </div>

      <!-- 右侧：结算区 -->
      <div class="pos-settle-panel">
        <a-card title="结算信息" size="small" style="height: 100%">
          <!-- 订单状态 -->
          <div class="settle-item">
            <span class="settle-label">订单状态</span>
            <StatusTag type="sales" :value="currentOrder?.status ?? 'PENDING'" />
          </div>

          <a-divider style="margin: 12px 0" />

          <!-- 金额信息 -->
          <div class="settle-item">
            <span class="settle-label">应收金额</span>
            <span class="settle-amount">¥{{ totalAmount }}</span>
          </div>
          <template v-if="paymentMethod === 'MEDICARE' && medicarePreview">
            <div class="settle-item">
              <span class="settle-label">医保统筹支付</span>
              <span style="color: #389e0d">-¥{{ medicarePreview.fund_pay.toFixed(2) }}</span>
            </div>
            <div v-if="medicarePreview.acct_pay > 0" class="settle-item">
              <span class="settle-label">个人账户支付</span>
              <span style="color: #389e0d">-¥{{ medicarePreview.acct_pay.toFixed(2) }}</span>
            </div>
          </template>
          <div v-else class="settle-item">
            <span class="settle-label">优惠金额</span>
            <span>¥{{ discountAmount }}</span>
          </div>
          <div class="settle-item">
            <span class="settle-label fw-bold">
              {{ paymentMethod === 'MEDICARE' && medicarePreview ? '个人现金' : '实收金额' }}
            </span>
            <span class="settle-actual">¥{{ actualAmount }}</span>
          </div>

          <a-divider style="margin: 12px 0" />

          <!-- 支付方式 -->
          <div class="settle-label" style="margin-bottom: 8px">支付方式</div>
          <a-radio-group
            v-model:value="paymentMethod"
            :disabled="!orderReady"
            style="width: 100%"
          >
            <a-space direction="vertical" style="width: 100%">
              <a-radio value="CASH">现金</a-radio>
              <a-radio value="ALIPAY">支付宝</a-radio>
              <a-radio value="WECHAT">微信支付</a-radio>
              <a-radio value="BANK_CARD">银行卡</a-radio>
              <a-radio value="MEDICARE">医保</a-radio>
            </a-space>
          </a-radio-group>

          <!-- 医保参数 (选择医保支付时显示) -->
          <template v-if="paymentMethod === 'MEDICARE'">
            <a-divider style="margin: 10px 0; font-size: 12px">医保参数</a-divider>
            <a-form size="small" layout="vertical" :disabled="!orderReady">
              <a-form-item label="就诊类型" style="margin-bottom: 8px">
                <a-select v-model:value="medicareOptions.med_type" style="width: 100%">
                  <a-select-option value="41">门诊</a-select-option>
                  <a-select-option value="11">住院</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="险种类型" style="margin-bottom: 8px">
                <a-select v-model:value="medicareOptions.insutype" style="width: 100%">
                  <a-select-option value="310">职工医保</a-select-option>
                  <a-select-option value="390">居民医保</a-select-option>
                </a-select>
              </a-form-item>
              <a-form-item label="使用账户余额" style="margin-bottom: 8px">
                <a-switch v-model:checked="medicareOptions.useAcctBalance" />
              </a-form-item>
              <a-form-item label="医保凭证号" style="margin-bottom: 8px">
                <a-input-search
                  v-model:value="medicareOptions.cert_no"
                  placeholder="输入身份证号或医保卡号"
                  enter-button="查询医保"
                  :loading="medicarePreviewLoading"
                  :disabled="!orderReady"
                  @search="handleMedicareQuery"
                />
              </a-form-item>
            </a-form>

            <!-- 医保查询结果 -->
            <template v-if="medicarePreview">
              <a-descriptions
                size="small"
                :column="1"
                bordered
                style="margin-top: 4px; font-size: 12px"
              >
                <a-descriptions-item label="参保人">
                  {{ medicarePreview.psn_name }}
                </a-descriptions-item>
                <a-descriptions-item label="统筹支付">
                  <span style="color: #389e0d">¥{{ medicarePreview.fund_pay.toFixed(2) }}</span>
                </a-descriptions-item>
                <a-descriptions-item v-if="medicarePreview.acct_pay > 0" label="账户支付">
                  <span style="color: #389e0d">¥{{ medicarePreview.acct_pay.toFixed(2) }}</span>
                </a-descriptions-item>
                <a-descriptions-item label="个人现金">
                  <span style="color: #e6550d; font-weight: 600">
                    ¥{{ medicarePreview.personal_cash.toFixed(2) }}
                  </span>
                </a-descriptions-item>
              </a-descriptions>
            </template>
            <a-alert
              v-else
              type="info"
              message="请输入医保凭证号并点击&quot;查询医保&quot;，验证参保信息后方可结算"
              banner
              style="margin-top: 4px; font-size: 12px"
            />
          </template>

          <a-divider style="margin: 12px 0" />

          <!-- 操作按钮 -->
          <div class="settle-actions">
            <!-- 待药师审核时显示状态 -->
            <a-alert
              v-if="currentOrder?.status === 'PENDING_REVIEW'"
              type="info"
              message="等待药师审核，审核通过后可结算"
              style="margin-bottom: 12px"
            />

            <a-button
              type="primary"
              block
              size="large"
              :disabled="!canSettle"
              :loading="settleLoading"
              @click="handleSettle"
            >
              {{ settleButtonText }}
            </a-button>

            <a-button
              block
              size="large"
              style="margin-top: 8px"
              @click="handleNewOrder"
            >
              新建订单
            </a-button>
          </div>
        </a-card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted } from 'vue'
import { message, Modal } from 'ant-design-vue'
import { useRoute } from 'vue-router'
import DrugSearchBox from '@/components/common/DrugSearchBox.vue'
import AiDrugRecommend from '@/components/common/AiDrugRecommend.vue'
import ScanInput from '@/components/common/ScanInput.vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  createSalesOrder,
  addSalesOrderItem,
  deleteSalesOrderItem,
  settleSalesOrder,
  cancelSalesOrder,
  getSalesOrderDetail,
  getMedicarePreview,
} from '@/api/sales'
import { validateTraceCode } from '@/api/inventory'
import type { SalesOrder, SalesOrderItem, PaymentMethod, MedicarePreviewResponse } from '@/types/sales'
import type { DrugSearchResult } from '@/types/drug'

const route = useRoute()
const searchMode = ref<'ai' | 'symptom' | 'scan'>('ai')
const currentOrder = ref<SalesOrder | null>(null)
const cartItems = ref<SalesOrderItem[]>([])
const paymentMethod = ref<PaymentMethod>('CASH')
const settleLoading = ref(false)
const lastAdded = ref<{ success: boolean; message: string } | null>(null)

const medicareOptions = ref({
  med_type: '41',
  insutype: '310',
  useAcctBalance: true,
  cert_no: '',
})
const medicarePreview = ref<MedicarePreviewResponse | null>(null)
const medicarePreviewLoading = ref(false)

// 是否包含处方药
const hasPrescriptionDrug = computed(() =>
  cartItems.value.some((item) => {
    // 需要从 order 中获取 need_audit 信息
    return currentOrder.value?.need_audit
  }),
)

// 合计金额
const totalAmount = computed(() => {
  const sum = cartItems.value.reduce((acc, item) => acc + parseFloat(item.unit_price), 0)
  return sum.toFixed(2)
})

const discountAmount = computed(() => '0.00')
const actualAmount = computed(() => {
  if (paymentMethod.value === 'MEDICARE' && medicarePreview.value) {
    return medicarePreview.value.personal_cash.toFixed(2)
  }
  return (parseFloat(totalAmount.value) - parseFloat(discountAmount.value)).toFixed(2)
})

// 订单是否处于可操作状态（用于控制支付方式、参数表单等 UI 的可用性）
const orderReady = computed(() => {
  if (!currentOrder.value) return false
  if (!['PENDING', 'APPROVED'].includes(currentOrder.value.status)) return false
  if (cartItems.value.length === 0) return false
  return true
})

// 结算按钮是否可点击（在 orderReady 基础上，医保方式还要求先查询）
const canSettle = computed(() => {
  if (!orderReady.value) return false
  if (paymentMethod.value === 'MEDICARE' && !medicarePreview.value) return false
  return true
})

const settleButtonText = computed(() => {
  if (!currentOrder.value) return '结算'
  if (currentOrder.value.status === 'PENDING_REVIEW') return '等待审核'
  if (paymentMethod.value === 'MEDICARE' && !medicarePreview.value) return '请先查询医保'
  return `确认结算 ¥${actualAmount.value}`
})

// 表格列定义
const cartColumns = [
  { title: '药品信息', key: 'drug_info', ellipsis: true },
  { title: '追溯码 / 批次', key: 'trace_info', width: 160 },
  { title: '单价', key: 'price', width: 80 },
  { title: '操作', key: 'action', width: 60 },
]

// 从 AI 搜索结果添加药品
async function handleDrugSelect(drug: DrugSearchResult) {
  try {
    if (!currentOrder.value) {
      // 自动创建订单
      const order = await createSalesOrder({
        items: [{ drug_id: drug.drug_id }],
      })
      currentOrder.value = order
      cartItems.value = order.items ?? []
    } else {
      // 添加到已有订单
      const item = await addSalesOrderItem(currentOrder.value.id, {
        drug_id: drug.drug_id,
      })
      cartItems.value.push(item)
      await refreshOrder()
    }
    showLastAdded(true, `已添加：${drug.common_name}`)
  } catch {
    showLastAdded(false, '添加失败，请检查库存')
  }
}

// 扫码追溯码添加
async function handleTraceScan(traceCode: string) {
  try {
    // 先验证追溯码
    const validation = await validateTraceCode({ trace_code: traceCode })
    if (!validation.is_available) {
      showLastAdded(false, validation.reason || '追溯码不可售')
      return
    }
    if (!validation.drug_id) {
      showLastAdded(false, '未找到药品信息')
      return
    }

    if (!currentOrder.value) {
      const order = await createSalesOrder({
        items: [{ drug_id: validation.drug_id, trace_code: traceCode }],
      })
      currentOrder.value = order
      cartItems.value = order.items ?? []
    } else {
      const item = await addSalesOrderItem(currentOrder.value.id, {
        drug_id: validation.drug_id!,
        trace_code: traceCode,
      })
      cartItems.value.push(item)
      await refreshOrder()
    }
    showLastAdded(true, `已扫码：${validation.drug_name}`)
  } catch {
    showLastAdded(false, '追溯码无效或已被占用')
  }
}

// 移除明细
async function handleRemoveItem(itemId: number) {
  if (!currentOrder.value) return
  try {
    await deleteSalesOrderItem(currentOrder.value.id, itemId)
    cartItems.value = cartItems.value.filter((i) => i.id !== itemId)
    await refreshOrder()
  } catch {
    // 错误已统一处理
  }
}

// 结算
async function handleSettle() {
  if (!currentOrder.value) return

  const isMedicare = paymentMethod.value === 'MEDICARE'
  const payMethodLabel: Record<string, string> = {
    CASH: '现金', ALIPAY: '支付宝', WECHAT: '微信支付', BANK_CARD: '银行卡', MEDICARE: '医保',
  }
  const insuLabel: Record<string, string> = { '310': '职工医保', '390': '居民医保' }
  const medTypeLabel: Record<string, string> = { '41': '门诊', '11': '住院' }

  let confirmContent: string
  if (isMedicare && medicarePreview.value) {
    const preview = medicarePreview.value
    confirmContent =
      `参保人：${preview.psn_name}\n` +
      `险种：${insuLabel[medicareOptions.value.insutype] ?? medicareOptions.value.insutype}，` +
      `就诊类型：${medTypeLabel[medicareOptions.value.med_type] ?? medicareOptions.value.med_type}\n` +
      `医保支付：¥${(preview.fund_pay + preview.acct_pay).toFixed(2)}，` +
      `个人现金：¥${preview.personal_cash.toFixed(2)}`
  } else {
    confirmContent = `确认收款 ¥${actualAmount.value}，支付方式：${payMethodLabel[paymentMethod.value] ?? paymentMethod.value}？`
  }

  Modal.confirm({
    title: '确认结算',
    content: confirmContent,
    okText: '确认结算',
    cancelText: '取消',
    onOk: async () => {
      settleLoading.value = true
      try {
        await settleSalesOrder(currentOrder.value!.id, {
          payment_method: paymentMethod.value,
          actual_amount: actualAmount.value,
          ...(isMedicare && medicarePreview.value && {
            use_medicare: true,
            med_type: medicareOptions.value.med_type,
            insutype: medicareOptions.value.insutype,
            acct_used_flag: medicareOptions.value.useAcctBalance ? '1' : '0',
            mdtrt_cert_no: medicareOptions.value.cert_no,
            psn_no: medicarePreview.value.psn_no,
            psn_name: medicarePreview.value.psn_name,
          }),
        })
        message.success('结算成功！')
        handleNewOrder()
      } finally {
        settleLoading.value = false
      }
    },
  })
}

// 取消订单
function handleCancelOrder() {
  if (!currentOrder.value) return
  Modal.confirm({
    title: '确认取消订单',
    content: '取消后已预占的追溯码将被释放，确认取消吗？',
    okText: '确认取消',
    okType: 'danger',
    cancelText: '保留订单',
    onOk: async () => {
      await cancelSalesOrder(currentOrder.value!.id)
      message.success('订单已取消')
      handleNewOrder()
    },
  })
}

// 医保查询
async function handleMedicareQuery() {
  if (!currentOrder.value) {
    message.warning('请先添加药品')
    return
  }
  if (!medicareOptions.value.cert_no.trim()) {
    message.warning('请输入医保凭证号')
    return
  }
  medicarePreviewLoading.value = true
  medicarePreview.value = null
  try {
    const preview = await getMedicarePreview(currentOrder.value.id, {
      mdtrt_cert_no: medicareOptions.value.cert_no.trim(),
      med_type: medicareOptions.value.med_type,
      insutype: medicareOptions.value.insutype,
      acct_used_flag: medicareOptions.value.useAcctBalance ? '1' : '0',
    })
    medicarePreview.value = preview
    message.success(`已查询到参保人：${preview.psn_name}`)
  } catch {
    message.error('医保查询失败，请检查凭证号是否正确')
  } finally {
    medicarePreviewLoading.value = false
  }
}

// 新建订单
function handleNewOrder() {
  currentOrder.value = null
  cartItems.value = []
  lastAdded.value = null
  paymentMethod.value = 'CASH'
  medicareOptions.value = { med_type: '41', insutype: '310', useAcctBalance: true, cert_no: '' }
  medicarePreview.value = null
}

// 清除医保预览：当关键参数变化时需重新查询
watch(
  [
    () => medicareOptions.value.cert_no,
    () => medicareOptions.value.insutype,
    () => medicareOptions.value.med_type,
    () => medicareOptions.value.useAcctBalance,
  ],
  () => { medicarePreview.value = null },
)
watch(paymentMethod, () => {
  if (paymentMethod.value !== 'MEDICARE') medicarePreview.value = null
})

async function refreshOrder() {
  if (!currentOrder.value) return
  currentOrder.value = await getSalesOrderDetail(currentOrder.value.id)
}

function showLastAdded(success: boolean, msg: string) {
  lastAdded.value = { success, message: msg }
  setTimeout(() => {
    lastAdded.value = null
  }, 3000)
}

// Load an existing order when navigated from the sales order list (e.g. ?orderId=123)
onMounted(async () => {
  const orderId = route.query.orderId
  if (orderId) {
    try {
      const order = await getSalesOrderDetail(Number(orderId))
      currentOrder.value = order
      cartItems.value = order.items ?? []
    } catch {
      // ignore — start fresh
    }
  }
})
</script>

<style scoped>
.pos-layout {
  display: grid;
  grid-template-columns: 320px 1fr 280px;
  gap: 16px;
  height: calc(100vh - 112px);
}

.pos-search-panel,
.pos-cart-panel,
.pos-settle-panel {
  height: 100%;
  overflow: hidden;
}

.settle-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.settle-label {
  color: #666;
  font-size: 14px;
}

.settle-amount {
  font-size: 16px;
  font-weight: 500;
}

.settle-actual {
  font-size: 24px;
  font-weight: 700;
  color: #e6550d;
}

.settle-actions {
  position: absolute;
  bottom: 16px;
  left: 16px;
  right: 16px;
}

.last-added-tip {
  margin-top: 12px;
}
</style>
