package lib

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"shortly-api-service/internal/redis"
	"time"

	"github.com/mssola/user_agent"
)

func ParseUserAgent(uaString string) (device, browser, os string) {

	ua := user_agent.New(uaString)
	name, version := ua.Browser()

	if ua.Mobile() {
		device = "Mobile"
	} else if ua.Bot() {
		device = "Bot"
	} else {
		device = "Desktop"
	}

	browser = name + " " + version
	os = ua.OS()
	return
}

type ipApiResponse struct {
	CountryName string `json:"country_name"`
}

func GetCountryFromIP(ip string) string {

	cacheKey := "ip-country:" + ip

	country, err := redis.RedisClient.Get(context.Background(), cacheKey).Result()

	if err == nil && country != "" {
		return country
	}

	if ip == "127.0.0.1" || ip == "::1" {
		return "Localhost"
	}

	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)

	client := http.Client{Timeout: 2 * time.Second}

	resp, err := client.Get(url)

	if err != nil {
		return "Unknown"
	}
	defer resp.Body.Close()

	var result ipApiResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil || result.CountryName == "" {
		return "Unknown"
	}

	_ = redis.RedisClient.Set(context.Background(), cacheKey, result.CountryName, 24*time.Hour).Err()
	return result.CountryName
}
