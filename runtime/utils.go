package runtime

import (
	"image"
	"math"
)

func RotatePoint(px, py, cx, cy, angle float64) (float64, float64) {
	s := math.Sin(angle)
	c := math.Cos(angle)

	// Punkt zurück zum Ursprung verschieben
	px -= cx
	py -= cy

	// Punkt rotieren
	newX := px*c - py*s
	newY := px*s + py*c

	// Punkt wieder zurück verschieben
	newX += cx
	newY += cy
	return newX, newY
}
func RotateRect(rect image.Rectangle, angle float64) image.Rectangle {
	// Mitte des Rechtecks finden
	cx := float64(rect.Min.X + rect.Dx()/2)
	cy := float64(rect.Min.Y + rect.Dy()/2)

	// Eckpunkte des Rechtecks
	points := []struct{ x, y float64 }{
		{float64(rect.Min.X), float64(rect.Min.Y)},
		{float64(rect.Max.X), float64(rect.Min.Y)},
		{float64(rect.Max.X), float64(rect.Max.Y)},
		{float64(rect.Min.X), float64(rect.Max.Y)},
	}

	// Rotierte Punkte
	rotatedPoints := make([]struct{ x, y float64 }, 4)
	for i, p := range points {
		rotatedPoints[i].x, rotatedPoints[i].y = RotatePoint(p.x, p.y, cx, cy, angle)
	}

	// Neue Minima und Maxima finden
	minX, minY := rotatedPoints[0].x, rotatedPoints[0].y
	maxX, maxY := rotatedPoints[0].x, rotatedPoints[0].y
	for _, p := range rotatedPoints[1:] {
		if p.x < minX {
			minX = p.x
		}
		if p.y < minY {
			minY = p.y
		}
		if p.x > maxX {
			maxX = p.x
		}
		if p.y > maxY {
			maxY = p.y
		}
	}

	return image.Rect(int(minX), int(minY), int(maxX), int(maxY))
}
