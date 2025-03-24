This is a hello-world Spring Boot web server generated from https://start.spring.io/.
The server is used to demonstrate memory thrashing causing EC2 EBSByteBalance% deplete. 
Therefore, even though it's a hello-world application, many dependencies are added, 
in the hope of loading jars into JVM. 
See the blog at [Memory thrashing depletes EBSByteBalance](). 

### Run the application
```sh
# run directly. 
./mvnw spring-boot:run -Dspring-boot.run.arguments="--server.port=8081"

# build jar and start the server
./mvnw clean package
java -jar target/HelloWorld-0.0.1-SNAPSHOT.jar --server.port=8081

# Hello
# Find the temporary security password from the startup log. Somethig like this
# "Using generated security password: 9d936074-d039-45d8-ab7e-eb14a2b5cba2"

curl -u hello:dummy http://localhost:8081/
Hello, World!%  
```

### Build OCI image and run in Docker
```sh
./mvnw clean package

docker build -t hello-world-ws . 

docker run -p 8081:8081 --rm hello-world-ws java -jar webserver.jar --server.port=8081

curl -u hello:dummy http://localhost:8081/
Hello, World!%                                                                      
```
