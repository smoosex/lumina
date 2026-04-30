# Lumina

中文 | [English](./README.md)

Lumina 是我的个人数字花园：包含一个作品集首页，以及一个用于发布本地 Markdown 读书笔记的栏目。

<p align="center">
  <img src="frontend/previews/preview1.png" width="45%" />
  <img src="frontend/previews/preview2.png" width="45%" />
  <img src="frontend/previews/preview3.png" width="45%" />
  <img src="frontend/previews/preview4.png" width="45%" />
  <img src="frontend/previews/preview5.png" width="45%" />
  <img src="frontend/previews/preview6.png" width="45%" />
  <img src="frontend/previews/preview7.png" width="45%" />
  <img src="frontend/previews/preview8.png" width="45%" />
</p>

当前项目是 monorepo 结构：

```text
lumina/
├── frontend/   # Vue 3 个人主页和读书笔记前端
└── backend/    # 读书笔记 Go API
```

## 功能

- 个人主页：Hero、技能、精选作品和联系方式。
- 读书笔记列表页和文章详情页。
- 本地 Markdown 上传流程，不做登录系统。
- 上传请求使用 HMAC 签名保护。
- 上传时通过 ISBN 获取书籍元数据。
- 后端下载、保存并提供书籍封面访问。

## 技术栈

- 前端：Vue 3、Vite、TypeScript、Tailwind CSS v4、Pinia、vue-i18n、GSAP、Three.js。
- 后端：Go、Fx、Gin、Gorm。
- 数据库：Supabase PostgreSQL。
- 开发工具：前端使用 Bun，后端使用 Go toolchain。

## 快速开始

安装前端依赖：

```bash
cd frontend
bun install
bun dev
```

启动后端：

```bash
cd backend
cp .env.example .env
go run ./cmd/api
```

运行后端测试：

```bash
cd backend
go test ./...
```

构建前端：

```bash
cd frontend
bun run build
```

## 配置

后端配置文件位于 `backend/.env`：

```text
SERVER_ADDR=127.0.0.1:8080
DATABASE_DSN=postgres://USER:PASSWORD@HOST:6543/postgres?sslmode=require
UPLOAD_SIGNATURE_SECRET=change-this-long-random-secret
UPLOAD_API_BASE_URL=http://127.0.0.1:8080
```

## 读书笔记

一篇读书笔记是带 frontmatter 的 Markdown 文件：

```markdown
---
slug: psychology-of-money
title: "读完《金钱心理学》，我开始少一点轻易评判别人"
isbn: "9787513941242"
status: published
tags:
  - 金钱
  - 心理学
---

# 读完《金钱心理学》，我开始少一点轻易评判别人

正文内容。
```

上传笔记：

```bash
cd backend
go run ./cmd/upload-note -file ../notes/psychology-of-money.md
```

后端会在上传时根据 ISBN 获取书籍元数据、渲染 Markdown、写入 PostgreSQL，并下载封面图片。

## 文档

- [前端 README](frontend/README.md)
- [前端 README 中文版](frontend/README.zh-CN.md)
