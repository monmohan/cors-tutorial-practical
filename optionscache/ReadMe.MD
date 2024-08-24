# Experiments with Caching Preflight Results in Browsers

This repository contains experiments to demonstrate how browsers cache CORS (Cross-Origin Resource Sharing) preflight results. The behavior is specified in the [Fetch Standard](https://fetch.spec.whatwg.org/#cors-preflight-cache).

## Key Points from the Specification

1. Preflight results can be cached by browsers.
2. The cache key should include:
   - Origin
   - URL
   - Credentials mode
   - Header names
   - HTTP method
3. The `max-age` value from the `Access-Control-Max-Age` header should be respected.
4. The exact caching mechanism and duration are left to the browser's discretion.

## Example Preflight Request

```http
OPTIONS /api/data?type=first HTTP/1.1
Accept: */*
Accept-Encoding: gzip, deflate, br, zstd
Accept-Language: en-GB,en-US;q=0.9,en;q=0.8
Access-Control-Request-Headers: content-type
Access-Control-Request-Method: GET
Connection: keep-alive
Host: ....
Origin: ....
Referer: .....
Sec-Fetch-Dest: empty
Sec-Fetch-Mode: cors
Sec-Fetch-Site: same-site
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/128.0.0.0 Safari/537.36
```

## Example Server Response

```http
HTTP/1.1 200 OK
Access-Control-Allow-Headers: Content-Type, X-Custom-Header-1, X-Custom-Header-2
Access-Control-Allow-Methods: GET, POST, OPTIONS
Access-Control-Allow-Origin: *
Access-Control-Max-Age: 3600 
Date: Sat, 24 Aug 2024 03:49:02 GMT
Content-Length: 19
Content-Type: text/plain; charset=utf-8
```

## Expected Behavior

1. I expected that in subsequent call to server in the next one hour from same browser (`Access-Control-Max-Age: 3600`), no OPTIONS call will be made


2. However, changes in the following will ignore the cache and trigger a new OPTIONS call:
   - URL query parameters (e.g., `/api/data?type=second`)
   - Use of different custom headers

3. Browser-specific behavior:
   - Chrome makes one OPTIONS call for any combination of the allowed custom headers, differing from the specification.

