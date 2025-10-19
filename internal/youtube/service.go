package youtube

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

// Service wraps the YouTube API service
type Service struct {
	client *youtube.Service
}

// Video represents a YouTube video
type Video struct {
	ID          string
	Title       string
	Description string
	PublishedAt string
}

// NewService creates a new YouTube service
func NewService(httpClient *http.Client) (*Service, error) {
	ctx := context.Background()
	service, err := youtube.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("error creating YouTube client: %w", err)
	}

	return &Service{client: service}, nil
}

// ListChannelVideos lists all videos from a channel
func (s *Service) ListChannelVideos(channelID string, maxResults int64) ([]Video, error) {
	// If no channel ID is provided, get the authenticated user's channel
	if channelID == "" {
		channelsCall := s.client.Channels.List([]string{"id"}).Mine(true)
		channelResponse, err := channelsCall.Do()
		if err != nil {
			return nil, fmt.Errorf("error retrieving channel: %w", err)
		}

		if len(channelResponse.Items) == 0 {
			return nil, fmt.Errorf("no channel found for authenticated user")
		}

		channelID = channelResponse.Items[0].Id
		fmt.Printf("Using channel ID: %s\n", channelID)
	}

	// Get uploads playlist ID
	channelsCall := s.client.Channels.List([]string{"contentDetails"}).Id(channelID)
	channelResponse, err := channelsCall.Do()
	if err != nil {
		return nil, fmt.Errorf("error retrieving channel details: %w", err)
	}

	if len(channelResponse.Items) == 0 {
		return nil, fmt.Errorf("channel not found")
	}

	uploadsPlaylistID := channelResponse.Items[0].ContentDetails.RelatedPlaylists.Uploads

	// Get videos from uploads playlist
	var videos []Video
	nextPageToken := ""

	for {
		call := s.client.PlaylistItems.List([]string{"snippet"}).
			PlaylistId(uploadsPlaylistID).
			MaxResults(50)

		if nextPageToken != "" {
			call = call.PageToken(nextPageToken)
		}

		response, err := call.Do()
		if err != nil {
			return nil, fmt.Errorf("error retrieving playlist items: %w", err)
		}

		for _, item := range response.Items {
			videos = append(videos, Video{
				ID:          item.Snippet.ResourceId.VideoId,
				Title:       item.Snippet.Title,
				Description: item.Snippet.Description,
				PublishedAt: item.Snippet.PublishedAt,
			})

			if maxResults > 0 && int64(len(videos)) >= maxResults {
				return videos, nil
			}
		}

		nextPageToken = response.NextPageToken
		if nextPageToken == "" {
			break
		}
	}

	return videos, nil
}

// GetChannelInfo returns information about a channel
func (s *Service) GetChannelInfo(channelID string) (*youtube.Channel, error) {
	call := s.client.Channels.List([]string{"snippet", "contentDetails", "statistics"})
	
	if channelID == "" {
		call = call.Mine(true)
	} else {
		call = call.Id(channelID)
	}

	response, err := call.Do()
	if err != nil {
		return nil, fmt.Errorf("error retrieving channel info: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, fmt.Errorf("channel not found")
	}

	return response.Items[0], nil
}

