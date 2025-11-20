package ui

import (
	"fmt"
	"strings"
	"time"
)

// LogEntry represents a single battle log entry
type LogEntry struct {
	Message   string
	Timestamp time.Time
	Type      LogEntryType
}

// LogEntryType categorizes log entries for color coding
type LogEntryType string

const (
	LogTypeDamage   LogEntryType = "damage"
	LogTypeHealing  LogEntryType = "healing"
	LogTypeStatus   LogEntryType = "status"
	LogTypeAction   LogEntryType = "action"
	LogTypeInfo     LogEntryType = "info"
	LogTypeWarning  LogEntryType = "warning"
	LogTypeError    LogEntryType = "error"
	LogTypeVictory  LogEntryType = "victory"
	LogTypeDefeat   LogEntryType = "defeat"
)

// BattleLog manages battle log entries
type BattleLog struct {
	Entries    []LogEntry
	MaxEntries int
}

// NewBattleLog creates a new battle log with a maximum entry limit
func NewBattleLog(maxEntries int) *BattleLog {
	if maxEntries <= 0 {
		maxEntries = 10
	}
	return &BattleLog{
		Entries:    make([]LogEntry, 0, maxEntries),
		MaxEntries: maxEntries,
	}
}

// Add adds a new entry to the battle log
func (bl *BattleLog) Add(message string, entryType LogEntryType) {
	entry := LogEntry{
		Message:   message,
		Timestamp: time.Now(),
		Type:      entryType,
	}

	bl.Entries = append(bl.Entries, entry)

	// Keep only the last MaxEntries
	if len(bl.Entries) > bl.MaxEntries {
		bl.Entries = bl.Entries[len(bl.Entries)-bl.MaxEntries:]
	}
}

// AddFromString adds a log entry from a plain string, auto-detecting type
func (bl *BattleLog) AddFromString(message string) {
	entryType := detectLogType(message)
	bl.Add(message, entryType)
}

// Clear clears all log entries
func (bl *BattleLog) Clear() {
	bl.Entries = make([]LogEntry, 0, bl.MaxEntries)
}

// GetRecent returns the most recent N entries
func (bl *BattleLog) GetRecent(n int) []LogEntry {
	if n <= 0 || n > len(bl.Entries) {
		return bl.Entries
	}
	return bl.Entries[len(bl.Entries)-n:]
}

// detectLogType attempts to detect the log entry type from the message
func detectLogType(message string) LogEntryType {
	msgLower := strings.ToLower(message)

	// Check for damage indicators
	if strings.Contains(msgLower, "damage") || strings.Contains(msgLower, "dealt") ||
		strings.Contains(msgLower, "hit") || strings.Contains(msgLower, "attacked") {
		return LogTypeDamage
	}

	// Check for healing indicators
	if strings.Contains(msgLower, "heal") || strings.Contains(msgLower, "restored") ||
		strings.Contains(msgLower, "recovered") {
		return LogTypeHealing
	}

	// Check for status indicators
	if strings.Contains(msgLower, "defended") || strings.Contains(msgLower, "passed") ||
		strings.Contains(msgLower, "sacrificed") || strings.Contains(msgLower, "switched") {
		return LogTypeStatus
	}

	// Check for victory/defeat
	if strings.Contains(msgLower, "victory") || strings.Contains(msgLower, "won") ||
		strings.Contains(msgLower, "wins") {
		return LogTypeVictory
	}
	if strings.Contains(msgLower, "defeat") || strings.Contains(msgLower, "lost") ||
		strings.Contains(msgLower, "loses") {
		return LogTypeDefeat
	}

	// Check for warnings
	if strings.Contains(msgLower, "warning") || strings.Contains(msgLower, "not enough") ||
		strings.Contains(msgLower, "cannot") {
		return LogTypeWarning
	}

	// Check for errors
	if strings.Contains(msgLower, "error") || strings.Contains(msgLower, "failed") {
		return LogTypeError
	}

	// Default to info
	return LogTypeInfo
}

// RenderBattleLog renders the battle log with color coding
func (r *Renderer) RenderBattleLog(entries []LogEntry, maxLines int) string {
	var result strings.Builder

	// Header
	result.WriteString("╔")
	result.WriteString(strings.Repeat("═", r.Width-2))
	result.WriteString("╗\n")

	headerText := "BATTLE LOG"
	headerPadding := (r.Width - len(headerText) - 2) / 2
	result.WriteString("║")
	result.WriteString(strings.Repeat(" ", headerPadding))
	if r.ColorSupport {
		result.WriteString(Colorize(headerText, Bold+ColorBrightCyan))
	} else {
		result.WriteString(headerText)
	}
	result.WriteString(strings.Repeat(" ", r.Width-len(headerText)-headerPadding-2))
	result.WriteString("║\n")

	result.WriteString("╠")
	result.WriteString(strings.Repeat("═", r.Width-2))
	result.WriteString("╣\n")

	// Determine which entries to show
	displayEntries := entries
	if maxLines > 0 && len(entries) > maxLines {
		displayEntries = entries[len(entries)-maxLines:]
	}

	// Render entries
	if len(displayEntries) == 0 {
		emptyMsg := "No battle events yet..."
		padding := (r.Width - len(emptyMsg) - 2) / 2
		result.WriteString("║")
		result.WriteString(strings.Repeat(" ", padding))
		if r.ColorSupport {
			result.WriteString(Colorize(emptyMsg, ColorGray))
		} else {
			result.WriteString(emptyMsg)
		}
		result.WriteString(strings.Repeat(" ", r.Width-len(emptyMsg)-padding-2))
		result.WriteString("║\n")
	} else {
		for _, entry := range displayEntries {
			result.WriteString(r.renderLogEntry(entry))
		}
	}

	// Footer
	result.WriteString("╚")
	result.WriteString(strings.Repeat("═", r.Width-2))
	result.WriteString("╝\n")

	return result.String()
}

