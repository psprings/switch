# switch

Random lunchtime experiment that I might end up using

## Purpose

Proxy a webhook endpoint e.g. Jenkins, Drone.io, etc
This allows you to take down a CI server and still receive webhooks (in a serverless model
or front many CI servers with one hook proxy) or even test webhooks while doing local development.

## How it works

Webhook hits `switch` (in server mode) which checks to see if the final destination hook is accessible
if it is accessible, the request is forwarded directly. If it is not available, the request
body and metadata are stored in a queue (such as SQS) and there would be an optional
capability to wake up a sleeping server.

For local webhook functionality, you can run `switch` in daemon mode and it will poll a given queue.
When it finds messages, it will attempt to `POST` with the URL (in this case `localhost`) and body content from the queue.