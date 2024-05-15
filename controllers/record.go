package controllers

import (
	"Week3/db"
	"Week3/forms"
	"Week3/helper/validator"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RecordController struct{}

func (h RecordController) CreateRecord(c *gin.Context){
	userNip := fmt.Sprintf("%s", c.MustGet("userNip"))
	nipInt, err := strconv.ParseInt(userNip, 10, 64)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var recordForm forms.CreateRecord
	if err := c.ShouldBindJSON(&recordForm); err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
		return
    }
	identityStr := strconv.FormatInt(recordForm.IdentityNumber, 10)
	if len(identityStr) != 16{
		c.JSON(400, gin.H{"message":"invalid identity"})
		return
	}
	fmt.Println(len(recordForm.Symptoms))
	if !validator.StringCheck(recordForm.Symptoms, 1, 2000){
		c.JSON(400, gin.H{"message":"invalid symptoms"})
		return	
	}
	if !validator.StringCheck(recordForm.Medications, 1, 2000){
		c.JSON(400, gin.H{"message":"invalid medications"})
		return
	}
	// check if patient exists
	query := "SELECT id FROM patient WHERE \"identityNumber\" = $1"
	var patientId string
	conn := db.CreateConn()
	err = conn.QueryRow(query, recordForm.IdentityNumber).Scan(&patientId)
	if err != nil {
		if err ==  sql.ErrNoRows {
			c.JSON(400, gin.H{"message":"identity doesn't exist"})
			return
		}
		c.JSON(500, gin.H{"message":"identity doesn't exist"})
		return
	}
	query = "INSERT INTO record (\"identityNumber\", nip, symptoms, medications) VALUES ($1,$2,$3,$4)"
	res, err := conn.Exec(query, recordForm.IdentityNumber, nipInt, recordForm.Symptoms, recordForm.Medications)
	if err != nil {
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	if rows, _ :=res.RowsAffected();rows == 0 {
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	c.JSON(201, gin.H{"message":"record created successfully"})
}

// func (h RecordController) GetAllRecord(c *gin.Context){
// 	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
// 	if errLimit != nil || limit < 0 {
// 		limit = 5
// 	}
// 	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
// 	if errOffset != nil || offset < 0 {
// 		offset = 0
// 	}
// 	identityNumber := c.Query("identityNumber")
	
// 	identityInt, err := strconv.ParseInt(identityNumber, 10,64)
// 	if err != nil{
// 		identityNumber = ""
// 	}
// 	userId := c.Query("userId")
// 	nip := c.Query("nip")
// 	createdAt := c.Query("createdAt")
// 	nipInt, err := strconv.ParseInt(nip, 10, 64)
// 	if err!= nil {
// 		nip = ""
// 	}
// 	_, err = validator.ExtractNIP(nipInt)
// 	if err !=nil{
// 		nip = ""
// 	}
// 	if createdAt != "asc" && createdAt != "desc"{
// 		createdAt = ""
// 	}
// 	baseQuery := "SELECT * FROM record"
// 	var args []interface{}
// 	var queryParams []string
// 	argIdx := 1
// 	if userId != ""{
// 		queryParams = append(queryParams, " id::text = $"+strconv.Itoa(argIdx) +" ") 
// 		args = append(args, userId)
// 		argIdx += 1
// 	}
// }