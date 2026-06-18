# Medicare Gateway

这是 PharmacyERP 的医保网关服务，部署在 ERP 主后端和本地医保模拟服务之间。它把 ERP 的业务请求包装成医保通用报文，并提供审计落库和可选 RabbitMQ 异步处理。

## 主要能力

- 对接本地医保模拟接口，默认目标为 `http://localhost:9090/csb`
- 支持签到并保存 `sign_no`
- 支持人员信息、预结算、结算、结算撤销、商品销售上传、结算清单查询等接口
- 将医保调用审计日志写入 PostgreSQL 的 `medicare` schema
- 可选 RabbitMQ 异步任务，RabbitMQ 不可用时同步接口仍可运行
- 服务日志输出到 stdout，便于本地调试或统一采集

## 配置

本地配置通过环境变量读取，示例见 `.env.example`。首次运行：

```powershell
Copy-Item .env.example .env
Copy-Item internal\config\config.go.example internal\config\config.go
```

然后修改 `.env` 和 `internal/config/config.go` 中的本地默认值。`internal/config/config.go` 可能包含数据库连接串和测试人员信息，已被根 `.gitignore` 和本目录 `.gitignore` 忽略。

常用环境变量：

| 变量 | 说明 |
| --- | --- |
| `MEDICARE_GATEWAY_ADDR` | 网关监听地址，默认 `:8088` |
| `MEDICARE_DATABASE_URL` | PostgreSQL 连接串，建议包含 `search_path=medicare` |
| `MEDICARE_BASE_URL` | 医保模拟服务地址 |
| `MEDICARE_ENABLE_RABBITMQ` | 是否启用 RabbitMQ |
| `MEDICARE_RABBIT_URL` | RabbitMQ 连接串 |
| `MEDICARE_LOG_LEVEL` | 日志级别 |

## 运行

```powershell
go mod download
go run .\cmd\gateway
```

默认监听：

```text
http://localhost:8088
```

启动时会自动执行 `internal/db/schema.sql` 初始化医保 schema，也可以手动执行：

```powershell
psql -h 127.0.0.1 -U your_db_user -d pharmacy_erp -f .\internal\db\schema.sql
```

## ERP 调用接口

签到：

```http
POST http://localhost:8088/api/sign-in
```

人员信息：

```http
POST http://localhost:8088/api/person
```

通用交易：

```http
POST http://localhost:8088/api/2101
POST http://localhost:8088/api/2102
POST http://localhost:8088/api/2103
POST http://localhost:8088/api/3505
POST http://localhost:8088/api/3201
POST http://localhost:8088/api/3202
```

请求体示例：

```json
{
  "erp_order_no": "ERP-SALE-20260612-0001",
  "async": false,
  "input": {
    "druginfo": {},
    "drugdetail": []
  }
}
```

`async=true` 主要用于结算、撤销、上传、查询类接口。RabbitMQ 可用时返回排队结果，后台 worker 调用医保接口；失败消息进入死信队列。

## 上传注意

不要提交 `.env`、`internal/config/config.go`、`gateway.exe`、`.gocache/`、`docs/generated/`、运行日志和本地测试 JSON。
