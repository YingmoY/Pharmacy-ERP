<template>
  <div class="page-content">
    <a-card :title="isNew ? '新建入库单' : '编辑入库单'" style="margin-bottom: 16px">
      <a-form :model="formData" layout="vertical" :rules="rules" ref="formRef">
        <a-row :gutter="24">
          <a-col :span="8">
            <a-form-item label="供应商" name="supplier_id" required>
              <a-select
                v-model:value="formData.supplier_id"
                placeholder="请选择供应商"
                show-search
                :filter-option="filterSupplier"
                @change="handleSupplierChange"
              >
                <a-select-option
                  v-for="s in suppliers"
                  :key="s.id"
                  :value="s.id"
                  :disabled="s.status === 0"
                >
                  {{ s.name }}（{{ s.supplier_code }}）
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="发票号" name="invoice_no">
              <a-input v-model:value="formData.invoice_no" placeholder="选填" />
            </a-form-item>
          </a-col>
          <a-col :span="8">
            <a-form-item label="备注" name="remark">
              <a-input v-model:value="formData.remark" placeholder="选填" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>
    </a-card>

    <!-- 入库明细 -->
    <a-card title="入库明细">
      <template #extra>
        <a-button type="dashed" @click="addDetailRow">
          <PlusOutlined /> 添加明细
        </a-button>
      </template>

      <a-table
        :columns="detailColumns"
        :data-source="detailRows"
        :pagination="false"
        size="small"
        row-key="key"
      >
        <template #bodyCell="{ column, record, index }">
          <template v-if="column.key === 'drug_id'">
            <a-select
              v-model:value="record.drug_id"
              placeholder="选择药品"
              show-search
              :filter-option="filterDrug"
              style="width: 100%"
              @change="() => calcAmount(index)"
            >
              <a-select-option v-for="d in drugs" :key="d.id" :value="d.id" :disabled="d.status === 0">
                {{ d.common_name }}（{{ d.specification }}）
              </a-select-option>
            </a-select>
          </template>
          <template v-if="column.key === 'batch_number'">
            <a-input v-model:value="record.batch_number" placeholder="批号" />
          </template>
          <template v-if="column.key === 'expire_date'">
            <a-date-picker v-model:value="record.expire_date" style="width: 100%" format="YYYY-MM-DD" />
          </template>
          <template v-if="column.key === 'planned_qty'">
            <a-input-number v-model:value="record.planned_qty" :min="1" style="width: 80px" @change="() => calcAmount(index)" />
          </template>
          <template v-if="column.key === 'unit_price'">
            <a-input-number v-model:value="record.unit_price" :min="0" :precision="2" prefix="¥" style="width: 100px" @change="() => calcAmount(index)" />
          </template>
          <template v-if="column.key === 'amount'">
            <span style="color: #e6550d">¥{{ record.amount }}</span>
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" danger size="small" @click="removeDetailRow(index)">删除</a-button>
          </template>
        </template>
      </a-table>

      <div style="text-align: right; margin-top: 12px; font-size: 16px; font-weight: 600">
        合计：¥{{ totalAmount }}
      </div>
    </a-card>

    <!-- 底部按钮 -->
    <div style="margin-top: 16px; display: flex; gap: 12px; justify-content: flex-end">
      <a-button @click="router.back()">取消</a-button>
      <a-button :loading="saveLoading" @click="handleSave">保存草稿</a-button>
      <a-button type="primary" :loading="submitLoading" @click="handleSubmit">保存并提交</a-button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { message, Modal } from 'ant-design-vue'
import dayjs, { type Dayjs } from 'dayjs'
import { PlusOutlined } from '@ant-design/icons-vue'
import { getSupplierList } from '@/api/suppliers'
import { getDrugList } from '@/api/drugs'
import {
  createInboundOrder,
  updateInboundOrder,
  addInboundDetail,
  getInboundOrderDetail,
  getInboundDetailList,
  submitInboundOrder,
} from '@/api/inbound'
import type { Supplier } from '@/types/supplier'
import type { DrugInfo } from '@/types/drug'

const route = useRoute()
const router = useRouter()
const isNew = route.name === 'InboundOrderNew'
const orderId = isNew ? null : Number(route.params.id)

const formRef = ref()
const saveLoading = ref(false)
const submitLoading = ref(false)
const suppliers = ref<Supplier[]>([])
const drugs = ref<DrugInfo[]>([])
const currentOrderId = ref<number | null>(orderId)

const formData = reactive({
  supplier_id: undefined as number | undefined,
  invoice_no: '',
  remark: '',
})

interface DetailRow {
  key: string
  drug_id: number | undefined
  batch_number: string
  expire_date: Dayjs | null
  planned_qty: number
  unit_price: number
  amount: string
}

