package services

import (
	"context"
	"encoding/base32"
	"fmt"
	"mime/multipart"
	entity "mydream_project/app/entities"
	"mydream_project/app/features/auth/repository"
	dependcy "mydream_project/config/dependency"
	"mydream_project/errorr"
	"mydream_project/helper"
	

	"github.com/go-playground/validator"
)

type (
	user struct {
		repo       repository.UserRepo
		validator *validator.Validate 
		dep      dependcy.Depend
	} 
    UserService interface {
		Login(ctx context.Context, req entity.LoginReq) (*entity.User, error) 
		Register(ctx context.Context, req entity.RegisterReq) error 
		VerifyEmail(ctx context.Context, verificationcode string) error
		ForgetPass(ctx context.Context, email string) error
		ResetPass(ctx context.Context, token string, newpass string) error
		GetProfile(ctx context.Context, id int) (*entity.User, error)
		Update(ctx context.Context, req entity.UpdateReq, file multipart.File) (*entity.User, error)
		Delete(ctx context.Context, id int) error 
	}
)

func NewUserService(repo repository.UserRepo, dep dependcy.Depend) UserService {
	return &user{repo: repo, dep: dep, validator: validator.New()}
}

func (u *user) Login(ctx context.Context, req entity.LoginReq) (*entity.User, error) {
	if err := u.validator.Struct(req); err != nil {
		u.dep.PromErr["ERROR"] = err.Error()
		u.dep.Log.Errorf("[ERROR] WHEN VALIDATE LOGIN REQ, Error: %v", err) 
		return nil, errorr.NewBad("Missing or Invalid Request Body")
	}
	user, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username) 
	if err != nil {
		u.dep.PromErr["error"] = err.Error() 
		return nil, err 
	}
	if err := helper.VerifyPassword(user.Password, req.Password); err != nil {
		u.dep.PromErr["error"] = err.Error() 
		u.dep.Log.Errorf("Error Service : %v", err) 
		return nil, errorr.NewBad("Wrong password")
	} 
	if user.IsVerified == false {
		u.dep.PromErr["error"] = "Email Not Verified" 
		return nil, errorr.NewBad("Email Not Verified")
	}
	return user, nil 
} 

func (u *user) Register(ctx context.Context, req entity.RegisterReq) error {
	if err := u.validator.Struct(req); err != nil {
		u.dep.PromErr["error"] = err.Error()
		u.dep.Log.Errorf("[ERROR] WHEN VALIDATE Regis REQ, Error: %v", err)
		return errorr.NewBad("Request body not valid")
	} 
	_, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username)
	if err == nil {
			u.dep.PromErr["error"] = "Email already registered"
			return errorr.NewBad("Username already registered")
	}
	user, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), req.Username) 
	if err == nil {
		if user.IsVerified == true {
			u.dep.PromErr["error"] = "Email already registered"
			return errorr.NewBad("Email already registered")
		}
	}
    passhash, err := helper.HashPassword(req.Password)
	if err != nil {
		u.dep.PromErr["error"] =  err.Error()
		u.dep.Log.Errorf("Erorr service: %v", err) 
		return errorr.NewBad("Register failed")
	}
	hashedEmailString := base32.StdEncoding.EncodeToString([]byte(req.Username))
	data := entity.User{
		Username: 					req.Username,
		Password: 					passhash,
		IsVerified:                 false,
		VerificationCode:           hashedEmailString,	
	}
	go func() {
		err := u.dep.Nsq.Publish("5", []byte(hashedEmailString))
		if err != nil {
			u.dep.PromErr["errpr"] = err.Error()
			u.dep.Log.Errorf("[FAILED] to publish to NSQ: %v", err)
			return 
		}
	}()
	if err := u.repo.Create(u.dep.Db.WithContext(ctx), data); err != nil {
		u.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil 
}

