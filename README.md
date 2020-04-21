This is a tool for processing JSON. You can learn more about it at https://ndjson-mcskyhcxcq-uc.a.run.app or read below:

## About
Prepend any url containing a JSON array with
```
https://ndjson-mcskyhcxcq-uc.a.run.app?url=
```
and it will come back as [newline delimited JSON](https://ndjson.org).

#### Why?
I work on a team with people who need to look at JSON from online sources, but not all can use the command line. 
With this we can paste urls to filtered data for each other, and also search by line.

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
