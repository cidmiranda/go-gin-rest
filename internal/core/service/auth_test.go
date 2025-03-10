package service_test

type loginTestedInput struct {
	email    string
	password string
}

type loginExpectedOutput struct {
	token string
	err   error
}
