package tasks

import (
	"context"

	"go.uber.org/zap"
)

type Foo struct {
	logger *zap.Logger
	msg    string
}

func NewFoo(logger *zap.Logger, msg string) *Foo {
	return &Foo{
		logger: logger,
		msg:    msg,
	}
}

func (f *Foo) Setup(ctx context.Context) error {
	f.logger.Info("Foo Setup")
	return nil
}

func (f *Foo) Run(context.Context) error {
	f.logger.Info("Foo Run", zap.String("message", f.msg))
	return nil
}

func (f *Foo) Cleanup(context.Context) error {
	f.logger.Info("Foo Cleanup")
	return nil
}
