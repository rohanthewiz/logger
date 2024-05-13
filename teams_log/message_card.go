package teams_log

const messageCardType = "MessageCard"
const messageCardContext = "http://schema.org/extensions"

type MessageCard struct {
	Type            string        `json:"@type"`             // hardwired to "MessageCard"
	Context         string        `json:"@context"`          //
	Summary         string        `json:"summary,omitempty"` // not really displayed in our case, but may be required
	Sections        []Section     `json:"sections"`
	PotentialAction []interface{} `json:"potentialAction"`
}

type Section struct {
	ActivityTitle    string `json:"activityTitle,omitempty"`
	ActivitySubtitle string `json:"activitySubtitle,omitempty"`
	ActivityImage    string `json:"activityImage,omitempty"`
	ActivityText     string `json:"activityText,omitempty"`
	Facts            []Fact `json:"facts"`
}

type Fact struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
