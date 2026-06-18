<template>
  <div class="page-content">
    <a-card class="search-card" size="small">
      <a-form layout="inline" :model="searchParams" @finish="fetchList">
        <a-form-item label="发票号"><a-input v-model:value="searchParams.invoice_no" allow-clear /></a-form-item>
        <a-form-item label="供应商"><a-input v-model:value="searchParams.supplier_name" allow-clear /></a-form-item>
        <a-form-item label="状态">
          <a-select v-model:value="searchParams.status" placeholder="全部" style="width: 120px" allow-clear>
            <a-select-option value="PENDING">待核验</a-select-option>
            <a-select-option value="VERIFIED">已核验</a-select-option>
            <a-select-option value="REJECTED">已驳回</a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="时间">
          <a-range-picker v-model:value="dateRange" format="YYYY-MM-DD" @change="handleDateChange" />
        </a-form-item>
        <a-form-item>
          <a-button type="primary" html-type="submit" :loading="loading">查询</a-button>
          <a-button style="margin-left: 8px" @click="handleReset">重置</a-button>
        </a-form-item>
      </a-form>
    </a-card>

    <a-card>
      <template #title>
        AI识别发票记录
        <a-tooltip title="通过AI智能识别入库发票，自动提取药品信息，减少人工录入错误">
          <QuestionCircleOutlined style="margin-left: 4px; color: #999" />
        </a-tooltip>
      </template>
      <template #extra>
        <a-button type="primary" :icon="h(UploadOutlined)" @click="showUploadModal = true">
          上传发票
        </a-button>
      </template>
      <a-table :columns="columns" :data-source="list" :loading="loading" :pagination="pagination" row-key="id" @change="handleTableChange">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'status'">
            <StatusTag type="invoice" :value="record.status" />
          </template>
          <template v-if="column.key === 'action'">
            <a-button type="link" size="small" @click="$router.push(`/ai/invoices/${record.id}`)">查看详情</a-button>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 上传发票弹窗 -->
    <a-modal
      v-model:open="showUploadModal"
      title="上传发票识别"
      :confirm-loading="uploading"
      ok-text="开始识别"
      cancel-text="取消"
      :ok-button-props="{ disabled: !pendingFile }"
      @ok="handleUpload"
      @cancel="resetUpload"
    >
      <a-upload-dragger
        :before-upload="handleBeforeUpload"
        :file-list="fileList"
        :max-count="1"
        accept=".pdf,.jpg,.jpeg,.png"
        :multiple="false"
        @remove="resetUpload"
      >
        <p class="ant-upload-drag-icon"><InboxOutlined /></p>
        <p class="ant-upload-text">点击或拖拽发票文件到此区域</p>
        <p class="ant-upload-hint">支持 PDF、JPG、PNG 格式，单次上传一份发票</p>
      </a-upload-dragger>
      <a-alert
        v-if="uploadError"
        type="error"
        :message="uploadError"
        style="margin-top: 12px"
        closable
        @close="uploadError = ''"
      />
    </a-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h, onMounted } from 'vue'
import { QuestionCircleOutlined, UploadOutlined, InboxOutlined } from '@ant-design/icons-vue'
import { message } from 'ant-design-vue'
import StatusTag from '@/components/common/StatusTag.vue'
import { getInvoiceRecordList, uploadAndRecognizeInvoice, getInvoiceDetail } from '@/api/ai'
import type { Dayjs } from 'dayjs'

const loading = ref(false)
const list = ref<Record<string, unknown>[]>([])
const dateRange = ref<[Dayjs, Dayjs] | null>(null)
const searchParams = reactive({ invoice_no: '', supplier_name: '', status: undefined as string | undefined, start_date: '', end_date: '' })
const pagination = reactive({ current: 1, pageSize: 20, total: 0, showTotal: (t: number) => `共 ${t} 条` })

