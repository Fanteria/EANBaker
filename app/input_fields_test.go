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
			field := NewInputField(tt.text, tt.fieldName, tt.suggestion)

			if field.name != tt.fieldName {
				t.Errorf("name = %v, want %v", field.name, tt.fieldName)
			}
			if field.suggestion != tt.suggestion {
				t.Errorf("suggestion = %v, want %v", field.suggestion, tt.suggestion)
			}
			if tt.text != "" && field.GetText() != tt.text {
				t.Errorf("GetText() = %v, want %v", field.GetText(), tt.text)
			}
		})
	}
}

func TestInputField_GetText(t *testing.T) {
	field := NewInputField("test value", "Test", "hint")
	
	got := field.GetText()
	if got != "test value" {
		t.Errorf("GetText() = %v, want %v", got, "test value")
	}
}

func TestInputField_SetText(t *testing.T) {
	field := NewInputField("initial", "Test", "hint")
	
	field.SetText("new value")
	got := field.GetText()
	if got != "new value" {
		t.Errorf("After SetText(), GetText() = %v, want %v", got, "new value")
	}
}

func TestInputField_SetText_Empty(t *testing.T) {
	field := NewInputField("initial", "Test", "hint")
	
	field.SetText("")
	got := field.GetText()
	if got != "" {
		t.Errorf("After SetText(''), GetText() = %v, want empty string", got)
	}
}

func TestInputField_SetText_Overwrite(t *testing.T) {
	field := NewInputField("first", "Test", "hint")
	
	field.SetText("second")
	field.SetText("third")
	
	got := field.GetText()
	if got != "third" {
		t.Errorf("GetText() = %v, want %v", got, "third")
	}
}

func TestInputField_SingleLine(t *testing.T) {
	field := NewInputField("", "Test", "hint")
	
	// Editor should be configured for single line
	if !field.editor.SingleLine {
		t.Error("Editor should be configured for single line")
	}
}

func TestInputField_EmptyInitialText(t *testing.T) {
	field := NewInputField("", "Test", "hint")
	
	got := field.GetText()
	if got != "" {
		t.Errorf("GetText() for empty initial = %v, want empty string", got)
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
			field := NewInputField(tt.text, "Test", "hint")
			got := field.GetText()
			if got != tt.text {
				t.Errorf("GetText() = %v, want %v", got, tt.text)
			}
		})
	}
}

func TestInputField_LongText(t *testing.T) {
	// Test with long text
	longText := ""
	for i := 0; i < 1000; i++ {
		longText += "x"
	}

	field := NewInputField(longText, "Test", "hint")
	got := field.GetText()
	if got != longText {
		t.Errorf("Long text not preserved, got length %d, want %d", len(got), len(longText))
	}
}
