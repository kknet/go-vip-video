package controllers

import (
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"go_vip_video/common"
	"go_vip_video/dto/m360k"
	"go_vip_video/dto/pc"
	"go_vip_video/service"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {

	var err error

	ca := common.GoCache
	dianying, found := ca.Get("index::dianying")
	if !found {
		dianying, err = service.GetPCList("dianying", "rankhot", 1)
		if err != nil {

		}
		ca.Set("index::dianying", dianying, cache.DefaultExpiration)
	}

	dianshi, found := ca.Get("index::dianshi")
	if !found {
		dianshi, err = service.GetPCList("dianshi", "rankhot", 1)
		if err != nil {

		}
		ca.Set("index::dianshi", dianshi, cache.DefaultExpiration)
	}

	zongyi, found := ca.Get("index::zongyi")
	if !found {
		zongyi, err = service.GetPCList("zongyi", "rankhot", 1)
		if err != nil {

		}
		ca.Set("index::zongyi", zongyi, cache.DefaultExpiration)
	}
	dongman, found := ca.Get("index::dongman")
	if !found {
		dongman, err = service.GetPCList("dongman", "rankhot", 1)
		if err != nil {

		}
		ca.Set("index::dongman", dongman, cache.DefaultExpiration)
	}

	//获取幻灯片
	swiper, found := ca.Get("index::swiper")
	if !found {
		sd, err := service.NewSwiperDocument()
		if err != nil {
			panic(err)
		}
		swiper = sd.SwiperResult()
		ca.Set("index::swiper", swiper, cache.DefaultExpiration)
	}

	c.Data["Swiper"] = swiper.([]*m360k.Swiper)
	c.Data["dianying"] = dianying.([]*pc.VideoItem)[:21]
	c.Data["dianshi"] = dianshi.([]*pc.VideoItem)[:21]
	c.Data["zongyi"] = zongyi.([]*pc.VideoItem)[:21]
	c.Data["dongman"] = dongman.([]*pc.VideoItem)[:21]
	c.TplName = "index.tpl"
}
