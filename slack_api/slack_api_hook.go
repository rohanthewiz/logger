package slack_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// SlackAPIHook is a logrus hook for sending logs to Slack via the Web API
type SlackAPIHook struct {
	Token          string
	Channel        string
	AcceptedLevels []logrus.Level
	Enabled        bool
	UseBlocks      bool // Whether to use rich block formatting or simple messages
}

// NewSlackAPIHook creates a new Slack API hook
func NewSlackAPIHook(token, channel string, acceptedLevels []logrus.Level, useBlocks bool) *SlackAPIHook {
	return &SlackAPIHook{
		Token:          token,
		Channel:        channel,
		AcceptedLevels: acceptedLevels,
		Enabled:        true,
		UseBlocks:      useBlocks,
	}
}

// Levels returns the log levels this hook should fire for
func (h *SlackAPIHook) Levels() []logrus.Level {
	if h.AcceptedLevels == nil {
		return logrus.AllLevels
	}
	return h.AcceptedLevels
}

// Fire is called when a log event is fired
func (h *SlackAPIHook) Fire(entry *logrus.Entry) error {
	if !h.Enabled {
		return nil
	}

	// Send asynchronously to avoid blocking
	go func() {
		var payload interface{}
		if h.UseBlocks {
			payload = h.createBlockMessage(entry)
		} else {
			payload = h.createSimpleMessage(entry)
		}

		if err := h.sendToSlack(payload); err != nil {
			// Print to stdout to avoid recursive logging
			fmt.Printf("Error sending log to Slack: %v\n", err)
		}
	}()

	return nil
}

// createSimpleMessage creates a simple text message for Slack
func (h *SlackAPIHook) createSimpleMessage(entry *logrus.Entry) map[string]interface{} {
	// Format the basic message
	text := fmt.Sprintf("*%s*: %s", strings.ToUpper(entry.Level.String()), entry.Message)

	// Add fields if present
	if len(entry.Data) > 0 {
		text += "\n\n*Fields:*"
		for key, value := range entry.Data {
			text += fmt.Sprintf("\n• `%s`: %v", key, value)
		}
	}

	return map[string]interface{}{
		"channel": h.Channel,
		"text":    text,
	}
}

// createBlockMessage creates a rich block-formatted message for Slack
func (h *SlackAPIHook) createBlockMessage(entry *logrus.Entry) map[string]interface{} {
	// Determine emoji and header based on level
	emoji, header := h.getEmojiAndHeader(entry.Level)
	
	// Create the blocks
	blocks := []interface{}{
		// Header block
		map[string]interface{}{
			"type": "header",
			"text": map[string]interface{}{
				"type": "plain_text",
				"text": fmt.Sprintf("%s %s", emoji, header),
			},
		},
	}

	// Add fields section if there are any
	if len(entry.Data) > 0 {
		fields := []map[string]interface{}{}
		
		// Extract common fields first
		commonFields := []string{"service", "environment", "error_type", "component", "module"}
		for _, field := range commonFields {
			if value, ok := entry.Data[field]; ok {
				fields = append(fields, map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("*%s:* `%v`", formatFieldName(field), value),
				})
			}
		}

		// Add timestamp
		fields = append(fields, map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*Timestamp:* `%s`", entry.Time.Format("2006-01-02 15:04:05 MST")),
		})

		// Add level
		fields = append(fields, map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*Level:* `%s`", strings.ToUpper(entry.Level.String())),
		})

		if len(fields) > 0 {
			blocks = append(blocks, map[string]interface{}{
				"type":   "section",
				"fields": fields,
			})
		}
	}

	// Add message section
	blocks = append(blocks, map[string]interface{}{
		"type": "section",
		"text": map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("*Message:* %s", entry.Message),
		},
	})

	// Add stack trace if present
	if stackTrace, ok := entry.Data["stack_trace"]; ok {
		blocks = append(blocks, 
			map[string]interface{}{
				"type": "divider",
			},
			map[string]interface{}{
				"type": "section",
				"text": map[string]interface{}{
					"type": "mrkdwn",
					"text": fmt.Sprintf("```\n%v\n```", stackTrace),
				},
			},
		)
	}

	// Add remaining fields as context
	var contextElements []map[string]interface{}
	for key, value := range entry.Data {
		// Skip already displayed fields
		if isCommonField(key) || key == "stack_trace" {
			continue
		}
		contextElements = append(contextElements, map[string]interface{}{
			"type": "mrkdwn",
			"text": fmt.Sprintf("`%s`: %v", key, value),
		})
	}

	if len(contextElements) > 0 {
		blocks = append(blocks, map[string]interface{}{
			"type":     "context",
			"elements": contextElements,
		})
	}

	// Add action buttons for error and fatal levels
	if entry.Level <= logrus.ErrorLevel {
		actions := []map[string]interface{}{}
		
		// Add log search button if log_url is provided
		if logURL, ok := entry.Data["log_url"]; ok {
			actions = append(actions, map[string]interface{}{
				"type": "button",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "View in Log System",
				},
				"url": fmt.Sprintf("%v", logURL),
			})
		}

		// Add incident button if incident_id is provided
		if incidentID, ok := entry.Data["incident_id"]; ok {
			actions = append(actions, map[string]interface{}{
				"type": "button",
				"text": map[string]interface{}{
					"type": "plain_text",
					"text": "View Incident",
				},
				"style": "danger",
				"value": fmt.Sprintf("incident_%v", incidentID),
			})
		}

		if len(actions) > 0 {
			blocks = append(blocks, map[string]interface{}{
				"type":     "actions",
				"elements": actions,
			})
		}
	}

	fallbackText := fmt.Sprintf("%s: %s", strings.ToUpper(entry.Level.String()), entry.Message)
	
	return map[string]interface{}{
		"channel": h.Channel,
		"text":    fallbackText, // Fallback for notifications
		"blocks":  blocks,
	}
}

// sendToSlack sends the payload to Slack
func (h *SlackAPIHook) sendToSlack(payload interface{}) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Slack API response
	var slackResp map[string]interface{}
	if err := json.Unmarshal(body, &slackResp); err != nil {
		return fmt.Errorf("failed to parse Slack response: %w", err)
	}

	// Check if Slack reported an error
	if ok, exists := slackResp["ok"].(bool); exists && !ok {
		if errMsg, exists := slackResp["error"].(string); exists {
			return fmt.Errorf("slack API error: %s", errMsg)
		}
		return fmt.Errorf("slack API returned ok=false")
	}

	return nil
}

// getEmojiAndHeader returns appropriate emoji and header text based on log level
func (h *SlackAPIHook) getEmojiAndHeader(level logrus.Level) (string, string) {
	switch level {
	case logrus.PanicLevel, logrus.FatalLevel:
		return "🚨", "Critical Error"
	case logrus.ErrorLevel:
		return "❌", "Error"
	case logrus.WarnLevel:
		return "⚠️", "Warning"
	case logrus.InfoLevel:
		return "ℹ️", "Information"
	case logrus.DebugLevel:
		return "🔍", "Debug"
	case logrus.TraceLevel:
		return "📝", "Trace"
	default:
		return "📋", "Log Entry"
	}
}

// formatFieldName converts snake_case to Title Case
func formatFieldName(field string) string {
	parts := strings.Split(field, "_")
	for i, part := range parts {
		if len(part) > 0 {
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}
	return strings.Join(parts, " ")
}

// isCommonField checks if a field is one of the common fields we display specially
func isCommonField(field string) bool {
	commonFields := []string{"service", "environment", "error_type", "component", "module"}
	for _, f := range commonFields {
		if f == field {
			return true
		}
	}
	return false
}