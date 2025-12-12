import os
import re

# Define the root directory to search
ROOT_DIR = "/workspaces/zitadel/docs"
CONTENT_DIR = os.path.join(ROOT_DIR, "content")
COMPONENTS_DIR = os.path.join(ROOT_DIR, "components")

# Define replacements
REPLACEMENTS = [
    # Fix raw-loader imports
    (r"(import\s+.*?from\s+['\"])!!raw-loader!(.*?)(['\"])", r"\1\2?raw\3"),
    
    # Fix @theme imports
    (r"['\"]@theme/Tabs['\"]", "'@/components/docusaurus/tabs'"),
    (r"['\"]@theme/CodeBlock['\"]", "'@/components/docusaurus/code-block'"),
    (r"['\"]@theme/ThemedImage['\"]", "'@/components/docusaurus/themed-image'"),
    (r"['\"]@theme/DocCardList['\"]", "'@/components/docusaurus/doc-card-list'"),
    (r"['\"]@theme/Admonition['\"]", "'@/components/docusaurus/admonition'"),
    
    # Fix @site imports
    (r"['\"]@site/(.*?)['\"]", r"'@/\1'"),
    
    # Fix specific relative imports that are broken
    # (This is harder to do generically, but we can try to fix common patterns)
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

    # Specific fix for frameworks.json in frameworks.jsx
    if filepath.endswith("frameworks.jsx"):
        content = content.replace('../../frameworks.json', '../frameworks.json')
        content = content.replace('../css/tile.module.css', '@/css/tile.module.css')

    # Specific fix for apicard.jsx
    if filepath.endswith("apicard.jsx"):
        content = content.replace('../css/apicard.module.css', '@/css/apicard.module.css')

    # Specific fix for player.jsx
    if filepath.endswith("player.jsx"):
        content = content.replace("import ReactPlayer from 'react-player'", "import ReactPlayer from 'react-player/lazy'")

    if content != original_content:
        print(f"Fixed imports in {filepath}")
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
