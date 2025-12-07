package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type createAlbumRequest struct {
	Name    string   `json:"name" binding:"required"`
	Mood    string   `json:"mood" binding:"required"`
	Links   []string `json:"links" binding:"required"`
	Creator string   `json:"creator" binding:"required"`
}

type linksPayload struct {
	Links []string `json:"links" binding:"required"`
}

func main() {
	router := gin.Default()

	router.POST("/auth/createUser", func(ctx *gin.Context) {
		var req authRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "MISSINGPARAMS"})
			return
		}

		id, err := createUser(req.Username, req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	})

	router.POST("/auth", func(ctx *gin.Context) {
		var req authRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "MISSINGPARAMS"})
			return
		}

		id, err := authenticate(req.Username, req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"id": id})
	})

	router.POST("/albums/create", func(ctx *gin.Context) {
		var req createAlbumRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := createSavedAlbum(req.Name, req.Mood, req.Links, req.Creator)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"id": id})
	})

	router.GET("/albums/:username", func(ctx *gin.Context) {
		username := ctx.Param("username")
		albums, err := getAllAlbums(username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"albums": albums})
	})

	router.DELETE("/albums/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BADID"})
			return
		}
		if err := deleteSavedAlbum(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.Status(http.StatusNoContent)
	})

	router.POST("/albums/:id/addsong", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BADID"})
			return
		}

		var payload linksPayload
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := addSongsToAlbum(id, payload.Links); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Status(http.StatusNoContent)
	})

	router.DELETE("/albums/:id/removesong", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BADID"})
			return
		}

		var payload linksPayload
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := removeSongsFromAlbum(id, payload.Links); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.Status(http.StatusNoContent)
	})

	router.Run(":8080")
}
