package valuer

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/simaqian"
)

type params struct {
	logger      simaqian.Logger
	expressions []expression
	http        *resty.Client
}

func newParams() *params {
	return &params{
		logger:      simaqian.Default(),
		expressions: make([]expression, 0),
		http:        resty.New(),
	}
}
