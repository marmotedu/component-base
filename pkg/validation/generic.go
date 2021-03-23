// Copyright 2020 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"

	"github.com/marmotedu/component-base/pkg/validation/field"
)

const (
	qnameCharFmt     string = "[A-Za-z0-9]"
	qnameExtCharFmt  string = "[-A-Za-z0-9_.]"
	qualifiedNameFmt string = "(" + qnameCharFmt + qnameExtCharFmt + "*)?" + qnameCharFmt
)

const (
	qualifiedNameErrMsg    string = "must consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"
	qualifiedNameMaxLength int    = 63
)

var qualifiedNameRegexp = regexp.MustCompile("^" + qualifiedNameFmt + "$")

// IsQualifiedName tests whether the value passed is what IAM calls a
// "qualified name". This is a format used in various places throughout the
// system. If the value is not valid, a list of error strings is returned.
// Otherwise an empty list (or nil) is returned.
func IsQualifiedName(value string) []string {
	var errs []string
	parts := strings.Split(value, "/")
	var name string
	switch len(parts) {
	case 1:
		name = parts[0]
	// nolint:gomnd // no need
	case 2:
		var prefix string
		prefix, name = parts[0], parts[1]
		if len(prefix) == 0 {
			errs = append(errs, "prefix part "+EmptyError())
		} else if msgs := IsDNS1123Subdomain(prefix); len(msgs) != 0 {
			errs = append(errs, prefixEach(msgs, "prefix part ")...)
		}
	default:
		return append(
			errs,
			"a qualified name "+RegexError(
				qualifiedNameErrMsg,
				qualifiedNameFmt,
				"MyName",
				"my.name",
				"123-abc",
			)+" with an optional DNS subdomain prefix and '/' (e.g. 'example.com/MyName')",
		)
	}

	if len(name) == 0 {
		errs = append(errs, "name part "+EmptyError())
	} else if len(name) > qualifiedNameMaxLength {
		errs = append(errs, "name part "+MaxLenError(qualifiedNameMaxLength))
	}
	if !qualifiedNameRegexp.MatchString(name) {
		errs = append(
			errs,
			"name part "+RegexError(qualifiedNameErrMsg, qualifiedNameFmt, "MyName", "my.name", "123-abc"),
		)
	}
	return errs
}

const labelValueFmt string = "(" + qualifiedNameFmt + ")?"

const labelValueErrMsg string = "a valid label must be an empty string or consist of alphanumeric characters, '-', '_' or '.', and must start and end with an alphanumeric character"

// LabelValueMaxLength is a label's max length.
const LabelValueMaxLength int = 63

var labelValueRegexp = regexp.MustCompile("^" + labelValueFmt + "$")

// IsValidLabelValue tests whether the value passed is a valid label value.  If
// the value is not valid, a list of error strings is returned.  Otherwise an
// empty list (or nil) is returned.
func IsValidLabelValue(value string) []string {
	var errs []string
	if len(value) > LabelValueMaxLength {
		errs = append(errs, MaxLenError(LabelValueMaxLength))
	}
	if !labelValueRegexp.MatchString(value) {
		errs = append(errs, RegexError(labelValueErrMsg, labelValueFmt, "MyValue", "my_value", "12345"))
	}
	return errs
}

const dns1123LabelFmt string = "[a-z0-9]([-a-z0-9]*[a-z0-9])?"

const dns1123LabelErrMsg string = "a DNS-1123 label must consist of lower case alphanumeric characters or '-', and must start and end with an alphanumeric character"

// DNS1123LabelMaxLength is a label's max length in DNS (RFC 1123).
const DNS1123LabelMaxLength int = 63

var dns1123LabelRegexp = regexp.MustCompile("^" + dns1123LabelFmt + "$")

// IsDNS1123Label tests for a string that conforms to the definition of a label in
// DNS (RFC 1123).
func IsDNS1123Label(value string) []string {
	var errs []string
	if len(value) > DNS1123LabelMaxLength {
		errs = append(errs, MaxLenError(DNS1123LabelMaxLength))
	}
	if !dns1123LabelRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123LabelErrMsg, dns1123LabelFmt, "my-name", "123-abc"))
	}
	return errs
}

const dns1123SubdomainFmt string = dns1123LabelFmt + "(\\." + dns1123LabelFmt + ")*"

const dns1123SubdomainErrorMsg string = "a DNS-1123 subdomain must consist of lower case alphanumeric characters, '-' or '.', and   must start and end with an alphanumeric character"

