package elasticsearch

import (
        "bytes"
        "net/http"
        "encoding/json"
        "strings"
        "time"

        log "github.com/Sirupsen/logrus"
        "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

type EPublisher struct {
}

func (f EPublisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
  policy := plugin.NewConfigPolicy()

  policy.AddNewStringRule([]string{""}, "uri", false, plugin.SetDefaultString("http://localhost:9200/snap/default"))
  return *policy, nil
}

type MetricToPublish struct {
	// The timestamp from when the metric was created.
	Key string
  Value float64
  Timestamp string
}

// Publish Elasticsearch publish function
func (f EPublisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {
  url, err := cfg.GetString("uri")
	if err != nil {
    log.Errorf("unable to parse elasticsearch uri from configs")
		return err
  }
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
