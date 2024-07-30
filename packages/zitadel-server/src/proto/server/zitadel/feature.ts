/* eslint-disable */

export const protobufPackage = "zitadel.feature.v1";

export enum InstanceFeature {
  INSTANCE_FEATURE_UNSPECIFIED = 0,
  INSTANCE_FEATURE_LOGIN_DEFAULT_ORG = 1,
  UNRECOGNIZED = -1,
}

export function instanceFeatureFromJSON(object: any): InstanceFeature {
  switch (object) {
    case 0:
    case "INSTANCE_FEATURE_UNSPECIFIED":
      return InstanceFeature.INSTANCE_FEATURE_UNSPECIFIED;
    case 1:
    case "INSTANCE_FEATURE_LOGIN_DEFAULT_ORG":
      return InstanceFeature.INSTANCE_FEATURE_LOGIN_DEFAULT_ORG;
    case -1:
    case "UNRECOGNIZED":
    default:
      return InstanceFeature.UNRECOGNIZED;
  }
}

export function instanceFeatureToJSON(object: InstanceFeature): string {
  switch (object) {
    case InstanceFeature.INSTANCE_FEATURE_UNSPECIFIED:
      return "INSTANCE_FEATURE_UNSPECIFIED";
    case InstanceFeature.INSTANCE_FEATURE_LOGIN_DEFAULT_ORG:
      return "INSTANCE_FEATURE_LOGIN_DEFAULT_ORG";
    case InstanceFeature.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
