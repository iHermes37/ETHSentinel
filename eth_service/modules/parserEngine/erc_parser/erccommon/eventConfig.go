package erccommon

type EventPaserConfig map[MethodName]EventParserFunc

type EventParserFunc func(data []byte) (TokenEvent, error)
