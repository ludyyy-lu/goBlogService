package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ludyyy-lu/goBlogService/global"
	"github.com/ludyyy-lu/goBlogService/pkg/app"
	"github.com/ludyyy-lu/goBlogService/pkg/email"
	"github.com/ludyyy-lu/goBlogService/pkg/errcode"
)

func Recovery() gin.HandlerFunc{
	defailtMailer := email.NewEmail(&email.SMTPInfo{
		Host: global.EmailSetting.Host,
		Port: global.EmailSetting.Port,
		IsSSL: global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From: global.EmailSetting.From,
	})
	return func(c *gin.Context) {
		defer func() {
			if err := recover();err != nil {
				s := "panic recover err: %v"
				global.Logger.WithCallersFrames().Errorf(s,err)

				err := defailtMailer.SendMail(
					global.EmailSetting.To,
					fmt.Sprintf("抛出异常，发生时间：%d",time.Now().Unix()),
					fmt.Sprintf("错误信息：%v",err),
				)
				if err != nil {
					global.Logger.Panicf("mail.SendMail err: %v",err)
				}
				app.NewResponse(c).ToErrorResponse(errcode.ServerError)
				c.Abort()
			}
		}()
		c.Next()
	}
}