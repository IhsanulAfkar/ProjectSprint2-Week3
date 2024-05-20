package validator

import (
	"Week3/models"
	"errors"
	"regexp"
	"strconv"
	"time"
)

func StringCheck(input string, minLength, maxLength int) bool {
	// minLength += 1
	length := len(input)
	// fmt.Println(length >= minLength && length <= maxLength)
	return length >= minLength && length <= maxLength
}

func ExtractNIP(nip int64) (*models.NIP, error) {
	nipStr := strconv.FormatInt(nip, 10)
	
	// check if length below 13
	if len(nipStr) < 13 || len(nipStr) > 15 {
		return nil, errors.New("incorrect nip length")
	}
	// get first 3 digits
	first3digits := nipStr[:3]
	// if first3digits != "615"{
	// 	return nil, errors.New("incorrect nip format")
	// }
	// check gender
	genderDigit := nipStr[3]
	if genderDigit != '1' && genderDigit != '2'{
		return nil, errors.New("incorrect nip format")
	}
	yearDigits:= nipStr[4:8]
	yearInt, err := strconv.Atoi(yearDigits)
	if err != nil {
		return nil, err
	}
	currentYear := time.Now().Year()
	if yearInt < 2000 || yearInt > currentYear {
		return nil, errors.New("invalid nip format")
	}  
	monthDigits:= nipStr[8:10]
	monthInt, err := strconv.Atoi(monthDigits)
	if err != nil {
		return nil, err
	}

	if monthInt < 1 || monthInt > 12 {
		return nil, errors.New("invalid nip format")
	} 
	lastDigits := nipStr[10:]

	var gender string
	if genderDigit == '1'{
		gender = "male"
	} else {
		gender = "female"
	}
	return &models.NIP{
		ToString: nipStr,
		ToInt: nip,
		First3Digits: first3digits,
		Gender: gender,
		Year: yearDigits,
		Month: monthDigits,
		EndDigits: lastDigits,
	}, nil
}

func IsURL(s string) bool {

	regex := `https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,4}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)$`
	pattern := regexp.MustCompile(regex)
	return pattern.MatchString(s)
}
func IsDateISO860(s string)bool {
	_, err := time.Parse("2006-01-02T15:04:05Z07:00", s)
	return err == nil
}