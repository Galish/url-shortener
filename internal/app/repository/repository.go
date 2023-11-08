package repository

import "context"

type Repository interface {
	Get(context.Context, string) (string, error)
	Set(context.Context, string, string) error
	SetBatch(context.Context, ...[2]string) error
	Has(context.Context, string) bool
	Ping(context.Context) (bool, error)
	Close() error
}
