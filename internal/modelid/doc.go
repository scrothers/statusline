// Package modelid decodes a raw Claude model identifier — from the
// first-party API, Amazon Bedrock, Google Vertex AI, OpenRouter, or a
// corporate API gateway that reorders the family/version segments — into a
// normalized, human-readable "<Family> <Version>" label. It is a pure leaf
// package: no I/O, no dependency on the session input schema, so it's
// reusable from both the model segment and the degraded fallback line.
package modelid
