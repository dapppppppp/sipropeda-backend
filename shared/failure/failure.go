package failure

// InternalError membungkus error bawaan sistem
func InternalError(err error) error {
	return err
}