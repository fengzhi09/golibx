package gox

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMaxMinTime(t *testing.T) {
	vals := []time.Time{AsTime(int64(20221228)), AsTime(int64(20230101)), AsTime(int64(20221229)), AsTime(int64(20221230)), AsTime(int64(20230104)), AsTime(int64(20230102)), AsTime(int64(20221231)), AsTime(int64(20221227)), AsTime(int64(20230103))}
	max, min := MaxMinTime(vals...)
	assert.Equal(t, AsTime(int64(20230104)), max)
	assert.Equal(t, AsTime(int64(20221227)), min)
}
