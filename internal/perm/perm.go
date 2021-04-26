package perm

import (
	"github.com/thoas/go-funk"
	"regexp"
	"strings"
)

// permRegex matches any string dot seperated with
// optional `*` wildcards. Example: `[-]this.is.*.permission`
var permRegex = regexp.MustCompile(`^-?(?:\w+|\*)(?:\.(\w+|\*))*$`)

func ValidatePermission(perm string) bool {
	return permRegex.MatchString(perm)
}

func GetPermissionRegexString(perm string) string {
	if !ValidatePermission(perm) {
		return perm
	}
	if strings.HasSuffix(perm, ".*") {
		perm = perm[:len(perm)-2]
		return GetPermissionRegexString(perm) + `(?:\.(\w+|\*))+`
	}
	// replace wildcard with regex
	perm = strings.ReplaceAll(perm, "*", `(?:\w+|\*)`)

	// escape the dots
	return strings.ReplaceAll(perm, ".", `\.`)
}

func GetPermissionRegex(perm string) (*regexp.Regexp, error) {
	return regexp.Compile(GetPermissionRegexString(perm))
}

type Validator struct {
	perms  []string
	regexs map[string]*regexp.Regexp
}

func NewValidator(perms []string) *Validator {
	v := &Validator{
		perms:  []string{},
		regexs: map[string]*regexp.Regexp{},
	}
	v.AppendSimple(perms)
	return v
}

func (v *Validator) Append(perms []string, superior bool) {
	if !superior {
		for _, perm := range perms {
			r, err := GetPermissionRegex(perm)
			if err != nil {
				continue
			}

			v.regexs[perm] = r
			v.perms = append(v.perms, perm)
		}
		return
	}

	v2 := NewValidator(perms)
	var toRemove []string
	for p, _ := range v.regexs {
		neg := strings.HasPrefix(p, "-")

		if neg && v2.ValidateRaw(strings.Replace(p, "-", "", 1)) {
			toRemove = append(toRemove, p)
		} else if !neg && v2.ValidateRaw("-"+p) {
			toRemove = append(toRemove, p)
		}
	}

	// remove from perms and regexes
	for _, s := range toRemove {
		delete(v.regexs, s)
	}
	v.perms = funk.SubtractString(v.perms, toRemove)

	v.Append(v2.perms, false)
}

func (v *Validator) AppendSimple(perms []string) {
	v.Append(perms, false)
}

func (v *Validator) Validate(perm string) bool {
	if !ValidatePermission(perm) {
		return false
	}

	// check for negatives
	if !v.ValidateRaw("-" + perm) {
		return false
	}
	return v.ValidateRaw(perm)
}

func (v *Validator) ValidateRaw(perm string) bool {
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

func (v *Validator) GetPermissions() []string {
	tmp := make([]string, len(v.perms))
	copy(tmp, v.perms)
	return tmp
}
