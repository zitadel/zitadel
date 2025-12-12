import os
import re

def add_frontmatter(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        if content.startswith('---'):
            return # Already has frontmatter
        
        filename = os.path.basename(filepath)
        title = filename.replace('-', ' ').replace('_', ' ').replace('.mdx', '').replace('.md', '').title()
        
        frontmatter = f"---\ntitle: {title}\n---\n\n"
        new_content = frontmatter + content
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(new_content)
        print(f"Added frontmatter to {filepath}")
    except Exception as e:
        print(f"Error adding frontmatter to {filepath}: {e}")

def fix_syntax(filepath):
    try:
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
        
        original_content = content
        
        # Fix 1: Unescaped } in text (not in code blocks)
        # This is tricky. A simple heuristic: if } is surrounded by spaces or at end of line, and not part of JSX expression.
        # But MDX treats { } as expressions.
        # If we have "Something { value }", it's an expression.
        # If we have "JSON object: { "a": 1 }", it fails if not valid JS expression (e.g. quotes).
        # The error "Unexpected token }" often happens when } is used as text.
        # We can try to escape } as \&rbrace; if it looks like text.
        # But safer is to wrap in {'}'} or just escape it.
        
        # Regex to find } that might be text.
        # We'll skip code blocks ``` ... ```
        
        parts = re.split(r'(```[\s\S]*?```|`[^`]*`)', content)
        for i in range(0, len(parts), 2): # Even parts are outside code blocks
            # Fix unclosed <p>
            # Simple heuristic: <p> followed by blank line or another block element without </p>
            # Actually, let's just replace <p> with nothing if it's just wrapping text loosely, or ensure it's closed.
            # But better: replace <p> with <br/> or just remove it if it's redundant.
            # The error "Expected a closing tag for <p>" suggests <p> is used but not closed.
            # Let's try to close it at the end of the paragraph (double newline).
            # This is hard with regex.
            
            # Fix "Unexpected token }"
            # Replace `}` with `\&rbrace;` if it's not part of a likely JS expression.
            # This is too risky to do blindly.
            pass

        # Specific fixes for known files/patterns
        
        # Fix: <p> tags - apply to all syntax files as it seems common
        # Remove <p> tags, keeping content.
        content = content.replace('<p>', '').replace('</p>', '')
        
        # Fix: <td> tags in _auth-type.mdx
        if '_auth-type.mdx' in filepath:
            # Check for <td> without </td>
            # Maybe <td>Text
            content = re.sub(r'<td>([^<]+)(?!\s*</td>)', r'<td>\1</td>', content)

        # Fix: Unexpected closing slash /
        # Often <br/> or <img ... /> is fine.
        # But maybe `</>` or `... / >`
        # The error "Unexpected closing slash / in tag"
        # In _test_setup.mdx:12:6-12:7
        # In _unlinked_oauth.mdx:6:6-6:7
        
        # Fix: Unexpected token }
        # In _application.mdx, _generate-key.mdx, _review-config.mdx
        if '_application.mdx' in filepath or '_generate-key.mdx' in filepath or '_review-config.mdx' in filepath:
             # Replace } with &rbrace; where it seems to be text
             # Or maybe it's inside a JSON block in text?
             # e.g. "Response: { ... }"
             # We can try to escape { and } if they look like JSON in text.
             pass

        if content != original_content:
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            print(f"Fixed syntax in {filepath}")
            
    except Exception as e:
        print(f"Error fixing syntax in {filepath}: {e}")

# List of files to process
frontmatter_files = [
    'content/docs/build-and-integrate/zitadel-apis/_accessing_zitadel_api.md',
    'content/docs/manage-and-govern/legal/annex-support-services.mdx',
    'content/docs/manage-and-govern/legal/data-processing-agreement.mdx',
    'content/docs/manage-and-govern/legal/policies/acceptable-use-policy.md',
    'content/docs/manage-and-govern/legal/policies/privacy-policy.mdx',
    'content/docs/manage-and-govern/legal/policies/vulnerability-disclosure-policy.mdx',
    'content/docs/manage-and-govern/legal/service-description/cloud-service-description.md',
    'content/docs/manage-and-govern/legal/service-description/service-level-description.md',
    'content/docs/manage-and-govern/legal/subprocessors.md',
    'content/docs/manage-and-govern/legal/terms-of-service.md',
    'content/docs/operate-and-self-host/manage/configure/_login.md'
]

syntax_files = [
    'content/docs/build-and-integrate/application/_application.mdx',
    'content/docs/build-and-integrate/application/_auth-type.mdx',
    'content/docs/build-and-integrate/application/_generate-key.mdx',
    'content/docs/build-and-integrate/application/_redirect-uris.mdx',
    'content/docs/build-and-integrate/application/_review-config.mdx',
    'content/docs/build-and-integrate/identity-providers/_test_setup.mdx',
    'content/docs/build-and-integrate/identity-providers/_unlinked_oauth.mdx',
    'content/docs/build-and-integrate/identity-providers/azure-ad-saml.mdx',
    'content/docs/references/openidoauth/authrequest.mdx'
]

for fp in frontmatter_files:
    add_frontmatter(fp)

for fp in syntax_files:
    fix_syntax(fp)
