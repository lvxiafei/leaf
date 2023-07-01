package tools

type JSONData struct {
	Type     string     `json:"type"`
	Version  int        `json:"version"`
	Source   string     `json:"source"`
	Elements []Elements `json:"elements"`
	AppState AppState   `json:"appState"`
	Files    Files      `json:"files"`
}
type StartBinding struct {
	ElementID string  `json:"elementId,omitempty"`
	Focus     float64 `json:"focus,omitempty"`
	Gap       float64 `json:"gap,omitempty"`
}
type EndBinding struct {
	ElementID string  `json:"elementId,omitempty"`
	Focus     float64 `json:"focus,omitempty"`
	Gap       float64 `json:"gap,omitempty"`
}

type Roundness struct {
	Type int `json:"type,omitempty"`
}

type BoundElements struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
type Elements struct {
	Type               string          `json:"type"`
	Version            int             `json:"version"`
	VersionNonce       int             `json:"versionNonce"`
	IsDeleted          bool            `json:"isDeleted"`
	ID                 string          `json:"id"`
	FillStyle          string          `json:"fillStyle"`
	StrokeWidth        int             `json:"strokeWidth"`
	StrokeStyle        string          `json:"strokeStyle"`
	Roughness          int             `json:"roughness"`
	Opacity            int             `json:"opacity"`
	Angle              int             `json:"angle"`
	X                  float64         `json:"x"`
	Y                  float64         `json:"y"`
	StrokeColor        string          `json:"strokeColor"`
	BackgroundColor    string          `json:"backgroundColor"`
	Width              float64         `json:"width"`
	Height             float64         `json:"height"`
	Seed               int             `json:"seed"`
	GroupIds           []interface{}   `json:"groupIds"`
	FrameID            interface{}     `json:"frameId"`
	Roundness          Roundness       `json:"roundness,omitempty"`
	BoundElements      []BoundElements `json:"boundElements"`
	Updated            int64           `json:"updated"`
	Link               interface{}     `json:"link"`
	Locked             bool            `json:"locked"`
	FontSize           float64         `json:"fontSize,omitempty"`
	FontFamily         int             `json:"fontFamily,omitempty"`
	Text               string          `json:"text,omitempty"`
	TextAlign          string          `json:"textAlign,omitempty"`
	VerticalAlign      string          `json:"verticalAlign,omitempty"`
	ContainerID        interface{}     `json:"containerId,omitempty"`
	OriginalText       string          `json:"originalText,omitempty"`
	LineHeight         float64         `json:"lineHeight,omitempty"`
	Baseline           int             `json:"baseline,omitempty"`
	StartBinding       StartBinding    `json:"startBinding,omitempty"`
	EndBinding         EndBinding      `json:"endBinding,omitempty"`
	LastCommittedPoint interface{}     `json:"lastCommittedPoint"`
	StartArrowhead     interface{}     `json:"startArrowhead"`
	EndArrowhead       interface{}     `json:"endArrowhead"` // not omitempty, null special
	IsFrameName        bool            `json:"isFrameName,omitempty"`
	Points             [][]float64     `json:"points,omitempty"`
}
type AppState struct {
	GridSize            interface{} `json:"gridSize"`
	ViewBackgroundColor string      `json:"viewBackgroundColor"`
}
type Files struct {
}
