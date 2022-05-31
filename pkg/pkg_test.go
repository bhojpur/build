package pkg

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocate(t *testing.T) {
	pcPath, _, err := locateGTK()
	assert.NoError(t, err)
	assert.Equal(t, "test/gtk/gtk+-3.0.pc", pcPath)
}

func TestLoad(t *testing.T) {
	pcPath, pc, err := locateGTK()
	if !assert.NoError(t, err) {
		return
	}
	err = pc.Load(pcPath, true)
	assert.NoError(t, err)
	expected := []string{
		"atk", "cairo", "cairo-gobject", "fontconfig", "freetype2", "gdk-3.0",
		"gdk-pixbuf-2.0", "gio-2.0", "glib-2.0", "gmodule-no-export-2.0", "gobject-2.0",
		"gthread-2.0", "gtk+-3.0", "pango", "pangocairo", "pangoft2",
	}
	assert.Equal(t, expected, pc.LoadedPkgNames())
}

func TestCFlagsGTK3(t *testing.T) {
	pcPath, pc, err := locateGTK()
	if !assert.NoError(t, err) {
		return
	}
	err = pc.Load(pcPath, true)
	assert.NoError(t, err)
	expected := []string{
		"-I/gtk/include/gtk-3.0", "-DGSEAL_ENABLE",
		"-I/gtk/include/pango-1.0", "-I/gtk/include/glib-2.0", "-I/gtk/lib/glib-2.0/include",
		"-pthread", "-I/gtk/include/cairo", "-I/gtk/include/gdk-pixbuf-2.0",
		"-I/gtk/include/atk-1.0", "-I/gtk/include/freetype2", "-I/gtk/include",
	}
	assert.Equal(t, expected, pc.CFlags())
}

func locateGTK() (string, Config, error) {
	pc, err := NewConfig([]string{"test/gtk"})
	if err != nil {
		return "", nil, err
	}
	pcPath, err := pc.Locate("gtk+-3.0")
	if err != nil {
		return "", nil, err
	}
	return pcPath, pc, nil
}
