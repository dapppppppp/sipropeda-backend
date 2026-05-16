package failure

// InternalError membungkus error bawaan sistem
func InternalError(err error) error {
	return err
}

// BadRequest membungkus error dari input yang tidak valid (400)
func BadRequest(err error) error {
	return err
}