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

func (l *lmClient) PutLog(prefix, data string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logReply, err := l.client.PutLog(ctx, &pb.PutLogRequest{
		Prefix: prefix,
		Data:   data,
	}); err != nil {
		log.Fatal(handleCallError("PutLog", err))
	} else {
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) PutLogStream(prefix string, n int) {
	ctx := context.Background()

	stream, err := l.client.PutLogStream(ctx)
	if err != nil {
		log.Fatal(handleCallError("PutLogStream", err))
		return
	}

	go func() {
		for i := 0; i < n; i++ {
			if err := stream.Send(&pb.PutLogRequest{
				Prefix: prefix,
				Data:   randomString(),
			}); err != nil {
				log.Fatal(handleCallError("PutLogStream.Send", err))
				break
			}
		}
		if err := stream.CloseSend(); err != nil {
			log.Fatal(handleCallError("PutLogStream.CloseSend", err))
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			return
		default:
		}

		logReply, recvErr := stream.Recv()
		if recvErr == io.EOF {
			return
		}

		if recvErr != nil {
			log.Fatal(handleCallError("PutLogStream.Recv", recvErr))
		}

		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) GetLog(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logReply, err := l.client.GetLog(ctx, &pb.KeyMessage{
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

	if logsReply, err := l.client.GetLogs(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	}); err != nil {
		log.Fatal(handleCallError("GetLogs", err))
	} else {
		log.Printf("%+v\n", logsReply)
	}
}

func (l *lmClient) GetLogsStream(prefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	stream, err := l.client.GetLogsStream(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatal(handleCallError("GetLogsStream", err))
		return
	}
	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(handleCallError("GetLogsStream.Recv", err))
			return
		}
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) TailLogStream(prefix string) {
	ctx := context.Background()

	stream, err := l.client.TailLogStream(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatal(handleCallError("TailLogStream", err))
		return
	}
	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(handleCallError("TailLogStream.Recv", err))
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

func (l *lmClient) ListKeys(prefix string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if logsReply, err := l.client.ListKeys(ctx, &pb.PrefixRequest{
		Prefix: prefix,
	}); err != nil {
		log.Fatal(handleCallError("ListKeys", err))
	} else {
		log.Printf("%+v\n", logsReply)
	}
}

func (l *lmClient) PutLogsUnary(prefix string, n int) {
	ctx := context.Background()

	funcs := 8
	funcsDone := 0
	doneC := make(chan bool)

	for i := 0; i < funcs; i++ {
		go floodLogsUnary(i, l.client, ctx, doneC, prefix, n/8)
	}

	for {
		select {
		case <-doneC:
			funcsDone = funcsDone + 1
			if funcsDone == funcs {
				os.Exit(0)
			}
		}
	}
}

func floodLogsUnary(funcId int, client pb.LoggerClient, ctx context.Context, done chan bool, prefix string, n int) {
	for i := 0; i < n; i++ {
		if res, err := client.PutLog(ctx, &pb.PutLogRequest{
			Prefix: prefix,
			Data:   randomString(),
		}); err != nil {
			handleCallError("PutLog", err)
		} else {
			log.Printf("funcId:%d %+v\n", funcId, res)
		}
	}
	done <- true
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
