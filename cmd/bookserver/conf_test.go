package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf(t *testing.T) {
	tt := []struct {
		conf  Config
		valid bool
	}{
		{Config{"", ""}, false},
		{Config{"80", ""}, false},
		{Config{"", "uri"}, false},
		{Config{"80", "uri"}, true},
	}
	for i, tc := range tt {
		t.Run(fmt.Sprintf("%v: %+v", i, tc.conf), func(t *testing.T) {
			assert.Equal(t, tc.valid, tc.conf.IsValid())
		})
	}
}
