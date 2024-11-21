package model

type Message struct {
	MessageId         string
	Body              string
	ReceiptHandle     string
	Attributes        map[string]string
	MessageAttributes map[string]MessageAttribute
}

type MessageAttribute struct {
	DataType    string
	StringValue *string
	BinaryValue []byte
}
