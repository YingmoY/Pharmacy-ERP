<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="药品编码">
          <a-input v-model:value="searchParams.drug_code" allow-clear />
        </a-form-item>
        <a-form-item label="通用名">
          <a-input v-model:value="searchParams.common_name" allow-clear />
        </a-form-item>
        <a-form-item label="厂家">
          <a-input v-model:value="searchParams.manufacturer" allow-clear />
        </a-form-item>
        <a-form-item label="类型">
          <a-select v-model:value="searchParams.is_prescription" placeholder="全部" style="width: 100px" allow-clear>
            <a-select-option :value="true">处方药</a-select-option>
            <a-select-option :value="false">非处方药</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 90px" allow-clear>
            <a-select-option :value="1">启用</a-select-option>
            <a-select-option :value="0">停用</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <div style="margin-bottom: 12px">
      <a-button type="primary" @click="showDrawer()">
        <PlusOutlined /> 新增药品
      </a-button>
    </div>

    <a-card>
      <a-table
        :columns="columns"
        :data-source="list"
        :loading="loading"
        :pagination="pagination"
        row-key="id"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'name'">
            <div style="font-weight: 500">{{ record.common_name }}</div>
            <div style="font-size: 12px; color: #666">{{ record.trade_name }}</div>
          </template>
          <template v-if="column.key === 'is_prescription'">
            <a-tag :color="record.is_prescription ? 'red' : 'blue'">
              {{ record.is_prescription ? '处方药' : '非处方药' }}
            </a-tag>
          </template>
          <template v-if="column.key === 'is_medicare'">
            <a-tag v-if="record.is_medicare" color="green">医保</a-tag>
            <span v-else style="color: #999">-</span>
          </template>
          <template v-if="column.key === 'retail_price'">
            <span style="color: #e6550d">¥{{ record.retail_price }}</span>
          </template>
          <template v-if="column.key === 'status'">
            <StatusTag type="enable" :value="String(record.status)" />
          </template>
          <template v-if="column.key === 'in_stock_count'">
            <span :style="{ color: (record.in_stock_count ?? 0) > 0 ? '#52c41a' : '#ff4d4f', fontWeight: 500 }">
              {{ record.in_stock_count ?? 0 }}
            </span>
          </template>
          <template v-if="column.key === 'action'">
            <div class="table-actions">
              <a-button type="link" size="small" @click="showInventory(record)">查看库存</a-button>
              <a-button type="link" size="small" @click="showDrawer(record)">编辑</a-button>
              <a-popconfirm
                :title="record.status === 1 ? '确认停用该药品吗？' : '确认启用该药品吗？'"
                @confirm="toggleStatus(record)"
              >
                <a-button type="link" size="small" :danger="record.status === 1">
                  {{ record.status === 1 ? '停用' : '启用' }}
                </a-button>
              </a-popconfirm>
            </div>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 库存明细弹窗 -->
    <a-modal
      v-model:open="inventoryVisible"
      :title="`库存明细 — ${inventoryDrug?.common_name ?? ''}`"
      width="800"
      :footer="null"
      @cancel="inventoryVisible = false"
    >
      <a-table
        :columns="inventoryColumns"
        :data-source="inventoryList"
        :loading="inventoryLoading"
        :pagination="inventoryPagination"
        row-key="id"
        size="small"
        @change="handleInventoryTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <a-tag :color="record.status === 'IN_STOCK' ? 'green' : record.status === 'RESERVED' ? 'orange' : 'default'">
              {{ statusLabel(record.status) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'expire_date'">
            {{ record.expire_date ? record.expire_date.slice(0, 10) : '-' }}
          </template>
        </template>
      </a-table>
    </a-modal>

    <!-- 新增/编辑抽屉 -->
    <a-drawer
      v-model:open="drawerVisible"
      :title="editingDrug ? '编辑药品' : '新增药品'"
      width="560"
      @close="handleDrawerClose"
    >
      <a-form :model="form" layout="vertical" ref="formRef" :rules="formRules">
        <a-row :gutter="16">
          <a-col :span="12">
            <a-form-item label="药品编码" name="drug_code">
              <a-input v-model:value="form.drug_code" :disabled="!!editingDrug" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="通用名" name="common_name">
              <a-input v-model:value="form.common_name" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="商品名" name="trade_name">
              <a-input v-model:value="form.trade_name" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="规格" name="specification">
              <a-input v-model:value="form.specification" placeholder="如：10mg×20片" />
            </a-form-item>
          </a-col>
          <a-col :span="24">
            <a-form-item label="厂家" name="manufacturer">
              <a-input v-model:value="form.manufacturer" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="剂型" name="dosage_form">
              <a-input v-model:value="form.dosage_form" placeholder="片剂/胶囊/注射液等" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="单位" name="unit">
              <a-input v-model:value="form.unit" placeholder="盒/支/瓶等" />
            </a-form-item>
          </a-col>
          <a-col :span="12">
            <a-form-item label="零售价" name="retail_price">
              <a-input-number v-model:value="form.retail_price" :min="0" :precision="2" prefix="¥" style="width: 100%" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="处方药">
              <a-switch v-model:checked="form.is_prescription" />
            </a-form-item>
          </a-col>
          <a-col :span="6">
            <a-form-item label="医保">
              <a-switch v-model:checked="form.is_medicare" />
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>

      <template #footer>
        <a-space>
          <a-button @click="handleDrawerClose">取消</a-button>
          <a-button type="primary" :loading="saveLoading" @click="handleSave">保存</a-button>
        </a-space>
      </template>
    </a-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getDrugList, createDrug, updateDrug, toggleDrugStatus, getDrugInventoryList } from '@/api/drugs'
import type { DrugInfo } from '@/types/drug'

const loading = ref(false)
const saveLoading = ref(false)
const list = ref<DrugInfo[]>([])
const drawerVisible = ref(false)
const editingDrug = ref<DrugInfo | null>(null)
const formRef = ref()

const searchParams = reactive({
  drug_code: '',
  common_name: '',
  manufacturer: '',
  is_prescription: undefined as boolean | undefined,
  status: undefined as number | undefined,
})

const form = reactive({
  drug_code: '',
  common_name: '',
  trade_name: '',
  specification: '',
  manufacturer: '',
  dosage_form: '',
  unit: '',
  retail_price: 0,
  is_prescription: false,
  is_medicare: false,
})

const formRules = {
  drug_code: [{ required: true, message: '请输入药品编码' }],
  common_name: [{ required: true, message: '请输入通用名' }],
  specification: [{ required: true, message: '请输入规格' }],
  manufacturer: [{ required: true, message: '请输入厂家' }],
  retail_price: [{ required: true, message: '请输入零售价' }],
}

const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

const columns = [
  { title: '药品编码', dataIndex: 'drug_code', width: 120 },
  { title: '药品名称', key: 'name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 140 },
  { title: '厂家', dataIndex: 'manufacturer', ellipsis: true },
  { title: '类型', key: 'is_prescription', width: 90 },
  { title: '医保', key: 'is_medicare', width: 60 },
  { title: '零售价', key: 'retail_price', width: 90 },
  { title: '在库数量', key: 'in_stock_count', width: 90 },
  { title: '状态', key: 'status', width: 70 },
  { title: '操作', key: 'action', width: 160 },
]

const inventoryColumns = [
  { title: '追溯码', dataIndex: 'trace_code', ellipsis: true },
  { title: '批号', dataIndex: 'batch_number', width: 140 },
  { title: '有效期', key: 'expire_date', width: 110 },
  { title: '货位', dataIndex: 'location_code', width: 120 },
  { title: '状态', key: 'status', width: 90 },
]

const inventoryVisible = ref(false)
const inventoryLoading = ref(false)
const inventoryDrug = ref<DrugInfo | null>(null)
const inventoryList = ref<any[]>([])
const inventoryPagination = reactive({ current: 1, pageSize: 20, total: 0 })

function statusLabel(s: string) {
  const map: Record<string, string> = {
    IN_STOCK: '在库',
    RESERVED: '已预留',
    SOLD: '已售出',
    RETURNED: '已退货',
    DAMAGED: '损耗',
    EXPIRED: '已过期',
  }
  return map[s] ?? s
}

async function loadInventory() {
  if (!inventoryDrug.value) return
  inventoryLoading.value = true
  try {
    const res = await getDrugInventoryList(inventoryDrug.value.id, {
      page: inventoryPagination.current,
      page_size: inventoryPagination.pageSize,
    })
    inventoryList.value = res.list
    inventoryPagination.total = res.total
  } finally {
    inventoryLoading.value = false
  }
}

function showInventory(drug: DrugInfo) {
  inventoryDrug.value = drug
  inventoryPagination.current = 1
  inventoryList.value = []
  inventoryVisible.value = true
  loadInventory()
}

function handleInventoryTableChange(pag: typeof inventoryPagination) {
  inventoryPagination.current = pag.current
  inventoryPagination.pageSize = pag.pageSize
  loadInventory()
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getDrugList({
      page: pagination.current,
      page_size: pagination.pageSize,
      drug_code: searchParams.drug_code || undefined,
      common_name: searchParams.common_name || undefined,
      manufacturer: searchParams.manufacturer || undefined,
      is_prescription: searchParams.is_prescription,
      status: searchParams.status,
    })
    list.value = res.list
    pagination.total = res.total
  } finally {
    loading.value = false
  }
}

function handleReset() {
  Object.assign(searchParams, { drug_code: '', common_name: '', manufacturer: '', is_prescription: undefined, status: undefined })
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  pagination.pageSize = pag.pageSize
  fetchList()
}

function showDrawer(drug?: DrugInfo) {
  editingDrug.value = drug ?? null
  if (drug) {
    Object.assign(form, drug, { retail_price: parseFloat(drug.retail_price) })
  } else {
    Object.assign(form, { drug_code: '', common_name: '', trade_name: '', specification: '', manufacturer: '', dosage_form: '', unit: '', retail_price: 0, is_prescription: false, is_medicare: false })
  }
  drawerVisible.value = true
}

function handleDrawerClose() {
  drawerVisible.value = false
  editingDrug.value = null
}

async function handleSave() {
  await formRef.value?.validate()
  saveLoading.value = true
  try {
    if (editingDrug.value) {
      await updateDrug(editingDrug.value.id, { ...form, retail_price: form.retail_price.toString() })
    } else {
      await createDrug({ ...form, retail_price: form.retail_price.toString() })
    }
    message.success('保存成功')
    handleDrawerClose()
    fetchList()
  } finally {
    saveLoading.value = false
  }
}

async function toggleStatus(drug: DrugInfo) {
  await toggleDrugStatus(drug.id, drug.status === 1 ? 0 : 1)
  message.success('状态已更新')
  fetchList()
}

onMounted(fetchList)
</script>
