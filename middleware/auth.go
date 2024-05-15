package middleware

import (
	"Week3/db"
	"Week3/helper/jwt"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)
func getBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New("bad header value given")
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New("incorrectly formatted authorization header")
	}

	return jwtToken[1], nil
}
func AdminAuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err!= nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	id, err := jwt.ParseToken(token)
	if err != nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message":err.Error()})
		return
	}
	if id.Role != "admin"{
		c.AbortWithStatusJSON(401, gin.H{"message":"access not allowed"})
		return
	}
	// find user
	conn := db.CreateConn()
	res, err := conn.Exec("SELECT 1 FROM public.user WHERE nip = $1 LIMIT 1",id.Nip)
	if err != nil{
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(500, gin.H{
			"message":"server error"})
			return
		}
	if rows,_:=res.RowsAffected();rows == 0 {
		c.AbortWithStatusJSON(404, gin.H{
			"message":"user not found"})
			return
	}
	c.Set("userNip",id.Nip)
	c.Next()
}

func AllAuthMiddleware(c *gin.Context) {
	token, err := getBearerToken(c.GetHeader("Authorization"))
	if err!= nil {
		c.AbortWithStatusJSON(401, gin.H{
			"message": err.Error()})
		return
	}
	id, err := jwt.ParseToken(token)
	if err != nil {
		fmt.Println(err.Error())
		c.AbortWithStatusJSON(401, gin.H{
			"message":err.Error()})
		return
	}
	
	// find user
	conn := db.CreateConn()
	var isExists bool
	err = conn.QueryRow("SELECT EXISTS (SELECT 1 FROM public.user WHERE nip = $1) AS is_exists",id.Nip).Scan(&isExists)
	if err != nil || !isExists{
		c.AbortWithStatusJSON(404, gin.H{"message":"user not found"})
		return
	}
	fmt.Println("token")
	fmt.Println(id)
	c.Set("userNip",id.Nip)
	c.Set("role",id.Role)
	c.Next()
}