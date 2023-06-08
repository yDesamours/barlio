package validator

import (
	"regexp"
	"strings"
)

type Validator struct {
	error map[string][]string
}

func (v *Validator) Valid() bool {
	return len(v.error) == 0
}

func (v *Validator) Check(expr bool, field, message string) {
	if !expr {
		v.error[field] = append(v.error[field], message)
	}
}

func (v *Validator) Error() map[string][]string {
	return v.error
}

func (v *Validator) NotEmpty(str, field, message string) {
	v.Check(len(str) > 0, field, message)
}

func (v *Validator) IsEmailValid(str, field, message string) {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	v.Check(emailRegex.MatchString(str), field, message)
}

func (v *Validator) Equal(str1, str2, field, message string) {
	v.Check(strings.Compare(str1, str2) == 0, field, message)
}
