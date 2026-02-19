# 作业管理系统 API 文档
注：仅完成部分学生相关接口，以下介绍接口并未全部完成

Apifox分享链接：https://s.apifox.cn/f31732a5-3958-4c5e-9988-3984b93822ae

## 基础信息
- 接口根路径：`http://localhost:8080`
- 请求格式：所有POST/PUT请求均为`application/json`格式
- 响应格式：统一JSON格式，结构如下：
  ```json
  {
    "code": 0,        // 错误码（0=成功，其他=失败）
    "message": "",    // 提示信息
    "data": null      // 业务数据（成功时返回，失败时为null）
  }
  ```
- 错误码说明：
  | 错误码 | 含义                 |
  |--------|----------------------|
  | 0      | 成功                 |
  | 10001  | 参数错误             |
  | 10002  | Token为空/无效/过期  |
  | 10003  | 无权限（非管理员/学生）|
  | 20001  | 数据库操作失败       |
  | 20002  | 账号/密码错误        |
  | 20003  | 截止时间设置错误     |

## 一、用户模块
### 1. 用户注册
- **接口路径**：`/user/register`
- **请求方法**：POST
- **是否需要认证**：否
- **请求参数**：
  | 字段名     | 类型   | 是否必传 | 说明                 |
  |------------|--------|----------|----------------------|
  | username   | string | 是       | 用户名（唯一）|
  | password   | string | 是       | 密码                 |
  | nickname   | string | 是       | 昵称                 |
  | department | string | 是       | 部门（backend/frontend/test） |
- **请求示例**：
  ```json
  {
    "username": "test01",
    "password": "123456",
    "nickname": "测试用户",
    "department": "backend"
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "username": "test01",
      "nickname": "测试用户",
      "role": "student",
      "department": "backend",
      "department_label": "后端"
    }
  }
  ```

### 2. 用户登录
- **接口路径**：`/user/login`
- **请求方法**：POST
- **是否需要认证**：否
- **请求参数**：
  | 字段名   | 类型   | 是否必传 | 说明 |
  |----------|--------|----------|------|
  | username | string | 是       | 用户名 |
  | password | string | 是       | 密码 |
- **请求示例**：
  ```json
  {
    "username": "test01",
    "password": "123456"
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", // 2小时过期
      "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...", // 7天过期
      "user": {
        "id": 1,
        "username": "test01",
        "role": "student"
      }
    }
  }
  ```

### 3. 获取个人信息
- **接口路径**：`/user/profile`
- **请求方法**：GET
- **是否需要认证**：是（请求头带Token）
- **请求头**：
  ```
  Authorization: Bearer {access_token}
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "username": "test01",
      "nickname": "测试用户",
      "role": "student",
      "department": "backend",
      "department_label": "后端",
      "email": ""
    }
  }
  ```

### 4. 注销账号/退出登录
- **接口路径**：`/user/account`
- **请求方法**：DELETE
- **是否需要认证**：是（请求头带Token）
- **请求头**：
  ```
  Authorization: Bearer {access_token}
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "退出成功",
    "data": null
  }
  ```

## 二、作业模块
### 1. 发布作业（仅管理员）
- **接口路径**：`/homework`
- **请求方法**：POST
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **请求参数**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | title       | string | 是       | 作业标题                      |
  | description | string | 否       | 作业描述                      |
  | department  | string | 是       | 作业所属部门（backend/frontend） |
  | deadline    | string | 是       | 截止时间（格式：2026-03-01T23:59:59+08:00） |
  | allow_late  | bool   | 否       | 是否允许迟交（默认false）|
- **请求示例**：
  ```json
  {
    "title": "Go接口开发作业",
    "description": "实现用户登录接口，要求JWT认证",
    "department": "backend",
    "deadline": "2026-03-01T23:59:59+08:00",
    "allow_late": false
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "发布成功",
    "data": {
      "id": 1,
      "title": "Go接口开发作业",
      "description": "实现用户登录接口，要求JWT认证",
      "department": "backend",
      "deadline": "2026-03-01T23:59:59+08:00",
      "allow_late": false,
      "creator_id": 1,
      "created_at": "2026-02-19T10:00:00+08:00"
    }
  }
  ```

### 2. 查询作业列表
- **接口路径**：`/homework`
- **请求方法**：GET
- **是否需要认证**：是
- **请求头**：
  ```
  Authorization: Bearer {access_token}
  ```
- **请求参数（URL参数）**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | department  | string | 否       | 筛选部门（不传查所有）|
  | page        | int    | 否       | 页码（默认1）|
  | size        | int    | 否       | 每页条数（默认10）|
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "title": "Go接口开发作业",
          "department": "backend",
          "deadline": "2026-03-01T23:59:59+08:00"
        }
      ],
      "total": 1,
      "page": 1,
      "size": 10
    }
  }
  ```

### 3. 查询作业详情
- **接口路径**：`/homework/{id}`
- **请求方法**：GET
- **是否需要认证**：是
- **请求头**：
  ```
  Authorization: Bearer {access_token}
  ```
- **路径参数**：`id` - 作业ID
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "id": 1,
      "title": "Go接口开发作业",
      "description": "实现用户登录接口，要求JWT认证",
      "department": "backend",
      "deadline": "2026-03-01T23:59:59+08:00",
      "allow_late": false,
      "creator": {
        "id": 1,
        "username": "admin",
        "nickname": "管理员"
      },
      "created_at": "2026-02-19T10:00:00+08:00"
    }
  }
  ```

