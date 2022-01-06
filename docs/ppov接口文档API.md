# 区块链管理员（前端）与mis-bc（区块链）通信加密

>注： 使用TCP进行通信，传输的buff前4个字节是，字节序为小端的int类型值。4个字节后面再接着是json字符串
>
>*与mis通信部分代码全部放置于Node/MinSupport.go文件中*

<u>只有mis前端发送给区块链的数据需要走加密接口，通过其他方式发送给区块链的数据不经过加密接口</u>

## 加解密通信流程及接口

### 1.通信建立

前端每次启动都会给区块链发送一个建立加密通信请求

过程如下：

1. 后台生成一个8位的随机密钥(用于aes对称加密)
2. 后台将生成的密钥使用本地的公钥进行加密
3. 将加密后的数据封装发送给区块链

 - request param:


```
 {
 "Type"    : "Setup",
 "Command": "SetUpConnection",
 "Data"   : "加密后的密钥数据"
 }

```

- response param:


```
 {
 "Code"    : 200,
 "Message": 加密后的uuid
 }
```


区块链收到之后利用本地的私钥对加密的数据进行解密得到前端生成的8位随机密钥，随机生成uuid存入SessionId字段中，利用随机密钥对其进行加密传给前端，后台解密数据得到uuid，至此前端与区块链通信加密建立成功

### 2.正常通信

 后台在这之后的每次和区块链通讯都要走SM4对称加密，默认sessionid超时时间是十分钟，十分钟内没有通讯，区块链就会失效这个sessionid，此时前端需要重新建立

- request param:


```
{
"Command": "根据去请求数据填",
"Type": "根据请求数据填",
"Data"   : "一串密文",
"IsEnc" : true
}
```


*这串密文是把正常请求的json转换成字符串，然后用SM4加密后的结果*

*只有前端与区块链通信时需要加入IsEnc字段，通过其他方式发送的数据暂时不走加密接口，无需加入此字段*

- response param:

` 加密后的返回数据`

*返回数据需要用SM4解密，解密后为下面的json结果*

## MIN-VPN与区块链通信接口 

### 1.IdentityRegistry 用户注册

 - request param:


```
 {
 "Type"        : "user-act",
 "Command"     : "Registry",
 "Prefix"      :"用于注册的前缀",
 "Timestamp"   :"时间戳",
 "Username"    :"用户名",
 "Realname"    :"用户的真实姓名",
 "Phone"       :"注册用户的手机号码",
 "IDcard"      :"注册用户的身份证号",
 "AboutMe"     :"关于我",
 "Face"        :"人脸信息",
 "Print"       :"指纹信息",
 "Iris"        :"密码",
 "Other"       :"其他信息",
 "Sig"         :"数字签名",
 "Pubkey"      :"公钥",
 "Action"      :"激活码"
 }
```


- response param:


```
 注册成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"注册成功"
 }
 注册失败
 用户已经注册
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库已经存在该NDN前缀”
 }
 手机号重复注册
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库已经存在该手机号”
 }
 公钥重复
 {
 StatusCode:400
 MessageType:“string”
 Message:“用户公钥重复，注册失败”
 }
 公钥格式错误
 {
 StatusCode:401
 MessageType:“string”
 Message:“用户公钥错误"
 }
```


### 2.IdentityDestroyByUsername 根据用户名删除用户

- request param:


```
{
"Type"    : "user-act",
"Command": "Destroy",
"Username":"用户名"
}
```


- response param:


```
 删除成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"成功"
 }
 用户名不存在删除失败
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库不存在该用户”
 }
```


### 3.IdentityDestroyByIdentityIdentifier 根据用户前缀删除用户

- request param:


```
{
"Type"    : "user-act",
"Command": "DestroyByPrefix",
"Prefix":"用户前缀"
}
```


- response param:


```
 删除成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"成功"
 }
 用户名不存在删除失败
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库不存在该用户”
 }
```


