package hw03frequencyanalysis

import (
	"reflect"
	"testing"
)

func TestText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "Топ-10 из простого текста",
			input: `apple banana apple cherry apple
					banana`,
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name: "С учётом регистра и пустых слов",
			input: `Apple apple APPLE
					orange`,
			expected: []string{"APPLE", "Apple", "apple", "orange"},
		},
		{
			name: "Текст с переносами",
			input: `Привет
    				мир привет
    				мир мир`,
			expected: []string{"мир", "Привет", "привет"},
		},
		{
			name:     "Одно слово много раз",
			input:    `test test test test test test test test test test`,
			expected: []string{"test"},
		},
		{
			name:     "Разные слова, одинаковая частота",
			input:    `a b c a b c a b c`,
			expected: []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Top10(tt.input)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ожидалось %v, получено %v", tt.expected, got)
			}
		})
	}
}
