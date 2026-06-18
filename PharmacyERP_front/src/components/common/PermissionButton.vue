<template>
  <!-- 权限按钮：根据权限点控制显示/隐藏 -->
  <template v-if="visible">
    <a-button v-bind="$attrs" :disabled="disabled" @click="$emit('click', $event)">
      <slot />
    </a-button>
  </template>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '@/store/auth'

const props = withDefaults(
  defineProps<{
    // 所需权限码，不传则只检查登录状态
    permission?: string
    // 额外的禁用条件（业务状态机控制）
    disabled?: boolean
  }>(),
  {
    permission: undefined,
    disabled: false,
  },
)

defineEmits<{ click: [event: MouseEvent] }>()

const authStore = useAuthStore()

// 根据权限决定是否显示
const visible = computed(() => {
  if (!props.permission) return true
  return authStore.hasPermission(props.permission)
})
</script>
