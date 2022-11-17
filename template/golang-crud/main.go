package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"
	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/crud"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	crud.RegisterRoutes(r, cfg)

	if err := r.Run(fmt.Sprintf(":%d", cfg.Port)); err != nil {
		panic(err)
	}
}
