/**
 * AssetService's public enum contract. Account-only build boundary: the native service
 * imports this enum from an administrative component. Keep that component out of Account.
 * build.mjs verifies these ordinal values against the unmodified native source.
 */
export enum Theme {
  DARK,
  LIGHT,
}
