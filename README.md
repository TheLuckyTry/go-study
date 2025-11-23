# 博客系统 (Blog System)

一个基于 Go-Zero 和 GORM 构建的现代化博客系统，提供完整的用户认证、文章管理和评论功能。

## 项目介绍

### 功能特性

- ✅ **用户认证系统** - JWT Token 认证
- ✅ **用户管理** - 注册、登录、个人信息
- ✅ **文章管理** - 创建、读取、更新、删除文章
- ✅ **评论系统** - 文章评论功能
- ✅ **权限控制** - 基于角色的访问控制
- ✅ **请求日志** - 完整的操作日志记录
- ✅ **错误处理** - 统一的错误处理机制

### 技术栈

- **后端框架**: Go-Zero
- **数据库**: SQLite (GORM)
- **认证**: JWT
- **API文档**: 自动生成
- **日志系统**: 结构化日志记录

## 项目结构

```
blog/
├── etc/                           # 配置文件目录
│   └── blog.yaml                 # 应用配置文件
├── internal/                     # 内部代码目录
│   ├── config/                   # 配置结构定义
│   │   └── config.go
│   ├── handler/                  # HTTP 处理器
│   │   ├── registerhandler.go   
│   │   ├── loginhandler.go		  
│   │   ├── createposthandler.go  
│   │   ├── getpostshandler.go    
│   │   ├── createcommenthandler.go  
│   │   └── getrequestlogshandler.go
│   ├── logic/                    # 业务逻辑层
│   │   ├── userlogic/
│   │   │   ├── registerlogic.go
│   │   │   ├── loginlogic.go
│   │   │   └── userinfologic.go
│   │   ├── postlogic/
│   │   │   ├── createpostlogic.go
│   │   │   ├── getpostsl ogic.go
│   │   │   ├── getpostlogic.go
│   │   │   ├── updatepostlogic.go
│   │   │   └── deletepostlogic.go
│   │   ├── commentslogic/
│   │   │   ├── createcommentlogic.go
│   │   │   └── getpostcommentslogic.go
│   │   └── requestloglogic/
│   │       └── getrequestlogsl ogic.go
│   ├── svc/                      # 服务上下文
│   │   └── servicecontext.go
│   ├── types/                    # 类型定义
│   │   └── types.go
│   └── utils/                    # 工具函数
│       ├── logger.go
│       ├── request_logger.go
│       └── jwt.go
├── model/                        # 数据模型
│   └── model.go
├── api/                          # API 定义文件
│   └── blog.api
├── main.go                       # 应用入口
└── README.md                     # 项目文档
```

## API 文档

### 认证接口

#### 用户注册

- **URL**: `POST /api/user/register`
- **认证**: 不需要
- **请求体**:

```json
{
  "username": "string",
  "password": "string", 
  "email": "string"
}
```

- **响应**:

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "id": 1,
    "username": "testuser",
    "email": "test@example.com"
  }
}
```

#### 用户登录

- **URL**: `POST /api/user/login`
- **认证**: 不需要
- **请求体**:

```json
{
  "username": "string",
  "password": "string"
}
```

- **响应**:

```json
{
  "code": 0,
  "message": "成功",
  "data": {
    "token": "jwt-token-string"
  }
}
```



### 文章接口

#### 获取文章列表

- **URL**: `GET /api/posts`
- **认证**: 不需要
- **响应**:

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "title": "文章标题",
      "content": "文章内容",
      "user_id": 1,
      "username": "作者名",
      "created_at": "2023-01-01T00:00:00Z",
      "updated_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### 获取单篇文章

- **URL**: `GET /api/posts/:id`
- **认证**: 不需要
- **响应**: 同文章列表中的单个文章对象

#### 创建文章

- **URL**: `POST /api/posts`
- **认证**: 需要 JWT Token
- **请求体**:

```json
{
  "title": "文章标题",
  "content": "文章内容"
}
```

- **响应**: 创建的文章信息

#### 更新文章

- **URL**: `PUT /api/posts/:id`
- **认证**: 需要 JWT Token (只能更新自己的文章)
- **请求体**:

```json
{
  "title": "新标题",
  "content": "新内容"
}
```

#### 删除文章

- **URL**: `DELETE /api/posts/:id`
- **认证**: 需要 JWT Token (只能删除自己的文章)

### 评论接口

#### 获取文章评论

- **URL**: `GET /api/posts/:postId/comments`
- **认证**: 不需要
- **响应**:

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "content": "评论内容",
      "user_id": 1,
      "username": "评论者",
      "post_id": 1,
      "created_at": "2023-01-01T00:00:00Z"
    }
  ]
}
```

