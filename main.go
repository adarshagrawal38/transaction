package main

import (
	"transaction/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	endPoint := controller.NewEndpoint()

	route.POST("/transactions", endPoint.AddTransaction)
	route.GET("/statistics", endPoint.Statistics)
	route.DELETE("/transactions", endPoint.DeleteTrans)
	route.POST("/setLoc", endPoint.SetLoc)
	route.POST("/resetLoc", endPoint.ResetLoc)
	route.Run("localhost:3000")
}
