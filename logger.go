package shorty

import (
	"log/slog"
	"os"
)

func DefaultLogger() *slog.Logger {
	return slog.New(
		slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}),
	)
}
