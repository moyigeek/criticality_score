package controller

import (
	"slices"
	"strconv"

	"github.com/HUSTSecLab/criticality_score/cmd/apiserver/internal/model"
	"github.com/HUSTSecLab/criticality_score/pkg/logger"
	"github.com/HUSTSecLab/criticality_score/pkg/storage"
	"github.com/HUSTSecLab/criticality_score/pkg/storage/repository"
	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

// @Summary Search score results by git link
// @Description Search score results by git link
// @Description NOTE: All details are ignored, should use /results/:scoreid to get details
// @Description NOTE: Maxium take count is 1000
// @Accept json
// @Produce json
// @Success 200 {object} model.PageDTO[model.ResultDTO]
// @Router /results [get]
// @Param q query string true "Search query"
// @Param start query int false "Skip count"
// @Param take query int false "Take count"
func resultsHandler(c *gin.Context) {
	r := repository.NewResultRepository(storage.GetDefaultAppDatabaseContext())

	type query struct {
		Search string `form:"q"`
		Skip   int    `form:"start"`
		Take   int    `form:"take"`
	}

	var q query = query{
		Skip: 0,
		Take: 100,
	}

	if err := c.ShouldBindQuery(&q); err != nil {
		c.JSON(400, "Invalid query parameters")
		return
	}

	if q.Take > 1000 {
		q.Take = 1000
	}

	if q.Search == "" {
		c.JSON(400, "Invalid query parameters")
		return
	}

	cnt, err := r.CountByLink(q.Search)
	if err != nil {
		c.JSON(500, "Error occurred when counting results")
		return
	}

	result, err := r.QueryByLink(q.Search, q.Skip, q.Take)
	if err != nil {
		c.JSON(500, "Error occurred when querying results")
		return
	}

	var items []model.ResultDTO = make([]model.ResultDTO, 0)

	for v := range result {
		items = append(items, *model.ResultDOToDTO(v))
	}

	c.JSON(200, model.NewPageDTO(cnt, q.Skip, items))
}

// @Summary Get score results
// @Description Get score results, including all details by scoreid
// @Accept json
// @Produce json
// @Success 200 {object} model.ResultDTO
// @Router /results/{scoreid} [get]
// @Param scoreid path int true "Score ID"
func resultHandler(c *gin.Context) {
	r := repository.NewResultRepository(storage.GetDefaultAppDatabaseContext())

	scoreidStr := c.Param("scoreid")
	scoreid, err := strconv.Atoi(scoreidStr)

	if err != nil {
		c.JSON(400, "Invalid query parameters")
		return
	}

	result, err := r.GetByScoreID(scoreid)
	if err != nil || result == nil {
		logger.Error("Error occurred when querying result", err)
		c.JSON(500, "Error occurred when querying result")
		return
	}

	gitDetails, err := r.QueryGitDetailsByScoreID(scoreid)
	if err != nil {
		logger.Error("Error occurred when querying git details", err)
		c.JSON(500, "Error occurred when querying git details")
		return
	}

	langDetails, err := r.QueryLangDetailsByScoreID(scoreid)
	if err != nil {
		logger.Error("Error occurred when querying lang details", err)
		c.JSON(500, "Error occurred when querying lang details")
		return
	}

	distDetails, err := r.QueryDistDetailsByScoreID(scoreid)
	if err != nil {
		logger.Error("Error occurred when querying dist details", err)
		c.JSON(500, "Error occurred when querying dist details")
		return
	}

	ret := model.ResultDOToDTO(result)
	ret.GitDetail = lo.Map(slices.Collect(gitDetails), func(v *repository.ResultGitDetail, i int) model.ResultGitMetadataDTO {
		return *model.ResultGitDetailDOToDTO(v)
	})
	ret.LangDetail = lo.Map(slices.Collect(langDetails), func(v *repository.ResultLangDetail, i int) model.ResultLangDetailDTO {
		return *model.ResultLangDetailDOToDTO(v)
	})
	ret.DistDetail = lo.Map(slices.Collect(distDetails), func(v *repository.ResultDistDetail, i int) model.ResultDistDetailDTO {
		return *model.ResultDistDetailDOToDTO(v)
	})

	c.JSON(200, ret)
}

func registResult(e gin.IRouter) {
	e.GET("/results", resultsHandler)
	e.GET("/results/:scoreid", resultHandler)
}
