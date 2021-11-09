# MySQL SQS

A small MySQL UDF library for making AWS SQS calls from a MySQL database. Currently it only sends messages, but more functions could be implemented if needed.

## Usage

### `sqs_send_message`

Sends a message to a queue.

```sql
`sqs_send_message` ( @QueueURL , @MessageBody , @DelaySeconds , @MessageAttributes , @MessageSystemAttributes , @MessageDeduplicationID , @MessageGroupID )
```

#### `@QueueURL` (string)

The URL of the Amazon SQS queue to which a message is sent.

Queue URLs and names are case-sensitive.

#### `@MessageBody` (string)

The message to send. The minimum size is one character. The maximum size is 256 KB.
> ### Warning:
> A message can include only XML, JSON, and unformatted text. The following Unicode characters are allowed:
>
> `#x9` | `#xA` | `#xD` | `#x20` to `#xD7FF` | `#xE000` to `#xFFFD` | `#x10000` to `#x10FFFF`
>
> Any characters not included in this list will be rejected. For more information, [see the W3C specification for characters](http://www.w3.org/TR/REC-xml/#charsets).

#### `@DelaySeconds` (integer)

The length of time, in seconds, for which to delay a specific message. Valid values: 0 to 900. Maximum: 15 minutes. Messages with a positive `DelaySeconds` value become available for processing after the delay period is finished. If you don’t specify a value, the default value for the queue applies.
> ### Note:
> When you set `FifoQueue`, you can’t set `DelaySeconds` per message. You can set this parameter only on a queue level.

#### `@MessageAttributes` (json)

Each message attribute consists of a `Name` , `Type` , and `Value` . For more information, see [Amazon SQS message attributes](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-metadata.html#sqs-message-attributes) in the *Amazon SQS Developer Guide* .

> Name -> (string)
>
> Value -> (structure)
>
> The user-specified message attribute value. For string data types, the `Value` attribute has the same restrictions on the content as the message body. For more information, see `SendMessage`.
>
> > `Name`, `type`, `value` and the message body must not be empty or null. All parts of the message attribute, including `Name`, `Type`, and `Value`, are part of the message size restriction (256 KB or 262,144 bytes).
>
> StringValue -> (string)
>
> > Strings are Unicode with UTF-8 binary encoding. For a list of code values, see [ASCII Printable Characters](http://en.wikipedia.org/wiki/ASCII#ASCII_printable_characters) .
>
> BinaryValue -> (blob)
>
> > Binary type attributes can store any binary data, such as compressed data, encrypted data, or images.
>
> StringListValues -> (list)
>
> > Not implemented. Reserved for future use.
> >
> > (string)
>
> BinaryListValues -> (list)
>
> > Not implemented. Reserved for future use.
> >
> > (blob)
>
> DataType -> (string)
>
> > Amazon SQS supports the following logical data types: `String`, `Number`, and `Binary`. For the `Number` data type, you must use `StringValue`.
> >
> > You can also append custom labels. For more information, see [Amazon SQS Message Attributes](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-metadata.html#sqs-message-attributes) in the *Amazon SQS Developer Guide* .

### Shorthand Syntax:
Unfortunately the shorthand version is not supported by this MySQL function, and the json syntax must be used instead.

### JSON Syntax:

```json
{"string": {
     "StringValue": "string",
     "BinaryValue": blob,
     "StringListValues": ["string", ...],
     "BinaryListValues": [blob, ...],
     "DataType": "string"
   }
 ...}
```

#### `@MessageSystemAttributes` (json)

The message system attribute to send. Each message system attribute consists of a `Name` , `Type` , and `Value` .

> ### Warning:
>
> -   Currently, the only supported message system attribute is `AWSTraceHeader` . Its type must be `String` and its value must be a correctly formatted X-Ray trace header string.
>
>
> -   The size of a message system attribute doesn't count towards the total size of a message.

> Name -> (string)
>
> Value -> (structure)
>
> The user-specified message system attribute value. For string data types, the `Value` attribute has the same restrictions on the content as the message body. For more information, see `SendMessage`.
>
> > `Name`, `type`, `value` and the message body must not be empty or null.
>
> StringValue -> (string)
>
> > Strings are Unicode with UTF-8 binary encoding. For a list of code values, see [ASCII Printable Characters](http://en.wikipedia.org/wiki/ASCII#ASCII_printable_characters) .
>
> BinaryValue -> (blob)
>
> > Binary type attributes can store any binary data, such as compressed data, encrypted data, or images.
>
> StringListValues -> (list)
>
> > Not implemented. Reserved for future use.
> >
> > (string)
>
> BinaryListValues -> (list)
>
> > Not implemented. Reserved for future use.
> >
> > (blob)
>
> DataType -> (string)
>
> > Amazon SQS supports the following logical data types: `String` , `Number` , and `Binary` . For the `Number` data type, you must use `StringValue` .
> >
> > You can also append custom labels. For more information, see [Amazon SQS Message Attributes](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-message-metadata.html#sqs-message-attributes) in the *Amazon SQS Developer Guide* .

### Shorthand Syntax:
Unfortunately the shorthand version is not supported by this MySQL function, and the json syntax must be used instead.

### JSON Syntax:

```json
{"AWSTraceHeader": {
     "StringValue": "string",
     "BinaryValue": blob,
     "StringListValues": ["string", ...],
     "BinaryListValues": [blob, ...],
     "DataType": "string"
   }
 ...}
```

#### `@MessageDeduplicationID` (string)

This parameter applies only to FIFO (first-in-first-out) queues.

The token used for deduplication of sent messages. If a message with a particular `MessageDeduplicationId` is sent successfully, any messages sent with the same `MessageDeduplicationId` are accepted successfully but aren't delivered during the 5-minute deduplication interval. For more information, see [Exactly-once processing](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/FIFO-queues-exactly-once-processing.html) in the *Amazon SQS Developer Guide*.

- Every message must have a unique `MessageDeduplicationId` ,

  - You may provide a `MessageDeduplicationId` explicitly.

  - If you aren't able to provide a `MessageDeduplicationId` and you enable `ContentBasedDeduplication` for your queue, Amazon SQS uses a SHA-256 hash to generate the `MessageDeduplicationId` using the body of the message (but not the attributes of the message).

  - If you don't provide a `MessageDeduplicationId` and the queue doesn't have `ContentBasedDeduplication` set, the action fails with an error.

- If the queue has `ContentBasedDeduplication` set, your `MessageDeduplicationId` overrides the generated one.

- When `ContentBasedDeduplication` is in effect, messages with identical content sent within the deduplication interval are treated as duplicates and only one copy of the message is delivered.

- If you send one message with `ContentBasedDeduplication` enabled and then another message with a `MessageDeduplicationId` that is the same as the one generated for the first `MessageDeduplicationId` , the two messages are treated as duplicates and only one copy of the message is delivered.

> ### Note:
>
> The `MessageDeduplicationId` is available to the consumer of the message (this can be useful for troubleshooting delivery issues).
>
> If a message is sent successfully but the acknowledgement is lost and the message is resent with the same `MessageDeduplicationId` after the deduplication interval, Amazon SQS can't detect duplicate messages.
>
> Amazon SQS continues to keep track of the message deduplication ID even after the message is received and deleted.
>
> The maximum length of `MessageDeduplicationId` is 128 characters. `MessageDeduplicationId` can contain alphanumeric characters (`a-z` , `A-Z` , `0-9` ) and punctuation (`!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~` ).
>
> For best practices of using `MessageDeduplicationId` , see [Using the MessageDeduplicationId Property](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-messagededuplicationid-property.html) in the *Amazon SQS Developer Guide* .

#### `@MessageGroupID` (string)

This parameter applies only to FIFO (first-in-first-out) queues.

The tag that specifies that a message belongs to a specific message group. Messages that belong to the same message group are processed in a FIFO manner (however, messages in different message groups might be processed out of order). To interleave multiple ordered streams within a single queue, use `MessageGroupId` values (for example, session data for multiple users). In this scenario, multiple consumers can process the queue, but the session data of each user is processed in a FIFO fashion.

- You must associate a non-empty `MessageGroupId` with a message. If you don't provide a `MessageGroupId`, the action fails.

- `ReceiveMessage` might return messages with multiple `MessageGroupId` values. For each `MessageGroupId`, the messages are sorted by time sent. The caller can't specify a `MessageGroupId`.

The length of `MessageGroupId` is 128 characters. Valid values: alphanumeric characters and punctuation `(!"#$%&'()*+,-./:;<=>?@[\]^_`{|}~)` .

For best practices of using `MessageGroupId` , see [Using the MessageGroupId Property](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/using-messagegroupid-property.html) in the *Amazon SQS Developer Guide* .

> ### Warning:
>
>`MessageGroupId` is required for FIFO queues. You can't use it for Standard queues.
>

## Example:

You can use `null` for values you don't want to include. This example sends a message with the specified message body, delay period, and message attributes, to the specified queue.

The odd syntax `into @_` is so we can run this in a trigger without the "Not allowed to return a result set from a trigger" error.

```sql
select `sqs_send_message`(
	'https://sqs.us-east-1.amazonaws.com/80398EXAMPLE/MyQueue',  -- Queue URL
    'Information about the largest city in Any Region.', -- Message Body
    10,   -- Delay Seconds
    '{"City": {"DataType": "String","StringValue": "Any City"}}', -- Message Attributes
	null, -- Message System Attributes
    null, -- Message Deduplication ID
    null  -- Message Group ID
) into @_;
```

## Output:

The output will be null if there was an error sending the message. If you want to see the actual error messages, they will appear in the MySQL error log file like so

```json
{
  "MD5OfMessageBody": "51b0a325...39163aa0",
  "MD5OfMessageAttributes": "00484c68...59e48f06",
  "MessageId": "da68f62c-0c07-4bee-bf5f-7e856EXAMPLE"
}
```

#### `MD5OfMessageBody` (string)

An MD5 digest of the non-URL-encoded message body string. You can use this attribute to verify that Amazon SQS received the message correctly. Amazon SQS URL-decodes the message before creating the MD5 digest. For information about MD5, see [RFC1321](https://www.ietf.org/rfc/rfc1321.txt).

#### `MD5OfMessageAttributes` (string)

An MD5 digest of the non-URL-encoded message attribute string. You can use this attribute to verify that Amazon SQS received the message correctly. Amazon SQS URL-decodes the message before creating the MD5 digest. For information about MD5, see [RFC1321](https://www.ietf.org/rfc/rfc1321.txt).

#### `MD5OfMessageSystemAttributes` (string)

An MD5 digest of the non-URL-encoded message system attribute string. You can use this attribute to verify that Amazon SQS received the message correctly. Amazon SQS URL-decodes the message before creating the MD5 digest.

#### `MessageId` (string)

An attribute containing the `MessageId` of the message sent to the queue. For more information, see [Queue and Message Identifiers](https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-queue-message-identifiers.html) in the *Amazon SQS Developer Guide*.

#### `SequenceNumber` (string)

This parameter applies only to FIFO (first-in-first-out) queues.

The large, non-consecutive number that Amazon SQS assigns to each message.

The length of `SequenceNumber` is 128 bits. `SequenceNumber` continues to increase for a particular `MessageGroupId` .

## Dependencies

You will need Golang, which you can get from here https://golang.org/doc/install.

You will also need to install the MySQL dev library:

### Debian / Ubuntu

```shell
sudo apt update
sudo apt install libmysqlclient-dev
```

## Installing

Know your MySQL plugin directory, which can be found by running this MySQL query:

```sql
select @@plugin_dir;
```

then replace `/usr/lib/mysql/plugin` below with your MySQL plugin directory.

```shell
cd ~ # or wherever you store your git projects
git clone https://github.com/StirlingMarketingGroup/mysql-sqs.git
cd mysql-sqs
go build -buildmode=c-shared -o mysql_sqs.so
sudo cp mysql_sqs.so /usr/lib/mysql/plugin/mysql_sqs.so # replace plugin dir here if needed
```

Enable the function in MySQL by running this MySQL query

```sql
create function`sqs_send_message`returns string soname'mysql_sqs.so';
```