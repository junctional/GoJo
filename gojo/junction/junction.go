package junction

import (
	"../patterns/binary"
	"../patterns/ternary"
	"../patterns/unary"
	"../types"
	"./controller"
	"errors"
)

type Junction struct {
	port       chan types.Packet
	JunctionId int
}

func NewJunction() *Junction {
	sender := make(chan types.Packet)

	controller.StartController(sender)

	return &Junction{sender, 1}
}

func NewAsyncSignal[T any](j *Junction) (types.SignalId, func(T)) {
	signalId := types.SignalId{
		ChannelType: types.AsyncSignal,
		Id:          getNewPortId(j),
		JunctionId:  (*j).JunctionId,
	}

	return signalId, func(data T) {
		(*j).port <- types.Packet{
			SignalId: signalId,
			Type:     types.MESSAGE,
			Payload: types.Payload{
				Msg: data,
			},
		}
	}
}

func NewSyncSignal[T any, R any](j *Junction) (types.SignalId, func(T) (R, error)) {
	signalId := types.SignalId{
		ChannelType: types.SyncSignal,
		Id:          getNewPortId(j),
		JunctionId:  (*j).JunctionId,
	}

	return signalId, func(data T) (R, error) {
		recvChannel := make(chan interface{})

		(*j).port <- types.Packet{
			SignalId: signalId,
			Type:     types.MESSAGE,
			Payload: types.Payload{
				Msg: data,
				Ch:  recvChannel,
			},
		}

		receivedData := <-recvChannel

		var returnData R

		switch t := receivedData.(type) {
		case R:
			returnData := t
			return returnData, nil
		default:
			return returnData, errors.New("invalid data type")
		}
	}
}

func getNewPortId(j *Junction) int {
	receiver := make(chan interface{})
	(*j).port <- types.Packet{Type: types.GetNewPortId, Payload: types.Payload{Ch: receiver}}
	signalId := <-receiver

	switch t := signalId.(type) {
	case int:
		return t
	}

	return 0
}

func NewUnaryAsyncJoinPattern[T any](j *Junction, signal types.SignalId) unary.AsyncPartialPattern[T] {
	return unary.AsyncPartialPattern[T]{
		Port:       (*j).port,
		Signals:    []types.SignalId{signal},
		JunctionId: (*j).JunctionId,
	}
}

func NewUnarySyncJoinPattern[T any](j *Junction, signal types.SignalId) unary.SyncPartialPattern[T] {
	return unary.SyncPartialPattern[T]{
		Port:       (*j).port,
		Signals:    []types.SignalId{signal},
		JunctionId: (*j).JunctionId,
	}
}

func NewBinaryAsyncJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.AsyncPartialPattern[T, R] {
	return binary.AsyncPartialPattern[T, R]{
		Port:       (*j).port,
		JunctionId: (*j).JunctionId,
		Signals:    []types.SignalId{signalOne, signalTwo},
	}
}

func NewBinarySyncJoinPattern[T any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId) binary.SyncPartialPattern[T, R] {
	return binary.SyncPartialPattern[T, R]{
		Port:       (*j).port,
		JunctionId: (*j).JunctionId,
		Signals:    []types.SignalId{signalOne, signalTwo},
	}
}

func NewTernaryAsyncJoinPattern[T any, S any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId, signalThree types.SignalId) ternary.AsyncPartialPattern[T, S, R] {
	return ternary.AsyncPartialPattern[T, S, R]{
		Port:       (*j).port,
		JunctionId: (*j).JunctionId,
		Signals:    []types.SignalId{signalOne, signalTwo, signalThree},
	}
}

func NewTernarySyncJoinPattern[T any, S any, R any](j *Junction, signalOne types.SignalId, signalTwo types.SignalId, signalThree types.SignalId) ternary.SyncPartialPattern[T, S, R] {
	return ternary.SyncPartialPattern[T, S, R]{
		Port:       (*j).port,
		JunctionId: (*j).JunctionId,
		Signals:    []types.SignalId{signalOne, signalTwo, signalThree},
	}
}
