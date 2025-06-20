package oauthService

import (
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/yunling101/ControllerManager/common"
	"github.com/yunling101/ControllerManager/models/q"
	"github.com/yunling101/ControllerManager/models/user"
	"github.com/yunling101/ControllerManager/pkg/cache"
	"github.com/yunling101/ControllerManager/pkg/ctrManager/controller"
	"github.com/yunling101/ControllerManager/pkg/plugin/oauth"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type handler struct {
	bind struct {
		user user.User
	}
	controller.Controller
}

func Render() *handler {
	return &handler{}
}

func (w handler) GetUserToken(request *restful.Request, response *restful.Response) {
	data := request.Request.FormValue("data")
	if data == "" {
		w.RenderFail(response, "参数不能为空")
		return
	}
	state := oauth.New().GetState()
	cache.Store(state["oauth_state"].(string), data)
	w.RenderSuccess(response, state)
}

func (w handler) GetUserAuthorize(request *restful.Request, response *restful.Response) {
	state := request.QueryParameter("state")
	redirectUri := request.QueryParameter("redirect_uri")
	cookie, err := request.Request.Cookie("oauth_grafana")
	if err != nil {
		w.RenderFail(response, "请重新登录")
		return
	}
	value, _ := url.QueryUnescape(cookie.Value)
	oauthGrafana, ok := cache.Load(value)
	if !ok {
		w.RenderFail(response, "oauth_grafana authentication failed")
		return
	}

	err = json.Unmarshal([]byte(oauthGrafana.(string)), &w.bind.user)
	if err != nil {
		w.RenderFail(response, "oauth_grafana unmarshal failed")
		return
	}

	redirect := fmt.Sprintf("%s?code=%v&state=%s", redirectUri, w.bind.user.Id, state)
	http.Redirect(response.ResponseWriter, request.Request, redirect, http.StatusFound)
}

func (w handler) GenerateAccessToken(request *restful.Request, response *restful.Response) {
	code := request.Request.FormValue("code")
	token := oauth.New().SetCode(code).GetToken()
	cache.Store(token["access_token"].(string), code)
	w.Render(response, token)
}

func (w handler) GetUserInfo(request *restful.Request, response *restful.Response) {
	authorization := request.Request.Header.Get("Authorization")
	if authorization == "" {
		w.RenderFail(response, "认证失败")
		return
	}
	w.Render(response, w.getInfo(authorization))
}

func (w *handler) getInfo(authorization string) (response common.Request) {
	auth := strings.Split(authorization, " ")
	if len(auth) == 2 {
		if auth[0] == "Bearer" {
			if id, ok := cache.Load(auth[1]); ok {
				i, err := strconv.Atoi(fmt.Sprintf("%v", id))
				if err == nil && i != 0 {
					err = q.Table(w.bind.user.TableName()).QueryOne(q.M{"id": i}, &w.bind.user)
					if err == nil {
						response = common.Request{
							"name":  w.bind.user.Nickname,
							"login": w.bind.user.Username,
							"id":    fmt.Sprintf("%v", w.bind.user.Id),
							"type":  "User",
							"email": w.bind.user.Email,
						}
					}
				}
			}
		}
	}
	return
}
