# Delivery Agent (Go)
Retrives postback objects from Redis and sends the outgoing request. Logs certain aspects of the delivery, as specificed below.

__Specification__

1. Continuously pull "postback" objects from Redis
2. Deliver each postback object to http endpoint:
  - Endpoint method: request.endpoint.method.
  - Endpoint url: request.endpoint.url, with {xxx} replaced with values from each request.endpoint.data.xxx element.
3. Log delivery time, response code, response time, and response body.

