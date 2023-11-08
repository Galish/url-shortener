package repository

type Repository interface {
	Get(string) (string, error)
	Set(string, string) error
	SetBatch(...[2]string) error
	Has(string) bool
	Ping() (bool, error)
	Close() error
}
