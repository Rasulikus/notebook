migrateup:
	migrate -path migrations -database postgres://admin:password@localhost:5432/notebook?sslmode=disable up
migratedown:
	migrate -path migrations -database postgres://admin:password@localhost:5432/notebook?sslmode=disable down