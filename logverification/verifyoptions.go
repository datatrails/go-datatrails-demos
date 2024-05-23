package logverification

type VerifyOptions struct {

	// tenantId is an optional tenant ID to use instead
	//  of the tenantId found on the eventJson.
	tenantId string

	// massifHeight is an optional massif height for the massif
	//  instead of the default.
	massifHeight uint8
}

type VerifyOption func(*VerifyOptions)

// WithTenantId is an optional tenant ID to use instead
//
//	of the tenantId found on the eventJson.
func WithTenantId(tenantId string) VerifyOption {
	return func(vo *VerifyOptions) { vo.tenantId = tenantId }
}

// WithMassifHeight is an optional massif height for the massif
//
//	instead of the default.
func WithMassifHeight(massifHeight uint8) VerifyOption {
	return func(vo *VerifyOptions) { vo.massifHeight = massifHeight }
}

// ParseOptions parses the given options into a VerifyOptions struct
func ParseOptions(options ...VerifyOption) VerifyOptions {
	verifyOptions := VerifyOptions{
		massifHeight: defaultMassifHeight, // set the default massif height first
	}

	for _, option := range options {
		option(&verifyOptions)
	}

	return verifyOptions
}
