# Base image for Python
FROM python:3.9-slim

# Instal dependencies
RUN apt-get update && \
    apt-get install -y libpq-dev gcc && \
    rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the requirements file and install the dependencies
COPY requirements.txt /app/
RUN pip install --no-cache-dir -r requirements.txt

# Copy the source code of the microservice
COPY . /app/

# Expose the port on which the service will run
EXPOSE 5010

# Command to run the service
CMD ["python", "app.py"]