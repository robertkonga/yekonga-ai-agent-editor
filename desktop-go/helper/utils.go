package helper

import (
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
	"yekonga-builder/console"
)

func IsMap(data interface{}) bool {
	val := reflect.ValueOf(data)

	// Dereference pointers to check the underlying value
	if val.Kind() == reflect.Ptr {
		if !val.IsNil() {
			val = val.Elem()
		}
	}

	if IsNotEmpty(data) {
		ok := val.Kind() == reflect.Map

		// Type assertion to check if data is a map[string]interface{}
		if ok {
			return true
		}
	}

	return false
}

func IsMapList(data interface{}) bool {
	val := reflect.ValueOf(data)

	// Dereference pointers to check the underlying value
	if val.Kind() == reflect.Ptr {
		if !val.IsNil() {
			val = val.Elem()
		}
	}

	if IsNotEmpty(data) {
		ok := val.Kind() == reflect.Array || val.Kind() == reflect.Slice

		// Type assertion to check if data is a map[string]interface{}
		if ok {
			return true
		}
	}

	return false
}

// FileExists checks if a file exists
func FileExists(filename string) bool {
	_, err := os.Stat(filename)

	if err != nil {
		// console.Error("FileExists", err)
		return false
	}

	return !os.IsNotExist(err)
}

// IsNumeric checks if a value is an int, float, or a numeric string
func IsNumeric(value interface{}) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return true // Directly numeric types

	case string:
		_, err := strconv.ParseFloat(v, 64)
		return err == nil // Returns true if string is a valid number

	default:
		return false // Not numeric
	}
}

func ToInt(value interface{}) int {
	number := 0
	if v, ok := value.(string); ok {
		n, err := strconv.Atoi(v)
		if err == nil {
			number = n
		}
	} else if v, ok := value.(int); ok {
		number = v
	} else if IsNumeric(value) {
		n, err := strconv.ParseInt(fmt.Sprintf("%v", value), 32, 64)
		if err == nil {
			number = int(n)
		}
	}

	return number
}

func ToFloat(value interface{}) float64 {
	// console.Log("ToFloat", fmt.Sprintf("%v", value))

	var number float64 = 0
	if v, ok := value.(string); ok {
		n, err := strconv.ParseFloat(v, 64)
		if err == nil {
			number = n
		}
	} else if v, ok := value.(float64); ok {
		number = v
	} else if v, ok := value.(int); ok {
		n, err := strconv.ParseFloat(strconv.Itoa(v), 64)
		if err == nil {
			number = n
		}
	} else if IsNumeric(value) {
		n, err := strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
		if err == nil {
			number = n
		}
	}

	// console.Log("ToFloat", GetType(number), number)

	return number
}

func ToFloat64(value interface{}) float64 {
	return ToFloat(value)
}

func CompareValues(a, b interface{}) int {
	// Convert both values to float64 for comparison
	aVal := ToFloat64(a)
	bVal := ToFloat64(b)

	if aVal < bVal {
		return -1
	} else if aVal > bVal {
		return 1
	}
	return 0
}

func ToJsonFormatted(data interface{}) string {
	jsonData, _ := json.MarshalIndent(data, "", "    ")

	return string(jsonData)
}

func ToJson(data interface{}) string {
	jsonData, _ := json.Marshal(data)

	return string(jsonData)
}

func ToByte(data interface{}) []byte {
	jsonData, _ := json.Marshal(data)

	return jsonData
}

// JSON file and converts it to a map
func ToMap[T any](data interface{}) map[string]T {
	var result map[string]T
	var dataByte []byte
	if s, ok := data.(string); ok {
		dataByte = []byte(s)
	} else {
		dataByte = ToByte(data)
	}

	if err := json.Unmarshal(dataByte, &result); err != nil {
		return nil
	}

	return result
}

func ToMapList[T any](data interface{}) []map[string]T {
	converted := []map[string]T{}
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		count := val.Len()
		converted = make([]map[string]T, count)

		for i := 0; i < count; i++ {
			elem := ToMap[T](val.Index(i).Interface())
			converted[i] = elem
		}
	}

	return converted
}

