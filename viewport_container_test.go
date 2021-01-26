package cockpit_stream

import (
	"fmt"
	"image"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewViewportContainer(t *testing.T) {
	container := NewViewportContainer()

	assert.IsType(t, &ViewportContainer{}, container)
}

func TestViewportContainer_Add(t *testing.T) {
	testId := "test_viewport"
	viewport := NewViewport(testId, 0, 0, 10, 10)

	container := NewViewportContainer()
	container.Add(testId, 0, 0, 10, 10)

	assert.Equal(t, viewport, container.data[testId])
}

func TestViewportContainer_Get(t *testing.T) {
	testId := "test_viewport"
	viewport := NewViewport(testId, 0, 0, 10, 10)

	container := NewViewportContainer()
	container.data[testId] = viewport

	fetchedVp, err := container.Get(testId)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, viewport, fetchedVp)
	assert.Equal(t, viewport.name, fetchedVp.name)
}

func TestViewportContainer_Has(t *testing.T) {
	testId := "test_viewport"
	viewport := &Viewport{name: testId}

	container := NewViewportContainer()
	container.data[testId] = viewport

	assert.True(t, container.Has(testId))
}

func TestViewportContainer_Count(t *testing.T) {
	container := NewViewportContainer()

	container.data["vp 1"] = &Viewport{name: "vp 1"}
	container.data["vp 2"] = &Viewport{name: "vp 2"}

	assert.Len(t, container.data, 2)
}

func TestViewportContainer_Bounds(t *testing.T) {
	for _, tc := range []struct {
		viewports      []*Viewport
		expectedBounds image.Rectangle
	}{
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			expectedBounds: image.Rect(0, 0, 200, 150),
		},
		{
			viewports: []*Viewport{
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			expectedBounds: image.Rect(0, 50, 200, 150),
		},
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
			},
			expectedBounds: image.Rect(0, 0, 100, 150),
		},
	} {
		container := NewViewportContainer()

		for _, vp := range tc.viewports {
			container.data[vp.name] = vp
		}
		container.recomputeBounds()

		assert.Equal(t, tc.expectedBounds, container.Bounds())
	}
}

func TestViewportContainer_Offset(t *testing.T) {
	for _, tc := range []struct {
		viewports      []*Viewport
		expectedOffset image.Point
	}{
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			expectedOffset: image.Point{X: 0, Y: 0},
		},
		{
			viewports: []*Viewport{
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			expectedOffset: image.Point{X: 0, Y: 50},
		},
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
			},
			expectedOffset: image.Point{X: 0, Y: 0},
		},
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			expectedOffset: image.Point{X: 50, Y: 0},
		},
	} {
		container := NewViewportContainer()

		for _, vp := range tc.viewports {
			container.data[vp.name] = vp
		}
		container.recomputeBounds()

		assert.Equal(t, tc.expectedOffset, container.Offset())
	}
}

func TestViewportContainer_ViewportOffset(t *testing.T) {
	for idx, tc := range []struct {
		viewports      []*Viewport
		testId         string
		expectedOffset image.Point
	}{
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			testId:         "green",
			expectedOffset: image.Point{X: 150, Y: 50},
		},
		{
			viewports: []*Viewport{
				NewViewport("blue", 0, 100, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			testId:         "blue",
			expectedOffset: image.Point{X: 0, Y: 50},
		},
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("blue", 0, 100, 50, 50),
			},
			testId:         "red",
			expectedOffset: image.Point{X: 50, Y: 0},
		},
		{
			viewports: []*Viewport{
				NewViewport("red", 50, 0, 50, 50),
				NewViewport("green", 150, 50, 50, 50),
			},
			testId:         "green",
			expectedOffset: image.Point{X: 100, Y: 50},
		},
	} {
		container := NewViewportContainer()

		for _, vp := range tc.viewports {
			container.data[vp.name] = vp
		}
		container.recomputeBounds()

		offset, err := container.ViewportOffset(container.data[tc.testId])
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, tc.expectedOffset, offset, fmt.Sprintf("test index %d \ntestId: %s \nbounds: %v \nboundsOffset: %v\n", idx, tc.testId, container.bounds, container.boundsOffset))
	}
}

func TestViewportContainer_Each(t *testing.T) {
	for _, tc := range []struct {
		viewports     []*Viewport
		expectedCalls int
	}{
		{
			viewports: []*Viewport{
				&Viewport{name: "vp 1"},
			},
			expectedCalls: 1,
		},
		{
			viewports: []*Viewport{
				&Viewport{name: "vp 1"},
				&Viewport{name: "vp 2"},
			},
			expectedCalls: 2,
		},
		{
			viewports: []*Viewport{
				&Viewport{name: "vp 1"},
				&Viewport{name: "vp 2"},
				&Viewport{name: "vp 3"},
			},
			expectedCalls: 3,
		},
	} {
		container := NewViewportContainer()
		for _, vp := range tc.viewports {
			container.data[vp.name] = vp
		}
		container.recomputeBounds()

		calls := 0
		container.Each(func(name string, viewport *Viewport) {
			calls++
		})

		assert.Equal(t, tc.expectedCalls, calls)
	}
}
