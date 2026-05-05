// Package validate provides utilities for validating .env file contents.
//
// It checks for common problems such as:
//
//   - Invalid key names (keys must match [A-Za-z_][A-Za-z0-9_]*)
//   - Values that consist entirely of whitespace
//   - Required keys that are absent or empty
//
// Basic usage:
//
//	env := map[string]string{
//		"APP_NAME": "myapp",
//		"DB_HOST":  "",
//	}
//
//	// Validate key format and value sanity
//	result := validate.Validate(env)
//	if !result.Valid() {
//		for _, issue := range result.Issues {
//			fmt.Println(issue)
//		}
//	}
//
//	// Ensure required keys are present and non-empty
//	required := []string{"APP_NAME", "DB_HOST", "DB_PORT"}
//	result = validate.ValidateKeys(env, required)
//	if !result.Valid() {
//		for _, issue := range result.Issues {
//			fmt.Println(issue)
//		}
//	}
package validate
