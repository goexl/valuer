package valuer

import (
	"github.com/go-resty/resty/v2"
	"github.com/goexl/simaqian"
)

type builder struct {
	params *params
}

func newBuilder() *builder {
	return &builder{
		params: newParams(),
	}
}

func (b *builder) Expression(expression expression) (builder *builder) {
	b.params.expressions = append(b.params.expressions, expression)
	builder = b

	return
}

func (b *builder) Logger(logger simaqian.Logger) (builder *builder) {
	b.params.logger = logger
	builder = b

	return
}

func (b *builder) Http(http *resty.Client) (builder *builder) {
	b.params.http = http
	builder = b

	return
}

func (b *builder) Build() *Parser {
	return newParser(b.params)
}