func ToList[T any](data interface{}) []T {
	converted := []T{}
	val := reflect.ValueOf(data)

	if val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
		count := val.Len()
		converted = make([]T, 0, count)

		for i := 0; i < count; i++ {
			var result T
			if err := json.Unmarshal(ToByte(val.Index(i).Interface()), &result); err != nil {
				console.Error("ToList", err.Error())
			} else {
				converted = append(converted, result)
			}
		}

	}

	return converted
}

func ToInterface(data interface{}) (interface{}, error) {
	var result interface{}

	// JSON file and converts it to a map
	if err := json.Unmarshal(ToByte(data), &result); err != nil {
		console.Error("ConvertTo", err.Error())
		return result, err
	}

	return result, nil
}

func IsSlice(v interface{}) bool {
	if v == nil {
		return false
	}

	if IsPointer(v) {
		return IsSlice(reflect.ValueOf(v).Elem().Interface())
	}

	t := reflect.TypeOf(v)

	kind := t.Kind()
	return kind == reflect.Array || kind == reflect.Slice
}

func IsList(v interface{}) bool {
	return IsSlice(v)
}

func IsArray(v interface{}) bool {
	return IsSlice(v)
}

// convertTo converts a map[string]interface{} to a struct of type T
func ConvertTo[T any](data interface{}) (T, error) {
	var result T

	// JSON file and converts it to a map
	if err := json.Unmarshal(ToByte(data), &result); err != nil {
		console.Error("ConvertTo", err.Error())
		return result, err
	}

	return result, nil
}

// setField sets the reflect.Value of a struct field to the provided value
func setField(field reflect.Value, value interface{}) error {
	// Convert the value to the field's type
	val := reflect.ValueOf(value)

	// Check if the value can be converted to the field's type
	if !val.Type().ConvertibleTo(field.Type()) {
		return fmt.Errorf("cannot convert %v to %v", val.Type(), field.Type())
	}

	// Perform the conversion and set the field
	field.Set(val.Convert(field.Type()))
	return nil
}

func IsPointer(v interface{}) bool {
	if v == nil {
		return false
	}

	return reflect.TypeOf(v).Kind() == reflect.Ptr
}

func CreateFile(data interface{}, filename string) error {
	return SaveToFile(data, filename)
}

func SaveToFile(data interface{}, filename string) error {
	var (
		rowData []byte
		err     error
	)
	// Convert to JSON
	if d, ok := data.(string); ok {
		rowData = []byte(d)
	} else {
		rowData, err = json.MarshalIndent(data, "", "  ") // Pretty print with indentation
		if err != nil {
			return err
		}
	}

	// Extract directory path
	dir := filepath.Dir(filename)

	// Create all folders if they don't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write to file
	err = os.WriteFile(filename, rowData, 0644) // 0644 is standard file permission
	if err != nil {
		return err
	}

	return nil
}

func CreateDirectory(dir string) error {
	// Create all folders if they don't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return nil
}

func CreateFolder(dir string) error {
	return CreateDirectory(dir)
}

// IsEmpty checks if a value is empty/zero for various types
func IsEmpty(v interface{}) bool {
	if v == nil {
		return true
	}

	val := reflect.ValueOf(v)

	// Dereference pointers to check the underlying value
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return true
		}
		// Get the element the pointer points to
		val = val.Elem()
	}

	switch val.Kind() {
	case reflect.String:
		return val.Len() == 0

	case reflect.Bool:
		return !val.Bool()

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return val.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return val.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return val.Complex() == 0

	case reflect.Interface:
		if val.IsNil() {
			return true
		}
		return IsEmpty(val.Interface())

	case reflect.Array, reflect.Slice:
		return val.Len() == 0

	case reflect.Map:
		return val.Len() == 0

	case reflect.Chan:
		return val.IsNil()

	case reflect.Struct:
		// Check if it's time.Time specifically
		if t, ok := val.Interface().(time.Time); ok {
			return t.IsZero()
		}
		// For other structs, compare with zero value
		return val.IsZero()

	case reflect.Func:
		return val.IsNil()

	default:
		// Use IsZero for any other types (available in Go 1.13+)
		return val.IsZero()
	}
}

func IsNotEmpty(value interface{}) bool {
	return !IsEmpty(value)
}

// CreateFuzzyRegex converts "R F Konga" into "(?i)R.*F.*Konga"
func CreateFuzzyRegex(input string) string {
	// 1. Split by whitespace
	parts := strings.Fields(input)

	// 2. Escape special characters in each part to prevent injection
	for i, part := range parts {
		parts[i] = regexp.QuoteMeta(part)
	}

	// 3. Join with .* and add the case-insensitive flag (?i)
	pattern := "(?i)" + strings.Join(parts, ".*")
	return pattern
}

