package chatgptgo

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/acheong08/ChatGPT-Go/config"
	http "github.com/bogdanfinn/fhttp"
	tls_client "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func NewChatbot(accessToken string) (Chatbot, error) {
	client, err := tls_client.NewHttpClient(
		tls_client.NewNoopLogger(),
		tls_client.WithClientProfile(profiles.Firefox_117),
		tls_client.WithRandomTLSExtensionOrder(),
	)
	if config.PUID != "" {
		client.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "chat.openai.com",
		}, []*http.Cookie{
			{
				Name:  "_puid",
				Value: config.PUID,
			},
		})
	}
	return Chatbot{
		AccessToken: accessToken,
		HTTPClient:  &client,
	}, err
}

func (c *Chatbot) GetHistory(offset, limit int) (*Conversations, error) {
	if limit == 0 {
		limit = 28
	}
	var conversations Conversations
	err := c.getRequest(
		fmt.Sprintf("https://chat.openai.com/backend-api/conversations?offset=%d&limit=%d&order=updated", offset, limit),
		&conversations,
	)
	return &conversations, err
}

func (c *Chatbot) GetConversation(conversationID string) (*Conversation, error) {
	var conversation Conversation
	err := c.getRequest(
		fmt.Sprintf("https://chat.openai.com/backend-api/conversation/%s", conversationID),
		&conversation,
	)
	return &conversation, err
}

func (c *Chatbot) getRequest(url string, obj any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.AccessToken))
	addHeaders(req)
	resp, err := (*c.HTTPClient).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("Received status code %d", resp.StatusCode)
	}
	err = json.NewDecoder(resp.Body).Decode(obj)
	return err

}

func addHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/118.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US")
	// req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://chat.openai.com/")
	req.Header.Set("DNT", "1")
	req.Header.Set("Alt-Used", "chat.openai.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Sec-Fetch-Dest", "empty")
	req.Header.Set("Sec-Fetch-Mode", "cors")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("TE", "trailers")
}
