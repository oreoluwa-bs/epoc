package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/oreoluwa-bs/epoc/web"
)

func main() {
	app := web.New(web.LoadConfig())

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	err := app.Start(ctx)
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
