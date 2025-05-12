# Formatting
Standard formatting guide.

## Basics
- *ALL* Variables, parameters, arguments, function names, etc., **MUST** be derived from the English language.
- *ALL* Variables, parameters, arguments, function names, etc., **MUST** use camelCase.
- *ALL* code files need to have the AGPL license header. See [LICENSE-Notice.txt](./LICENSE-Notice.txt)

## Complexity
I am generally lenient on complexity in regards to linters. Complexity is subjective, and linters treat it as objective. Generally speaking *all* code flow should be understandable within 1~2 passes. If it takes multiple passes to fully understand. Refactor it to be less complex.

## Parameter Naming Conventions
- Parameters should be short (~1 [one] word).
- Parameters should be shortened by removing vowels (destination ❌ -> dest ✅).
    - Prioritize *words* over *1 (one) letter*.

## Variable Naming Conventions
- Global variables **SHOULD** use longer forms of words (dest ❌ -> destination ✅).
    - Certain keywords can and should be shortened (Error ❌ -> Err ✅).
    - Global variables should *always* be discernible without *context.*
- Local variables should **NOT** use *1 (one) letter*.
- Local variables **SHOULD** use longer forms of words (dest ❌ -> destination ✅).
    - Local variables with *2 (two) or more words* should be shortened to *1 (one) word*, or by by removing vowels (destination ❌ -> dest ✅)

### Good (✅) Example
```go
// Function name uses camel case ✅
func MyFunction(src string) { // Paramter name is short ✅
    source := src // Local variable is long ✅
    // ...
}
```

### Bad (❌) Example
```go
// Function name is not camel case (Myfunction ❌ -> MyFunction ✅)
func Myfunction(source string) { // Parameter not shortened (source ❌ -> src ✅)
    src := source // Local variable shortened (src ❌ -> source ✅)
    // ...
}
```