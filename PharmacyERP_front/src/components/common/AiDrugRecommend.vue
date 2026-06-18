<template>
  <div class="ai-recommend-box">
    <a-input-group compact>
      <a-textarea
        v-model:value="query"
        :placeholder="placeholder"
        :auto-size="{ minRows: 2, maxRows: 3 }"
        :disabled="loading"
        style="width: calc(100% - 72px); font-size: 13px"
        @pressEnter.prevent="handleSearch"
      />
      <a-button
        type="primary"
        style="width: 72px; height: auto"
        :loading="loading"
        @click="handleSearch"
      >
        {{ loading ? '' : 'AI推荐' }}
      </a-button>
    </a-input-group>

    <!-- Loading state -->
    <div v-if="loading" class="loading-bar">
      <a-spin size="small" />
      <span class="loading-text">{{ loadingHint }}</span>
    </div>

    <!-- AI explanation badge -->
    <div v-if="!loading && explanation" class="explanation-bar">
      <RobotOutlined style="margin-right: 4px; color: #722ed1" />
      <span>{{ explanation }}</span>
    </div>

    <!-- Results -->
    <div v-if="!loading && results.length > 0" class="search-results">
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
          <a-tag v-if="drug.match_reason" color="geekblue" size="small" style="font-size: 10px">{{ drug.match_reason }}</a-tag>
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

    <!-- No result -->
    <div v-if="!loading && searched && results.length === 0" class="no-result">
      <a-empty description="未找到推荐药品，请尝试其他描述" :image-style="{ height: '40px' }" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onUnmounted } from 'vue'
import { RobotOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import type { DrugSearchResult } from '@/types/drug'

const props = withDefaults(
  defineProps<{
    placeholder?: string
    onlyAvailable?: boolean
  }>(),
  {
    placeholder: '描述症状或药品名（可有错别字），例如：头痛发烧、消炎止痛、阿莫灵',
    onlyAvailable: true,
  },
)

const emit = defineEmits<{
  select: [drug: DrugSearchResult]
}>()

const query = ref('')
const results = ref<DrugSearchResult[]>([])
const explanation = ref('')
const loading = ref(false)
const searched = ref(false)
const loadingHint = ref('AI正在分析您的描述...')

let activeWs: WebSocket | null = null

function buildWsUrl(): string {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  return `${proto}//${window.location.host}/ai/api/v1/drugs/recommend/ws`
}

function closeWs() {
  if (activeWs) {
    activeWs.close()
    activeWs = null
  }
}

onUnmounted(closeWs)

function handleSearch() {
  const q = query.value.trim()
  if (!q || loading.value) return

  closeWs()
  loading.value = true
  searched.value = true
  explanation.value = ''
  results.value = []
  loadingHint.value = 'AI正在分析您的描述...'

  const ws = new WebSocket(buildWsUrl())
  activeWs = ws

  ws.onopen = () => {
    ws.send(JSON.stringify({
      query: q,
      limit: 10,
      filters: { only_available: props.onlyAvailable },
    }))
  }

  ws.onmessage = (event: MessageEvent) => {
    try {
      const msg = JSON.parse(event.data as string) as {
        type: 'progress' | 'result' | 'error'
        message?: string
        data?: { items: DrugSearchResult[]; explanation: string }
      }
      if (msg.type === 'progress') {
        loadingHint.value = msg.message ?? loadingHint.value
      } else if (msg.type === 'result') {
        results.value = msg.data?.items ?? []
        explanation.value = msg.data?.explanation ?? ''
        loading.value = false
        activeWs = null
      } else if (msg.type === 'error') {
        message.error(`AI推荐失败：${msg.message ?? '未知错误'}`)
        loading.value = false
        activeWs = null
      }
    } catch {
      // ignore malformed frames
    }
  }

  ws.onerror = () => {
    message.error('AI推荐服务连接失败，请稍后重试')
    loading.value = false
    activeWs = null
  }

  ws.onclose = (evt: CloseEvent) => {
    // If we're still in loading state when the socket closes unexpectedly
    if (loading.value) {
      if (!evt.wasClean) {
        message.error('AI推荐连接中断，请重试')
      }
      loading.value = false
      activeWs = null
    }
  }
}

function handleSelect(drug: DrugSearchResult) {
  emit('select', drug)
  query.value = ''
  results.value = []
  explanation.value = ''
  searched.value = false
}
</script>

<style scoped>
.ai-recommend-box {
  position: relative;
}

.loading-bar {
  margin-top: 10px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  background: #f5f0ff;
  border-radius: 6px;
  border: 1px solid #d9c8ff;
}

.loading-text {
  font-size: 13px;
  color: #722ed1;
  animation: fade-cycle 0.4s ease;
}

@keyframes fade-cycle {
  from { opacity: 0.4; }
  to   { opacity: 1; }
}

.explanation-bar {
  margin-top: 6px;
  padding: 6px 10px;
  background: #f9f0ff;
  border-radius: 4px;
  font-size: 12px;
  color: #531dab;
  line-height: 1.4;
}

.term-tags {
  margin-top: 6px;
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.search-results {
  margin-top: 6px;
  background: #fff;
  border: 1px solid #d9d9d9;
  border-radius: 6px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  max-height: 380px;
  overflow-y: auto;
}

.search-result-item {
  padding: 10px 12px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  transition: background 0.15s;
}

.search-result-item:last-child {
  border-bottom: none;
}

.search-result-item:hover {
  background: #f5f0ff;
}

.drug-name {
  font-size: 13px;
  font-weight: 500;
  margin-bottom: 4px;
  display: flex;
  align-items: center;
  gap: 4px;
  flex-wrap: wrap;
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
