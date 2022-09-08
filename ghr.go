package ghr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Time time.Time

func (impl Time) String() string {
	return time.Time(impl).String()
}

func (impl Time) GoTime() time.Time {
	return time.Time(impl)
}

func (impl *Time) IsToday() bool {
	return time.Since(impl.GoTime()) < 24*time.Hour
}

func (impl *Time) UnmarshalJSON(b []byte) error {
	t := time.Time{}
	err := json.Unmarshal(b, &t)
	if err != nil {
		return err
	}
	*impl = Time(t)
	return nil
}

type Author struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Uploader struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

type Asset struct {
	URL                string   `json:"url"`
	ID                 int      `json:"id"`
	NodeID             string   `json:"node_id"`
	Name               string   `json:"name"`
	Label              string   `json:"label"`
	Uploader           Uploader `json:"uploader"`
	ContentType        string   `json:"content_type"`
	State              string   `json:"state"`
	Size               int      `json:"size"`
	DownloadCount      int      `json:"download_count"`
	CreatedAt          Time     `json:"created_at"`
	UpdatedAt          Time     `json:"updated_at"`
	BrowserDownloadURL string   `json:"browser_download_url"`
}

type Reaction struct {
	URL        string `json:"url"`
	TotalCount int    `json:"total_count"`
	Num1       int    `json:"+1"`
	Num10      int    `json:"-1"`
	Laugh      int    `json:"laugh"`
	Hooray     int    `json:"hooray"`
	Confused   int    `json:"confused"`
	Heart      int    `json:"heart"`
	Rocket     int    `json:"rocket"`
	Eyes       int    `json:"eyes"`
}

type Release struct {
	URL             string   `json:"url"`
	AssetsURL       string   `json:"assets_url"`
	UploadURL       string   `json:"upload_url"`
	HTMLURL         string   `json:"html_url"`
	ID              int      `json:"id"`
	Author          Author   `json:"author"`
	NodeID          string   `json:"node_id"`
	TagName         string   `json:"tag_name"`
	TargetCommitish string   `json:"target_commitish"`
	Name            string   `json:"name"`
	Draft           bool     `json:"draft"`
	Prerelease      bool     `json:"prerelease"`
	CreatedAt       Time     `json:"created_at"`
	PublishedAt     Time     `json:"published_at"`
	Assets          []Asset  `json:"assets"`
	TarballURL      string   `json:"tarball_url"`
	ZipballURL      string   `json:"zipball_url"`
	Body            string   `json:"body"`
	MentionsCount   int      `json:"mentions_count,omitempty"`
	Reactions       Reaction `json:"reactions,omitempty"`
}

type Releases []Release

type GHR struct {
	Repo   string
	client *http.Client
}

func NewGHR(repo string) (*GHR, error) {
	if !strings.Contains(repo, "/") {
		return nil, errors.New("bad format")
	}
	return &GHR{
		Repo:   repo,
		client: http.DefaultClient,
	}, nil
}

// https://api.github.com/repos/go-swagger/go-swagger/releases

func (impl *GHR) GetReleases() (Releases, error) {
	const endPoint = "https://api.github.com/repos/%s/releases"

	resp, err := impl.client.Get(fmt.Sprintf(endPoint, impl.Repo))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("bad status: %d, %s", resp.StatusCode, resp.Status)
	}

	ret := Releases{}
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("json parse err: %w", err)
	}

	return ret, nil
}

type Tag struct {
	Name       string `json:"name"`
	ZipballURL string `json:"zipball_url"`
	TarballURL string `json:"tarball_url"`
	Commit     Commit `json:"commit"`
	NodeID     string `json:"node_id"`
}

type Commit struct {
	Sha string `json:"sha"`
	URL string `json:"url"`
}

type Tags []Tag

func (impl *GHR) GetTags() (Tags, error) {
	const endPoint = "https://api.github.com/repos/%s/tags"

	resp, err := impl.client.Get(fmt.Sprintf(endPoint, impl.Repo))
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	default:
		return nil, fmt.Errorf("bad status: %d, %s", resp.StatusCode, resp.Status)
	}

	ret := Tags{}
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, fmt.Errorf("json parse err: %w", err)
	}

	return ret, nil
}
