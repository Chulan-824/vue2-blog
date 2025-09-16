# React Client (Vite 7 + React 18 + Tailwind CSS v4 + Ant Design 5)

该目录包含基于 Vite 7、React 18、Tailwind CSS v4 以及 Ant Design 5 组件库重新构建的前端工程。

## 可用脚本

- `npm run dev` - 启动本地开发服务器。
- `npm run build` - 产出生产环境构建文件。
- `npm run preview` - 预览生产构建结果。
- `npm run lint` - 运行 ESLint 进行代码检查。
- `npm run format` - 使用 Prettier 校验代码风格。

## 目录结构

```
client-react/
├─ index.html
├─ package.json
├─ src/
│  ├─ App.tsx
│  ├─ main.tsx
│  ├─ components/
│  ├─ layouts/
│  ├─ pages/
│  └─ styles/
└─ tailwind.config.ts
```

`src/pages/HomePage.tsx` 为示例仪表盘页面，演示了如何同时使用 Tailwind CSS 的原子类与 Ant Design 组件进行快速开发。
