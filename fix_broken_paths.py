import os
import re

# Define the root directory to search
ROOT_DIR = "/workspaces/zitadel/docs"
CONTENT_DIR = os.path.join(ROOT_DIR, "content")
COMPONENTS_DIR = os.path.join(ROOT_DIR, "components")

# Define replacements
REPLACEMENTS = [
    # Fix concepts path in manage-and-govern/console
    (r"['\"]../../../concepts/(.*?)['\"]", r"'../../concepts/\1'"),
    
    # Fix apis path in build-and-integrate
    (r"['\"]../../apis/(.*?)['\"]", r"'../references/apis/\1'"),
    
    # Fix integrate path in manage-and-govern/console
    (r"['\"]../../integrate/(.*?)['\"]", r"'../../build-and-integrate/\1'"),
    
    # Fix frameworks.json path
    (r"['\"]../../frameworks.json['\"]", r"'../frameworks.json'"),
    
    # Fix CSS paths in components (if not already fixed)
    (r"['\"]../css/(.*?)['\"]", r"'@/css/\1'"),
    
    # Fix relative imports that might be missing a level or have too many
    # This is risky, so we target specific known issues
    
    # Fix token-exchange imports
    (r"['\"]../../apis/openidoauth/(.*?)['\"]", r"'../references/apis/openidoauth/\1'"),
]

def process_file(filepath):
    try:
        with open(filepath, "r", encoding="utf-8") as f:
            content = f.read()
    except Exception as e:
        print(f"Error reading {filepath}: {e}")
        return

    original_content = content
    
    for pattern, replacement in REPLACEMENTS:
        content = re.sub(pattern, replacement, content)

    if content != original_content:
        print(f"Fixed paths in {filepath}")
        with open(filepath, "w", encoding="utf-8") as f:
            f.write(content)

def main():
    # Walk content directory
    for root, dirs, files in os.walk(CONTENT_DIR):
        for file in files:
            if file.endswith(".mdx") or file.endswith(".md") or file.endswith(".js") or file.endswith(".tsx"):
                process_file(os.path.join(root, file))

    # Walk components directory
    for root, dirs, files in os.walk(COMPONENTS_DIR):
        for file in files:
            if file.endswith(".jsx") or file.endswith(".tsx") or file.endswith(".js"):
                process_file(os.path.join(root, file))

if __name__ == "__main__":
    main()
