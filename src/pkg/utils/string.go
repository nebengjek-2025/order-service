package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ConvertString to convert any data type to String
func ConvertString(v interface{}) string {
	result := ""
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		result = v.(string)
	case int:
		result = strconv.Itoa(v.(int))
	case int64:
		result = strconv.FormatInt(v.(int64), 10)
	case bool:
		result = strconv.FormatBool(v.(bool))
	case float64:
		result = strconv.FormatFloat(v.(float64), 'f', -1, 64)
	case []uint8:
		result = string(v.([]uint8))
	default:
		resultJSON, err := json.Marshal(v)
		if err == nil {
			result = string(resultJSON)
		} else {
			log.Println("Error on lib/converter ConvertString() ", err)
		}
	}

	return result
}

// ConvertInt to convert any date type to Int
func ConvertInt(v interface{}) int {
	result := int(0)
	switch v.(type) {
	case string:
		str := strings.TrimSpace(v.(string))
		result, _ = strconv.Atoi(str)
	case int:
		result = int(v.(int))
	case int64:
		result = int(v.(int64))
	case float64:
		result = int(v.(float64))
	case []byte:
		result, _ = strconv.Atoi(string(v.([]byte)))
	default:
		result = int(0)
	}

	return result
}

// ConvertInt64 to convert any date type to Int64
func ConvertInt64(v interface{}) int64 {
	result := int64(0)
	switch v.(type) {
	case string:
		str := strings.TrimSpace(v.(string))
		result, _ = strconv.ParseInt(str, 10, 64)
	case int:
		result = int64(v.(int))
	case int64:
		result = int64(v.(int64))
	case float64:
		result = int64(v.(float64))
	case []byte:
		result, _ = strconv.ParseInt(string(v.([]byte)), 10, 64)
	default:
		result = int64(0)
	}

	return result
}

// GetLocalTime to retrieve current local time
func GetLocalTime() time.Time {
	return time.Now().Local()
}

func GenerateUUID() uuid.UUID {
	// Generate Random uuID
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println(fmt.Errorf("failed to generate UUID: %w", err))
		return uuid.Nil
	}

	return id
}

func ConvertStringUuid(v string) uuid.UUID {
	return uuid.MustParse(v)
}

func GenerateToken(email string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(email), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	hasher := md5.New()
	hasher.Write(hash)
	return hex.EncodeToString(hasher.Sum(nil))
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func FormatPrice(price float64) string {
	formatted := fmt.Sprintf("%.2f", price)
	parts := strings.Split(formatted, ".")

	integerPart := parts[0]
	decimalPart := parts[1]
	integerPartWithSeparator := ""

	for i := len(integerPart); i > 0; i -= 3 {
		if i-3 >= 0 {
			integerPartWithSeparator = "." + integerPart[i-3:i] + integerPartWithSeparator
		} else {
			integerPartWithSeparator = integerPart[:i] + integerPartWithSeparator
		}
	}

	return "Rp " + integerPartWithSeparator[1:] + "," + decimalPart
}

func FormatDuration(minutes int) string {
	if minutes < 60 {
		return fmt.Sprintf("%d menit", minutes)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60

	if remainingMinutes > 0 {
		return fmt.Sprintf("%d jam %d menit", hours, remainingMinutes)
	}

	return fmt.Sprintf("%d jam", hours)
}
