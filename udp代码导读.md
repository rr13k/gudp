reliability.go 函数解析

NewServer：用于创建一个新的服务器实例，初始化一些默认选项，如传输编码器、加密方法等。

SetHandler：用于设置当从客户端接收到自定义记录类型时调用的处理函数。

Serve：用于启动服务器，监听UDP端口并处理接收到的数据。

Stop：用于停止服务器。

handleRecord：用于处理接收到的记录，根据记录类型进行不同的处理，如握手、心跳等。

handleHandshakeRecord：处理握手记录，包括客户端和服务器之间的握手协议。

handlePingRecord：处理ping记录并发送pong响应。

handleCustomRecord：处理自定义记录类型，并调用设置的处理函数来处理数据。

registerClient：注册客户端，生成会话ID并保存客户端相关信息。

findClientBySessionID、findClientIDBySessionID、findClientByAddr：根据不同的标识查找客户端。

sendToAddr：向指定地址发送数据。

sendToClient：向客户端发送加密后的数据。

SendToClientByID：根据客户端ID向特定客户端发送数据。

clientGarbageCollection：定期清理过期的客户端记录。

BroadcastToClients：向所有已注册的客户端广播数据。

parseRecord：解析接收到的字节为记录结构。

parseSessionID：从解密后的记录体中解析会话ID。