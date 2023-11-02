package builder

import (
	"github.com/goexl/http"
	"github.com/goexl/log"
	"github.com/goexl/valuer/internal/core"
	"github.com/goexl/valuer/internal/internal"
	"github.com/goexl/valuer/internal/param"
)

type Parser struct {
	params *param.Parser
}

func NewParser() *Parser {
	return &Parser{
		params: param.NewParser(),
	}
}

func (p *Parser) Expression(expression internal.Expression) (builder *Parser) {
	p.params.Expressions = append(p.params.Expressions, expression)
	builder = p

	return
}

func (p *Parser) Logger(logger log.Logger) (builder *Parser) {
	p.params.Logger = logger
	builder = p

	return
}

func (p *Parser) Http(http *http.Client) (builder *Parser) {
	p.params.Http = http
	builder = p

	return
}

func (p *Parser) Build() *core.Parser {
	return core.NewParser(p.params)
}
