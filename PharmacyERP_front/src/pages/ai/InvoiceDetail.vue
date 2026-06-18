<template>
  <div class="page-content">
    <a-spin :spinning="loading">
      <template v-if="invoice">
        <div class="status-card" style="margin-bottom: 16px">
          <a-row align="middle" justify="space-between">
            <a-col>
              <div style="font-size: 18px; font-weight: 600">发票号：{{ invoice.result?.invoice_no ?? '—' }}</div>
              <div style="margin-top: 4px">
                <StatusTag type="invoice" :value="invoice.status" />
                <span style="margin-left: 12px; font-size: 13px; color: #666">
                  AI置信度：{{ Math.round((invoice.result?.confidence ?? 0) * 100) }}%
                </span>
              </div>
            </a-col>
            <a-col>
              <a-space>
                <a-button v-if="invoice.status === 'COMPLETED'" type="primary" @click="openConvertModal">
                  转为入库单
                </a-button>
                <a-button @click="$router.push('/ai/invoices')">返回列表</a-button>
              </a-space>
            </a-col>
          </a-row>
        </div>

        <a-tabs v-model:activeKey="activeTab">
          <a-tab-pane key="info" tab="发票信息">
            <a-row :gutter="16">
              <a-col :span="14">
                <a-descriptions bordered :column="2" size="small">
                  <a-descriptions-item label="供应商">{{ invoice.result?.recognized_supplier_name }}</a-descriptions-item>
                  <a-descriptions-item label="发票金额">¥{{ invoice.result?.total_amount }}</a-descriptions-item>
                  <a-descriptions-item label="发票日期">{{ invoice.result?.invoice_date }}</a-descriptions-item>
                  <a-descriptions-item label="识别药品数">{{ invoice.result?.items?.length ?? 0 }}</a-descriptions-item>
                  <a-descriptions-item label="上传时间" :span="2">{{ invoice.created_at }}</a-descriptions-item>
                </a-descriptions>

                <a-divider>识别药品明细</a-divider>
                <a-table
                  :columns="drugColumns"
                  :data-source="invoice.result?.items ?? []"
                  :pagination="false"
                  row-key="row_index"
                  size="small"
                >
                  <template #bodyCell="{ column, record }">
                    <template v-if="column.key === 'confidence'">
                      <a-progress
                        :percent="Math.round((record.confidence ?? 0) * 100)"
                        size="small"
                        :status="(record.confidence ?? 0) >= 0.9 ? 'success' : 'active'"
                      />
                    </template>
                    <template v-if="column.key === 'is_matched'">
                      <a-tag :color="record.matched_drug_id ? 'green' : 'orange'">
                        {{ record.matched_drug_id ? '已匹配' : '未匹配' }}
                      </a-tag>
                    </template>
                  </template>
                </a-table>
              </a-col>
              <a-col :span="10">
                <a-card title="供应商候选" size="small">
                  <a-empty v-if="!(invoice.result?.supplier_candidates?.length)" description="无候选供应商" />
                  <a-list v-else :data-source="invoice.result?.supplier_candidates ?? []" size="small">
                    <template #renderItem="{ item: sc }">
                      <a-list-item>
                        <a-list-item-meta
                          :title="sc.name"
                          :description="`置信度: ${Math.round(sc.confidence * 100)}%`"
                        />
                      </a-list-item>
                    </template>
                  </a-list>
                </a-card>
              </a-col>
            </a-row>
          </a-tab-pane>
        </a-tabs>
      </template>
    </a-spin>

    <!-- ===== 转为入库单 - 编辑确认弹窗 ===== -->
    <a-modal
      v-model:open="convertModalOpen"
      title="确认并编辑识别结果"
      width="980px"
      :footer="null"
      destroy-on-close
    >
      <!-- 供应商 & 发票号 -->
      <a-row :gutter="12" style="margin-bottom: 12px">
        <a-col :span="15">
          <div class="modal-form-row">
            <span class="modal-label">供应商 <span class="required">*</span></span>
            <a-select
              v-model:value="editSupplierId"
              show-search
              :loading="supplierSearching"
              :options="editSupplierOptions"
              :filter-option="false"
              placeholder="搜索供应商名称..."
              style="flex: 1"
              @search="debouncedSearchSupplier"
            />
          </div>
        </a-col>
        <a-col :span="9">
          <div class="modal-form-row">
            <span class="modal-label">发票号</span>
            <a-input v-model:value="editInvoiceNo" placeholder="发票号" style="flex: 1" />
          </div>
        </a-col>
      </a-row>

      <a-alert
        v-if="!editSupplierId"
        type="warning"
        message="请选择供应商，否则无法创建入库单"
        show-icon
        style="margin-bottom: 10px"
      />
      <a-alert
        v-if="unmatchedCount > 0"
        type="info"
        :message="`${unmatchedCount} 个药品未匹配 — 请在「匹配药品」列搜索并选择对应药品，或取消勾选跳过`"
        show-icon
        style="margin-bottom: 10px"
      />

      <!-- 可编辑明细表格 -->
      <a-table
        :data-source="editItems"
        :columns="editColumns"
        row-key="row_index"
        size="small"
        :pagination="false"
        :scroll="{ x: 860, y: 400 }"
      >
        <template #bodyCell="{ column, record }">
          <!-- 包含勾选 -->
          <template v-if="column.key === 'include'">
            <a-checkbox
              v-model:checked="record.include"
              :disabled="!record.matched_drug_id"
            />
          </template>

          <!-- AI识别名称 -->
          <template v-else-if="column.key === 'ai_name'">
            <a-tooltip :title="record.ai_name">
              <span :style="{ color: record.matched_drug_id ? 'inherit' : '#faad14', fontSize: '12px' }">
                {{ record.ai_name || '—' }}
              </span>
            </a-tooltip>
          </template>

          <!-- 匹配药品下拉搜索 -->
          <template v-else-if="column.key === 'matched_drug'">
            <a-select
              v-model:value="record.matched_drug_id"
              show-search
              allow-clear
              :filter-option="false"
              :loading="record.searching"
              :options="record.drug_options"
              placeholder="输入药品名搜索..."
              size="small"
              style="width: 100%"
              @search="(v) => debouncedSearchDrug(record, v)"
              @change="(v) => onDrugChange(record, v)"
            />
          </template>

          <!-- 批号 -->
          <template v-else-if="column.key === 'batch_number'">
            <a-input v-model:value="record.batch_number" size="small" style="width: 100%" />
          </template>

          <!-- 有效期 -->
          <template v-else-if="column.key === 'expire_date'">
            <a-date-picker
              v-model:value="record.expire_date"
              value-format="YYYY-MM-DD"
              size="small"
              style="width: 118px"
            />
          </template>

          <!-- 数量 -->
          <template v-else-if="column.key === 'planned_qty'">
            <a-input-number
              v-model:value="record.planned_qty"
              :min="1"
              size="small"
              style="width: 68px"
            />
          </template>

          <!-- 单价 -->
          <template v-else-if="column.key === 'unit_price'">
            <a-input-number
              v-model:value="record.unit_price"
              :min="0"
              :precision="2"
              size="small"
              style="width: 82px"
            />
          </template>
        </template>
      </a-table>

      <!-- 弹窗底部操作 -->
      <div style="display: flex; justify-content: space-between; align-items: center; margin-top: 14px">
        <span style="font-size: 12px; color: #888">
          已勾选 {{ includedCount }} / {{ editItems.length }} 行
        </span>
        <a-space>
          <a-button @click="convertModalOpen = false">取消</a-button>
          <a-button
            type="primary"
            :loading="convertLoading"
            :disabled="!canConvert"
            @click="doConvert"
          >
            创建入库单（{{ includedCount }} 个药品）
          </a-button>
        </a-space>
      </div>
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import {
  getInvoiceRecord,
  convertInvoiceToInbound,
  searchDrugsWithAI,
} from '@/api/ai'
import { getSupplierList, getSupplierDetail } from '@/api/suppliers'
import type { InvoiceRecord } from '@/api/ai'

