package middlewares

import (
	"github.com/DmitriySolopenkov/distribCalc.rev2/internal/orchestrator/repositories"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/jwt"
	"github.com/DmitriySolopenkov/distribCalc.rev2/pkg/response"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := repositories.UserRepository().GetById(jwt.New().JwtUserId(c)); err != nil {
			response.BadRequest(c, "the current user does not exist or has been logged out")
			c.Abort()
			return
		}

		c.Next()
	}
}
