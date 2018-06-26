# times

golang的time辅助包

1. 写此包的原因：
   Go原生的time包封装程度不够，用起来没那么顺手，而且Format很是奇怪，必须基于2006年1月2日15:04:05；此包是对time包的一些封装
2. 实现方式：参照PHP的date()函数和strtotime()函数实现，使用方式也与这两个函数类似。


# 使用实例

- 字符串转为time.Time类型：

```
t := times.StrToLocalTime("2012-11-12 23:32:01")
t := times.StrToLocalTime("2012-11-12")
```

原生的Go包这么实现：

```
t := time.Date(2012, 11, 12, 23, 32, 01, 0, time.Local)
t := time.Date(2012, 11, 12, 0, 0, 0, 0, time.Local)
```
  
- time.Time类型格式化为字符串：

```
now := time.Now()
strTime := times.Format("Y-m-d H:i:s", now)
```

原生的Go包这么实现：

```
strTime := time.Now().Format("2006-01-02 15:04:05")
```
