# Usage
Standard usage guide.

```json
{
    "mods": [
        {
            "path": "pd2mm/mods", // The storage of mod archives (e.g. pd2mm/mods/MyMod.zip)
            "export": "pd2mm/export/mods", // Where mods will be exported to. (This is what the game will read)
            "extract": "pd2mm/extract/mods", // Where mods will be extracted to temporarily.
            "include": [
                {
                    "path": "MyMod", // The path to look for.
                    "to": "{export}/MyMod" // Where to copy it to.
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
                    //  copy and rename both support variables such as {path}, {export}, {extract}
                    //  which will be replaced with the corresponding key in the JSON.
                    "from": "pd2mm/patches/mods/MyMod",
                    "to": "{export}/MyMod"
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