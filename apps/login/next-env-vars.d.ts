declare namespace NodeJS {
  interface ProcessEnv {
    /**
     * Multitenancy: The system api url
     */
    AUDIENCE: string;

    /**
     * Multitenancy: The service user id
     */
    SYSTEM_USER_ID: string;

    /**
     * Multitenancy: The service user private key
     */
    SYSTEM_USER_PRIVATE_KEY: string;

    /**
     * Self hosting: The instance url
     */
    ZITADEL_API_URL: string;

    /**
     * Self hosting: The service user id
     */
    ZITADEL_USER_ID: string;
    /**
     * Self hosting: The service user token
     */
    ZITADEL_USER_TOKEN: string;

    /**
     * Optional: wheter a user must have verified email
     */
    EMAIL_VERIFICATION: string;
  }
}
