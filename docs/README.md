# API 文档说明

本目录包含后端核心接口的 OpenAPI 3.0 规范文件，方便通过 Swagger UI、Redoc 等工具生成可交互的在线文档。

## 文件结构

- `openapi.yaml`：Vue2 Blog Service 的完整接口描述，覆盖 Auth、Article、Message 等核心模块。

## 本地查看方式

### 使用 Docker 运行 Swagger UI

1. 确保本地已安装 Docker。
2. 在仓库根目录执行以下命令启动 Swagger UI：

   ```bash
   docker run --rm -p 8081:8080 \
     -e SWAGGER_JSON=/tmp/openapi.yaml \
     -v "$(pwd)/docs/openapi.yaml:/tmp/openapi.yaml" \
     swaggerapi/swagger-ui
   ```

3. 打开浏览器访问 [http://localhost:8081](http://localhost:8081) 即可查看文档。

### 使用 Node 工具（可选）

如果希望使用本地静态文件服务器，也可以借助 `redoc-cli` 生成 HTML 页面：

```bash
npm install -g redoc-cli
redoc-cli serve docs/openapi.yaml
```

默认会在 `http://localhost:8080` 提供交互式文档。

## 更新规范

- 修改接口后请同步更新 `openapi.yaml`，确保字段描述、请求参数与实际实现保持一致。
- 建议使用 [Swagger Editor](https://editor.swagger.io/) 对 YAML 文件进行校验，以保证语法正确。
