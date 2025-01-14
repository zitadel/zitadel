declare namespace NodeJS {
  interface ProcessEnv {
    /**
     * The system api url
     */
    AUDIENCE: string;

    /**
     * The system api service user ID
     */
    SYSTEM_USER_ID: string;

    /**
     * The service user key
     */
    SYSTEM_USER_PRIVATE_KEY: string;

    /**
     * The instance url
     */
    ZITADEL_API_URL: string;

    /**
     * The service user id for the instance
     */
    ZITADEL_USER_ID: string;

    /**
     * The service user token for the instance
     */
    ZITADEL_USER_TOKEN: string;
  }
}