// RenderBattleLogSimple renders battle log from string slice (for compatibility)
func (r *Renderer) RenderBattleLogSimple(messages []string, maxLines int) string {
	entries := make([]LogEntry, len(messages))
	for i, msg := range messages {
		entries[i] = LogEntry{
			Message:   msg,
			Timestamp: time.Now(),
			Type:      detectLogType(msg),
		}
	}
	return r.RenderBattleLog(entries, maxLines)
}

// renderLogEntry renders a single log entry with color coding
func (r *Renderer) renderLogEntry(entry LogEntry) string {
	var result strings.Builder

	// Format timestamp
	timestamp := entry.Timestamp.Format("15:04:05")
	
	// Prefix with timestamp
	prefix := fmt.Sprintf("[%s] ", timestamp)
	
	// Color code based on type
	message := entry.Message
	if r.ColorSupport {
		switch entry.Type {
		case LogTypeDamage:
			message = Colorize(message, ColorRed)
		case LogTypeHealing:
			message = Colorize(message, ColorGreen)
		case LogTypeStatus:
			message = Colorize(message, ColorYellow)
		case LogTypeAction:
			message = Colorize(message, ColorCyan)
		case LogTypeWarning:
			message = Colorize(message, ColorBrightYellow)
		case LogTypeError:
			message = Colorize(message, ColorBrightRed)
		case LogTypeVictory:
			message = Colorize(message, Bold+ColorBrightGreen)
		case LogTypeDefeat:
			message = Colorize(message, Bold+ColorBrightRed)
		case LogTypeInfo:
			message = Colorize(message, ColorWhite)
		}
		prefix = Colorize(prefix, ColorGray)
	}

	fullMessage := prefix + message

	// Word wrap if message is too long
	maxWidth := r.Width - 4 // Account for borders and padding
	lines := wordWrap(stripANSI(fullMessage), maxWidth)

	for i, line := range lines {
		result.WriteString("║ ")
		
		if i == 0 {
			// First line with full message (including ANSI codes)
			displayLen := len(stripANSI(fullMessage))
			if displayLen <= maxWidth {
				result.WriteString(fullMessage)
				result.WriteString(strings.Repeat(" ", maxWidth-displayLen))
			} else {
				// Truncate with ANSI codes preserved
				result.WriteString(truncateWithANSI(fullMessage, maxWidth-3))
				result.WriteString("...")
			}
		} else {
			// Continuation lines
			result.WriteString("    " + line)
			result.WriteString(strings.Repeat(" ", maxWidth-len(line)-4))
		}
		
		result.WriteString(" ║\n")
	}

	return result.String()
}

// wordWrap wraps text to fit within the specified width
func wordWrap(text string, width int) []string {
	if len(text) <= width {
		return []string{text}
	}

	var lines []string
	words := strings.Fields(text)
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len() == 0 {
			currentLine.WriteString(word)
		} else if currentLine.Len()+1+len(word) <= width {
			currentLine.WriteString(" ")
			currentLine.WriteString(word)
		} else {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// truncateWithANSI truncates a string with ANSI codes to a visible length
func truncateWithANSI(s string, visibleLength int) string {
	var result strings.Builder
	visibleCount := 0
	inEscape := false

	for i := 0; i < len(s); i++ {
		if s[i] == '\033' {
			inEscape = true
			result.WriteByte(s[i])
			continue
		}

		if inEscape {
			result.WriteByte(s[i])
			if s[i] == 'm' {
				inEscape = false
			}
			continue
		}

		if visibleCount >= visibleLength {
			break
		}

		result.WriteByte(s[i])
		visibleCount++
	}

	// Add reset code if we're in the middle of an escape sequence
	if inEscape || strings.Contains(s, "\033[") {
		result.WriteString(ColorReset)
	}

	return result.String()
}

// RenderCompactLog renders a compact version of the battle log (inline)
func (r *Renderer) RenderCompactLog(entries []LogEntry, maxEntries int) string {
	var result strings.Builder

	displayEntries := entries
	if maxEntries > 0 && len(entries) > maxEntries {
		displayEntries = entries[len(entries)-maxEntries:]
	}

	for i, entry := range displayEntries {
		message := entry.Message
		if r.ColorSupport {
			switch entry.Type {
			case LogTypeDamage:
				message = Colorize("▸ "+message, ColorRed)
			case LogTypeHealing:
				message = Colorize("▸ "+message, ColorGreen)
			case LogTypeStatus:
				message = Colorize("▸ "+message, ColorYellow)
			default:
				message = "▸ " + message
			}
		} else {
			message = "> " + message
		}

		result.WriteString(message)
		if i < len(displayEntries)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}
