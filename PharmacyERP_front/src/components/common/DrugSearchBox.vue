<template>
  <div class="drug-search-box">
    <a-input-group compact>
      <a-input
        v-model:value="keyword"
        :placeholder="placeholder"
        style="width: calc(100% - 80px)"
        size="large"
        allow-clear
        @input="handleInput"
        @pressEnter="handleSearch"
      >
        <template #prefix><SearchOutlined /></template>
      </a-input>
      <a-button type="primary" size="large" style="width: 80px" :loading="loading" @click="handleSearch">
        搜索
      </a-button>
    </a-input-group>

    <!-- 搜索结果列表 -->
    <div v-if="results.length > 0" class="search-results">
      <div v-if="usingFallback" class="fallback-tip">
        <a-tag color="orange" style="margin: 4px 8px; font-size: 11px">AI服务不可用，已切换为普通搜索</a-tag>
      </div>
      <div
        v-for="drug in results"
        :key="drug.drug_id"
        class="search-result-item"
        @click="handleSelect(drug)"
      >
        <div class="drug-name">
          {{ drug.common_name }}
          <a-tag v-if="drug.is_prescription" color="red" size="small">处方药</a-tag>
          <a-tag v-if="(drug.inventory?.near_expire_available_qty ?? 0) > 0" color="orange" size="small">近效期</a-tag>
        </div>
        <div class="drug-info">
          <span>{{ drug.specification }}</span>
          <span class="divider">·</span>
          <span>{{ drug.manufacturer }}</span>
          <span class="divider">·</span>
          <span class="price">¥{{ drug.retail_price }}</span>
          <span class="divider">·</span>
          <span :class="(drug.inventory?.available_qty ?? 0) > 0 ? 'stock-ok' : 'stock-empty'">
            可售 {{ drug.inventory?.available_qty ?? 0 }} 盒
          </span>
        </div>
      </div>
    </div>

    <!-- 无结果提示 -->
    <div v-if="searched && results.length === 0 && !loading" class="no-result">
      <a-empty description="未找到匹配药品" :image-style="{ height: '40px' }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { SearchOutlined } from '@ant-design/icons-vue'
import { searchDrugsWithAI } from '@/api/ai'
import { getDrugList } from '@/api/drugs'
import type { DrugSearchResult } from '@/types/drug'

const props = withDefaults(
  defineProps<{
    placeholder?: string
    onlyAvailable?: boolean
  }>(),
  {
    placeholder: '输入药品名称、拼音或症状词智能搜索',
    onlyAvailable: true,
  },
)

const emit = defineEmits<{
  select: [drug: DrugSearchResult]
}>()

const keyword = ref('')
const results = ref<DrugSearchResult[]>([])
const loading = ref(false)
const searched = ref(false)
const usingFallback = ref(false)

let debounceTimer: ReturnType<typeof setTimeout> | null = null

function handleInput() {
  if (debounceTimer) clearTimeout(debounceTimer)
  if (!keyword.value.trim()) {
    results.value = []
    searched.value = false
    usingFallback.value = false
    return
  }
  debounceTimer = setTimeout(() => {
    handleSearch()
  }, 400)
}

async function handleSearch() {
  if (!keyword.value.trim()) return
  loading.value = true
  searched.value = true
  usingFallback.value = false
  try {
    const aiResults = await searchDrugsWithAI({
      query: keyword.value,
      search_mode: 'HYBRID',
      limit: 10,
      filters: { only_available: props.onlyAvailable },
    })
    results.value = aiResults
  } catch {
    // AI 服务不可用，降级到普通药品搜索
    await fallbackSearch()
  } finally {
    loading.value = false
  }
}

async function fallbackSearch() {
  usingFallback.value = true
  try {
    const res = await getDrugList({
      keyword: keyword.value,
      status: 1,
      page_size: 10,
    })
    results.value = res.list.map((drug) => ({
      drug_id: drug.id,
      drug_code: drug.drug_code,
      common_name: drug.common_name,
      trade_name: drug.trade_name,
      specification: drug.specification,
      dosage_form: drug.dosage_form,
      manufacturer: drug.manufacturer,
      unit: drug.unit,
      retail_price: drug.retail_price,
      is_prescription: drug.is_prescription,
      is_medicare: drug.is_medicare,
      score: 0,
      inventory: {
        available_qty: (drug as any).in_stock_count ?? 0,
        near_expire_available_qty: 0,
      },
    } satisfies DrugSearchResult))
  } catch {
    results.value = []
  }
}

function handleSelect(drug: DrugSearchResult) {
  emit('select', drug)
  keyword.value = ''
  results.value = []
  searched.value = false
  usingFallback.value = false
}
</script>

<style scoped>
.drug-search-box {
  position: relative;
}

.search-results {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 0 0 6px 6px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.12);
  z-index: 1000;
  max-height: 400px;
  overflow-y: auto;
}

.fallback-tip {
  border-bottom: 1px solid #f0f0f0;
  background: #fffbe6;
}

.search-result-item {
  padding: 10px 16px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background 0.15s;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: #f0f7ff;
}

.drug-name {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 6px;
}

.drug-info {
  font-size: 12px;
  color: #666;
  display: flex;
  align-items: center;
  gap: 4px;
}

.divider {
  color: #ccc;
}

.price {
  color: #e6550d;
  font-weight: 500;
}

.stock-ok {
  color: #52c41a;
}

.stock-empty {
  color: #ff4d4f;
}

.no-result {
  padding: 16px;
  text-align: center;
}
</style>
