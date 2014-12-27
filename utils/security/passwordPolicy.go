package security


type passwordValidator func(string) bool

type passwordPolicy struct {
    MinimumLength int
    Validations []passwordValidator
}

var (
    minimumLength = 6
    PasswordPolicy = &passwordPolicy{
        Validations: []passwordValidator{
            meetsMinLength,
        },
    }
)

/*
 * Declare password validation functions below.
 */

func meetsMinLength(password string) bool {
    return len(password) >= minimumLength
}
