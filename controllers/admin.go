package controllers

import (
	"Week3/db"
	"Week3/forms"
	"Week3/helper"
	"Week3/helper/hash"
	"Week3/helper/jwt"
	"Week3/helper/validator"
	"Week3/models"
	"database/sql"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type AdminController struct{}

func (h AdminController) SignUp(c *gin.Context) {
	var adminRegister forms.AdminRegister
	if err := c.ShouldBindJSON(&adminRegister); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	// validate inputs
	if !validator.StringCheck(adminRegister.Name, 5, 50) {
		c.JSON(400, gin.H{"message":"invalid name"})
		return 
	}
	if !validator.StringCheck(adminRegister.Password, 5, 33) {
		c.JSON(400, gin.H{"message":"invalid password"})
		return 
	}
	nip, err := validator.ExtractNIP(adminRegister.Nip)
	if err!= nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "615"{
		c.JSON(400,gin.H{"message":"invalid nip"})
		return
	}
	// check duplicates nip
	conn := db.CreateConn()
	isExist := helper.CheckNIP(conn, nip.ToInt)
	if isExist {
		c.JSON(409, gin.H{"message":"cannot have duplicate nip"})
		return
	}
	hashed_password, _ := hash.HashPassword(adminRegister.Password)
	query := "INSERT INTO public.user (nip, name, password) VALUES ($1, $2, $3) RETURNING *"
	var admin models.User
	err = conn.QueryRowx(query, nip.ToInt,adminRegister.Name, hashed_password).StructScan(&admin)
	if err != nil {
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}

	accessToken := jwt.SignJWT(nip.ToString,"admin")

	c.JSON(201, gin.H{"message":"success", "data":gin.H{
		"userId":admin.Id,
		"nip": admin.Nip,
		"name":admin.Name,
		"accessToken":accessToken}})
}
func (h AdminController) SignIn(c *gin.Context){
	var adminLogin forms.AdminLogin
	if err := c.ShouldBindJSON(&adminLogin); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if !validator.StringCheck(adminLogin.Password, 5, 33) {
		c.JSON(400, gin.H{"message":"invalid password"})
		return 
	}
	nip, err := validator.ExtractNIP(adminLogin.Nip)
	if err!= nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "615"{
		c.JSON(404,gin.H{"message":"user not found"})
		return
	}
	var admin models.User
	conn := db.CreateConn()
	err = conn.QueryRowx("SELECT * FROM public.user WHERE nip = $1",adminLogin.Nip).StructScan(&admin)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"user not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if !hash.CheckPasswordHash(adminLogin.Password, *admin.Password){
		c.JSON(400, gin.H{"message":"invalid password"})
		return
	}
	accessToken := jwt.SignJWT(nip.ToString, "admin")
	c.JSON(200, gin.H{"message":"success","data":gin.H{
		"userId":admin.Id,
		"nip":admin.Nip,
		"name":admin.Name,
		"accessToken":accessToken,
	}})
}

func (h AdminController) GetAllUsers(c *gin.Context) {
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	userId := c.Query("userId")
	name := strings.ToLower(c.Query("name"))
	nip := c.Query("nip")
	role := c.Query("role")
	createdAt := c.Query("createdAt")
	// check nip
	_, err := strconv.ParseInt(nip, 10, 64)
	if err!= nil {
		nip = ""
	}
	// _, err = validator.ExtractNIP(nipInt)
	// if err !=nil{
	// 	nip = ""
	// }
	if role != "it" && role !="nurse"{
		role = ""
	}
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	baseQuery := "SELECT id, nip, name, \"createdAt\" FROM public.user"

	var args []interface{}
	var queryParams []string
	argIdx := 1
	if userId != ""{
		queryParams = append(queryParams, " id::text = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, userId)
		argIdx += 1
	}
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams," name LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if nip != ""{
		nipWildcard := nip + "%"
		queryParams = append(queryParams, " nip::text LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nipWildcard)
		argIdx += 1
	}
	// lazy & slow approach, should be optimized
	if role != "" {
		if role == "it"{
			queryParams = append(queryParams, " nip::text LIKE '615%' ") 
		} else {
			queryParams = append(queryParams, " nip::text LIKE '303%' ") 
		}
	}
	if len(queryParams) > 0 {
		allQuery := strings.Join(queryParams, " AND")
		baseQuery += " WHERE " + allQuery
	}
	baseQuery += " ORDER BY "
	if createdAt == "" {
		baseQuery += " \"createdAt\" DESC"
	} else {
		if createdAt == "asc"{
			baseQuery += " \"createdAt\" ASC"
		}
	}

	baseQuery +=  " LIMIT " + strconv.Itoa(limit) + " OFFSET " + strconv.Itoa(offset)
	conn := db.CreateConn()
	 
	users  := make([]models.GetUser,0)
	err = conn.Select(&users, baseQuery,args...)
	if err != nil {
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success","data":users})
}
