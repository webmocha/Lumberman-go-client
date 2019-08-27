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
		lmc put-log <Prefix> <Log Data>

Get Log by key:
		lmc get-log <Key>

Get all Logs by prefix:
		lmc get-logs <Prefix>

Get all Logs as stream by prefix:
		lmc get-logs-stream <Prefix>

Tail Logs as stream by prefix:
		lmc tail-logs-stream <Prefix>

List Log by prefixes:
		lmc list-prefixes

List Log keys by prefix:
		lmc list-keys <Prefix>

Flood logs for prefix:
		lmc log-flood <Prefix>

Flood logs for prefix as stream:
		lmc put-log-stream <Prefix>
		default n: 1000
		lmc -n 2500 put-log-stream <Prefix>
`

var (
	serverAddr = flag.String("server_addr", "127.0.0.1:9090", "The server address in the format of host:port")
	putStreamN = flag.Int("n", 1000, "number of logs to put (default: 1000)")
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
		case "tail-logs-stream":
			lmc.TailLogStream(flag.Args()[1])
		case "list-keys":
			lmc.ListKeys(flag.Args()[1])
		case "log-flood":
			lmc.LogFlood(flag.Args()[1])
		case "put-logs-stream":
			lmc.PutLogStream(flag.Args()[1], *putStreamN)
		default:
			fmt.Println(help)
		}

	case 3:
		switch flag.Args()[0] {
		case "put-log":
			lmc.PutLog(flag.Args()[1], flag.Args()[2])
		default:
			fmt.Println(help)
		}
	default:
		fmt.Println(help)
	}
}
