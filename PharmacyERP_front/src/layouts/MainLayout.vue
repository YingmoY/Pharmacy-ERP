<template>
  <a-layout style="min-height: 100vh">
    <!-- 侧边栏 -->
    <a-layout-sider
      v-model:collapsed="appStore.collapsed"
      collapsible
      :width="220"
      :collapsed-width="64"
      style="background: #001529; overflow: auto; height: 100vh; position: fixed; left: 0; top: 0; bottom: 0; z-index: 100"
    >
      <!-- Logo -->
      <div class="logo" @click="router.push('/dashboard')">
        <span class="logo-icon">💊</span>
        <span v-if="!appStore.collapsed" class="logo-text">智慧药店 ERP</span>
      </div>

      <!-- 菜单 -->
      <a-menu
        v-model:selectedKeys="selectedKeys"
        v-model:openKeys="openKeys"
        mode="inline"
        theme="dark"
        :inline-collapsed="appStore.collapsed"
        @click="handleMenuClick"
      >
        <a-menu-item key="/dashboard">
          <template #icon><DashboardOutlined /></template>
          <span>首页看板</span>
        </a-menu-item>

        <!-- 收银销售 -->
        <a-sub-menu key="sales">
          <template #icon><ShoppingCartOutlined /></template>
          <template #title>收银销售</template>
          <a-menu-item key="/pos">
            <CreditCardOutlined />收银台
          </a-menu-item>
          <a-menu-item key="/sales/orders">销售单</a-menu-item>
          <a-menu-item key="/sales/refunds">退货处理</a-menu-item>
        </a-sub-menu>

        <!-- 药师审核 -->
        <a-menu-item key="/pharmacist/reviews">
          <template #icon><AuditOutlined /></template>
          <span>药师审核</span>
          <a-badge
            v-if="appStore.activeAlertCount > 0 && !appStore.collapsed"
            :count="pendingReviewCount"
            :overflow-count="99"
            style="margin-left: 8px"
          />
        </a-menu-item>

        <!-- 采购入库 -->
        <a-sub-menu key="inbound">
          <template #icon><InboxOutlined /></template>
          <template #title>采购入库</template>
          <a-menu-item key="/inbound/orders">入库单</a-menu-item>
          <a-menu-item key="/ai/invoices">AI 发票识别</a-menu-item>
        </a-sub-menu>

        <!-- 库存追溯 -->
        <a-sub-menu key="inventory">
          <template #icon><DatabaseOutlined /></template>
          <template #title>库存追溯</template>
          <a-menu-item key="/inventory">库存总览</a-menu-item>
          <a-menu-item key="/trace">追溯码查询</a-menu-item>
          <a-menu-item key="/inventory/near-expire">近效期库存</a-menu-item>
          <a-menu-item key="/inventory/pending-shelving">待上架库存</a-menu-item>
          <a-menu-item key="/inventory-adjustments">库存调整记录</a-menu-item>
        </a-sub-menu>

        <!-- 上架调拨 -->
        <a-sub-menu key="shelving">
          <template #icon><AppstoreOutlined /></template>
          <template #title>上架调拨</template>
          <a-menu-item key="/shelving/workbench">上架工作台</a-menu-item>
          <a-menu-item key="/shelving/relocate">货位调拨</a-menu-item>
          <a-menu-item key="/shelving/mix-check">混放检查</a-menu-item>
        </a-sub-menu>

        <!-- 盘库管理 -->
        <a-sub-menu key="inventory-tasks">
          <template #icon><FileSearchOutlined /></template>
          <template #title>盘库管理</template>
          <a-menu-item key="/inventory-tasks">盘库任务</a-menu-item>
        </a-sub-menu>

        <!-- 扫码作业 -->
        <a-sub-menu key="scan-tasks">
          <template #icon><ScanOutlined /></template>
          <template #title>扫码作业</template>
          <a-menu-item key="/scan-tasks">扫码任务</a-menu-item>
        </a-sub-menu>

        <!-- 预警通知 -->
        <a-sub-menu key="alerts">
          <template #icon><BellOutlined /></template>
          <template #title>
            预警通知
            <a-badge
              v-if="appStore.activeAlertCount > 0 && !appStore.collapsed"
              :count="appStore.activeAlertCount"
              :overflow-count="99"
              style="margin-left: 8px"
            />
          </template>
          <a-menu-item key="/alerts">预警中心</a-menu-item>
          <a-menu-item key="/notifications">通知消息</a-menu-item>
        </a-sub-menu>

        <!-- 报表中心 -->
        <a-sub-menu key="reports">
          <template #icon><BarChartOutlined /></template>
          <template #title>报表中心</template>
          <a-menu-item key="/reports/sales">销售报表</a-menu-item>
          <a-menu-item key="/reports/inbound">入库报表</a-menu-item>
          <a-menu-item key="/reports/inventory">库存报表</a-menu-item>
          <a-menu-item key="/reports/trace-log">追溯日志报表</a-menu-item>
          <a-menu-item key="/reports/export-tasks">导出任务</a-menu-item>
        </a-sub-menu>

        <!-- 基础资料 -->
        <a-sub-menu key="master">
          <template #icon><ProfileOutlined /></template>
          <template #title>基础资料</template>
          <a-menu-item key="/master/drugs">药品资料</a-menu-item>
          <a-menu-item key="/master/suppliers">供应商资料</a-menu-item>
          <a-menu-item key="/master/locations">货位资料</a-menu-item>
        </a-sub-menu>

        <!-- 系统管理 -->
        <a-sub-menu key="system">
          <template #icon><SettingOutlined /></template>
          <template #title>系统管理</template>
          <a-menu-item key="/system/users">用户管理</a-menu-item>
          <a-menu-item key="/system/roles">角色管理</a-menu-item>
          <a-menu-item key="/system/permissions">权限管理</a-menu-item>
          <a-menu-item key="/audit/login-logs">登录日志</a-menu-item>
          <a-menu-item key="/audit/operation-logs">操作日志</a-menu-item>
          <a-menu-item key="/audit/security-events">安全事件</a-menu-item>
          <a-menu-item key="/audit/data-change-logs">数据变更日志</a-menu-item>
        </a-sub-menu>
      </a-menu>
    </a-layout-sider>

    <!-- 右侧主体 -->
    <a-layout :style="{ marginLeft: appStore.collapsed ? '64px' : '220px', transition: 'all 0.2s' }">
      <!-- 顶部导航 -->
      <a-layout-header style="background: #fff; padding: 0 24px; display: flex; align-items: center; box-shadow: 0 1px 4px rgba(0,0,0,0.08); position: sticky; top: 0; z-index: 99">
        <MenuFoldOutlined
          v-if="!appStore.collapsed"
          class="trigger"
          @click="appStore.toggleCollapsed"
        />
        <MenuUnfoldOutlined
          v-else
          class="trigger"
          @click="appStore.toggleCollapsed"
        />

        <!-- 面包屑 -->
        <a-breadcrumb style="margin-left: 16px; flex: 1">
          <a-breadcrumb-item>智慧药店ERP</a-breadcrumb-item>
          <a-breadcrumb-item>{{ currentPageTitle }}</a-breadcrumb-item>
        </a-breadcrumb>

        <!-- 右侧操作区 -->
        <div style="display: flex; align-items: center; gap: 16px">
          <!-- 通知 -->
          <a-badge :count="appStore.unreadNotificationCount" :overflow-count="99">
            <BellOutlined
              style="font-size: 18px; cursor: pointer"
              @click="router.push('/notifications')"
            />
          </a-badge>

          <!-- 用户信息 -->
          <a-dropdown>
            <div style="display: flex; align-items: center; gap: 8px; cursor: pointer">
              <a-avatar style="background-color: #1677ff">
                {{ authStore.user?.real_name?.charAt(0) || 'U' }}
              </a-avatar>
              <span>{{ authStore.user?.real_name || authStore.user?.username }}</span>
            </div>
            <template #overlay>
              <a-menu>
                <a-menu-item key="profile" @click="router.push('/system/users')">
                  <UserOutlined /> 个人信息
                </a-menu-item>
                <a-menu-divider />
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined /> 退出登录
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </a-layout-header>

      <!-- 页面内容 -->
      <a-layout-content style="min-height: calc(100vh - 64px); background: #f5f5f5">
        <router-view />
      </a-layout-content>
    </a-layout>
  </a-layout>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Modal } from 'ant-design-vue'
