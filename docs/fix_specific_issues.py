import os

def fix_auth_type():
    filepath = 'content/docs/build-and-integrate/application/_auth-type.mdx'
    if not os.path.exists(filepath):
        print(f"File not found: {filepath}")
        return

    with open(filepath, 'r', encoding='utf-8') as f:
        content = f.read()
    
    # Fix corrupted text
    content = content.replace('a</td>n</td>y', 'any')
    
    # Fix stray </td>
    content = content.replace('<td>\n      </td><', '<td>\n      <')
    content = content.replace('<td>\n      </td>', '<td>\n      ') # Just in case
    
    with open(filepath, 'w', encoding='utf-8') as f:
        f.write(content)
    print(f"Fixed {filepath}")

def fix_frontmatter():
    files = [
        'content/docs/manage-and-govern/legal/annex-support-services.mdx',
        'content/docs/manage-and-govern/legal/data-processing-agreement.mdx',
        'content/docs/manage-and-govern/legal/policies/acceptable-use-policy.md',
        'content/docs/manage-and-govern/legal/policies/privacy-policy.mdx',
        'content/docs/manage-and-govern/legal/policies/vulnerability-disclosure-policy.mdx',
        'content/docs/manage-and-govern/legal/service-description/cloud-service-description.md',
        'content/docs/manage-and-govern/legal/service-description/service-level-description.md',
        'content/docs/manage-and-govern/legal/subprocessors.md',
        'content/docs/manage-and-govern/legal/terms-of-service.md'
    ]
    
    for filepath in files:
        if not os.path.exists(filepath):
            print(f"File not found: {filepath}")
            continue
            
        with open(filepath, 'r', encoding='utf-8') as f:
            content = f.read()
            
        # Remove trailing space after ---
        # Regex might be safer
        import re
        content = re.sub(r'^---\s+$', '---', content, flags=re.MULTILINE)
        
        with open(filepath, 'w', encoding='utf-8') as f:
            f.write(content)
        print(f"Fixed frontmatter in {filepath}")

if __name__ == '__main__':
    fix_auth_type()
    fix_frontmatter()
