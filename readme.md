# 将 go 作为 IDL 来使用

## 背景

idl 文件描述接口，但是本身缺少版本控制和模块控制

## 目的

实现 版本控制和模块控制

## 思路

1. go 语法和 thrift 很像
2. go 需要 generator，thrift 也需要。没有引入新的包袱
3. go 有 go mod，可以做模块引用和版本管理；go mod 能力的上限就是我们 idl 的上限

## 临时实现

时间有限，只做了一个 go 2 thrift 的翻译程序（本意并不是这个），然后再用翻译得到的 thrift 文件去使用 thrift 生态的工具

使用参考本项目下的 `example/idl_test.go` （记得带上 build tag `idl`

## 现在可以

当我们创建不同的 git branch，并且服务于同一个服务接口时，版本管理时不可缺少的

现在我们可以用 git tag 给不同的 branch 打标签，并通过 go mod 根据版本需要去指定依赖