// Utility function to check if slice contains an element
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Reverse[T interface{}](slice []T) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func SortMap[T interface{}](options map[string]T) map[string]T {
	// 1. Get all keys into a slice
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}

	// 2. Sort the keys based on the map values
	slices.SortFunc(keys, func(a, b string) int {
		return cmp.Compare(ToString(options[a]), ToString(options[b]))
	})

	// 3. Print the results in order
	fmt.Println("Sorted by Value:", keys)
	newOptions := make(map[string]T)
	for _, k := range keys {
		fmt.Printf("%s: %s\n", k, options[k])
		newOptions[k] = options[k]
	}

	return newOptions
}

func SortedKeys[T interface{}](options map[string]T) []string {
	// 1. Get all keys into a slice
	keys := make([]string, 0, len(options))
	for k := range options {
		keys = append(keys, k)
	}

	// 2. Sort the keys based on the map values
	slices.SortFunc(keys, func(a, b string) int {
		return cmp.Compare(ToString(options[a]), ToString(options[b]))
	})

	return keys
}

// ToCamelCase converts a string to CamelCase
func ToCamelCase(s string) string {
	s = ToUnderscore(s)
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsNumber(r)
	})

	for i := range words {
		words[i] = strings.Title(strings.ToLower(words[i]))
	}

	return strings.Join(words, "")
}

// ToSlug converts a string to a URL-friendly slug
func ToSlug(s string) string {
	// Convert to lowercase
	s = ToUnderscore(s)

	// Replace non-alphanumeric characters with hyphens
	re := regexp.MustCompile(`[^a-z0-9]+`)
	s = re.ReplaceAllString(s, "-")

	// Trim leading and trailing hyphens
	s = strings.Trim(s, "-")

	return s
}

func ToString(s interface{}) string {
	return fmt.Sprintf("%v", s)
}

// camelToSnake converts camelCase or PascalCase to snake_case
func CamelToSnake(s string) string {
	var (
		result    []rune
		prevLower bool
		prevUpper bool
		prevDigit bool
	)
	for _, r := range s {

		isLower := unicode.IsLower(r)
		isUpper := unicode.IsUpper(r)
		isDigit := unicode.IsDigit(r)

		if (prevLower && isUpper) ||
			(prevDigit && (isLower || isUpper)) ||
			(isDigit && (prevLower || prevUpper)) {
			result = append(result, '_')
		}

		// if isUpper {
		// 	// Add underscore before uppercase letter (except first char)
		// 	result = append(result, '_')
		// }
		result = append(result, unicode.ToLower(r))

		prevLower = isLower
		prevUpper = isUpper
		prevDigit = isDigit
	}

	if len(result) == 0 {
		_ = fmt.Sprint(prevDigit, prevLower, prevUpper)
	}

	return string(result)
}

// ToUnderscore converts a string into snake_case format.
// It handles camelCase, PascalCase, and kebab-case.
func ToUnderscore(text string) string {
	if text == "" {
		return ""
	}

	// 1. Insert underscore before capital letters (camelCase/PascalCase support)
	t := CamelToSnake(text)

	// 2. Convert to lowercase
	t = strings.ToLower(t)

	// 3. Replace spaces and hyphens with a single underscore
	// Example: "hello-world thing" -> "hello_world_thing"
	reSeparator := regexp.MustCompile("[\\s-]+")
	t = reSeparator.ReplaceAllString(t, "_")

	// 4. Remove any multiple consecutive underscores resulting from the previous steps
	reMultipleUnderscores := regexp.MustCompile("_+")
	t = reMultipleUnderscores.ReplaceAllString(t, "_")

	// 5. Remove leading and trailing underscores
	t = strings.Trim(t, "_")

	return t
}

func ToVariable(s string) string {
	s = ToCamelCase(s)

	s = strings.ToLower(s[0:1]) + s[1:]

	return s
}

