package valuer

var _ = New

// New 创建
func New() *builder {
	return newBuilder()
}
