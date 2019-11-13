# gredis
golang实现redis客户端

## redis 协议简述

```
*<参数数量> CR LF
$<参数 1 的字节数量> CR LF
<参数 1 的数据> CR LF
...
$<参数 N 的字节数量> CR LF
<参数 N 的数据> CR LF
假设需要set一个val,则协议侧的实现如下:
```

```bash
➜  ~ telnet 127.0.0.1 6379
Trying 127.0.0.1...
Connected to localhost.
Escape character is '^]'.
*3
$3
set
$5
mykey
$2
ok
+OK
```