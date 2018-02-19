package controllers

import (

	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/jaimejorge/go-cast/api"
	"github.com/jaimejorge/go-cast/events"
	"github.com/jaimejorge/go-cast/log"
	"github.com/jaimejorge/go-cast/net"
)

type UrlController struct {
	interval       time.Duration
	channel        *net.Channel
	eventsCh       chan events.Event
	DestinationID  string
	UrlSessionID int
}

const NamespaceUrl = "urn:x-cast:com.url.cast"

type UrlCommand struct {
	net.PayloadHeaders
	UrlSessionID int `json:"mediaSessionId"`
}

type LoadUrlCommand struct {
	net.PayloadHeaders
	Url       UrlItem   `json:"media"`
	CurrentTime int         `json:"currentTime"`
	Autoplay    bool        `json:"autoplay"`
	CustomData  interface{} `json:"customData"`
}

type UrlItem struct {
	ContentId   string `json:"contentId"`
	StreamType  string `json:"streamType"`
	ContentType string `json:"contentType"`
}

type UrlStatusUrl struct {
	ContentId   string  `json:"contentId"`
	StreamType  string  `json:"streamType"`
	ContentType string  `json:"contentType"`
	Duration    float64 `json:"duration"`
}

type DashUrlCommand struct {
        net.PayloadHeaders
        Url       string   `json:"url"`
        Type  string   `json:"type"`
}


func NewUrlController(conn *net.Connection, eventsCh chan events.Event, sourceId, destinationID string) *UrlController {
	controller := &UrlController{
		channel:       conn.NewChannel(sourceId, destinationID, NamespaceUrl),
		eventsCh:      eventsCh,
		DestinationID: destinationID,
	}

	controller.channel.OnMessage("MEDIA_STATUS", controller.onStatus)

	return controller
}

func (c *UrlController) SetDestinationID(id string) {
	c.channel.DestinationId = id
	c.DestinationID = id
}

func (c *UrlController) sendEvent(event events.Event) {
	select {
	case c.eventsCh <- event:
	default:
		log.Printf("Dropped event: %#v", event)
	}
}

func (c *UrlController) onStatus(message *api.CastMessage) {
	response, err := c.parseStatus(message)
	if err != nil {
		log.Errorf("Error parsing status: %s", err)
	}

	for _, status := range response.Status {
		c.sendEvent(*status)
	}
}

func (c *UrlController) parseStatus(message *api.CastMessage) (*UrlStatusResponse, error) {
	response := &UrlStatusResponse{}

	err := json.Unmarshal([]byte(*message.PayloadUtf8), response)

	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal status message:%s - %s", err, *message.PayloadUtf8)
	}

	for _, status := range response.Status {
		c.UrlSessionID = status.UrlSessionID
	}

	return response, nil
}

type UrlStatusResponse struct {
	net.PayloadHeaders
	Status []*UrlStatus `json:"status,omitempty"`
}

type UrlStatus struct {
	net.PayloadHeaders
	UrlSessionID         int                    `json:"mediaSessionId"`
	PlaybackRate           float64                `json:"playbackRate"`
	PlayerState            string                 `json:"playerState"`
	CurrentTime            float64                `json:"currentTime"`
	SupportedUrlCommands int                    `json:"supportedUrlCommands"`
	Volume                 *Volume                `json:"volume,omitempty"`
	Url                  *UrlStatusUrl      `json:"media"`
	CustomData             map[string]interface{} `json:"customData"`
	RepeatMode             string                 `json:"repeatMode"`
	IdleReason             string                 `json:"idleReason"`
}

func (c *UrlController) Start(ctx context.Context) error {
	//_, err := c.GetStatus(ctx)
	return nil
}
func (c *UrlController) LoadUrl(ctx context.Context, url string, Force bool,Reload bool,Reload_time int ) (*api.CastMessage, error) {
        c.channel.Request(ctx, &DashUrlCommand{
                Url: url,
                Type: "loc",
        })


        return nil, nil
}
