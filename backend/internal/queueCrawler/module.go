package queueCrawler


type Module struct{
	Handler *Handler
}

func NewModule() *Module{
	service:=NewService()
	handler:=NewHandler(service)
	return &Module{Handler: handler}
}
