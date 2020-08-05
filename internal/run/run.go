package run

import (
	"context"
	"os"
)

func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	return 0
}