### 4.IdentifierGenerate 发布标识

- request param:


```
{
"Type"    : "user-act",
"Command": "Generate",
"Prefix":"用户前缀"
"U_identifier":
"L_identifier":
"Content_Hash":
"Timestamp":时间戳
"AboutMe":
"Other":
"Sig":
}
```


- response param:


```
 发布标识成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"成功"
 }
 发布标识失败
 用户未登录
 {
 StatusCode:400
 MessageType:“string”
 Message:“用户未登录”
 }
 用户登录超时
 {
 StatusCode:401
 MessageType:“string”
 Message:“登录超时”
 }
 发布的标识已经存在
 {
 StatusCode:402
 MessageType:“string”
 Message:“数据库已经存在该NDN标识”
 }
 需要生成的标识没有以用户前缀名为前缀
 {
 StatusCode:403
 MessageType:“string”
 Message:“待生成的标识没有以前缀名为前缀”
 }
```


### 5.IdentifierDelete 删除用户发布的标识

- request param:


```
{
"Type"    : "user-act",
"Command": "Delete",
"Prefix":"用户前缀"
"U_identifier":
"Sig":
}
```


- response param:


```
 删除标识成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"成功"
 }
 删除标识失败
 要删除的前缀不存在
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库不存在该NDN标识”
 }
 证书校验错误
 {
 StatusCode:401
 MessageType:“string”
 Message:“证书校验错误,结果为0”
 }
 证书格式错误
 {
 StatusCode:402
 MessageType:“string”
 Message:“证书格式错误”
 }
```


### 6.IdentityResetPassword 重置用户密码

- request param:


```
{
"Type"    : "user-act",
"Command": "ResetPassword",
"Username":"用户名"
"Iris":"用户的新密码"
"Previous":"用户的原密码"
}
```


- response param:


```
 更改密码成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"成功"
 }
 更改密码失败
 要修改密码的用户不存在
 {
 StatusCode:400
 MessageType:“string”
 Message:“数据库不存在该用户”
 }
 原密码错误
 {
 StatusCode:400
 MessageType:“string”
 Message:“原密码错误”
 }
```


### 7.GetAllIdentityImmutableInf 得到所有MIN用户信息

- request param:


```
{
"Type"    : "user-act",
"Command": "getUser",
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"json"
 Message:"所有用户信息数组"
 }
```


*用户信息包括Username、AboutMe、Type、Iris、Phone、Realname等字段*

### 8.GetOneIdentityMutableInfByPrefix 得到所有NDN用户信息

- request param:


```
{
"Type"    : "user-act",
"Command": "getUserInf",
"Prefix":"用户前缀"
}
```


- response param:

```
 {
 StatusCode:200
 MessageType:"json"
 Message:"用户对应的北斗定位信息"
 }
```


### 9.GetIdentityByPage 分页请求用户数据

- request param:


```
{
"Type"    : "user-act",
"Command": "getUserByPage",
"PageSize":每页大小
"PageNum":第几页
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"json"
 Message:"返回某页的用户数据"
 }
```


### 10.GetAllIdentityAllInf 分页请求用户数据

- request param:


```
{
"Type"    : "user-act",
"Command": "GetUserAllMsg",
}
```


- response param:

```
 {
 StatusCode:200
 MessageType:"json"
 Message:"返回所有的用户数据"
 }
```


### 11.GetKeyLocatorByUsername 根据用户名获取用户证书

- request param:


```
{
"Type"    : "user-act",
"Command": "GetKeyLocatorByUsername",
"Username":"用户名"
}

```

- response param:


```
 获取证书成功
 {
 StatusCode:200
 MessageType:"string"
 Message:"返回用户证书"
 }
 获取证书失败
 用户不存在
 {
 StatusCode:400
 MessageType:"string"
 Message:"not have such user"
 }
 服务器故障
 {
 StatusCode:500
 MessageType:"string"
 Message:"server error"
 }
```


