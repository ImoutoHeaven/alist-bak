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
	log.Debugf("[alist_v3] Making %s request to: %s", method, url)
	
	req := base.RestyClient.R()
	req.SetHeader("Authorization", d.Token)
	log.Debugf("[alist_v3] Using authorization token: %s", d.Token)
	
	if callback != nil {
		log.Debug("[alist_v3] Executing callback function")
		callback(req)
	}
	
	log.Debugf("[alist_v3] Executing request: %s %s", method, url)
	res, err := req.Execute(method, url)
	
	if err != nil {
		code := 0
		if res != nil {
			code = res.StatusCode()
		}
		log.Debugf("[alist_v3] Request error: %v, status code: %d", err, code)
		return nil, code, err
	}
	
	log.Debugf("[alist_v3] Response status: %s", res.Status())
	log.Debugf("[alist_v3] Response body: %s", res.String())
	
	if res.StatusCode() >= 400 {
		log.Debugf("[alist_v3] HTTP error status: %s", res.Status())
		return nil, res.StatusCode(), fmt.Errorf("request failed, status: %s", res.Status())
	}
	
	code := utils.Json.Get(res.Body(), "code").ToInt()
	log.Debugf("[alist_v3] Response code from body: %d", code)
	
	if code != 200 {
		message := utils.Json.Get(res.Body(), "message").ToString()
		log.Debugf("[alist_v3] Request failed with code: %d, message: %s", code, message)
		
		if (code == 401 || code == 403) && !utils.IsBool(retry...) {
			log.Debug("[alist_v3] Authentication error, attempting to re-login")
			err = d.login()
			if err != nil {
				log.Debugf("[alist_v3] Re-login failed: %v", err)
				return nil, code, err
			}
			log.Debug("[alist_v3] Re-login successful, retrying request")
			return d.request(api, method, callback, true)
		}
		return nil, code, fmt.Errorf("request failed, code: %d, message: %s", code, message)
	}
	
	log.Debug("[alist_v3] Request completed successfully")
	return res.Body(), 200, nil
}
