## 构建c++编译环境

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

## 开始编译:

```sh
# 进入编译位置
mkdir -p gns/lib/GameNetworkingSockets/build && cd gns/lib/GameNetworkingSockets/build

# 构建Ninja生成器配置文件, 其中DCMAKE_INSTALL_PREFIX 用来指定安装路径
cmake -G Ninja -DCMAKE_INSTALL_PREFIX=/Users/zhouyuan11/work/pulseHero-udp/gns/lib/GameNetworkingSockets/build ..

# 使用ninja进行编译
ninja

# 编译后进行安装
sudo ninja install
```


# udp 要保证心跳
消息要区分为
Client : 1
loginVerify: 2
Server: 3
Ping: 4
Pong: 5