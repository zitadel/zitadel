package domain

import "time"

const (
	UsersAssetPath  = "users"
	AvatarAssetPath = "/avatar"

	policyPrefix          = "policy"
	LabelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = LabelPolicyPrefix + "/logo"
	labelPolicyIconPrefix = LabelPolicyPrefix + "/icon"
	labelPolicyFontPrefix = LabelPolicyPrefix + "/font"
	Dark                  = "dark"

	CssPath              = LabelPolicyPrefix + "/css"
	CssVariablesFileName = "variables.css"

	LabelPolicyLogoPath = labelPolicyLogoPrefix
	LabelPolicyIconPath = labelPolicyIconPrefix
	LabelPolicyFontPath = labelPolicyFontPrefix
)

type AssetInfo struct {
	Bucket          string
	Key             string
	ETag            string
	Size            int64
	LastModified    time.Time
	Location        string
	VersionID       string
	Expiration      time.Time
	AutheticatedURL string
	ContentType     string
}

func GetHumanAvatarAssetPath(userID string) string {
	return UsersAssetPath + "/" + userID + AvatarAssetPath
}
