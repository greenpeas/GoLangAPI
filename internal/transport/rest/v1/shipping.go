package v1

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"seal/internal/domain/shipping"
	"seal/internal/transport"
	"seal/pkg/app_error"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// List for swagger only
type shippingList struct {
	RecordsTotal    int                 `json:"records_total"`
	RecordsFiltered int                 `json:"records_filtered"`
	Data            []shipping.Shipping `json:"data"`
}

func (h *Handler) registerShippingHandler(api *gin.RouterGroup) {
	group := api.Group("/shipping")
	{
		group.GET(":id", h.shipping)
		group.GET("modem/:imei", h.shippingByModemImei)
		group.GET(":id/route", h.shippingRoute)
		group.POST(":id/files/:type", h.shippingUploadFiles)
		group.GET(":id/files/:name", h.shippingDownloadFile)
		group.DELETE(":id/files/:name", h.shippingDeleteFile)
		group.PUT(":id/files/:name", h.shippingUpdateFile)
		group.PUT(":id/modem/:imei", h.shippingModemSet)
		group.PUT(":id", h.shippingUpdate)
		group.PUT(":id/start", h.shippingStart)
		group.GET(":id/coordinates", h.shippingCoordinates)
		group.GET(":id/telemetry", h.shippingTelemetry)
		group.PUT(":id/end", h.shippingEnd)
		group.DELETE(":id", h.shippingDelete)
		group.GET("", h.shippingList)
		group.POST("", h.shippingCreate)
	}
}

