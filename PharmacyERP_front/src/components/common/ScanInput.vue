<template>
  <div class="scan-input-wrapper">
    <a-input
      ref="inputRef"
      v-model:value="inputValue"
      :placeholder="placeholder"
      :disabled="disabled"
      size="large"
      allow-clear
      :class="{ 'scan-input-active': isFocused }"
      @focus="isFocused = true"
      @blur="isFocused = false"
      @pressEnter="handleSubmit"
      @change="handleChange"
    >
      <template #prefix>
        <ScanOutlined style="color: #1677ff; font-size: 18px" />
      </template>
      <template #suffix>
        <span v-if="inputValue" style="color: #999; font-size: 12px">回车确认</span>
      </template>
    </a-input>
    <!-- 扫码提示 -->
    <div v-if="hint" class="scan-hint">{{ hint }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { ScanOutlined } from '@ant-design/icons-vue'

const props = withDefaults(
  defineProps<{
    modelValue?: string
    placeholder?: string
    disabled?: boolean
    // 是否自动聚焦
    autoFocus?: boolean
    // 提交后是否自动清空
    clearAfterSubmit?: boolean
    // 输入框提示文字
    hint?: string
  }>(),
  {
    modelValue: '',
    placeholder: '请扫描或输入追溯码，回车确认',
    disabled: false,
    autoFocus: true,
    clearAfterSubmit: true,
    hint: undefined,
  },
)

const emit = defineEmits<{
  'update:modelValue': [value: string]
  submit: [value: string]
  change: [value: string]
}>()

const inputRef = ref()
const inputValue = ref(props.modelValue)
const isFocused = ref(false)

watch(() => props.modelValue, (val) => {
  inputValue.value = val
})

function handleChange() {
  emit('update:modelValue', inputValue.value)
  emit('change', inputValue.value)
}

function handleSubmit() {
  const val = inputValue.value.trim()
  if (!val) return
  emit('submit', val)
  if (props.clearAfterSubmit) {
    inputValue.value = ''
    emit('update:modelValue', '')
  }
}

// 自动聚焦
function focus() {
  inputRef.value?.focus()
}

onMounted(() => {
  if (props.autoFocus) {
    setTimeout(() => focus(), 100)
  }
})

defineExpose({ focus })
</script>

<style scoped>
.scan-input-wrapper {
  width: 100%;
}

.scan-hint {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
  padding-left: 4px;
}
</style>
