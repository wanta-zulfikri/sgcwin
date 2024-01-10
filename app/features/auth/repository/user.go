package repository

import (
	 entity"mydream_project/app/entities"
	"mydream_project/errorr"
	"reflect"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	user struct {
		log *logrus.Logger 
	}
	UserRepo interface {
		Create(db *gorm.DB, user entity.User) error 
		FindByEmail(db *gorm.DB, email string) (*entity.User, error) 
		VerifyEmail(db *gorm.DB, verificationcode string)  error 
		InsertForgotPassToken(db *gorm.DB, req entity.ForgotPass) error 
		ResetPass(db *gorm.DB, newpass string, token string) error
		FindByUsername(db *gorm.DB, username string) (*entity.User, error) 
		GetById(db *gorm.DB, id int) (*entity.User, error) 
		Update(db *gorm.DB, user entity.User) (*entity.User, error) 
		Delete(db *gorm.DB, user entity.User) error
	}
) 

func NewUserRepo(log *logrus.Logger) UserRepo {
	return &user{log}
}

func (u *user) Create(db *gorm.DB, user entity.User) error {
	if err := db.Create(&user).Error; err != nil {
		u.log.Errorf("[ERROR] WHEN CREATE USER, Error: %v ", err) 
		return errorr.NewInternal("Internal Server Error")
	} 
	return nil 
}


func (u *user) FindByEmail(db *gorm.DB, email string) (*entity.User,error) {
	res := entity.User{}
	err := db.Where("email = ?", email).Find(&res).Error 
	if res.Username == "" {
		return nil, errorr.NewBad("Email not registered")
	} 
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR] WHEN FIND BY EMAIL, Error: %v ", err)
			return nil, errorr.NewInternal(err.Error())
		} else {
			u.log.Errorf("error Db: %v", err) 
			return nil, errorr.NewBad(err.Error())
		}
	}
	return &res, nil 
}


func (u *user) FindByUsername(db *gorm.DB, username string) (*entity.User, error) {
	res := entity.User{} 
	err := db.Where("username = ?" , username).First(&res).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR] WHEN FIND USERNAME,Error : %v ", err) 
			return nil, errorr.NewInternal(err.Error())
		} else {
			u.log.Errorf("error Db: %v", err) 
			return nil, errorr.NewInternal(err.Error())
		}
	}
	return &res, nil 
}

func (u *user) Delete(db *gorm.DB, user entity.User) error {
	if err := db.Where("user_id=?", user.ID).First(&entity.User{}).Error; err == nil {
		return errorr.NewBad("You still have active user")
	}
	if err := db.Delete(&user).Error; err != nil {
		u.log.Errorf("[ERROR] WHEN DELETE USER, Error: %v ", err)
		return errorr.NewInternal(err.Error())
	}
	return nil 
}


func (u *user) Update(db *gorm.DB, user entity.User) (*entity.User, error) {
	newdata := entity.User{}
	if err := db.First(&newdata, user.ID).Error; err == gorm.ErrRecordNotFound {
		u.log.Errorf("[ERROR]WHEN UPDATE USER, Error: %v", err) 
		return nil, errorr.NewBad("Id Not Found")
	}
	v := reflect.ValueOf(user) 
	n := reflect.ValueOf(&newdata).Elem() 
	for i := 0; i < v.NumField(); i++ {
		if val, ok := v.Field(i).Interface().(string); ok {
			if val != "" {
				n.Field(i).SetString(val)
			}
		}
	}
	if user.IsVerified {
		newdata.IsVerified = false 
	}
	if err := db.Save(&newdata).Error; err != nil {
		u.log.Errorf("errorr Db : %v", err) 
		return nil, errorr.NewInternal("error update user")
	} 
	return &newdata, nil
}

func (u *user) GetById(db *gorm.DB, id int) (*entity.User, error) {
	res := entity.User{}
	err := db.Find(&res, id).Error
	if res.Username == "" {
		return nil, errorr.NewBad("Id not found")
	}
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			u.log.Errorf("[ERROR] WHEN GET USER DATA BY ID, Error: %v", err) 
			return nil, errorr.NewInternal(err.Error())
		} else {
			u.log.Errorf("[ERROR]WHEN GET USER DATA BY ID, Error: %v", err)
			return nil, errorr.NewBad(err.Error())
		}
	}
	return &res, nil
} 

func (u *user) VerifyEmail(db *gorm.DB, verificationcode string) error {
	if err := db.Model(&entity.User{}).Where("verification_code = ?", verificationcode).Update("is_verified", true).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errorr.NewBad("Data not found")
		}
		u.log.Errorf("[ERROR] When Verify Email, Error: %v", err) 
		return errorr.NewInternal("Internal Server Error")
	} 
	return nil 
}

func (u *user) InsertForgotPassToken(db *gorm.DB, req entity.ForgotPass) error {
	if err := db.Create(req).Error; err != nil {
		u.log.Errorf("[ERROR]entering the password reset token, error:%v", err)
	}
	return nil
}

func (u *user) ResetPass(db *gorm.DB, newpass string, token string) error {
	return db.Transaction(func(db *gorm.DB) error{
		userdata := entity.ForgotPass{} 
		if err := db.Where("token = ? AND delete_at IS NULL", token).Find(&userdata).Error; err != nil {
			u.log.Errorf("[ERROR]WHEN Getting user information with forget token, error:%v", err)
			return errorr.NewInternal("Internal Server Error")
		} 
		if userdata.Email == "" {
			return errorr.NewBad("Data Not Found")
		}
		if err := db.Model(&entity.User{}).Where("email=?", userdata.Email).Update("password", newpass).Error; err != nil {
			u.log.Errorf("[ERROR]When entering the password reset token, error:%v", err)
			return errorr.NewInternal("Internal Server Error")
		}
		db.Where("token =? AND deleted_at IS NULL", token).Delete(&userdata)
		return nil 
	})
}