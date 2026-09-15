# Goal: 
# When going through each repo, check each csproj file (C# project file), then Bump specific packages, 
# if any csproj was changed, go through all files and replace a multi-line string to another string
# Note: If the selected repositories definetly have a reference to the package, we don't need to check if the csproj files 
# have been changed

import re
from pathlib import Path
import xml.etree.ElementTree as ET

TARGET_PACKAGES = [
    "ABC.Events.SingleRegion",
    "ABC.Events.MultiRegion",
]

NEW_VERSION = '2.0.0'

repo_root = Path('.')
csproj_changed = False

for csproj in repo_root.rglob("*.csproj"):
    tree = ET.parse(csproj)
    root = tree.getroot()
    file_updated = False

    for elem in root.iter():
        if elem.tag.endswith("PackageReference"):
            pkg = elem.attrib.get("Include", "")
            if pkg in TARGET_PACKAGES:
                elem.set("Version", NEW_VERSION)
                file_updated = True

    if file_updated:
        #write the new xml file and ensure a white space at the end of the file, because xml will remove it.

        # check whether a trailing white space exist
        orig_text = csproj.read_text(encoding='utf-8')
        trailing_whitespace = re.search(r'(\s*\Z)', orig_text).group(1)

        # Update the xml
        tree.write(csproj, encoding='utf-8', xml_declaration=False)

        # add the trailing white space if it existed in the original csproj
        if trailing_whitespace:
            with csproj.open("a", encoding='utf-8') as f:
                f.write(trailing_whitespace)

        print(f"Bumped versions in {csproj}")
        csproj_changed = True

if csproj_changed:
    # replace a specific multi-line service injection to a single line in all files
    pat = re.compile(
        r"services.AddTransient<\s*IABCEvent,\s*ABCEvent\s*>"
        , flags=re.IGNORECASE
    )
    repl = "services.AddTransient<IABCEvent, ABCEvent>"

    for src in repo_root.rglob("*.*"):
        if not src.is_file():
            continue
        try:
            text = src.read_text(encoding='utf-8', errors='ignore')
        except Exception:
            continue
        new_text = pat.sub(repl, text)
        if new_text != text:
            src.write_text(new_text, encoding='utf-8')
            print(f" patched in {src}")