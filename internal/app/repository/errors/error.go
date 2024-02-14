// Package contains custom repository errors.
package repoerrors

// RepoError is a custom repository error object.
type RepoError struct {
	Err         error
	ShortURL    string
	OriginalURL string
}

// Unwrap refers to the initial error Error method.
func (er *RepoError) Error() string {
	return er.Err.Error()
}

// Unwrap refers to the initial error Unwrap method.
func (er *RepoError) Unwrap() error {
	return er.Err
}

// New returns a custom repository error instance.
func New(err error, shortURL, originalURL string) error {
	return &RepoError{
		Err:         err,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
