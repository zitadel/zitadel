import * as jspb from 'google-protobuf'

import * as protoc$gen$openapiv2_options_annotations_pb from '../../../protoc-gen-openapiv2/options/annotations_pb'; // proto import: "protoc-gen-openapiv2/options/annotations.proto"
import * as zitadel_settings_v2beta_settings_pb from '../../../zitadel/settings/v2beta/settings_pb'; // proto import: "zitadel/settings/v2beta/settings.proto"


export class BrandingSettings extends jspb.Message {
  getLightTheme(): Theme | undefined;
  setLightTheme(value?: Theme): BrandingSettings;
  hasLightTheme(): boolean;
  clearLightTheme(): BrandingSettings;

  getDarkTheme(): Theme | undefined;
  setDarkTheme(value?: Theme): BrandingSettings;
  hasDarkTheme(): boolean;
  clearDarkTheme(): BrandingSettings;

  getFontUrl(): string;
  setFontUrl(value: string): BrandingSettings;

  getHideLoginNameSuffix(): boolean;
  setHideLoginNameSuffix(value: boolean): BrandingSettings;

  getDisableWatermark(): boolean;
  setDisableWatermark(value: boolean): BrandingSettings;

  getResourceOwnerType(): zitadel_settings_v2beta_settings_pb.ResourceOwnerType;
  setResourceOwnerType(value: zitadel_settings_v2beta_settings_pb.ResourceOwnerType): BrandingSettings;

  getThemeMode(): ThemeMode;
  setThemeMode(value: ThemeMode): BrandingSettings;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BrandingSettings.AsObject;
  static toObject(includeInstance: boolean, msg: BrandingSettings): BrandingSettings.AsObject;
  static serializeBinaryToWriter(message: BrandingSettings, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BrandingSettings;
  static deserializeBinaryFromReader(message: BrandingSettings, reader: jspb.BinaryReader): BrandingSettings;
}

export namespace BrandingSettings {
  export type AsObject = {
    lightTheme?: Theme.AsObject,
    darkTheme?: Theme.AsObject,
    fontUrl: string,
    hideLoginNameSuffix: boolean,
    disableWatermark: boolean,
    resourceOwnerType: zitadel_settings_v2beta_settings_pb.ResourceOwnerType,
    themeMode: ThemeMode,
  }
}

export class Theme extends jspb.Message {
  getPrimaryColor(): string;
  setPrimaryColor(value: string): Theme;

  getBackgroundColor(): string;
  setBackgroundColor(value: string): Theme;

  getWarnColor(): string;
  setWarnColor(value: string): Theme;

  getFontColor(): string;
  setFontColor(value: string): Theme;

  getLogoUrl(): string;
  setLogoUrl(value: string): Theme;

  getIconUrl(): string;
  setIconUrl(value: string): Theme;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Theme.AsObject;
  static toObject(includeInstance: boolean, msg: Theme): Theme.AsObject;
  static serializeBinaryToWriter(message: Theme, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Theme;
  static deserializeBinaryFromReader(message: Theme, reader: jspb.BinaryReader): Theme;
}

export namespace Theme {
  export type AsObject = {
    primaryColor: string,
    backgroundColor: string,
    warnColor: string,
    fontColor: string,
    logoUrl: string,
    iconUrl: string,
  }
}

export enum ThemeMode { 
  THEME_MODE_UNSPECIFIED = 0,
  THEME_MODE_AUTO = 1,
  THEME_MODE_LIGHT = 2,
  THEME_MODE_DARK = 3,
}
