package mask

import "strings"

const (
	def        = "CENSORED"
	cvv        = "***"
	cardNumber = "CARD NUMBER"
	cardHolder = "CARD HOLDER"
)

func Default(_ string) string {
	return def
}

// CVV returns mask for CVV/CVS code.
// Example: ***.
func CVV(_ string) string {
	return cvv
}

// CardNumber returns mask for card number.
// Example: 400680******6735.
func CardNumber(input string) string {
	// remove all spaces
	input = strings.ReplaceAll(input, " ", "")

	// create mask if card number is valid
	if len(input) == 16 {
		return input[:6] + strings.Repeat("*", 6) + input[12:]
	}

	return cardNumber
}

// CardHolder returns "CARD HOLDER".
func CardHolder(_ string) string {
	return cardHolder
}

// IP returns mask for IP.
// Example v4: 127.***.***.1.
// Example v6: 2345:0425:****:****:****:****:****:23b5.
func IP(input string) string {
	const (
		sepv4  = "."
		maskv4 = "***"
		sepv6  = ":"
		maskv6 = "****"

		fakeIP   = "fake IP"
		fakeIPv4 = "fake IPv4"
		fakeIPv6 = "fake IPv6"
	)

	// IPv4
	if len(input) >= 7 || len(input) <= 15 {
		parts := strings.Split(input, sepv4)

		if len(parts) != 4 {
			return fakeIPv4
		}

		parts[1], parts[2] = maskv4, maskv4

		return strings.Join(parts, sepv4)
	}

	// IPv6
	if len(input) == 39 {
		parts := strings.Split(input, sepv6)

		if len(parts) != 8 {
			return fakeIPv6
		}

		parts[2], parts[3], parts[4], parts[5], parts[6] = maskv6, maskv6, maskv6, maskv6, maskv6

		return strings.Join(parts, sepv6)
	}

	return fakeIP
}
