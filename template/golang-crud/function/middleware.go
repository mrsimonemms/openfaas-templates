package function

import "github.com/gin-gonic/gin"

// Middleware is run before every CRUD endpoint. It is not run on any additional endpoints
// @link https://gin-gonic.com/docs/examples/custom-middleware
var Middleware []gin.HandlerFunc = []gin.HandlerFunc{}
