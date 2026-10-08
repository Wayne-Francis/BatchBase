# BatchBase

BatchBase is my capstone project for Boot.dev. It is a simple, lightweight application aimed at solving one major problem I experienced during my time as a Product Scientist in the pharmaceutical industry: data wrangling across multiple systems to get a clear understanding of batch history.

## The Problem

A major problem I found working as a Product Scientist in the pharmaceutical industry was trying to efficiently and accurately gather all the data on a batch, or a number of batches, to perform an analysis or write a report.

Usually, all the data I needed was stored across multiple systems, each with different ways of pulling or accessing the data and often in different formats.

This was a major part of my job, but it was mostly unseen. A significant amount of time could be spent gathering, organising and connecting the data before I could actually start doing the work that was visible, such as an analysis, presentation or report.

## The Goal

The goal of BatchBase is to create a lightweight SQL database that is flexible and allows the data needed for an investigation to be queried and accessed in one place.

## Planned Features

The planned features for V1 are:

* A relational SQL database containing mock pharmaceutical data
* Multiple SQL queries for retrieving related batch information
* A simple search function
* Basic filtering
* A very simple user interface
* Basic data visualisation/plots
* Relationships between batch, raw material, manufacturing and QC data

## Technology

BatchBase uses:

* Go
* PostgreSQL
* sqlc
* Goose

## Installation

### Prerequisites

BatchBase requires:

* Git
* Go
* PostgreSQL

### 1. Clone the repository

Clone the repository and move into the BatchBase directory:

```bash
git clone https://github.com/Wayne_Francis/BatchBase.git
cd BatchBase
```

### 2. Install Go

Install Go if it is not already installed.

### 3. Install PostgreSQL

BatchBase uses PostgreSQL as its database.

If PostgreSQL is not already installed, download and install PostgreSQL from the official PostgreSQL website.

During installation:

1. Install the PostgreSQL server and command-line tools.
2. When prompted, create a password for the default PostgreSQL user, typically just:

```text
postgres
```

3. Remember the password you choose. It will be required when configuring BatchBase.
4. Leave the PostgreSQL port as the default:

```text
5432
```

### 4. Create the BatchBase Database

BatchBase uses a PostgreSQL database called:

```text
batchbase
```

Open the PostgreSQL command-line tool (`psql`) and connect using the PostgreSQL user:

```bash
psql -U postgres
```

Enter the password you created during the PostgreSQL installation.

Create the BatchBase database:

```sql
CREATE DATABASE batchbase;
```

Then exit `psql`:

```sql
\q
```

### 5. Configure BatchBase

BatchBase stores its local configuration outside the project directory.

Create the following file in your home directory:

```text
~/.BatchBaseconfig.json
```

The file should contain:

```json
{
  "db_url": "postgres://postgres:YOUR_PASSWORD@localhost:5432/batchbase",
  "current_user_name": ""
}
```

Replace `YOUR_PASSWORD` with the PostgreSQL password you created during installation.

For example, if your PostgreSQL password is `mypassword`:

```json
{
  "db_url": "postgres://postgres:mypassword@localhost:5432/batchbase",
  "current_user_name": ""
}
```

The `current_user_name` field can be left empty initially. BatchBase will set the current user when you register or switch users.

**Important:** The `.BatchBaseconfig.json` file contains local configuration and should not be shared.

### 6. Install the Go Dependencies

From the BatchBase directory:

```bash
go mod download
```

### 7. Install Goose

BatchBase uses Goose for database migrations.

Install Goose with:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Check that Goose is installed:

```bash
goose -version
```

### 8. Run the Database Migrations

Move into the schema directory:

```bash
cd sql/schema
```

Run the migrations:

```bash
goose postgres postgres://postgres:YOUR_PASSWORD@localhost:5432/batchbase up
```

Replace `YOUR_PASSWORD` with the same PostgreSQL password used in your BatchBase configuration.

Then return to the BatchBase directory:

```bash
cd ../..
```

### 9. Run BatchBase

