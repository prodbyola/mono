package mono

import "errors"

// Base URL for Mono v2 API
const BASE_URL = "https://api.withmono.com/v2/"

// BvnVerificationMethodType represents the type of verification method used for BVN verification.
// It defines constants for various verification methods, such as email, phone methods, and alternate phone.
//
// Please check Step 2 on https://docs.mono.co/docs/bvn-lookup-integration-guide
type BvnVerificationMethodType int

const (
	// email
	EmailMethod BvnVerificationMethodType = iota

	// phone
	PhoneMethod

	// phone_1
	PhoneMethod1

	// alternative_phone
	AlternatePhoneMethod
)

func (vm BvnVerificationMethodType) String() string {
	switch vm {
	case PhoneMethod:
		return "phone"
	case PhoneMethod1:
		return "phone_1"
	case AlternatePhoneMethod:
		return "alternative_phone"
	default:
		return "email"
	}
}

// str must be one of the of the following: "email", "phone", "phone_1", "alternate_phone"
//
// Please check Step 2 on https://docs.mono.co/docs/bvn-lookup-integration-guide
func (vm *BvnVerificationMethodType) FromString(str string) error {
	verificationMethods := []BvnVerificationMethodType{
		EmailMethod,
		PhoneMethod,
		PhoneMethod1,
		AlternatePhoneMethod,
	}

	valid := false
	for _, method := range verificationMethods {
		if str == method.String() {
			vm = &method
			valid = true
		}
	}

	if !valid {
		return errors.New("Invalid method type")
	}

	return nil
}
