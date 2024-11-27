package main

type ContextString string

var ContextWasResponed ContextString

const (
	NotWritten           = "not_written"
	AlreadyResponseError = "already_responsed"
	Written              = "written"
)
