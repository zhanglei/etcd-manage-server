package server

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/etcd-manage/etcd-manage-server/program/models"
	"github.com/etcd-manage/etcdsdk/etcdv3"
	"github.com/etcd-manage/etcdsdk/model"
	"github.com/gin-gonic/gin"
)

// ServerController etcd服务列表相关操作
type ServerController struct {
}

// List 获取etcd服务列表，全部
func (api *ServerController) List(c *gin.Context) {
	list, err := new(models.EtcdServersModel).All()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusOK, list)
}

// Restore 修复v1版本或e3w对目录的标记
func (api *ServerController) Restore(c *gin.Context) {
	etcdId := c.Query("etcd_id")
	var err error
	defer func() {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()

	etcdIdNum, _ := strconv.Atoi(etcdId)
	etcdOne := new(models.EtcdServersModel)
	etcdOne, err = etcdOne.FirstById(int32(etcdIdNum))
	if err != nil {
		return
	}
	if etcdOne.Version != model.ETCD_VERSION_V3 {
		err = errors.New("Only V3 version is allowed to be repaired")
		return
	}
	// 连接etcd
	cfg := &model.Config{
		Version:   etcdOne.Version,
		Address:   strings.Split(etcdOne.Address, ","),
		TlsEnable: etcdOne.TlsEnable == "true",
		CertFile:  etcdOne.CaFile,
		KeyFile:   etcdOne.KeyFile,
		CaFile:    etcdOne.CaFile,
		Username:  etcdOne.Username,
		Password:  etcdOne.Password,
	}
	client, err := etcdv3.NewClient(cfg)
	if err != nil {
		return
	}
	clientV3, ok := client.(*etcdv3.EtcdV3Sdk)
	if ok == false {
		err = errors.New("Connecting etcd V3 service error")
		return
	}
	err = clientV3.Restore()
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, "ok")
}