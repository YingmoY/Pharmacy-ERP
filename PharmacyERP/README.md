# PharmacyERP Backend

这是 PharmacyERP 的 Go 主后端，基于 Gin、GORM、PostgreSQL、RabbitMQ 和 JWT，提供 ERP 核心业务 API。

## 主要能力

- 登录、用户、角色、权限和 RBAC
- 药品、供应商、库位等基础资料
- 入库、库存、扫码任务、盘点任务、上架与移库
- 销售、退货、药师审核、医保调用
- AI 发票识别记录与入库转换
- 追溯查询、报表、通知告警和审计日志

## 配置

本地运行默认读取：

```text
configs/config.local.yaml
```

该文件包含数据库、RabbitMQ、JWT 等敏感配置，已被忽略。首次运行请复制示例：

```powershell
Copy-Item configs\config.example.yaml configs\config.local.yaml
```

也可以通过环境变量指定配置文件：

```powershell
$env:PHARMACY_CONFIG_FILE="configs/config.local.yaml"
```

## 运行

```powershell
go mod download
go run .\cmd\server
```

开发热更新可使用 Air：

```powershell
air
```

默认监听 `0.0.0.0:8080`，主要 API 前缀为 `/api/v1`。

## 测试

```powershell
go test ./...
```

## 相关服务

- AI 服务默认地址：`http://127.0.0.1:9080`
- 医保网关默认地址：`http://localhost:8088`
- RabbitMQ 日志队列默认名：`operation_log_queue`

## 上传注意

不要提交 `configs/config.local.yaml`、`tmp/`、`.gotmp/`、运行日志和编译出的 `*.exe`。
