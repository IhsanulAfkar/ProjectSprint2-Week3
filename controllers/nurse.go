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
	"time"

	"github.com/gin-gonic/gin"
)

type NurseController struct{}

func (h NurseController) CreateNurse(c *gin.Context){
	var nurseRegister forms.NurseRegister
	if err := c.ShouldBindJSON(&nurseRegister); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	// validate inputs
	if  !validator.StringCheck(nurseRegister.Name,5,50) {
		c.JSON(400, gin.H{"message":"invalid name"})
		return
	}
	if !validator.IsURL(nurseRegister.IdentityCardScanImg){
		c.JSON(400, gin.H{"message":"invalid identityCardScanImg"})
		return
	}
	nip, err := validator.ExtractNIP(nurseRegister.Nip)
	if err!= nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "303"{
		c.JSON(400,gin.H{"message":"invalid nip"})
		return
	}
	conn := db.CreateConn()
	isExist := helper.CheckNIP(conn, nip.ToInt)
	if isExist {
		c.JSON(409, gin.H{"message":"cannot have duplicate nip"})
		return
	}
	query := "INSERT INTO public.user (nip, name, \"identityCardScanImg\") VALUES ($1, $2, $3) RETURNING *"
	var nurse models.User
	err = conn.QueryRowx(query,nip.ToInt, nurseRegister.Name, nurseRegister.IdentityCardScanImg).StructScan(&nurse)
	if err != nil {
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{"message":"success", "data":gin.H{
		"userId":nurse.Id,
		"nip": nurse.Nip,
		"name":nurse.Name,
		}})
}

func (h NurseController) NurseLogin(c *gin.Context){
	var nurseLogin forms.AdminLogin
	if err := c.ShouldBindJSON(&nurseLogin); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	// if !validator.StringCheck(nurseLogin.Password, 5, 33) {
	// 	c.JSON(400, gin.H{"message":"invalid password"})
	// 	return 
	// }
	nip, err := validator.ExtractNIP(nurseLogin.Nip)
	if err!= nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "303"{
		c.JSON(404,gin.H{"message":"nurse not found"})
		return
	}
	var nurse models.User
	conn := db.CreateConn()
	err = conn.QueryRowx("SELECT * FROM public.user WHERE nip = $1",nurseLogin.Nip).StructScan(&nurse)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(400, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	// check if nurse allowed
	if nurse.Password == nil{
		c.JSON(400, gin.H{"message":"nurse not have access"})
		return
	}
	if !hash.CheckPasswordHash(nurseLogin.Password, *nurse.Password){
		c.JSON(400, gin.H{"message":"invalid password"})
		return
	}

	accessToken := jwt.SignJWT(nip.ToString, "nurse")
	c.JSON(200, gin.H{"message":"success","data":gin.H{
		"userId":nurse.Id,
		"nip":nurse.Nip,
		"name":nurse.Name,
		"accessToken":accessToken,
	}})
}

func (h NurseController)UpdateNurse(c *gin.Context){
	var nurseUpdate forms.NurseUpdate
	if err := c.ShouldBindJSON(&nurseUpdate); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	nurseId := c.Param("nurseId")
	
	conn := db.CreateConn()
	var nurse models.User
	query := "SELECT * FROM public.user WHERE id::text = $1"
	err := conn.QueryRowx(query, nurseId).StructScan(&nurse)
	if err != nil {
		if err == sql.ErrNoRows{ 
			c.JSON(404, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if strconv.FormatInt(nurse.Nip, 10)[:3] != "303"{
		c.JSON(404, gin.H{"message":"nurse not found"})
		return
	}
	if !validator.StringCheck(nurseUpdate.Name, 5,50){
		c.JSON(400, gin.H{"message":"invalid name"})
		return
	}
	nip, err := validator.ExtractNIP(nurseUpdate.Nip)
	if err != nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "303"{
		c.JSON(404,gin.H{"message":"invalid nip"})
		return
	}
	if nurse.Nip != nip.ToInt {
		res := helper.CheckNIP(conn, nip.ToInt)
		if res {
			c.JSON(409, gin.H{"message":"cannot have duplicate nip"})
			return
		}
	}
	
	// update nurse
	query = "UPDATE public.user SET nip = $1, name = $2, \"updatedAt\" = $3 WHERE id = $4"
	res, err := conn.Exec(query, nurseUpdate.Nip, nurseUpdate.Name, time.Now(), nurseId)
	if err != nil {
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0{
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success update nurse"})
}

func (h NurseController)DeleteNurse(c *gin.Context){
	nurseId:= c.Param("nurseId")

	// check if nurse exist
	conn := db.CreateConn()
	query := "SELECT * FROM public.user WHERE id = $1"
	var user models.User
	err := conn.QueryRowx(query, nurseId).StructScan(&user)
	if err!= nil{
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	nipStr := strconv.FormatInt(user.Nip, 10)

	if nipStr[:3] != "303"{
		c.JSON(404, gin.H{"message":"no nurse found"})
		return
	}
	// delete 
	query = "DELETE FROM public.user WHERE id = $1"

	res, _ := conn.Exec(query, nurseId)
	rows, _  := res.RowsAffected()
	if rows == 0 {
		c.JSON(500, gin.H{"message":"server error, no record deleted"})
		return
	}
	c.JSON(200, gin.H{"message":"delete nurse success"})
}

func (h NurseController)GrantAccess(c *gin.Context){
	nurseId := c.Param("nurseId")
	var input struct {Password string `json:"password"`}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	if !validator.StringCheck(input.Password, 5, 33) {
		c.JSON(400, gin.H{"message":"invalid password"})
		return
	}
	
	// check nurse is exist
	conn:=db.CreateConn()
	query:= "SELECT * FROM public.user WHERE id::text = $1"
	var nurse models.User
	err := conn.QueryRowx(query, nurseId).StructScan(&nurse)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(404, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	hashed_password, _ := hash.HashPassword(input.Password)
	query = "UPDATE public.user SET password = $1,  \"updatedAt\" = $2 WHERE id = $3"
	res, err := conn.Exec(query, hashed_password, time.Now(), nurseId)
	if err!= nil{
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	rows, _ := res.RowsAffected()
	if rows == 0{
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success grant access to nurse "+ nurse.Id})
}