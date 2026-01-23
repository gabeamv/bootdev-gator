# Guided Boot.dev blog aggregator project.
Gator is an [RSS (Really Simple Syndication)](https://en.wikipedia.org/wiki/RSS) feed aggregator which allows users to add and follow RSS feeds off the internet using the feeds' URLs. Users can control how often they want their feed updated with the latest post from the publisher. Users can also browse the latest posts the feed aggegator receives. To run a command, type in **bootdev-gator** then the desired command along with their arguments if they take any. Run the command **bootdev-gator help** to see commands, descriptions, and what arguments they receive.

**Installation**
1. In order to run this project, install the latest version of [Go](https://go.dev/doc/install).
2. You will also need Postgres. 
    * On mac, run **brew install postgres@15** or another suitable version.
    * On Linux / WSL (Debian), run **sudo apt update**, then **sudo apt install postgresql postgresql-contrib**
3. Ensure you have the postgres by running **psql --version**
    * If you are on Linux, update postgres password: **sudo passwd postgres**
    * Don't forget the password.
4. Enter the psql shell: **psql postgres**
5. Create a new database called '**gator**': **CREATE DATABASE gator;**
6. Set the user password (Linux only): **ALTER USER postgres PASSWORD 'postgres';**
7. Manually create a configuration file in your home directory called **.gatorconfig.json** with the contents:
    * {"db_url":"postgres://{db_username}:{db_password}@localhost:5432/gator?sslmode=disable","current_username":""}
    * replace '{db_username} and {db_password}' with the credentials you set
8. Install the application using: **go install github.com/gabeamv/bootdev-gator**
    * The executable will be installed in your bin go/bin folder







