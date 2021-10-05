
export enum IdpCreateType {
  OIDC = 'OIDC',
  JWT = 'JWT',
}

export interface RadioItemIdpType {
  createType: IdpCreateType;
  titleI18nKey: string;
  mdi?: string;
}

export const OIDC = {
  titleI18nKey: 'IDP.OIDC.TITLE',
  mdi: 'mdi_openid',
  createType: IdpCreateType.OIDC,
};

export const JWT = {
  titleI18nKey: 'IDP.JWT.TITLE',
  mdi: 'mdi_jwt',
  createType: IdpCreateType.JWT,
};
