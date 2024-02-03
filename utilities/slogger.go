package utilities

import (
	"log/slog"
	"os"
)

var Slogger *slog.Logger

func InitSlog() {
	Slogger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(Slogger)
}
