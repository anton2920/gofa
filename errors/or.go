package errors

func Or(errs ...error) error {
	for i := 0; i < len(errs); i++ {
		if errs[i] != nil {
			return errs[i]
		}
	}
	return nil
}
