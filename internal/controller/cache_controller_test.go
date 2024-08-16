package controller

import (
	"testing"

	"github.com/ernestngugi/todo/internal/mocks"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCacheController(t *testing.T) {

	redisProvider := mocks.NewMockRedisProvider()

	cacheController := NewTestCacheController(redisProvider)

	Convey("TestCacheController", t, func() {

		Convey("can add something to cache", func() {

			err := cacheController.CacheValue("key", "value")
			So(err, ShouldBeNil)
		})

		Convey("can get a cached value", func() {

			err := cacheController.CacheValue("key1", "value1")
			So(err, ShouldBeNil)

			var value string

			err = cacheController.GetCachedValue("key1", &value)
			So(err, ShouldBeNil)

			So(value, ShouldEqual, "value1")
		})

		Convey("can check if a key already exists", func() {

			err := cacheController.CacheValue("key", "value")
			So(err, ShouldBeNil)

			exist, err := cacheController.Exists("key")
			So(err, ShouldBeNil)
			So(exist, ShouldBeTrue)
		})

		Convey("can remove key from cache", func() {

			key := "key"

			err := cacheController.CacheValue(key, "value1")
			So(err, ShouldBeNil)

			err = cacheController.RemoveFromCache(key)
			So(err, ShouldBeNil)

			exists, err := cacheController.Exists(key)
			So(err, ShouldBeNil)
			So(exists, ShouldBeFalse)
		})
	})
}
