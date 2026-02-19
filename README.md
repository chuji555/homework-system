# 作业管理系统 (Homework System)
一个基于 Go + Gin + GORM 构建的轻量级作业管理系统，支持用户认证、作业发布/提交/批改、角色权限控制等核心功能，专为学生和管理员设计的高效作业协作工具。

## 一、项目简介
本项目是一套面向校园场景的作业管理解决方案，核心目标是简化作业发布、提交、批改的全流程：
- 管理员可发布作业、设置截止时间、批改作业、标记优秀作业；
- 学生可查看作业列表、提交作业、查看自己的提交记录；
- 基于 JWT 实现无状态认证，结合角色中间件实现细粒度权限控制；
- 采用分层架构设计，代码结构清晰，易于扩展和维护。

## 二、技术栈说明
| 技术/框架       | 版本       | 用途                     |
|----------------|------------|--------------------------|
| Go             | 1.20+      | 核心开发语言             |
| Gin            | v1.9.1     | HTTP Web 框架            |
| GORM           | v2.0       | ORM 框架，操作 MySQL 数据库 |
| JWT (golang-jwt) | v4.5.0   | 用户认证（生成/校验 Token） |
| MySQL          | 8.0+       | 关系型数据库             |
| YAML           | -          | 配置文件管理             |

## 三、项目结构说明
```
homework-system/
├── configs/                # 配置文件目录
│   └── config.yaml         # 核心配置（数据库、JWT、端口等）
├── dao/                    # 数据访问层（数据库操作）
│   ├── homework.go         # 作业相关数据库操作
│   ├── submission.go       # 提交相关数据库操作
│   └── user.go             # 用户相关数据库操作
├── handler/                # 接口处理器层（接收请求/返回响应）
│   ├── homework.go         # 作业模块接口
│   ├── submission.go       # 提交模块接口
│   └── user.go             # 用户模块接口
├── middleware/             # 中间件层（认证/权限/日志等）
│   ├── admin.go            # 管理员权限中间件
│   ├── auth.go             # JWT 认证中间件
│   └── student.go          # 学生权限中间件
├── models/                 # 数据模型层（数据库表映射）
│   ├── homework.go         # 作业模型
│   ├── submission.go       # 提交模型
│   └── user.go             # 用户模型
├── pkg/                    # 公共工具包
│   ├── errcode/            # 自定义错误码
│   │   └── errcode.go
│   └── utils/              # 通用工具函数（可选）
├── router/                 # 路由层（接口注册/分组）
│   └── router.go           # 路由初始化
├── service/                # 业务逻辑层（核心逻辑封装）
│   ├── homework.go         # 作业业务逻辑
│   ├── submission.go       # 提交业务逻辑
│   └── user.go             # 用户业务逻辑
├── go.mod                  # Go 模块依赖
├── go.sum                  # 依赖版本锁
├── main.go                 # 项目入口（启动服务）
└── README.md               # 项目说明文档
```

### 分层设计说明
1. **Router 层**：统一注册接口、划分路由组、绑定中间件，是接口的入口；
2. **Handler 层**：接收 HTTP 请求、校验参数、调用 Service 层、返回标准化响应；
3. **Service 层**：封装核心业务逻辑，是项目的核心层，解耦 Handler 和 DAO；
4. **DAO 层**：仅负责数据库操作，隔离数据库层与业务层，便于切换数据库；
5. **Models 层**：定义数据库表对应的结构体，统一数据结构；
6. **Middleware 层**：抽离通用逻辑（认证、权限、日志），复用性强。

## 四、已实现功能清单
### 1. 用户模块
| 功能                | 接口路径          | 请求方法 | 权限要求       | 说明                     |
|---------------------|-------------------|----------|----------------|--------------------------|
| 用户注册            | /user/register    | POST     | 公开           | 学生/管理员注册账号      |
| 用户登录            | /user/login       | POST     | 公开           | 返回 AccessToken/RefreshToken |
| Token 刷新          | /user/refresh     | POST     | 公开           | 用 RefreshToken 刷新 AccessToken |
| 获取用户信息        | /user/profile     | GET      | 已登录         | 获取当前登录用户的信息    |
| 退出登录            | /user/account     | DELETE   | 已登录         | 前端丢弃 Token（后端可选拉黑） |

### 2. 作业模块
| 功能                | 接口路径          | 请求方法 | 权限要求       | 说明                     |
|---------------------|-------------------|----------|----------------|--------------------------|
| 创建作业            | /homework         | POST     | 管理员         | 发布新作业，设置标题/截止时间等 |
| 修改作业            | /homework/:id     | PUT      | 管理员         | 修改指定 ID 的作业信息    |
| 删除作业            | /homework/:id     | DELETE   | 管理员         | 删除指定 ID 的作业        |
| 作业列表查询        | /homework         | GET      | 已登录         | 分页查询作业列表          |
| 作业详情查询        | /homework/:id     | GET      | 已登录         | 查询指定 ID 的作业详情    |

