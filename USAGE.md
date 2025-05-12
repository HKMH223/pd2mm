# Usage
Standard usage guide.

```jsonc
{
    "mods": [
        {
            "path": "pd2mm/mods", // The storage of mod archives (e.g. pd2mm/mods/MyMod.zip)
            "output": "pd2mm/output/mods", // Where mods will be output to. (This is what the game will read)
            "extract": "pd2mm/extract/mods", // Where mods will be extracted to temporarily.
            "export": "pd2mm/export/mods", // Export is not considered a temp directory and will only ever copy files to it, never delete.
            "include": [
                {
                    "path": "MyMod", // The path to look for.
                    "to": "{output}/MyMod" // Where to copy it to.
                }
            ],
            "exclude": [
                "MyMod/MyOtherPath", // A path to exclude.
            ],
            "expects": [
                {
                    "path": [ // A folder or file the program will "expect" to find."
                        "mod.txt"
                    ],
                    "require": [ // Parts of a path the expected path requires. (e.g. partA/partB/mod.txt)
                        "partA", 
                        "partB"
                    ],
                    "base": 0 // [ADVANCED] - The root path index.
                }
            ],
            // Unlike include, copy does not look through the extract folder
            "copy": [
                { 
                    // Copy from point a to point b.
                    //  copy and rename both support variables such as {path}, {output}, {extract}
                    //  which will be replaced with the corresponding key in the JSON.
                    "from": "pd2mm/patches/mods/MyMod",
                    "to": "{output}/MyMod"
                }
            ],
            "rename": [
                {
                    "path": [ // The path to look for.
                        "MyMod"
                    ],
                    "from": [ // What to replace.
                        "MyMod"
                    ],
                    "to": [ // What to replace it with.
                        "MyCoolMod"
                    ]
                },
            ]
        }
    ]
}
```