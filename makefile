

clean:
	@rm *.db -r 

run:
	@go run app
runtestDB:
	@go run test/testdb.go
