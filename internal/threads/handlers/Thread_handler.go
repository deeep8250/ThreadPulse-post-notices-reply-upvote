package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	service "github.com/threadpulse/internal/threads/services"
	"github.com/threadpulse/models"
)

type ThreadHandler struct {
	services service.ServiceInterface
}

func NewThreadHandler(serv service.ServiceInterface) *ThreadHandler {
	return &ThreadHandler{services: serv}
}

func (h *ThreadHandler) CreateThreadHandler(c *gin.Context) {
	var Thread models.CreateThread
	err := c.ShouldBindJSON(&Thread)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized user",
		})
		return
	}

	err = h.services.CreateThread(userID.(int), Thread)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "something went wrong while creating the thread",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"response": "thread created",
	})

}

func (h *ThreadHandler) GetAllThreadHandler(c *gin.Context) {
	limit := c.DefaultQuery("limit", "10")
	LimitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}
	page := c.DefaultQuery("page", "1")
	PageInt, err := strconv.Atoi(page)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	threads, count, err := h.services.GetAllThreads(PageInt, LimitInt)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"threads": threads,
		"count":   count,
	})

}

func (h *ThreadHandler) GetThreadByIdHandler(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	Thread, err := h.services.GetThreadById(idInt)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": Thread,
	})

}

func (h *ThreadHandler) UpdateThreadHandler(c *gin.Context) {

	id := c.Param("id")
	IdInt, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	var input models.UpdateThread

	err = c.ShouldBindJSON(&input)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.Error(err)
		c.Abort()
		return
	}

	err = h.services.UpdateThread(IdInt, userID.(int), input)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "thread updated",
	})

}

func (h *ThreadHandler) DeleteThreadHandler(c *gin.Context) {
	id := c.Param("id")
	IdInt, err := strconv.Atoi(id)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	userID, ok := c.Get("userID")
	if !ok {
		c.Error(err)
		c.Abort()
		return
	}
	err = h.services.DeleteThread(IdInt, userID.(int))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"response": "thread deleted",
	})

}

func (h *ThreadHandler) GetHotThreads(c *gin.Context) {
	limit := c.DefaultQuery("limit", "2")
	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ctx := c.Request.Context()

	hotThreads, err := h.services.GetHotThreadsService(ctx, limitInt)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hot threads": hotThreads,
	})

}
