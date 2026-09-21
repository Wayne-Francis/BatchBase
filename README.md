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

-A relational SQL database containing mock pharmaceutical data
-Multiple SQL queries for retrieving related batch information
-A simple search function
-Basic filtering
-A very simple user interface
-A basic data visualisation/plot
-Relationships between batch, raw material, manufacturing and QC data

## Technology

BatchBase will use:

-Go
-PostgreSQL
-sqlc
-Goose

## Project Status

Currently in development.

This project is being developed as part of the Boot.dev Capstone Project.