#### 创建评论

- **URL**: `POST /api/comments/create`
- **认证**: 需要 JWT Token
- **请求体**:

```json
{
  "post_id": 1,
  "content": "评论内容"
}
```

#### 删除评论

- **URL**: `DELETE /api/comments/:id`
- **认证**: 需要 JWT Token (只能删除自己的评论)

### 管理接口

#### 获取请求日志

- **URL**: `POST /api/log/page`
- **认证**: 需要 JWT Token
- **查询参数**:
    - `page`: 页码 (默认: 1)
    - `pageSize`: 每页大小 (默认: 20)
    - `method`: 请求方法 (可选)
    - `path`: 请求路径 (可选)
    - `username`: 用户名 (可选)
    - `userId`: 用户ID (可选)
- **响应**:

```json
{
  "code": 0,
  "message": "成功",
  "data": [
    {
      "id": 1,
      "method": "POST",
      "path": "/api/user/register",
      "status_code": 200,
      "username": "testuser",
      "ip_address": "127.0.0.1",
      "duration": 150,
      "created_at": "2023-01-01 00:00:00"
    }
  ],
  "total": 100
}
```

## 快速开始

### 环境要求

- Go 1.18+
- SQLite3

### 安装和运行

1. **克隆项目**

```bash
git clone <repository-url>
cd go-study
```

2. **安装依赖**

```bash
go mod tidy
```

3. **生成代码**

```bash
goctl api go -api blog.api -dir .
```

4. **运行服务**

```bash
go run blog.go
```

5. **访问服务**

```
服务地址: http://localhost:8087
```

### 配置说明

修改 `etc/blog.yaml` 文件进行配置：

```yaml
  Name: blog
  Host: 0.0.0.0
  Port: 8087
  Timeout: 60000

  Database:
    Driver: github.com/glebarez/sqlite
    Source: blog.db

  JWT:
    AccessSecret: a-simple-jwt-secret-key
    AccessExpire: 86400 # 24小时

  Log:
    ServiceName: blog-api
    Mode: file
    Level: info
    Path: logs
    Encoding: json
    TimeFormat: "2006-01-02T15:04:05.000Z07:00"
    Stat: false
    KeepDays: 7
    MaxBackups: 3
    MaxSize: 128
```

## 开发指南

### 添加新的 API

1. 在 `blog.api` 中定义新的接口
2. 运行 `goctl api go -api blog.api -dir .` 生成代码
3. 在对应的 logic 文件中实现业务逻辑

### 日志系统

系统使用结构化的日志记录，包含：

- **系统日志**: 系统运行状态
- **业务日志**: 用户操作记录
- **错误日志**: 异常和错误信息
- **安全日志**: 安全相关事件
- **请求日志**: HTTP 请求记录（存储在数据库中）

### 错误处理

系统使用统一的错误处理机制，所有错误都会返回标准格式：

```json
{
  "code": 1001,
  "message": "错误描述"
}
```

## 部署说明

### 生产环境配置

1. 修改 JWT Secret Key
2. 配置数据库连接
3. 设置合适的日志级别
4. 配置反向代理和 SSL

### 数据库迁移

系统使用 GORM AutoMigrate 自动创建和更新数据库表结构。

