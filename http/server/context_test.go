package server

import (
	"context"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	"github.com/RavenHuo/daenerys/log"
	"github.com/magiconair/properties/assert"
)

func TestContext_SetGet(t *testing.T) {
	ctx := Context{}

	ctx.Set("foo", "bar")

	val, ok := ctx.Get("foo")
	assert.Equal(t, true, ok)
	assert.Equal(t, "bar", val)

	val, ok = ctx.Get("foo2")
	assert.Equal(t, false, ok)
	assert.Equal(t, nil, val)

	now := time.Now()
	ctx.Set("foo-time", now)
	timeVal := ctx.GetTime("foo-time")
	assert.Equal(t, now, timeVal)

	duration := 10 * time.Second
	ctx.Set("foo-duration", duration)
	durationVal := ctx.GetDuration("foo-duration")
	assert.Equal(t, duration, durationVal)

	slice := []string{"1", "2", "3"}
	ctx.Set("foo-slice", slice)
	sliceVal := ctx.GetStringSlice("foo-slice")
	assert.Equal(t, slice, sliceVal)

	stringMap := map[string]interface{}{"1": "1", "2": "2", "3": "3"}
	ctx.Set("foo-map", stringMap)
	mapVal := ctx.GetStringMap("foo-map")
	assert.Equal(t, stringMap, mapVal)
}

func TestContext_MustGet(t *testing.T) {
	defer func() {
		if rc := recover(); rc != nil {
			// "TestContext_MustGet got panic, stacks:%s"
			//log.Logger.Errorf(context.Background(), context.Background(), string(debug.Stack()))
			log.Errorf(context.Background(), "TestContext_MustGet got panic, stacks:%s", string(debug.Stack()))
			val, ok := rc.(string)
			if ok && strings.Contains(val, "does not exist") {
				t.Logf("panic error:%s", val)
			} else {
				t.Fail()
			}
		}
	}()

	ctx := Context{}

	ctx.Set("foo", "bar")

	val := ctx.MustGet("foo")
	assert.Equal(t, "bar", val)

	val = ctx.MustGet("foo2")
}
