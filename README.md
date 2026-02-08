# url-resolver
Stream URL resolver
Behaviour

Extracts the video_id from the URL path.

Validates and forwards the bearer token.

Calls the Identity Service to determine user roles.

Calls the Availability Service to validate the availability window.

Retrieves video metadata from an in-memory map.

Returns a JSON playback response.

Error handling:

401 Unauthorized – missing/invalid Authorization header

403 Forbidden – video outside availability window

404 Not Found – unknown video ID

502 Bad Gateway – downstream service failure

Potential Improvements

With more time, the following enhancements could be considered:

Introduce integration-level tests for the handler using real HTTP servers (httptest.Server) to exercise the full request lifecycle.

Execute Identity and Availability calls concurrently in the handler to reduce overall latency, product decision of whether to default to standard level on identity service failure or not.

Replace the in-memory catalog with a dedicated metadata service or persistent store.

Externalise configuration (e.g. base URLs, playback base URL) via environment variables.

Add structured logging and request tracing for better observability.