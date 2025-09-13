package router

import (
	"dilu/modules/search/apis"

	"github.com/gin-gonic/gin"
)

func init() {
	routerNoCheckRole = append(routerNoCheckRole, registerProxyRouter)
}

// 默认需登录认证的路由
func registerProxyRouter(v1 *gin.RouterGroup) {
	r := v1.Group("")
	{
		r.GET("/search", apis.ApiSearchApi.SearchGet)
		r.POST("/search", apis.ApiSearchApi.SearchPost)
		r.GET("test", apis.ApiSearchApi.Test)
	}
}
