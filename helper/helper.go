package helper

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

func CheckNIP(conn *sqlx.DB, nip int64) bool {
	query := "SELECT EXISTS (SELECT 1 FROM admin WHERE nip = $1 UNION SELECT 1 FROM nurse WHERE nip = $2) AS nip_exists"
	var isExist bool
	err := conn.QueryRow(query,nip,nip).Scan(&isExist)
	if err!= nil{
		fmt.Println(err.Error())
		return false
	}
	return isExist
}
func FormatToIso860(s string)string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		
		return ""
	}

	// Format the time object into ISO 8601 format
	return t.Format("2006-01-02T15:04:05Z07:00")
}
func Includes(target string, array []string)bool{
	for _, value := range array {
        if value == target {
            return true
        }
    }
    return false
}

