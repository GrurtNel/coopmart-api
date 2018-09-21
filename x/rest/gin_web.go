package rest

import (
	"github.com/gin-gonic/gin"
)

const STATUS_OK = 200

type JsonRender struct {
}

func (r *JsonRender) SendData(ctx *gin.Context, data interface{}) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   data,
		"status": "ok",
	})
}
func (r *JsonRender) Success(ctx *gin.Context) {
	ctx.JSON(STATUS_OK, map[string]interface{}{
		"data":   nil,
		"status": "ok",
	})
}
