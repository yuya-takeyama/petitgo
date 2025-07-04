#!/usr/bin/env python3
"""
Check if all files end with newlines and fix them if needed.
"""

import os
import glob

def check_and_fix_newlines(directory):
    """Check and fix newlines for all files in directory"""
    files = glob.glob(os.path.join(directory, "*"))
    files = [f for f in files if os.path.isfile(f)]
    
    results = []
    
    for file_path in sorted(files):
        with open(file_path, 'rb') as f:
            content = f.read()
        
        filename = os.path.basename(file_path)
        
        if len(content) == 0:
            results.append(f"❌ {filename}: Empty file")
            continue
            
        if content[-1:] != b'\n':
            results.append(f"❌ {filename}: Missing newline at end")
            # Fix it
            with open(file_path, 'ab') as f:
                f.write(b'\n')
            results.append(f"✅ {filename}: Fixed - added newline")
        else:
            results.append(f"✅ {filename}: OK - ends with newline")
    
    return results

def main():
    print("🔍 Checking newlines in examples/ directory...")
    print()
    
    results = check_and_fix_newlines("examples")
    
    for result in results:
        print(result)
    
    print()
    print("🎉 Check complete!")

if __name__ == "__main__":
    main()