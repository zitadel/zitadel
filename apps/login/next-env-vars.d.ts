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
     * The Zitadel API url
     */
    ZITADEL_API_URL: string;

    /**
     * The service user token
     * If ZITADEL_SERVICE_USER_TOKEN is set, its value is used.
     * If ZITADEL_SERVICE_USER_TOKEN is not set but ZITADEL_SERVICE_USER_TOKEN_FILE is set, the application blocks until the file is created.
     * As soon as the file exists, its content is read and ZITADEL_SERVICE_USER_TOKEN is set.
     */
    ZITADEL_SERVICE_USER_TOKEN: string;

    /**
     * Optional: wheter a user must have verified email
     */
    EMAIL_VERIFICATION: string;

    /**
     * Optional: custom request headers to be added to every request
     * Split by comma, key value pairs separated by colon
     * For example: to call the Zitadel API at an internal address, you can set:
     * `CUSTOM_REQUEST_HEADERS=Host:http://zitadel-internal:8080`
     */
    CUSTOM_REQUEST_HEADERS?: string;

    /**
     * The base path the app is served from, e.g. /ui/v2/login
     */
    NEXT_PUBLIC_BASE_PATH: string;

    /**
     * Optional: The application name shown in the login and invite emails
     */
    NEXT_PUBLIC_APPLICATION_NAME?: string;

    /**
     * Optional: override the redirect URI after successful login
     * If set to "/", it will redirect back to the host of the request (useful for rewrites)
     * Otherwise, it will use the value as an absolute redirect URI.
     * Takes precedence over organization settings.
     */
    DEFAULT_REDIRECT_URI?: string;

    /**
     * Optional: Comma-separated list of additional allowed origins for Server Actions (e.g., 'localhost:3000,zitadel.com').
     * If not set, it defaults to an empty list, allowing only same-origin requests.
     */
    SERVER_ACTION_ALLOWED_ORIGINS?: string;
  }
}