### 12.RevokeIdentityCert 注销用户证书

- request param:


```
{
"Type"    : "user-act",
"Command": "RevokeIdentityCert",
"Username":"用户名"
}
```


- response param:


```
注销证书成功
{
StatusCode:200
MessageType:"string"
Message:"success"
}
注销证书失败
{
StatusCode:400
MessageType:"string"
Message:"not have such user"
}
```


### 13.GetOneIdentityByUsername 通过用户名得到该用户的数据

- request param:


```
{
"Type"    : "user-act",
"Command": "GetUserByUsername",
"Username":"用户名"
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"json"
 Message:"返回该用户数据"
 }
```


### 14.MINgetUserByUsername 通过用户名得到该用户的数据

- request param:

```
{
"Type"    : "user-act",
"Command": "getUserByUsername",
"Username":"用户名"
}
```


- response param:


```
{
 StatusCode:200
 MessageType:"json"
 Message:"返回该用户数据"
 }
```


### 15.GetOneIdentityPublicKey 得到MIN用户的公钥

- request param:


```
{
"Type"    : "user-act",
"Command": "GetPubkey",
"Prefix":"用户前缀"
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"string"
 Message:"返回该用户公钥"
 }
```


### 16.查询函数

#### 1.MINQuery 查询用户所有标识

- request param:


```
{
"Type"    : "user-act",
"Command": "Query",
"QueryCode":1,
"Prefix":"用户前缀"
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"string"
 Message:"返回该用户所有标识"
 }
```


*返回的信息有 U_identifier和AboutMe字段*

#### 2.MINQuery3 查询用户face和密码

- request param:


```
{
"Type"    : "user-act",
"Command": "Query",
"QueryCode":3,
"Prefix":"用户前缀"
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"string"
 Message:"返回该用户face和密码"
 }

```

*返回的信息有 Face、Print和Iris字段*

#### 3.MINQuery4 查询用户信息是否存在于数据库

- request param:


```
{
"Type"    : "user-act",
"Command": "Query",
"QueryCode":4,
"Prefix":"用户前缀"
}
```


- response param:


```
 {
 StatusCode:200
 MessageType:"string"
 Message:"true/false"
 }
```


### 17.MINLogin 用户登录

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Login",
 "Prefix":"用户前缀",
 "Password":"用户密码",
 "Timestamp":"时间戳"
 }
```


- response param:

```

登陆成功
{
StatusCode:200,
MessageType:"string",
Message:"登录成功"
}
密码错误
{
StatusCode:400,
MessageType:"string",
Message:"密码错误"
}
用户不存在
{
StatusCode:401,
MessageType:"string",
Message:"用户不存在"
}
```


### 18.IdentifierResolve 解析用户发布的标识

- request param:


```
 {
 "Type"    : "user-act",
 "Command": "Resolve",
 "U_identifier":"用户发布的标识",
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"标识解析结果"
  }
```


### 19.MINMachine 设备注册

- request param:


```
 {
 “Type” : “user-act”,
 “Command”: “AddMachine”,
 “Prefix”:“用于注册的前缀”
 “Timestamp”:“时间戳”
 “Username”:“机器名”
 “Realname”:“”
 “Phone”:“”
 “IDcard”:“”
 “AboutMe”:""
 “Face”:""
 “Print”:""
 “Iris”:""
 “Other”:""
 “Sig”:""
 “Pubkey”:“公钥”
 “Action”:“激活码”
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"机器记录成功"
  }
```


### 20.GetMachine 获取所有设备的信息

- request param:


```
 {
 “Type” : “user-act”,
 “Command”: “GetMachine”,
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"json",
  Message:"所有机器的信息"
  }
```


### 21.RemoveMachine 删除设备

- request param:


```
 {
 “Type” : “user-act”,
 “Command”: “DeleteMachine”,
 “Prefix”:“设备的前缀”
 }
