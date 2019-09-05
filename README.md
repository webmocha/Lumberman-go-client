<p align="center">
  <img src="https://user-images.githubusercontent.com/132562/63731122-1730de80-c823-11e9-833f-3e4c91670a46.png" alt="Lumberman go client" />
</p>

<h1 align="center">Lumberman go client</h1>

<p align="center">
  <strong><a href="https://github.com/webmocha/Lumberman">Lumberman</a> client reference implementation for go</strong>
</p>

## Reference

This repo contains runnable examples for connecting to [Lumberman](https://github.com/webmocha/Lumberman) grpc server and making calls to each function.

All of the examples are in [client.go](./client.go).

## Lumberman Service Definition

see [lumberman.proto](https://github.com/webmocha/Lumberman/blob/master/lumberman.proto)


## Quick Guide

### Create a grpc client

```go
conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
if err != nil {
  log.Fatalf("fail to dial: %v", err)
}
defer conn.Close()
client := pb.NewLoggerClient(conn)

lmc := &lmClient{
  client: client,
}
```

### List keys for Prefix

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

logsReply, err := lmc.ListKeys(ctx, &pb.PrefixRequest{
  Prefix: prefix,
})
if err != nil {
	if s, ok := status.FromError(err); !ok {
		log.Fatal(status.Errorf(codes.Internal, "client.ListKeys <- server Unknown Internal Error('%s')", s.Message()))
	} else {
		log.Fatal(status.Errorf(s.Code(), "client.ListKeys<-server.Error('%s')", s.Message()))
	}
}

log.Printf("%+v\n", logsReply)
```

### Get all logs from a stream

```javascript
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

stream, err := l.client.GetLogsStream(ctx, &pb.PrefixRequest{
  Prefix: prefix,
})
if err != nil {
  if s, ok := status.FromError(err); !ok {
    log.Fatal(status.Errorf(codes.Internal, "client.GetLogsStream <- server Unknown Internal Error('%s')", s.Message()))
  } else {
    log.Fatal(status.Errorf(s.Code(), "client.GetLogsStream<-server.Error('%s')", s.Message()))
  }
}

for {
  logReply, err := stream.Recv()
  if err == io.EOF {
    break
  }
  if err != nil {
    if s, ok := status.FromError(err); !ok {
      log.Fatal(status.Errorf(codes.Internal, "client.GetLogsStream.Recv <- server Unknown Internal Error('%s')", s.Message()))
    } else {
      log.Fatal(status.Errorf(s.Code(), "client.GetLogsStream.Recv<-server.Error('%s')", s.Message()))
    }
  }

  log.Printf("%+v\n", logReply)
}
```

## Usage

### Installation

```sh
go get -u github.com/webmocha/lumberman-go-client
ln -s $GOPATH/bin/lumberman-go-client $GOPATH/bin/lmc
```

### Options

Define server address

_default is `127.0.0.1:9090`_

override with `-server_addr` flag

example:

```sh
lmc -server_addr=172.17.0.14:5000 list-prefixes
```

### Write to Log

```sh
lmc put-log <Prefix> <Log Object>
```

examples:

```sh
lmc put-log user-search 'cat'
```

```sh
lmc put-log user-click '{ "href": "/login" }'
```

```sh
lmc put-log player-move '{ "x": 20, "y": -42, "z": 1 }'
```

returns: `key`

### Get Log by key

```sh
lmc get-log <Key Name>
```

example:

```sh
lmc get-log 'user-click|2019-08-26T06:19:02.662282619Z'
```

output:

```
2019/08/26 16:19:27 key:"user-click|2019-08-26T06:19:02.662282619Z" timestamp:<... > data:"{ \"href\": \"/login\" }"
```

### Get all Logs by prefix

```sh
lmc get-logs <Prefix>
```

example:

```sh
lmc get-logs user-search
```

output:

```
2019/08/26 16:20:49 logs:<key:"user-search|2019-08-26T01:30:42.620978567Z" timestamp:<... > data:"cat" > logs:<key:"user-search|2019-08-26T01:31:38.844208133Z" timestamp:<... > data:"doggo" > logs:<key:"user-search|2019-08-26T01:31:42.385940486Z" timestamp:<... > data:"birb" >
```

### Get all Logs as stream by prefix

```sh
lmc get-logs-stream <Prefix>
```

example:

```sh
lmc get-logs-stream user-search
```

output:

```
2019/08/26 16:20:49 logs:<key:"user-search|2019-08-26T01:30:42.620978567Z" timestamp:<... > data:"cat" >
2019/08/26 16:20:49 logs:<key:"user-search|2019-08-26T01:31:38.844208133Z" timestamp:<... > data:"doggo" >
l2019/08/26 16:20:49 ogs:<key:"user-search|2019-08-26T01:31:42.385940486Z" timestamp:<... > data:"birb" >
```

### Tail Logs as stream by prefix

_stream stays open, tailing new log events_

```sh
lmc tail-logs-stream <Prefix>
```

example:

```sh
lmc tail-logs-stream user-click
```

output:

```
2019/08/26 16:23:08 key:"user-click|2019-08-26T06:19:00.062988065Z" timestamp:<... > data:"{ \"href\": \"/login\" }"
2019/08/26 16:23:10 key:"user-click|2019-08-26T06:19:02.662282619Z" timestamp:<... > data:"{ \"href\": \"/forgot-password\" }"
```

### List Log prefixes

```sh
lmc list-prefixes
```

example output:

```
2019/08/26 16:25:03 prefixes:"user-search" prefixes:"user-click" prefixes:"player-move"
```

### List Log keys by prefix

```sh
lmc list-keys <Prefix>
```

example:

```sh
lmc list-keys user-click
```

output:

```
2019/08/26 16:26:15 keys:"user-click|2019-08-26T06:15:37.24192515Z" keys:"user-click|2019-08-26T06:19:00.062988065Z" keys:"user-click|2019-08-26T06:19:02.662282619Z"
```

### Put n logs by calling n unary rpc calls

Default `n: 1000`

```sh
lmc put-logs-unary <Prefix>
lmc -n <CallsCount> put-logs-unary <Prefix>
```

example:

```sh
lmc put-logs-unary 'test-1000'
lmc -n 2500 put-logs-unary 'test-2500'
```

### Put n logs using single bidirectional rpc stream

Default `n: 1000`

```sh
lmc put-logs-stream <Prefix> (n:1000)
lmc -n <CallsCount> put-logs-stream <Prefix> (n:1000)
```

example:

```sh
lmc put-logs-stream 'test-1000'
lmc -n 2500 put-logs-stream 'test-2500'
```
