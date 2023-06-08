package validator

type Validator struct {
	error map[string][]string
}

func (v *Validator) Valid() bool {
	return len(v.error) == 0
}

func (v *Validator) Check(expr bool, field string, message string) {
	if !expr {
		v.error[field] = append(v.error[field], message)
	}
}

func (v *Validator) Error() map[string][]string {
	return v.error
}

func (v *Validator) NotEmpty(str string, field string, message string) {
	v.Check(len(str) > 0, field, message)
}
