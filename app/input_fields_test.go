package app

import (
	"testing"
)

func TestNewInputField(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		fieldName  string
		suggestion string
	}{
		{
			name:       "Empty text",
			text:       "",
			fieldName:  "Test",
			suggestion: "Enter value",
		},
		{
			name:       "With initial text",
			text:       "initial value",
			fieldName:  "Name",
			suggestion: "Enter name",
		},
		{
			name:       "All fields set",
			text:       "value",
			fieldName:  "Field",
			suggestion: "Hint",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storedValue := tt.text
			msg := &Message{}
			field := NewInputField(
				tt.fieldName,
				tt.suggestion,
				msg,
				func(value string) error {
					storedValue = value
					return nil
				},
				func() string {
					return storedValue
				},
			)

			if field.name != tt.fieldName {
				t.Errorf("name = %v, want %v", field.name, tt.fieldName)
			}
			if field.suggestion != tt.suggestion {
				t.Errorf("suggestion = %v, want %v", field.suggestion, tt.suggestion)
			}
			if tt.text != "" && field.editor.Text() != tt.text {
				t.Errorf("editor.Text() = %v, want %v", field.editor.Text(), tt.text)
			}
		})
	}
}

func TestInputField_SingleLine(t *testing.T) {
	msg := &Message{}
	value := ""
	field := NewInputField(
		"Test",
		"hint",
		msg,
		func(v string) error {
			value = v
			return nil
		},
		func() string {
			return value
		},
	)

	// Editor should be configured for single line
	if !field.editor.SingleLine {
		t.Error("Editor should be configured for single line")
	}
}

func TestInputField_EmptyInitialText(t *testing.T) {
	msg := &Message{}
	value := ""
	field := NewInputField(
		"Test",
		"hint",
		msg,
		func(v string) error {
			value = v
			return nil
		},
		func() string {
			return value
		},
	)

	if field.editor.Text() != "" {
		t.Errorf("editor.Text() for empty initial = %v, want empty string", field.editor.Text())
	}
}

func TestInputField_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name string
		text string
	}{
		{name: "Unicode", text: "Příliš žluťoučký kůň"},
		{name: "Emoji", text: "🎉 Party 🎊"},
		{name: "Special chars", text: "Test™ ® © €"},
		{name: "Whitespace", text: "  spaces  "},
		{name: "Tabs", text: "tab\there"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &Message{}
			value := tt.text
			field := NewInputField(
				"Test",
				"hint",
				msg,
				func(v string) error {
					value = v
					return nil
				},
				func() string {
					return value
				},
			)
			got := field.editor.Text()
			if got != tt.text {
				t.Errorf("editor.Text() = %v, want %v", got, tt.text)
			}
		})
	}
}

func TestInputField_LongText(t *testing.T) {
	// Test with long text
	longText := ""
	for range 1000 {
		longText += "x"
	}

	msg := &Message{}
	value := longText
	field := NewInputField(
		"Test",
		"hint",
		msg,
		func(v string) error {
			value = v
			return nil
		},
		func() string {
			return value
		},
	)
	got := field.editor.Text()
	if got != longText {
		t.Errorf("Long text not preserved, got length %d, want %d", len(got), len(longText))
	}
}
