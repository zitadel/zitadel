declare namespace NodeJS {
  interface ProcessEnv {
    // Allow any environment variable that matches the pattern
    [key: `${string}_AUDIENCE`]: string; // The system api url
    [key: `${string}_AUDIENCE`]: string; // The service user id
    [key: `${string}_AUDIENCE`]: string; // The service user private key

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
