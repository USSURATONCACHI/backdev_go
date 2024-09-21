package main

import (
	"net/http"

	"backdev_go/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	Model model.Model
	GinEngine *gin.Engine
}

func CreateServer(mdl model.Model) Server {
	server := Server {
		Model: mdl,
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
	UserUuid uuid.UUID `json:"user_uuid"`
}
type AuthorizeResponse struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func ServerAuthorize(c *gin.Context, mdl model.Model) {
	var request AuthorizeRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	
	tokens, err := mdl.CreateToken(request.UserUuid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	response := AuthorizeResponse {
		AccessToken: tokens.JwtToken,
		RefreshToken: tokens.RefreshToken.String(),
	}

	c.IndentedJSON(http.StatusOK, response)
}


type ValidateRequest struct {
	AccessToken string `json:"access_token"`
}
func ServerValidate(c *gin.Context, mdl model.Model) {
	var request ValidateRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err)
		return
	}

	result, err := mdl.ValidateToken(request.AccessToken)

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