```


- response param:


```
 删除成功
 {
 StatusCode:200,
 MessageType:"string",
 Message:"成功"
 }
 删除失败
 该设备不存在
 {
 StatusCode:400,
 MessageType:"string",
 Message:"数据库不存在该machine"
 }
```


## 前端获取指定高度区块组接口

### 1.GetBGOfCertainHeight 从区块链中获取某高度区块组信息

- request param:


```
 {
 “Type” : “Front-end”,
 “Command”: “getHeight”,
 “Height”:“区块组高度”
 }
```


- response param:

`  区块组信息`

## 网络接口部分

### 1.MINLog 记录用户日志

- request param:


```
 {
 "Type"    : "network",
 "Command": "Log",
 "Prefix":"用户前缀",
 "Action":"激活码",
 "Sig":"",
 "Timestamp":"时间戳"
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"日志记录成功"
  }
```


### 2.GetAllLogByIdentity 通过用户名得到用户所有的日志信息

- request param:


```
 {
 "Type"    : "network",
 "Command": "GetAllLogByIdentity",
 "Username":"用户名",
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"json"
  Message:"用户所有的日志信息"
  }
```


### 3.GetAllLogByTimestamp 得到时间段内所有的日志信息

- request param:


```
 {
 "Type"    : "network",
 "Command": "GetAllLogByTimestamp",
 "Start":"时间段起始时间",
 "End":"时间段结束时间",
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"该时间段内所有的日志信息"
  }
```


### 4.GetPageLogByIdentity 根据用户名分页请求用户日志信息

- request param:


```
 {
 "Type"    : "network",
 "Command": "GetPageLogByIdentity",
 "Username":"用户名",
 "PageSize":"页面大小"
 "PageNum":"页码"
 }

```

- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"用户在某页所有的日志信息"
  }
```


### 5.GetPageLogByTimestamp 分页得到某时间段内的所有日志信息

- request param:


```
 {
 "Type"    : "network",
 "Command": "GetPageLogByTimestamp",
 "Start":”时间段起始时间“
 "End":”时间段结束时间“
 "PageSize":"页面大小"
 "PageNum":"页码"
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"该时间段内所有日志信息"
  }
```


### 6.GetPageHighLogByTimestamp 获取某段时间内日志告警信息

- request param:


```
 {
 "Type"    : "network",
 "Command": "GetPageHighLogByTimestamp"    
 "Start":”时间段起始时间“
 "End":”时间段结束时间“
 "PageSize":"页面大小"
 "PageNum":"页码"
 }
```


- response param:


```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"时间段内的日志告警信息"
  }
```


### 7.AddManagement 添加管理域

- request param:


```
 {
 "Type"    : "network",
 "Command": "AddManagement",
 "Prefix":"管理域前缀"
 "Action":""
 "PubList":""
 "AboutMe":""
 }
```


- response param:


```
  管理域添加成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"上级域同意添加下级管理域"
  }
  管理域添加失败
  数据库不存在该管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库不存在该管理域"
  }
  管理域已经被分配
  {
  StatusCode:400,
  MessageType:"string",
  Message:"管理域已经被分配"
  }
```


### 8.DeleteManagement 删除管理域

- request param:


```
 {
 "Type"    : "network",
 "Command": "DeleteManagement",
 "Prefix":"用户发布的标识",
 "Action":"激活码",
 "PubList":"",
 "AboutMe":""
 }

```

- response param:


```
  管理域删除成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"上级域同意删除下级管理域"
  }
  管理域删除失败
  数据库不存在该管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库不存在该管理域"
  }
  管理域已经被分配
  {
  StatusCode:400,
  MessageType:"string",
  Message:"管理域已经被分配"
  }
```


## 白名单相关的接口

### 1.whitePaperAdd 添加白名单

- request param:


```
 {
 "Type"    : "whiteUser",
 "Command": "addWhiteUser",
 "Phone":"手机号",
 "Action":"激活码"
 "Realname":"真实姓名"
 "IDcard":"身份证号"
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"白名单添加成功"
  }
```


