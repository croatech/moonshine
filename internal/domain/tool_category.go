package domain

type ToolCategory struct {
	Model
	Name      string      `json:"name"`
	Type      string      `json:"type"`
	ToolItems []*ToolItem `json:"tool_items,omitempty"`
}
