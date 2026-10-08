# Integration tests

Backend unit tests currently live beside their Go packages and can be run with
`cd backend && go test ./...`.

Place cross-component tests in this directory as the API and moderation
workflows grow. Integration tests should use an isolated MongoDB instance and
must not depend on production data or credentials.
