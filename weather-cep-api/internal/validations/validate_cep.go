package validations

import "regexp"


func IsZipcodeValid(zipcode string) (bool,string) {
	re := regexp.MustCompile(`\D`)
	zipcode = re.ReplaceAllString(zipcode, "")

	return len(zipcode) == 8,zipcode
}
