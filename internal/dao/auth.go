package dao

import "github.com/ludyyy-lu/goBlogService/internal/model"

//获取认证
func (d *Dao) GetAuth(appKey, appSecret string) (model.Auth, error) {
	auth := model.Auth{AppKey:appKey,AppSecret: appSecret}
	return auth.Get(d.engine)
}