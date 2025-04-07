package commit

// CommitTypeItem represents a commit type option
type CommitTypeItem struct {
	TypeTitle       string
	TypeDescription string
}

// Title implements the list.Item interface
func (i CommitTypeItem) Title() string { return i.TypeTitle }
func (i CommitTypeItem) Description() string { return i.TypeDescription }
func (i CommitTypeItem) FilterValue() string { return i.TypeTitle }
