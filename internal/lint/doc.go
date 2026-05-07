// Package lint provides heuristic quality checks for parsed .env files.
//
// # Overview
//
// After parsing a .env file with the envfile package, callers can pass the
// resulting []envfile.Entry slice to [Run] to receive a list of [Finding]
values describing any detected issues.
//
// # Checks performed
//
//   - Duplicate keys: the same key appears more than once (error).
//   - Shadowed OS variables: keys such as PATH or HOME that override
//     well-known operating-system variables (warn).
//   - Overly long values: values exceeding 512 characters (warn).
//   - Whitespace-only values: a value that is non-empty but consists
//     entirely of whitespace characters (warn).
//
// # Usage
//
//	entries, _ := envfile.Parse("production.env")
//	findings := lint.Run(entries)
//	for _, f := range findings {
//		fmt.Println(f)
//	}
package lint
