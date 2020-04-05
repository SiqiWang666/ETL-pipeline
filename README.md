
Team Members:
Mingyu Ma & Siqi Wang & Weizhao Li

Contribution:
Weizhao Li: MS - Data - Cleaning
Siqi Wang: MS-Visitor-Counts
Mingyu Ma: MS-Website-Counts/MS-Browser-Counts



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