package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"go_vip_video/common"
	"go_vip_video/models"
	"log"
	"time"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) Oauth() {
	toUrl := c.GetString("toUrl", "/user")
	redirectUrl := fmt.Sprintf(`https://open.weixin.qq.com/connect/oauth2/authorize?appid=wxcb331d5bde931fd0&redirect_uri=http://new.qiandao.name/login&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect`, toUrl)
	c.Ctx.Redirect(301, redirectUrl)
}

func (c *UserController) Login() {
	code := c.GetString("code")
	state := c.GetString("state")

	wa := common.WechatAccount
	oa := wa.GetOauth()

	resAccessToken, err := oa.GetUserAccessToken(code)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("获取到的授权数据: %v", resAccessToken))
	info, err := oa.GetUserInfo(resAccessToken.AccessToken, resAccessToken.OpenID, "")
	if err != nil {
		panic(err)
	}
	//创建账户
	user := &models.User{
		Nickname:    info.Nickname,
		OpenId:      info.OpenID,
		HeadImgURL:  info.HeadImgURL,
		Sex:         info.Sex,
		City:        info.City,
		Province:    info.Province,
		Unionid:     info.Unionid,
		CreatedTime: time.Now().Unix(),
		UpdatedTime: time.Now().Unix(),
	}

	if err = user.FirstOrCreateByOpenId(models.GlobalORMDB); err != nil {
		log.Fatal(err)
	}

	c.SetSession("uid", user.ID)
	c.Ctx.Redirect(301, state)
}

//用户中心
func (c *UserController) UserCenter() {
	uid := c.GetSession("uid").(int64)
	//根据uid 查找用户信息
	user := &models.User{ID: uid}
	if err := user.LoadById(models.GlobalORMDB); err != nil {
		panic(err)
	}
	c.Data["UserInfo"] = user

	c.TplName = "user.tpl"
}
