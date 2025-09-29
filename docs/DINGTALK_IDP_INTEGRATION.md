# DingTalk IDP Integration for Zitadel

本文档描述了为Zitadel添加DingTalk（钉钉）Identity Provider支持的完整实现。

## 概述

DingTalk IDP集成允许用户通过钉钉账号登录Zitadel管理的应用程序。该实现基于OAuth 2.0协议，遵循Zitadel现有的IDP架构模式。

## 实现内容

### 1. 后端核心组件

#### 1.1 DingTalk Provider (`internal/idp/providers/dingtalk/`)
- **dingtalk.go**: 主要的Provider实现
  - 实现了`idp.Provider`接口
  - 使用OAuth 2.0协议与钉钉API交互
  - 支持自定义授权端点和用户信息端点

- **DingTalk API端点**:
  - 授权端点: `https://login.dingtalk.com/oauth2/auth`
  - Token端点: `https://api.dingtalk.com/v1.0/oauth2/userAccessToken`
  - 用户信息端点: `https://api.dingtalk.com/v1.0/contact/users/me`

#### 1.2 用户信息映射 (`User` 结构体)
```go
type User struct {
    UnionID   string              // 钉钉用户唯一ID
    Nick      string              // 用户昵称
    Email     domain.EmailAddress // 邮箱地址
    Mobile    string              // 手机号
    AvatarURL string              // 头像URL
    // ... 其他字段
}
```

**字段映射**:
- `GetID()` → `UnionID` (钉钉全局唯一用户ID)
- `GetDisplayName()` → `Nick` (用户昵称)
- `GetPreferredUsername()` → `Nick` (使用昵称作为首选用户名)
- `GetEmail()` → `Email` (邮箱地址)
- `GetPhone()` → `Mobile` (手机号)
- `GetAvatarURL()` → `AvatarURL` (头像链接)
- `GetPreferredLanguage()` → `Chinese` (默认中文)

### 2. Proto定义更新

#### 2.1 IDP类型枚举 (`proto/zitadel/idp/v2/idp.proto`)
```protobuf
enum IDPType {
  // ... 其他类型
  IDP_TYPE_DINGTALK = 13;
}
```

#### 2.2 配置消息定义
```protobuf
message DingTalkConfig {
  string client_id = 1;
  repeated string scopes = 2;
}

message IDPConfig {
  oneof config {
    // ... 其他配置
    DingTalkConfig dingtalk = 14;
  }
}
```

### 3. Domain模型扩展

#### 3.1 IDP类型支持 (`internal/domain/idp.go`)
```go
const (
    // ... 其他类型
    IDPTypeDingTalk
)

func (t IDPType) GetCSSClass() string {
    case IDPTypeDingTalk:
        return "dingtalk"
}

func (t IDPType) DisplayName() string {
    case IDPTypeDingTalk:
        return "DingTalk"
}
```

#### 3.2 V3 Domain模型 (`backend/v3/domain/id_provider.go`)
```go
type DingTalk struct {
    ClientID     string
    ClientSecret *crypto.CryptoValue
    Scopes       []string
}

type IDPDingTalk struct {
    *IdentityProvider
    DingTalk
}
```

### 4. 数据库模式扩展

#### 4.1 Projection表定义 (`internal/query/projection/idp_template.go`)
- 新增表: `IDPTemplateDingTalkTable`
- 新增列定义:
  - `DingTalkIDCol`
  - `DingTalkClientIDCol`
  - `DingTalkClientSecretCol` 
  - `DingTalkScopesCol`

#### 4.2 查询模板 (`internal/query/idp_template.go`)
```go
type DingTalkIDPTemplate struct {
    IDPID        string
    ClientID     string
    ClientSecret *crypto.CryptoValue
    Scopes       database.TextArray[string]
}
```

### 5. API转换器更新

#### 5.1 类型映射 (`internal/api/grpc/idp/v2/query.go`)
```go
func idpTypeToPb(idpType domain.IDPType) idp_pb.IDPType {
    case domain.IDPTypeDingTalk:
        return idp_pb.IDPType_IDP_TYPE_DINGTALK
}

func dingtalkConfigToPb(idpConfig *idp_pb.IDPConfig, template *query.DingTalkIDPTemplate) {
    idpConfig.Config = &idp_pb.IDPConfig_Dingtalk{
        Dingtalk: &idp_pb.DingTalkConfig{
            ClientId: template.ClientID,
            Scopes:   template.Scopes,
        },
    }
}
```

