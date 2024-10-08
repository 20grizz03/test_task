package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	Song "test_task/models"
)

var db *gorm.DB

// Получение всех песен с фильтрацией и пагинацией
func GetSongs(w http.ResponseWriter, r *http.Request) {
	group := r.URL.Query().Get("group")
	song := r.URL.Query().Get("song")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	var songs []Song.Song
	query := db.Limit(limit).Offset((page - 1) * limit)

	if group != "" {
		query = query.Where("group = ?", group)
	}
	if song != "" {
		query = query.Where("song = ?", song)
	}

	query.Find(&songs)
	json.NewEncoder(w).Encode(songs)
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// Функция для получения детализированной информации о песне из внешнего API
func FetchSongDetails(group, song string) (SongDetails, error) {
	var details SongDetails

	// Создаем URL для запроса к API с параметрами "group" и "song"
	apiURL := "http://example.com/info" //  реальный адрес вашего API
	params := url.Values{}
	params.Add("group", group)
	params.Add("song", song)

	// Формируем полный URL с параметрами
	reqURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Выполняем GET-запрос к API
	resp, err := http.Get(reqURL)
	if err != nil {
		return details, fmt.Errorf("failed to fetch song details: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что API вернул успешный ответ
	if resp.StatusCode != http.StatusOK {
		return details, fmt.Errorf("API returned non-OK status: %s", resp.Status)
	}

	// Читаем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return details, fmt.Errorf("failed to read API response: %v", err)
	}

	// Декодируем JSON-ответ в структуру SongDetails
	err = json.Unmarshal(body, &details)
	if err != nil {
		return details, fmt.Errorf("failed to parse API response: %v", err)
	}

	// Возвращаем детализированную информацию о песне
	return details, nil
}

// Добавление новой песни
func CreateSong(w http.ResponseWriter, r *http.Request) {
	var newSong Song.Song
	json.NewDecoder(r.Body).Decode(&newSong)

	// Получение информации через внешний API
	details, err := FetchSongDetails(newSong.Group, newSong.Song)
	if err != nil {
		http.Error(w, "Failed to fetch song details", http.StatusInternalServerError)
		return
	}

	newSong.ReleaseDate = details.ReleaseDate
	newSong.Text = details.Text
	newSong.Link = details.Link

	db.Create(&newSong)
	json.NewEncoder(w).Encode(newSong)
}

// Получение текста песни с пагинацией по куплетам
func GetSongLyrics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))

	var song Song.Song
	db.First(&song, id)

	verses := paginateLyrics(song.Text, page)
	json.NewEncoder(w).Encode(verses)
}

// Обновление данных песни
func UpdateSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song Song.Song
	db.First(&song, id)

	json.NewDecoder(r.Body).Decode(&song)
	db.Save(&song)

	json.NewEncoder(w).Encode(song)
}

// Удаление песни
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var song Song.Song
	db.Delete(&song, id)

	w.WriteHeader(http.StatusNoContent)
}

// Функция для пагинации текста песни
func paginateLyrics(text string, page int) []string {
	verses := strings.Split(text, "\n\n")
	limit := 2
	start := (page - 1) * limit
	end := start + limit

	if start >= len(verses) {
		return []string{}
	}
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end]
}
