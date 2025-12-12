import os
import re

def fix_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()
    
    original_content = content
    
    # 1. Fix quotes in Callout title (Specific case found in logs)
    content = content.replace('title="Never store and commit secrets in the ".env" or settings.py file"', "title='Never store and commit secrets in the \".env\" or settings.py file'")
    
    # 2. Fix <p> in Callout title (Specific case found in logs)
    content = content.replace('title="<p>"', 'title="Note"')
    
    # 3. Remove stray </p> inside Callout if it matches the pattern we saw
    # This regex matches </p> on a line by itself, possibly with whitespace
    content = re.sub(r"^\s*</p>\s*$", "", content, flags=re.MULTILINE)
    
    # 4. Fix !include (if any)
    # content = re.sub(r"^!(.*)", r"{/* !\1 */}", content, flags=re.MULTILINE)
    
    if content != original_content:
        print(f"Fixed JSX in {filepath}")
        with open(filepath, 'w') as f:
            f.write(content)

# Walk
for root, dirs, files in os.walk('docs/content/docs'):
    for file in files:
        if file.endswith('.md') or file.endswith('.mdx'):
            fix_file(os.path.join(root, file))
