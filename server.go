package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Model Model
	GinEngine *gin.Engine
}

func CreateServer(model Model) Server {
	server := Server {
		Model: model,
		GinEngine: gin.Default(),
	}


	server.GinEngine.POST(
		"/authorize", 
		func(ctx *gin.Context) { ServerAuthorize(ctx, server.Model) },
	)

	server.GinEngine.POST(
		"/validate", 
		func(ctx *gin.Context) { ServerValidate(ctx, server.Model) },
	)

	return server
}

func (server Server) Run(ip string) {
	server.GinEngine.Run(ip)
}


type AuthorizeRequest struct {
	UserUuid string `json:"user_uuid"`
}
type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ServerAuthorize(c *gin.Context, model Model) {
	var request AuthorizeRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}
	
	token, err := model.CreateTokenString(request.UserUuid)
	
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := AuthorizeResponse {
		AccessToken: token,
		RefreshToken: "none",
	}

	c.IndentedJSON(http.StatusOK, response)
}


type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}
func ServerValidate(c *gin.Context, model Model) {
	var request ValidateRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	result, err := model.ValidateToken(request.AccessToken)

	if err != nil {
		c.AbortWithError(http.StatusUnauthorized, err)
		return
	}
	if result {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}