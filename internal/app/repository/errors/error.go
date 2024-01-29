package repoerrors

// RepoErr is a custom repository error object.
type RepoErr struct {
	Err         error
	ShortURL    string
	OriginalURL string
}

func (er *RepoErr) Error() string {
	return er.Err.Error()
}

func (er *RepoErr) Unwrap() error {
	return er.Err
}

func New(err error, shortURL, originalURL string) error {
	return &RepoErr{
		Err:         err,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}
