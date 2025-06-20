package ctrManager

import (
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/routers"
	"github.com/yunling101/ControllerManager/pkg/plugin/jwt"
	"log"
	"net/http"
	"time"
)

func server(listen string) (server *http.Server) {
	server = &http.Server{
		Addr:           listen,
		Handler:        restful.DefaultContainer,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 16 << 20,
	}

	log.Printf("listen to %s", listen)
	return server
}

func Run(listen string) error {
	restful.DefaultContainer.Add(routers.WebService())
	//restful.DefaultContainer.Add(routers.APIsService().Filter(jwt.AuthToken))
	return server(listen).ListenAndServe()
}

func RunPrometheus() error {
	//auth := jwt.File(common.Flags().Prometheus.SecretKeyFile)
	restful.DefaultContainer.Add(routers.APIsService().Filter(jwt.Env().AuthToken))
	return server(common.Flags().Prometheus.ListenAddress).ListenAndServe()
}

func RunAlertmanager() error {
	//auth := jwt.File(common.Flags().Alertmanager.SecretKeyFile)
	restful.DefaultContainer.Add(routers.AlertmanagerService().Filter(jwt.Env().AuthToken))
	return server(common.Flags().Alertmanager.ListenAddress).ListenAndServe()
}
