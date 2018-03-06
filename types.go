package travisci

// https://developer.travis-ci.com/explore/#explorer
// https://developer.travis-ci.com/format

type APIObject struct {
	Type           string          `json:"@type"`
	HREF           string          `json:"@href"`
	Representation string          `json:"@representation"`
	Pagination     *Pagination     `json:"@pagination"`
	Permissions    map[string]bool `json:"@permissions"`
}

type Page struct {
	HREF   string `json:"@href"`
	Offset int
	Limit  int
}

type Pagination struct {
	Count  int // overall number of entries in the collection
	Limit  int // maximum number of entries included in the current subset
	Offset int // number of entries preceding the first entry of the subset

	IsFirst bool `json:"is_first"`
	IsLast  bool `json:"is_last"`

	First *Page
	Next  *Page
	Last  *Page
}

type Repository struct {
	APIObject

	ID int

	Name string
	Slug string

	Description string

	// TODO: Always seems to be `null`
	// GithubLanguage string `json:"github_language"`

	Active  bool
	Starred bool

	Private bool

	DefaultBranch Branch `json:"default_branch"`
}

type Repositories struct {
	APIObject

	Repositories []Repository
}

type Build struct {
	APIObject

	ID int

	Number        string
	State         string
	Duration      int    // seconds
	EventType     string `json:"event_type"`
	PreviousState string `json:"previous_state"`

	PullRequestTitle  string `json:"pull_request_title"`
	PullRequestNumber int    `json:"pull_request_number"`

	StartedAt  string `json:"started_at"`
	FinishedAt string `json:"finished_at"`
}

type Branch struct {
	APIObject

	Name string
}

type EnvVar struct {
	APIObject

	ID    string
	Name  string
	Value string

	Public bool
}

type EnvVars struct {
	APIObject

	EnvVars []EnvVar `json:"env_vars"`
}

type User struct {
	APIObject

	ID        int
	Login     string
	Name      string
	GitHubID  int    `json:"github_id"`
	AvatarURL string `json:"avatar_url"`
	IsSyncing bool   `json:"is_syncing"`
	SyncedAt  string `json:"synced_at"`
}