import {
  DashboardOutlined,
  ShoppingCartOutlined,
  CreditCardOutlined,
  AuditOutlined,
  InboxOutlined,
  DatabaseOutlined,
  AppstoreOutlined,
  FileSearchOutlined,
  ScanOutlined,
  BellOutlined,
  BarChartOutlined,
  ProfileOutlined,
  SettingOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  UserOutlined,
  LogoutOutlined,
} from '@ant-design/icons-vue'
import { useAuthStore } from '@/store/auth'
import { useAppStore } from '@/store/app'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const appStore = useAppStore()

// 当前选中菜单
const selectedKeys = ref<string[]>([route.path])
// 展开的子菜单
const openKeys = ref<string[]>([])

// 从路由路径推断当前页面标题
const currentPageTitle = computed(() => route.meta.title as string || '')

// 待审核数量（从 store 读取）
const pendingReviewCount = ref(0)

// 监听路由变化，同步菜单选中状态
watch(
  () => route.path,
  (path) => {
    selectedKeys.value = [path]
    // 推断应该展开哪个子菜单
    if (path.startsWith('/sales') || path === '/pos') openKeys.value = ['sales']
    else if (path.startsWith('/pharmacist')) openKeys.value = ['pharmacist']
    else if (path.startsWith('/inbound') || path.startsWith('/ai')) openKeys.value = ['inbound']
    else if (path.startsWith('/inventory') && !path.startsWith('/inventory-task')) openKeys.value = ['inventory']
    else if (path.startsWith('/shelving')) openKeys.value = ['shelving']
    else if (path.startsWith('/inventory-task')) openKeys.value = ['inventory-tasks']
    else if (path.startsWith('/scan-task')) openKeys.value = ['scan-tasks']
    else if (path.startsWith('/alert') || path.startsWith('/notification')) openKeys.value = ['alerts']
    else if (path.startsWith('/report')) openKeys.value = ['reports']
    else if (path.startsWith('/master')) openKeys.value = ['master']
    else if (path.startsWith('/system') || path.startsWith('/audit')) openKeys.value = ['system']
  },
  { immediate: true },
)

// 点击菜单
function handleMenuClick({ key }: { key: string }) {
  router.push(key)
}

// 退出登录
function handleLogout() {
  Modal.confirm({
    title: '确认退出',
    content: '确定要退出登录吗？',
    okText: '退出',
    cancelText: '取消',
    onOk: async () => {
      await authStore.logout()
      router.push('/login')
    },
  })
}
</script>

<style scoped>
.logo {
  height: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  cursor: pointer;
  padding: 0 16px;
  overflow: hidden;
}

.logo-icon {
  font-size: 24px;
  flex-shrink: 0;
}

.logo-text {
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  white-space: nowrap;
  overflow: hidden;
}

.trigger {
  font-size: 18px;
  cursor: pointer;
  color: rgba(0, 0, 0, 0.65);
  padding: 4px;
  border-radius: 4px;
  transition: color 0.2s;
}

.trigger:hover {
  color: #1677ff;
  background: rgba(22, 119, 255, 0.1);
}
</style>
