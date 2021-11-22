package lib

func WithIgnoreUnknownField(cs *HTTPClientService) {
	cs.PbDiscardUnknown = true
}

func WithHTTPRequestPreflight(f HTTPRequestPreflightHandler) {
	return func(hs *HTTPClientService) {
		hs.HttpRequestPreflight = f
	}
}

func WithHTTPResponseValidator(f HTTPResponseValidatorHandler) HTTPServiceConstructorArg {
	return func(hs *HTTPClientService) {
		hs.HttpResponseValidator = f
	}
}

func WithResponseValidator(f HTTPClientMethodValidatorHandler) HTTPServiceConstructorArg {
	return func(hs *HTTPClientService) {
		hs.ResponseValidator = f
	}
}
