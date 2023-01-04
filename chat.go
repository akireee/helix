package helix

import "errors"

type GetChatChattersParams struct {
	BroadcasterID string `query:"broadcaster_id"`
	ModeratorID   string `query:"moderator_id"`
	After         string `query:"after"`
	First         string `query:"first"`
}

type ChatChatter struct {
	UserLogin string `json:"user_login"`
	UserID    string `json:"user_id"`
	Username  string `json:"user_name"`
}

type ManyChatChatters struct {
	Chatters   []ChatChatter `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

type GetChatChattersResponse struct {
	ResponseCommon
	Data ManyChatChatters
}

// Required scope: moderator:read:chatters
func (c *Client) GetChannelChatChatters(params *GetChatChattersParams) (*GetChatChattersResponse, error) {
	if params.BroadcasterID == "" || params.ModeratorID == "" {
		return nil, errors.New("error: broadcaster and moderator identifiers must be provided")
	}
	resp, err := c.get("/chat/chatters", &ManyChatChatters{}, params)
	if err != nil {
		return nil, err
	}

	chatters := &GetChatChattersResponse{}
	resp.HydrateResponseCommon(&chatters.ResponseCommon)
	chatters.Data.Chatters = resp.Data.(*ManyChatChatters).Chatters

	return chatters, nil
}

type GetChatBadgeParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetChatBadgeResponse struct {
	ResponseCommon
	Data ManyChatBadge
}

type ManyChatBadge struct {
	Badges []ChatBadge `json:"data"`
}

type ChatBadge struct {
	SetID    string         `json:"set_id"`
	Versions []BadgeVersion `json:"versions"`
}

type BadgeVersion struct {
	ID         string `json:"id"`
	ImageUrl1x string `json:"image_url_1x"`
	ImageUrl2x string `json:"image_url_2x"`
	ImageUrl4x string `json:"image_url_4x"`
}

func (c *Client) GetChannelChatBadges(params *GetChatBadgeParams) (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges", &ManyChatBadge{}, params)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

func (c *Client) GetGlobalChatBadges() (*GetChatBadgeResponse, error) {
	resp, err := c.get("/chat/badges/global", &ManyChatBadge{}, nil)
	if err != nil {
		return nil, err
	}

	channels := &GetChatBadgeResponse{}
	resp.HydrateResponseCommon(&channels.ResponseCommon)
	channels.Data.Badges = resp.Data.(*ManyChatBadge).Badges

	return channels, nil
}

type GetChannelEmotesParams struct {
	BroadcasterID string `query:"broadcaster_id"`
}

type GetEmoteSetsParams struct {
	EmoteSetIDs []string `query:"emote_set_id"` // Minimum: 1. Maximum: 25.
}

type SendChatAnnouncementParams struct {
	BroadcasterID string `query:"broadcaster_id"` // required
	ModeratorID   string `query:"moderator_id"`   // required
	Message       string `json:"message"`         // upto 500 chars, thereafter str is truncated
	// blue || green || orange || purple are valid, default 'primary' or empty str result in channel accent color.
	Color string `json:"color"`
}

type SendChatAnnouncementResponse struct {
	ResponseCommon
}

type GetChannelEmotesResponse struct {
	ResponseCommon
	Data ManyEmotes
}

type GetEmoteSetsResponse struct {
	ResponseCommon
	Data ManyEmotesWithOwner
}

type ManyEmotes struct {
	Emotes []Emote `json:"data"`
}

type ManyEmotesWithOwner struct {
	Emotes []EmoteWithOwner `json:"data"`
}

type Emote struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Images     EmoteImage `json:"images"`
	Tier       string     `json:"tier"`
	EmoteType  string     `json:"emote_type"`
	EmoteSetId string     `json:"emote_set_id"`
}

type EmoteWithOwner struct {
	Emote
	OwnerID string `json:"owner_id"`
}

type EmoteImage struct {
	Url1x string `json:"url_1x"`
	Url2x string `json:"url_2x"`
	Url4x string `json:"url_4x"`
}

func (c *Client) GetChannelEmotes(params *GetChannelEmotesParams) (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes", &ManyEmotes{}, params)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

func (c *Client) GetGlobalEmotes() (*GetChannelEmotesResponse, error) {
	resp, err := c.get("/chat/emotes/global", &ManyEmotes{}, nil)
	if err != nil {
		return nil, err
	}

	emotes := &GetChannelEmotesResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotes).Emotes

	return emotes, nil
}

// GetEmoteSets
func (c *Client) GetEmoteSets(params *GetEmoteSetsParams) (*GetEmoteSetsResponse, error) {
	resp, err := c.get("/chat/emotes/set", &ManyEmotesWithOwner{}, params)
	if err != nil {
		return nil, err
	}

	emotes := &GetEmoteSetsResponse{}
	resp.HydrateResponseCommon(&emotes.ResponseCommon)
	emotes.Data.Emotes = resp.Data.(*ManyEmotesWithOwner).Emotes

	return emotes, nil
}

// SendChatAnnouncement sends an announcement to the broadcaster’s chat room.
// Required scope: moderator:manage:announcements
func (c *Client) SendChatAnnouncement(params *SendChatAnnouncementParams) (*SendChatAnnouncementResponse, error) {
	resp, err := c.postAsJSON("/chat/announcements", nil, params)
	if err != nil {
		return nil, err
	}

	chatResp := &SendChatAnnouncementResponse{}
	resp.HydrateResponseCommon(&chatResp.ResponseCommon)

	return chatResp, nil
}

type GetChatSettingsParams struct {
	// Required, the ID of the broadcaster whose chat settings you want to get
	BroadcasterID string `query:"broadcaster_id"`

	// Optional, can be specified if you want the `non_moderator_chat_delay` and `non_moderator_chat_delay_duration` fields in the response. The ID should be a user that has moderation privileges in the broadcaster's chat.
	// The ID must match the specified User Access Token & the User Access Token must have the `moderator:read:chat_settings` scope
	ModeratorID string `query:"moderator_id,omitempty"`
}

type ChatSettings struct {
	BroadcasterID string `json:"broadcaster_id"`

	EmoteMode bool `json:"emote_mode"`

	FollowerMode bool `json:"follower_mode"`
	// Follower mode duration in minutes
	FollowerModeDuration int `json:"follower_mode_duration"`

	SlowMode bool `json:"slow_mode"`
	// Slow mode wait time in seconds
	SlowModeWaitTime int `json:"slow_mode_wait_time"`

	SubscriberMode bool `json:"subscriber_mode"`

	UniqueChatMode bool `json:"unique_chat_mode"`

	// Only included if the user access token includes the `moderator:read:chat_settings` scope
	ModeratorID string `json:"moderator_id"`

	// Boolean value denoting whether the "Non moderator chat delay" setting is enabled.
	// Only included if the request specifies a user access token that includes the moderator:read:chat_settings scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelay bool `json:"non_moderator_chat_delay"`
	// The amount of time, in seconds, that messages are delayed before appearing in chat.
	// Only included if the request specifies a user access token that includes the moderator:read:chat_settings scope and the user in the moderator_id query parameter is one of the broadcaster’s moderators.
	NonModeratorChatDelayDuration int `json:"non_moderator_chat_delay_duration"`
}

type ManyChatSettings struct {
	Settings []ChatSettings `json:"data"`
}

type GetChatSettingsResponse struct {
	ResponseCommon
	Data ManyChatSettings
}

// GetChatSettings gets the chat settings for the broadcaster's chat room.
// Optional scope: moderator:read:chat_settings
func (c *Client) GetChatSettings(params *GetChatSettingsParams) (*GetChatSettingsResponse, error) {
	if params.BroadcasterID == "" {
		return nil, errors.New("error: broadcaster id must be specified")
	}
	resp, err := c.get("/chat/settings", &ManyChatSettings{}, params)
	if err != nil {
		return nil, err
	}

	settings := &GetChatSettingsResponse{}
	resp.HydrateResponseCommon(&settings.ResponseCommon)
	settings.Data.Settings = resp.Data.(*ManyChatSettings).Settings

	return settings, nil
}
