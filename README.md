Run `go run ./cmd/migrations` to run schema. Make sure you have `DATABASE_URL` env var set as a conn string.

Register and hit the `/images` endpoint to upload an image. It should appear in your folder specified in `VOLUME_PATH`
environment variable.
