package main

import (
	"testing"
)

func TestParseGeometry(t *testing.T) {
	geometry, err := ParseGeometry("w=200,h=150,q=85,c=true,o=true")
	if err != nil {
		t.Error("error is not nil")
	}
	if geometry.Width != 200 {
		t.Errorf("width must be 200. actual %d", geometry.Width)
	}
	if geometry.Height != 150 {
		t.Errorf("height must be 150. actual %d", geometry.Height)
	}
	if geometry.Quality != 85 {
		t.Errorf("quality must be 85. actual %d", geometry.Quality)
	}
	if geometry.NeedsAutoCrop != true {
		t.Errorf("needsAutoCrop must be true. actual %t", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != true {
		t.Errorf("needsOriginalImage must be true. actual %t", geometry.NeedsOriginalImage)
	}

	geometry, err = ParseGeometry("w=200")
	if err != nil {
		t.Error("error is not nil")
	}
	if geometry.Width != 200 {
		t.Errorf("width must be 200. actual. %d", geometry.Width)
	}
	if geometry.Height != 0 {
		t.Errorf("height must be 0. actual. %d", geometry.Height)
	}
	if geometry.Quality != DEFAULT_QUALITY {
		t.Errorf("quality must be %d. actual %d", DEFAULT_QUALITY, geometry.Quality)
	}
	if geometry.NeedsAutoCrop != false {
		t.Errorf("needsAutoCrop must be false. actual %t", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != false {
		t.Errorf("needsOriginalImage must be false. actual %t", geometry.NeedsOriginalImage)
	}

	geometry, err = ParseGeometry("h=200")
	if err != nil {
		t.Error("error is not nil")
	}
	if geometry.Width != 0 {
		t.Errorf("width must be 200. actual. %d", geometry.Width)
	}
	if geometry.Height != 200 {
		t.Errorf("height must be 2000. actual. %d", geometry.Height)
	}
	if geometry.Quality != DEFAULT_QUALITY {
		t.Errorf("quality must be %d. actual %d", DEFAULT_QUALITY, geometry.Quality)
	}
	if geometry.NeedsAutoCrop != false {
		t.Errorf("needsAutoCrop must be false. actual %t", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != false {
		t.Errorf("needsOriginalImage must be false. actual %t", geometry.NeedsOriginalImage)
	}

	geometry, err = ParseGeometry("")
	if geometry != nil {
		t.Error("geometry must be nil")
	}
	if err == nil {
		t.Error("error must not be nil")
	}

	geometry, err = ParseGeometry("invalid")
	if geometry != nil {
		t.Error("geometry must be nil")
	}
	if err == nil {
		t.Error("error must not be nil")
	}
}
