package main

import (
	"bytes"
	"context"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/yunomu/domprinter"
	"github.com/yunomu/domprinter/lambda/events"
)

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(logger)
}

type handler struct {
	domprinter *domprinter.DomPrinter
}

func (h *handler) handle(ctx context.Context, req *events.Request) ([]byte, error) {
	zap.L().Info("Handler", zap.Any("request", req))

	var buf bytes.Buffer
	if err := h.domprinter.Print(ctx, req.Url, &buf); err != nil {
		zap.L().Error("handler/Print", zap.Error(err), zap.Any("request", req))
		return nil, err
	}

	return buf.Bytes(), nil
}

func main() {
	ctx := context.Background()

	h := &handler{
		domprinter: domprinter.New(),
	}
	lambda.StartWithContext(ctx, h.handle)
}
