package param

import (
	"github.com/goexl/http"
	"github.com/goexl/log"
	"github.com/goexl/valuer/internal/internal"
)

type Parser struct {
	Logger      log.Logger
	Expressions []internal.Expression
	Http        *http.Client
}

func NewParser() *Parser {
	return &Parser{
		Logger:      log.New().Apply(),
		Expressions: make([]internal.Expression, 0),
		Http:        http.New().Build(),
	}
}