func ToTitle(s string) string {
	if s == "" {
		return ""
	}

	// Step 1: camelCase → snake_case
	s = ToUnderscore(s)

	// Step 2: collapse multiple underscores → single space
	re := regexp.MustCompile(`_+`)
	s = re.ReplaceAllString(s, " ")

	// Step 3: trim and split into words
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}

	words := strings.Split(s, " ")
	var titleWords []string

	for _, word := range words {
		if word == "" {
			continue
		}
		// Title case: first rune uppercase, rest unchanged
		runes := []rune(word)
		runes[0] = unicode.ToUpper(runes[0])
		titleWords = append(titleWords, string(runes))
	}

	return strings.Join(titleWords, " ")
}

// Pluralization rules
var pluralRules = []struct {
	pattern *regexp.Regexp
	replace string
}{
	{regexp.MustCompile("([^aeiou])y$"), "${1}ies"},  // e.g., city → cities
	{regexp.MustCompile("(f|fe)$"), "ves"},           // e.g., knife → knives
	{regexp.MustCompile("(s|sh|ch|x|z)$"), "${1}es"}, // e.g., bus → buses, box → boxes
	{regexp.MustCompile("$"), "s"},                   // Default rule: add "s"
}

// Singularization rules
var singularRules = []struct {
	pattern *regexp.Regexp
	replace string
}{
	{regexp.MustCompile("ies$"), "y"},                // e.g., cities → city
	{regexp.MustCompile("ves$"), "f"},                // e.g., knives → knife
	{regexp.MustCompile("(s|sh|ch|x|z)es$"), "${1}"}, // e.g., boxes → box
	{regexp.MustCompile("s$"), ""},                   // Default rule: remove "s"
}

// Pluralize converts a singular noun to its plural form
func Pluralize(word string) string {
	word = Singularize(word)

	for _, rule := range pluralRules {
		if rule.pattern.MatchString(word) {
			return rule.pattern.ReplaceAllString(word, rule.replace)
		}
	}
	return word
}

// Singularize converts a plural noun to its singular form
func Singularize(word string) string {
	for _, rule := range singularRules {
		if rule.pattern.MatchString(word) {
			return rule.pattern.ReplaceAllString(word, rule.replace)
		}
	}
	return word
}

func ReadFile(filename string) string {
	return LoadFile(filename)
}

// LoadFile reads a JSON file and converts it to a map
func LoadFile(filename string) string {
	filename = GetPath(filename)

	file, err := os.Open(filename)
	if err != nil {
		return ""
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return ""
	}

	return string(bytes)
}

// LoadJSONFile reads a JSON file and converts it to a map
func LoadJSONFile(filename string) (map[string]interface{}, error) {
	filename = GetPath(filename)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(bytes, &data); err != nil {
		return nil, err
	}

	return data, nil
}

func GetPath(relativePath string) string {
	// 1. Get the path of the executable
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	// 2. Get the directory of the executable
	exPath := filepath.Dir(ex)

	// 3. Join the executable's directory with the relative path
	absolutePath := filepath.Join(exPath, relativePath)

	if FileExists(absolutePath) {
		return absolutePath
	}

	if filepath.IsAbs(relativePath) {
		if FileExists(relativePath) {
			return relativePath
		}
	}

	if FileExists(relativePath) {
		absPath, err := filepath.Abs(relativePath)
		if err != nil {
			return relativePath
		}
		return absPath
	}

	return absolutePath
}

func GetValueOfString(data interface{}, key string) string {
	return GetMapString(data, key)
}

func GetMapString(data interface{}, key string) string {
	if value, ok := GetMapValue(data, key).(string); ok {
		return value
	}

	return ""
}

func GetValueOfInt(data interface{}, key string) int {
	return GetMapInt(data, key)
}

func GetMapInt(data interface{}, key string) int {
	v := GetMapValue(data, key)

	if value, ok := v.(int); ok {
		return value
	} else if IsNumeric(v) {
		return ToInt(v)
	}

	return 0
}

func GetValueOfFloat(data interface{}, key string) float64 {
	return GetMapFloat(data, key)
}

func GetMapFloat(data interface{}, key string) float64 {
	v := GetMapValue(data, key)
	if value, ok := v.(int); ok {
		return ToFloat(value)
	} else if value, ok := v.(float64); ok {
		return value
	} else if IsNumeric(v) {
		return ToFloat64(v)
	}

	return 0
}

func GetValueOfBoolean(data interface{}, key string) bool {
	return GetMapBoolean(data, key)

}

