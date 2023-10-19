package repository

type Repository interface {
	Get(string) (string, error)
	Set(string, string)
	Has(string) bool
}

func New() Repository {
	return newKVStore()
}
