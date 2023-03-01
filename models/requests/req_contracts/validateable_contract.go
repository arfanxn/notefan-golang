package requests

type ValidateableContract interface {
	Validate() error
}
