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
)

type lmClient struct {
	client pb.LoggerClient
}

func (l *lmClient) Log(prefix, data string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logReply, err := l.client.Log(ctx, &pb.LogRequest{
		Prefix: prefix,
		Data:   data,
	})
	if err != nil {
		log.Fatalf("%v.Log(_) = _, %v: ", l.client, err)
		return err
	}

	log.Printf("%+v\n", logReply)
	return nil
}

func (l *lmClient) GetLog(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logReply, err := l.client.GetLog(ctx, &pb.GetLogRequest{
		Key: key,
	})
	if err != nil {
		log.Fatalf("%v.GetLog(_) = _, %v: ", l.client, err)
		return err
	}

	log.Printf("%+v\n", logReply)
	return nil
}

func (l *lmClient) GetLogs(prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logsReply, err := l.client.GetLogs(ctx, &pb.GetLogsRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatalf("%v.GetLogs(_) = _, %v: ", l.client, err)
		return err
	}

	log.Printf("%+v\n", logsReply)
	return nil
}

func (l *lmClient) StreamLogs(prefix string) {
	ctx := context.Background()

	stream, err := l.client.StreamLogs(ctx, &pb.GetLogsRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatalf("%v.StreamLogs(_) = _, %v: ", l.client, err)
		return
	}
	for {
		logReply, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.StreamLogs(_) = _, %v", l.client, err)
		}
		log.Printf("%+v\n", logReply)
	}
}

func (l *lmClient) ListPrefixes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	prefixesReply, err := l.client.ListPrefixes(ctx, new(empty.Empty))
	if err != nil {
		log.Fatalf("%v.ListPrefixes(_) = _, %v: ", l.client, err)
		return err
	}

	log.Printf("%+v\n", prefixesReply)
	return nil
}

func (l *lmClient) ListLogs(prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	logsReply, err := l.client.ListLogs(ctx, &pb.ListLogsRequest{
		Prefix: prefix,
	})
	if err != nil {
		log.Fatalf("%v.ListLogs(_) = _, %v: ", l.client, err)
		return err
	}

	log.Printf("%+v\n", logsReply)
	return nil
}

func (l *lmClient) LogFlood(prefix string) {
	ctx := context.Background()

	go func() {
		for {
			if _, err := l.client.Log(ctx, &pb.LogRequest{
				Prefix: prefix,
				Data:   randomString(),
			}); err != nil {
				log.Fatalf("%v.ListLogs(_) = _, %v: ", l.client, err)
				return
			}
		}
	}()
	time.Sleep(10 * time.Minute)
	log.Println("DONE")
	os.Exit(0)
}

func randomString() string {
	bytes := make([]byte, 10)
	for i := 0; i < 10; i++ {
		bytes[i] = byte(97 + rand.Intn(122-97))
	}
	return string(bytes)
}
