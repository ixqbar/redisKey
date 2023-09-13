package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func RedisKeyScan(cliContext *cli.Context) error {
	redisCtx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cliContext.String("host"), cliContext.Int("port")),
		Password: cliContext.String("password"),
		DB:       cliContext.Int("db"),
	})

	_, err := rdb.Ping(redisCtx).Result()
	if err != nil {
		return err
	}

	//扫描key
	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = rdb.Scan(redisCtx, cursor, cliContext.Args().Get(0), 10).Result()
		if err != nil {
			return err
		}

		for i := range keys {
			fmt.Printf("%s\n", keys[i])
		}

		if cliContext.Bool("delete") {
			rdb.Del(redisCtx, keys...)
		}

		n += len(keys)
		if cursor == 0 {
			break
		}
	}

	return nil
}

func main() {
	cli.HelpFlag = &cli.BoolFlag{
		Name: "help",
	}

	app := &cli.App{
		Name:  "redisKey",
		Usage: "按规则扫描过滤key",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "host",
				Value:   "127.0.0.1",
				Aliases: []string{"h"},
			},
			&cli.IntFlag{
				Name:    "port",
				Value:   6379,
				Aliases: []string{"p"},
			},
			&cli.IntFlag{
				Name:    "db",
				Value:   0,
				Aliases: []string{"d"},
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "",
			},
			&cli.BoolFlag{
				Name:     "delete",
				Value:    false,
				Required: false,
			},
		},
	}

	app.Action = func(cliContext *cli.Context) error {
		if cliContext.Args().Len() != 1 || len(cliContext.Args().Get(0)) == 0 {
			return cli.ShowAppHelp(cliContext)
		}
		return RedisKeyScan(cliContext)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
