package routes

import (
	"Week3/controllers"
	"Week3/db"
	"Week3/middleware"
	"Week3/models"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Init() *gin.Engine {
	adminController := new(controllers.AdminController)
	nurseController := new(controllers.NurseController)
	patientController := new(controllers.PatientController)
	recordController := new(controllers.RecordController)
	mediaController := new(controllers.MediaController)
	router := gin.New()
	router.GET("/", func(c *gin.Context) {
		conn := db.CreateConn()
		var listAdmin []models.Admin
		err := conn.Select(&listAdmin,"SELECT * FROM admin")
		if err != nil{
			fmt.Println(err.Error())
		}
		c.JSON(200,gin.H{"data":listAdmin})
	})
	v1 := router.Group("/v1")

	{
		user := v1.Group("/user")
		{
			admin := user.Group("/it")
			{
				admin.POST("/register", adminController.SignUp)
				admin.POST("/login", adminController.SignIn)
			}
			nurse := user.Group("/nurse")
			{
				nurse.POST("/login",nurseController.NurseLogin)
				nurse.Use(middleware.AdminAuthMiddleware)
				nurse.POST("/register", nurseController.CreateNurse)
				nurse.PUT("/:nurseId", nurseController.UpdateNurse)
				nurse.DELETE("/:nurseId", nurseController.DeleteNurse)
				nurse.POST("/:nurseId/access", nurseController.GrantAccess)
			}
			user.Use(middleware.AdminAuthMiddleware)
			user.GET("/", adminController.GetAllUsers)
		}
		medical := v1.Group("/medical")
		{
			medical.Use(middleware.AllAuthMiddleware)
			patient := medical.Group("/patient")
			{
				patient.POST("/",patientController.CreatePatient)
				patient.GET("/",patientController.GetAllPatient)
			}
			record := medical.Group("/record")
			{
				record.POST("/", recordController.CreateRecord)
			}
			v1.POST("/image",mediaController.UploadImage)
		}
	}
	return router
}