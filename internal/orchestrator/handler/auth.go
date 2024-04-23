package handler

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/services"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/jwt"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/response"
	"github.com/gin-gonic/gin"
)

type Auth struct {
	Route *gin.RouterGroup
}

type AuthRequest struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Register a new user
// @Tags Auth
// @ID auth-register
// @Accept json
// @Param input body AuthRequest true "fields"
// @Produce json
// @Success 200
// @Router /register [post]
func (a *Auth) Register(ctx *gin.Context) {
	var request AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "invalid data")
		return
	}

	_, err := services.UserService().Create(request.Login, request.Password)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Data(ctx, "successful")
}

// @Summary Authorization a user
// @Tags Auth
// @ID auth-login
// @Accept json
// @Param input body AuthRequest true "fields"
// @Produce json
// @Success 200
// @Router /login [post]
func (a *Auth) Login(ctx *gin.Context) {
	var request AuthRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.BadRequest(ctx, "invalid data")
		return
	}

	userId := services.UserService().Authorization(request.Login, request.Password)

	if userId == 0 {
		response.BadRequest(ctx, "authorization failed")
		return
	}

	token, err := jwt.New().CreateUserToken(userId)
	if err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}

	response.Data(ctx, map[string]interface{}{
		"token":   token.Token,
		"user_id": userId,
	})
}
