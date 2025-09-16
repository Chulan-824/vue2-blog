# 后端重构规划（Go + Gin + MongoDB + JWT）

## 1. 项目现状概述
- 前端：Vue2 + Element UI，依赖 axios 调用现有 REST 接口（详见 `client/src/api/index.js`）。
- 后端：Node.js + Express + Mongoose，使用 Session + connect-mongo 维持登录态，接口集中在 `server/routes/*`。
- 数据存储：MongoDB 数据库 `blog`，集合涵盖 `article`、`articleInfo`、`user`、`message`、`visitor` 等。
- 功能：文章列表/详情/搜索、热门文章、留言及子留言、注册登录、验证码、头像上传、最近访客记录。

> 重构要求：后端迁移至 Go + Gin；保留现有 MongoDB 表结构；认证机制改为 JWT。

## 2. 重构目标
1. **技术栈升级**：将 Express 服务迁移为 Go + Gin，采用模块化分层架构，提升性能与可维护性。
2. **认证机制重构**：移除 Session，采用 JWT（Access Token + 可选 Refresh Token）完成身份认证与授权。
3. **兼容现有数据**：不修改 MongoDB 既有集合结构，确保重构后服务可以直接操作旧数据。
4. **业务能力对齐**：文章、留言、访客、上传、注册登录等业务接口与旧版保持一致或增强。
5. **工程化**：完善配置管理、日志监控、错误处理、测试体系与 CI/CD 支撑。

## 3. 技术选型与系统架构
### 3.1 服务架构
- **框架**：`github.com/gin-gonic/gin`
- **配置**：`spf13/viper`（或 `env` + YAML/JSON）
- **日志**：`uber-go/zap`（结构化日志，支持输出至 stdout/File）
- **数据库驱动**：`go.mongodb.org/mongo-driver`，使用 Context + Client Pool
- **认证**：`github.com/golang-jwt/jwt/v5`
- **校验**：`github.com/go-playground/validator/v10`
- **依赖注入（可选）**：`google/wire` 或自定义初始化

### 3.2 目录推荐
```
vue2-blog-main/
└── server-go/
    ├── cmd/
    │   └── api/
    │       └── main.go              # 应用入口
    ├── configs/                     # 配置文件（dev/staging/prod）
    ├── internal/
    │   ├── server/                  # gin Server、路由注册、全局中间件
    │   ├── handler/                 # 控制器 (article, auth, message, visitor, upload)
    │   ├── service/                 # 业务逻辑层
    │   ├── repository/              # 数据访问层，封装 Mongo 操作
    │   ├── model/                   # 结构体、DTO、验证规则
    │   ├── auth/                    # JWT 签发/校验、权限相关
    │   ├── middleware/              # 鉴权、日志、恢复、跨域
    │   ├── config/                  # 应用配置结构与加载
    │   └── pkg/                     # 通用工具（响应封装、错误码、分页等）
    ├── pkg/
    │   └── logger/                  # 日志初始化
    ├── docs/                        # Swagger/OpenAPI 等文档
    └── test/                        # 集成测试或测试数据
```

### 3.3 分层职责
- **Handler**：解析请求、调用 Service、返回统一响应（JSON + 错误码）。
- **Service**：封装业务规则、组合多个 Repository 调用、处理事务或操作链。
- **Repository**：直接操作 Mongo 集合，提供 CRUD、查询、分页等能力。
- **Model**：
  - `entity`：映射 Mongo `bson` 字段
  - `dto/request`、`dto/response`：接口输入输出
  - `vo`：对前端暴露的数据结构
- **Middleware**：鉴权、日志、恢复、跨域、限流、防刷等。

## 4. 数据库设计（保持现有结构）
以下结构均与 Node 版本保持一致，仅描述字段含义与 Go 结构体的参考写法。

### 4.1 Article 集合（`article`）
| 字段 | 类型 | 说明 |
|------|------|------|
| `_id` | ObjectID | 主键 |
| `type` | String | 文章类型（原创/转载） |
| `title` | String | 标题 |
| `content` | String | HTML 内容 |
| `tag` | String | 分类标签（HTML&Css / JS / Node / Vue&React / Other）|
| `updateDate` | Date | 更新时间 |
| `date` | Date | 创建时间 |
| `surface` | String | 封面图 URL |
| `pv` | Number | 浏览量 |
| `comment` | Array<ObjectID> | 评论引用（暂未使用） |

