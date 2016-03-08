# Postback Delivery
Requirements:

- Go
- Redis

Installation:

1. Clone
2. Move src/ingest.php to some endpoint on the server (i.e. /var/www/html/ingest.php)
3. Either use `go install && bin/deliveryagent path/to/file.log` if on path, or `cd src/github.com/rickbau5/deliveryagent && go run main.go path/to/file.log`

To set it all up to run automatically, follow your favorite method for your box. 

Documentation can be found in the [doc](doc) folder, including [versions](doc/Info.md#Versions) used for all relevant modules and the [development log](doc/DevelopmentLog.md)

## Specification
Build a webapp to function as a small scale simulation of synchronizing data with third-party partners.

### Instructions
1. Provision provided linux server (see Resources - Server) with software stack required to complete project.
2. Build a php application to ingest http requests, and a go application to deliver http responses. Use Redis to host a job queue between them.
3. Reach out to (Resources - Contact) once your project is ready to demo, or if you encounter a block during development. Pursue independent troubleshooting prior to escalating questions with contact resource.
4. Maintain development notes, provide support documentation, and commit your project (application code / stack config) to Github.

### Data Flow
1. Web request (see sample request) >
2. "Ingestion Agent" (php) >
3. "Delivery Queue" (redis)
4. "Delivery Agent" (go) >
5. Web response (see sample response)

### App Operation
#### [Ingestion Agent (PHP)](doc/IngestionAgent.md)
1. Accept incoming http request
2. Push a "postback" object to Redis for each "data" object contained in accepted request.

#### [Delivery Agent (Go)](doc/DeliveryAgent.md)
1. Continuously pull "postback" objects from Redis
2. Deliver each postback object to http endpoint:
  - Endpoint method: request.endpoint.method.
  - Endpoint url: request.endpoint.url, with {xxx} replaced with values from each request.endpoint.data.xxx element.
3. Log delivery time, response code, response time, and response body.

#### Sample Request 
See [Ingestion Agent](doc/IngestionAgent.md) doc for example and command to create the request (using cURL).

#### Sample Response
`http://sample_domain_endpoint.com/data?key=Phyllobates&value=Terribilis&foo=`

### Extra Merit
Most of these ignored, as didn't fit into the ["6-8 hour"](http://i.imgur.com/KU6K3d8.jpg) timeline. Would be worth implementing if I choose to pursue the project any further.
- [ ] Clean, descriptive VCS commit history.
- [ ] Clean, easy-to-follow support documentation for an engineer attempting to troubleshoot your system.
- [ ] All services should be configured to run automatically, and service should remain functional after system restarts.
- [ ] High availability infrastructure considerations.
- [ ] Data integrity considerations, including safe shutdown.
- [ ] Modular code design.
- [ ] Configurable default value for unmatched url {key}s.
- [ ] Performance of system under external load.
- [ ] Performance of system with single request in infinite loop.
- [ ] Minimal bandwidth utilization between ingestion and delivery servers.
- [ ] Configurable response delivery retry attempts. 
- [ ] Ingestion endpoint functional at /i in addition to /ingest.php.
- [ ] Data validation / error handling.
- [ ] Ability to deliver POST (as well as GET) responses.
- [ ] Service monitoring / application profiling.
- [ ] Delivery volume / success / failure visualizations.
- [ ] Internal benchmarking tool.
