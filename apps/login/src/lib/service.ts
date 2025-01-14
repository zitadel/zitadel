import { Client, createClientFor } from "@zitadel/client";
import { createServerTransport } from "@zitadel/client/node";
import { IdentityProviderService } from "@zitadel/proto/zitadel/idp/v2/idp_service_pb";
import { OIDCService } from "@zitadel/proto/zitadel/oidc/v2/oidc_service_pb";
import { OrganizationService } from "@zitadel/proto/zitadel/org/v2/org_service_pb";
import { SessionService } from "@zitadel/proto/zitadel/session/v2/session_service_pb";
import { SettingsService } from "@zitadel/proto/zitadel/settings/v2/settings_service_pb";
import { UserService } from "@zitadel/proto/zitadel/user/v2/user_service_pb";
import { getInstanceUrl, systemAPIToken } from "./api";

type ServiceClass =
  | typeof IdentityProviderService
  | typeof UserService
  | typeof OrganizationService
  | typeof SessionService
  | typeof OIDCService
  | typeof SettingsService;

export async function createServiceForHost<T extends ServiceClass>(
  service: T,
  host: string,
) {
  let instanceUrl;
  try {
    instanceUrl = await getInstanceUrl(host);
  } catch (error) {
    console.error(
      "Could not get instance url, fallback to ZITADEL_API_URL",
      error,
    );
    instanceUrl = process.env.ZITADEL_API_URL;
  }

  const systemToken = await systemAPIToken();

  const transport = createServerTransport(systemToken, {
    baseUrl: instanceUrl,
  });

  return createClientFor<T>(service)(transport);
}

export class ServiceInitializer {
  public idpService: Client<typeof IdentityProviderService> | null = null;
  public orgService: Client<typeof OrganizationService> | null = null;
  public sessionService: Client<typeof SessionService> | null = null;
  public userService: Client<typeof UserService> | null = null;
  public oidcService: Client<typeof OIDCService> | null = null;
  public settingsService: Client<typeof SettingsService> | null = null;

  private static instance: ServiceInitializer;

  constructor(private host: string) {
    this.initializeServices();
  }

  public static async getInstance(host: string): Promise<ServiceInitializer> {
    if (!ServiceInitializer.instance) {
      ServiceInitializer.instance = new ServiceInitializer(host);
      await ServiceInitializer.instance.initializeServices();
    }
    return ServiceInitializer.instance;
  }

  async initializeServices() {
    this.idpService = await createServiceForHost(
      IdentityProviderService,
      this.host,
    );
    this.orgService = await createServiceForHost(
      OrganizationService,
      this.host,
    );
    this.sessionService = await createServiceForHost(SessionService, this.host);
    this.userService = await createServiceForHost(UserService, this.host);
    this.oidcService = await createServiceForHost(OIDCService, this.host);
    this.settingsService = await createServiceForHost(
      SettingsService,
      this.host,
    );
  }

  public getSettingsService(): Client<typeof SettingsService> {
    if (!this.settingsService) {
      throw new Error("SettingsService is not initialized");
    }
    return this.settingsService;
  }

  public getUserService(): Client<typeof UserService> {
    if (!this.userService) {
      throw new Error("UserService is not initialized");
    }
    return this.userService;
  }

  public getOrgService(): Client<typeof OrganizationService> {
    if (!this.orgService) {
      throw new Error("OrganizationService is not initialized");
    }
    return this.orgService;
  }

  public getSessionService(): Client<typeof SessionService> {
    if (!this.sessionService) {
      throw new Error("SessionService is not initialized");
    }
    return this.sessionService;
  }

  public getIDPService(): Client<typeof IdentityProviderService> {
    if (!this.idpService) {
      throw new Error("IDPService is not initialized");
    }
    return this.idpService;
  }

  public getOIDCService(): Client<typeof OIDCService> {
    if (!this.oidcService) {
      throw new Error("OIDCService is not initialized");
    }
    return this.oidcService;
  }
}