### 2.whitePaperGetAll 获取所有白名单

- request param:


```
 {
 "Type"    : "whiteUser",
 "Command": "getAllUser",
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"所有白名单用户"
  }
```


### 3.getWhiteUserByPhone 通过手机号获取白名单用户信息

- request param:


```
 {
 "Type"    : "whiteUser",
 "Command": "getUserByPhone",
 "Phone":"用户手机号"
 }
```

- response param:

```
  获取白名单成功
  {
  StatusCode:200,
  MessageType:"json",
  Message:"白名单用户信息"
  }
  获取白名单失败
  {
  StatusCode:400,
  MessageType:"json",
  Message:"空数据"
  }
```


### 4.WhitePaperDelete 删除白名单

- request param:


```
 {
 "Type" : "whiteUser",
 "Command": "deleteUser",
 "Phone":"用户手机号"
 }

```

- response param:


```
  删除白名单成功
  {
  StatusCode:200,
  MessageType:"json",
  Message:"成功"
  }
  删除白名单失败
  {
  StatusCode:400,
  MessageType:"json",
  Message:"数据库不存在该user"
  }

```

#  vpn-management（后台）与ppov（区块链）不加密走NDN接口

*该部分代码全部放置于Node/NdnSupport.go文件中，通过ndn路由转发不走加密接口*

## MIN-VPN与区块链通信接口 

### 1.IdentityRegistry 用户注册

 - request param:


```
  {
  "Type"    : "user-act",
  "Command": "Registry",
  "Prefix":"用于注册的前缀"
  "Level":"指示几级链"
  "Timestamp":"时间戳"
  "Username":"用户名"
  "Realname":"用户的真实姓名"
  "Phone":"注册用户的手机号码"
  "IDcard":"注册用户的身份证号"
  "AboutMe":"关于我"
  "Face":"人脸信息"
  "Print":"指纹信息"
  "Iris":"密码"
  "Other":"其他信息"
  "Sig":"数字签名"
  "Pubkey":"公钥"
  "Action":"激活码"
  "EMail":"注册邮箱"
  }
```


*除了以上信息，走ndn接口进行注册的用户信息中还加入了北斗信息，储存在前缀中，信息格式如下*


```
  {
  "MID" :"手机设备号"
  "Mac":"Mac地址"
  "BeiDou":"北斗位置"
  "GPS":"GPS定位信息"
  "SIM":"SIM卡标志信息"
  }
```


- response param:


```
  注册成功
  {
  StatusCode:200
  MessageType:"string"
  Message:"注册成功"
  }
  注册失败
  用户已经注册
  {
  StatusCode:400
  MessageType:“string”
  Message:“数据库已经存在该NDN前缀”
  }
  手机号重复注册
  {
  StatusCode:400
  MessageType:“string”
  Message:“数据库已经存在该手机号”
  }
  公钥重复
  {
  StatusCode:400
  MessageType:“string”
  Message:“用户公钥重复，注册失败”
  }
  公钥格式错误
  {
  StatusCode:401
  MessageType:“string”
  Message:“用户公钥错误"
  }
```


###  2.MINLogin 用户登录

- request param:


```
 {
 "Type"    : "user-act",
 "Command": "Login",
 "Prefix":"用户前缀",
 "Password":"用户密码",
 "Timestamp":"时间戳"
 }
```


*除了以上信息，走ndn接口进行登录的用户信息中还加入了北斗信息，储存在前缀中，信息格式如下*


```
  {
  "MID" :"手机设备号"
  "Mac":"Mac地址"
  "BeiDou":"北斗位置"
  "GPS":"GPS定位信息"
  "SIM":"SIM卡标志信息"
  }
```


- response param:

