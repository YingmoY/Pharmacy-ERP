<template>
  <div class="page-content">
    <!-- 搜索框 -->
    <a-card style="margin-bottom: 16px">
      <a-input-search
        v-model:value="searchCode"
        :default-value="route.params.trace_code as string"
        placeholder="输入或扫描追溯码，回车查询"
        size="large"
        enter-button="查询"
        :loading="loading"
        @search="handleSearch"
      >
        <template #prefix><ScanOutlined /></template>
      </a-input-search>
    </a-card>

    <a-spin :spinning="loading">
      <template v-if="inventory">
        <!-- 当前状态卡片 -->
        <a-row :gutter="16" style="margin-bottom: 16px">
          <a-col :span="24">
            <a-card>
              <a-row :gutter="24">
                <a-col :span="4">
                  <a-statistic
                    title="当前状态"
                    :value="' '"
                  >
                    <template #formatter>
                      <StatusTag type="trace" :value="inventory.status" />
                      <a-tag v-if="inventory.is_reserved" color="orange" style="margin-left: 4px">已预占</a-tag>
                    </template>
                  </a-statistic>
                </a-col>
                <a-col :span="5">
                  <a-statistic title="追溯码" :value="inventory.trace_code" />
                </a-col>
                <a-col :span="5">
                  <a-statistic title="药品名称" :value="inventory.drug_name" />
                </a-col>
                <a-col :span="4">
                  <a-statistic title="批号" :value="inventory.batch_number" />
                </a-col>
                <a-col :span="3">
                  <a-statistic
                    title="有效期"
                    :value="inventory.expire_date"
                    :value-style="{ color: isNearExpire ? '#ff4d4f' : 'inherit' }"
                  />
                </a-col>
                <a-col :span="3">
                  <a-statistic title="当前货位" :value="inventory.location_code || '未上架'" />
                </a-col>
              </a-row>
            </a-card>
          </a-col>
        </a-row>

        <!-- 追溯时间线 -->
        <a-card title="追溯链路">
          <TraceTimeline :logs="traceLogs" />
        </a-card>
      </template>

      <a-empty v-else-if="searched && !loading" description="未找到该追溯码" />
      <a-empty v-else-if="!searched" description="请输入追溯码进行查询" />
    </a-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ScanOutlined } from '@ant-design/icons-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import TraceTimeline from '@/components/common/TraceTimeline.vue'
import { getTraceCodeDetail, getTraceLog } from '@/api/inventory'
import type { DrugTraceInventory, DrugTraceLog } from '@/types/inventory'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const searchCode = ref(route.params.trace_code as string || '')
const loading = ref(false)
const searched = ref(false)
const inventory = ref<DrugTraceInventory | null>(null)
const traceLogs = ref<DrugTraceLog[]>([])

const isNearExpire = computed(() => {
  if (!inventory.value?.expire_date) return false
  return dayjs(inventory.value.expire_date).diff(dayjs(), 'day') <= 30
})

async function handleSearch(code: string) {
  const traceCode = code.trim()
  if (!traceCode) return
  loading.value = true
  searched.value = true
  try {
    const [inv, logs] = await Promise.allSettled([
      getTraceCodeDetail(traceCode),
      getTraceLog(traceCode),
    ])
    inventory.value = inv.status === 'fulfilled' ? inv.value : null
    traceLogs.value = logs.status === 'fulfilled' ? logs.value.list : []
    // 更新 URL 但不刷新页面
    router.replace(`/trace/${traceCode}`)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  if (searchCode.value) {
    handleSearch(searchCode.value)
  }
})
</script>
