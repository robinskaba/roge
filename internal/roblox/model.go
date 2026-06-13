package roblox

import "time"

type CreatorType string

const (
	CreatorTypeUser  CreatorType = "USER"
	CreatorTypeGroup CreatorType = "GROUP"
)

func (t CreatorType) IsValid() bool {
	return t == CreatorTypeUser || t == CreatorTypeGroup
}

type Creator struct {
	Type CreatorType
	Id   string
}

type Version struct {
	Id              int
	Time            time.Time
	ModerationState *string // TODO find out all possible values
}

type Asset struct {
	Id          string
	Version     Version
	Name        string
	Description string
	Creator     Creator
}
