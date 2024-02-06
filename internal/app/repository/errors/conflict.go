package repoerrors

import "errors"

// ErrConflict represents repository conflict error.
var ErrConflict = errors.New("data conflict")

// AsErrConflict checks whether the error is a repository conflict type.
func AsErrConflict(err error) *RepoErr {
	var repoErr *RepoErr

	if err != nil &&
		errors.Is(err, ErrConflict) &&
		errors.As(err, &repoErr) {
		return repoErr
	}

	return nil
}
