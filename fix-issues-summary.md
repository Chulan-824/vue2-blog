# Go 后端重构待修复问题汇总

- **请求方式兼容性**：`server-go/internal/server/server.go` 中旧接口（例如 `/article`、`/article/getShow`、`/message/getList`、`/visitor` 等）被改成 GET，并且 `/article/getShow` 错误路由到详情接口；旧前端 (`client/src/api/index.js`) 依旧通过 POST 访问上述路径并携带分页参数，现状会直接 404 或返回错误数据。
- **请求体解析方式不匹配**：多个 Handler 使用 `ShouldBindJSON`（如 `internal/handler/auth_handler.go`, `message_handler.go`, `article_handler.go`, `visitor_handler.go`），而原前端仍发送 `Content-Type: application/x-www-form-urlencoded`；需要支持表单解析或调整前端，否则所有核心接口都会因为“请求参数错误”失败。
- **留言子回复参数映射错误**：兼容路由中 `/message/childCommit` 仅读取 `messageId` 字段并写回 `c.Set("id", ...)`，但前端发送的是 `parentId`；导致后端无法定位父留言，子回复接口必然失败。
- **登录态检查返回结构变化**：`AuthHandler.CheckLogin` 走统一响应包装（`pkg.Success`），旧前端期望的是 `{ userInfo: ... }` 的裸对象；当前实现会让前端判断始终视为未登录，需要恢复兼容返回格式。
- **头像上传链路失效**：
  - Handler 读取字段名 `avatar`，前端组件依旧上传字段 `file` (`client/src/components/Avatar.vue`)，请求会被判定为未上传文件。
  - 服务端保存路径为 `./public/avatar`，但静态资源仅映射 `./public/img` (`internal/server/server.go`)，最终返回的 `.../img/upload/avatar/...` 无法访问。
- **功能缺失**：
  - `/article/extend` 延伸阅读接口在新项目中未实现。
  - 注册验证码接口 `/register/vcode`、`/register/checkVcode` 仅留 `TODO`，与重构计划“保留验证码逻辑”不符，导致现有注册流程不可用。

> 建议按以上问题优先级逐项修复，以确保新后端能够与现有前端平滑对接。
