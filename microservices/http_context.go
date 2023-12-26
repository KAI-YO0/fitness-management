package microservices

type (
	// HttpContext implement IContext it is context for HTTP service
	HttpContext struct {
		ms *Microservice
	}
	IHttpContext interface {
		IContext
	}
)

// NewHttpContext is the constructor function for HttpContext
func NewHttpContext(ms *Microservice) *HttpContext {
	return &HttpContext{
		ms: ms,
	}
}
