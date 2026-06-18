# PharmacyERP AI Service

这是 PharmacyERP 的 Python AI 子服务，基于 FastAPI，负责 AI 药品搜索、药品推荐、发票识别和供应商/药品匹配。

## API 前缀

```text
/ai/api/v1
```

常见路由包括：

- `/health`
- `/drugs`
- `/drugs/recommend/ws`
- `/invoices/recognize`
- `/matching`

## 配置

本服务使用 `ai-services/.env` 和 `app/config.py`。这两个文件可能包含数据库密码和 LLM key，已被忽略。首次运行请复制示例：

```powershell
Copy-Item .env.example .env
Copy-Item app\config.example.py app\config.py
```

然后修改：

- `DB_HOST`、`DB_PORT`、`DB_USER`、`DB_PASSWORD`、`DB_NAME`
- `LLM_BASE_URL`、`LLM_API_KEY`、`LLM_MODEL`
- `SERVICE_HOST`、`SERVICE_PORT`

## 运行

推荐使用项目自带脚本：

```powershell
.\start.ps1
```

手动启动：

```powershell
python -m venv .venv
.\.venv\Scripts\pip install -r requirements.txt
.\.venv\Scripts\python main.py
```

默认监听：

```text
http://localhost:9080
```

健康检查：

```text
http://localhost:9080/ai/api/v1/health
```

## 上传注意

不要提交 `.env`、`.venv/`、`app/config.py`、`__pycache__/` 和 `*.pyc`。
