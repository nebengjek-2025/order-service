package utils

import (
	"crypto/md5"
	"crypto/rand"
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

var entityPrefixes = map[string]string{
	"user":    "USR",
	"driver":  "DRV",
	"trip":    "TRP",
	"order":   "ORD",
	"wallet":  "WLT",
	"payment": "PAY",
}

// ConvertString to convert any data type to String
func ConvertString(v interface{}) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case bool:
		return strconv.FormatBool(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case []uint8:
		return string(val)
	default:
		resultJSON, err := json.Marshal(v)
		if err != nil {
			log.Println("Error on lib/converter ConvertString() ", err)
			return ""
		}
		return string(resultJSON)
	}
}

// ConvertInt to convert any date type to Int
func ConvertInt(v interface{}) int {
	switch val := v.(type) {
	case string:
		str := strings.TrimSpace(val)
		result, _ := strconv.Atoi(str)
		return result

	case int:
		return val

	case int64:
		return int(val)

	case float64:
		return int(val)

	case []byte:
		result, _ := strconv.Atoi(string(val))
		return result

	default:
		return 0
	}
}

// ConvertInt64 to convert any date type to Int64
func ConvertInt64(v interface{}) int64 {
	switch val := v.(type) {
	case string:
		val = strings.TrimSpace(val)
		result, _ := strconv.ParseInt(val, 10, 64)
		return result

	case int:
		return int64(val)

	case int64:
		return val

	case float64:
		return int64(val)

	case []byte:
		result, _ := strconv.ParseInt(string(val), 10, 64)
		return result

	default:
		return 0
	}
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

func GenerateUniqueIDWithPrefix(entityType string) string {
	entityTypeLower := strings.ToLower(entityType)
	prefix, ok := entityPrefixes[entityTypeLower]
	if !ok {
		prefix = "GEN"
	}

	timestamp := time.Now().UTC().Format("20060102T150405")
	randomBytes := make([]byte, 4)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomHex := strings.ToUpper(hex.EncodeToString(randomBytes))

	return fmt.Sprintf("NBJ_%s_%s_%s", prefix, timestamp, randomHex)
}
