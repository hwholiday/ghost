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
	err          error
}

func NewDoveRes() *build {
	return &build{
		dove:         nil,
		doveMetadata: nil,
		doveBody:     nil,
		err:          nil,
	}
}

func (b *build) Result() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	if b.doveMetadata.GetAckId() == 0 {
		return nil, errors.New("metadata is empty")
	}
	if b.doveBody.GetCode() == 0 {
		return nil, errors.New("body code is  ")
	}
	res := api.Dove{
		Metadata: b.doveMetadata,
		Body:     b.doveBody,
	}

	return proto.Marshal(&res)
}

func (b *build) MetadataSeq(seq string) *build {
	if b.err != nil {
		return b
	}
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
	if b.err != nil {
		return b
	}
	if b.doveMetadata != nil {
		b.doveMetadata.CrcId = crcId
		b.doveMetadata.AckId = ackId
	} else {
		b.doveMetadata = &api.DoveMetadata{
			CrcId:     crcId,
			AckId:     ackId,
			Seq:       uuid.NewString(),
			Timestamp: time.Now().UnixMilli(),
		}
	}

	return b
}

func (b *build) BodyMsg(msg string) *build {
	if b.err != nil {
		return b
	}
	if b.doveBody != nil {
		b.doveBody.Msg = msg
	} else {
		b.doveBody = &api.DoveBody{
			Msg: msg,
		}
	}
	return b
}

func (b *build) BodyOk(code ...uint64) *build {
	if b.err != nil {
		return b
	}
	var c = DefaultDoveBodyCodeOK
	if len(code) > 1 {
		c = code[0]
	}
	return b.BodyCode(c)
}

func (b *build) BodyCode(code uint64) *build {
	if b.err != nil {
		return b
	}
	if b.doveBody != nil {
		b.doveBody.Code = code
	} else {
		b.doveBody = &api.DoveBody{
			Code: code,
		}
	}
	return b
}
func (b *build) BodyData(data []byte) *build {
	if b.err != nil {
		return b
	}
	if b.doveBody != nil {
		b.doveBody.Data = data
	} else {
		b.doveBody = &api.DoveBody{
			Data: data,
		}
	}
	return b
}
func (b *build) BodyPbData(data proto.Message) *build {
	if b.err != nil {
		return b
	}
	var byt []byte
	byt, b.err = proto.Marshal(data)
	if b.err != nil {
		return b
	}
	return b.BodyData(byt)
}
func (b *build) BodyExpand(expand []byte) *build {
	if b.err != nil {
		return b
	}
	if b.doveBody != nil {
		b.doveBody.Expand = expand
	} else {
		b.doveBody = &api.DoveBody{
			Expand: expand,
		}
	}
	return b
}