// DNS1123SubdomainMaxLength is a subdomain's max length in DNS (RFC 1123).
const DNS1123SubdomainMaxLength int = 253

var dns1123SubdomainRegexp = regexp.MustCompile("^" + dns1123SubdomainFmt + "$")

// IsDNS1123Subdomain tests for a string that conforms to the definition of a
// subdomain in DNS (RFC 1123).
func IsDNS1123Subdomain(value string) []string {
	var errs []string
	if len(value) > DNS1123SubdomainMaxLength {
		errs = append(errs, MaxLenError(DNS1123SubdomainMaxLength))
	}
	if !dns1123SubdomainRegexp.MatchString(value) {
		errs = append(errs, RegexError(dns1123SubdomainErrorMsg, dns1123SubdomainFmt, "example.com"))
	}
	return errs
}

// IsValidPortNum tests that the argument is a valid, non-zero port number.
func IsValidPortNum(port int) []string {
	if 1 <= port && port <= 65535 {
		return nil
	}
	return []string{InclusiveRangeError(1, 65535)}
}

// IsInRange tests that the argument is in an inclusive range.
func IsInRange(value int, min int, max int) []string {
	if value >= min && value <= max {
		return nil
	}
	return []string{InclusiveRangeError(min, max)}
}

// IsValidIP tests that the argument is a valid IP address.
func IsValidIP(value string) []string {
	if net.ParseIP(value) == nil {
		return []string{"must be a valid IP address, (e.g. 10.9.8.7)"}
	}
	return nil
}

// IsValidIPv4Address tests that the argument is a valid IPv4 address.
func IsValidIPv4Address(fldPath *field.Path, value string) field.ErrorList {
	var allErrors field.ErrorList
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() == nil {
		allErrors = append(allErrors, field.Invalid(fldPath, value, "must be a valid IPv4 address"))
	}
	return allErrors
}

// IsValidIPv6Address tests that the argument is a valid IPv6 address.
func IsValidIPv6Address(fldPath *field.Path, value string) field.ErrorList {
	var allErrors field.ErrorList
	ip := net.ParseIP(value)
	if ip == nil || ip.To4() != nil {
		allErrors = append(allErrors, field.Invalid(fldPath, value, "must be a valid IPv6 address"))
	}
	return allErrors
}

const (
	percentFmt    string = "[0-9]+%"
	percentErrMsg string = "a valid percent string must be a numeric string followed by an ending '%'"
)

var percentRegexp = regexp.MustCompile("^" + percentFmt + "$")

// IsValidPercent checks that string is in the form of a percentage.
func IsValidPercent(percent string) []string {
	if !percentRegexp.MatchString(percent) {
		return []string{RegexError(percentErrMsg, percentFmt, "1%", "93%")}
	}
	return nil
}

// MaxLenError returns a string explanation of a "string too long" validation
// failure.
func MaxLenError(length int) string {
	return fmt.Sprintf("must be no more than %d characters", length)
}

// RegexError returns a string explanation of a regex validation failure.
func RegexError(msg string, fmt string, examples ...string) string {
	if len(examples) == 0 {
		return msg + " (regex used for validation is '" + fmt + "')"
	}
	msg += " (e.g. "
	for i := range examples {
		if i > 0 {
			msg += " or "
		}
		msg += "'" + examples[i] + "', "
	}
	msg += "regex used for validation is '" + fmt + "')"
	return msg
}

// EmptyError returns a string explanation of a "must not be empty" validation
// failure.
func EmptyError() string {
	return "must be non-empty"
}

func prefixEach(msgs []string, prefix string) []string {
	for i := range msgs {
		msgs[i] = prefix + msgs[i]
	}
	return msgs
}

// InclusiveRangeError returns a string explanation of a numeric "must be
// between" validation failure.
func InclusiveRangeError(lo, hi int) string {
	return fmt.Sprintf(`must be between %d and %d, inclusive`, lo, hi)
}

const (
	minPassLength = 8
	maxPassLength = 16
)

// IsValidPassword validate password.
func IsValidPassword(password string) error {
	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			hasNumber = true
			passLen++
		case unicode.IsUpper(ch):
			hasUpper = true
			passLen++
		case unicode.IsLower(ch):
			hasLower = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			hasSpecial = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}

	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !hasLower {
		appendError("lowercase letter missing")
	}
	if !hasUpper {
		appendError("uppercase letter missing")
	}
	if !hasNumber {
		appendError("at least one numeric character required")
	}
	if !hasSpecial {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(
			fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength),
		)
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}

	return nil
}
