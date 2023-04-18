UserFinder
==================
This is a microservice design to create and get users and their passwords from a database in planetscale

Building the Docker Image
-------------------------

To build the Docker image, run the following command:


`docker build -t my-microservice .`

This will create a Docker image with the tag `my-microservice`.

Running the Docker Container
----------------------------

To run the Docker container, use the following command:


`docker run -p 3030:3030 -e DSN=<insert DSN here> -e SECRET_KEY=<insert secret key here> my-microservice`

This will start the container and map port 3030 from the container to port 3030 on your host machine. The `-e` flag is used to set the `DSN` and `SECRET_KEY` environment variables.


Contributing
------------

Please submit bug reports and feature requests to the GitHub issue tracker.

License
-------

This project is licensed under the MIT License. 