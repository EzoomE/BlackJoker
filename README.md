# BlackJoker  
C2服务器  
该服务器当前提供了基础的反向Shell模块,基于HTTP通信,可以实现僵尸网络的基础搭建和布局    
———————————————————————————————————————————  
Go的交叉编译可以使服务端运行在任何平台,如Linux  
客户端由Python编写,在Windows运行  
  
使用说明:  
1.客户端的监听IP需要自己设置,在BlackJoker/Client/main.pyx文件的开头ip变量
