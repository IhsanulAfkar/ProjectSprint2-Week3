package controllers

import (
	"Week3/db"
	"Week3/forms"
	"Week3/helper"
	"Week3/helper/validator"
	"Week3/models"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type PatientController struct{}

func (h PatientController) CreatePatient(c *gin.Context) {
	var patientForm forms.PatientCreate
	if err := c.ShouldBindJSON(&patientForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	fmt.Println(patientForm)
	identityStr := strconv.FormatInt(patientForm.IdentityNumber, 10)
	if len(identityStr) != 16 {
		c.JSON(400, gin.H{"message":"incorrect identityNumber"})
		return
	}
	if patientForm.PhoneNumber[:3] != "+62" || !validator.StringCheck(patientForm.PhoneNumber, 10, 15){
		c.JSON(400, gin.H{"message":"incorrect phoneNumber"})
		return
	}
	if !validator.StringCheck(patientForm.Name, 3, 30){
		c.JSON(400, gin.H{"message":"incorrect name"})
		return
	}
	if !validator.IsDateISO860(patientForm.BirthDate) {
		c.JSON(400, gin.H{"message":"incorrect birthDate"})
		return
	}
	if !helper.Includes(patientForm.Gender, models.Gender[:]){
		c.JSON(400, gin.H{"message":"incorrect gender"})
		return
	}
	if !validator.IsURL(patientForm.IdentityCardScanImg){
		c.JSON(400, gin.H{"message":"incorrect identity card"})
		return
	}
	// check if identity exist
	conn := db.CreateConn()
	res, _ := conn.Exec("SELECT * FROM patient WHERE \"identityNumber\" = $1 ", patientForm.IdentityNumber)
	if rows,_:= res.RowsAffected(); rows > 0 {
		c.JSON(409, gin.H{"message":"identityNumber exist"})
		return
	}

	query := `INSERT INTO patient ("identityNumber", name, "phoneNumber", "birthDate", gender, "identityCardScanning") VALUES ($1,$2,$3,$4,$5,$6)`
	res, err := conn.Exec(query, patientForm.IdentityNumber, patientForm.Name, patientForm.PhoneNumber, patientForm.BirthDate, patientForm.Gender, patientForm.IdentityCardScanImg)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if rows, _ := res.RowsAffected(); rows ==0{
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{"message":"patient added successfully"})
}

func (h PatientController) GetAllPatient(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	identityNumber := c.Query("identityNumber")
	name := strings.ToLower(c.Query("name"))
	phoneNumber := c.Query("phoneNumber")
	createdAt := c.Query("createdAt")
	identityNumInt, err := strconv.ParseInt(identityNumber, 10, 64)
	if err != nil{
		identityNumber = ""
	}
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	baseQuery := "SELECT * FROM patient"

	var args []interface{}
	// queryParams := make(map[string]interface{})
	var queryParams []string
	argIdx := 1
	fmt.Println("name")
	fmt.Println(name)
	if name != ""{
		nameWildcard := "%" + name +"%"
		queryParams = append(queryParams," name ILIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nameWildcard)
		argIdx += 1
	}
	if phoneNumber!=""{
		phoneWildcard := "+" + phoneNumber + "%"
		queryParams = append(queryParams, " \"phoneNumber\" LIKE $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, phoneWildcard)
		argIdx += 1
	}
	if identityNumber != ""{
		queryParams = append(queryParams, " \"identityNumber\" = $"+strconv.Itoa(argIdx) + " ")
		args = append(args, identityNumInt)
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
	fmt.Println(baseQuery)
	patients := make([]models.Patient,0)
	err = conn.Select(&patients, baseQuery, args...)
	if err != nil{
		fmt.Println(err.Error())
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(200, gin.H{"message":"success","data":patients})
}

