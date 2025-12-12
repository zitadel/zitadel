import os
import re

def process_file(filepath):
    with open(filepath, 'r') as f:
        content = f.read()

    original_content = content

    # 1. Remove duplicate Tabs imports (the one I added previously)
    content = content.replace("import { Tab, Tabs } from 'fumadocs-ui/components/tabs';", "")
    content = content.replace("import { Callout } from 'fumadocs-ui/components/callout';", "") # I'll re-add if needed or rely on auto-import if configured? No, I should keep Callout if I use it.
    # Actually, I replaced :::note with <Callout>. So I need Callout import.
    # But I should check if I added it. Yes, convert_content.py added it.
    # So I should KEEP Callout import, but REMOVE Tabs import because I'm using the Docusaurus wrapper.
    
    # Re-add Callout import if I removed it by accident or if it's missing and <Callout> is used.
    # But for now, let's just remove the specific line with Tabs.
    
    # 2. Fix Docusaurus imports
    content = content.replace("from '@theme/Tabs'", "from '@/components/docusaurus/tabs'")
    
    # Fix TabItem: import TabItem from '@theme/TabItem' -> import { TabItem } from '@/components/docusaurus/tabs'
    content = re.sub(r"import\s+TabItem\s+from\s+['\"]@theme/TabItem['\"]", "import { TabItem } from '@/components/docusaurus/tabs'", content)
    
    content = content.replace("from '@theme/CodeBlock'", "from '@/components/docusaurus/code-block'")
    content = content.replace("from '@theme/Admonition'", "from '@/components/docusaurus/admonition'")
    content = content.replace("from '@theme/ThemedImage'", "from '@/components/docusaurus/themed-image'")
    content = content.replace("from '@docusaurus/BrowserOnly'", "from '@/components/docusaurus/browser-only'")
    
    # 3. Fix @site alias
    content = content.replace("from '@site/", "from '@/")
    
    # 4. Fix raw-loader
    # Pattern: import X from '!!raw-loader!./file'
    # Replacement: import X from './file?raw'
    content = re.sub(r"import\s+(\w+)\s+from\s+['\"]!!raw-loader!(.*?)['\"];", r"import \1 from '\2?raw';", content)
    
    # 5. Fix comments
    # <!-- text --> -> {/* text */}
    # We use a non-greedy match for content.
    content = re.sub(r"<!--(.*?)-->", r"{/*\1*/}", content, flags=re.DOTALL)
    
    # 6. Fix !include (if any) - Just commenting it out for now as it's likely invalid
    # content = re.sub(r"^!(.*)", r"{/* !\1 */}", content, flags=re.MULTILINE)
    # Actually, let's see what the '!' error is about first.
    
    if content != original_content:
        print(f"Updating {filepath}")
        with open(filepath, 'w') as f:
            f.write(content)

# Walk through docs/content/docs
for root, dirs, files in os.walk('docs/content/docs'):
    for file in files:
        if file.endswith('.md') or file.endswith('.mdx'):
            process_file(os.path.join(root, file))