Go 结构示例：
```go
type Article struct {
    ID         primitive.ObjectID   `bson:"_id,omitempty"`
    Type       string               `bson:"type"`
    Title      string               `bson:"title"`
    Content    string               `bson:"content"`
    Tag        string               `bson:"tag"`
    UpdateDate time.Time            `bson:"updateDate"`
    Date       time.Time            `bson:"date"`
    Surface    string               `bson:"surface"`
    PV         int64                `bson:"pv"`
    CommentIDs []primitive.ObjectID `bson:"comment"`
}
```
索引建议：`tag + pv` 组合索引用于热门文章；`title`/`tag` 文本索引用于搜索。

### 4.2 ArticleInfo 集合（`articleInfo`）
| 字段 | 类型 | 说明 |
|------|------|------|
| `_id` | ObjectID | 主键 |
| `tags` | Array<String> | 标签列表 |
| `num` | Number | 文章总数统计 |

### 4.3 User 集合（`user`）
| 字段 | 类型 | 说明 |
|------|------|------|
| `_id` | ObjectID | 主键 |
| `user` | String | 用户名 |
| `pwd` | String | 密码哈希（需将明文替换为 BCrypt/Argon2 哈希）|
| `regDate` | Number | 注册时间戳（ms）|
| `photo` | String | 头像 URL |
| `disabled` | Boolean | 是否禁用 |
| `admin` | Boolean | 是否管理员 |

> **密码迁移**：上线前编写一次性脚本将现有明文密码转换为哈希存储，或在首次登陆时检测并迁移。

### 4.4 Message 集合（`message`）
| 字段 | 类型 | 说明 |
|------|------|------|
| `_id` | ObjectID | 主键 |
| `user` | ObjectID | 父留言用户 ID（引用 `user` 集合）|
| `content` | String | 留言 HTML 内容 |
| `date` | Date | 留言时间 |
| `children` | Array<Object> | 子留言数组 |
| `children[].user` | ObjectID | 子留言用户 ID |
| `children[].content` | String | 子留言内容 |
| `children[].reUser` | String | 被回复者用户名（字符串）|
| `children[].date` | Date | 子留言时间 |

### 4.5 Visitor 集合（`visitor`）
| 字段 | 类型 | 说明 |
|------|------|------|
| `_id` | ObjectID | 主键 |
| `user` | ObjectID | 访客用户 ID（引用 `user` 集合）|
| `date` | Date | 访问时间 |

### 4.6 上传与静态资源
- 头像上传原逻辑：保存至服务器 `public/img/upload/avatar`，更新用户 `photo` 字段。
- 重构后可改为：Gin 接收文件 -> 保存本地或云存储（七牛/OSS/S3），仍然更新 Mongo 中的 `photo` 字段。

## 5. 业务模块与接口
### 5.1 认证模块（Auth）
- `POST /api/v1/auth/register`：注册用户（含验证码校验，可保留现有验证码逻辑，或改为邮件/短信）。
- `POST /api/v1/auth/login`：用户名密码登录，返回 Access Token（必要时返回 Refresh Token）。
- `POST /api/v1/auth/refresh`：使用 Refresh Token 换取新的 Access Token（可选）。
- `POST /api/v1/auth/logout`：前端删除 Token；若需要强制退出，在服务端维护黑名单/刷新 Token。
- `GET /api/v1/auth/profile`：获取当前登录用户信息，替代旧接口 `/login/ifLogin`。

### 5.2 文章模块（Article）
- `POST /api/v1/articles/search`：关键词搜索（保留正则匹配标题/标签）。
- `GET /api/v1/articles/hot?limit=8`：热门文章（按 `pv` 排序）。
- `GET /api/v1/articles/info`：返回 `tags`, `num`。
- `GET /api/v1/articles`：分页 + 标签过滤（`?tag=Vue&React&skip=0&limit=5`）。
- `GET /api/v1/articles/:id`：文章详情 + `pv` 自增。
- （可选）`POST /api/v1/articles`、`PUT /api/v1/articles/:id`：后台管理接口。

### 5.3 留言模块（Message）
- `POST /api/v1/messages`：提交留言（登录态 required）。
- `POST /api/v1/messages/:id/replies`：提交子留言。
- `GET /api/v1/messages?skip=0&limit=5`：分页获取留言列表，含用户信息与子留言。

### 5.4 访客模块（Visitor）
- `GET /api/v1/visitors?limit=12`：最近访客列表（含用户信息）。
- `POST /api/v1/visitors`：登录成功后写入当前用户到访客列表（Service 内部调用）。

