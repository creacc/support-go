package session

// 通用
var (
	NoError = Status{0, "正常", nil}
	Unknown = Status{-1, "未知", nil}

	NotFound  = StatusFactory{code: 1}
	FormatErr = StatusFactory{code: 2}

	/*************************************************模块*****************************************************/
	BusinessModule   = root.shift(1, "业务模块", 100)
	foundationModule = root.shift(9, "基础模块", 100)
	/*************************************************模块*****************************************************/

	/*************************************************组件*****************************************************/
	httpSrvComponent = foundationModule.shift(1, "http服务端组件", 100)
	httpCliComponent = foundationModule.shift(2, "http客户端组件", 100)
	mySQLComponent   = foundationModule.shift(3, "MySQL组件", 100)
	redisComponent   = foundationModule.shift(4, "Redis组件", 100)
	ymlComponent     = foundationModule.shift(5, "YAML组件", 100)
	jsonComponent    = foundationModule.shift(6, "Json组件", 100)
	/*************************************************组件*****************************************************/
)

// MySQL数据库
var (
	SQLSyntaxErr = mySQLComponent.factory("SQL语法错误")
)

// Redis
var (
	RedisDialErr = redisComponent.factory("Redis链接错误")
	RedisSetErr  = redisComponent.factory("Redis写入错误")
	RedisGetErr  = redisComponent.factory("Redis读取错误")
)

// http server
var (
	ServerReqErr                 = httpSrvComponent.factory("服务端请求错误")
	ServerReqLackErr             = httpSrvComponent.factory("服务端请求缺少参数")
	ServerReqBuildErr            = httpSrvComponent.factory("服务端请求构建错误")
	ServerReqExecuteErr          = httpSrvComponent.factory("服务端请求执行错误")
	ServerReqParamTypeUnknownErr = httpSrvComponent.factory("服务端未知类型参数")
	ServerReqQueryParamErr       = httpSrvComponent.factory("Query参数错误")
	ServerReqFormParamErr        = httpSrvComponent.factory("表单参数错误")
	ServerReqPathParamErr        = httpSrvComponent.factory("路径参数错误")
	ServerReqFileOpenErr         = httpSrvComponent.factory("文件打开错误")
	ServerReqFileCopyErr         = httpSrvComponent.factory("文件复制错误")
	ServerReqFileEmptyParamErr   = httpSrvComponent.factory("空文件错误")
	ServerRespError              = httpSrvComponent.factory("服务端响应错误")
)

var (
	ClientReqNewErr = httpCliComponent.factory("客户端请求构建错误")
	ClientReqDoErr  = httpCliComponent.factory("客户端请求执行错误")

	ClientReqJsonParamErr           = httpCliComponent.factory("客户端请求Json参数错误")
	ClientReqMultipartParamFieldErr = httpCliComponent.factory("客户端请求Multipart参数错误")
	ClientReqFileParamFileOpenErr   = httpCliComponent.factory("客户端请求文件参数打开错误")
	ClientReqFileParamFileCreateErr = httpCliComponent.factory("客户端请求文件参数创建错误")
	ClientReqFileParamFileCopyErr   = httpCliComponent.factory("客户端请求文件参数复制错误")
	ClientReqFileParamFileEmptyErr  = httpCliComponent.factory("客户端请求文件参数为空")
)

var (
	YMLReadFileErr = httpCliComponent.factory("客户端请求执行错误")
)

var (
	JsonArrayKeyTypeErr   = httpCliComponent.factory("客户端请求执行错误")
	JsonArrayValueTypeErr = httpCliComponent.factory("客户端请求执行错误")
)
