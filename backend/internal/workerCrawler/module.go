package workerCrawler

type Module struct{
	Handler *Handler
}

func NewModule()*Module{
	service:=NewService()
	handler:=NewHandler(service)
	return &Module{handler}
}