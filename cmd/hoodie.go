package main

import (
	"github.com/plamendelchev/hoodie/internal/router"
)

func main() {
	router.Init().Run(":8080")
}
