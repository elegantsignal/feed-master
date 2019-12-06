package proc

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/umputun/feed-master/app/feed"
)

func TestNewTelegramClientIfTokenEmpty(t *testing.T) {
	client, err := NewTelegramClient("", 0)

	assert.Nil(t, err)
	assert.Nil(t, client.Bot)
}

func TestNewTelegramClientCheckTimeout(t *testing.T) {
	cases := []struct {
		timeout, expected time.Duration
	}{
		{0, 600},
		{300, 300},
		{100500, 100500},
	}

	//nolint:scopelint
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			client, err := NewTelegramClient("", tc.timeout)

			assert.Nil(t, err)
			assert.Equal(t, tc.expected, client.Timeout)
		})
	}
}

func TestSendIfBotIsNil(t *testing.T) {
	client, err := NewTelegramClient("", 0)

	got := client.Send("@channel", feed.Item{})

	assert.Nil(t, err)
	assert.Nil(t, got)
}

func TestSendIfChannelIDEmpty(t *testing.T) {
	client := TelegramClient{
		Bot: &tb.Bot{},
	}

	got := client.Send("", feed.Item{})

	assert.Nil(t, got)
}

func TestTagLinkOnlySupport(t *testing.T) {
	html := `
<li>Особое канадское искусство. </li>
<li>Результаты нашего странного эксперимента.</li>
<li>Теперь можно и в <a href="https://t.me/uwp_podcast">телеграмме</a></li>
<li>Саботаж на местах.</li>
<li>Их нравы: кумовство и коррупция.</li>
<li>Вопросы и ответы</li>
</ul>
<p><a href="https://podcast.umputun.com/media/ump_podcast437.mp3"><em>аудио</em></a></p>`

	htmlExpected := `
Особое канадское искусство. 
Результаты нашего странного эксперимента.
Теперь можно и в <a href="https://t.me/uwp_podcast">телеграмме</a>
Саботаж на местах.
Их нравы: кумовство и коррупция.
Вопросы и ответы

<a href="https://podcast.umputun.com/media/ump_podcast437.mp3">аудио</a>`

	client := TelegramClient{}

	got := client.tagLinkOnlySupport(html)

	assert.Equal(t, got, htmlExpected, "support only html tag a")
}

func TestGetMessageHTML(t *testing.T) {
	item := feed.Item{
		Title:       "\tPodcast\n\t",
		Description: "<p>News <a href='#'>Podcast Link</a></p>\n",
		Enclosure: feed.Enclosure{
			URL: "https://example.com",
		},
	}

	expected := "Podcast\n\nNews <a href=\"#\">Podcast Link</a>\n\nhttps://example.com"

	client := TelegramClient{}
	got := client.getMessageHTML(item)

	assert.Equal(t, got, expected)
}

func TestRecipientChannelIDNotStartWithAt(t *testing.T) {
	cases := []string{"channel", "@channel"}
	expected := "@channel"

	for _, channelID := range cases {
		t.Run("", func(t *testing.T) {
			got := recipient{chatID: channelID} //nolint

			assert.Equal(t, got.Recipient(), expected)
		})
	}
}

func TestGetFilenameByURL(t *testing.T) {
	cases := []struct {
		url, expected string
	}{
		{"https://example.com/100500/song.mp3", "song.mp3"},
		{"https://example.com//song.mp3", "song.mp3"},
		{"https://example.com/song.mp3", "song.mp3"},
		{"https://example.com/song.mp3/", ""},
		{"https://example.com/", ""},
	}

	//nolint:scopelint
	for _, tc := range cases {
		t.Run("", func(t *testing.T) {
			client := TelegramClient{}
			got := client.getFilenameByURL(tc.url)

			assert.Equal(t, got, tc.expected)
		})
	}
}
