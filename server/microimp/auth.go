package microimp

import (
	db "github.com/DQFSN/blog/server/db"
	"github.com/DQFSN/blog/server/model"
	"github.com/DQFSN/blog/server/util"
	"strings"

	mpb "github.com/DQFSN/blog/proto/micro"
	"context"
	"fmt"
)

type AuthHandler struct{}

func (auth AuthHandler) LogIn(ctx context.Context, in *mpb.LogInRequest, out *mpb.LogInReply) error {

	fmt.Printf("新请求--->%v\n", in)
	mysqlDB := db.DB()
	user := model.User{}
	err := mysqlDB.Find(&user, "email = ?", in.Email).Error

	if err != nil {
		out.Status = fmt.Sprintf("LogIn : %s %s", in.Email, err)
		return  err
	}


	if util.ComparePasswords(user.Password,[]byte(in.Password)) {
		out.Status =  "ok: " + in.Email + " " + in.Password
		return  nil
	}

	out.Status = "wrong : " + in.Email + " " + in.Password
	return  nil

}

func (auth AuthHandler) SignUp(ctx context.Context, in *mpb.SignUpRequest, out *mpb.SignUpReply) error {
	if strings.Contains(in.Email, "@") && len(in.Password) > 0 && in.Password == in.PasswordCheck {
		mysqlDB := db.DB()

		hashPassword := util.HashAndSalt([]byte(in.Password))
		user := model.User{Email: in.Email, Password: hashPassword}
		err := mysqlDB.Create(&user).Error
		if err != nil {
			out.Status = fmt.Sprintf("insert: %s %s", in.Email, err)
			return  err
		}
		out.Status = "ok: " + in.Email + " " + in.Password
		return nil
	}
	out.Status = "wrong : " + in.Email + " " + in.Password
	return nil
}

func (auth AuthHandler) ModifyUser(ctx context.Context, in *mpb.ModifyUserRequest, out *mpb.ModifyUserReply) error {

	if strings.Contains(in.EmailNow, "@") && in.EmailPre != in.EmailNow && in.PasswordPre != in.PasswordNow {
		mysqlDB := db.DB()
		user := model.User{Email: in.EmailPre, Password: in.PasswordPre}
		mysqlDB.First(&user)

		//更新
		user.Email = in.EmailNow
		user.Password = util.HashAndSalt([]byte(in.PasswordNow))
		err := mysqlDB.Save(&user).Error
		if err != nil {
			out.Status = fmt.Sprintf("update: %s %s", in.EmailPre, err)
			return err
		}
		out.Status = "ok: " + in.EmailNow + " " + in.PasswordNow
		return nil
	}
	out.Status = "wrong : " + in.EmailNow + " " + in.PasswordNow
	return nil
}