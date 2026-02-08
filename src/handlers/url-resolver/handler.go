package main

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"urlresolver/src/dataobjects/domainobjects"
	"urlresolver/src/services/availability"
	"urlresolver/src/services/identity"
)

var videos = map[string]domainobjects.Video{
	"001": {
		ID:                    "46325",
		Title:                 "Example Video 001",
		StandardDefinitionURL: "example001",
		PremiumDefinitionURL:  "example001-premium",
	},
	"002": {
		ID:                    "46326",
		Title:                 "Example Video 002",
		StandardDefinitionURL: "example002",
		PremiumDefinitionURL:  "example002-premium",
	},
}

type PlaybackResponse struct {
	VideoID           string `json:"video_id"`
	Title             string `json:"title"`
	PlaybackBaseURL   string `json:"playback_baseurl"`
	PlaybackFilename  string `json:"playback_filename"`
	PlaybackExtension string `json:"playback_extension"`
}

type IdentityService interface {
	GetUserInfo(ctx context.Context, token string) (identity.UserInfo, error)
}

type AvailabilityService interface {
	GetAvailability(ctx context.Context, videoID, token string) (availability.Availability, error)
}

type Handler struct {
	identitySvc     IdentityService
	availabilitySvc AvailabilityService
	playbackBaseURL string
}

func main() {
	http.HandleFunc("/video/", NewHandler().GetVideo)
	http.ListenAndServe(":8080", nil)
}

func NewHandler() *Handler {
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}

	baseURL := "https://bytestream.notreal"
	playbackBaseURL := "https://s3.eu-west-1.amazonaws.com/bytestreamfake"

	identitySvc := identity.NewService(baseURL, httpClient)
	availabilitySvc := availability.NewService(baseURL, httpClient)

	return &Handler{
		identitySvc:     identitySvc,
		availabilitySvc: availabilitySvc,
		playbackBaseURL: playbackBaseURL,
	}
}

func (h *Handler) GetVideo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Extract video ID from path: /video/{id}
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 2 || parts[0] != "video" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	videoID := parts[1]

	// Extract bearer token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	token := strings.TrimSpace(authHeader[len("bearer "):])

	video, ok := videos[videoID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	user, err := h.identitySvc.GetUserInfo(ctx, token)
	if err != nil {
		http.Error(w, "identity service error", http.StatusBadGateway)
		return
	}

	availability, err := h.availabilitySvc.GetAvailability(ctx, videoID, token)
	if err != nil {
		http.Error(w, "availability service error", http.StatusBadGateway)
		return
	}

	// Check availability window
	if !availability.IsValid(time.Now().UTC()) {
		http.Error(w, "video not available", http.StatusForbidden)
		return
	}

	// Select filename
	filename := video.StandardDefinitionURL
	if user.IsPremium() {
		filename = video.PremiumDefinitionURL
	}

	response := PlaybackResponse{
		VideoID:           video.ID,
		Title:             video.Title,
		PlaybackBaseURL:   h.playbackBaseURL,
		PlaybackFilename:  filename,
		PlaybackExtension: ".mp4",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
