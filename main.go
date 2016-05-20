package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

func IsUsernameValid(username string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]*$")
	return re.MatchString(username)
}

func GetDateData(username string) (map[string]int, []string, error) {

	dates := make([]string, 0)
	dateCounts := make(map[string]int)
	if !IsUsernameValid(username) {
		err := errors.New("Username not valid")
		return dateCounts, dates, err
	}

	resp, err := http.Get("https://github.com/" + username)
	if err != nil {
		return dateCounts, dates, err
	}
	root, err := html.Parse(resp.Body)
	if err != nil {
		return dateCounts, dates, err
	}

	matcher := func(n *html.Node) bool {
		return scrape.Attr(n, "class") == "day"
	}
	dayNodes := scrape.FindAllNested(root, matcher)
	for _, dayNode := range dayNodes {
		date := scrape.Attr(dayNode, "data-date")
		count, err := strconv.Atoi(scrape.Attr(dayNode, "data-count"))
		if err != nil {
			count = 0
		}
		dates = append(dates, date)
		dateCounts[date] = count
	}
	sort.Strings(dates)
	return dateCounts, dates, nil
}

func GetCurrentStreak(data map[string]int, dates []string) (int, []string) {
	streakDates := make([]string, 0)
	count := 0
	for i := len(dates) - 1; i >= 0; i-- {
		d := dates[i]
		if data[d] == 0 {
			break
		}
		count = count + 1
		streakDates = append(streakDates, d)
	}
	return count, streakDates
}

func GetLongestStreak(data map[string]int, dates []string) (int, []string) {
	currentStreak := make([]string, 0)
	currentCount := 0
	longestStreak := make([]string, 0)
	longestCount := 0
	for i := len(dates) - 1; i >= 0; i-- {
		d := dates[i]
		if data[d] == 0 {
			if currentCount > longestCount {
				longestCount = currentCount
				longestStreak = currentStreak
			}
			currentCount = 0
			currentStreak = make([]string, 0)
			continue
		}
		currentCount = currentCount + 1
		currentStreak = append(currentStreak, d)
	}
	if currentCount > longestCount {
		longestCount = currentCount
		longestStreak = currentStreak
	}
	return longestCount, longestStreak
}

type Response struct {
	Username           string   `json:"username"`
	CurrentStreakCount int      `json:"currentStreakCount"`
	CurrentStreakDates []string `json:"currentStreakDates"`
	LongestStreakCount int      `json:"longestStreakCount"`
	LongestStreakDates []string `json:"longestStreakDates"`
}

func UsernameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	if !IsUsernameValid(username) {
		http.Error(w, "An unexpected error occurred when fetching your profile", http.StatusBadRequest)
		return
	}
	res, dates, err := GetDateData(username)
	if err != nil {
		http.Error(w, "An unexpected error occurred when fetching your profile", http.StatusBadRequest)
		return
	}
	currentCount, currentStreak := GetCurrentStreak(res, dates)
	longestCount, longestStreak := GetLongestStreak(res, dates)
	resp := &Response{
		Username:           username,
		CurrentStreakCount: currentCount,
		CurrentStreakDates: currentStreak,
		LongestStreakCount: longestCount,
		LongestStreakDates: longestStreak,
	}
	resJSON, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(resJSON))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

func main() {
	fmt.Printf("Start Server")
	r := mux.NewRouter()
	r.HandleFunc("/github-user/{username:[a-zA-Z0-9]+}", UsernameHandler)
	r.HandleFunc("/{username:[a-zA-Z0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "public/user.html")
	})
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))
	http.Handle("/", r)
	http.ListenAndServe(":9898", nil)
}
