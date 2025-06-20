package routers

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller/alertmanager"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller/oauthService"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller/prometheus"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller/webService"
)

func WebService() *restful.WebService {
	service := new(restful.WebService)
	service.Route(service.GET("/").To(webService.Controller{}.Index))
	service.Route(service.POST("/token").To(oauthService.Render().GetUserToken))
	service.Route(service.GET("/grafana/authorize").To(oauthService.Render().GetUserAuthorize))
	service.Route(service.POST("/grafana/access_token").To(oauthService.Render().GenerateAccessToken))
	service.Route(service.GET("/grafana/user").To(oauthService.Render().GetUserInfo))
	return service
}

func APIsService() *restful.WebService {
	service := new(restful.WebService)
	//service.Path("/api").Produces(restful.MIME_JSON)
	service.Route(service.POST("/sync_rules").To(prometheus.Controller{}.SyncRules))
	service.Route(service.DELETE("/delete_rules").To(prometheus.Controller{}.DeleteRules))
	return service
}

func AlertmanagerService() *restful.WebService {
	service := new(restful.WebService)
	service.Route(service.GET("/get_config").To(alertmanager.Controller{}.GetConfig))
	service.Route(service.POST("/overlay_config").To(alertmanager.Controller{}.OverlayConfig))
	service.Route(service.POST("/alert_template").To(alertmanager.Controller{}.AlertTemplate))
	return service
}
