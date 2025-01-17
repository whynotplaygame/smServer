package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"smServer/constant"
	"smServer/server/common"
	"smServer/server/web/logic"
	"smServer/server/web/model"
)

var DefaultAccountController = &AccountController{}

type AccountController struct {
}

func (a *AccountController) Register(ctx *gin.Context) {
	fmt.Println("AccountController Register")
	/*
		1, 获取请求参数
		2，根据用户名 查询数据是否存在，有，用户名已存在，没有就注册
		3，告诉起那段，注册成功即可

	*/
	rq := &model.RegisterReq{}
	err := ctx.ShouldBind(rq)
	if err != nil {
		log.Println("参数格式不合法", err)
		ctx.JSON(http.StatusOK, common.Error(constant.InvalidParam, "参数不合法"))
		return
	}
	// 一般web服务，错误格式自定义
	err = logic.DefaultAccount.Register(rq)
	if err != nil {
		log.Println("注册业务出错", err)
		ctx.JSON(http.StatusOK, common.Error(err.(*common.MyError).Code(), err.Error()))
		return
	}
	// 成功返回ok
	ctx.JSON(http.StatusOK, common.Success(constant.OK, nil))

}
