# BlackJokerC2  
该服务器当前提供了基础的反向Shell模块,基于HTTP通信,可以实现僵尸网络的基础搭建和布局  
———————————————————————————————————————————  
Go的交叉编译可以使服务端运行在任何平台,如Linux  
客户端由Python编写,在Windows运行  
  
使用说明:  
1.客户端的监听IP需要自己设置,在BlackJoker/Client/main.pyx文件的开头ip变量  
2.如需编译go为二进制文件则在Blackarch下执行`go build -o`  
3.该项目使用Mysql存储Session信息数据,需要首先手动在Blackarch/Mysql/jsonOSsystemDB.go中配置数据库信息,!!!以及首先创建一个名为jsonossystem的数据库!!!  
  
环境要求:  
1.golang语言环境  
2.Python语言环境(方便Client DeBug)  
3.Server可编译为任意平台可执行文件  