### 5.5 上传模块（Upload）
- `POST /api/v1/upload/avatar`：头像上传（登录态 required），保存文件并更新用户 `photo` 字段。

### 5.6 验证码模块（可选）
- `POST /api/v1/captcha/register`：生成注册验证码（SVG）。
- `POST /api/v1/captcha/register/verify`：验证验证码。
> 若采用无状态服务，可将验证码答案临时存储至 Redis 或使用签名 token。

## 6. JWT 认证设计
### 6.1 Token 签发
- **Access Token**
  - 载荷（Claims）：`userId`, `username`, `admin`, `iat`, `exp`
  - 有效期：15–30 分钟（可配置）
- **Refresh Token（可选）**
  - 长有效期（7–30 天）
  - 存储策略：
    - 方案 A：签发 JWT，存于 HttpOnly Cookie；若需失效控制，将 token ID 保存至 Mongo `refresh_tokens` 集合或 Redis，并设置 TTL。
    - 方案 B：仅返回 Access Token，由客户端在失效后重新登录。

### 6.2 中间件流程
1. 从 `Authorization: Bearer <token>` 或 HttpOnly Cookie 中提取 Token。
2. 验证签名（HMAC SHA256）。
3. 校验 `exp`、`iat`、黑名单（若实现）。
4. 将用户信息写入 Gin `Context`（键如 `ctx.Set("user", *UserClaims)`）。
5. 下游 Handler 读取 Context 实现鉴权/授权。

### 6.3 Token 失效策略
- 主动登出：
  - 短时 Token：客户端删除即可。
  - 长时 Token：将 Token JTI 写入黑名单（Redis + TTL）。
- 被动失效：
  - 修改密码/禁用用户时，可增加版本号 `tokenVersion`（存于 User 集合），签发 Token 时带上，校验时比较。

### 6.4 前端配合
- 登录成功后保存 Token（优先 HttpOnly Cookie + SameSite=Lax/Strict）。
- 请求 API 时在 Header 添加 `Authorization`。
- 处理 401/403 状态码：完成自动跳转至登录页或刷新 Token 流程。

## 7. 中间件与通用能力
- **日志**：请求日志（method、path、status、duration）、业务日志按模块分类。
- **异常恢复**：Gin `Recovery` + 自定义 panic 处理。
- **跨域**：开放给 Vite + React 前端（配置允许 origin/headers/methods）。
- **速率限制**（可选）：访问敏感接口（验证码、留言、登录）加限流。
- **输入验证**：使用 validator tag 对 DTO 做格式校验，统一错误返回。
- **响应包装**：定义统一响应结构 `{ code, message, data }`，方便前端处理。
- **安全**：
  - 富文本内容保存时可做 Server 端白名单过滤。
  - 上传文件限制大小/类型。
  - 搜索接口避免正则注入（使用 `$regex` 时转义特殊字符）。

## 8. 重构实施步骤
1. **准备阶段**
   - 创建 Go 模块，配置基础依赖（Gin、Mongo 驱动、JWT、Viper、Zap）。
   - 建立配置文件模板（Mongo URI、JWT Secret、Token 过期时间、上传目录等）。
   - 初始化日志、配置、Mongo 客户端，并实现健康检查接口。

2. **基础设施搭建**
   - 编写统一响应格式、错误码、验证器、分页封装。
   - 实现 JWT 模块（签发、解析、BlackList stub）。
   - 编写全局中间件：日志、Recovery、CORS、JWT、请求 ID。

3. **数据访问层**
   - 按照现有集合定义 Go 结构体，封装 Repository 方法（查询、分页、插入、更新）。
   - 配置必要索引（如 `article` 的 `tag`, `pv`; `visitor` 的 `date`）。

4. **业务模块迁移**
   1. **Auth 模块**：注册（含验证码校验逻辑）、登录（密码哈希验证 + JWT）、Profile、Refresh。
   2. **Article 模块**：列表、详情、搜索、热门、统计、`pv` 自增。
   3. **Message 模块**：留言、子留言、分页查询（`populate` 逻辑可通过 `lookup` 或后续批量查询用户信息实现）。
   4. **Visitor 模块**：登录后记录访客、访客列表。
   5. **Upload 模块**：头像上传，更新用户头像。

5. **前后端联调**
   - 提供 Swagger/OpenAPI 文档或接口说明。
   - React 前端切换到新 API，完成 Token 管理与接口适配。

6. **测试与验收**
   - 单元测试：Service/Repository 层使用 Mongo Test（或 mock）验证。
   - 集成测试：利用 `httptest` 或外部测试框架覆盖核心接口流程。
   - 手动/自动化测试：配合前端验证文章浏览、注册登录、留言、上传等流程。

