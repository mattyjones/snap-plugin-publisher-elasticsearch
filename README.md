### Elasticsearch Publisher for Snap Telemetry

### Supported Platforms
x86_64

### Known Issues
- [ ] Split configuration into multiple options
- [ ] Add tests
- [ ] Update naming to properly reflect Elasticsearch

### Snap Version dependencies
      1.0

### Installation
At the moment you can use docker to build the Dockerfile.  You can "docker cp" with your container id to obtain the /go/bin/snap-plugin-publisher-elasticsearch binary.  Travis CI coming shortly.

### Usage
There is a single option, "uri" which should be set to an Elasticsearch server's http port(typically tcp/9200).

### Contributors
@rakah

### License
Apache-2.0
