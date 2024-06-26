package controllers

import (
	"Week3/db"
	"Week3/forms"
	"Week3/helper/validator"
	"Week3/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type RecordController struct{}

func (h RecordController) CreateRecord(c *gin.Context){
	userNip := fmt.Sprintf("%s", c.MustGet("userNip"))
	nipInt, err := strconv.ParseInt(userNip, 10, 64)
	if err != nil {
		c.JSON(400, gin.H{
			"message":err.Error()})
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
			c.JSON(404, gin.H{"message":"identity doesn't exist"})
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

func (h RecordController) GetAllRecord(c *gin.Context){
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if errLimit != nil || limit < 0 {
		limit = 5
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil || offset < 0 {
		offset = 0
	}
	identityNumber := c.Query("identityDetail.identityNumber")
	
	identityInt, err := strconv.ParseInt(identityNumber, 10,64)
	
	if err != nil{
		identityNumber = ""
	}
	userId := c.Query("createdBy.userId")
	nip := c.Query("createdBy.nip")
	createdAt := c.Query("createdAt")
	nipInt, err := strconv.ParseInt(nip, 10, 64)
	if err!= nil {
		nip = ""
	}
	_, err = validator.ExtractNIP(nipInt)
	if err !=nil{
		nip = ""
	}
	if createdAt != "asc" && createdAt != "desc"{
		createdAt = ""
	}
	baseQuery := `select record."identityNumber" as "identityNumber", patient."phoneNumber" as "phoneNumber", patient.name as "identityName",
	patient."birthDate" as "birthDate", patient.gender as gender,
	patient."identityCardScanImg" as "identityCardScanImg", 
	record.symptoms, record.medications, record."createdAt", public.user.nip as nip, public.user.name as "userName", public.user.id as "userId"
	from record join public.user on record.nip = public.user.nip join patient on record."identityNumber" = patient."identityNumber"`
	var args []interface{}
	var queryParams []string
	argIdx := 1
	if userId != ""{
		queryParams = append(queryParams, " public.user.id::text = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, userId)
		argIdx += 1
	}
	if nip != "" {
		queryParams = append(queryParams, " public.user.nip = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, nipInt)
		argIdx += 1
	}
	if identityNumber != ""{
		queryParams = append(queryParams, " record.\"identityNumber\" = $"+strconv.Itoa(argIdx) +" ") 
		args = append(args, identityInt)
		argIdx += 1
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
	 
	records := make([]models.GetRecord,0)
	rows, err := conn.Query(baseQuery, args...)
	if err != nil {
		 
		c.JSON(500, gin.H{"message":"server error"})
		return
	}
	defer rows.Close()
	for rows.Next(){
		var record models.GetRecord 
		err = rows.Scan(&record.IdentityDetail.IdentityNumber, &record.IdentityDetail.PhoneNumber, &record.IdentityDetail.Name, &record.IdentityDetail.BirthDate, &record.IdentityDetail.Gender, &record.IdentityDetail.IdentityCardScanImg, &record.Symptoms, &record.Medications, &record.CreatedAt, &record.CreatedBy.Nip, &record.CreatedBy.Name, &record.CreatedBy.UserId)
		if err != nil {
			 
			c.JSON(500, gin.H{"message":"server error"})
			return
		}
		records = append(records, record)
	}
	// rapiin data
	
	c.JSON(200, gin.H{"message":"success","data":records})
}	