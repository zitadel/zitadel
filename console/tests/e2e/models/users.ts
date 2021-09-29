import { 
    API_CALLS_DOMAIN,
    ORG_OWNER_PW,
    ORG_OWNER_VIEWER_PW,
    ORG_PROJECT_CREATOR_PW,
} from "./env";

export interface User {
    username: string
    password: string
}

export const ORG_OWNER: User = {
    username: username("org_owner"),
    password: ORG_OWNER_PW,
}

export const ORG_OWNER_VIEWER: User = {
    username: username("org_owner_viewer"),
    password: ORG_OWNER_VIEWER_PW,
}

export const ORG_PROJECT_CREATOR: User = {
    username: username("org_project_creator"),
    password: ORG_PROJECT_CREATOR_PW,
}

function username(short: string): string {
    return `${short}_user_name@caos-demo.${API_CALLS_DOMAIN}`
}