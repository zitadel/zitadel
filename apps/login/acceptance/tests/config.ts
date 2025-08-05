import fs from "fs";
import path from "path";

export interface Config {
    adminToken: string;
    zitadelApiUrl: string;
    zitadelAdminUser: string;
    sinkNotificationUrl: string;
}

export class ConfigReader {
  private readonly _config: Config = {
    adminToken: "",
    zitadelApiUrl: "",
    zitadelAdminUser: "",
    sinkNotificationUrl: "",
  };

  constructor() {
    this.load();
  }

  private load() {
    this._config.adminToken = process.env.ZITADEL_ADMIN_TOKEN || "";
    this._config.zitadelApiUrl = process.env.ZITADEL_API_URL || "http://localhost:8080";
    this._config.zitadelAdminUser = process.env.ZITADEL_ADMIN_USER || "zitadel-admin@zitadel.localhost";
    this._config.sinkNotificationUrl = process.env.SINK_NOTIFICATION_URL || "http://localhost:3333/notification";

    if (!this._config.adminToken) {
      const file = process.env.ZITADEL_ADMIN_TOKEN_FILE;
      if (!file) {
        throw new Error("ZITADEL_ADMIN_TOKEN_FILE is not set in the environment variables.");
      }
      const filePath = path.resolve(file);
      if (!fs.existsSync(filePath)) {
        throw new Error(`ZITADEL_ADMIN_TOKEN_FILE not found at path: ${filePath}`);
      }
      const token = fs.readFileSync(filePath, "utf-8").trim();
      if (!token) {
        throw new Error("ZITADEL_ADMIN_TOKEN_FILE is empty.");
      }
      this._config.adminToken = token;
    }
  }

  get config(): Config {
    return this._config;
  }
}
