package param

import (
	"github.com/goexl/http"
	"github.com/goexl/log"
	"github.com/goexl/valuer/internal/core"
)

type Parser struct {
	Logger      log.Logger
	Expressions []core.Expression
	Http        *http.Client
}

func NewParser() *Parser {
	return &Parser{
		Logger:      log.New().Apply(),
		Expressions: make([]core.Expression, 0),
		Http:        http.New().Build(),
	}
}
