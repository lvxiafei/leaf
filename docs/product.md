用户角度：

1. 定位：一款抹平底层网络，提供上层网络设计、自检的奢侈品工具
2. 回答：数据包 到哪里去？走什么路？走到哪了？

产品功能：

1. 支持五元组+网口、源目mac、路由、vlan id等，再也不用argue到哪去，下一站，从哪走这些事了
    - 通道中的唯一标志：五元组
    - 端点处的唯一标志：十元组，mac、路由、网口
2. 支持netfilter问题精准定位，可以清楚知道那个表，那个链，那条规则的问题，再也不用一个个查了
    1. 支持kube-proxy、自定义链等
    2. 支持pod
    3. 更加轻量级，不涉及内核模块insmod
3. 支持自动画真实流量图，定义设计图，diff流量图，在网络设计和自检时较为奢侈
4. 支持多场景，如node2eip、pod2pod、pod2eip等
5. 支持隧道、vlan等常见流量
6. 支持展示网络包在内核协议栈各函数各层耗时
7. 支持自定义hook函数，更加轻量
8. 支持ringbuf，保证输出顺序
9. 支持进程名
10. ...