From the main BatchBase directory:

```bash
go run .
```

Running BatchBase without a command launches the interactive dashboard.

The dashboard provides access to:

* Users
* Materials
* Manufacturing
* Quality Control
* Batch History
* Plots

## Development Functions

BatchBase contains several development functions used for testing and populating the database. These functions are not included in the dashboard but can be run from the command line at any time.

### Seed

The seed function can be run with:

```bash
go run . seed
```

The seed function first clears the development database and then creates a complete set of mock data.

See the BatchBase Word document for information on the mock drug. It is not based on any real drug combination and is used for demonstration purposes only.

This includes:

* Users — a user is need to fill the `created_by` field, so `seed_user` is created, a new user can be created in the dashboard.
* Materials and material lots
* In-process batches
* Blends
* Fills
* Finished products
* Assemblies
* IPC QC results
* QC release results
* Specifications

The seed data is designed to demonstrate the relationships between materials, manufacturing stages and QC data.

### Reset Functions

Individual reset functions are also available for development and testing.

Because the database contains relationships and foreign key constraints, data should generally be deleted from the end of the manufacturing/QC chain backwards. For example, QC data associated with a batch must be removed before attempting to remove the manufacturing data that it references.

Examples include:

```bash
go run . resetmaterials
go run . resetmaterialusage
go run . resetblends
go run . resetfills
go run . resetfinishedproducts
go run . resetassemblies
go run . resetipcqcresults
go run . resetqcreleaseresults
go run . resetspecs
```

## Project Status

V1 development is complete.

The V1 application contains the database structure, relationships, batch history queries, manufacturing workflow, QC functionality, development seed data, basic data visualisation and an interactive CLI dashboard.

## V2 Development Goals

The next stage of development would focus on moving BatchBase from a CLI-based capstone project towards a more complete application.

### Proper Web Front End

The CLI dashboard could be replaced with a proper web application.

Potential features include:

* Browser-based interface
* Interactive dashboards
* Improved navigation
* Forms for entering and editing data
* Batch searching and filtering
* User authentication and permissions

### Better Data Display and Reporting

The current application focuses on storing and retrieving the data correctly. A V2 would improve how that data is presented for actual analysis.

Potential features include:

* Analysis-friendly CSV file output
* Exportable batch history
* Better filtering and sorting
* Combined manufacturing and QC datasets
* Investigation/report-ready data
* More useful batch summaries

### More Interactive Data Visualisation

The current version contains basic plots demonstrating the QC data.

A V2 could expand this into more interactive data visualisation, including:

* Interactive charts
* Batch-to-batch comparisons
* QC result distributions
* Trend analysis
* Manufacturing and QC data displayed together
* Filtering data directly from charts and dashboards
* More flexible plotting options

### Future Development

Other potential future development could include:

* Batch release workflows
* More advanced QC/specification checks
* Investigation support
* Audit history
* More detailed user permissions
* Production-ready reporting
* Integration with additional data sources

## Notes on AI Usage

The database structure and CRUD functionality were designed by me and were heavily based on the SQL projects from Boot.dev. I borrowed a significant amount of the code structure and patterns I had previously developed during those projects.

AI was used to help troubleshoot errors and, in some cases, clean up or improve the organisation of files.

The seed function was created heavily with the help of AI. I felt that manually creating each piece of mock data and ensuring sufficient variation was not a valuable use of my time. I provided the basic design and asked that the existing CRUD structure be maintained when creating the seed data.

In addition, some of the output formatting, such as the batch history output, was formatted with the help of AI. I did not want to spend a significant amount of time on CLI formatting when the underlying functionality was the more important part of the project.

Finally, in the interest of completing the project and making V1 demonstrable, most of the dashboard functionality was developed with significant assistance from AI. I had reached the point where I wanted to finish the project with a functional, demonstrable interface rather than spend additional time developing the CLI presentation layer.

The underlying project design, database structure, relationships, CRUD functionality and overall direction were my own work, with AI being used primarily as a development and troubleshooting tool.