### 4. 修改作业（仅管理员）
- **接口路径**：`/homework/{id}`
- **请求方法**：PUT
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **路径参数**：`id` - 作业ID
- **请求参数**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | title       | string | 否       | 作业标题（不传则不修改）|
  | description | string | 否       | 作业描述（不传则不修改）|
  | deadline    | string | 否       | 截止时间（不传则不修改）|
  | allow_late  | bool   | 否       | 是否允许迟交（不传则不修改）|
- **请求示例**：
  ```json
  {
    "title": "Go接口开发作业（修改版）",
    "deadline": "2026-03-02T23:59:59+08:00"
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "修改成功",
    "data": {
      "id": 1,
      "title": "Go接口开发作业（修改版）",
      "deadline": "2026-03-02T23:59:59+08:00"
    }
  }
  ```

### 5. 删除作业（仅管理员）
- **接口路径**：`/homework/{id}`
- **请求方法**：DELETE
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **路径参数**：`id` - 作业ID
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "删除成功",
    "data": null
  }
  ```

## 三、提交模块
### 1. 提交作业（仅学生）
- **接口路径**：`/submission`
- **请求方法**：POST
- **是否需要认证**：是（学生Token）
- **请求头**：
  ```
  Authorization: Bearer {student_access_token}
  ```
- **请求参数**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | homework_id | int64  | 是       | 作业ID                        |
  | content     | string | 是       | 提交内容（作业答案/说明）|
  | file_url    | string | 否       | 附件URL（可选）|
- **请求示例**：
  ```json
  {
    "homework_id": 1,
    "content": "我完成了登录接口的开发，使用JWT认证",
    "file_url": ""
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "提交成功",
    "data": {
      "id": 1,
      "homework_id": 1,
      "user_id": 2,
      "content": "我完成了登录接口的开发，使用JWT认证",
      "file_url": "",
      "status": "unreviewed", // unreviewed/reviewed
      "created_at": "2026-02-19T11:00:00+08:00"
    }
  }
  ```

### 2. 查询我的提交（仅学生）
- **接口路径**：`/submission/my`
- **请求方法**：GET
- **是否需要认证**：是（学生Token）
- **请求头**：
  ```
  Authorization: Bearer {student_access_token}
  ```
- **请求参数（URL参数）**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | homework_id | int64  | 否       | 筛选指定作业的提交（不传查所有） |
  | page        | int    | 否       | 页码（默认1）|
  | size        | int    | 否       | 每页条数（默认10）|
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "homework_id": 1,
          "homework_title": "Go接口开发作业",
          "content": "我完成了登录接口的开发，使用JWT认证",
          "status": "unreviewed",
          "score": 0,
          "created_at": "2026-02-19T11:00:00+08:00"
        }
      ],
      "total": 1,
      "page": 1,
      "size": 10
    }
  }
  ```

### 3. 查询作业提交列表（仅管理员）
- **接口路径**：`/submission/homework/{homework_id}`
- **请求方法**：GET
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **路径参数**：`homework_id` - 作业ID
- **请求参数（URL参数）**：
  | 字段名   | 类型   | 是否必传 | 说明                          |
  |----------|--------|----------|-------------------------------|
  | page     | int    | 否       | 页码（默认1）|
  | size     | int    | 否       | 每页条数（默认10）|
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "user": {
            "id": 2,
            "username": "test01",
            "nickname": "测试用户"
          },
          "content": "我完成了登录接口的开发，使用JWT认证",
          "status": "unreviewed",
          "score": 0,
          "created_at": "2026-02-19T11:00:00+08:00"
        }
      ],
      "total": 1,
      "page": 1,
      "size": 10
    }
  }
  ```

### 4. 批改作业（仅管理员）
- **接口路径**：`/submission/{id}/review`
- **请求方法**：PUT
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **路径参数**：`id` - 提交ID
- **请求参数**：
  | 字段名   | 类型   | 是否必传 | 说明                          |
  |----------|--------|----------|-------------------------------|
  | score    | int    | 是       | 分数（0-100）|
  | comment  | string | 否       | 批改评语                      |
- **请求示例**：
  ```json
  {
    "score": 95,
    "comment": "接口逻辑正确，JWT认证实现规范，优秀！"
  }
  ```
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "批改成功",
    "data": {
      "id": 1,
      "score": 95,
      "comment": "接口逻辑正确，JWT认证实现规范，优秀！",
      "status": "reviewed"
    }
  }
  ```

### 5. 标记优秀作业（仅管理员）
- **接口路径**：`/submission/{id}/excellent`
- **请求方法**：PUT
- **是否需要认证**：是（管理员Token）
- **请求头**：
  ```
  Authorization: Bearer {admin_access_token}
  ```
- **路径参数**：`id` - 提交ID
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "标记优秀作业成功",
    "data": {
      "id": 1,
      "is_excellent": true
    }
  }
  ```

### 6. 查询优秀作业
- **接口路径**：`/submission/excellent`
- **请求方法**：GET
- **是否需要认证**：是
- **请求头**：
  ```
  Authorization: Bearer {access_token}
  ```
- **请求参数（URL参数）**：
  | 字段名      | 类型   | 是否必传 | 说明                          |
  |-------------|--------|----------|-------------------------------|
  | department  | string | 否       | 筛选部门（不传查所有）|
  | page        | int    | 否       | 页码（默认1）|
  | size        | int    | 否       | 每页条数（默认10）|
- **响应示例（成功）**：
  ```json
  {
    "code": 0,
    "message": "success",
    "data": {
      "list": [
        {
          "id": 1,
          "homework_title": "Go接口开发作业",
          "user_nickname": "测试用户",
          "content": "我完成了登录接口的开发，使用JWT认证",
          "score": 95,
          "comment": "接口逻辑正确，JWT认证实现规范，优秀！"
        }
      ],
      "total": 1,
      "page": 1,
      "size": 10
    }

  }
