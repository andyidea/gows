// method.go

package ws

type HandlerFunc func(*Context)

type Method struct {
	router      string
	handlerFunc HandlerFunc
}

type Methods []Method

var methods Methods

func (self Methods) Get(r string) (hf HandlerFunc, ok bool) {
	ok = false
	for _, method := range self {
		if method.router == r {
			ok = true
			hf = method.handlerFunc
			return
		}
	}

	return
}
