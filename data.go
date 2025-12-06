package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"math/rand"
	"strings"
	"time"
)

var server = "bensql67.database.windows.net"
var port = 1433
var username = "benjamin"
var password = "Cs2060684#"
var database = "wavedata"

var connString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
	server, username, password, port, database)

var db *sql.DB

func ensureDB() error {
	if db != nil {
		return nil
	}

	opened, err := sql.Open("sqlserver", connString)
	if err != nil {
		return err
	}

	if err := opened.Ping(); err != nil {
		return err
	}

	db = opened
	return nil
}

type user struct {
	id       int
	username string
	password string
}

type savedAlbum struct {
	id                int
	name              string
	mood              string
	spotifyLinkString string
	creator           string
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

func (u *user) newUser(newusername string, newpassword string) {
	hash := sha256.Sum256([]byte(newpassword))

	// use 9 digits to stay within SQL Server INT range
	u.id = generateRandomDigits(9)
	u.username = newusername
	u.password = hex.EncodeToString(hash[:])
}

func (s *savedAlbum) newAlbum(newname string, newmood string, newspotifyLinkString string, newcreator string) {

	s.id = generateRandomDigits(8)
	s.name = newname
	s.mood = newmood
	s.spotifyLinkString = newspotifyLinkString
	s.creator = newcreator
}

func createUser(username string, password string) (int, error) {
	if err := ensureDB(); err != nil {
		return 0, err
	}

	var newuser *user = &user{}
	newuser.newUser(username, password)

	_, err := db.Exec("INSERT INTO users (id, username, password) VALUES (@p1, @p2, @p3)", newuser.id, newuser.username, newuser.password)

	if err != nil {
		return 0, err
	}

	return newuser.id, nil
}

func authenticate(username string, password string) (int, error) {
	if err := ensureDB(); err != nil {
		return 0, err
	}

	var id int

	hash := sha256.Sum256([]byte(password))

	hashedPassword := hex.EncodeToString(hash[:])

	err := db.QueryRow("SELECT id FROM users WHERE username=@p1 AND password=@p2", username, hashedPassword).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, err
	}

	if err != nil {
		return 0, err
	}

	return id, nil
}

func createSavedAlbum(name string, mood string, spotifyLinks []string, creator string) (int, error) {
	if err := ensureDB(); err != nil {
		return 0, err
	}

	spotifyLinkString := ""

	for i := 0; i < len(spotifyLinks); i += 1 {
		if i == len(spotifyLinks)-1 {
			spotifyLinkString += spotifyLinks[i]
		} else {
			spotifyLinkString += spotifyLinks[i] + "."
		}
	}

	newsavedAlbum := &savedAlbum{}
	newsavedAlbum.newAlbum(name, mood, spotifyLinkString, creator)

	_, err := db.Exec("INSERT INTO savedAlbums (id, name, mood, spotifyLinkString, creator) VALUES (@p1, @p2, @p3, @p4, @p5)", newsavedAlbum.id, newsavedAlbum.name, newsavedAlbum.mood, newsavedAlbum.spotifyLinkString, newsavedAlbum.creator)

	if err != nil {
		return 0, err
	}

	return newsavedAlbum.id, nil
}

func createSavedAlbumWithSpotifyLinkStringAndID(id int, name string, mood string, spotifyLinkString string, creator string) (int, error) {
	if err := ensureDB(); err != nil {
		return 0, err
	}

	var newsavedAlbum *savedAlbum = &savedAlbum{}
	newsavedAlbum.newAlbum(name, mood, spotifyLinkString, creator)
	newsavedAlbum.id = id

	_, err := db.Exec("INSERT INTO savedAlbums (id, name, mood, spotifyLinkString, creator) VALUES (@p1, @p2, @p3, @p4, @p5)", newsavedAlbum.id, newsavedAlbum.name, newsavedAlbum.mood, newsavedAlbum.spotifyLinkString, newsavedAlbum.creator)

	if err != nil {
		return 0, nil
	}

	return newsavedAlbum.id, nil
}

func deleteSavedAlbum(id int) error {
	if err := ensureDB(); err != nil {
		return err
	}

	_, err := db.Exec("DELETE FROM savedAlbums WHERE id=@p1", id)

	if err != nil {
		return err
	}

	return nil
}

func addSongsToAlbum(id int, links []string) error {
	if err := ensureDB(); err != nil {
		return err
	}

	var name, mood, spotifyLinkString, creator string

	err := db.QueryRow("SELECT name, mood, spotifyLinkString, creator FROM savedAlbums WHERE id=@p1", id).Scan(&name, &mood, &spotifyLinkString, &creator)

	if err != nil {
		return err
	}

	err = deleteSavedAlbum(id)

	for i := 0; i < len(links); i += 1 {
		if i == len(links)-1 {
			spotifyLinkString += links[i]
		} else {
			spotifyLinkString += links[i] + "."
		}
	}

	if err != nil {
		return err
	}

	var replacementSavedAlbum *savedAlbum = &savedAlbum{}

	replacementSavedAlbum.newAlbum(name, mood, spotifyLinkString, creator)
	replacementSavedAlbum.id = id

	_, err = createSavedAlbumWithSpotifyLinkStringAndID(replacementSavedAlbum.id, replacementSavedAlbum.name, replacementSavedAlbum.mood, replacementSavedAlbum.spotifyLinkString, replacementSavedAlbum.creator)

	if err != nil {
		return err
	}

	return nil
}

func removeSongsFromAlbum(id int, links []string) error {
	if err := ensureDB(); err != nil {
		return err
	}

	var name, mood, spotifyLinkString, creator string

	err := db.QueryRow("SELECT name, mood, spotifyLinkString, creator FROM savedAlbums WHERE id=@p1", id).Scan(&name, &mood, &spotifyLinkString, &creator)

	if err != nil {
		return err
	}

	err = deleteSavedAlbum(id)

	if err != nil {
		return err
	}
	var result string
	for i := 0; i < len(links); i += 1 {
		result = strings.ReplaceAll(spotifyLinkString, links[i]+".", "")
	}

	var replacementSavedAlbum *savedAlbum = &savedAlbum{}

	replacementSavedAlbum.newAlbum(name, mood, result, creator)
	replacementSavedAlbum.id = id

	_, err = createSavedAlbumWithSpotifyLinkStringAndID(replacementSavedAlbum.id, replacementSavedAlbum.name, replacementSavedAlbum.mood, replacementSavedAlbum.spotifyLinkString, replacementSavedAlbum.creator)

	if err != nil {
		return err
	}

	return nil
}


func getAllAlbums(username string) []savedAlbum {
	if err := ensureDB(); err != nil {
		return []savedAlbum{}
	}

	rows, err := db.Query("SELECT id, name, mood, spotifyLinkString, creator FROM savedAlbums WHERE creator=@p1", username)
	if err != nil {
		return []savedAlbum{}
	}
	defer rows.Close()

	var albums []savedAlbum

	for rows.Next() {
		var a savedAlbum
		if err := rows.Scan(&a.id, &a.name, &a.mood, &a.spotifyLinkString, &a.creator); err != nil {
			return []savedAlbum{}
		}
		albums = append(albums, a)
	}

	return albums
}
