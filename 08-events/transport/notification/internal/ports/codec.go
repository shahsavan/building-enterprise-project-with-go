package ports

import (
	"bytes"

	"github.com/hamba/avro/v2"
)

var assignmentSchema = avro.MustParse(`{ /* same schema as in /schemas/assignment.avsc */ }`)

func EncodeAssignment(a AssignmentCreated) ([]byte, error) {
	var buf bytes.Buffer
	err := avro.NewEncoderForSchema(assignmentSchema, &buf).Encode(a)
	return buf.Bytes(), err
}

func DecodeAssignment(b []byte) (AssignmentCreated, error) {
	var out AssignmentCreated
	err := avro.NewDecoderForSchema(assignmentSchema, bytes.NewReader(b)).Decode(&out)
	return out, err
}
