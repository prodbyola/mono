package mono

import (
	"errors"
)

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

var verificationMethods = map[string]BvnVerificationMethodType{
	"email":           EmailMethod,
	"phone":           PhoneMethod,
	"phone_1":         PhoneMethod1,
	"alternate_phone": AlternatePhoneMethod,
}

func (vm BvnVerificationMethodType) String() string {
	for key, method := range verificationMethods {
		if method == vm {
			return key
		}
	}

	return "email"
}

// str must be one of the of the following: "email", "phone", "phone_1", "alternate_phone"
//
// Please check Step 2 on https://docs.mono.co/docs/bvn-lookup-integration-guide
func (vm *BvnVerificationMethodType) FromString(str string) error {
	method, found := verificationMethods[str]
	if !found {
		return errors.New("Invalid method type")
	}

	*vm = method
	return nil
}
