package controllers

import (
	"../../deps/lessgo/data/rdo"
	"../../deps/lessgo/data/rdo/base"
	"../../deps/lessgo/net/email"
	"../../deps/lessgo/pagelet"
	"../../deps/lessgo/pass"
	"../../deps/lessgo/utils"
	"../conf"
	"../models/login"
	"../reg/signup"
	"fmt"
	"io"
	"time"
)

type Reg struct {
	*pagelet.Controller
}

func (c Reg) IndexAction() {

}

func (c Reg) SignUpAction() {
	c.ViewData["continue"] = c.Params.Get("continue")
}

func (c Reg) SignUpRegAction() {

	c.AutoRender = false

	var rsp struct {
		ResponseJson
		Data struct {
			Continue    string `json:"continue"`
			AccessToken string `json:"access_token"`
		} `json:"data"`
	}
	rsp.ApiVersion = apiVersion
	rsp.Status = 400
	rsp.Message = "Bad Request"
	rsp.Data.Continue = "/ids"

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	uuid := utils.StringNewRand36(8)
	uname := utils.StringNewRand36(8)
	c.Params.Set("uname", uname)

	if err := signup.Validate(c.Params); err != nil {
		rsp.Message = err.Error()
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsu, err := dcn.Base.Query(q)
	if err == nil && len(rsu) == 1 {
		rsp.Message = "The `Email` already exists, please choose another one"
		return
	}

	pass, err := pass.HashDefault(c.Params.Get("passwd"))
	if err != nil {
		return
	}

	item := map[string]interface{}{
		"uuid":     uuid,
		"uname":    uname,
		"email":    c.Params.Get("email"),
		"pass":     pass,
		"name":     c.Params.Get("name"),
		"status":   1,
		"roles":    "100",
		"timezone": "UTC",                    // TODO
		"created":  base.TimeNow("datetime"), // TODO
		"updated":  base.TimeNow("datetime"), // TODO
	}
	if _, err := dcn.Base.Insert("ids_login", item); err != nil {
		rsp.Status = 500
		rsp.Message = "Can not write to database"
		return
	}

	rsp.Status = 200
	rsp.Message = ""
}

func (c Reg) ForgotPassAction() {
}

func (c Reg) ForgotPassPutAction() {

	c.AutoRender = false

	var rsp struct {
		ResponseJson
		Data struct {
			Continue string `json:"continue"`
		} `json:"data"`
	}
	rsp.ApiVersion = apiVersion
	rsp.Status = 400
	rsp.Message = "Bad Request"
	rsp.Data.Continue = "/ids"

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	if err := login.EmailSetValidate(c.Params); err != nil {
		rsp.Message = err.Error()
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsl, err := dcn.Base.Query(q)
	if err != nil || len(rsl) != 1 {
		rsp.Message = "Email can not found"
		return
	}

	id := utils.StringNewRand36(24)
	item := map[string]interface{}{
		"id":      id,
		"status":  0,
		"email":   c.Params.Get("email"),                 // TODO
		"expired": base.TimeNowAdd("datetime", "+3600s"), // TODO
	}
	if _, err := dcn.Base.Insert("ids_resetpass", item); err != nil {
		rsp.Status = 500
		rsp.Message = "Can not write to database"
		return
	}

	mr, err := email.MailerPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	cfg := conf.ConfigFetch()

	// TODO tempate
	body := fmt.Sprintf(`<html>
<body>
<div>You recently requested a password reset for your %s account. To create a new password, click on the link below:</div>
<br>
<a href="http://%s/ids/reg/pass-reset?id=%s">Reset My Password</a>
<br>
<div>This request was made on %s.</div>
<br>
<div>Regards,</div>
<div>%s Account Services</div>

<div>********************************************************</div>
<div>Please do not reply to this message. Mail sent to this address cannot be answered.</div>
</body>
</html>`, cfg.ServiceName, c.Request.Host, id, base.TimeNow("datetime"), cfg.ServiceName)

	err = mr.SendMail(c.Params.Get("email"), c.T("Reset your password"), body)

	rsp.Status = 200
	rsp.Message = ""
}

func (c Reg) PassResetAction() {

	if c.Params.Get("id") == "" {
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		return
	}

	q := base.NewQuerySet().From("ids_resetpass").Limit(1)
	q.Where.And("id", c.Params.Get("id"))
	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		return
	}

	expired := rsr[0].Field("expired").TimeParse("datetime") //, base.TimeParse(rsr[0]["expired"].(string), "datetime")
	if expired.Before(time.Now()) {
		return
	}

	c.ViewData["pass_reset_id"] = c.Params.Get("id")
}

func (c Reg) PassResetPutAction() {

	c.AutoRender = false

	var rsp ResponseJson

	rsp.ApiVersion = apiVersion
	rsp.Status = 400
	rsp.Message = "Bad Request"

	defer func() {
		if rspj, err := utils.JsonEncode(rsp); err == nil {
			io.WriteString(c.Response.Out, rspj)
		}
	}()

	if c.Params.Get("id") == "" {
		rsp.Message = "Token can not be null"
		return
	}

	if err := login.PassSetValidate(c.Params); err != nil {
		rsp.Message = err.Error()
		return
	}

	dcn, err := rdo.ClientPull("def")
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	q := base.NewQuerySet().From("ids_resetpass").Limit(1)
	q.Where.And("id", c.Params.Get("id"))
	rsr, err := dcn.Base.Query(q)
	if err != nil || len(rsr) != 1 {
		rsp.Message = "Token not found"
		return
	}

	expired := rsr[0].Field("expired").TimeParse("datetime") // base.TimeParse(rsr[0]["expired"].(string), "datetime")
	if expired.Before(time.Now()) {
		rsp.Message = "Token expired"
		return
	}

	if rsr[0].Field("email").String() != c.Params.Get("email") {
		rsp.Message = "Email or Birthday is not valid"
		return
	}

	q = base.NewQuerySet().From("ids_login").Limit(1)
	q.Where.And("email", c.Params.Get("email"))
	rsl, err := dcn.Base.Query(q)
	if err != nil || len(rsl) != 1 {
		rsp.Message = "User can not found"
		return
	}

	q = base.NewQuerySet().From("ids_profile").Limit(1)
	q.Where.And("uid", rsl[0].Field("uid").Int())
	rspf, err := dcn.Base.Query(q)
	if err != nil || len(rspf) != 1 {
		rsp.Message = "User can not found"
		return
	}
	if fmt.Sprintf("%v", rspf[0].Field("birthday").String()) != c.Params.Get("birthday") {
		rsp.Message = "Email or Birthday is not valid"
		return
	}

	pass, err := pass.HashDefault(c.Params.Get("passwd"))
	if err != nil {
		rsp.Message = "Internal Server Error"
		return
	}

	itemlogin := map[string]interface{}{
		"pass":    pass,
		"updated": base.TimeNow("datetime"),
	}
	ft := base.NewFilter()
	ft.And("uid", rsl[0].Field("uid").Int())
	dcn.Base.Update("ids_login", itemlogin, ft)

	//
	delfr := base.NewFilter()
	delfr.And("id", c.Params.Get("id"))
	dcn.Base.Delete("ids_resetpass", delfr)

	rsp.Status = 200
	rsp.Message = "Successfully Updated"
}
