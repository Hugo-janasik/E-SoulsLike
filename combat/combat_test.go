package combat

import (
	"testing"
)

func TestGetDistance(t *testing.T) {
	tests := []struct {
		name     string
		x1, y1   float64
		x2, y2   float64
		expected float64
	}{
		{"Same point", 0, 0, 0, 0, 0},
		{"Horizontal distance", 0, 0, 10, 0, 10},
		{"Vertical distance", 0, 0, 0, 10, 10},
		{"Diagonal distance", 0, 0, 3, 4, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetDistance(tt.x1, tt.y1, tt.x2, tt.y2)
			if result != tt.expected {
				t.Errorf("GetDistance(%v, %v, %v, %v) = %v; want %v",
					tt.x1, tt.y1, tt.x2, tt.y2, result, tt.expected)
			}
		})
	}
}

func TestCheckCollision(t *testing.T) {
	tests := []struct {
		name     string
		x1, y1   float64
		r1       float64
		x2, y2   float64
		r2       float64
		expected bool
	}{
		{"No collision", 0, 0, 10, 100, 0, 10, false},
		{"Touching circles exactly", 0, 0, 10, 20, 0, 10, false}, // Juste toucher = pas de collision
		{"Overlapping circles", 0, 0, 10, 15, 0, 10, true},
		{"Same position", 0, 0, 10, 0, 0, 10, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckCollision(tt.x1, tt.y1, tt.r1, tt.x2, tt.y2, tt.r2)
			if result != tt.expected {
				t.Errorf("CheckCollision() = %v; want %v", result, tt.expected)
			}
		})
	}
}

func TestGetDirection(t *testing.T) {
	tests := []struct {
		name      string
		x1, y1    float64
		x2, y2    float64
		expectedX float64
		expectedY float64
	}{
		{"Right direction", 0, 0, 10, 0, 1, 0},
		{"Left direction", 10, 0, 0, 0, -1, 0},
		{"Up direction", 0, 10, 0, 0, 0, -1},
		{"Down direction", 0, 0, 0, 10, 0, 1},
		{"Same point", 0, 0, 0, 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dx, dy := GetDirection(tt.x1, tt.y1, tt.x2, tt.y2)

			// Tolérance pour les calculs flottants
			tolerance := 0.0001

			if abs(dx-tt.expectedX) > tolerance || abs(dy-tt.expectedY) > tolerance {
				t.Errorf("GetDirection(%v, %v, %v, %v) = (%v, %v); want (%v, %v)",
					tt.x1, tt.y1, tt.x2, tt.y2, dx, dy, tt.expectedX, tt.expectedY)
			}
		})
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
