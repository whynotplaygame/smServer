package gate

import (
	"smServer/net"
	"smServer/server/gate/controller"
)

var Router = &net.Router{}

func Init() {
	initRouter()
}

func initRouter() {
	controller.GateHandler.Router(Router)
}
