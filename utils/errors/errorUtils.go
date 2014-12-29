package errorUtils

type GeneralError struct {
    Message string
}

func (err *GeneralError) Error() string { return err.Message }

type InvalidFieldsError struct {
	Message string
	Fields  []string
}

func (err *InvalidFieldsError) Error() string { return err.Message }