```
 登陆成功
 {
 StatusCode:200,
 MessageType:"string",
 Message:"登录成功"
 }
 密码错误
 {
 StatusCode:400,
 MessageType:"string",
 Message:"密码错误"
 }
 用户不存在
 {
 StatusCode:401,
 MessageType:"string",
 Message:"用户不存在"
 }
 用户附加设备信息格式错误（北斗信息）
 {
 StatusCode:402,
 MessageType:"string",
 Message:"用户附加设备信息格式错误"
 }
 用户附加设备信息字段不全（北斗信息）
 {
 StatusCode:403,
 MessageType:"string",
 Message:"用户附加设备信息字段不全"
 }
 用户被封禁
 {
 StatusCode:404,
 MessageType:"string",
 Message:"用户被封禁"
 }
 MID不一致
 {
 StatusCode:404,
 MessageType:"string",
 Message:"MID不一致"
 }
 Mac不一致
 {
 StatusCode:405,
 MessageType:"string",
 Message:"Mac不一致"
 }
 SIM不一致
 {
 StatusCode:406,
 MessageType:"string",
 Message:"SIM不一致"
 }
```

### 3.IdentifierGenerate 发布标识

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Generate",
 "Prefix":"用户前缀"
 "U_identifier":
 "L_identifier":
 "Content_Hash":
 "Timestamp":时间戳
 "AboutMe":
 "Other":
 "Sig":
 }
```

- response param:

```
  发布标识成功
  {
  StatusCode:200
  MessageType:"string"
  Message:"成功"
  }
  发布标识失败
  用户未登录
  {
  StatusCode:400
  MessageType:“string”
  Message:“用户未登录”
  }
  用户登录超时
  {
  StatusCode:401
  MessageType:“string”
  Message:“登录超时”
  }
  发布的标识已经存在
  {
  StatusCode:402
  MessageType:“string”
  Message:“数据库已经存在该NDN标识”
  }
  需要生成的标识没有以用户前缀名为前缀
  {
  StatusCode:403
  MessageType:“string”
  Message:“待生成的标识没有以前缀名为前缀”
  }
```

### 4.GetAllIdentityImmutableInf 得到所有MIN用户信息

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "getUser",
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"json"
  Message:"所有用户信息数组"
  }
```

*用户信息包括Username、AboutMe、Type、Iris、Phone、Realname等字段*

### 5.GetOneIdentityPublicKey 得到MIN用户的公钥

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "GetPubkey",
 "Prefix":"用户前缀"
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"string"
  Message:"返回该用户公钥"
  }
```

### 6.IdentifierDelete 删除用户发布的标识

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Delete",
 "Prefix":"用户前缀"
 "U_identifier":
 "Sig":
 }
```

- response param:

```
  删除标识成功
  {
  StatusCode:200
  MessageType:"string"
  Message:"成功"
  }
  删除标识失败
  要删除的前缀不存在
  {
  StatusCode:400
  MessageType:“string”
  Message:“数据库不存在该NDN标识”
  }
  证书校验错误
  {
  StatusCode:401
  MessageType:“string”
  Message:“证书校验错误,结果为0”
  }
  证书格式错误
  {
  StatusCode:402
  MessageType:“string”
  Message:“证书格式错误”
  }
```

### 7.IdentityDestroyByIdentityIdentifier 根据用户前缀删除用户

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "DestroyByPrefix",
 "Prefix":"用户前缀"
 }
```

- response param:

```
  删除成功
  {
  StatusCode:200
  MessageType:"string"
  Message:"成功"
  }
  用户名不存在删除失败
  {
  StatusCode:400
  MessageType:“string”
  Message:“数据库不存在该用户”
  }
```

###  8.IdentifierResolve 解析用户发布的标识

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Resolve",
 "U_identifier":"用户发布的标识",
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"标识解析结果"
  }
```

### 9.GetOneIdentityMutableInfByPrefix 得到所有NDN用户信息

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "getUserInf",
 "Prefix":"用户前缀"
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"json"
  Message:"用户对应的北斗定位信息"
  }
