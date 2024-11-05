# Overview

The main goal of this project it's to provide a cash flow system manager, it's store debit and credit transactions and deliver a daily report, it can be a complete report or a partial report, it's depends if the day it's over or not.

The system has been splitted in five distinguiguished services to manage each specific responsability.

#### Transaction-service

> This is a golang service responsible for Handling transactions and stores in DB

#### Daily-summary-service

> This is a golang service that is in charge for daily report generation.

#### Database

> A Postgres SQL Database, stores all the transactions.

#### KeyCloak

> (third party app) Manage authentication and authorization.

#### Cache

> A Redis Cache service that provide fast report querying.

