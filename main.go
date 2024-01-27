package main

import (
	"context"
	"fmt"

	"github.com/oreoluwa-bs/epoc/web"
)

func main() {
	app := web.New()

	// ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	// defer cancel()

	err := app.Start(context.TODO())
	if err != nil {
		fmt.Println("failed to start app:", err)
	}
}
