package server

import (
	"github.com/Sora233/DDC/docs"
	"github.com/Sora233/DDC/server/config"
	"github.com/Sora233/DDC/server/middleware"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RunServer() error {
	g := gin.New()
	g.Use(
		gin.Recovery(),
		middleware.RequestIDMiddleware(),
		middleware.WrapperMiddleware("ddc-server"),
	)
	router := g.Group(config.Global.APIPrefix)
	v1 := router.Group("/v1")

	docs.SwaggerInfo.BasePath = config.Global.APIPrefix

	if config.Global.EnableSwagger {
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	v1.GET("/ping", Ping)
	v1.GET("/danmu", GetSearchDanmu)
	v1.GET("/money", GetSearchMoney)

	return g.Run(config.Global.Addr)
}
