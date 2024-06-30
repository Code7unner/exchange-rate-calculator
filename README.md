# Exchange Rate Calculator API

## Description

This project is an Exchange Rate Calculator API designed to provide real-time conversion rates between cryptocurrency and fiat currencies.
The API accepts a source and target currency along with an amount, and returns the converted amount using the latest exchange rates.
The exchange rates are updated every minute through a background task.

## Features

- **Real-Time Conversion**: Converts between specified source and target currencies in real-time.
- **Background Rate Updates**: Exchange rates are updated every minute in the background.
- **Database-Driven**: Exchange rates are fetched from a PostgreSQL database.
- **Currency Availability**: The API only processes conversions for available currencies. Unavailable currencies will return an error.
- **Supported Currencies**: EUR, USD, CNY, USDT, USDC, ETH.

## Technical Requirements

- **Language**: Go 1.20 or higher
- **Logging**: Using `zerolog`
- **Data Structures**: Go data structures for storing currency pairs
- **Framework**: Go Fiber for the API
- **Performance**: API response time should not exceed 100ms
- **Exchange Rates API**: Using [FastForex](https://www.fastforex.io/) for obtaining exchange rates
- **Database**: PostgreSQL for storing exchange rates

## Clean Architecture

Project was completed via the [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) practises

![The Clean Architecture](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

Relations between Clean Architecture layers and project folders:

| Template path           | Clean Architecture Layer Color       |
|-------------------------|--------------------------------------|
| internal/entity         | Yellow \(Enterprise Business Rules\) |
| internal/usecase        | Red \(Application Business Rules\)   |
| internal/adapter        | Green \(Interface Adapters\)         |
| internal/infrastructure | Blue \(Frameworks & Drivers\)        |

## Open API Spec

- [api/openapi-spec/httpapi.openapi.yaml](api/openapi-spec/httpapi.openapi.yaml)

## Dependencies

**Infrastructure:**
- Postgres 12
- Fast Forex API

## Start

- Add `FAST_FOREX_API_KEY` to the .env file **(See: .env.example file)**

```shell
make start
```

- Wait for the "INFO Server started on: ${...}" message
```shell
make logs-httpapi
```

```shell
make start-swagger
```

## Update currency pairs
If you want to update currency pairs then update - [internal/infrastructure/config/currency-pairs.json](internal/infrastructure/config/currency-pairs.json) file 

It is an embedded file that processed when app inited