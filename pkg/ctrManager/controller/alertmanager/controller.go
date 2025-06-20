package alertmanager

import (
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller"
	"github.com/yunling101/ControllerManager/pkg/plugin/monitor"
)

type Controller struct {
	controller.Controller
}

func (w Controller) set() *monitor.Client {
	v := common.Flags().Alertmanager
	return monitor.NewAlertManager(v.AlertmanagerBaseDir).
		SetPrefix(v.AlertmanagerConfigName).
		SetBasicAuth(v.AlertmanagerBasicUsername, v.AlertmanagerBasicPassword).
		Listen(v.AlertmanagerListenAddress).
		Reload(v.AlertmanagerReload)
}

func (w Controller) GetConfig(request *restful.Request, response *restful.Response) {
	cfg, err := w.set().LoadFile()
	if err != nil {
		w.RenderFail(response, err.Error())
		return
	}
	w.RenderSuccess(response, cfg)
}

func (w Controller) OverlayConfig(request *restful.Request, response *restful.Response) {
	body, err := w.ReadBody(request, response)
	if err != nil {
		return
	}

	err = w.set().ModifyFile(body)
	if err != nil {
		w.Error(request, err.Error()).RenderFail(response, err.Error())
		return
	}

	w.RenderSuccess(response, nil)
}

func (w Controller) AlertTemplate(request *restful.Request, response *restful.Response) {
	body, err := w.ReadBody(request, response)
	if err != nil {
		return
	}

	t := make(map[string]string)
	if err = json.Unmarshal(body, &t); err != nil {
		w.RenderFail(response, err.Error())
		return
	}

	err = w.set().Templates(t["filename"], t["content"])
	if err != nil {
		w.Error(request, err.Error()).RenderFail(response, err.Error())
		return
	}

	w.RenderSuccess(response, nil)
}
