package main

import (
	"fmt"

	"github.com/mrsimonemms/openfaas-templates/template/golang-crud/pkg/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}
