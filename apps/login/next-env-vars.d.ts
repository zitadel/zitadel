declare namespace NodeJS {
  interface ProcessEnv {
    // Allow any environment variable that matches the pattern
    [key: `${string}_AUDIENCE`]: string; // The system api url
    [key: `${string}_SYSTEM_USER_ID`]: string; // The service user id
    [key: `${string}_SYSTEM_USER_PRIVATE_KEY`]: string; // The service user private key

    AUDIENCE: string; // The fallback system api url
    SYSTEM_USER_ID: string; // The fallback service user id
    SYSTEM_USER_PRIVATE_KEY: string; // The fallback service user private key

    /**
     * Self hosting: The Zitadel API url
     */
    ZITADEL_API_URL: string;

    /**
     * Takes effect only if ZITADEL_API_URL is not empty.
     * This is only relevant if Zitadels runtime has the ZITADEL_INSTANCEHOSTHEADERS config changed.
     * The default is x-zitadel-instance-host.
     * Most users don't need to set this variable.
     */
    ZITADEL_INSTANCE_HOST_HEADER: string;

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
