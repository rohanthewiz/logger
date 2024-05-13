package teams_log

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/rohanthewiz/serr"
)

func SendLog(msg MessageCard, url string) (err error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return serr.Wrap(err, "Unable to marshal message card")
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(msgBytes))
	if err != nil {
		return serr.Wrap(err, "Post to Teams connector failed")
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		rb, err := io.ReadAll(resp.Body)
		if err != nil {
			return serr.Wrap(err, "when", "error marshalling response body", "code", strconv.Itoa(resp.StatusCode))
		}
		return serr.New("Non-200 response code", "code", strconv.Itoa(resp.StatusCode), "body", string(rb))
	}

	return
}