func GetMapBoolean(data interface{}, key string) bool {
	if value, ok := GetMapValue(data, key).(bool); ok {
		return value
	}

	return false
}

func GetValueOfDate(data interface{}, key string) time.Time {
	return GetMapDate(data, key)
}

func GetMapDate(data interface{}, key string) time.Time {
	value := GetMapValue(data, key)
	// console.Log("value", GetType(value), value)

	return GetTimestamp(value)
}

func GetValueOfMap(data interface{}, key string) map[string]interface{} {
	return GetMap(data, key)
}

func GetMap(data interface{}, key string) map[string]interface{} {
	v := GetMapValue(data, key)

	if IsNotEmpty(v) {
		return ToMap[interface{}](v)
	}

	return nil
}

func GetValueOf(data interface{}, key string) interface{} {
	return GetMapValue(data, key)
}

// func getMapValueItem() {
// 	if vi, oki := v[first]; oki {
// 		if len(keys[1:]) == 0 {
// 			return vi
// 		} else {
// 			last := strings.Join(keys[1:], ".")
// 			return GetMapValue(vi, last)
// 		}
// 	}
// }

func GetMapValue(data interface{}, key string) interface{} {
	if str, ok := data.(string); ok {
		var d map[string]interface{}
		if err := json.Unmarshal([]byte(str), &d); err == nil {
			data = d
		} else {
			console.Info(d, err.Error())
			console.Info(data)
		}
	}

	keys := strings.Split(key, ".")
	first := keys[0]
	var localData interface{}

	if IsNotEmpty(data) {
		if IsPointer(data) {
			v := reflect.ValueOf(data).Elem()
			if v.IsValid() {
				localData = v.Interface()
			}
		} else {
			localData = data
		}
	}

	if v, ok := localData.([]interface{}); ok {
		if IsNumeric(first) {
			pos := ToInt(first)

			if vi := v[pos]; vi != nil {
				if len(keys[1:]) == 0 {
					return vi
				} else {
					last := strings.Join(keys[1:], ".")
					return GetMapValue(vi, last)
				}
			}
		}
	}

	return nil
}

func GetMapArray(data interface{}, source string, keys map[string]string) []interface{} {
	list := []interface{}{}
	dataList := GetMapValue(data, source)

	if v, ok := dataList.([]interface{}); ok {
		for _, vi := range v {
			data := map[string]interface{}{}
			for kii, vii := range keys {
				data[kii] = GetMapValue(vi, vii)
			}

			list = append(list, data)
		}
	}

	return list
}

func GetTimestamp(value interface{}) time.Time {
	result := StringToDatetime(value)

	if result != nil {
		return (*result).UTC()
	}

	return time.Now().UTC()
}

func GetTimestampString(value interface{}) string {
	return GetTimestamp(value).Format(time.RFC3339)
}

func ToTimestampString(value interface{}, layout string) time.Time {
	if IsEmpty(layout) {
		layout = "2006"
	}
	if v, ok := value.(string); ok {
		parsedTime, _ := time.Parse(layout, v)

		return parsedTime.UTC()
	} else if v, ok := value.(time.Time); ok {
		return v
	}

	return time.Now().UTC()
}

func StringToDatetime(value interface{}) *time.Time {
	if strValue, ok := value.(string); ok {
		var t time.Time
		var err error

		t, err = time.Parse(time.DateOnly, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.DateTime, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.UnixDate, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.RFC3339, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.RFC822, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.TimeOnly, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.RFC850, strValue)
		if err == nil {
			return &t
		}

		t, err = time.Parse(time.UnixDate, strValue)
		if err == nil {
			return &t
		}

		ISO_8601 := "2006-01-02T15:04:05Z07:00"
		RFC_1123Z := "Tue Dec 30 2025 11:00:59 GMT+0300 (East Africa Time)"

		t, err = time.Parse(ISO_8601, strValue)
		if err == nil {
			return &t
		}
		// 		console.Log("ISO_8601", strValue, t, err)
		t, err = time.Parse(RFC_1123Z, strValue)
		if err == nil {
			return &t
		}

		// Try custom date parsing as a last resort
		t, err = DateParse(strValue, GetTimestamp(nil))
		// console.Log("DateParse", strValue, t, err)
		if err == nil {
			return &t
		}
	} else if strValue, ok := value.(time.Time); ok {
		return &strValue
	}

	return nil
}
