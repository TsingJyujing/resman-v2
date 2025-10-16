generate_sql:
	rm -rf contoller/dao
	sqlc generate

start_server: generate_sql
	go run main.go server