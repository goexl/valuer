package valuer

import (
	"github.com/goexl/valuer/internal/builder"
)

var _ = New

// New 创建
func New() *builder.Parser {
	return builder.NewParser()
}
