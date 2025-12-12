import os

# We are running from /workspaces/zitadel
root_dir = "docs/content/docs"

print(f"Scanning {root_dir}...")

for dirpath, dirnames, filenames in os.walk(root_dir):
    for f in filenames:
        if f.startswith("_") and f.endswith(".mdx"):
            filepath = os.path.join(dirpath, f)
            try:
                with open(filepath, "r") as file:
                    content = file.read()
                
                if not content.strip().startswith("---"):
                    print(f"Adding frontmatter to {filepath}")
                    # Infer title from filename
                    title = f.replace("_", " ").replace(".mdx", "").strip().title()
                    
                    new_content = f"---\ntitle: {title}\n---\n\n{content}"
                    
                    with open(filepath, "w") as file:
                        file.write(new_content)
            except Exception as e:
                print(f"Error processing {filepath}: {e}")
