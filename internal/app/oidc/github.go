/*
 * Copyright 2022. The FlowTable Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package oidc

import (
	"encoding/json"
	"github.com/emicklei/go-restful/v3"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"log"
	"net/http"
	"os"
)

var (
	clientID     = os.Getenv("GITHUB_CLIENT_ID")
	clientSecret = os.Getenv("GITHUB_CLIENT_SECRET")
)

type Authenticator struct {
	clientConfig oauth2.Config
	ctx          context.Context
}

func NewGithubOauthService() *restful.WebService {
	ws := new(restful.WebService)
	githubOauthService := newAuthenticator()
	ws.Path("/auth")
	ws.Route(ws.GET("/login/github").To(githubOauthService.Login))
	ws.Route(ws.GET("/callback/github").To(githubOauthService.HandleCallback))
	return ws
}
func newAuthenticator() *Authenticator {
	ctx := context.Background()
	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  "http://localhost:8088/auth/callback/github", /* FIXME: leon 2022-06-30 use config value for domain from yaml or db */
		Scopes:       []string{"user:email", "read:user"},
	}

	return &Authenticator{
		clientConfig: config,
		ctx:          ctx,
	}
}

func (a *Authenticator) Login(request *restful.Request, response *restful.Response) {
	http.Redirect(response.ResponseWriter, request.Request, a.clientConfig.AuthCodeURL("state"), http.StatusTemporaryRedirect)
}

func (a *Authenticator) HandleCallback(request *restful.Request, response *restful.Response) {
	if request.QueryParameter("state") != "state" {
		err := response.WriteErrorString(http.StatusBadRequest, "state did not match")
		if err != nil {
			return
		}
		return
	}
	token, err := a.clientConfig.Exchange(a.ctx, request.QueryParameter("code"))
	if err != nil {
		log.Printf("no token found: %v", err)
		response.WriteHeader(http.StatusUnauthorized)
		return
	}
	// TODO: login or register an account then redirect to console home
	resp := struct {
		OAuth2Token *oauth2.Token
	}{token}
	data, err := json.MarshalIndent(resp, "", "    ")
	if err != nil {
		err := response.WriteErrorString(http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
		return
	}
	_, err = response.Write(data)
	if err != nil {
		return
	}
}
