# Ice for Go

自动生成golang代码，只生成客户端sdk。
注意：对参数有严格要求，只支持如下格式。bytes可以用json或者protobuf。

Automatic generation of golang code, client-side sdk only.
Note: There are strict requirements for parameters and only the following formats are supported. bytes can be json or protobuf.

```
sequence<byte> bytes;
void Echo(bytes req, out bytes res);
```


## Demo

### server

```
cd demo/server
make
./server
```

### client

```
cd demo
../slice2go ./Printer.ice
go run client.go
```
