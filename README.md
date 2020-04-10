# Go Gin Web include Casbin

include gin; jwt; casbin; gorm

### Conf

You should modify `conf/app.ini`

```
[database]
Type = mysql
User = root
Password =
Host = 
Name = 
TablePrefix = 
...
```

## How to run

```
$ cd $GOMODUEL/casbin
$go build
$casbin.exe
```

Project information and existing API

```
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:	export GIN_MODE=release
 - using code:	gin.SetMode(gin.ReleaseMode)

[GIN-debug] GET    /export/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
[GIN-debug] HEAD   /export/*filepath         --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (5 handlers)
...
[GIN-debug] GET    /api/v1/role/users        --> NULL/casbin/routers/api/v1.GetUsersForRole (7 handlers)
[GIN-debug] GET    /api/v1/role/all          --> NULL/casbin/routers/api/v1.GetAllRoles (7 handlers)
[info] start http server listening :8880

```
##Restful doc

+ 打开浏览器，在如下 url 中查看 API doc
```
http://localhost:8880
```
+ 注意使用 postman 模拟api/v1/... 请求时:
  + 需要在 header 中设置 Authorization = [token] 
  + 设置 cookie 为 mysession = [sessionid]

## Features

+ 使用中间件，通过 jwt 校验 token 认证，通过 auth 检查权限。
+ 注意：
  + 用户登录后，会生成 token 返回前台，同时会将用户 role 缓存到 session 中
  + 由于用户可能存在多种角色身份，故 session 中存储的是 []role 


## CasbinRule权限结构
   
   ```
   type CasbinRule struct {
      PType string          //策略类型:p-策略;g:角色组
      V0    string          //PType=p:角色代码||PType=g:用户ID
      V1    string          //PType=p:请求的path||PType=g:角色代码
      V2    string          //PType=p:请求的方法
      V3    string 
      V4    string 
      V5    string 
     }
   ```

## casbin rbac配置文件

### 请求的规则
+ r 是规则的名称，sub 为请求的实体，obj 为资源的名称, act 为请求的实际操作动作
```
[request_definition]
r = sub, obj, act
```

### 策略的规则
+ p 是策略的名称，sub 为请求的实体，obj 为资源的名称, act 为请求的实际操作动作
```
[policy_definition]
p = sub, obj, act
```

### 角色的定义
+ g 角色的名称，第一个位置为用户，第二个位置为角色，第三个位置为域（在多租户场景下使用）
```
[role_definition]
g = _, _
```

### 策略使用规范
+ 任意一条 policy rule 满足, 则最终结果为 allow
```
[policy_effect]
e = some(where (p.eft == allow))
```

### 规则匹配
+ 前三个用来匹配上面定义的请求的规则， 最后一个或条件为：如果实体是root 直接通过， 不验证权限
```
[matchers]
m = r.sub == p.sub && keyMatch(r.obj, p.obj) && regexMatch(r.act, p.act) || r.sub == "root"
```