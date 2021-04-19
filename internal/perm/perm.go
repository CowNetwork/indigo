package perm

import (
	"regexp"
	"strings"
)

// permRegex matches any string dot seperated with
// optional `*` wildcards. Example: `[-]this.is.*.permission`
var permRegex = regexp.MustCompile("^-?(?:\\w+|\\*)(?:\\.(\\w+|\\*))*$")

func ValidatePermission(perm string) bool {
	return permRegex.MatchString(perm)
}

func GetPermissionRegexString(perm string) string {
	if !ValidatePermission(perm) {
		return perm
	}
	if strings.HasSuffix(perm, ".*") {
		perm = perm[:len(perm)-2]
		return GetPermissionRegexString(perm) + "(?:\\.(\\w+|\\*))+"
	}
	// replace wildcard with regex
	perm = strings.ReplaceAll(perm, "*", "(?:\\w+|\\*)")

	// escape the dots
	return strings.ReplaceAll(perm, ".", "\\.")
}

func GetPermissionRegex(perm string) (*regexp.Regexp, error) {
	return regexp.Compile(perm)
}

type Validator struct {
	perms  []string
	regexs map[string]*regexp.Regexp
}

func NewValidator(perms []string) *Validator {
	// TODO sort perms by range and negative permissions
	// TODO compact perms when permission A overlap/negate B
	regexs := make(map[string]*regexp.Regexp, len(perms))
	for _, perm := range perms {
		r, err := GetPermissionRegex(perm)
		if err == nil {
			regexs[perm] = r
		}
	}
	return &Validator{
		perms:  perms,
		regexs: regexs,
	}
}

func (v *Validator) Validate(perm string) bool {
	if !ValidatePermission(perm) {
		return false
	}

	for _, regex := range v.regexs {
		if regex.MatchString(perm) {
			return true
		}
	}
	return false
}
