# PharmacyERP Frontend

这是 PharmacyERP 的前端管理端，基于 Vue 3、Vite、TypeScript、Pinia、Vue Router 和 Ant Design Vue。

## 功能页面

- 登录与权限控制
- 仪表盘
- 药品、供应商、库位管理
- 入库、库存、上架、移库、盘点
- 销售、退货、药师审核
- AI 发票识别、AI 药品推荐
- 通知告警、审计日志、报表
- 移动端扫码任务页面

## 安装与运行

```powershell
pnpm install
pnpm run dev
```

默认开发地址：

```text
http://localhost:5173
```

## 开发代理

Vite 已配置代理：

| 前端路径 | 目标服务 |
| --- | --- |
| `/api` | `http://localhost:8080` |
| `/ai/api/` | `http://localhost:9080` |

因此本地开发时需要先启动主后端和 AI 服务。

## 常用命令

```powershell
pnpm run build
pnpm run preview
pnpm run type-check
```

## 上传注意

不要提交 `node_modules/`、`dist/` 和自动生成的 `components.d.ts`。
