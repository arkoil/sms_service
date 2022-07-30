package sms_ru

type Options func(api APIHandler) APIHandler

func WithTest() Options {
	return func(api APIHandler) APIHandler {
		api.test = true
		return api
	}
}
func JSONFormat() Options {
	return func(api APIHandler) APIHandler {
		api.jsonResponse = true
		return api
	}
}
