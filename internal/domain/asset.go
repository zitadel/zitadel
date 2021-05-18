package domain

import "time"

const (
	UsersAssetPath  = "users"
	AvatarAssetPath = "/avatar"

	policyPrefix          = "policy"
	labelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = labelPolicyPrefix + "/logo"
	labelPolicyIconPrefix = labelPolicyPrefix + "/icon"
	labelPolicyFontPrefix = labelPolicyPrefix + "/font"
	Dark                  = "dark"

	CssPath              = labelPolicyPrefix + "/css"
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
