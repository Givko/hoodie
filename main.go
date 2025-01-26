package main

import (
	"github.com/givko/hoodie/internal/api/router"
)

func main() {
	router.Init().Run(":8080")
}
