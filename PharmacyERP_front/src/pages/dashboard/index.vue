<template>
  <div class="page-content">
    <!-- 顶部统计卡片 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <a-col :span="5" v-for="stat in statsCards" :key="stat.key">
        <a-card hoverable @click="stat.onClick">
          <a-statistic
            :title="stat.title"
            :value="stat.value"
            :prefix="stat.prefix"
            :value-style="{ color: stat.color, fontSize: '28px' }"
          />
          <div style="margin-top: 8px; font-size: 12px; color: #999">{{ stat.desc }}</div>
        </a-card>
      </a-col>
    </a-row>

    <!-- 中间区域 -->
    <a-row :gutter="16" style="margin-bottom: 16px">
      <!-- 销售趋势图 -->
      <a-col :span="16">
        <a-card title="近7日销售趋势" :loading="trendLoading">
          <div style="height: 280px">
            <v-chart :option="chartOption" autoresize />
          </div>
        </a-card>
      </a-col>

      <!-- 待办事项 -->
      <a-col :span="8">
        <a-card title="待办事项">
          <a-list :data-source="todoItems" size="small">
            <template #renderItem="{ item }">
              <a-list-item
                style="cursor: pointer; padding: 10px 0"
                @click="router.push(item.route)"
              >
                <a-list-item-meta>
                  <template #avatar>
                    <a-badge :count="item.count" :overflow-count="999">
                      <div
                        class="todo-icon"
                        :style="{ background: item.bg }"
                      >
                        <component :is="item.icon" style="color: #fff; font-size: 16px" />
                      </div>
                    </a-badge>
                  </template>
                  <template #title>{{ item.title }}</template>
                  <template #description>{{ item.desc }}</template>
                </a-list-item-meta>
                <RightOutlined style="color: #999" />
              </a-list-item>
            </template>
          </a-list>
        </a-card>
      </a-col>
    </a-row>

    <!-- 底部区域 -->
    <a-row :gutter="16">
      <!-- 热销药品排行 -->
      <a-col :span="12">
        <a-card title="热销药品排行（今日）" :loading="topDrugsLoading">
          <a-table
            :columns="topDrugsColumns"
            :data-source="topDrugs"
            :pagination="false"
            size="small"
            row-key="drug_id"
          />
        </a-card>
      </a-col>

      <!-- 库存异常概览 -->
      <a-col :span="12">
        <a-card title="库存异常概览">
          <a-row :gutter="16">
            <a-col :span="12" v-for="item in abnormalItems" :key="item.key">
              <div
                class="abnormal-card"
                :style="{ borderColor: item.color }"
                @click="router.push(item.route)"
              >
                <div class="abnormal-number" :style="{ color: item.color }">
                  {{ item.count }}
                </div>
                <div class="abnormal-label">{{ item.label }}</div>
              </div>
            </a-col>
          </a-row>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import VChart from 'vue-echarts'
import {
  AppstoreOutlined,
  BellOutlined,
  FileSearchOutlined,
  FieldTimeOutlined,
  RightOutlined,
} from '@ant-design/icons-vue'
import {
  getDashboardOverview,
  getSalesTrend,
  getTopSellingDrugs,
} from '@/api/dashboard'
import type { DashboardOverview, TopDrugItem } from '@/types/report'
import dayjs from 'dayjs'

// 注册 ECharts 组件
use([CanvasRenderer, LineChart, GridComponent, TooltipComponent, LegendComponent])

const router = useRouter()

const stats = ref<DashboardOverview | null>(null)
const trendPoints = ref<{ date: string; sales_amount: number; sales_count: number }[]>([])
const topDrugs = ref<TopDrugItem[]>([])

const trendLoading = ref(true)
const topDrugsLoading = ref(true)

// 统计卡片
const statsCards = computed(() => [
  {
    key: 'sales',
    title: '今日销售额',
    value: stats.value?.today_sales_amount ?? 0,
    prefix: '¥',
    color: '#1677ff',
    desc: `${stats.value?.today_sales_count ?? 0} 笔销售`,
    onClick: () => router.push('/reports/sales'),
  },
  {
    key: 'inbound',
    title: '今日入库单数',
    value: stats.value?.today_inbound_count ?? 0,
    color: '#52c41a',
    desc: '今日完成入库',
    onClick: () => router.push('/inbound/orders'),
  },
  {
    key: 'stock',
    title: '当前在库数量',
    value: stats.value?.in_stock_count ?? 0,
    color: '#13c2c2',
    desc: '追溯码总量',
    onClick: () => router.push('/inventory'),
  },
  {
    key: 'alert',
    title: '活跃预警数',
    value: stats.value?.active_alert_count ?? 0,
    color: '#fa8c16',
    desc: '需要处理',
    onClick: () => router.push('/alerts'),
  },
  {
    key: 'nearExpire',
    title: '近效期数量',
    value: stats.value?.near_expire_count ?? 0,
    color: '#ff4d4f',
    desc: '30天内到期',
    onClick: () => router.push('/inventory/near-expire'),
  },
])

