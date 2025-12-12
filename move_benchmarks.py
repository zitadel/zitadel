import os
import shutil
import re

DOCS_ROOT = "."
CONTENT_ROOT = os.path.join(DOCS_ROOT, "content/docs/references/benchmarks")
DATA_ROOT = os.path.join(DOCS_ROOT, "src/data/benchmarks")

def move_benchmarks():
    if not os.path.exists(CONTENT_ROOT):
        print(f"Content root not found: {CONTENT_ROOT}")
        return

    for root, dirs, files in os.walk(CONTENT_ROOT):
        for file in files:
            if file == "output.json" or file == "_output.json":
                # Found a benchmark output file
                source_path = os.path.join(root, file)
                
                # Calculate relative path from benchmarks root
                rel_path = os.path.relpath(root, CONTENT_ROOT)
                
                # Construct destination path
                dest_dir = os.path.join(DATA_ROOT, rel_path)
                dest_path = os.path.join(dest_dir, "output.json") # Always name it output.json in data
                
                # Create destination directory
                os.makedirs(dest_dir, exist_ok=True)
                
                # Move file
                print(f"Moving {source_path} to {dest_path}")
                shutil.move(source_path, dest_path)
                
                # Update index.mdx
                index_path = os.path.join(root, "index.mdx")
                if os.path.exists(index_path):
                    update_import(index_path, rel_path)
                else:
                    print(f"Warning: No index.mdx found in {root}")

def update_import(file_path, rel_path):
    with open(file_path, 'r') as f:
        content = f.read()
    
    # Construct new import path
    # rel_path is like "v2.65.0/machine_jwt_profile_grant"
    # We want "@/src/data/benchmarks/v2.65.0/machine_jwt_profile_grant/output.json"
    new_import_path = f"@/src/data/benchmarks/{rel_path}/output.json"
    
    # Regex to find the import
    # import data from './output.json' or './_output.json'
    pattern = r"import\s+(\w+)\s+from\s+['\"](\./_?output\.json)['\"]"
    
    def replace(match):
        var_name = match.group(1)
        print(f"Updating import in {file_path}: {match.group(0)} -> import {var_name} from '{new_import_path}'")
        return f"import {var_name} from '{new_import_path}'"
    
    new_content = re.sub(pattern, replace, content)
    
    if new_content != content:
        with open(file_path, 'w') as f:
            f.write(new_content)
    else:
        print(f"No import match found in {file_path}")

if __name__ == "__main__":
    move_benchmarks()
