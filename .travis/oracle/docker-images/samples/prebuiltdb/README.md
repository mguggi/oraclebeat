# Example of creating an image with pre-built DB
***Warning: The description below requires changes in the dockerfile related to the version and edition in question (docker-images/OracleDatabase/dockerfiles/$VERSION/Dockerfile.$EDITION). It's recommended that you revert the changes after you have completed creating image without VOLUME (see next section).***

## 1. Create an image without a VOLUME

When creating an image with a pre-build DB, we first need an image without a VOLUME. To create such an image you'll need to comment out the VOLUME command in the dockerfile in use (docker-images/OracleDatabase/dockerfiles/$VERSION/Dockerfile.$EDITION).

Comment out the following line:
```
#VOLUME ["$ORACLE_BASE/oradata"]
```
In the example below I have edited the Dockerfile for version 12.2.0.1 and Enterprise Edition (ee) (docker-images/OracleDatabase/dockerfiles/12.2.0.1/Dockerfile.ee).

If you want to keep your existing image, you need to rename it:
```
docker tag oracle/database:12.2.0.1-ee oracle/database-org:12.2.0.1-ee
```

Then create an image without a VOLUME:
```
cd docker-images/OracleDatabase/dockerfiles
sh buildDockerImage.sh -v 12.2.0.1 -e
```

## 2. Startup a container and create the database

We first need to start up a container to get the database created.
**Note:**  First make sure you have built **oracle/database:12.2.0.1-ee**. 
Now start the container as follow:
```
docker run --name oracle-build -p 1521:1521 -p 5500:5500 oracle/database:12.2.0.1-ee
```

## 3. Reset passwords (optional)

It's recommended to reset passwords before creating the image. This way you don't have to do it everytime you create a new container.

Now you should connect to the container and reset passwords:
```
docker exec oracle-build ./setPassword.sh <newPassword>
```
## 4. Stop the running container

Stop the container (and therefore also the database) before generating your new prebuilt image:
```
docker stop -t 30 oracle-build
```

## 5. Create the image with the prebuilt database

Open a new console on your host and run the following command:
```
docker commit -m "Image with prebuilt database" oracle-build oracle/db-prebuilt:12.2.0.1-ee
```

## 5. Clean up

Remove the temporary container:
```
docker rm oracle-build
```
Remove the image without VOLUME and rename the old image (if exists):
```
docker rmi oracle/database:12.2.0.1-ee
docker tag oracle/database-org:12.2.0.1-ee oracle/database:12.2.0.1-ee
```

Uncomment the volume instruction in your Dockerfile:
```
VOLUME ["$ORACLE_BASE/oradata"]
```
## 6. Ready to use your image with prebuild database

Run your prebuild image:

```
docker run --name <container-name> -p 1521:1521 -p 5500:5500 oracle/db-prebuilt:12.2.0.1-ee
```

After the container is up and running you can connect to the new database.

You can also run your new image from a docker compose.
Create a directory, example "ora12c-db01", and add the following docker-compose.yml file:
```
version: '2'
services:
  orcl-node:
    image: oracle/db-prebuilt:12.2.0.1-ee
    ports:
      - "1521:1521"
      - "5500:5500"
```
And run:
```
docker-compose up
```

# Copyright
Copyright (c) 2014-2017 Oracle and/or its affiliates. All rights reserved.
