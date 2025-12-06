package main

import (
	"net/http"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.POST("/auth/createUser", func(ctx *gin.Context) {
		username := ctx.Query("username")
		password := ctx.Query("password")

		if username != "" && password != "" {
			id, err := createUser(username, password)

			if err != nil {

				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			} else {
				ctx.JSON(http.StatusCreated, gin.H{
					"id": id,
				})
				return
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "MISSINGPARAMS",
			})
			return
		}
	})

	router.POST("/auth/", func(ctx *gin.Context) {
		username := ctx.Query("username")
		password := ctx.Query("password")

		if username != "" && password != "" {
			id, err := authenticate(username, password)

			if err != nil {
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": err,
				})
				return
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"id": id,
				})
			}
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "MISSINGPARAMS",
			})
			return
		}
	})

	router.POST("/albums/create", func(ctx *gin.Context) {
		// Bind into a concrete struct to avoid nil dereference
		album := savedAlbum{}
		err := ctx.ShouldBindJSON(&album)

		if album.id != 0 {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "BADPARAM",
			})
			return
		}

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		} else {
			tempAlbum := &savedAlbum{}
			tempAlbum.newAlbum(album.name, album.mood, album.spotifyLinkString, album.creator)
			strs := strings.Split(album.spotifyLinkString, ".")
			id, err := createSavedAlbum(tempAlbum.name, tempAlbum.mood, strs, tempAlbum.creator)

			if err != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": err,
				})
				return
			} else {
				ctx.JSON(http.StatusCreated, gin.H{
					"id": id,
				})
				return
			}
		}
	})

	router.GET("/albums/:username", func(ctx *gin.Context) {
		username := ctx.Param("username")
		albums := getAllAlbums(username)
		ctx.JSON(http.StatusOK, gin.H{
			"albums": albums,
		})
	})

	router.DELETE("/albums/:id", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BADID"})
			return
		}
		if err := deleteSavedAlbum(id); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ctx.Status(http.StatusNoContent)
	})

	type linksPayload struct {
		Links []string `json:"links"`
	}

	router.POST("/albums/:id/addsong", func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "BADID"})
			return
		}

		var payload linksPayload
		if err := ctx.ShouldBindJSON(&payload); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := addSongsToAlbum(id, payload.Links); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		if err := removeSongsFromAlbum(id, payload.Links); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}

		ctx.Status(http.StatusNoContent)
	})

	router.Run(":8080")
}
