# MPC Design Pattern Implemented In Go

An E-Commerce microservice supporting both HTTP and gRPC communication protocols built in Go to demo the concept of **MPC Design Pattern** which is explained in this [article](https://medium.com/@m.a.faried/mvc-or-mpc-e907f39f9e35)

- **Presentations** Containing the implementation of both HTTP & gRPC protocols.
- **Models** The api is built using the a data access layer and DTO's located in the models folder for the accessing the database.
- **Controllers** for coupling the routes to their data sources which are the same in this case.

### To Run The App

You need 'go' installed on your machine and you need to cd into the src folder and run

```
go mod tidy
go run .
```