### 3. 提交模块
| 功能                | 接口路径                          | 请求方法 | 权限要求       | 说明                     |
|---------------------|-----------------------------------|----------|----------------|--------------------------|
| 提交作业            | /submission                       | POST     | 学生           | 提交指定作业的答案        |
| 查看我的提交        | /submission/my                    | GET      | 学生           | 查询当前学生的提交记录    |
| 按作业查提交列表    | /submission/homework/:homework_id | GET      | 管理员         | 查询指定作业的所有提交    |
| 批改作业            | /submission/:id/review            | PUT      | 管理员         | 给指定提交打分/写评语     |
| 标记优秀作业        | /submission/:id/excellent         | PUT      | 管理员         | 将指定提交标记为优秀      |
| 查看优秀作业        | /submission/excellent             | GET      | 已登录         | 所有登录用户可查看        |

## 五、进阶功能说明
当前项目未实现进阶功能

## 六、本地运行指南
### 前置条件
1. 安装 Go 1.20+（推荐 1.21）：https://golang.org/dl/
2. 安装 MySQL 8.0+，并创建数据库（如 `homework_system`）；
3. 配置环境（可选）：确保 GOPATH/GOMOD 已正确配置。

### 运行步骤
#### 1. 克隆项目（本地开发可跳过，直接打开项目目录）
```bash
git clone https://github.com/chuji555/homework-system.git
cd homework-system
```

#### 2. 配置文件修改
编辑 `configs/config.yaml`，填写本地 MySQL 信息：
```yaml
# config.yaml 示例
server:
  port: 8080  # 服务端口
mysql:
  dsn: "root:你的密码@tcp(127.0.0.1:3306)/homework_system?charset=utf8mb4&parseTime=True&loc=Local"
jwt:
  secret: "redrock-homework-system-2024"  # JWT 签名密钥
  access_expire: 2h                       # AccessToken 过期时间
  refresh_expire: 7d                      # RefreshToken 过期时间
```

#### 3. 初始化数据库
- 手动创建数据库：
  ```sql
  CREATE DATABASE IF NOT EXISTS homework_system DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
  ```
- 项目启动时，GORM 会自动根据 Models 层创建表（需在 main.go 中开启自动迁移）：
  ```go
  // main.go 中添加
  if err := dao.DB.AutoMigrate(&models.User{}, &models.Homework{}, &models.Submission{}); err != nil {
      log.Fatal("数据库自动迁移失败：", err)
  }
  ```

#### 4. 安装依赖
```bash
go mod tidy
```

#### 5. 启动服务
```bash
# 方式1：直接运行
go run main.go

# 方式2：编译后运行
go build -o homework-system main.go
./homework-system  # Windows 执行：homework-system.exe
```

#### 6. 验证运行
服务启动后，访问 `http://localhost:8080`，若返回 404 则说明服务正常启动（无根路径接口），可通过接口测试工具调用 `/user/register` 验证。

## 七、API 文档
### 接口测试方式
#### 方式1：使用 Apifox/Postman 导入（推荐）
1. 下载 API 文档 JSON 文件（可手动导出/编写）；
2. 打开 Apifox/Postman → 导入 → 选择 JSON 文件 → 即可直接测试所有接口。

#### 方式2：使用 curl 命令测试（无需下载工具）
##### 示例1：用户注册
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"test01\",\"password\":\"123456\",\"nickname\":\"测试用户\",\"department\":\"backend\"}" http://localhost:8080/user/register
```

##### 示例2：用户登录
```bash
curl -X POST -H "Content-Type: application/json" -d "{\"username\":\"test01\",\"password\":\"123456\"}" http://localhost:8080/user/login
```

##### 示例3：获取用户信息（需替换 Token）
```bash
curl -H "Authorization: Bearer 你的AccessToken" http://localhost:8080/user/profile
```

### 通用响应格式
所有接口返回标准化 JSON 响应：
```json
{
  "code": 0,        // 错误码（0=成功，其他=失败）
  "message": "",    // 提示信息
  "data": {}        // 业务数据（成功时返回，失败时为 null）
}
```

### 错误码说明
| 错误码 | 含义               | 常见场景                     |
|--------|--------------------|------------------------------|
| 0      | 成功               | 接口调用成功                 |
| 10001  | 参数错误           | 缺少必传参数、参数格式错误   |
| 10002  | Token 为空         | 认证接口未传 Token           |
| 10003  | Token 格式错误     | Token 非 Bearer 格式         |
| 10004  | Token 过期/无效    | Token 签名错误、已过期       |
| 10005  | 角色权限不足       | 学生访问管理员接口           |
| 20001  | 数据库错误         | 数据库连接失败、操作失败     |
| 20002  | 截止时间错误       | 作业截止时间早于当前时间     |
```
