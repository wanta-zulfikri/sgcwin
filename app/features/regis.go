package features

import (
	userserv "mydream_project/app/features/auth/services"
	userrepo "mydream_project/app/features/auth/repository"

	"go.uber.org/dig"
)

func RegisterRepo(C *dig.Container) error {
	if err := C.Provide(userrepo.NewUserRepo); err != nil {
		return err
	}
	return nil
}


func RegisterService(C *dig.Container) error {
	if err := C.Provide(userserv.NewUserService); err != nil {
		return err
	}
	return nil
}