// 上传相关
const showUploadModal = ref(false)
const uploading = ref(false)
const pendingFile = ref<File | null>(null)
const fileList = ref<{ uid: string; name: string; status: string }[]>([])
const uploadError = ref('')

const columns = [
  { title: '发票号', dataIndex: 'invoice_no', width: 160 },
  { title: '供应商', dataIndex: 'supplier_name', ellipsis: true },
  { title: '识别药品数', dataIndex: 'drug_count', width: 100 },
  { title: '发票金额', dataIndex: 'total_amount', width: 120 },
  { title: '状态', key: 'status', width: 90 },
  { title: 'AI置信度', dataIndex: 'confidence', width: 100, customRender: ({ text }: { text: number }) => `${Math.round(text * 100)}%` },
  { title: '上传时间', dataIndex: 'created_at', width: 160 },
  { title: '操作', key: 'action', width: 90 },
]

function handleDateChange(dates: [Dayjs, Dayjs] | null) {
  if (dates) {
    searchParams.start_date = dates[0].format('YYYY-MM-DD')
    searchParams.end_date = dates[1].format('YYYY-MM-DD')
  } else {
    searchParams.start_date = ''
    searchParams.end_date = ''
  }
}

async function fetchList() {
  loading.value = true
  try {
    const res = await getInvoiceRecordList({ page: pagination.current, page_size: pagination.pageSize, ...searchParams })
    list.value = res.list
    pagination.total = res.total
  } finally { loading.value = false }
}

function handleReset() {
  Object.assign(searchParams, { invoice_no: '', supplier_name: '', status: undefined, start_date: '', end_date: '' })
  dateRange.value = null
  fetchList()
}

function handleTableChange(pag: typeof pagination) {
  pagination.current = pag.current
  fetchList()
}

function handleBeforeUpload(file: unknown) {
  // AntD Vue 4 的 before-upload 传入 UploadFile 包装对象，真正的 File 在 originFileObj
  const f = file as { originFileObj?: File; name?: string } & File
  const rawFile: File = f.originFileObj ?? f
  pendingFile.value = rawFile
  fileList.value = [{ uid: '-1', name: rawFile.name, status: 'done' }]
  return false // 阻止自动上传
}

function resetUpload() {
  pendingFile.value = null
  fileList.value = []
  uploadError.value = ''
}

async function handleUpload() {
  if (!pendingFile.value) return
  uploading.value = true
  uploadError.value = ''
  try {
    const record = await uploadAndRecognizeInvoice(pendingFile.value)
    showUploadModal.value = false
    resetUpload()
    await fetchList()
    // 若 AI 仍在识别中，后台轮询直到完成
    if (record.status === 'PROCESSING' || record.status === 'PENDING') {
      message.loading({ content: 'AI 识别中，请稍候…', key: 'ai-poll', duration: 0 })
      pollUntilDone(record.id)
    } else {
      message.success('发票识别完成')
    }
  } catch (e: unknown) {
    const err = e as { message?: string }
    uploadError.value = err?.message || '上传失败，请检查文件格式后重试'
  } finally {
    uploading.value = false
  }
}

function pollUntilDone(id: number, attempts = 0) {
  if (attempts > 40) { // 最多轮询 40 次（约 10 分钟）
    message.destroy('ai-poll')
    message.warning('AI 识别超时，请手动刷新查看结果')
    return
  }
  setTimeout(async () => {
    try {
      const rec = await getInvoiceDetail(id)
      if (rec.status === 'COMPLETED') {
        message.destroy('ai-poll')
        message.success('AI 识别完成，请查看结果')
        await fetchList()
      } else if (rec.status === 'FAILED') {
        message.destroy('ai-poll')
        message.error(`AI 识别失败：${rec.error_message || '未知错误'}`)
        await fetchList()
      } else {
        pollUntilDone(id, attempts + 1)
      }
    } catch {
      pollUntilDone(id, attempts + 1)
    }
  }, 15000) // 每 15 秒轮询一次
}

onMounted(fetchList)
</script>
