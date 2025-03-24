package utils

import (
	"bytes"
	"echo-fullstack/app/entity"
	"echo-fullstack/config"
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/labstack/echo/v4"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Res struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func TitleCase(str string) string {
	tc := cases.Title(language.Indonesian)
	return tc.String(str)
}

func ObjectToString(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func Respond(code int, data interface{}, message string) (res Res) {
	return Res{
		Status:  code,
		Message: message,
		Data:    data,
	}
}

func Average(numbers []float64) float64 {
	var sum float64
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

func GenerateRandomPIN() string {
	rand.Seed(time.Now().UnixNano())
	charset := "0123456789"
	randomBytes := make([]byte, 6)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(randomBytes)
	return randomString
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomBytes := make([]byte, length)
	for i := range randomBytes {
		randomBytes[i] = charset[rand.Intn(len(charset))]
	}
	randomString := string(randomBytes)
	return randomString
}

func StripTags(str string) string {
	str = strip.StripTags(str)
	return str
}

func StripTagsFromStruct(input interface{}) {
	structValue := reflect.ValueOf(input).Elem()

	for i := 0; i < structValue.NumField(); i++ {
		fieldValue := structValue.Field(i)

		if fieldValue.Kind() == reflect.String {
			originalValue := fieldValue.String()
			strippedValue := strip.StripTags(originalValue)
			fieldValue.SetString(strippedValue)
		} else if fieldValue.Kind() == reflect.Struct {
			StripTagsFromStruct(fieldValue.Addr().Interface())
		}
	}
}

func CompactJSON(data []byte) string {
	var js map[string]interface{}
	if json.Unmarshal(data, &js) != nil {
		return string(data)
	}

	result := new(bytes.Buffer)
	if err := json.Compact(result, data); err != nil {
		return ""
	}
	return result.String()
}

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = charset[rand.Intn(len(charset))]
	}
	return string(bytes)
}

func PopulatePaging(c echo.Context, custom string) (param entity.ReqPaging) {
	customval := c.QueryParam(custom)
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit == 0 {
		limit = 10
	}
	if limit >= 50 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	page, _ := strconv.Atoi(c.QueryParam("page"))
	if page == 0 && offset == 0 {
		page = 1
		offset = 0
	}
	if page >= 1 && offset == 0 {
		offset = (page - 1) * limit
	}
	draw, _ := strconv.Atoi(c.QueryParam("draw"))
	if draw == 0 {
		draw = 1
	}
	order := c.QueryParam("sort")
	if strings.ToLower(order) == "asc" {
		order = "ASC"
	} else {
		order = "DESC"
	}
	sort := c.QueryParam("order")
	if sort == "" {
		sort = "id"
	}
	param = entity.ReqPaging{
		Search: c.QueryParam("search"),
		Order:  order,
		Limit:  limit,
		Offset: offset,
		Sort:   sort,
		Custom: customval,
		Page:   page}
	return
}

func PopulateResPaging(param *entity.ReqPaging, data interface{}, totalResult int64) (output entity.ResPaging) {
	totalPages := int(totalResult) / param.Limit
	if int(totalResult)%param.Limit > 0 {
		totalPages++
	}

	currentPage := param.Offset/param.Limit + 1
	next := false
	back := false
	if currentPage < totalPages {
		next = true
	}
	if currentPage <= totalPages && currentPage > 1 {
		back = true
	}
	output = entity.ResPaging{
		Status:          200,
		Draw:            1,
		Data:            data,
		Search:          param.Search,
		Order:           param.Order,
		Limit:           param.Limit,
		Offset:          param.Offset,
		Sort:            param.Sort,
		Next:            next,
		Back:            back,
		TotalData:       int(totalResult),
		RecordsFiltered: int(totalResult),
		CurrentPage:     currentPage,
		TotalPage:       totalPages,
	}
	return
}

func ConvertToCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := range parts {
		parts[i] = TitleCase(parts[i])
	}
	return strings.Join(parts, "")
}

func IsStringInArray(str string, array []string) bool {
	for _, s := range array {
		if s == str {
			return true
		}
	}
	return false
}

func LastId(table string) (id int64) {
	type OnlyId struct {
		ID int64
	}
	var last OnlyId
	config.DB.Table(table).Order("id desc").Limit(1).Scan(&last)

	id = last.ID + 1
	return
}

type FavoriteResponse struct {
	Response string `gorm:"column:response"`
}

func FavoriteValueOfColomn(table string, column string, where string) (response FavoriteResponse, err error) {
	err = config.DB.Table(table).
		Select(column + " AS response").
		Where(where).
		Group(column).
		Order("COUNT(*) DESC").
		Limit(1).
		Scan(&response).Error

	return response, err
}

func GenerateKeyStruct(data interface{}) string {
	value := reflect.ValueOf(data)
	key := ""
	for i := 0; i < value.NumField(); i++ {
		fieldValue := value.Field(i).Interface()
		key += fmt.Sprintf(":%v", fieldValue)
	}
	return key
}

func IsStringNotEmpty(value interface{}) bool {
	if str, ok := value.(string); ok && str != "" {
		return true
	}
	return false
}

func ConvertToSnakeCase(input string) string {
	var result []rune

	input = strings.ToLower(input)

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
		} else if len(result) > 0 && result[len(result)-1] != '_' {
			result = append(result, '_')
		}
	}

	return string(result)
}

func ConvertToKebabCase(input string) string {
	var result []rune

	input = strings.ToLower(input)

	for _, r := range input {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, r)
		} else if len(result) > 0 && result[len(result)-1] != '-' {
			result = append(result, '-')
		}
	}

	return string(result)
}

func ContainsString(s, substr string) bool {
	return strings.Contains(s, substr)
}

func MakeKey(values ...interface{}) string {
	var keyParts []string

	for _, value := range values {
		switch v := value.(type) {
		case string:
			keyParts = append(keyParts, v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			keyParts = append(keyParts, fmt.Sprintf("%v", v))
		case bool:
			keyParts = append(keyParts, fmt.Sprintf("%t", v))
		case []string:
			keyParts = append(keyParts, value.([]string)...)
		default:
			val := reflect.ValueOf(value)
			if val.Kind() == reflect.Struct {
				for i := 0; i < val.NumField(); i++ {
					field := val.Field(i)
					keyParts = append(keyParts, fmt.Sprintf("%v", field.Interface()))
				}
			}
		}
	}

	return strings.Join(keyParts, "_")
}
