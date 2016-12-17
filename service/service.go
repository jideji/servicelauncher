package service

type Services map[string]Service

type ServiceLoader func() Services

type Service interface {
	IsRunning() (bool, error)
	Name() string
	Pid() (int, error)
	Start() error
	Stop() error
}
