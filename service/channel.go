package service

import (
	"encoding/json"
	"fmt"
	"go_vip_video/dto"
	"io/ioutil"
	"net/http"
)

type channelDataReq struct {
	url   string
	count int
	start int
	tid   int
	cid   int
}

func ChannelDataService(tid int, count int, start int, cid int) *channelDataReq {
	return &channelDataReq{count: count, start: start, tid: tid, cid: cid, url: fmt.Sprintf("https://ios.api.360kan.com/iphone/channel/datas?count=%d&start=%d&tid=%d&cid=%d", count, start, tid, cid)}
}

func (c *channelDataReq) Do() (*dto.ChannelDataResp, error) {
	response, err := http.Get(c.url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	//检出结果集
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	channelDataResp := &dto.ChannelDataResp{}
	body = body[32:]
	if err := json.Unmarshal(body, &channelDataResp); err != nil {
		return nil, err
	}
	return channelDataResp, nil
}
