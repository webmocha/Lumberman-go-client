<p align="center">
  <img src="https://user-images.githubusercontent.com/132562/63731122-1730de80-c823-11e9-833f-3e4c91670a46.png" alt="Lumberman go client" />
</p>

<h1 align="center">Lumberman go client</h1>

<p align="center">
  <strong><a href="https://github.com/webmocha/Lumberman">Lumberman</a> client reference implementation for go</strong>
</p>

## Requirements

```sh
go get
```

## Lumberman LogService

see [lumberman.proto](https://github.com/webmocha/Lumberman/blob/master/lumberman.proto)

## Installation

```sh
go get -u github.com/webmocha/lumberman-go-client
ln -s $GOPATH/bin/lumberman-go-client $GOPATH/bin/lmc
```

## Usage

### Define server address

_default is `127.0.0.1:9090`_

override with `-server_addr` flag

example:

```sh
lmc -server_addr=172.17.0.14:5000 list-prefixes
```

### Write to Log

```sh
lmc <Prefix> <Log Object>
```

examples:

```sh
lmc user-search 'cat'
```

```sh
lmc user-click '{ "href": "/login" }'
```

```sh
lmc player-move '{ "x": 20, "y": -42, "z": 1 }'
```

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

### Stream Logs by prefix

_stream stays open, tailing new log events_

```sh
lmc stream-logs <Prefix>
```

example:

```sh
lmc stream-logs user-click
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
lmc list-logs <Prefix>
```

example:

```sh
lmc list-logs user-click
```

output:

```
2019/08/26 16:26:15 keys:"user-click|2019-08-26T06:15:37.24192515Z" keys:"user-click|2019-08-26T06:19:00.062988065Z" keys:"user-click|2019-08-26T06:19:02.662282619Z"
```

### Flood logs for prefix

floods log for 10 minutes

```sh
lmc log-flood <Prefix>
```
