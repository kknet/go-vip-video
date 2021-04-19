package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/patrickmn/go-cache"
	"go_vip_video/dto"
	"go_vip_video/service"
	"go_vip_video/vcache"
	"strconv"
)

var lines = []dto.Lines{{
	Name: "线路一",
	Url:  "https://z1.m1907.cn/?jx=",
}, {
	Name: "线路二",
	Url:  "https://jx.qiandao.name/pangu/index.php?url=",
}}

type DetailController struct {
	beego.Controller
	vId  string //视频id
	cat  int    //视频分类
	site string //选中站点
	num  string //选中剧集
	jxID int    //解析id

	detail *dto.Detail //详情
	sites  []*dto.Site //站点
	links  []*dto.Link //剧集
}

//id+cat+站点+剧集  可以定位到具体url
func (c *DetailController) init() {
	//请求参数
	c.vId = c.Ctx.Input.Param(":id")
	cat, _ := strconv.Atoi(c.Ctx.Input.Param(":cat"))
	c.cat = cat

	//获取详情
	c.detail = c.getDetail().(*dto.Detail)

	//获取全部站点信息
	sites := c.getSites()
	c.sites = sites
	c.site = c.GetString("site")
	//设置默认站点
	if c.site == "" {
		c.site = sites[0].Code
	}
	//获取视频链接列表
	if c.cat != 1 {
		c.links = c.getLinks().([]*dto.Link)
	}
	c.num = c.GetString("num", c.links[0].Num)

	//获取解析id
	c.jxID, _ = c.GetInt("jxId", 0)
}
func (c *DetailController) Get() {
	c.init()
	link := c.getLinkBySite()
	c.Data["Detail"] = c.detail.Data
	c.Data["Id"] = c.vId
	c.Data["Cat"] = c.cat
	c.Data["Link"] = link
	c.Data["Site"] = c.site
	c.Data["Links"] = c.links
	c.Data["Sites"] = c.sites
	c.Data["Num"] = c.num
	//线路
	c.Data["JxLines"] = lines
	c.Data["JxID"] = c.jxID
	c.Data["JxUrl"] = lines[c.jxID].Url
	c.TplName = "detail.tpl"
}

//获取站点信息
func (c *DetailController) getSites() []*dto.Site {
	var sites []*dto.Site

	bySite, err := service.NewSite(c.vId, c.cat)
	if err != nil {
		panic(err)
	}
	sites, c.links, err = bySite.Do()
	if err != nil {
		panic(err)
	}
	return sites
}

//获取剧集信息
func (c *DetailController) getLinks() interface{} {
	links, found := vcache.GoCache.Get(fmt.Sprintf("links::cat:%d::site:%s::vid:%s", c.cat, c.site, c.detail.Data.Rpt.VideoID))
	if !found {

		getLink, err := service.NewGetLink(c.detail.Data.Rpt.VideoID, c.cat, c.site)
		if err != nil {
			panic(err)
		}
		do, err := getLink.Do()
		if err != nil {
			panic(err)
		}
		linkHtml := do.Data
		links, err = service.Parse(linkHtml, c.cat)
		if err != nil {
			panic(err)
		}
		vcache.GoCache.Set(fmt.Sprintf("links::cat:%d::site:%s::vid:%s", c.cat, c.site, c.detail.Data.Rpt.VideoID), links, cache.DefaultExpiration)

	}
	return links
}

func (c *DetailController) getLinkBySite() string {
	for _, link := range c.links {
		if c.num == "" {
			return link.Url
		}
		if c.num == link.Num {
			return link.Url
		}
	}
	return ""
}

func (c *DetailController) getDetail() interface{} {
	var err error
	detail, found := vcache.GoCache.Get(fmt.Sprintf("detail::cat:%d::vid:%s", c.cat, c.vId))
	if !found {
		detail, err = service.NewDetail(c.cat, c.vId).Do()
		if err != nil {
			panic(err)
		}
		vcache.GoCache.Set(fmt.Sprintf("detail::cat:%d::vid:%s", c.cat, c.vId), detail, cache.DefaultExpiration)
	}

	return detail
}
