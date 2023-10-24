package lookup

type LookUp struct {
	Bvn BvnLookUp
}

// New is a factory function for creating a LookUp instance.
// It initializes and returns a LookUp struct with the provided API key and a BvnLookUp instance.
//
// Parameters:
// - apiKey: A string representing the API key to be used for LookUp service.
//
// Returns:
// - LookUp: A LookUp struct with the specified API key and a BvnLookUp instance.
func New(apiKey string) LookUp {
	return LookUp{
		Bvn: NewBvnLookUp(apiKey),
	}
}
