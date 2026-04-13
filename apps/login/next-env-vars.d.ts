declare namespace NodeJS {
  interface ProcessEnv {
    // Allow any environment variable that matches the pattern
    [key: `${string}_AUDIENCE`]: string; // The system api url
    [key: `${string}_SYSTEM_USER_ID`]: string; // The service account id
    [key: `${string}_SYSTEM_USER_PRIVATE_KEY`]: string; // The service account private key
    [key: `${string}_SYSTEM_USER_PRIVATE_KEY_FILE`]: string; // The service account private key file path

    AUDIENCE: string; // The fallback system api url
    SYSTEM_USER_ID: string; // The fallback service account id
    SYSTEM_USER_PRIVATE_KEY: string; // The fallback service account private key
    SYSTEM_USER_PRIVATE_KEY_FILE: string; // The fallback service account private key file path

    /**
     * The Zitadel API url
     */
    ZITADEL_API_URL: string;

    /**
     * The service account token
     * If ZITADEL_SERVICE_USER_TOKEN is set, its value is used.
     * If ZITADEL_SERVICE_USER_TOKEN is not set but ZITADEL_SERVICE_USER_TOKEN_FILE is set, the application blocks until the file is created.
     * As soon as the file exists, its content is read and ZITADEL_SERVICE_USER_TOKEN is set.
     */
    ZITADEL_SERVICE_USER_TOKEN: string;

    /**
     * Path to a private key file for login client JWT authentication.
     * When set, the login service reads this key and signs JWTs with a
     * hardcoded subject of "login-client".
     * AUDIENCE defaults to ZITADEL_API_URL if not explicitly set.
     */
    ZITADEL_LOGINCLIENT_KEYFILE?: string;

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
     * Optional: override the redirect URI after successful login.
     * If the value starts with "/", it will be used as a relative path and prepended with the host of the request (useful for rewrites).
     * Otherwise, it will use the value as an absolute redirect URI.
     * Takes precedence over organization settings.
     */
    DEFAULT_REDIRECT_URI?: string;

    /**
     * Optional: Comma-separated list of additional allowed origins for Server Actions.
     * Origins should include the protocol, e.g., 'https://zitadel.com,http://localhost:3000'.
     * If not set, it defaults to an empty list, allowing only same-origin requests.
     */
    SERVER_ACTION_ALLOWED_ORIGINS?: string;

    /**
     * Optional: Enable automatic code submission on page load.
     * Set to "true" to auto-submit verification codes (e.g. email verification, Email OTP).
     * Default behavior (undefined) requires users to click a Submit button,
     * which is safer for environments with enterprise email link scanners.
     */
    NEXT_PUBLIC_AUTO_SUBMIT_CODE?: string;

    /**
     * Optional: Enable the SWR in-memory cache for API requests globally.
     * Defaults to true. Set to "false" to disable completely.
     */
    API_CACHE_ENABLED?: string;

    /**
     * Optional: JSON string to configure the cache TTLs (in minutes) and size limits for specific backend API routes or global fallbacks.
     * Example: '{"defaultMinutes": 15, "longMinutes": 60, "maxSize": 200, "getBrandingSettings": 120}'
     * 
     * Properties:
     * - \`defaultMinutes\`: The globally utilized default TTL in minutes (falls back to 15 if not set).
     * - \`longMinutes\`: The TTL utilized string for long-cached routes, like branding/translation (falls back to 60 if not set).
     * - \`maxSize\`: Maximum number of entries the in-memory cache may hold. Oldest entries are evicted when capacity is reached (defaults to 100).
     * - \`[route_name]\`: Explicit overrides per specific API method (e.g., \`getHostedLoginTranslation\`).
     */
    API_CACHE_CONFIG?: string;

    /**
     * Optional: Disable OpenTelemetry instrumentation.
     * Set to "true" to bypass OTEL initialization.
     * In local development (NODE_ENV=development), it is disabled by default unless explicitly set to "false".
     */
    OTEL_SDK_DISABLED?: string;
  }
}