7. **部署与迁移**
   - 构建 Docker 镜像（含多阶段构建）。
   - 配置环境变量：`MONGO_URI`、`JWT_SECRET`、`TOKEN_EXPIRE`、上传路径等。
   - 灰度/并行部署，验证新旧服务数据一致性。
   - 切流后下线旧 Node 服务。

## 9. 测试策略
- **单元测试**：
  - JWT 模块（签发/校验/刷新/黑名单）。
  - Service 层逻辑（留言提交、PV 自增、访客记录）。
- **集成测试**：
  - 使用临时 Mongo 实例/内存库（如 `mongo-test-server`）模拟数据库。
  - 覆盖登录 -> 获取文章列表 -> 留言 -> 回复等场景。
- **端到端测试**：
  - 前后端联调后，利用 Cypress/Playwright 编写 UI 测试。
- **性能测试**：
  - 使用 k6/ab 对热门接口（文章列表、详情）压测，评估并发能力。

## 10. 运维与安全
- **配置管理**：区分 `dev/staging/prod`，敏感信息通过环境变量或密钥管理（K8s Secret/HashiCorp Vault）。
- **日志采集**：stdout -> Loki/ELK；关键操作添加业务日志。
- **监控告警**：
  - 指标：请求成功率、延迟、Mongo 连接数、JWT 校验失败数。
  - 健康检查：`GET /healthz` 返回服务状态。
- **备份策略**：MongoDB 定期备份。
- **安全审计**：
  - 自定义管理员操作审计。
  - 定期扫描依赖漏洞（Go toolchain + `govulncheck`）。

## 11. 风险与对策
| 风险 | 说明 | 缓解措施 |
|------|------|----------|
| JWT 与 Session 行为差异 | 原前端依赖 Session 行为（自动刷新、服务器端登出） | 前端改造 Token 管理；实现短期 Access Token + 手动登出逻辑；必要时引入 Refresh Token/黑名单 |
| 明文密码遗留 | `user` 集合当前存储明文密码 | 上线前全量加密；或在登录逻辑中检测明文并转换 |
| 富文本安全 | 留言/文章内容为 HTML，存在 XSS 风险 | 在后端引入白名单过滤（例如 bluemonday），前端渲染时也做转义 |
| 子留言用户信息聚合 | Mongo `populate` 在 Go 中实现复杂 | 使用 `aggregation pipeline` + `$lookup`，或先查询留言再批量查询用户映射 |
| 上传文件管理 | 需要处理并发/安全问题 | 限制文件大小、类型；生成唯一文件名；引入 CDN/对象存储 |
| 项目切换期接口差异 | 前端 React 重构同步进行 | 为前端提供 Mock/旧版兼容层，确保切换阶段业务不受影响 |

## 12. 里程碑
1. **M1 – 基础框架搭建（1~2 周）**
   - 完成 Go 项目骨架、配置、日志、Mongo 驱动与 JWT 基础
   - 编写健康检查、基础中间件

2. **M2 – 核心业务迁移（2~3 周）**
   - 实现文章、留言、访客、上传模块及对应接口
   - 完成密码加密策略与验证码重构

3. **M3 – 前后端联调与测试（1~2 周）**
   - 与 React 前端完成联调
   - 补充单元/集成/端到端测试
   - 准备部署脚本与文档

4. **M4 – 部署上线与验收（1 周）**
   - 灰度上线，监控验证
   - 收集反馈并优化调优

## 13. 文档附录
- **现有接口参考**（Node 版本）：
  - 文章：`/article`, `/article/getShow`, `/article/getHot`, `/article/getInfo`, `/article/extend`, `/article/search`
  - 注册：`/register`, `/register/vcode`, `/register/checkVcode`
  - 登录：`/login`, `/login/ifLogin`, `/login/logout`
  - 留言：`/message/commit`, `/message/childCommit`, `/message/getList`
  - 访客：`/visitor`
  - 上传：`/upload/avatar`

- **参考库**：
  - Gin 官方示例：https://github.com/gin-gonic/examples
  - Mongo Go Driver 文档：https://pkg.go.dev/go.mongodb.org/mongo-driver
  - JWT 指南：https://github.com/golang-jwt/jwt
  - Go 项目布局参考：https://github.com/golang-standards/project-layout

---
> 本文档覆盖后端重构关键决策、数据库映射、JWT 方案以及执行路线，可作为项目实施与团队协作的基础文档。
