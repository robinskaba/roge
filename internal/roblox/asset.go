package roblox

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"strings"

	"github.com/robinskaba/roge/internal/conversion"
	"github.com/robloxapi/rbxfile"
)

func fetchAssetMetadata(apiKey string, assetId string, version int) (assetDeliveryMetaDto, error) {
	var metadata assetDeliveryMetaDto

	metaURL := fmt.Sprintf("https://apis.roblox.com/asset-delivery-api/v1/assetId/%s", assetId)
	if version > 0 {
		metaURL += fmt.Sprintf("/version/%d", version)
	}

	body, err := Fetch(metaURL, apiKey)
	if err != nil {
		return metadata, err
	}

	if err := json.Unmarshal(body, &metadata); err != nil {
		return metadata, fmt.Errorf("parsing metadata: %w", err)
	}

	return metadata, nil
}

func buildMultipart(metaJSON []byte, fileData []byte) (*bytes.Buffer, *multipart.Writer) {
	body := &bytes.Buffer{}
	w := multipart.NewWriter(body)
	{
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="request"`)
		h.Set("Content-Type", "application/json")
		p, _ := w.CreatePart(h)
		p.Write(metaJSON)
	}
	{
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="fileContent"; filename="model.rbxm"`)
		h.Set("Content-Type", "model/x-rbxm")
		p, _ := w.CreatePart(h)
		p.Write(fileData)
	}
	w.Close()
	return body, w
}

func formatAssetVersions(assetVersions []assetVersionDto) ([]Version, error) {
	result := []Version{}

	for _, v := range assetVersions {
		pkgVersion, err := v.toVersion()
		if err != nil {
			fmt.Println(fmt.Errorf("package version has an invalid format: %w", err))
			continue
		}

		result = append(result, pkgVersion)
	}

	return result, nil
}

func GetAsset(apiKey string, assetId string) (Asset, error) {
	url := fmt.Sprintf("https://apis.roblox.com/assets/v1/assets/%s", assetId)

	var assetDto assetDto

	body, err := Fetch(url, apiKey)
	if err != nil {
		return Asset{}, err
	}

	err = json.Unmarshal(body, &assetDto)
	if err != nil {
		return Asset{}, err
	}

	asset, err := assetDto.toAsset()
	if err != nil {
		return Asset{}, err
	}

	return asset, nil
}

func Pull(apiKey string, assetId string, version int) (*rbxfile.Instance, error) {
	metadata, err := fetchAssetMetadata(apiKey, assetId, version)
	if err != nil {
		return nil, err
	}
	if metadata.Location == "" {
		return nil, errors.New("no location in metadata")
	}

	binary, err := fetchBinaryFromCdn(metadata.Location)
	if err != nil {
		return nil, err
	}

	file, err := conversion.DecodeRbxm(binary)
	if err != nil {
		return nil, err
	}

	return file, nil
}

type PushConfig struct {
	// ALL
	ApiKey string
	Rbxm   []byte

	// PATCH only
	AssetId string

	// POST
	AuthorId   string
	AuthorType CreatorType

	// POST + PATCH with update
	Name        string
	Description string
}

func Push(cfg PushConfig) (string, int, error) {
	// setup creator if set
	var meta assetUploadMetaDto = assetUploadMetaDto{
		AssetType:   "Model",
		DisplayName: cfg.Name,
		Description: cfg.Description,
		AssetId:     cfg.AssetId, // if not set in cfg then basically not set in here thx to 'omitempty'
	}
	if cfg.AuthorId != "" {
		if !cfg.AuthorType.IsValid() {
			return "", 0, fmt.Errorf("invalid asset author")
		}
		if cfg.AssetId != "" {
			return "", 0, fmt.Errorf("asset ID and author ID cannot both be set at once")
		}

		var creatorDto creatorDto
		switch cfg.AuthorType {
		case CreatorTypeUser:
			creatorDto.UserId = cfg.AuthorId
		case CreatorTypeGroup:
			creatorDto.GroupId = cfg.AuthorId
		}
		meta.CreationContext = &creationContextDto{
			Creator: creatorDto,
		}
	}

	// setup request
	method := "POST"
	url := "https://apis.roblox.com/assets/v1/assets"

	if cfg.AssetId != "" {
		// set as PATCH if assetId set
		method = "PATCH"
		url += fmt.Sprintf("/%s", cfg.AssetId)

		// add updateMask flag with updated fields if set
		if cfg.Name != "" || cfg.Description != "" {
			updatedFields := []string{}
			if cfg.Name != "" {
				updatedFields = append(updatedFields, "displayName")
			}
			if cfg.Description != "" {
				updatedFields = append(updatedFields, "description")
			}
			url += "?updateMask=" + strings.Join(updatedFields, ",")
		}
	}

	// payload
	metaJson, err := json.Marshal(meta)
	if err != nil {
		return "", 0, err
	}
	body, w := buildMultipart(metaJson, cfg.Rbxm)
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("x-api-key", cfg.ApiKey)

	// send request
	res, err := doUpload(req)
	if err != nil {
		return "", 0, err
	}
	if !res.Done {
		// missing path?
		if res.Path == "" {
			raw, _ := json.MarshalIndent(res, "", "  ")
			return "", 0, fmt.Errorf("done=false with no operation path:\n%s", raw)
		}

		// poll until ready
		final, err := pollOperation(cfg.ApiKey, res.Path)
		if err != nil {
			return "", 0, err
		}
		res = final
	}

	// response data
	var result assetPushResultDto
	if err = json.Unmarshal(res.Result, &result); err != nil {
		return "", 0, fmt.Errorf("failed to parse push response: %w", err)
	}
	version, err := strconv.Atoi(result.RevisionId)
	if err != nil {
		version = 1
	}

	return result.AssetId, version, nil
}

func GetVersions(apiKey string, assetId string) ([]Version, error) {
	versions := []Version{}
	pageToken := ""

	for {
		// build request
		url := fmt.Sprintf(
			"https://apis.roblox.com/assets/v1/assets/%s/versions?maxPageSize=50",
			assetId,
		)
		if pageToken != "" {
			url += "&pageToken=" + pageToken
		}

		var page assetVersionPage

		body, err := Fetch(url, apiKey)
		if err != nil {
			return nil, err
		}

		// unmarshal page
		if err := json.Unmarshal(body, &page); err != nil {
			return nil, fmt.Errorf("parsing response: %w", err)
		}

		// convert to versions
		converted, err := formatAssetVersions(page.AssetVersions)
		if err != nil {
			return nil, fmt.Errorf("cannot format to package version: %w", err)
		}
		versions = append(versions, converted...) // add to result

		// go to next page
		if page.NextPageToken == "" {
			break
		}
		pageToken = page.NextPageToken
	}
	return versions, nil
}
