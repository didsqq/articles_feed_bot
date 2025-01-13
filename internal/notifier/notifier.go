package notifier

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/didsqq/news_feed_bot/internal/botkit/markup"
	"github.com/didsqq/news_feed_bot/internal/model"
	"github.com/go-shiori/go-readability"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ArticleProvider interface {
	AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error)
	MarkAsPosted(ctx context.Context, article model.Article) error
}

type Summarizer interface {
	Summarize(text string) (string, error)
}

type UsersProvider interface {
	GetAll(ctx context.Context) ([]model.User, error)
}

type Notifier struct {
	articles         ArticleProvider
	users            UsersProvider
	summarizer       Summarizer
	bot              *tgbotapi.BotAPI
	sendInterval     time.Duration
	lookupTimeWindow time.Duration
	channelID        int64
}

func New(
	articleProvider ArticleProvider,
	users UsersProvider,
	summarizer Summarizer,
	bot *tgbotapi.BotAPI,
	sendInterval time.Duration,
	lookupTimeWindow time.Duration,
	channelID int64,
) *Notifier {
	return &Notifier{
		articles:         articleProvider,
		users:            users,
		summarizer:       summarizer,
		bot:              bot,
		sendInterval:     sendInterval,
		lookupTimeWindow: lookupTimeWindow,
		channelID:        channelID,
	}
}

func (n *Notifier) Start(ctx context.Context) error {
	log.Print("[INFO] notifier start")
	ticker := time.NewTicker(n.sendInterval)
	defer ticker.Stop()

	if err := n.SelectAndSendArticle(ctx); err != nil {
		return err
	}

	for {
		select {
		case <-ticker.C:
			if err := n.SelectAndSendArticle(ctx); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (n *Notifier) SelectAndSendArticle(ctx context.Context) error {
	topOneArticles, err := n.articles.AllNotPosted(ctx, time.Now().Add(-(n.lookupTimeWindow)), 2)
	if err != nil {
		return err
	}

	if len(topOneArticles) == 0 {
		return nil
	}

	article := topOneArticles[0]

	summary, err := n.extractSummary(article)
	if err != nil {
		log.Printf("[ERROR] failed to extract summary: %v", err)
	}

	users, err := n.users.GetAll(ctx)
	if err != nil {
		log.Printf("[ERROR] failed to get users: %v", err)
	}

	var errCh = make(chan error)
	var wg sync.WaitGroup

	for _, user := range users {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := n.sendArticle(user, article, summary)
			if err != nil {
				errCh <- err
				return
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		log.Printf("[ERROR] Ошибка при отправке article пользователю:%v", err)
	}

	return n.articles.MarkAsPosted(ctx, article)
}

func (n *Notifier) extractSummary(article model.Article) (string, error) {
	var r io.Reader

	if article.Summary != "" {
		r = strings.NewReader(article.Summary)
	} else {
		resp, err := http.Get(article.Link)
		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		r = resp.Body
	}

	doc, err := readability.FromReader(r, nil)
	if err != nil {
		return "", err
	}

	summary, err := n.summarizer.Summarize(cleanupText(doc.TextContent))
	if err != nil {
		return "", err
	}

	return "\n\n" + summary, nil
}

var redundantNewLines = regexp.MustCompile(`\n{3,}`)

func cleanupText(text string) string {
	return redundantNewLines.ReplaceAllString(text, "\n")
}

func (n *Notifier) sendArticle(user model.User, article model.Article, summary string) error {
	approach := false
	if user.Keywords[0] != "" {
		for _, keyword := range user.Keywords {
			if strings.Contains(strings.ToLower(summary), keyword) || strings.Contains(strings.ToLower(article.Title), keyword) {
				approach = true
				break
			}
		}
	} else {
		approach = true
	}

	if approach {
		const msgFormat = "*%s*%s\n\n%s"

		msg := tgbotapi.NewMessage(user.ChatID, fmt.Sprintf(
			msgFormat,
			markup.EscapeForMarkdown(article.Title),
			markup.EscapeForMarkdown(summary),
			markup.EscapeForMarkdown(article.Link),
		))
		msg.ParseMode = "MarkdownV2"

		_, err := n.bot.Send(msg)
		if err != nil {
			return err
		}
	} else {
		msg := tgbotapi.NewMessage(user.ChatID, "[Тестово]эта статья вам не подошла")
		msg.ParseMode = "MarkdownV2"

		_, err := n.bot.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
