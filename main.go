package main

// #include <string.h>
// #include <stdbool.h>
// #include <mysql.h>
// #cgo CFLAGS: -O3 -I/usr/include/mysql -fno-omit-frame-pointer
import "C"
import (
	"encoding/json"
	"log"
	"os"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
)

func msg(message *C.char, s string) {
	m := C.CString(s)
	defer C.free(unsafe.Pointer(m))

	C.strcpy(message, m)
}

var l = log.New(os.Stderr, "sqs-send-message: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Llongfile)

//export sqs_send_message_init
func sqs_send_message_init(initid *C.UDF_INIT, args *C.UDF_ARGS, message *C.char) C.bool {
	if args.arg_count != 7 {
		msg(message, "`sqs_send_message` requires 7 parameters")
		return C.bool(true)
	}

	argsTypes := (*[7]uint32)(unsafe.Pointer(args.arg_type))

	argsTypes[0] = C.STRING_RESULT
	argsTypes[1] = C.STRING_RESULT
	argsTypes[2] = C.INT_RESULT
	argsTypes[3] = C.STRING_RESULT
	argsTypes[4] = C.STRING_RESULT
	argsTypes[5] = C.STRING_RESULT
	argsTypes[6] = C.STRING_RESULT

	initid.maybe_null = C.bool(true)

	return C.bool(false)
}

//export sqs_send_message
func sqs_send_message(initid *C.UDF_INIT, args *C.UDF_ARGS, result *C.char, length *uint64, isNull *C.char, message *C.char) *C.char {
	c := 7
	argsArgs := (*[1 << 30]*C.char)(unsafe.Pointer(args.args))[:c:c]
	argsLengths := (*[1 << 30]uint64)(unsafe.Pointer(args.lengths))[:c:c]

	intArgs := (*[1 << 30]*uint64)(unsafe.Pointer(args.args))

	var queueURL *string
	if argsArgs[0] != nil {
		queueURL = aws.String(C.GoStringN(argsArgs[0], C.int(argsLengths[0])))
	}

	var messageBody *string
	if argsArgs[1] != nil {
		messageBody = aws.String(C.GoStringN(argsArgs[1], C.int(argsLengths[1])))
	}

	var delaySeconds *int64
	if argsArgs[2] != nil {
		delaySeconds = aws.Int64(int64(*intArgs[2]))
	}

	var messageAttributes map[string]*sqs.MessageAttributeValue
	if argsArgs[3] != nil {
		err := json.Unmarshal([]byte(C.GoStringN(argsArgs[3], C.int(argsLengths[3]))), &messageAttributes)
		if err != nil {
			l.Println(errors.Wrapf(err, "failed to unmarshal message attributes").Error())

			*length = 0
			*isNull = 1
			return nil
		}
	}

	var messageSystemAttributes map[string]*sqs.MessageSystemAttributeValue
	if argsArgs[4] != nil {
		err := json.Unmarshal([]byte(C.GoStringN(argsArgs[4], C.int(argsLengths[4]))), &messageSystemAttributes)
		if err != nil {
			l.Println(errors.Wrapf(err, "failed to unmarshal message attributes").Error())

			*length = 0
			*isNull = 1
			return nil
		}
	}

	var messageDeduplicationID *string
	if argsArgs[5] != nil {
		messageDeduplicationID = aws.String(C.GoStringN(argsArgs[5], C.int(argsLengths[5])))
	}

	var messageGroupID *string
	if argsArgs[6] != nil {
		messageGroupID = aws.String(C.GoStringN(argsArgs[6], C.int(argsLengths[6])))
	}

	sess, err := session.NewSession()
	if err != nil {
		l.Println(errors.Wrapf(err, "failed to create aws session").Error())

		*length = 0
		*isNull = 1
		return nil
	}

	out, err := sqs.New(sess).SendMessage(&sqs.SendMessageInput{
		MessageBody:             messageBody,
		QueueUrl:                queueURL,
		DelaySeconds:            delaySeconds,
		MessageAttributes:       messageAttributes,
		MessageSystemAttributes: messageSystemAttributes,
		MessageDeduplicationId:  messageDeduplicationID,
		MessageGroupId:          messageGroupID,
	})
	if err != nil {
		l.Println(errors.Wrapf(err, "failed send sqs message").Error())

		*length = 0
		*isNull = 1
		return nil
	}

	j, _ := json.Marshal(out)

	*length = uint64(len(j))
	*isNull = 0
	return C.CString(string(j))
}

func main() {}
