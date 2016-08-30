package main

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"io"
)

var (
	ErrUnknownFormat = errors.New("archive format is unknown")
)

const TweetCSVFilename string = "tweets.csv"

type Archive struct {
	z *zip.ReadCloser
	r io.ReadCloser
	c *csv.Reader
}

func NewArchive(filename string) (*Archive, error) {
	z, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}

	for _, file := range z.File {
		if file.Name == TweetCSVFilename {
			r, err := file.Open()
			if err != nil {
				z.Close()
				return nil, err
			}

			return &Archive{
				z: z,
				r: r,
				c: csv.NewReader(r),
			}, nil
		}
	}

	z.Close()
	return nil, ErrUnknownFormat
}

func (a *Archive) Close() error {
	// close zip file
	if a.r != nil {
		err := a.r.Close()
		if err != nil {
			return err
		}
		a.r = nil
	}

	// close zip archive
	if a.z != nil {
		err := a.z.Close()
		if err != nil {
			return err
		}
		a.z = nil
	}

	return nil
}

func (a *Archive) Next() (*Tweet, error) {
	for {
		tweet, err := readTweet(a.c)
		if err == ErrUnknownFormat {
			continue
		} else if err != nil {
			return nil, err
		}

		return tweet, nil
	}

	return nil, nil
}
