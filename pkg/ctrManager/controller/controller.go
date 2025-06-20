package controller

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/logf"
	"io"
)

const defaultSkip = 4

type Controller struct {
	skip int
}

func (w *Controller) Skip(skip int) *Controller {
	w.skip = skip
	return w
}

func (w *Controller) getSkip() (skip int) {
	skip = w.skip
	if skip == 0 {
		skip = defaultSkip
	}
	return
}

func (w *Controller) Render(response *restful.Response, value interface{}) {
	_ = response.WriteJson(value, "application/json; charset=UTF-8")
}

func (w *Controller) Info(request *restful.Request, value interface{}) *Controller {
	// log.Printf("%s path:%s %v", color.CyanString("[INFO]"), request.Request.RequestURI, value)
	logf.Logger().Skip(w.getSkip()).InfoF("path:%s %v", request.Request.RequestURI, value)
	return w
}

func (w *Controller) Error(request *restful.Request, value interface{}) *Controller {
	// log.Printf("%s path:%s %v", color.RedString("[ERROR]"), request.Request.RequestURI, value)
	logf.Logger().Skip(w.getSkip()).ErrorF("path:%s %v", request.Request.RequestURI, value)
	return w
}

func (w *Controller) RenderFail(response *restful.Response, value interface{}) {
	w.Render(response, map[string]interface{}{"code": 0, "message": value})
}

func (w *Controller) RenderSuccess(response *restful.Response, value interface{}) {
	w.Render(response, map[string]interface{}{"code": 1, "data": value})
}

func (w *Controller) ReadBody(request *restful.Request, response *restful.Response) (body []byte, err error) {
	body, err = io.ReadAll(request.Request.Body)
	if err != nil {
		w.Error(request, err.Error()).RenderFail(response, err.Error())
		return
	}
	if len(body) == 0 {
		err = logf.Error("%s", "request body cannot be empty")
		w.RenderFail(response, err.Error())
	}
	return
}
