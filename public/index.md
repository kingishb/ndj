## About
Prepend any url containing a JSON array with
```
http://ndjson.sarriaking.com?url=
```
and it will come back as [newline delimitied JSON](http://ndjson.org). 

### Example
[http://bit.ly/2PMxvrj](http://bit.ly/2PMxvrj) contains some JSON records.

[http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj](http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj) are the same records streamed out as newline delimited JSON.

### Optional query parameters
There are some additional query params for filtering, selecting, and sampling records. After url you can pass one of:

|Param|Description|Example|
|-----|-----------|-------|
|**&head=**| Integer returning the **first n records** from a JSON array | [http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&head=3](http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&head=3)|
|**&sample=**| Integer returning a **n% sample** of records from a JSON array | [http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&sample=20](http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&sample=20) |
|**&nth=**| Positive integer returning the **nth record** from a JSON array |  [http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&nth=3](http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&nth=3) |
|**&filter=**| String returning any record **containing the substring s** | [http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&filter=keats](http://ndjson.sarriaking.com?url=http://bit.ly/2PMxvrj&filter=keats) |
