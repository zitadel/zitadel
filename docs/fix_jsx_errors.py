import os
import re

def fix_unclosed_p_tags(content):
    content = re.sub(r'<p>([\s\S]*?)\s*(</td>)', r'<p>\1</p>\n\2', content)
    content = re.sub(r'<p>([\s\S]*?)\s*(</tr>)', r'<p>\1</p>\n\2', content)
    content = re.sub(r'<p>([\s\S]*?)\Z', r'<p>\1</p>', content)
    return content

def process_file(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    original_content = content
    
    if 'python-django.mdx' in filepath:
        content = content.replace('title="Additional permission checks could be done with "permission_required" from "django.contrib.auth.decorators" also described in the [Django documentation](https://docs.djangoproject.com/en/5.0/topics/auth/customizing/#custom-permissions)."', 
                                  "title='Additional permission checks could be done with \"permission_required\" from \"django.contrib.auth.decorators\" also described in the [Django documentation](https://docs.djangoproject.com/en/5.0/topics/auth/customizing/#custom-permissions).'")
    
    if 'cloudflare-oidc.mdx' in filepath:
        content = content.replace('title="Cloudflare will return an error "User email was not returned. API permissions are likely incorrect". Enable to send the user information inside the token on your client settings.">',
                                  "title='Cloudflare will return an error \"User email was not returned. API permissions are likely incorrect\". Enable to send the user information inside the token on your client settings.'>")

    if '_create-user.mdx' in filepath:
        content = content.replace('title="If you started with Zitadel before version 3, you might have the "Human User [deprecated]" UI.">',
                                  "title='If you started with Zitadel before version 3, you might have the \"Human User [deprecated]\" UI.'>")

    content = fix_unclosed_p_tags(content)
    
    if content != original_content:
        print(f"Fixed {filepath}")
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)

files_to_check = [
    'docs/content/docs/build-and-integrate/application/_auth-type.mdx',
    'docs/content/docs/build-and-integrate/application/_redirect-uris.mdx',
    'docs/content/docs/build-and-integrate/identity-providers/_test_setup.mdx',
    'docs/content/docs/build-and-integrate/identity-providers/_unlinked_oauth.mdx',
    'docs/content/docs/build-and-integrate/identity-providers/azure-ad-saml.mdx',
    'docs/content/docs/references/openidoauth/authrequest.mdx',
    'docs/content/docs/build-and-integrate/examples/login/python-django.mdx',
    'docs/content/docs/build-and-integrate/services/cloudflare-oidc.mdx',
    'docs/content/docs/manage-and-govern/console/_create-user.mdx'
]

for relative_path in files_to_check:
    full_path = os.path.join('/workspaces/zitadel', relative_path)
    if os.path.exists(full_path):
        process_file(full_path)
    else:
        print(f"File not found: {full_path}")
