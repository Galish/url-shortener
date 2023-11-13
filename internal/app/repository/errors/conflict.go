package repoErrors

import "errors"

var ErrConflict = errors.New("data conflict")

func AsErrConflict(err error) *RepoErr {
	var repoErr *RepoErr

	if err != nil &&
		errors.Is(err, ErrConflict) &&
		errors.As(err, &repoErr) {
		return repoErr
	}

	return nil
}