// 待办事项（来自看板概览数据）
const todoItems = computed(() => [
  {
    key: 'shelving',
    title: '待上架',
    desc: '入库完成等待上架',
    count: stats.value?.pending_shelving_count ?? 0,
    icon: AppstoreOutlined,
    bg: '#52c41a',
    route: '/inventory/pending-shelving',
  },
  {
    key: 'alert',
    title: '待处理预警',
    desc: '库存异常需要处理',
    count: stats.value?.active_alert_count ?? 0,
    icon: BellOutlined,
    bg: '#fa8c16',
    route: '/alerts',
  },
  {
    key: 'loss',
    title: '盘亏候选',
    desc: '盘库未扫到的追溯码',
    count: stats.value?.loss_candidate_count ?? 0,
    icon: FileSearchOutlined,
    bg: '#ff4d4f',
    route: '/inventory-tasks',
  },
  {
    key: 'nearExpire',
    title: '近效期药品',
    desc: '30天内到期库存',
    count: stats.value?.near_expire_count ?? 0,
    icon: FieldTimeOutlined,
    bg: '#722ed1',
    route: '/inventory/near-expire',
  },
])

// 热销排行表格列
const topDrugsColumns = [
  { title: '排名', dataIndex: 'rank', width: 60 },
  { title: '药品名称', dataIndex: 'common_name', ellipsis: true },
  { title: '规格', dataIndex: 'specification', width: 100, ellipsis: true },
  { title: '销售数量', dataIndex: 'total_quantity', width: 80 },
  {
    title: '销售金额',
    dataIndex: 'total_amount',
    width: 100,
    customRender: ({ value }: { value: number }) => `¥${Number(value).toFixed(2)}`,
  },
]

// 异常概览
const abnormalItems = computed(() => [
  {
    key: 'misplaced',
    label: '错架',
    count: 0,
    color: '#fa541c',
    route: '/inventory?status=MISPLACED',
  },
  {
    key: 'lossCandidate',
    label: '盘亏候选',
    count: stats.value?.loss_candidate_count ?? 0,
    color: '#ff4d4f',
    route: '/inventory?status=LOSS_CANDIDATE',
  },
  {
    key: 'nearExpire',
    label: '近效期',
    count: stats.value?.near_expire_count ?? 0,
    color: '#faad14',
    route: '/inventory/near-expire',
  },
  {
    key: 'pending',
    label: '待上架',
    count: stats.value?.pending_shelving_count ?? 0,
    color: '#1677ff',
    route: '/inventory/pending-shelving',
  },
])

// 图表配置
const chartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  grid: { left: 60, right: 20, top: 20, bottom: 30 },
  xAxis: {
    type: 'category',
    data: trendPoints.value.map((t) => t.date.slice(5)), // 只显示月-日
    axisLabel: { fontSize: 12 },
  },
  yAxis: { type: 'value', name: '金额（元）', nameTextStyle: { fontSize: 11 } },
  series: [
    {
      name: '销售额',
      type: 'line',
      data: trendPoints.value.map((t) => Number(t.sales_amount)),
      smooth: true,
      areaStyle: { opacity: 0.1 },
      itemStyle: { color: '#1677ff' },
      lineStyle: { color: '#1677ff', width: 2 },
    },
  ],
}))

onMounted(async () => {
  // 计算近7天日期范围
  const today = dayjs().format('YYYY-MM-DD')
  const sevenDaysAgo = dayjs().subtract(6, 'day').format('YYYY-MM-DD')

  // 并行加载数据
  const [statsRes] = await Promise.allSettled([getDashboardOverview()])
  if (statsRes.status === 'fulfilled') stats.value = statsRes.value

  // 加载趋势数据
  try {
    trendLoading.value = true
    const trendData = await getSalesTrend({ start_date: sevenDaysAgo, end_date: today, granularity: 'DAY' })
    trendPoints.value = trendData.points ?? []
  } catch {
    trendPoints.value = []
  } finally {
    trendLoading.value = false
  }

  // 加载热销排行
  try {
    topDrugsLoading.value = true
    topDrugs.value = await getTopSellingDrugs({ start_date: sevenDaysAgo, end_date: today, top_n: 5 })
  } catch {
    topDrugs.value = []
  } finally {
    topDrugsLoading.value = false
  }
})
</script>

<style scoped>
.todo-icon {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.abnormal-card {
  border: 2px solid;
  border-radius: 8px;
  padding: 16px;
  text-align: center;
  cursor: pointer;
  margin-bottom: 16px;
  transition: transform 0.15s;
}

.abnormal-card:hover {
  transform: translateY(-2px);
}

.abnormal-number {
  font-size: 32px;
  font-weight: 700;
  line-height: 1;
  margin-bottom: 6px;
}

.abnormal-label {
  font-size: 13px;
  color: #666;
}
</style>