const route = useRoute()
const router = useRouter()
const invoiceId = Number(route.params.id)

const loading = ref(false)
const activeTab = ref('info')
const invoice = ref<InvoiceRecord | null>(null)

const drugColumns = [
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 130 },
  { title: '批号', dataIndex: 'batch_number', width: 130 },
  { title: '数量', dataIndex: 'quantity', width: 70 },
  { title: '单价', dataIndex: 'unit_price', width: 90 },
  { title: '匹配状态', key: 'is_matched', width: 90 },
  { title: 'AI置信度', key: 'confidence', width: 140 },
]

// ===== 转换弹窗状态 =====
interface EditItem {
  row_index: number
  ai_name: string
  matched_drug_id: number | null
  batch_number: string
  expire_date: string
  planned_qty: number
  unit_price: number
  include: boolean
  drug_options: Array<{ value: number; label: string }>
  searching: boolean
}

const convertModalOpen = ref(false)
const convertLoading = ref(false)
const editItems = ref<EditItem[]>([])
const editSupplierId = ref<number | null>(null)
const editSupplierOptions = ref<Array<{ value: number; label: string }>>([])
const editInvoiceNo = ref('')
const supplierSearching = ref(false)

const includedCount = computed(() =>
  editItems.value.filter(i => i.include && i.matched_drug_id != null).length,
)
const unmatchedCount = computed(() =>
  editItems.value.filter(i => !i.matched_drug_id).length,
)
const canConvert = computed(
  () => editSupplierId.value != null && includedCount.value > 0,
)

