package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/LaansDole/go-git-tui/internal/ui/common"
)

func TruncateText(text string, maxLength int, ellipsis string) string {
	return common.TruncateText(text, maxLength, ellipsis)
}

func TruncatePath(path string, maxLength int, prefixChars int, suffixChars int) string {
	return common.TruncatePath(path, maxLength, prefixChars, suffixChars)
}

func RenderText(text string, style lipgloss.Style) string {
	return common.RenderText(text, style)
}
