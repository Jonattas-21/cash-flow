# Overview

The main goal of this project it's to provide a cash flow system manager, it's store debit and credit transactions and deliver a daily report, it can be a complete report or a partial report, it's depends if the day it's over or not.

### Activity Diagram
Bellow folow two use cases that illustrate what it's possible to operate in the system.

![alt text](https://raw.githubusercontent.com/Jonattas-21/cash-flow/refs/heads/main/docs/activity_diagram.png "USe Case")

### Container Diagram
The system has been splitted in five distinguiguished services to manage each specific responsability.

![alt text](https://raw.githubusercontent.com/Jonattas-21/cash-flow/refs/heads/main/docs/container_diagram.png "Container Diagram")


#### Transaction-service

> This is a golang service responsible for Handling transactions and stores in DB,

#### Daily-summary-service

> This is a golang service that is in charge for daily report generation.

#### Database

> A Postgres SQL Database, stores all the transactions.

#### KeyCloak

> (third party app) Manage authentication and authorization.

#### Cache

> A Redis Cache service that provide fast report querying.


### How to run
> Install make-> make run
