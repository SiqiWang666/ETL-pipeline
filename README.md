
A server log analysis built in microservices architecture.
```shell
.
├── README.md
├── config.yaml
├── logs
├── monolith
├── ms-api-gateway     # service, api gate way
├── ms-browser-counts  # service, browser-counts
├── ms-data-cleaning   # service, data-cleaning
├── ms-line-count      # service, line-count
├── ms-visitor-counts  # service, visitor-counts
└── ms-website-counts  # service, website-counts
```

Message format
```
{
  "statusCode": 201,
  "message": "Success",
  "data": {
    "log_a.txt": 100
  }
}
```



Error format
```
{
  "statusCode": 404,
  "error": "File log_b.txt not found!"
}
```