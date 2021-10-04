import { 
    E2E_API_CALLS_DOMAIN,
    E2E_ORG_OWNER_PW,
    E2E_ORG_OWNER_VIEWER_PW,
    E2E_ORG_PROJECT_CREATOR_PW,
} from "../playwright.config";

export interface User {
    username: string
    password: string
}

export const ORG_OWNER: User = {
    username: username("org_owner"),
    password: E2E_ORG_OWNER_PW,
}

export const ORG_OWNER_VIEWER: User = {
    username: username("org_owner_viewer"),
    password: E2E_ORG_OWNER_VIEWER_PW,
}

export const ORG_PROJECT_CREATOR: User = {
    username: username("org_project_creator"),
    password: E2E_ORG_PROJECT_CREATOR_PW,
}

function username(short: string): string {
    return `${short}_user_name@caos-demo.${E2E_API_CALLS_DOMAIN}`
}