const editColumns = [
  { title: '包含', key: 'include', width: 55 },
  { title: 'AI识别名称', key: 'ai_name', width: 120, ellipsis: true },
  { title: '匹配药品', key: 'matched_drug', width: 215 },
  { title: '批号', key: 'batch_number', width: 110 },
  { title: '有效期', key: 'expire_date', width: 130 },
  { title: '数量', key: 'planned_qty', width: 80 },
  { title: '单价(¥)', key: 'unit_price', width: 100 },
]

async function openConvertModal() {
  const result = invoice.value?.result
  if (!result) return

  editInvoiceNo.value = result.invoice_no ?? ''
  editSupplierId.value = result.matched_supplier_id ?? null

  editSupplierOptions.value = (result.supplier_candidates ?? []).map(s => ({
    value: s.supplier_id,
    label: `${s.name}（匹配度 ${Math.round(s.confidence * 100)}%）`,
  }))

  // 若已匹配供应商不在候选列表中，单独拉取名称
  if (
    editSupplierId.value &&
    !editSupplierOptions.value.find(o => o.value === editSupplierId.value)
  ) {
    try {
      const s = await getSupplierDetail(editSupplierId.value)
      editSupplierOptions.value.unshift({ value: s.id, label: s.name })
    } catch {
      editSupplierOptions.value.unshift({
        value: editSupplierId.value,
        label: `供应商 #${editSupplierId.value}`,
      })
    }
  }

  editItems.value = (result.items ?? []).map(item => {
    const candidates = item.drug_candidates ?? []
    const drug_options = candidates.map(c => ({
      value: c.drug_id,
      label: [c.common_name, c.specification, c.manufacturer].filter(Boolean).join(' · '),
    }))

    return {
      row_index: item.row_index,
      ai_name: item.drug_name ?? '',
      matched_drug_id: item.matched_drug_id ?? null,
      batch_number: item.batch_number ?? '',
      expire_date: item.expire_date ?? '',
      planned_qty: Math.max(1, parseInt(item.quantity ?? '1', 10) || 1),
      unit_price: parseFloat(item.unit_price ?? '0') || 0,
      include: !!item.matched_drug_id,
      drug_options,
      searching: false,
    }
  })

  convertModalOpen.value = true
}

// 药品搜索（各行独立防抖）
const _drugTimers: Record<number, ReturnType<typeof setTimeout>> = {}

function debouncedSearchDrug(item: EditItem, keyword: string) {
  if (!keyword.trim()) return
  clearTimeout(_drugTimers[item.row_index])
  _drugTimers[item.row_index] = setTimeout(async () => {
    item.searching = true
    try {
      const results = await searchDrugsWithAI({ query: keyword, search_mode: 'HYBRID', limit: 15 })
      item.drug_options = results.map(r => ({
        value: r.drug_id,
        label: [r.common_name, r.specification, r.manufacturer].filter(Boolean).join(' · '),
      }))
    } catch {
      // ignore
    } finally {
      item.searching = false
    }
  }, 350)
}

function onDrugChange(item: EditItem, drugId: number | null) {
  item.matched_drug_id = drugId
  item.include = drugId != null
}

// 供应商搜索（防抖）
let _supplierTimer: ReturnType<typeof setTimeout> | null = null

function debouncedSearchSupplier(keyword: string) {
  if (_supplierTimer) clearTimeout(_supplierTimer)
  _supplierTimer = setTimeout(async () => {
    supplierSearching.value = true
    try {
      const res = await getSupplierList({ name: keyword, page: 1, page_size: 20 })
      editSupplierOptions.value = res.list.map(s => ({ value: s.id, label: s.name }))
    } catch {
      // ignore
    } finally {
      supplierSearching.value = false
    }
  }, 300)
}

async function doConvert() {
  const validItems = editItems.value.filter(i => i.include && i.matched_drug_id != null)
  if (!editSupplierId.value) { message.warning('请选择供应商'); return }
  if (validItems.length === 0) { message.warning('请至少勾选一个已匹配的药品'); return }

  convertLoading.value = true
  try {
    const res = await convertInvoiceToInbound(invoiceId, {
      supplier_id: editSupplierId.value,
      invoice_no: editInvoiceNo.value || undefined,
      items: validItems.map(i => ({
        drug_id: i.matched_drug_id!,
        batch_number: i.batch_number,
        expire_date: i.expire_date,
        planned_qty: i.planned_qty,
        unit_price: i.unit_price,
      })),
    })
    message.success(`入库单 ${res['order_no']} 已创建`)
    router.push(`/inbound/orders/${res['id']}`)
  } finally {
    convertLoading.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    invoice.value = await getInvoiceRecord(invoiceId)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.modal-form-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.modal-label {
  white-space: nowrap;
  font-size: 13px;
  color: #333;
  min-width: 52px;
}

.required {
  color: #ff4d4f;
}
</style>