const detailRows = ref<DetailRow[]>([])

const rules = {
  supplier_id: [{ required: true, message: '请选择供应商' }],
}

const detailColumns = [
  { title: '药品', key: 'drug_id', width: 240 },
  { title: '批号', key: 'batch_number', width: 120 },
  { title: '有效期', key: 'expire_date', width: 140 },
  { title: '计划数量', key: 'planned_qty', width: 100 },
  { title: '单价', key: 'unit_price', width: 130 },
  { title: '金额', key: 'amount', width: 100 },
  { title: '操作', key: 'action', width: 60 },
]

const totalAmount = computed(() => {
  return detailRows.value.reduce((sum, row) => sum + parseFloat(row.amount || '0'), 0).toFixed(2)
})

function addDetailRow() {
  detailRows.value.push({
    key: Date.now().toString(),
    drug_id: undefined,
    batch_number: '',
    expire_date: null,
    planned_qty: 1,
    unit_price: 0,
    amount: '0.00',
  })
}

function removeDetailRow(index: number) {
  detailRows.value.splice(index, 1)
}

function calcAmount(index: number) {
  const row = detailRows.value[index]
  row.amount = (row.planned_qty * row.unit_price).toFixed(2)
}

function filterSupplier(input: string, option: { children: string }) {
  return option.children?.toLowerCase().includes(input.toLowerCase())
}

function filterDrug(input: string, option: { children: string }) {
  return option.children?.toLowerCase().includes(input.toLowerCase())
}

function handleSupplierChange() {
  // 供应商切换时的处理
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (!currentOrderId.value) {
      const order = await createInboundOrder({
        supplier_id: formData.supplier_id!,
        invoice_no: formData.invoice_no || undefined,
        remark: formData.remark || undefined,
      })
      currentOrderId.value = order.id
      // 添加明细
      for (const row of detailRows.value) {
        if (!row.drug_id || !row.batch_number || !row.expire_date) continue
        await addInboundDetail(order.id, {
          drug_id: row.drug_id,
          batch_number: row.batch_number,
          expire_date: row.expire_date.format('YYYY-MM-DD'),
          planned_qty: row.planned_qty,
          unit_price: row.unit_price,
        })
      }
    } else {
      await updateInboundOrder(currentOrderId.value, {
        supplier_id: formData.supplier_id,
        invoice_no: formData.invoice_no || undefined,
        remark: formData.remark || undefined,
      })
    }
    message.success('保存成功')
    router.push(`/inbound/orders/${currentOrderId.value}`)
  } finally {
    saveLoading.value = false
  }
}

async function handleSubmit() {
  await formRef.value?.validate()
  if (detailRows.value.length === 0) {
    message.warning('请至少添加一条入库明细')
    return
  }
  Modal.confirm({
    title: '确认提交',
    content: '提交后明细不可再修改，确认提交吗？',
    okText: '确认',
    onOk: async () => {
      submitLoading.value = true
      try {
        await handleSave()
        if (currentOrderId.value) {
          await submitInboundOrder(currentOrderId.value)
          message.success('已提交，等待扫码确认')
          router.push(`/inbound/orders/${currentOrderId.value}`)
        }
      } finally {
        submitLoading.value = false
      }
    },
  })
}

onMounted(async () => {
  // 加载供应商和药品列表
  const [suppRes, drugRes] = await Promise.allSettled([
    getSupplierList({ page: 1, page_size: 1000, status: 1 }),
    getDrugList({ page: 1, page_size: 1000, status: 1 }),
  ])
  if (suppRes.status === 'fulfilled') suppliers.value = suppRes.value.list
  if (drugRes.status === 'fulfilled') drugs.value = drugRes.value.list

  // 编辑模式：加载现有数据
  if (!isNew && orderId) {
    const [orderRes, detailRes] = await Promise.allSettled([
      getInboundOrderDetail(orderId),
      getInboundDetailList(orderId),
    ])
    if (orderRes.status === 'fulfilled') {
      const order = orderRes.value
      formData.supplier_id = order.supplier_id
      formData.invoice_no = order.invoice_no ?? ''
      formData.remark = order.remark ?? ''
    }
    if (detailRes.status === 'fulfilled') {
      detailRows.value = detailRes.value.map((d) => ({
        key: d.id.toString(),
        drug_id: d.drug_id,
        batch_number: d.batch_number,
        expire_date: dayjs(d.expire_date),
        planned_qty: d.planned_qty,
        unit_price: parseFloat(d.unit_price),
        amount: d.amount,
      }))
    }
  } else {
    addDetailRow()
  }
})
</script>
