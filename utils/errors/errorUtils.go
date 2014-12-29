package errorUtils

type GeneralError struct {
    Message string
}

type InvalidFieldsError struct {
	GeneralError
	Fields  []string
}

func (err *GeneralError) Error() string { return err.Message }