### 6. 前端支持

#### 6.1 类型映射 (`apps/login/src/lib/idp.ts`)
```typescript
case IDPType.IDP_TYPE_DINGTALK:
  return IdentityProviderType.IDENTITY_PROVIDER_TYPE_DINGTALK;
```

#### 6.2 身份提供商类型定义更新
- `proto/zitadel/settings/v2/login_settings.proto`
- `proto/zitadel/settings/v2beta/login_settings.proto`

## 配置示例

### 管理员配置
```json
{
  "name": "钉钉登录",
  "type": "DINGTALK",
  "config": {
    "client_id": "your-dingtalk-app-id",
    "scopes": ["openid", "profile", "email"]
  },
  "options": {
    "isLinkingAllowed": true,
    "isCreationAllowed": true,
    "isAutoCreation": false,
    "isAutoUpdate": true
  }
}
```

### 钉钉应用配置
1. 在钉钉开放平台创建应用
2. 配置回调URL: `https://your-zitadel-domain/ui/login/login/externalidp/callback/{idp-id}`
3. 获取App ID和App Secret
4. 配置应用权限（获取用户信息）

## OAuth 2.0流程

1. **授权请求**: 用户点击"钉钉登录"按钮
   - 重定向到: `https://login.dingtalk.com/oauth2/auth`
   - 参数: `client_id`, `response_type=code`, `redirect_uri`, `scope`, `state`

2. **授权回调**: 钉钉返回授权码
   - 回调URL: `{configured_callback_url}?code={auth_code}&state={state}`

3. **令牌交换**: 使用授权码换取访问令牌
   - POST: `https://api.dingtalk.com/v1.0/oauth2/userAccessToken`
   - 参数: `client_id`, `client_secret`, `code`, `grant_type`

4. **获取用户信息**: 使用访问令牌获取用户详情
   - GET: `https://api.dingtalk.com/v1.0/contact/users/me`
   - Header: `Authorization: Bearer {access_token}`

## 安全特性

- **令牌加密**: 客户端密钥使用Zitadel的加密系统存储
- **状态验证**: 使用CSRF保护防止跨站请求伪造
- **PKCE支持**: 可选择启用PKCE以增强安全性
- **作用域限制**: 支持限制请求的OAuth作用域

## 测试

### 单元测试
```bash
go test ./internal/idp/providers/dingtalk/...
```

### 集成测试
1. 配置测试钉钉应用
2. 运行端到端测试验证完整认证流程
3. 验证用户信息映射正确性

## 部署注意事项

1. **DNS配置**: 确保钉钉可以访问回调URL
2. **SSL证书**: 钉钉要求使用HTTPS回调URL
3. **网络访问**: 确保Zitadel服务器可以访问钉钉API端点
4. **监控配置**: 添加钉钉IDP相关的监控和日志

## 故障排除

### 常见问题
1. **回调URL不匹配**: 检查钉钉应用配置中的回调URL
2. **作用域权限不足**: 确保钉钉应用具有必要的API权限
3. **网络连接问题**: 验证到钉钉API的网络连接
4. **令牌过期**: 检查访问令牌的有效期设置

### 调试日志
启用调试模式查看详细的OAuth流程日志:
```yaml
logger:
  level: debug
  format: text
```

## 后续优化

1. **缓存支持**: 实现用户信息缓存以减少API调用
2. **批量用户同步**: 支持从钉钉批量同步用户信息
3. **组织架构映射**: 支持钉钉部门和角色的映射
4. **SSO增强**: 支持钉钉的单点登录协议
5. **移动端支持**: 优化移动端钉钉登录体验

## 相关文档

- [钉钉开放平台文档](https://developers.dingtalk.com/)
- [Zitadel IDP配置指南](https://zitadel.com/docs/guides/integrate/identity-providers)
- [OAuth 2.0 RFC 6749](https://tools.ietf.org/html/rfc6749)