<template>
  <div class="page-content">
    <a-row :gutter="16">
      <a-col :span="10">
        <a-card title="混放检查">
          <a-form layout="vertical">
            <a-form-item label="货位编码">
              <ScanInput v-model="locationCode" placeholder="扫描或输入货位码" :auto-focus="true" :clear-after-submit="false" @submit="handleCheck" />
            </a-form-item>
            <a-button type="primary" block :loading="loading" @click="handleCheck(locationCode)" :disabled="!locationCode">
              执行检查
            </a-button>
          </a-form>
        </a-card>
      </a-col>

      <a-col :span="14">
        <a-card title="检查结果">
          <a-empty v-if="!result" description="请输入货位码并执行检查" />
          <div v-else>
            <a-alert
              :type="result.has_mixed_drugs ? 'warning' : 'success'"
              :message="result.has_mixed_drugs ? '该货位存在混放情况' : '该货位无混放情况'"
              style="margin-bottom: 16px"
              show-icon
            />

            <div v-if="(result.drugs ?? []).length > 0">
              <a-divider>货位内药品明细</a-divider>
              <a-table :columns="drugColumns" :data-source="result.drugs ?? []" :pagination="false" row-key="drug_id" size="small">
              </a-table>
            </div>
          </div>
        </a-card>
      </a-col>
    </a-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import ScanInput from '@/components/common/ScanInput.vue'
import { checkMixPlacement } from '@/api/locations'
import type { MixCheckResult } from '@/types/location'

const locationCode = ref('')
const loading = ref(false)
const result = ref<MixCheckResult | null>(null)

const drugColumns = [
  { title: '药品名称', dataIndex: 'drug_name', ellipsis: true },
  { title: '数量', dataIndex: 'count', width: 80 },
]

async function handleCheck(code: string) {
  if (!code) return
  locationCode.value = code
  loading.value = true
  try {
    result.value = await checkMixPlacement(code)
  } catch {
    message.error('检查失败，请确认货位码是否正确')
  } finally {
    loading.value = false
  }
}
</script>
