package roblox

type creatorDto struct {
	UserId  string `json:"userId,omitempty"`
	GroupId string `json:"groupId,omitempty"`
}

type creationContextDto struct {
	Creator creatorDto `json:"creator"`
}

type moderationResultDto struct {
	ModerationState string `json:"moderationState"`
}

type assetDto struct {
	Path               string              `json:"path"`
	RevisionId         string              `json:"revisionId"`
	RevisionCreateTime string              `json:"revisionCreateTime"`
	Id                 string              `json:"assetId"`
	DisplayName        string              `json:"displayName"`
	Description        string              `json:"description"`
	AssetType          string              `json:"assetType"`
	CreationContext    creationContextDto  `json:"creationContext"`
	ModerationResult   moderationResultDto `json:"moderationResult"`
	State              string              `json:"state"`
}

type assetDeliveryMetaDto struct {
	Location string `json:"location"`
}

type assetUploadMetaDto struct {
	AssetType       string              `json:"assetType"`
	DisplayName     string              `json:"displayName"`
	Description     string              `json:"description"`
	AssetId         string              `json:"assetId,omitempty"`
	CreationContext *creationContextDto `json:"creationContext,omitempty"`
}

type assetVersionDto struct {
	Path             string               `json:"path"`
	CreationTime     string               `json:"createTime"`
	CreationContext  creationContextDto   `json:"creationContext"`
	ModerationResult *moderationResultDto `json:"moderationResult"`
	Published        bool                 `json:"published"`
}

type assetVersionPage struct {
	AssetVersions []assetVersionDto `json:"assetVersions"`
	NextPageToken string            `json:"nextPageToken"`
}

type assetPushResultDto struct {
	AssetId    string `json:"assetId"`
	RevisionId string `json:"revisionId"`
}
