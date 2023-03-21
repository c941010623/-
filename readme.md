# Skywatch 作業

使用 GO 並且將 JSON 轉換成 MessagePack 格式

## Description

撰寫 Json to MessagePack 格式程式，並含有 unitTest，
MessagePack to Json 尚未完成。

## Getting Started

### Dependencies

* 使用簡易 JSON 來達到轉換程式的設計
```
{"name":"Alice","age":20,"score":[80,85,90]}
```

### Run

在目錄下執行以下指令，可得到轉換結果，並列印出 MessagePack 格式結果。
```
go run app.go
```

## Unit Test

在目錄下執行以下指令：
```
go test -v 
```

## 改進的地方

目前已知 Json 的順序會導致編碼不同，例如程式中的:
```
{"name":"Alice","age":20,"score":[80,85,90]}
```
在運作中有可能會變成:
```
{"age":20,"score":[80,85,90],"name":"Alice"}
```
故在此 UnixTest 會有時正常，有時失敗