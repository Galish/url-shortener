package repoerrors

// RepoErr is a custom repository error object.
type RepoErr struct {
	Err         error
	ShortURL    string
	OriginalURL string
}

// Unwrap refers to the initial error Error method.
func (er *RepoErr) Error() string {
	return er.Err.Error()
}

// Unwrap refers to the initial error Unwrap method.
func (er *RepoErr) Unwrap() error {
	return er.Err
}

// New returns a custom repository error instance.
func New(err error, shortURL, originalURL string) error {
	return &RepoErr{
		Err:         err,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
