package handler

import (
	"net/http"
)

func (s *shortenerService) getFullLink(w http.ResponseWriter, r *http.Request) {
	fullLink, err := s.store.Get(r.URL.Path[1:])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	http.Redirect(w, r, fullLink, http.StatusTemporaryRedirect)
}
