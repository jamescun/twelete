package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/mrjones/oauth"
	"github.com/urfave/cli"
)

var (
	ErrNoArchive = errors.New("please specify the path to your twitter archive zip file --archive=<path>")

	ErrNoBefore = errors.New("please specify either --before=<YYYY_MM_DD> or --before-id=<tweet id>")

	ErrNotFound = errors.New("tweet not found")

	ErrRateLimit = errors.New("rate limited. try again later with --before-id=<last tweet id>")
)

var App = &cli.App{
	Name:     "twelete",
	HelpName: "twelete",
	Usage:    "delete tweets from a twitter archive",
	Version:  "1.0.0",
	Action:   Twelete,
	Writer:   os.Stdout,

	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "archive",
			Usage: "path to twitter archive zip file",
		},
		cli.StringFlag{
			Name:  "before",
			Usage: "delete tweets before this date (YYYY-MM-DD)",
		},
		cli.Int64Flag{
			Name:  "before-id",
			Usage: "delete tweets before an id",
		},
		cli.BoolFlag{
			Name:  "keep-replies",
			Usage: "don't delete replies",
		},
		cli.BoolFlag{
			Name:  "keep-retweets",
			Usage: "don't delete retweets",
		},
		cli.IntFlag{
			Name:  "limit",
			Usage: "maximum number of tweets to delete",
			Value: 1000,
		},
		cli.DurationFlag{
			Name:  "pause",
			Usage: "pause between deletes to prevent throttling by Twitter's API",
			Value: 10 * time.Second,
		},
		cli.StringFlag{
			Name:  "consumer-key",
			Usage: "twitter app consumer key",
		},
		cli.StringFlag{
			Name:  "consumer-secret",
			Usage: "twitter app consumer secret",
		},
		cli.StringFlag{
			Name:  "access-token",
			Usage: "twitter user access token",
		},
		cli.StringFlag{
			Name:  "access-secret",
			Usage: "twitter user access secret",
		},
	},
}

func Twelete(c *cli.Context) error {
	if c.String("archive") == "" {
		return ErrNoArchive
	}

	d := &Deleter{
		BeforeId: uint64(c.Uint64("before-id")),
		Retweets: !c.Bool("keep-retweets"),
		Replies:  !c.Bool("keep-replies"),
	}

	if b := c.String("before"); b != "" {
		t, err := time.Parse("2006-01-02", b)
		if err != nil {
			return err
		}

		d.Before = t
	}

	if d.Before.IsZero() && d.BeforeId < 1 {
		return ErrNoBefore
	}

	twitter := oauth.NewConsumer(c.String("consumer-key"), c.String("consumer-secret"), oauth.ServiceProvider{})
	twitterClient, _ := twitter.MakeHttpClient(&oauth.AccessToken{
		Token:  c.String("access-token"),
		Secret: c.String("access-secret"),
	})

	a, err := NewArchive(c.String("archive"))
	if err != nil {
		return err
	}

	total := 0
	limit := c.Int("limit")
	pause := c.Duration("pause")
	for {
		if total >= limit {
			break
		}

		tweet, err := a.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if d.Delete(tweet) {
			log.Printf("tweet: id=%d text=%s\n", tweet.Id, tweet.Timestamp.Format(time.RFC3339))
			err := deleteTweet(twitterClient, tweet.Id)
			if err == ErrNotFound {
				log.Println("tweet: already deleted")
				continue
			} else if err != nil {
				return err
			}
			total++
			time.Sleep(pause)
		}
	}

	log.Println("deleted", total, "tweets")

	return nil
}

type HTTPError int

func (e HTTPError) Error() string {
	return fmt.Sprintf("HTTP Error %d", e)
}

func deleteTweet(c *http.Client, id uint64) error {
	idStr := strconv.FormatUint(id, 10)
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/statuses/destroy/"+idStr+".json", nil)
	if err != nil {
		return err
	}

	res, err := c.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		if res.StatusCode == 404 {
			return ErrNotFound
		} else if res.StatusCode == 429 {
			return ErrRateLimit
		}

		return HTTPError(res.StatusCode)
	}

	return nil
}

func main() {
	err := App.Run(os.Args)
	if err != nil {
		log.Fatalln("error:", err)
	}
}
