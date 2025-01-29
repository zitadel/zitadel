declare namespace NodeJS {
  interface ProcessEnv {
    /**
     * Multitenancy: The system api url
     */
    QA_AUDIENCE: string;

    /**
     * Multitenancy: The service user id
     */
    QA_SYSTEM_USER_ID: string;

    /**
     * Multitenancy: The service user private key
     */
    QA_SYSTEM_USER_PRIVATE_KEY: string;

    /**
     * Multitenancy: The system api url for prod environment
     */
    PROD_AUDIENCE: string;

    /**
     * Multitenancy: The service user id for prod environment
     */
    PROD_SYSTEM_USER_ID: string;

    /**
     * Multitenancy: The service user private key for prod environment
     */
    PROD_SYSTEM_USER_PRIVATE_KEY: string;

    /**
     * Self hosting: The instance url
     */
    ZITADEL_API_URL: string;

    /**
     * Self hosting: The service user id
     */
    ZITADEL_SERVICE_USER_ID: string;
    /**
     * Self hosting: The service user token
     */
    ZITADEL_SERVICE_USER_TOKEN: string;

    /**
     * Optional: wheter a user must have verified email
     */
    EMAIL_VERIFICATION: string;
  }
}
