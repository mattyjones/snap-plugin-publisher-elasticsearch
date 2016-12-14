package elasticsearch

import (
        "bytes"
        "net/http"
        "encoding/json"
        "fmt"
        "strings"
        "time"
        "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
      	log "github.com/Sirupsen/logrus"
)

// FPublisher is a testing publisher.
type FPublisher struct {
}

/*
        GetConfigPolicy() returns the configPolicy for your plugin.

        A config policy is how users can provide configuration info to
        plugin. Here you define what sorts of config info your plugin
        needs and/or requires.
*/
func (f FPublisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
  policy := plugin.NewConfigPolicy()
  return *policy, nil
}

type MetricToPublish struct {
	// The timestamp from when the metric was created.
	Key string
  Value float64
  Timestamp string
}

// Publish test publish function
func (f FPublisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {
  url := "http://localhost:9200/snap/test"
	var tagsForPrefix []string

  tagConfigs, err := cfg.GetString("prefix_tags")
  	if err == nil {
  		tagsForPrefix = strings.Split(tagConfigs, ",")
  }
  // log.Debug("Start elastic publish")
  for _, m := range mts {
    key := strings.Join(m.Namespace.Strings(), ".")
    for _, tag := range tagsForPrefix {
      nextTag, ok := m.Tags[tag]
      if ok {
        key = nextTag + "." + key
      }
    }
    s := MetricToPublish{key, m.Data.(float64), m.Timestamp.Format(time.RFC3339)}
    // log.Debug(s)
    jsonStr, err := json.Marshal(s)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
  }

  return nil
}
