
## NeeRPC 报文格式
```text
| Option{MagicNumber: xxx, CodecType: xxx} | Header{ServiceMethod ...} | Body interface{} |
| <------      固定 JSON 编码      ------>  | <-------   编码方式由 CodeType 决定   ------->|
```
- `Option`固定在报文的最开始, `Header` 和 `Body`可以有多个

```text
| Option | Header1 | Body1 | Header2 | Body2 | ...
```

## 启动 `NeeRPC` 服务
```go
lis, _ := net.Listen("tcp", ":9999")
geerpc.Accept(lis)
```

## Call 的设计
对`net/rpc` 而言，一个函数需要能够被远程调用，需要满足如下五个条件
- the method’s type is exported. 
- the method is exported. 
- the method has two arguments, both exported (or builtin) types. 
- the method’s second argument is a pointer. 
- the method has return type error.

example

```go

func (t *T) MethodName(argType T1, replyType *T2) error
```


## 超时处理
### 客户端处理超时的地方
- 与服务端建立连接，导致的超时
- 发送请求到服务端，写报文导致的超时
- 等待服务端处理时，等待处理导致的超时（比如服务端已挂死， 迟迟不响应）
- 从服务端接收响应时，读报文导致的超时

### 服务端处理超时的地方
- 读取客户端请求报文时，读报文导致的超时
- 发送响应碑文时，写报文导致的超时
- 掉哟过映射服务的方法时，处理报文导致的超时

### 超时处理机制
1. 客户端创建连接时
2. 客户端 `Client.Call()`整个过程导致的超时（包含发送报文，等待处理，接收报文所有截断）
3. 服务端处理报文，即`Server.HandleRequest`超时

### 用户可以使用`context.WithTimeout`创建具备超时检测能力的`context`对象来控制
```go
ctx, _ := context.WithTimeout(context.Background(), time.Secnoud)
var reply int
err := client.Call(ctx, "Foo.Sum", &Args{1,2}, &reply)
...
```

## 服务端支持`http`协议
- 客户端向`rpc`服务器发送`CONNECT` 请求
```http request
CONNECT x.x.x.x:xxxx/_neerpc_ http/1.0
```
- `RPC`服务器返回 `HTTP 200` 状态码表示连接建立
```http request
HTTP/1.0 200 Connected to Nee RPC
```
- 客户端使用创建好的连接发送`RPC`报文，先发送`Option`，再发送`N`个请求报文，服务端处理`RPC`请求并响应


## 负载均衡
- 随机选择
- `Round Robin`轮询调度算法

## 服务发现和注册中心
客户端和服务端都只需要感知注册中心的存在，而无需感知对方的存在
1. 服务端启动后，向注册中心发送注册信息，注册中心得知该服务已经启动，处于可用状态。一般来说，服务端还需要定期向注册中心发送心跳，证明自己还活着。
2. 客户端向注册中心询问，当前那个服务是可用的，注册中心将可用的服务列表返回给客户端
3. 客户端根据注册中心得到的服务列表，选择其中一个发起调用

注册中心的功能，比如：

- 配置的动态同步
- 通知机制等

