package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var server = "bensql67.database.windows.net"
var port = 1433
var username = "benjamin"
var password = os.Getenv("DB_PASSWORD")
var database = "Cs2060684#"

var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, username, password, port, database)

var db, _ = sql.Open("sqlserver", connString)

type user struct {
	id int
	username string
	password string
}

type savedAlbum struct {
	id           int
	name         string
	mood         string
	spotifyLinkString string
	creator string
}

func dbInit() {

}

func generateRandomDigits(digits int) int {
    if digits <= 0 {
        return 0 
    }

    rand.Seed(time.Now().UnixNano())

    min := intPow(10, digits-1)
    max := intPow(10, digits) - 1

    return rand.Intn(max-min+1) + min
}

func intPow(base, exp int) int {
    result := 1
    for i := 0; i < exp; i++ {
        result *= base
    }
    return result
}

func (u *user) newUser(username string, password string) user {
	hash := sha256.Sum256([]byte(password))
	return user{
		id: generateRandomDigits(15),
		username: username,
		password: hex.EncodeToString((hash[:])),
	}
}

func (s *savedAlbum) newAlbum(name string, mood string, spotifyLinkString string, creator string) savedAlbum {
	
	return savedAlbum{
		id:           generateRandomDigits(8),
		name:         name,
		mood:         mood,
		spotifyLinkString: spotifyLinkString,
		creator: creator,
	}
}


func createUser(username string, password string) (int, error) {
	var newuser *user
	newuser.newUser(username, password)

	_, err := db.Exec("INSERT INTO users (id, username, password) VALUES (?, ?, ?)", newuser.id, newuser.username, newuser.password)

	if err != nil {
		return 0, err
	}

	return newuser.id, nil
}

func authenticate(username string, password string) (int, error) {
	var id int;

	err := db.QueryRow("SELECT * FROM users WHERE username=? AND password=?", username, password).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func createSavedAlbum(name string, mood string, spotifyLinks []string, creator string) (int, error) {
	spotifyLinkString := ""
	
	for i := 0; i < len(spotifyLinks); i += 1 {
		if i == len(spotifyLinks) - 1 {
			spotifyLinkString += spotifyLinks[i]
		} else {
			spotifyLinkString += spotifyLinks[i] + "."
		}
	}

	var newsavedAlbum *savedAlbum

	_, err := db.Exec("INSERT INTO savedAlbums (id, name, mood, spotifyLinkString, creator) VALUES (?, ?, ?, ?, ?)", newsavedAlbum.id, newsavedAlbum.name, newsavedAlbum.mood, newsavedAlbum.spotifyLinkString, newsavedAlbum.creator)

	if  err != nil {
		return 0, err
	}

	return newsavedAlbum.id, nil
}

func createSavedAlbumWithSpotifyLinkStringAndID(id int, name string, mood string, spotifyLinkString string, creator string,) (int, error) {
	var newsavedAlbum *savedAlbum
	newsavedAlbum.newAlbum(name, mood, spotifyLinkString, creator)
	newsavedAlbum.id = id

	_, err := db.Exec("INSERT INTO savedAlbums (id, name, mood, spotifyLinkString, creator) VALUES (?, ?, ?, ?, ?)", newsavedAlbum.id, newsavedAlbum.name, newsavedAlbum.mood, newsavedAlbum.spotifyLinkString, newsavedAlbum.creator)

	if err != nil {
		return 0, nil
	}

	return newsavedAlbum.id, nil
}


func deleteSavedAlbum(id int) error {
	_, err := db.Exec("DELETE FROM savedAlbums WHERE id=?", id)

	if err != nil {
		return err
	}

	return nil;
}

func addSongsToAlbum(id int, links []string) error {
	var name, mood, spotifyLinkString, creator string

	err := db.QueryRow("SELECT * FROM savedAlbums WHERE id=?", id).Scan(&name, &mood, &spotifyLinkString, &creator)

	if err != nil {
		return err
	}
	
	err = deleteSavedAlbum(id)

	
	for i := 0; i < len(links); i += 1 {
		if i == len(links) - 1 {
			spotifyLinkString += links[i]
		} else {
			spotifyLinkString += links[i] + "."
		}
	}

	if err != nil {
		return err
	}

	var replacementSavedAlbum *savedAlbum

	replacementSavedAlbum.newAlbum(name, mood, spotifyLinkString, creator)
	replacementSavedAlbum.id = id

	_, err = createSavedAlbumWithSpotifyLinkStringAndID(replacementSavedAlbum.id, replacementSavedAlbum.name, replacementSavedAlbum.mood, replacementSavedAlbum.spotifyLinkString, replacementSavedAlbum.creator)

	if err != nil {
		return err
	}

	return nil
}

func removeSongsFromAlbum(id int, links []string) error {
	var name, mood, spotifyLinkString, creator string

	err := db.QueryRow("SELECT * FROM savedAlbums WHERE id=?", id).Scan(&name, &mood, &spotifyLinkString, &creator)

	if err != nil {
		return err 
	}

	err = deleteSavedAlbum(id)

	if err != nil {
		return err 
	}
	var result string
	for i := 0; i < len(links); i += 1 {
		result = strings.ReplaceAll(spotifyLinkString, links[i] + ".", "")
	}

	var replacementSavedAlbum *savedAlbum;

	replacementSavedAlbum.newAlbum(name, mood, result, creator)
	replacementSavedAlbum.id = id
	
	_, err = createSavedAlbumWithSpotifyLinkStringAndID(replacementSavedAlbum.id, replacementSavedAlbum.name, replacementSavedAlbum.mood, replacementSavedAlbum.spotifyLinkString, replacementSavedAlbum.creator)
	
	if err != nil {
		return err
	}

	return nil
}