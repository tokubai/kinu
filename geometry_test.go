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
		t.Error("width must be 200. actual", geometry.Width)
	}
	if geometry.Height != 150 {
		t.Error("height must be 150. actual", geometry.Height)
	}
	if geometry.Quality != 85 {
		t.Error("quality must be 85. actual", geometry.Quality)
	}
	if geometry.NeedsAutoCrop != true {
		t.Error("needsAutoCrop must be true. actual", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != true {
		t.Error("needsOriginalImage must be true. actual", geometry.NeedsOriginalImage)
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
		t.Error("quality must be %d. actual %d", DEFAULT_QUALITY, geometry.Quality)
	}
	if geometry.NeedsAutoCrop != false {
		t.Error("needsAutoCrop must be false. actual", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != false {
		t.Error("needsOriginalImage must be false. actual", geometry.NeedsOriginalImage)
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
		t.Error("quality must be %d. actual %d", DEFAULT_QUALITY, geometry.Quality)
	}
	if geometry.NeedsAutoCrop != false {
		t.Error("needsAutoCrop must be false. actual", geometry.NeedsAutoCrop)
	}
	if geometry.NeedsOriginalImage != false {
		t.Error("needsOriginalImage must be false. actual", geometry.NeedsOriginalImage)
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
