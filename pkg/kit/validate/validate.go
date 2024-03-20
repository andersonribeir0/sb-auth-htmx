package validate

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"unicode"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,8}$`)
	urlRegex   = regexp.MustCompile(`^(http(s)?://)?([\da-z\.-]+)\.([a-z\.]{2,6})([/\w \.-]*)*/?$`)
)

type RuleFunc func() RuleSet

type RuleSet struct {
	Name         string
	RuleValue    any
	FieldValue   any
	FieldName    any
	MessageFunc  func(RuleSet) string
	ValidateFunc func(RuleSet) bool
}

type Fields map[string][]RuleSet

type Messages map[string]string

func Password() RuleSet {
	return RuleSet{
		Name: "password",
		//RuleValue: n,
		ValidateFunc: func(set RuleSet) bool {
			str, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			_, ok = ValidatePassword(str)
			return ok
		},
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s should be valid", set.FieldName)
		},
	}
}

func Required() RuleSet {
	return RuleSet{
		Name: "required",
		MessageFunc: func(set RuleSet) string {
			return fmt.Sprintf("%s is a required field", set.FieldName)
		},
		ValidateFunc: func(rule RuleSet) bool {
			str, ok := rule.FieldValue.(string)
			if !ok {
				return false
			}
			return len(str) > 0
		},
	}
}

func Message(msg string) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "message",
			RuleValue: msg,
		}
	}
}

func Url() RuleSet {
	return RuleSet{
		Name: "url",
		MessageFunc: func(set RuleSet) string {
			return "not a valid url"
		},
		ValidateFunc: func(set RuleSet) bool {
			u, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return urlRegex.MatchString(u)
		},
	}
}

func Email() RuleSet {
	return RuleSet{
		Name: "email",
		MessageFunc: func(set RuleSet) string {
			return "email address is invalid"
		},
		ValidateFunc: func(set RuleSet) bool {
			email, ok := set.FieldValue.(string)
			if !ok {
				return false
			}
			return emailRegex.MatchString(email)
		},
	}
}

func Equal(s string) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "max",
			RuleValue: s,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return str == s
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be equal %s", set.FieldName, s)
			},
		}
	}
}

func Max(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "max",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) <= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be maximum %d characters long", set.FieldName, n)
			},
		}
	}
}

func Min(n int) RuleFunc {
	return func() RuleSet {
		return RuleSet{
			Name:      "min",
			RuleValue: n,
			ValidateFunc: func(set RuleSet) bool {
				str, ok := set.FieldValue.(string)
				if !ok {
					return false
				}
				return len(str) >= n
			},
			MessageFunc: func(set RuleSet) string {
				return fmt.Sprintf("%s should be at least %d characters long", set.FieldName, n)
			},
		}
	}
}

func Rules(rules ...RuleFunc) []RuleSet {
	ruleSets := make([]RuleSet, len(rules))
	for i := 0; i < len(ruleSets); i++ {
		ruleSets[i] = rules[i]()
	}
	return ruleSets
}

type Validator struct {
	data   any
	fields Fields
}

func New(data any, fields Fields) *Validator {
	return &Validator{
		fields: fields,
		data:   data,
	}
}

func Validate(in any, out any, fields Fields) bool {
	return true
}

func (v *Validator) Validate(target any) bool {
	ok := true
	for fieldName, ruleSets := range v.fields {
		// reflect panics on un-exported variables.
		if !unicode.IsUpper(rune(fieldName[0])) {
			continue
		}
		fieldValue := getFieldValueByName(v.data, fieldName)
		for _, set := range ruleSets {
			set.FieldValue = fieldValue
			set.FieldName = fieldName
			if set.Name == "message" {
				setErrorMessage(target, fieldName, set.RuleValue.(string))
				continue
			}
			if !set.ValidateFunc(set) {
				msg := set.MessageFunc(set)
				setErrorMessage(target, fieldName, msg)
				ok = false
			}
		}
	}
	return ok
}

func setErrorMessage(v any, fieldName string, msg string) {
	if v == nil {
		return
	}
	switch t := v.(type) {
	case map[string]string:
		t[fieldName] = msg
	default:
		structVal := reflect.ValueOf(v)
		if structVal.Kind() != reflect.Ptr || structVal.IsNil() {
			return
		}
		structVal = structVal.Elem()
		field := structVal.FieldByName(fieldName)
		field.Set(reflect.ValueOf(msg))
	}
}

func getFieldValueByName(v any, name string) any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return nil
	}
	fieldVal := val.FieldByName(name)
	if !fieldVal.IsValid() {
		return nil
	}
	return fieldVal.Interface()
}

// validatePassword checks if the password is strong and meets the criteria:
// - At least 8 characters long
// - Contains at least one digit
// - Contains at least one lowercase letter
// - Contains at least one uppercase letter
// - Contains at least one special character
func ValidatePassword(password string) (string, bool) {
	var (
		hasUpper     = false
		hasLower     = false
		hasNumber    = false
		hasSpecial   = false
		specialRunes = "!@#$%^&*"
	)

	if len(password) < 8 {
		return "Password must contain at least 8 characters", false
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char) || strings.ContainsRune(specialRunes, char):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return "Password must contain at least 1 uppercase character", false
	}
	if !hasLower {
		return "Password must contain at least 1 lowercase character", false
	}
	if !hasNumber {
		return "Password must contain at least 1 numeric character (0, 1, 2, ...)", false
	}
	if !hasSpecial {
		return "Password must contain at least 1 special character (@, ;, _, ...)", false
	}
	return "", true
}
