import json

# Mappings for SIDEBAR RESOURCES to FUMADOCS PATHS (Relative to docs/content/docs/apis)
# Format: "docs/apis/resources/X" -> "../reference/api/Y"
# or "apis/X" -> "X" if local.

def map_path(p):
    if not isinstance(p, str):
        return p
    
    # Local files in apis/
    if p == "apis/introduction":
        return "introduction"
    if p.startswith("apis/observability/"):
        return p.replace("apis/observability/", "observability/")
    if p.startswith("apis/assets/"):
        return p.replace("apis/assets/", "assets/")
    if p == "apis/scim2":
        return "scim2"
    if p == "apis/statuscodes":
        return "statuscodes"
    if p == "apis/v2":
        return "v2"
    if p.startswith("apis/benchmarks/"):
        return p.replace("apis/benchmarks/", "benchmarks/")
    if p == "apis/migration_v1_to_v2":
        return "migration_v1_to_v2"
    
    # Generated V2 Resources
    # "docs/apis/resources/user_service_v2/sidebar.ts" -> "../reference/api/user"
    if "user_service_v2" in p: return "../reference/api/user"
    if "session_service_v2" in p: return "../reference/api/session"
    if "oidc_service_v2" in p: return "../reference/api/oidc"
    if "saml_service_v2" in p: return "../reference/api/saml"
    if "settings_service_v2" in p: return "../reference/api/settings"
    if "feature_service_v2" in p: return "../reference/api/feature"
    if "org_service_v2" in p: return "../reference/api/org" # Note: sidebars might say org_service/v2 logic? script says org_service_v2
    if "idp_service_v2" in p: return "../reference/api/idp"
    if "webkey_service_v2" in p: return "../reference/api/webkey"
    if "action_service_v2" in p: return "../reference/api/action"
    if "instance_service_v2" in p: return "../reference/api/instance"
    if "project_service_v2" in p: return "../reference/api/project"
    if "application_service_v2" in p: return "../reference/api/application"
    if "authorization_service_v2" in p: return "../reference/api/authorization"
    if "internal_permission_service_v2" in p: return "../reference/api/internal_permission"

    # Generated V1 Resources
    if "resources/auth" in p: return "../reference/api-v1/auth"
    if "resources/mgmt" in p: return "../reference/api-v1/management"
    if "resources/admin" in p: return "../reference/api-v1/admin"
    if "resources/system" in p: return "../reference/api-v1/system"

    # Fallback: remove apis/ prefix if present
    if p.startswith("apis/"):
        return p.replace("apis/", "")
    
    return p

apis_sidebar = [
    "apis/introduction",
    {
      "type": "category",
      "label": "Core Resources",
      "items": [
        {
          "type": "category",
          "label": "V2",
          "items": [
            { "type": "category", "label": "User", "items": ["docs/apis/resources/user_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Session", "items": ["docs/apis/resources/session_service_v2/sidebar.ts"] },
            { "type": "category", "label": "OIDC", "items": ["docs/apis/resources/oidc_service_v2/sidebar.ts"] },
            { "type": "category", "label": "SAML", "items": ["docs/apis/resources/saml_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Settings", "items": ["docs/apis/resources/settings_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Feature", "items": ["docs/apis/resources/feature_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Organization", "items": ["docs/apis/resources/org_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Identity Provider", "items": ["docs/apis/resources/idp_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Web Key", "items": ["docs/apis/resources/webkey_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Action", "items": ["docs/apis/resources/action_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Instance", "items": ["docs/apis/resources/instance_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Project", "items": ["docs/apis/resources/project_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Application", "items": ["docs/apis/resources/application_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Authorizations", "items": ["docs/apis/resources/authorization_service_v2/sidebar.ts"] },
            { "type": "category", "label": "Internal Permissions", "items": ["docs/apis/resources/internal_permission_service_v2/sidebar.ts"] },
          ],
        },
        {
          "type": "category",
          "label": "V1",
          "items": [
            { "type": "category", "label": "Authenticated User", "items": ["docs/apis/resources/auth/sidebar.ts"] },
            { "type": "category", "label": "Organization Objects", "items": ["docs/apis/resources/mgmt/sidebar.ts"] },
            { "type": "category", "label": "Instance Objects", "items": ["docs/apis/resources/admin/sidebar.ts"] },
            { "type": "category", "label": "Instance Lifecycle", "items": ["docs/apis/resources/system/sidebar.ts"] },
            "apis/migration_v1_to_v2"
          ],
        },
        {
          "type": "category",
          "label": "Assets",
          "items": ["apis/assets/assets"],
        },
      ],
    },
    {
      "type": "category",
      "label": "Observability",
      "items": [
        "apis/observability/metrics",
        "apis/observability/health",
      ],
    },
    {
      "type": "category",
      "label": "Provision Users",
      "items": ["apis/scim2"],
    },
    {
      "type": "doc",
      "label": "gRPC Status Codes",
      "id": "apis/statuscodes",
    },
    {
      "type": "link",
      "label": "Rate Limits (Cloud)",
      "href": "/legal/policies/rate-limit-policy",
    },
    {
      "type": "category",
      "label": "Benchmarks",
      "items": [
        {
          "type": "category",
          "label": "v2.65.0",
          "items": ["apis/benchmarks/v2.65.0/machine_jwt_profile_grant/index"],
        },
        {
          "type": "category",
          "label": "v2.66.0",
          "items": ["apis/benchmarks/v2.66.0/machine_jwt_profile_grant/index"],
        },
        {
          "type": "category",
          "label": "v2.70.0",
          "items": [
            "apis/benchmarks/v2.70.0/machine_jwt_profile_grant/index",
            "apis/benchmarks/v2.70.0/oidc_session/index",
          ],
        },
        {
          "type": "category",
          "label": "v4",
          "items": [
            "apis/benchmarks/v4/add_session/index",
            "apis/benchmarks/v4/human_password_login/index",
            "apis/benchmarks/v4/introspect/index",
            "apis/benchmarks/v4/machine_client_credentials_login/index",
            "apis/benchmarks/v4/machine_jwt_profile_grant/index",
            "apis/benchmarks/v4/machine_pat_login/index",
            "apis/benchmarks/v4/manipulate_user/index",
            "apis/benchmarks/v4/oidc_session/index",
            "apis/benchmarks/v4/otp_session/index",
            "apis/benchmarks/v4/password_session/index",
            "apis/benchmarks/v4/user_info/index",
          ],
        },
      ],
    },
]

def map_item(item):
    if isinstance(item, str):
        return map_path(item)
    if item["type"] == "category":
        pages = []
        for i in item.get("items", []):
            mapped = map_item(i)
            # Flatten if checking sidebar.ts which maps to a single folder path
            if isinstance(mapped, str) and mapped.startswith("../reference"):
                 # The item was a sidebar.ts import which we mapped to a folder path.
                 # In Fumadocs, including a folder path includes its children.
                 # sidebar.js had { items: sidebar_api_... } which put the Children of that sidebar into items.
                 # So we should just include the folder path.
                 pages.append(mapped)
            else:
                 pages.append(mapped)
        return {
            "title": item["label"],
            "pages": pages
        }
    if item["type"] == "link":
        return {
            "title": item["label"],
            "url": item["href"]
        }
    if item["type"] == "doc":
         return {
             "title": item["label"],
             "url": map_path(item["id"])
         }
    return None

meta_pages = [map_item(i) for i in apis_sidebar]

print(json.dumps({"root": True, "pages": meta_pages}, indent=2))
