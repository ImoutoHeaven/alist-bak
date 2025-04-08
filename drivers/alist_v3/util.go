package alist_v3

import (
	"fmt"
	"net/http"

	"github.com/alist-org/alist/v3/drivers/base"
	"github.com/alist-org/alist/v3/internal/op"
	"github.com/alist-org/alist/v3/pkg/utils"
	"github.com/alist-org/alist/v3/server/common"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
)

func (d *AListV3) login() error {
	if d.Username == "" {
		return nil
	}
	var resp common.Resp[LoginResp]
	_, _, err := d.request("/auth/login", http.MethodPost, func(req *resty.Request) {
		req.SetResult(&resp).SetBody(base.Json{
			"username": d.Username,
			"password": d.Password,
		})
	})
	if err != nil {
		return err
	}
	d.Token = resp.Data.Token
	op.MustSaveDriverStorage(d)
	return nil
}

func (d *AListV3) request(api, method string, callback base.ReqCallback, retry ...bool) ([]byte, int, error) {
	url := d.Address + "/api" + api
	req := base.RestyClient.R()
	req.SetHeader("Authorization", d.Token)
	if callback != nil {
		callback(req)
	}
	res, err := req.Execute(method, url)
	if err != nil {
		code := 0
		if res != nil {
			code = res.StatusCode()
		}
		return nil, code, err
	}
	log.Debugf("[alist_v3] response body: %s", res.String())
	if res.StatusCode() >= 400 {
		// 对于HTTP状态码403的请求，我们不返回错误，而是正常继续
		if res.StatusCode() == 403 {
			log.Warnf("[alist_v3] Received 403 response from %s but continuing anyway", api)
			return res.Body(), 200, nil
		}
		return nil, res.StatusCode(), fmt.Errorf("request failed, status: %s", res.Status())
	}
	code := utils.Json.Get(res.Body(), "code").ToInt()
	if code != 200 {
		// 对于API返回码403的请求，我们不返回错误，而是正常继续
		if code == 403 {
			log.Warnf("[alist_v3] Received code 403 from %s but continuing anyway", api)
			return res.Body(), 200, nil
		}
		
		if (code == 401 || code == 403) && !utils.IsBool(retry...) {
			err = d.login()
			if err != nil {
				return nil, code, err
			}
			return d.request(api, method, callback, true)
		}
		return nil, code, fmt.Errorf("request failed,code: %d, message: %s", code, utils.Json.Get(res.Body(), "message").ToString())
	}
	return res.Body(), 200, nil
}