func (u *user)VerifyEmail(ctx context.Context, verificationcode string) error {
	if err := u.repo.VerifyEmail(u.dep.Db.WithContext(ctx), verificationcode); err != nil {
		u.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil 
}

func (u *user) ForgetPass(ctx context.Context, email string) error {
	user, err := u.repo.FindByEmail(u.dep.Db.WithContext(ctx), email)
	if err != nil {
		u.dep.PromErr["error"] = err.Error() 
		return errorr.NewBad("Email not registered")
	}
	if user.IsVerified == false {
		u.dep.PromErr["error"] = "Email Not Verified" 
		return errorr.NewBad("Email not verified")
	}
	hashedEmailString := base32.StdEncoding.EncodeToString([]byte(user.Username))
	if err := u.repo.InsertForgotPassToken(u.dep.Db.WithContext(ctx), entity.ForgotPass{Token: hashedEmailString, Email: user.Username}); err != nil {
		u.dep.PromErr["error"] = err.Error()
		return err
	}
	go func() {
		err := u.dep.Nsq.Publish("6", []byte(hashedEmailString))
		if err != nil {
			u.dep.PromErr["error"] = err.Error()
			u.dep.Log.Errorf("[FAILED] to publish to NSQ: %v", err)
			return 
		}		
	}()
		return nil 
} 

func (u *user) ResetPass(ctx context.Context, token string, newpass string) error {
	if err := u.repo.ResetPass(u.dep.Db.WithContext(ctx), newpass, token); err != nil {
		u.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil 
} 

func (u *user) GetProfile(ctx context.Context, id int) (*entity.User, error) {
	res, err := u.repo.GetById(u.dep.Db.WithContext(ctx), id)
	if err != nil {
		u.dep.PromErr["error"] = err.Error()
		return nil, err
	}
	return res, nil 
}

func (u *user) Update(ctx context.Context, req entity.UpdateReq, file multipart.File) (*entity.User, error) {
	data := entity.User{} 
	if req.Password != "" {
		passhash, err := helper.HashPassword(req.Password) 
		if err != nil {
			u.dep.PromErr["error"] =  err.Error()
			u.dep.Log.Errorf("[ERROR] WHEN HASHING PASSWORD, Error: %v", err)
			return nil, errorr.NewBad("Registerfailed")
		}
		req.Password = passhash
	}
	if req.Username != "" {
		_, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username)
		if err == nil {
			u.dep.PromErr["error"] = "Username already registered" 
			return nil, errorr.NewBad("Username already registered")
		}
	}
	if req.Username != "" {
		user, err := u.repo.FindByUsername(u.dep.Db.WithContext(ctx), req.Username) 
		if err == nil {
			if user.IsVerified == true {
				u.dep.PromErr["error"] = "Email already registered" 
				return nil, errorr.NewBad("Email already registered")
			}
		}
		hashedEmailString := base32.StdEncoding.EncodeToString([]byte(req.Username))
		go func() {
			err := u.dep.Nsq.Publish("7", []byte(hashedEmailString))
			if err != nil {
				u.dep.PromErr["error"] = err.Error()
				u.dep.Log.Errorf("[FAILED] to publish to NSQ: %v", err)
				return 
			}
		}()
		data.VerificationCode = hashedEmailString
		data.IsVerified = true 
	}
	if file != nil {
		filename := fmt.Sprintf("%s_%s", "User", req.Username)
		if err1 := u.dep.Gcp.UploadFile(file, filename); err1 != nil {
			u.dep.PromErr["error"] = err1.Error()
			u.dep.Log.Errorf("Error Service : %v", err1)
			return nil, errorr.NewBad("Failed to upload image")
		}
		req.Username = filename 
		file.Close()
	}
	data.Username  = req.Username 
	data.Password  = req.Password 
	data.ID = uint(req.Id) 
	res, err := u.repo.Update(u.dep.Db.WithContext(ctx), data)
	if err != nil {
		u.dep.PromErr["error"] = err.Error() 
		return nil, err
	}
	return res, nil 
}

func (u *user) Delete(ctx context.Context, id int ) error {
	data := entity.User{} 
	data.ID = uint(id)
	err := u.repo.Delete(u.dep.Db.WithContext(ctx), data) 
	if err != nil {
		u.dep.PromErr["error"] = err.Error()
		return err
	}
	return nil
}