package service

type Services map[string]Service

type Service interface {
	IsRunning() (bool, error)
	Name() string
	Pid() (int, error)
	Start() error
	Stop() error
}
