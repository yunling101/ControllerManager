package prometheus

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller"
	"github.com/yunling101/ControllerManager/pkg/plugin/monitor"
)

type Controller struct {
	controller.Controller
}

func (w Controller) set() *monitor.Client {
	v := common.Flags().Prometheus
	return monitor.NewPrometheus(v.PrometheusRulesStoreDir).
		SetPrefix(v.PrometheusRulesPrefix).
		SetSuffix(v.PrometheusRulesSuffix).
		Listen(v.PrometheusListenAddress).
		SetBasicAuth(v.PrometheusBasicUsername, v.PrometheusBasicPassword).
		Reload(v.PrometheusReload)
}

func (w Controller) SyncRules(request *restful.Request, response *restful.Response) {
	body, err := w.ReadBody(request, response)
	if err != nil {
		return
	}

	err = w.set().AddRules(string(body))
	if err != nil {
		w.Error(request, err.Error()).RenderFail(response, err.Error())
		return
	}

	w.RenderSuccess(response, nil)
}

func (w Controller) DeleteRules(request *restful.Request, response *restful.Response) {
	body, err := w.ReadBody(request, response)
	if err != nil {
		return
	}

	err = w.set().DeleteRule(string(body))
	if err != nil {
		w.Error(request, err.Error()).RenderFail(response, err.Error())
		return
	}

	w.RenderSuccess(response, nil)
}
