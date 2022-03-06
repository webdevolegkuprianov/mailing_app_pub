package store

//Store
type Store interface {
	Serv() ServRepository
	Data() DataRepository
}
