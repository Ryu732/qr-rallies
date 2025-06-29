// main.go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    
    router.GET("/test", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status": "ok",
        })
    })

    // ポート 8080 で起動
    router.Run(":8080")
}
