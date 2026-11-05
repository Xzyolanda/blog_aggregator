package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestParsePublishedAt(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    time.Time
		valid   bool
		wantErr bool
	}{
		{
			name:  "RFC1123 with numeric timezone",
			value: "Mon, 02 Jan 2006 15:04:05 -0700",
			want:  time.Date(2006, time.January, 2, 22, 4, 5, 0, time.UTC),
			valid: true,
		},
		{
			name:  "RFC822 with numeric timezone",
			value: "02 Jan 06 15:04 +0200",
			want:  time.Date(2006, time.January, 2, 13, 4, 0, 0, time.UTC),
			valid: true,
		},
		{
			name:  "RFC3339",
			value: "2006-01-02T15:04:05+02:00",
			want:  time.Date(2006, time.January, 2, 13, 4, 5, 0, time.UTC),
			valid: true,
		},
		{
			name:  "date without timezone",
			value: "2006-01-02",
			want:  time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC),
			valid: true,
		},
		{
			name:  "empty date",
			value: "",
			valid: false,
		},
		{
			name:    "unsupported date",
			value:   "not a date",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := parsePublishedAt(test.value)
			if (err != nil) != test.wantErr {
				t.Fatalf("parsePublishedAt() error = %v, wantErr %v", err, test.wantErr)
			}
			if test.wantErr {
				return
			}
			if got.Valid != test.valid {
				t.Fatalf("parsePublishedAt() Valid = %v, want %v", got.Valid, test.valid)
			}
			if test.valid && !got.Time.Equal(test.want) {
				t.Fatalf("parsePublishedAt() = %v, want %v", got.Time, test.want)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFetchFeedUnescapesFieldsAndSetsUserAgent(t *testing.T) {
	previousClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
		if got := req.Header.Get("User-Agent"); got != "gator" {
			t.Errorf("User-Agent = %q, want %q", got, "gator")
		}
		body := `<rss><channel><title>Tech &amp;ldquo;News&amp;rdquo;</title><description>Blog &amp;amp; news</description><item><title>Post &amp;ldquo;one&amp;rdquo;</title><description>First &amp;amp; foremost</description></item></channel></rss>`
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(body)),
			Request:    req,
		}, nil
	})}
	defer func() { http.DefaultClient = previousClient }()

	feed, err := fetchFeed(context.Background(), "https://feed.test/rss")
	if err != nil {
		t.Fatalf("fetchFeed() error = %v", err)
	}
	if got, want := feed.Channel.Title, "Tech “News”"; got != want {
		t.Errorf("channel title = %q, want %q", got, want)
	}
	if got, want := feed.Channel.Description, "Blog & news"; got != want {
		t.Errorf("channel description = %q, want %q", got, want)
	}
	if len(feed.Channel.Item) != 1 {
		t.Fatalf("got %d items, want 1", len(feed.Channel.Item))
	}
	if got, want := feed.Channel.Item[0].Title, "Post “one”"; got != want {
		t.Errorf("item title = %q, want %q", got, want)
	}
	if got, want := feed.Channel.Item[0].Description, "First & foremost"; got != want {
		t.Errorf("item description = %q, want %q", got, want)
	}
}
