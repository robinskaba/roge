package roblox

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (dto creationContextDto) toCreator() (Creator, error) {
	var creator Creator

	// must have a creator
	if dto.Creator.UserId == "" && dto.Creator.GroupId == "" {
		return creator, fmt.Errorf("no creator present in creation context")
	}

	// if has user id set as user
	if dto.Creator.UserId != "" {
		creator.Type = CreatorTypeUser
		creator.Id = dto.Creator.UserId
	} else { // is group
		creator.Type = CreatorTypeGroup
		creator.Id = dto.Creator.GroupId
	}
	return creator, nil
}

func (dto assetVersionDto) toVersion() (Version, error) {
	var version Version

	// extract version number
	var err error
	parts := strings.Split(dto.Path, "/")
	version.Id, err = strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return version, err
	}

	// extract creation time
	extrTime, err := time.Parse(time.RFC3339Nano, dto.CreationTime)
	if err != nil {
		return version, err
	}
	version.Time = extrTime

	// extract moderation state
	if dto.ModerationResult != nil && dto.ModerationResult.ModerationState != "" {
		version.ModerationState = &(dto.ModerationResult.ModerationState)
	}

	return version, nil
}

func (dto assetDto) toAsset() (Asset, error) {
	var asset Asset

	// map basics
	asset.Id = dto.Id
	asset.Name = dto.DisplayName
	asset.Description = dto.Description

	// map creator
	creator, err := dto.CreationContext.toCreator()
	if err != nil {
		return asset, err
	}
	asset.Creator = creator

	// map version
	var version Version

	id, err := strconv.Atoi(dto.RevisionId)
	if err != nil {
		return asset, err
	}
	version.Id = id

	time, err := time.Parse(time.RFC3339Nano, dto.RevisionCreateTime)
	if err != nil {
		return asset, err
	}
	version.Time = time
	version.ModerationState = &dto.ModerationResult.ModerationState
	asset.Version = version

	return asset, nil
}
