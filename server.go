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

	server.GinEngine.POST(
		"/refresh", 
		func(ctx *gin.Context) { ServerRefresh(ctx, server.Model) },
	)

	server.GinEngine.Static("/", "./static")

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
type ErrorJson struct {
	Error string `json:"error"`
}

func ServerAuthorize(c *gin.Context, mdl model.Model) {
	var request AuthorizeRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorJson {
			Error: err.Error(),
		})
		return
	}
	
	tokens, err := mdl.CreateToken(request.UserUuid, c.ClientIP())
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJson {
			Error: err.Error(),
		})
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
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorJson {
			Error: err.Error(),
		})
		return
	}

	result, err := mdl.ValidateToken(request.AccessToken)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorJson {
			Error: err.Error(),
		})
		return
	}
	if result {
		c.AbortWithStatus(http.StatusOK)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

}


type RefreshRequest struct {
	AccessToken string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	EmulateIp string `json:"emulate_ip"`
}
type RefreshResponse AuthorizeResponse;
func ServerRefresh(c *gin.Context, mdl model.Model) {
	var request RefreshRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorJson { Error: err.Error() })
		return
	}

	refreshTokenUuid, err := uuid.Parse(request.RefreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorJson { Error: err.Error() })
		return
	}

	ip := c.ClientIP()
	if len(request.EmulateIp) > 0 {
		ip = request.EmulateIp
	}

	tokens, clientErr, serverErr := mdl.RefreshToken(request.AccessToken, refreshTokenUuid, ip)

	if clientErr != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorJson { Error: clientErr.Error() })
		return
	}
	if serverErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorJson { Error: serverErr.Error() })
		return
	}

	response := RefreshResponse {
		AccessToken: tokens.JwtToken,
		RefreshToken: tokens.RefreshToken.String(),
	}

	c.IndentedJSON(http.StatusOK, response)
}