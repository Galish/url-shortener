package repository

import "errors"

var ErrConflict = errors.New("data conflict")

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

func NewRepoError(err error, shortURL, originalURL string) error {
	return &RepoErr{
		Err:         err,
		ShortURL:    shortURL,
		OriginalURL: originalURL,
	}
}

func AsErrConflict(err error) *RepoErr {
	var repoErr *RepoErr

	if err != nil &&
		errors.Is(err, ErrConflict) &&
		errors.As(err, &repoErr) {
		return repoErr
	}

	return nil
}
