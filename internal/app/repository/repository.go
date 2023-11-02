package repository

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
	Has(string) bool
}
