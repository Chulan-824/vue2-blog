# Vue2 Blog Server (Go版本)

这是原Vue2博客系统的Go语言重构版本，使用Gin框架和JWT认证替换原有的Express+Session架构。

## 技术栈

- **框架**: Gin (Go Web框架)
- **数据库**: MongoDB
- **认证**: JWT (替换Session)
- **日志**: Zap (结构化日志)
- **配置**: Viper (配置管理)
- **密码加密**: bcrypt

## 项目结构

```
server-go/
├── cmd/api/              # 应用入口
├── configs/              # 配置文件
├── internal/
│   ├── auth/            # JWT认证
│   ├── config/          # 配置模块
│   ├── handler/         # 控制器层
│   ├── middleware/      # 中间件
│   ├── model/           # 数据模型
│   ├── pkg/             # 工具包
│   ├── repository/      # 数据访问层
│   ├── server/          # 服务器配置
│   └── service/         # 业务逻辑层
├── pkg/logger/          # 日志工具
├── test/                # 测试文件
├── Dockerfile           # Docker构建文件
├── Makefile            # 构建脚本
└── go.mod              # Go模块文件
```

## 功能特性

### API接口
- ✅ 用户注册/登录 (JWT认证)
- ✅ 文章列表/详情/搜索/热门
- ✅ 留言系统 (支持子留言)
- ✅ 访客记录
- ✅ 头像上传
- ✅ 文章信息统计

### 兼容性
- 完全兼容现有MongoDB数据结构
- 支持明文密码自动迁移到bcrypt加密
- 提供新版API(/api/v1)和旧版兼容接口

### 新特性
- JWT无状态认证
- 结构化日志
- 优雅关闭
- Docker支持
- 健康检查接口

## 快速开始

### 环境要求
- Go 1.21+
- MongoDB
- (可选) Docker

### 本地开发

1. 克隆项目并进入目录
```bash
cd server-go
```

2. 安装依赖
```bash
make deps
```

3. 配置数据库
编辑 `configs/config.dev.yaml` 文件中的MongoDB连接信息

4. 运行开发服务器
```bash
make dev
```

服务器将在 http://localhost:8080 启动

### 生产部署

1. 构建应用
```bash
make build
```

2. 配置生产环境
编辑 `configs/config.yaml` 文件

3. 运行应用
```bash
./bin/vue2-blog-server
```

### Docker部署

1. 构建镜像
```bash
make docker-build
```

2. 运行容器
```bash
make docker-run
```

## API文档

### 新版API (推荐)

#### 认证
- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `POST /api/v1/auth/logout` - 用户登出
- `GET /api/v1/auth/profile` - 获取用户信息
- `POST /api/v1/auth/refresh` - 刷新Token

#### 文章
- `GET /api/v1/articles` - 文章列表
- `GET /api/v1/articles/{id}` - 文章详情
- `GET /api/v1/articles/hot` - 热门文章
- `GET /api/v1/articles/info` - 文章统计
- `POST /api/v1/articles/search` - 搜索文章

#### 留言
- `GET /api/v1/messages` - 留言列表
- `POST /api/v1/messages` - 创建留言 (需认证)
- `POST /api/v1/messages/{id}/replies` - 回复留言 (需认证)

#### 访客
- `GET /api/v1/visitors` - 访客列表

#### 上传
- `POST /api/v1/upload/avatar` - 头像上传 (需认证)

### 兼容接口
所有原有的Node.js版本接口路径仍然支持，无需修改前端代码。

## 配置说明

### 配置文件
主配置文件为 `configs/config.yaml`，开发环境可使用 `configs/config.dev.yaml`

### 环境变量
可通过环境变量覆盖配置：
- `SERVER_HOST` - 服务器地址
- `SERVER_PORT` - 服务器端口
- `DATABASE_URI` - MongoDB连接串
- `JWT_SECRET` - JWT密钥

## 开发工具

```bash
# 代码格式化
make fmt

# 静态检查
make lint

# 运行测试
make test

# 清理构建文件
make clean
```

## 数据迁移

### 密码迁移
系统会自动检测明文密码并在用户登录时转换为bcrypt加密，无需手动迁移。

### 数据结构
完全兼容现有MongoDB数据结构，无需修改现有数据。

## 监控和日志

### 健康检查
```bash
curl http://localhost:8080/healthz
```

### 日志配置
支持JSON和Console两种日志格式，可在配置文件中调整日志级别。

## 性能优化

- 使用连接池管理MongoDB连接
- JWT无状态认证减少数据库查询
- 静态文件服务优化
- 优雅关闭确保请求完整处理

## 安全特性

- JWT Token认证
- bcrypt密码加密
- CORS跨域保护
- 文件上传类型和大小限制
- SQL注入防护

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查MongoDB是否运行
   - 验证连接字符串配置

2. **JWT认证失败**
   - 确认JWT密钥配置
   - 检查Token是否过期

3. **文件上传失败**
   - 确认上传目录权限
   - 检查文件大小限制

## 贡献指南

1. Fork项目
2. 创建功能分支
3. 提交代码
4. 创建Pull Request

## 许可证

MIT License