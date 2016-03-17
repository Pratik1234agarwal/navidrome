package api

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/deluan/gosonic/api/responses"
	"github.com/deluan/gosonic/engine"
	"github.com/deluan/gosonic/utils"
)

type MediaAnnotationController struct {
	BaseAPIController
	scrobbler engine.Scrobbler
}

func (c *MediaAnnotationController) Prepare() {
	utils.ResolveDependencies(&c.scrobbler)
}

func (c *MediaAnnotationController) Scrobble() {
	id := c.RequiredParamString("id", "Required id parameter is missing")
	time := c.ParamTime("time", time.Now())
	submission := c.ParamBool("submission", false)
	if submission {
		mf, err := c.scrobbler.Register(id, time, true)
		if err != nil {
			beego.Error("Error scrobbling:", err)
			c.SendError(responses.ERROR_GENERIC, "Internal error")
		}
		beego.Info(fmt.Sprintf(`Scrobbled (%s) "%s" at %v`, id, mf.Title, time))
	}

	response := c.NewEmpty()
	c.SendResponse(response)
}