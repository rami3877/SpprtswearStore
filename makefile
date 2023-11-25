

clean:
	@rm DataBase -r 

run:
	@go run app
runtestDB:
	@go run test/testdb.go
