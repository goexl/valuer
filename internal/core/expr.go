package core

import (
	"os"
	"regexp"

	"github.com/goexl/exception"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
)

func (p *Parser) file(args ...any) (result any, err error) {
	name := ""
	if 0 == len(args) {
		err = exception.New().Message("必须传入参数").Field(field.New("args", args)).Build()
	} else {
		name = gox.ToString(args[0])
	}
	if nil != err {
		return
	}

	fields := gox.Fields[any]{
		field.New("filename", name),
	}
	if bytes, re := os.ReadFile(name); nil != re {
		p.params.Logger.Error("读取文件出错", field.Error(re), fields...)
	} else {
		result = string(bytes)
		p.params.Logger.Debug("读取文件成功", field.New("content", result), fields...)
	}

	return
}

func (p *Parser) url(args ...any) (result any, err error) {
	url := ""
	if 0 == len(args) {
		err = exception.New().Message("必须传入参数").Field(field.New("args", args)).Build()
	} else {
		url = gox.ToString(args[0])
		err = gox.If(p.isHttp(url), exception.New().Message("必须是URL地址").Field(field.New("url", url))).Build()
	}
	if nil != err {
		return
	}

	fields := gox.Fields[any]{
		field.New("url", url),
	}
	if rsp, re := p.params.Http.R().Get(url); nil != re {
		p.params.Logger.Error("读取端点出错", field.Error(re), fields...)
	} else if rsp.IsError() {
		codeField := field.New("code", rsp.StatusCode())
		bodyField := field.New("body", rsp.Body())
		p.params.Logger.Warn("远端服务器返回错误", codeField, fields.Add(bodyField)...)
	} else {
		result = string(rsp.Body())
		p.params.Logger.Debug("读取端点成功", field.New("content", result), fields...)
	}

	return
}

func (p *Parser) match(args ...any) (result any, err error) {
	if 2 != len(args) {
		err = exception.New().Message("参数错误").Field(
			field.New("args", args),
			field.New("need", 2),
			field.New("real", 1),
		).Build()
	}
	if nil != err {
		return
	}

	reg := regexp.MustCompile(args[1].(string))
	result = reg.FindStringSubmatch(args[0].(string))

	return
}
