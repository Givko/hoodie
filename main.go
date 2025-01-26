package main

import (
	"github.com/plamendelchev/hoodie/internal/api/router"
)

func main() {
	router.Init().Run(":8080")
}
