package elasticsearch

import (
  "bytes"
  "net/http"
  "encoding/json"
  "strings"
  "time"
  "fmt"
  "io/ioutil"

  log "github.com/Sirupsen/logrus"
  "github.com/intelsdi-x/snap-plugin-lib-go/v1/plugin"
)

type EPublisher struct {
}

//NewElasticsearchPublisher returns an instance of the Elasticsearch publisher
func NewElasticsearchPublisher() *EPublisher {
        return &EPublisher{}
}

func (f EPublisher) GetConfigPolicy() (plugin.ConfigPolicy, error) {
  policy := plugin.NewConfigPolicy()

  policy.AddNewStringRule([]string{""}, "uri", false, plugin.SetDefaultString("http://localhost:9200/snap/default"))
  return *policy, nil
}

type Index struct {
  Index string `json:"_index"`
  Type  string `json:"_type"`
}

type ESBulkMetadata struct {
  Index Index `json:"index"`
}

type MetricToPublish struct {
  // The timestamp from when the metric was created.
  Key string
  Value float64
  Timestamp string
}

// Publish Elasticsearch publish function
func (f EPublisher) Publish(mts []plugin.Metric, cfg plugin.Config) error {
  var sampleBuffer bytes.Buffer
  var body []byte
  url, err := cfg.GetString("uri")
        if err != nil {
    log.Errorf("unable to parse elasticsearch uri from configs")
                return err
  }
  i := ESBulkMetadata{
     Index{
        Index: "snap",
        Type:    "default",
    },
  }
  indexTargetStr, err := json.Marshal(i)
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
    MetricJSONStr, err := json.Marshal(s)
    if err != nil {
      log.Errorf("Unable to marshal metric. Error: %s", err)
      return fmt.Errorf("Unable to marshal metric. Error: %s", err)
    }
    sampleBuffer.WriteString(string(indexTargetStr))
    sampleBuffer.WriteString("\n")
    //log.Errorf("Build! %i", len(sampleBuffer.String()))
    sampleBuffer.WriteString(string(MetricJSONStr))
    sampleBuffer.WriteString("\n")
    if len(sampleBuffer.String()) > 2000 {
      log.Errorf("%s", sampleBuffer.String())
      req, err := http.NewRequest("POST", url, &sampleBuffer)
      req.Header.Set("Content-Type", "application/json")
      timeout := time.Duration(30 * time.Second)
      client := http.Client{
        Timeout: timeout,
      }
      resp, err := client.Do(req)
      if err != nil {
        log.Errorf("%s", err)
      }
      defer resp.Body.Close()
      body, _ = ioutil.ReadAll(resp.Body)
      log.Errorf("%s", string(body))
      if err != nil {
        log.Errorf("Unable to send metrics. Error: %s", err)
        return fmt.Errorf("Unable to send metrics. Error: %s", err)
      }
    sampleBuffer.Reset()
    }
  }
  return nil
}
