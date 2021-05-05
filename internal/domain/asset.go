package domain

import "time"

const (
	UsersAssetPath  = "users"
	AvatarAssetPath = "/avatar"

	orgPrefix             = "org"
	iamPrefix             = "iam"
	policyPrefix          = "/policy"
	labelPolicyPrefix     = policyPrefix + "/label"
	labelPolicyLogoPrefix = labelPolicyPrefix + "/logo"
	labelPolicyIconPrefix = labelPolicyPrefix + "/icon"
	labelPolicyFontPrefix = labelPolicyPrefix + "/font"
	Dark                  = "/Dark"

	DefaultLabelPolicyLogoPath = iamPrefix + labelPolicyLogoPrefix
	DefaultLabelPolicyIconPath = iamPrefix + labelPolicyIconPrefix
	DefaultLabelPolicyFontPath = iamPrefix + labelPolicyFontPrefix

	OrgLabelPolicyLogoPath = orgPrefix + labelPolicyLogoPrefix
	OrgLabelPolicyIconPath = orgPrefix + labelPolicyIconPrefix
	OrgLabelPolicyFontPath = orgPrefix + labelPolicyFontPrefix
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
}

func GetHumanAvatarAssetPath(userID string) string {
	return UsersAssetPath + "/" + userID + AvatarAssetPath
}