```

###  10.查询函数

#### 1.MINQuery 查询用户所有标识

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Query",
 "QueryCode":1,
 "Prefix":"用户前缀"
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"string"
  Message:"返回该用户所有标识"
  }
```

*返回的信息有 U_identifier和AboutMe字段*

#### 2.MINQuery3 查询用户face和密码

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Query",
 "QueryCode":3,
 "Prefix":"用户前缀"
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"string"
  Message:"返回该用户face和密码"
  }
```

*返回的信息有 Face、Print和Iris字段*

#### 3.MINQuery4 查询用户信息是否存在于数据库

- request param:

```
 {
 "Type"    : "user-act",
 "Command": "Query",
 "QueryCode":4,
 "Prefix":"用户前缀"
 }
```

- response param:

```
  {
  StatusCode:200
  MessageType:"string"
  Message:"true/false"
  }
```

## 网络接口模块

### 1.minLogThroughNDN min用户登录

- request param:

```
 {
 "Type"    : "network",
 "Command": "Log",
 "Prefix":"用户前缀",
 "Action":"激活码",
 "Sig":"",
 "Timestamp":"时间戳"
 "Level":"指示几级链"
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"日志记录成功"
  }
```

### 2.minAddDomain 添加管理域

- request param:

```
 {
 "Type"    : "network",
 "Command": "AddDomain",
 "Prefix":"用户前缀",
 }
```

- response param:

```
  管理域添加成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理添加成功"
  }
  管理域添加失败
  数据库已经存在该管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库已经存在该NDN前缀管理"
  }
  添加的管理域不是一级域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"本域无法处理一级前缀"
  }
```

### 3.minAddSubDomain 添加子管理域

- request param:

```
 {
 "Type"    : "network",
 "Command": "AddSubDomain",
 "Prefix":"用户前缀",
 "Action":"邀请码"
 }
```

- response param:

```
  管理域添加成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理添加成功"
  }
  管理域添加失败
  数据库已经存在该管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库已经存在该NDN前缀管理"
  }
  添加的管理域是一级域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"本域无法处理一级前缀"
  }
  服务器错误
  {
  StatusCode:500,
  MessageType:"string",
  Message:"服务器错误"
  }
```

  *在添加子级管理域过程中，子管理域会向上级管理域发送注册请求，得到上级域通过则注册成功*

### 4.minShowTopo 查询管理域拓扑结构

- request param:

```
 {
 "Type"    : "network",
 "Command": "ShowTopo",
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理域拓扑信息"
  }
```

### 5.minDeleteDomain 删除管理域

- request param:

```
 {
 "Type"    : "network",
 "Command": "DeleteDomain",
 "Prefix":"用户前缀",
 }
```

- response param:

```
  管理域删除成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理删除成功"
  }
  管理域删除失败
  数据库不存在该管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库不存在该NDN前缀管理"
  }
  删除的管理域不是一级域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"本域无法处理一级前缀"
  }
```

### 6.minDeleteSubDomain 删除子管理域

- request param:

```
 {
 "Type"    : "network",
 "Command": "DeleteSubDomain",
 "Prefix":"用户前缀",
 "Action":"激活码"
 }
```

- response param:

```
  子管理域删除成功
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理删除成功"
  }
  管理域删除失败
  数据库不存在该子管理域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"数据库不存在该NDN前缀管理"
  }
  删除的管理域是一级域
  {
  StatusCode:400,
  MessageType:"string",
  Message:"本域无法处理一级前缀"
  }
  服务器错误
  {
  StatusCode:500,
  MessageType:"string",
  Message:"服务器错误"
  }
```

### 7.minShowInfo 展示信息

- request param:

```
 {
 "Type"    : "network",
 "Command": "ShowInfo",
 }
```

- response param:

```
  {
  StatusCode:200,
  MessageType:"string",
  Message:"管理域信息"
  }
```