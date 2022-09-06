## About
Prepend any url containing a JSON array with
```
https://ndjson-mcskyhcxcq-uc.a.run.app?url=
```
and it will come back as [newline delimited JSON](https://ndjson.org).

### Example
[https://bit.ly/2PMxvrj](http://bit.ly/2PMxvrj) contains some JSON records.

[
https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj) are the same records streamed out as newline delimited JSON.

### Optional query parameters
There are some additional query params for filtering, selecting, and sampling records. After url you can pass one of:

|Param|Description|Example|
|-----|-----------|-------|
|**&head=**| Integer returning the **first n records** from a JSON array | [https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&head=3](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&head=3)|
|**&sample=**| Integer returning a **n% sample** of records from a JSON array | [https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&sample=20](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&sample=20) |
|**&nth=**| Positive integer returning the **nth record** from a JSON array |  [https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&nth=3](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&nth=3) |
|**&filter=**| String returning any record **containing the substring s** | [https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&filter=keats](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&filter=keats) |
|**&jq=**| String containing a **[jq](https://stedolan.github.io/jq/)** style filter per record | [https://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&jq=.name](http://ndjson-mcskyhcxcq-uc.a.run.app?url=http://bit.ly/2PMxvrj&jq=.name) |
