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
	"fmt"
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
	query := "INSERT INTO nurse (nip, name, \"identityCardScanning\") VALUES ($1, $2, $3) RETURNING *"
	var nurse models.Nurse
	err = conn.QueryRowx(query,nip.ToInt, nurseRegister.Name, nurseRegister.IdentityCardScanImg).StructScan(&nurse)
	if err != nil {
		fmt.Println(err.Error())
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
	if !validator.StringCheck(nurseLogin.Password, 5, 33) {
		c.JSON(400, gin.H{"message":"invalid password"})
		return 
	}
	nip, err := validator.ExtractNIP(nurseLogin.Nip)
	if err!= nil{
		c.JSON(400,gin.H{"message":err.Error()})
		return
	}
	if nip.First3Digits != "303"{
		c.JSON(404,gin.H{"message":"nurse not found"})
		return
	}
	var nurse models.Nurse
	conn := db.CreateConn()
	err = conn.QueryRowx("SELECT * FROM nurse WHERE nip = $1",nurseLogin.Nip).StructScan(&nurse)
	if err != nil {
		if err == sql.ErrNoRows{
			c.JSON(400, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	// check if granted
	if !nurse.IsGranted {
		c.JSON(400, gin.H{"message":"nurse not have access"})
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
		c.JSON(400,gin.H{"message":"invalid nip"})
		return
	}
	conn := db.CreateConn()
	var nurse models.Nurse
	query := "SELECT * FROM nurse WHERE id::text = $1"
	err = conn.QueryRowx(query, nurseId).StructScan(&nurse)
	if err != nil {
		if err == sql.ErrNoRows{ 
			c.JSON(404, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
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
	query = "UPDATE nurse SET nip = $1, name = $2, \"updatedAt\" = $3 WHERE id = $4"
	res, err := conn.Exec(query, nurseUpdate.Nip, nurseUpdate.Name, time.Now(), nurseId)
	if err != nil {
		fmt.Println(err.Error())
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
	query := "SELECT id FROM nurse WHERE id = $1"
	var id string
	err := conn.QueryRow(query, nurseId).Scan(&id)
	if err!= nil{
		if err == sql.ErrNoRows{
			c.JSON(404, gin.H{"message":"nurse not found"})
			return
		}
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	// delete 
	query = "DELETE FROM nurse WHERE id = $1"
	res, err := conn.Exec(query, nurseId)
	if err != nil {
		fmt.Println(err.Error())
	}
	rows,_ :=res.RowsAffected()
	if rows == 0 {
		c.JSON(500, gin.H{"message":"server error, no record deleted"})
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
	query:= "SELECT * FROM nurse WHERE id::text = $1"
	var nurse models.Nurse
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
	query = "UPDATE nurse SET password = $1, \"isGranted\" = true, \"updatedAt\" = $2 WHERE id = $3"
	res, err := conn.Exec(query, hashed_password, time.Now(), nurseId)
	if err!= nil{
		fmt.Println(err.Error())
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