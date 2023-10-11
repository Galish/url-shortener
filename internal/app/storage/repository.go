package storage

type Repository interface {
	Get(string) (string, error)
	Set(string, string)
	Has(string) bool
}
