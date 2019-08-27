package main

import (
	"flag"
	"fmt"
	"log"

	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc"
)

const help = `Usage

Write to Log:
		lmc log <Prefix> <Log Object>

Get Log by key:
		lmc get-log <Key>

Get all Logs by prefix:
		lmc get-logs <Prefix>

Get all Logs as stream by prefix:
		lmc get-logs-stream <Prefix>

Stream Logs by prefix:
		lmc stream-logs <Prefix>

List Log by prefixes:
		lmc list-prefixes

List Log keys by prefix:
		lmc list-logs <Prefix>

Flood logs for prefix:
		lmc log-flood <Prefix>
`

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:9090", "The server address in the format of host:port")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewLoggerClient(conn)

	lmc := &lmClient{
		client: client,
	}

	flag.Args()
	switch len(flag.Args()) {
	case 1:
		switch flag.Args()[0] {
		case "list-prefixes":
			lmc.ListPrefixes()
		default:
			fmt.Println(help)
		}

	case 2:
		switch flag.Args()[0] {
		case "get-log":
			lmc.GetLog(flag.Args()[1])
		case "get-logs":
			lmc.GetLogs(flag.Args()[1])
		case "get-logs-stream":
			lmc.GetLogsStream(flag.Args()[1])
		case "stream-logs":
			lmc.StreamLogs(flag.Args()[1])
		case "list-logs":
			lmc.ListLogs(flag.Args()[1])
		case "log-flood":
			lmc.LogFlood(flag.Args()[1])
		default:
			fmt.Println(help)
		}

	case 3:
		switch flag.Args()[0] {
		case "log":
			lmc.Log(flag.Args()[1], flag.Args()[2])
		default:
			fmt.Println(help)
		}
	default:
		fmt.Println(help)
	}
}
