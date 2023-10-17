package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/Yanhao/bedrock-cli/clients/proxy"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
)

func kvset(ctx *cli.Context) error {
	key := ctx.String("key")
	value := ctx.String("value")

	resp, err := getProxyConn().KvSet(context.TODO(), &proxy.KvSetRequest{
		StorageId: 0,
		Key:       []byte(key),
		Value:     []byte(value),
	})

	if err != nil || resp.Err != proxy.Error_OK {
		log.Printf("kvset failed, err: %v", err)
		return err
	}

	return nil
}

func kvget(ctx *cli.Context) error {
	key := ctx.String("key")

	resp, err := getProxyConn().KvGet(context.TODO(), &proxy.KvGetRequest{
		StorageId: 0,
		Key:       []byte(key),
	})

	if err != nil || resp.Err != proxy.Error_OK {
		log.Printf("kvget failed, err: %v", err)
		return err
	}

	log.Printf("value: %s", string(resp.Value))

	return nil
}

func kvdel(ctx *cli.Context) error {
	key := ctx.String("key")

	resp, err := getProxyConn().KvDelete(context.TODO(), &proxy.KvDeleteRequest{
		StorageId: 0,
		Key:       []byte(key),
	})

	if err != nil || resp.Err != proxy.Error_OK {
		log.Printf("kvget failed, err: %v", err)
		return err
	}

	log.Println("delete success")

	return nil
}

var (
	proxyCli     proxy.ProxyServiceClient
	proxyCliOnce sync.Once
)

func getProxyConn() proxy.ProxyServiceClient {
	proxyCliOnce.Do(func() {
		addr := os.Getenv("PROXY_URL")
		if addr == "" {
			panic("cannot find env `PROXY_URL`")
		}

		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			panic(err)
		}

		proxyCli = proxy.NewProxyServiceClient(conn)

	})

	return proxyCli

}

func main() {
	app := &cli.App{

		Commands: []*cli.Command{
			{
				Name:   "kvset",
				Action: kvset,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "key",
					},
					&cli.StringFlag{
						Name: "value",
					},
				},
			},
			{
				Name:   "kvget",
				Action: kvget,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "key",
					},
				},
			},
			{
				Name:   "kvdel",
				Action: kvdel,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name: "key",
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
