package webService

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller"
)

type Controller struct {
	controller.Controller
}

func (w Controller) Index(request *restful.Request, response *restful.Response) {
	w.RenderSuccess(response, "Controller Manager Service Version "+common.Version)
}
