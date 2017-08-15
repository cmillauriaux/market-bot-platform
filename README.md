# Market Bot Platform (MBP)

MBP is an Open Source platform to develop bots capable of buying and selling shares or cryptomonials in major marketplaces.
Any marketplace can be added by implements an Interface. A coinbase implementation is provided, you can make PR to add new implementations.
A training mode is available where your bot received a real history of the marketplace, make orders and you can measure his gain.

## Implement a new marketplace

TODO

## Where to find bitcoin history for a marketplace
You can fin bitcoin history for each marketplace as a CSV format here : `http://api.bitcoincharts.com/v1/csv/`. Use a CSV history allow to load all of history in many seconds, whereas with the API it asks many hours.

## Implement a bot

TODO

## Train a bot

TODO

## Play your robot for real

TODO

## How to profile
go tool pprof http://localhost:6060/debug/pprof/profile

## How to contributes