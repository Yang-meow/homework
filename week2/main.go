package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/get_user", LogicDemo)
	r.Run()
}

func LogicDemo(c *gin.Context) {
	uidStr := c.Query("user_id")
	userID, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, "invalid uid")
		return
	}
	user, err := GetUserByUserID(context.Background(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, "not found")
			return
		} else {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, err)
			return
		}
	}

	// deal with users
	fmt.Println(user.Name, " ", user.Age)
	c.JSON(http.StatusOK, "success")
}
