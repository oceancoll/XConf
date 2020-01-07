package source

import (
	"context"

	"github.com/micro/go-micro/config/source"
)

type baseURL struct{}

func WithURL(u string) source.Option {
	return func(o *source.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, baseURL{}, u)
	}
}
