# Use a base image with Node.js
FROM node:16

# Set the working directory
WORKDIR /app

# Copy the package.json file and install dependencies
COPY package.json /app/
RUN npm install

# Copy the source code of the microservice
COPY . /app/

# Expose the port on which the service will run
EXPOSE 8080

# Command to run the service
CMD ["node", "server.js"]