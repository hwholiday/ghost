package dove

import (
	"errors"
	"github.com/google/uuid"
	api "github.com/hwholiday/ghost/dove/api/dove"
	"google.golang.org/protobuf/proto"
	"time"
)

type build struct {
	dove         *api.Dove
	doveMetadata *api.DoveMetadata
	doveBody     *api.DoveBody
}

func NewDoveRes() *build {
	return &build{}
}

func (b *build) Result() ([]byte, error) {
	if b.doveMetadata == nil {
		return nil, errors.New("metadata is empty")
	}
	if b.doveBody == nil {
		return nil, errors.New("body is empty")
	}
	res := api.Dove{
		Metadata: b.doveMetadata,
		Body:     b.doveBody,
	}
	return proto.Marshal(&res)
}

func (b *build) MetadataSeq(seq string) *build {
	if b.doveMetadata != nil {
		b.doveMetadata.Seq = seq
	} else {
		b.doveMetadata = &api.DoveMetadata{
			Seq:       seq,
			Timestamp: time.Now().UnixMilli(),
		}
	}
	return b
}

func (b *build) Metadata(crcId, ackId uint64) *build {
	b.doveMetadata = &api.DoveMetadata{
		CrcId:     crcId,
		AckId:     ackId,
		Seq:       uuid.NewString(),
		Timestamp: time.Now().UnixMilli(),
	}
	return b
}

func (b *build) BodyErr(code uint64, msg string) *build {
	b.doveBody = &api.DoveBody{
		Msg:  msg,
		Code: code,
	}
	return b
}

func (b *build) BodyOk(code ...uint64) *build {
	var c = DefaultDoveBodyCodeOK
	if len(code) > 1 {
		c = code[0]
	}
	b.doveBody = &api.DoveBody{
		Code: c,
	}
	return b
}
func (b *build) BodyOkWithData(data []byte, code ...uint64) *build {
	var c = DefaultDoveBodyCodeOK
	if len(code) > 1 {
		c = code[0]
	}
	b.doveBody = &api.DoveBody{
		Code: c,
		Data: data,
	}
	return b
}
