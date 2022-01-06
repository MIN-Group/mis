### 1.UpLoadFace 上传人脸
 - prefix:/Face(without encrypt) /face(encrypted)
 - request param:


```JSON
 {
 "Type"        : "user-act",
 "Command"     : "UpLoadFace",
 "Id"          : "string(username)",
 "Window"      : "int(the number of slices)",
 "Seq"         : "int(sequence number)",
 "Data"        : "string(face slice data)"
 }
```


- response param:

上传成功:
```JSON
 {
 "StatusCode":200,
 "MessageType":"string",
 "Message":"上传成功"
 }
```

### 2.GetFace 获取人脸

 - prefix:/Face(without encrypt) /face(encrypted)
 - request param:


```JSON
 {
 "Type"        : "user-act",
 "Command"     : "GetFace",
 "Username"    : "string(username)",
 "Seq"         : "int(sequence number, init=1)"
 }
```


- response param:

获取成功:
```Json
 {
 "StatusCode":200,
 "MessageType":"json",
 "Message":{"Seq":3,"Data":"xxxx","Window":4}
 }
```

用户不存在:
```Json
 {
 "StatusCode":400,
 "MessageType":"string",
 "Message":"用户不存在"
 }
```

序列号大于窗口长度:
```Json
 {
 "StatusCode":401,
 "MessageType":"string",
 "Message":"序列号大于窗口长度"
 }
```