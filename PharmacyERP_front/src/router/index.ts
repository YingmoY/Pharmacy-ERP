import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { setupRouterGuards } from './guards'

const routes: RouteRecordRaw[] = [
  // 登录页
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/pages/login/index.vue'),
    meta: { requiresAuth: false, title: '登录' },
  },

  // 移动端扫码（独立布局）
  {
    path: '/m',
    component: () => import('@/layouts/MobileLayout.vue'),
    children: [
      {
        path: 'scan-tasks',
        name: 'MobileScanTasks',
        component: () => import('@/pages/mobile/ScanTaskList.vue'),
        meta: { title: '扫码任务' },
      },
      {
        path: 'scan-tasks/:id',
        name: 'MobileScanWorkbench',
        component: () => import('@/pages/mobile/ScanWorkbench.vue'),
        meta: { title: '扫码工作台' },
      },
    ],
  },

  // 主布局
  {
    path: '/',
    component: () => import('@/layouts/MainLayout.vue'),
    redirect: '/dashboard',
    meta: { requiresAuth: true },
    children: [
      // 首页看板
      {
        path: 'dashboard',
        name: 'Dashboard',
        component: () => import('@/pages/dashboard/index.vue'),
        meta: { title: '首页看板' },
      },

      // 收银销售
      {
        path: 'pos',
        name: 'POS',
        component: () => import('@/pages/pos/index.vue'),
        meta: { title: '收银台', permission: 'sales.orders.create' },
      },
      {
        path: 'sales/orders',
        name: 'SalesOrders',
        component: () => import('@/pages/sales/OrderList.vue'),
        meta: { title: '销售单' },
      },
      {
        path: 'sales/orders/:id',
        name: 'SalesOrderDetail',
        component: () => import('@/pages/sales/OrderDetail.vue'),
        meta: { title: '销售单详情' },
      },
      {
        path: 'sales/orders/:id/refund',
        name: 'SalesRefund',
        component: () => import('@/pages/sales/Refund.vue'),
        meta: { title: '退货处理', permission: 'sales.orders.refund' },
      },
      {
        path: 'sales/refunds',
        redirect: '/sales/orders',
      },

      // 药师审核
      {
        path: 'pharmacist/reviews',
        name: 'PharmacistReviews',
        component: () => import('@/pages/pharmacist/ReviewList.vue'),
        meta: { title: '药师审核', permission: 'pharmacist.review.read' },
      },
      {
        path: 'pharmacist/reviews/:id',
        name: 'PharmacistReviewDetail',
        component: () => import('@/pages/pharmacist/ReviewDetail.vue'),
        meta: { title: '审核详情', permission: 'pharmacist.review.read' },
      },

      // 采购入库
      {
        path: 'inbound/orders',
        name: 'InboundOrders',
        component: () => import('@/pages/inbound/OrderList.vue'),
        meta: { title: '入库单' },
      },
      {
        path: 'inbound/orders/new',
        name: 'InboundOrderNew',
        component: () => import('@/pages/inbound/OrderEdit.vue'),
        meta: { title: '新建入库单', permission: 'inbound.orders.create' },
      },
      {
        path: 'inbound/orders/:id',
        name: 'InboundOrderDetail',
        component: () => import('@/pages/inbound/OrderDetail.vue'),
        meta: { title: '入库单详情' },
      },
      {
        path: 'inbound/orders/:id/edit',
        name: 'InboundOrderEdit',
        component: () => import('@/pages/inbound/OrderEdit.vue'),
        meta: { title: '编辑入库单', permission: 'inbound.orders.update' },
      },
      {
        path: 'inbound/orders/:id/confirm',
        name: 'InboundConfirm',
        component: () => import('@/pages/inbound/ScanConfirm.vue'),
        meta: { title: '入库扫码确认', permission: 'inbound.orders.update' },
      },

      // AI 发票识别
      {
        path: 'ai/invoices',
        name: 'AIInvoices',
        component: () => import('@/pages/ai/InvoiceList.vue'),
        meta: { title: 'AI 发票识别' },
      },
      {
        path: 'ai/invoices/:id',
        name: 'AIInvoiceDetail',
        component: () => import('@/pages/ai/InvoiceDetail.vue'),
        meta: { title: '发票识别详情' },
      },

      // 库存追溯
      {
        path: 'inventory',
        name: 'Inventory',
        component: () => import('@/pages/inventory/InventoryList.vue'),
        meta: { title: '库存总览' },
      },
      {
        path: 'inventory/near-expire',
        name: 'NearExpire',
        component: () => import('@/pages/inventory/NearExpire.vue'),
        meta: { title: '近效期库存' },
      },
      {
        path: 'inventory/pending-shelving',
        name: 'PendingShelving',
        component: () => import('@/pages/inventory/PendingShelving.vue'),
        meta: { title: '待上架库存' },
      },
      {
        path: 'inventory-adjustments',
        name: 'InventoryAdjustments',
        component: () => import('@/pages/inventory/AdjustmentList.vue'),
        meta: { title: '库存调整记录' },
      },
      {
        path: 'trace',
        name: 'TraceQuery',
        component: () => import('@/pages/inventory/TraceQuery.vue'),
        meta: { title: '追溯码查询' },
      },
      {
        path: 'trace/:trace_code',
        name: 'TraceDetail',
        component: () => import('@/pages/inventory/TraceQuery.vue'),
        meta: { title: '追溯码查询' },
      },

      // 上架调拨
      {
        path: 'shelving/workbench',
        name: 'ShelvingWorkbench',
        component: () => import('@/pages/shelving/Workbench.vue'),
        meta: { title: '上架工作台', permission: 'shelving.create' },
      },
      {
        path: 'shelving/relocate',
        name: 'ShelvingRelocate',
        component: () => import('@/pages/shelving/Relocate.vue'),
        meta: { title: '货位调拨', permission: 'shelving.create' },
      },
      {
        path: 'shelving/mix-check',
        name: 'MixCheck',
        component: () => import('@/pages/shelving/MixCheck.vue'),
        meta: { title: '混放检查' },
      },

      // 盘库管理
      {
        path: 'inventory-tasks',
        name: 'InventoryTasks',
        component: () => import('@/pages/inventoryTask/TaskList.vue'),
        meta: { title: '盘库任务' },
      },
      {
        path: 'inventory-tasks/:id',
        name: 'InventoryTaskDetail',
        component: () => import('@/pages/inventoryTask/TaskDetail.vue'),
        meta: { title: '盘库详情' },
      },
      {
        path: 'inventory-tasks/:id/scan',
        name: 'InventoryTaskScan',
        component: () => import('@/pages/inventoryTask/ScanWorkbench.vue'),
        meta: { title: '盘库扫码', permission: 'inventory.tasks.scan' },
      },
      {
        path: 'inventory-tasks/:id/loss-candidates',
        name: 'LossCandidates',
        component: () => import('@/pages/inventoryTask/LossCandidates.vue'),
        meta: { title: '盘亏候选处理' },
      },
      {
        path: 'inventory-tasks/:id/misplaced',
        name: 'MisplacedItems',
        component: () => import('@/pages/inventoryTask/MisplacedItems.vue'),
        meta: { title: '错架处理' },
      },

      // 扫码作业
      {
        path: 'scan-tasks',
        name: 'ScanTasks',
        component: () => import('@/pages/scanTask/TaskList.vue'),
        meta: { title: '扫码任务' },
      },
      {
        path: 'scan-tasks/:id',
        name: 'ScanTaskDetail',
        component: () => import('@/pages/scanTask/TaskDetail.vue'),
        meta: { title: '扫码任务详情' },
      },

      // 预警通知
      {
        path: 'alerts',
        name: 'Alerts',
        component: () => import('@/pages/alert/AlertCenter.vue'),
        meta: { title: '预警中心' },
      },
      {
        path: 'notifications',
        name: 'Notifications',
        component: () => import('@/pages/alert/Notifications.vue'),
        meta: { title: '通知消息' },
      },

      // 报表中心
      {
        path: 'reports/sales',
        name: 'SalesReport',
        component: () => import('@/pages/reports/SalesReport.vue'),
        meta: { title: '销售报表' },
      },
      {
        path: 'reports/inbound',
        name: 'InboundReport',
        component: () => import('@/pages/reports/InboundReport.vue'),
        meta: { title: '入库报表' },
      },
      {
        path: 'reports/inventory',
        name: 'InventoryReport',
        component: () => import('@/pages/reports/InventoryReport.vue'),
        meta: { title: '库存报表' },
      },
      {
        path: 'reports/trace-log',
        name: 'TraceLogReport',
        component: () => import('@/pages/reports/TraceLogReport.vue'),
        meta: { title: '追溯日志报表' },
      },
      {
        path: 'reports/export-tasks',
        name: 'ExportTasks',
        component: () => import('@/pages/reports/ExportTasks.vue'),
        meta: { title: '导出任务' },
      },

      // 基础资料
      {
        path: 'master/drugs',
        name: 'MasterDrugs',
        component: () => import('@/pages/master/Drugs.vue'),
        meta: { title: '药品资料' },
      },
      {
        path: 'master/suppliers',
        name: 'MasterSuppliers',
        component: () => import('@/pages/master/Suppliers.vue'),
        meta: { title: '供应商资料' },
      },
      {
        path: 'master/locations',
        name: 'MasterLocations',
        component: () => import('@/pages/master/Locations.vue'),
        meta: { title: '货位资料' },
      },

      // 系统管理
      {
        path: 'system/users',
        name: 'SystemUsers',
        component: () => import('@/pages/system/Users.vue'),
        meta: { title: '用户管理', permission: 'users.read' },
      },
      {
        path: 'system/roles',
        name: 'SystemRoles',
        component: () => import('@/pages/system/Roles.vue'),
        meta: { title: '角色管理', permission: 'system.role.read' },
      },
      {
        path: 'system/permissions',
        name: 'SystemPermissions',
        component: () => import('@/pages/system/Permissions.vue'),
        meta: { title: '权限管理', permission: 'system.permission.read' },
      },

      // 审计日志
      {
        path: 'audit/login-logs',
        name: 'LoginLogs',
        component: () => import('@/pages/audit/LoginLogs.vue'),
        meta: { title: '登录日志', permission: 'audit.login.read' },
      },
      {
        path: 'audit/operation-logs',
        name: 'OperationLogs',
        component: () => import('@/pages/audit/OperationLogs.vue'),
        meta: { title: '操作日志', permission: 'audit.operation.read' },
      },
      {
        path: 'audit/security-events',
        name: 'SecurityEvents',
        component: () => import('@/pages/audit/SecurityEvents.vue'),
        meta: { title: '安全事件', permission: 'audit.security.read' },
      },
      {
        path: 'audit/data-change-logs',
        name: 'DataChangeLogs',
        component: () => import('@/pages/audit/DataChangeLogs.vue'),
        meta: { title: '数据变更日志', permission: 'audit.data_change.read' },
      },
    ],
  },

  // 404
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    redirect: '/dashboard',
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior: () => ({ top: 0 }),
})

setupRouterGuards(router)

export default router
