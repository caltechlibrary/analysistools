
# TODO

## Bugs

- [ ] Missing documentation show examples and integration choices

## Up next

- [ ] exclude list matches using `strings.Contains()`, this might not be what we want
- [ ] token matching right not appears to be line limited, need to write a test to confirm, if true I need to track line number but work on proximity by stream to find all occurences
- [ ] the phrase check doesn't check with prefix versus exact match, should handle the case better where it is an exact matchin and include
- [ ] need option to save in SQLite3 database rather than just output a CSV, this will allow for additional processing and analysis
