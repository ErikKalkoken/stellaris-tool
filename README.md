# Stellaris Tool

## Parser

This is general parser for Paradox save files. It is mostly tested with Stellaris, but should work with save files from other Paradox games too.

The parser takes a reader for a safe file stream and returns a map.

Here are some noteworthy special cases and how the parser handles them:

- "none" values are converted to nil and will be nulls in the JSON output
- Identifiers used as keys in objects (incl. "none) are converted to strings
- Arrays of numbers are converted to float arrays
- Keyword keys and integer keys are converted to strings. So that all map keys are strings.
