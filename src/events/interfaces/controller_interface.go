package interfaces

import (
	"github.com/gin-gonic/gin"
)

type EventController interface {
	Create(*gin.Context)
	Get(*gin.Context)
	//GetAll(*gin.Context)
	//Update(*gin.Context)
	//Delete(*gin.Context)
}
