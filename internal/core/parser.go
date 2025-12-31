package core

import (
	"encoding/json"
	"strings"

	"github.com/drone/envsubst"
	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/vm"
	"github.com/goexl/env"
	"github.com/goexl/gox"
	"github.com/goexl/gox/check"
	"github.com/goexl/gox/field"
	"github.com/goexl/valuer/internal/constant"
	"github.com/goexl/valuer/internal/param"
)

type Parser struct {
	vm      *vm.VM
	options []expr.Option
	params  *param.Parser
}

func NewParser(params *param.Parser) (g *Parser) {
	g = new(Parser)
	g.vm = new(vm.VM)
	g.params = params
	g.options = []expr.Option{
		expr.AllowUndefinedVariables(),
		expr.Function(constant.FuncFile, g.file),
		expr.Function(constant.FuncUrl, g.url),
		expr.Function(constant.FuncHttp, g.url),
		expr.Function(constant.FuncMatch, g.match),
	}
	for _, expression := range params.Expressions {
		g.options = append(g.options, expr.Function(expression.Name(), expression.Exec))
	}

	return
}

func (p *Parser) Parse(key string) (value string) {
	defer p.recover()

	if got := env.Get(key); "" != strings.TrimSpace(got) {
		value = got
	}
	if got := p.eval(key); "" != strings.TrimSpace(got) {
		value = got
	}
	if "" == value { // 如果环境变量取值没有改变，证明键没有环境变量，需要将键值赋值
		value = key
	}

	size := len(value)
	if constant.JsonObjectStart == (value)[0:1] && constant.JsonObjectEnd == (value)[size-1:size] {
		value = p.fixJsonObject(value)
	} else if constant.JsonArrayStart == (value)[0:1] && constant.JsonArrayEnd == (value)[size-1:size] {
		value = p.fixJsonArray(value)
	} else {
		value = p.expr(value)
	}

	return
}

func (p *Parser) expr(from string) (to string) {
	fields := gox.Fields[any]{
		field.New("Expression", from),
	}
	if program, ce := expr.Compile(from, p.options...); nil != ce {
		to = from
		p.params.Logger.Debug("表达式编译出错", field.Error(ce), fields...)
	} else if result, re := p.vm.Run(program, nil); nil != re {
		to = from
		p.params.Logger.Debug("表达式运算出错", field.Error(re), fields...)
	} else {
		to = gox.ToString(result)
		p.params.Logger.Debug("表达式运算成功", field.New("result", to), fields...)
	}

	return
}

func (p *Parser) fixJsonObject(from string) (to string) {
	object := make(map[string]any)
	if ue := json.Unmarshal([]byte(from), &object); nil != ue {
		to = from
	} else {
		p.fixObjectExpr(object)
	}

	if from == to {
		// 不需要进行转换
	} else if bytes, me := json.Marshal(object); nil != me {
		to = from
	} else {
		to = string(bytes)
	}

	return
}

func (p *Parser) fixJsonArray(from string) (to string) {
	array := make([]any, 0)
	if ue := json.Unmarshal([]byte(from), &array); nil != ue {
		to = from
	} else {
		p.fixArrayExpr(&array)
	}

	if from == to {
		// 不需要进行转换
	} else if bytes, me := json.Marshal(array); nil != me {
		to = from
	} else {
		to = string(bytes)
	}

	return
}

func (p *Parser) fixArrayExpr(array *[]any) {
	for index, value := range *array {
		switch vt := value.(type) {
		case string:
			(*array)[index] = p.expr(vt)
		case []any:
			p.fixArrayExpr(&vt)
		case map[string]any:
			p.fixObjectExpr(vt)
		}
	}
}

func (p *Parser) fixObjectExpr(object map[string]any) {
	for key, value := range object {
		switch vt := value.(type) {
		case string:
			object[key] = p.expr(vt)
		case []any:
			p.fixArrayExpr(&vt)
		case map[string]any:
			p.fixObjectExpr(vt)
		}
	}
}

func (p *Parser) eval(from string) (to string) {
	to = from
	if !strings.Contains(to, constant.Dollar) {
		return
	}

	count := 0
	for {
		if value, ee := envsubst.Eval(to, env.Get); nil == ee {
			to = value
		}

		if count >= 2 || !strings.Contains(to, constant.Dollar) {
			break
		}
		if strings.Contains(to, constant.Dollar) {
			count++
		}
	}

	return
}

func (p *Parser) isHttp(url string) bool {
	return check.New().
		Any().
		String(url).
		Items(constant.PrefixHttpProtocol, constant.PrefixHttpsProtocol).
		Prefix().
		Check()
}

func (p *Parser) recover() {
	if ctx := recover(); nil != ctx {
		switch value := ctx.(type) {
		case error:
			p.params.Logger.Warn("获取器执行出错", field.Error(value))
		}
	}
}
