# my-chat-web (Vue 3)

这是一个与当前 Go 后端（Gin）接口匹配的简易前端示例：

- 登录 / 注册
- 自动在请求头加入 `Authorization: Bearer <token>`
- token 过期时自动调用 `/api/v1/refresh-token` 刷新并重试请求
- 会话列表 / 联系人列表 / 群组入口
- 聊天页：拉取历史 + WebSocket 收发消息

## 启动

1) 先启动后端（默认 `http://localhost:8080`）。
2) 启动前端：

```bash
cd web
npm install
npm run dev
```

然后访问：`http://localhost:5173`

## 接口约定

- 后端统一响应：`{ code, message, data }`
- BaseURL：`/api/v1`
- 登录：`POST /login` body: `{ telephone, password }` data: `{ token, refresh_token, nickname }`
- 刷新：`POST /refresh-token` body: `{ refresh_token }` data: `{ access_token, refresh_token }`
- WebSocket：`GET /ws`（前端使用 `ws(s)://host/api/v1/ws?token=...`）

## 说明（重要）

当前后端的 login 返回里没有 `userId/uuid`，但 WS 发送消息需要 `send_id`。
因此前端在 Home/Chat 页提供了一个手动填写 userId 的输入框。

更推荐的后端改法：在 login 的 `data` 里把 user.uuid 返回给前端。

