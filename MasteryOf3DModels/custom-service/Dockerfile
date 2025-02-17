# Use the official Ruby base image
FROM ruby:3.3.7

# Install necessary dependencies
RUN apt-get update -qq && apt-get install -y nodejs postgresql-client

# Set the working directory for the app
WORKDIR /app

# Copy the Gemfile and Gemfile.lock to install the Ruby dependencies
COPY Gemfile /app/Gemfile
COPY Gemfile.lock /app/Gemfile.lock

# Install Ruby gems
RUN bundle install

# Copy the rest of the application files
COPY . /app

# Expose the port the app will run on
EXPOSE 5007

# Command to run the application
CMD ["ruby", "app.rb"]