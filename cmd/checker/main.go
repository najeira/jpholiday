// based on https://github.com/shogo82148/holidays-jp
// MIT License Copyright (c) 2021 Ichinose Shogo
// Copyright (c) 2022 Tetsuhiro Ueda
// downloader for syukujitsu.csv
// https://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html
package main

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/najeira/jpholiday"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// 内閣府ホーム  >  内閣府の政策  >  制度  >  国民の祝日について
// https://www8.cao.go.jp/chosei/shukujitsu/gaiyou.html
const syukujitsuURL = "https://www8.cao.go.jp/chosei/shukujitsu/syukujitsu.csv"

func main() {
	if err := _main(); err != nil {
		log.Fatal(err)
	}
}

func _main() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	rawData, err := download(ctx)
	if err != nil {
		return err
	}

	holidays, err := formatHolidays(rawData)
	if err != nil {
		return err
	}

	const startYear = 1955
	endYear := time.Now().Year() + 1

	var fails int
	date := time.Date(startYear, 1, 1, 0, 0, 0, 0, jpholiday.JST)
	for {
		y, m, d := date.Date()
		if y > endYear {
			break
		}

		dateStr := fmt.Sprintf("%04d-%02d-%02d", y, m, d)
		csvName := holidays[dateStr]
		libName := jpholiday.Name(date)

		if csvName != "" {
			if libName == "" {
				fails++
				fmt.Printf("missing %s %s\n", dateStr, csvName)
			}
		} else {
			if libName != "" {
				fails++
				fmt.Printf("extra %s %s\n", dateStr, libName)
			}
		}

		date = date.AddDate(0, 0, 1)
	}

	fmt.Printf("%d fails: %d-01-01 to %d-12-31\n", fails, startYear, endYear)
	return nil
}

func download(ctx context.Context) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, syukujitsuURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "https://github.com/shogo82148/holidays-jp")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// continue to download
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func formatHolidays(rawData []byte) (map[string]string, error) {
	dec := japanese.ShiftJIS.NewDecoder()
	br := bytes.NewReader(rawData)
	tr := transform.NewReader(br, dec)
	cr := csv.NewReader(tr)

	// skip 国民の祝日・休日月日,国民の祝日・休日名称 line
	if _, err := cr.Read(); err != nil {
		return nil, err
	}

	holidays := make(map[string]string)
	for {
		record, err := cr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		date := formatDate(record[0])
		holidays[date] = record[1]
	}
	return holidays, nil
}

// 2021/1/1 -> 2021-01-01
func formatDate(s string) string {
	date := strings.Split(s, "/")
	y, err := strconv.Atoi(date[0])
	if err != nil {
		panic(err)
	}
	m, err := strconv.Atoi(date[1])
	if err != nil {
		panic(err)
	}
	d, err := strconv.Atoi(date[2])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%04d-%02d-%02d", y, m, d)
}