// ItemShipping godoc
// @Summary      Shipping info
// @Description  shipping info
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id} [get]
// @Security 	 BearerAuth
func (h *Handler) shipping(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Shipping.GetById(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ShippingRoute godoc
// @Summary      Shipping route info
// @Description  shipping route info
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	route.Route
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/route [get]
// @Security 	 BearerAuth
func (h *Handler) shippingRoute(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Shipping.Route(id); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ListShipping godoc
// @Summary      List shipping
// @Description  get shipping
// @Tags         shipping
// @Accept       json
// @Param        find    	  query     string  false  "search string"
// @Param        find_type    query     int     false  "search type (0 - '=', 1 - 'like', 2 = 'ilike')"	Enums(0, 1, 2)
// @Param        limit        query     int     false  "limit"	minimum(0)	maximum (100)
// @Param        offset       query     int     false  "offset"	minimum(0)	maximum (32767)
// @Param        status       query     []string  false  "status"	collectionFormat(multi)	enums(0,1,2)
// @Success      200	{object}	shippingList
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping [get]
// @Security 	 BearerAuth
func (h *Handler) shippingList(c *gin.Context) {
	var queryParams shipping.QueryParams
	if err := c.ShouldBind(&queryParams); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if list, err := h.Usecase.Shipping.List(queryParams); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, list)
	}
}

// CreateShipping godoc
// @Summary      Create shipping
// @Description  add shipping
// @Tags         shipping
// @Accept       json
// @Produce      json
// @Param		 data	body	shipping.CreateRequest	true	"data"
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping [post]
// @Security 	 BearerAuth
func (h *Handler) shippingCreate(c *gin.Context) {
	var fromRequest shipping.CreateRequest
	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	userId := c.GetInt("userId")
	if userId == 0 {
		c.Error(app_error.InternalServerError(errors.New("can't get user")))
		return
	}

	if data, err := h.Usecase.Shipping.Create(fromRequest, userId); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// UpdateShipping godoc
// @Summary      Update shipping
// @Description  update shipping
// @Tags         shipping
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Param		 data	body	shipping.UpdateRequest	true	"data"
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id} [put]
// @Security 	 BearerAuth
func (h *Handler) shippingUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest shipping.UpdateRequest

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if data, err := h.Usecase.Shipping.Update(id, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// StartShipping godoc
// @Summary      Start shipping
// @Description  Start shipping
// @Tags         shipping
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/start [put]
// @Security 	 BearerAuth
func (h *Handler) shippingStart(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var shipping shipping.Db

	if shipping, err = h.Usecase.Shipping.GetDbById(id); err != nil {
		c.Error(app_error.ErrNotFound)
		return
	}

	if data, err := h.Usecase.Shipping.Start(shipping); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// StartShipping godoc
// @Summary      End shipping
// @Description  End shipping
// @Tags         shipping
// @Accept       json
// @Produce      json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/end [put]
// @Security 	 BearerAuth
func (h *Handler) shippingEnd(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var shipping shipping.Db

	if shipping, err = h.Usecase.Shipping.GetDbById(id); err != nil {
		c.Error(app_error.ErrNotFound)
		return
	}

	if data, err := h.Usecase.Shipping.End(shipping); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// DeleteShipping godoc
// @Summary      Shipping delete
// @Description  shipping delete
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     false  "id"	minimum(0)	maximum (32767)
// @Success      200	{object}	transport.DeleteResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id} [delete]
// @Security 	 BearerAuth
func (h *Handler) shippingDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Shipping.DeleteById(id); err != nil {
		c.Error(err)
	} else if data {
		c.JSON(http.StatusOK, transport.DeleteResponse{Success: true})
	} else {
		c.Error(app_error.ErrNotFound)
	}
}

// UploadShippingFiles godoc
// @Summary      Upload shipping files
// @Description  upload shipping files
// @Tags         shipping
// @Accept       multipart/form-data
// @Param        id       path     int     true  "id"		minimum(0)	maximum (32767)
// @Param        type     path     int     true  "type" 	minimum(0)	maximum (3)
// @Param        seal_id  formData int     false "seal_id"
// @Param        comment  formData string  false "comment"
// @Param        file   formData    []file true  "upload files"
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/files/{type} [post]
// @Security 	 BearerAuth
func (h *Handler) shippingUploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()

	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	files := form.File["file"]

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	fileType, err := strconv.Atoi(c.Param("type"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	var fromRequest shipping.AddFileRequest

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	filesDirectory := h.Usecase.Shipping.GetFilesDirectory(id)

	var shippingFiles []shipping.File
	for _, file := range files {
		fileName, filePath, err := getFileName(file.Filename, filesDirectory, 0)
		if err != nil {
			c.Error(err)
			return
		}

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			h.Usecase.Shipping.RemoveFilesFromDisk(id, shippingFiles)
			c.Error(err)
			return
		} else {
			checksum, err := crc(file)
			if err != nil {
				h.Logger.Error(fmt.Sprintf("Error calculate crc: %v", err))
			}

			shippingFiles = append(shippingFiles, shipping.File{
				Title:    file.Filename,
				Name:     fileName,
				Type:     fileType,
				SealId:   fromRequest.SealId,
				Comment:  fromRequest.Comment,
				Checksum: checksum,
			})
		}
	}

	if data, err := h.Usecase.Shipping.AddFiles(id, shippingFiles); err != nil {
		h.Usecase.Shipping.RemoveFilesFromDisk(id, shippingFiles)
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

func getFileName(fileName, filesDirectory string, i int) (string, string, error) {
	now := time.Now().UnixNano()
	ext := filepath.Ext(fileName)
	name := fmt.Sprintf("%d%s", now, ext)

	if i > 0 {
		name = fmt.Sprintf("%d_%d%s", now, i, ext)
	}

	filePath := fmt.Sprintf("%s/%s", filesDirectory, name)

	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		return name, filePath, nil
	} else if i > 500 {
		return "", "", fmt.Errorf("file exists %s", filePath)
	}

	return getFileName(fileName, filesDirectory, i+1)
}

func crc(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()

	if err != nil {
		return "", err
	}

	defer f.Close()

	crc := md5.New()

	if _, err := io.Copy(crc, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(crc.Sum(nil)), nil
}

// ShippingDownloadFile godoc
// @Summary      Download shipping file
// @Description  download shipping file
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     true  "id"		minimum(0)	maximum (32767)
// @Param        name     path     string  true  "name"
// @Success      200	{file} 		binary
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/files/{name} [get]
// @Security 	 BearerAuth
func (h *Handler) shippingDownloadFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	fileName := c.Param("name")

	if fileName == "" {
		c.Error(app_error.ErrNotFound)
		return
	}

	fileInfo, err := h.Usecase.Shipping.GetFileInfo(id, fileName)
	if err != nil {
		c.Error(err)
		return
	}

	filesDirectory := h.Usecase.Shipping.GetFilesDirectory(id)
	filePath := fmt.Sprintf("%s/%s", filesDirectory, fileName)
	c.FileAttachment(filePath, fileInfo.Title)
}

// ShippingDeleteFile godoc
// @Summary      Delete shipping file
// @Description  delete shipping file
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     true  "id"		minimum(0)	maximum (32767)
// @Param        name     path     string  true  "name"
// @Success      200	{object}	transport.SuccessResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/files/{name} [delete]
// @Security 	 BearerAuth
func (h *Handler) shippingDeleteFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	fileName := c.Param("name")

	if fileName == "" {
		c.Error(app_error.ErrNotFound)
		return
	}

	if err := h.Usecase.Shipping.RemoveFile(id, fileName); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, transport.SuccessResponse{Success: true})
}

// UpdateShippingFileInfo godoc
// @Summary      Update shipping file info
// @Description  Update shipping file info
// @Tags         shipping
// @Accept       json
// @Produce      json
// @Param        id     path     	int     true  "id"		minimum(0)	maximum (32767)
// @Param        name   path 		string  true  "name"
// @Param		 data	body	    shipping.UpdateFileRequest	true	"data"
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/files/{name} [put]
// @Security 	 BearerAuth
func (h *Handler) shippingUpdateFile(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	fileName := c.Param("name")

	if fileName == "" {
		c.Error(app_error.ErrNotFound)
		return
	}

	var fromRequest shipping.UpdateFileRequest

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if data, err := h.Usecase.Shipping.UpdateFileInfo(id, fileName, fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ShippingModemSet godoc
// @Summary      Add modem to shipping (by id or imei)
// @Description  add modem to shipping (by id or imei)
// @Tags         shipping
// @Accept       json
// @Param        id       path     int     true  "id"		minimum(0)	maximum (32767)
// @Param        imei   path     int     true  "imei or id" 	minimum(0)	maximum (999999999999999)
// @Success      200	{object}	transport.SuccessResponse
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/modem/{imei} [put]
// @Security 	 BearerAuth
func (h *Handler) shippingModemSet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	imei, err := strconv.ParseUint(c.Param("imei"), 10, 64)
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if imei > math.MaxInt32 {
		if data, err := h.Usecase.Shipping.SetModemByImei(id, imei); err != nil {
			c.Error(err)
		} else {
			c.JSON(http.StatusOK, data)
		}

		return
	}

	if data, err := h.Usecase.Shipping.SetModemById(id, int(imei)); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, transport.SuccessResponse{Success: data})
	}
}

// ShippingCoordinates godoc
// @Summary      Shipping coordinates
// @Description  shipping coordinates
// @Tags         shipping
// @Accept       json
// @Param        id          path    int     true  "id"		minimum(0)	maximum (32767)
// @Param        from	     query	 string	 false "from"
// @Param        limit       query	 int	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]shipping.trackResponseCoordinate
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/coordinates [get]
// @Security 	 BearerAuth
func (h *Handler) shippingCoordinates(c *gin.Context) {
	var fromRequest shipping.TrackQueryParams
	var err error

	if fromRequest.Id, err = strconv.Atoi(c.Param("id")); err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if data, err := h.Usecase.Shipping.Coordinates(fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ShippingTelemetry godoc
// @Summary      Shipping telemetry
// @Description  shipping telemetry
// @Tags         shipping
// @Accept       json
// @Param        id          path    int     true  "id"		minimum(0)	maximum (32767)
// @Param        from	     query	 string	 false "from"
// @Param        limit       query	 int	 false "limit"
// @Param        order_desc  query	 bool	 false "order_desc"
// @Success      200	{object}	[]shipping.trackResponseTelemetry
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/{id}/telemetry [get]
// @Security 	 BearerAuth
func (h *Handler) shippingTelemetry(c *gin.Context) {
	var fromRequest shipping.TrackQueryParams
	var err error

	if fromRequest.Id, err = strconv.Atoi(c.Param("id")); err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if err := c.ShouldBind(&fromRequest); err != nil {
		h.Logger.Debug(err.Error())
		c.Error(app_error.BadRequestError(err))
		return
	}

	if errs := h.Validator.Struct(fromRequest); errs != nil {
		h.Logger.Debug("Ошибки валидации", errs)
		c.Error(app_error.ValidationError(errs))
		return
	}

	if data, err := h.Usecase.Shipping.Telemetry(fromRequest); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}

// ActiveShippingByModemImei godoc
// @Summary      Get active shipping by modem imei
// @Description  get active shipping by modem imei
// @Tags         shipping
// @Accept       json
// @Param        imei   path     int     true  "imei" 	minimum(0)	maximum (999999999999999)
// @Success      200	{object}	shipping.Shipping
// @Failure      400	{object}	app_error.AppError
// @Failure      401	{object}	app_error.AppError
// @Failure      422	{object}	app_error.AppError
// @Failure      500	{object}	app_error.AppError
// @Router       /shipping/modem/{imei} [get]
// @Security 	 BearerAuth
func (h *Handler) shippingByModemImei(c *gin.Context) {
	imei, err := strconv.ParseUint(c.Param("imei"), 10, 64)
	if err != nil {
		h.Logger.Error(err.Error())
		c.Error(err)
		return
	}

	if data, err := h.Usecase.Shipping.GetActiveByModemImei(imei); err != nil {
		c.Error(err)
	} else {
		c.JSON(http.StatusOK, data)
	}
}
