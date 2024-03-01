package handler

import (
	"github.com/gin-gonic/gin"
	"west.garden/template/common/render"
)

func Pong(c *gin.Context) {
	render.Json(c, render.Ok, "Pong,This is Server!")
}
