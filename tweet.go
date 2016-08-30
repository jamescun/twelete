package main

import (
	"encoding/csv"
	"strconv"
	"time"
)

type Tweet struct {
	Id uint64

	Text string

	ReplyId     string
	ReplyUserId string

	RetweetId     string
	RetweetUserId string

	Timestamp time.Time
}

const TwitterTimestamp string = "2006-01-02 15:04:05 -0700"

func readTweet(r *csv.Reader) (*Tweet, error) {
	record, err := r.Read()
	if err != nil {
		return nil, err
	}

	if len(record) < 10 || record[0] == "tweet_id" {
		return nil, ErrUnknownFormat
	}

	id, err := strconv.ParseUint(record[0], 10, 64)
	if err != nil {
		return nil, err
	}

	ts, err := time.Parse(TwitterTimestamp, record[3])
	if err != nil {
		return nil, err
	}

	return &Tweet{
		Id:            id,
		Text:          record[5],
		ReplyId:       record[1],
		ReplyUserId:   record[2],
		RetweetId:     record[6],
		RetweetUserId: record[7],
		Timestamp:     ts,
	}, nil
}
