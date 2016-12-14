// +build medium

package elasticsearch

import (
        "testing"
        "time"

        "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"

        . "github.com/smartystreets/goconvey/convey"
)

func TestFilePublisher(t *testing.T) {
        fp := FPublisher{}

        Convey("Test publish", t, func() {
                Convey("Publish without a config file", func() {
                        metrics := []plugin.Metric{
                                plugin.Metric{
                                        Namespace: plugin.NewNamespace("x", "y", "z"),
                                        Config:    map[string]interface{}{"pw": "123aB"},
                                        Data:      3,
                                        Tags:      map[string]string{"hello": "world"},
                                        Unit:      "int",
                                        Timestamp: time.Now(),
                                },
                        }
                        err := fp.Publish(metrics, plugin.Config{})
                        So(err, ShouldEqual, plugin.ErrConfigNotFound)
                })
                Convey("Publish with a config file", func() {
                        metrics := []plugin.Metric{
                                plugin.Metric{
                                        Namespace: plugin.NewNamespace("x", "y", "z"),
                                        Config:    map[string]interface{}{"pw": "abc123"},
                                        Data:      3,
                                        Tags:      map[string]string{"hello": "world"},
                                        Unit:      "int",
                                        Timestamp: time.Now(),
                                },
                        }
                        err := fp.Publish(metrics, plugin.Config{"file": "/tmp/file_publisher_test.log"})
                        So(err, ShouldBeNil)
                })
                Convey("Test GetConfigPolicy", func() {
                        fp := FPublisher{}
                        _, err := fp.GetConfigPolicy()

                        Convey("No error returned", func() {
                                So(err, ShouldBeNil)
                        })
                })
        })
}
