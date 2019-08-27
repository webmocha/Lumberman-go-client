package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	empty "github.com/golang/protobuf/ptypes/empty"
	pb "github.com/webmocha/lumberman/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type lmClient struct {
	client pb.LoggerClient
}

func (l *lmClient) Log(prefix, data string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logReply, err := l.client.Log(ctx, &pb.LogRequest{
		Prefix: prefix,
		Data:   data,
	}); err != nil {
		log.Fatal(handleCallError("Log", err))
	} else {
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) GetLog(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logReply, err := l.client.GetLog(ctx, &pb.GetLogRequest{
		Key: key,
	}); err != nil {
		log.Fatal(handleCallError("GetLog", err))
	} else {
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) GetLogs(prefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logsReply, err := l.client.GetLogs(ctx, &pb.GetLogsRequest{
		Prefix: prefix,
	}); err != nil {
		log.Fatal(handleCallError("GetLogs", err))
	} else {
		log.Printf("%+v\n", logsReply)
	}
}

func (l *lmClient) StreamLogs(prefix string) {
	ctx := context.Background()

	stream, err := l.client.StreamLogs(ctx, &pb.GetLogsRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatal(handleCallError("StreamLogs", err))
		return
	}
	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(handleCallError("StreamLogs.Recv", err))
			return
		}
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) ListPrefixes() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if prefixesReply, err := l.client.ListPrefixes(ctx, new(empty.Empty)); err != nil {
		log.Fatal(handleCallError("ListPrefixes", err))
	} else {
		log.Printf("%+v\n", prefixesReply)
	}
}

func (l *lmClient) ListLogs(prefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logsReply, err := l.client.ListLogs(ctx, &pb.ListLogsRequest{
		Prefix: prefix,
	}); err != nil {
		log.Fatal(handleCallError("ListLogs", err))
	} else {
		log.Printf("%+v\n", logsReply)
	}
}

func (l *lmClient) LogFlood(prefix string) {
	ctx := context.Background()

	for i := 0; i < 8; i++ {
		go floodLogs(i, l.client, ctx, prefix)
	}

	time.Sleep(10 * time.Minute)
	log.Println("DONE")
	os.Exit(0)
}

func floodLogs(funcId int, client pb.LoggerClient, ctx context.Context, prefix string) {
	for {
		if res, err := client.Log(ctx, &pb.LogRequest{
			Prefix: prefix,
			Data:   randomString(),
		}); err != nil {
			handleCallError("Log", err)
		} else {
			log.Printf("funcId:%d %+v\n", res)
		}
	}
}

func randomString() string {
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(97 + rand.Intn(122-97))
	}
	return string(bytes)
}

func handleCallError(rpcFunc string, err error) error {
	if s, ok := status.FromError(err); !ok {
		return status.Errorf(codes.Internal, "client.%s <- server Unknown Internal Error('%s')", rpcFunc, s.Message())
	} else {
		return status.Errorf(s.Code(), "client.%s<-server.Error('%s')", rpcFunc, s.Message())
	}
}
