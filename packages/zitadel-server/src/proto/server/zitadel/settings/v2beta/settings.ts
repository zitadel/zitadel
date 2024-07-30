/* eslint-disable */

export const protobufPackage = "zitadel.settings.v2beta";

export enum ResourceOwnerType {
  RESOURCE_OWNER_TYPE_UNSPECIFIED = 0,
  RESOURCE_OWNER_TYPE_INSTANCE = 1,
  RESOURCE_OWNER_TYPE_ORG = 2,
  UNRECOGNIZED = -1,
}

export function resourceOwnerTypeFromJSON(object: any): ResourceOwnerType {
  switch (object) {
    case 0:
    case "RESOURCE_OWNER_TYPE_UNSPECIFIED":
      return ResourceOwnerType.RESOURCE_OWNER_TYPE_UNSPECIFIED;
    case 1:
    case "RESOURCE_OWNER_TYPE_INSTANCE":
      return ResourceOwnerType.RESOURCE_OWNER_TYPE_INSTANCE;
    case 2:
    case "RESOURCE_OWNER_TYPE_ORG":
      return ResourceOwnerType.RESOURCE_OWNER_TYPE_ORG;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ResourceOwnerType.UNRECOGNIZED;
  }
}

export function resourceOwnerTypeToJSON(object: ResourceOwnerType): string {
  switch (object) {
    case ResourceOwnerType.RESOURCE_OWNER_TYPE_UNSPECIFIED:
      return "RESOURCE_OWNER_TYPE_UNSPECIFIED";
    case ResourceOwnerType.RESOURCE_OWNER_TYPE_INSTANCE:
      return "RESOURCE_OWNER_TYPE_INSTANCE";
    case ResourceOwnerType.RESOURCE_OWNER_TYPE_ORG:
      return "RESOURCE_OWNER_TYPE_ORG";
    case ResourceOwnerType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}
