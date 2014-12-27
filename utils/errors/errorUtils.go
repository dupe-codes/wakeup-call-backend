package errorUtils


type InvalidFieldsError struct {
    Message string
    Fields []string
}

func (err *InvalidFieldsError) Error() string { return err.Message }
