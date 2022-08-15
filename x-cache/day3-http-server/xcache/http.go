package xcache

const basePath = "/_xcache/"

type HttpPool struct {
	self     string
	basePath string
}

func NewHttpPool(self string) *HttpPool {
	return &HttpPool{
		self:     self,
		basePath: basePath,
	}
}

func (p *HttpPool) Log(format string, v ...interface{}) {

}
