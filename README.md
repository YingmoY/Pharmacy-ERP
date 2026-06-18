# PharmacyERP

PharmacyERP 是一个面向药房业务的 ERP 示例项目，包含 Go 主后端、Vue 前端、Python AI 服务和医保网关服务。项目支持药品基础资料、供应商与库位、入库、库存、扫码盘点、销售、药师审核、追溯、报表、审计日志、通知告警、AI 发票识别和医保接口网关等功能。

## 目录结构

| 路径 | 说明 |
| --- | --- |
| `PharmacyERP/` | Go + Gin 主后端服务，提供 ERP 核心 API |
| `PharmacyERP_front/` | Vue 3 + Vite 前端管理端 |
| `ai-services/` | FastAPI AI 服务，用于发票识别、药品匹配与推荐 |
| `medicare-gateway/` | Go 医保网关，隔离 ERP 与本地医保模拟接口 |
| `references/` | OpenAPI、SQL、业务文档和示例资料 |
| `MedicareAPISimulation/` | 第三方医保模拟服务，本仓库通过 `.gitignore` 整目录排除 |

## 本地依赖

- Go 1.21+，医保网关建议 Go 1.22+
- Node.js 18+ 与 pnpm
- Python 3.10+
- PostgreSQL
- RabbitMQ，可选但推荐，用于异步日志和医保网关队列

## 首次配置

本仓库不会提交本机凭据、LLM key、数据库密码、二进制运行产物和第三方医保模拟程序。首次 clone 后按需复制示例文件：

```powershell
Copy-Item PharmacyERP\configs\config.example.yaml PharmacyERP\configs\config.local.yaml
Copy-Item ai-services\.env.example ai-services\.env
Copy-Item ai-services\app\config.example.py ai-services\app\config.py
Copy-Item medicare-gateway\.env.example medicare-gateway\.env
Copy-Item medicare-gateway\internal\config\config.go.example medicare-gateway\internal\config\config.go
Copy-Item service_manager.example.py service_manager.py
```

然后把复制出来的本地文件改成自己的数据库、RabbitMQ、JWT 和 LLM 配置。

## 数据库初始化

参考 `references/` 目录中的 SQL：

- `references/public.sql`：ERP 主业务库结构与演示数据
- `references/ai.sql`：AI 服务相关表
- `references/medicare.sql`：医保网关 schema

具体导入方式取决于你的 PostgreSQL 用户和数据库名，例如：

```powershell
psql -h 127.0.0.1 -U your_db_user -d pharmacy_erp -f references\public.sql
psql -h 127.0.0.1 -U your_db_user -d pharmacy_erp -f references\ai.sql
psql -h 127.0.0.1 -U your_db_user -d pharmacy_erp -f references\medicare.sql
```

## 启动顺序

1. 启动 PostgreSQL 和 RabbitMQ。
2. 启动 AI 服务：

```powershell
cd ai-services
.\start.ps1
```

3. 启动主后端：

```powershell
cd PharmacyERP
go run .\cmd\server
```

4. 启动医保网关，可选：

```powershell
cd medicare-gateway
go run .\cmd\gateway
```

5. 启动前端：

```powershell
cd PharmacyERP_front
pnpm install
pnpm run dev
```

默认访问地址：

- 前端：http://localhost:5173
- 主后端：http://localhost:8080
- AI 服务：http://localhost:9080
- 医保网关：http://localhost:8088

## 文档入口

- 主后端：[PharmacyERP/README.md](PharmacyERP/README.md)
- 前端：[PharmacyERP_front/README.md](PharmacyERP_front/README.md)
- AI 服务：[ai-services/README.md](ai-services/README.md)
- 医保网关：[medicare-gateway/README.md](medicare-gateway/README.md)
