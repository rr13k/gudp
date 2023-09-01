## gudp

## 简介

gudp 为godot game 所实现，提供轻量快速的可靠udp实现, 用户通过简单的函数即可快速搭建可靠的udp服务。

使用mit协议开源，支持通过fork进行贡献。

## 特性
1. 支持可靠消息、不可靠消息
2. 支持身份认证
3. 内置支持重发和心跳检测
4. 使用proto作为通用协议

## 使用方式
```sh
# 检查cpp环境是否安装
g++ --version

# 检查是否安装了 openssl
openssl version

# mac安装openssl
brew install openssl

# 安装构建需要的ninja
brew install ninja

```

