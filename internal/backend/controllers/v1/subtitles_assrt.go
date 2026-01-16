package v1

import (
	"net/http"

	"github.com/ChineseSubFinder/ChineseSubFinder/pkg/logic/sub_supplier/assrt"
	"github.com/ChineseSubFinder/ChineseSubFinder/pkg/types/backend"
	"github.com/gin-gonic/gin"
)

func (cb *ControllerBase) AssrtSubtitlesSearch(c *gin.Context) {
	var err error
	defer func() {
		cb.ErrorProcess(c, "AssrtSubtitlesSearch", err)
	}()

	req := backend.ReqAssrtSubtitlesSearch{}
	err = c.ShouldBindJSON(&req)
	if err != nil {
		return
	}

	supplier := assrt.NewSupplier(cb.cronHelper.FileDownloader)
	subInfos, err := supplier.GetSubListFromFile(req.VideoFPath, req.IsMovie)
	if err != nil {
		return
	}

	reply := backend.ReplyAssrtSubtitlesSearch{
		Subtitles: make([]backend.ReplyAssrtSubtitle, 0, len(subInfos)),
	}
	for _, subInfo := range subInfos {
		reply.Subtitles = append(reply.Subtitles, backend.ReplyAssrtSubtitle{
			SubSha256: subInfo.GetUID(),
			Title:     subInfo.Name,
			Language:  int(subInfo.Language),
			Ext:       subInfo.Ext,
			Season:    subInfo.Season,
			Episode:   subInfo.Episode,
		})
	}

	c.JSON(http.StatusOK, reply)
}

func (cb *ControllerBase) AssrtSubtitlesDownload(c *gin.Context) {
	var err error
	defer func() {
		cb.ErrorProcess(c, "AssrtSubtitlesDownload", err)
	}()

	uid := c.Query("uid")
	if uid == "" {
		c.JSON(http.StatusBadRequest, backend.ReplyCommon{Message: "uid is empty"})
		return
	}

	found, subInfo, err := cb.cronHelper.FileDownloader.CacheCenter.DownloadFileGet(uid)
	if err != nil {
		return
	}
	if found == false {
		c.JSON(http.StatusNotFound, backend.ReplyCommon{Message: "subtitle not found"})
		return
	}

	c.Data(http.StatusOK, "application/octet-stream", subInfo.Data)
}
