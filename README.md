# mes-actions-bourse

## Description :

This small app is requesting stock options information from the website finance.yahoo.com

## CLI options :

```shell
  -cac
        Get the values for few selected index and crypto-money: CAC40, SBF120, DAX, EuroStoxx50, Nasdaq, S&P500, Bitcoin and Ethereum

  -sym "symbols names"
        Need to provide a list of symbols separated by a comma, for example BTC-USD,BNP.PA,LI.PA,etc...

  -test
        Use my local URL, for dev/testing purpose

  -ver
        Print version and exit

  -verb
        Print more verbose informations from the application run

  -xls "my file xlsx"
        Need to provide an xlsx file with the stock options to request
        The First row of the file must be column title from this: Symbol, Price (PRU), Quantity, TargetSell, TargetBuy
```

## Requirements if using the xlsx file :
- Each column much start (at the first cell) with the column name
- The xlsx file much contain at least column "Symbol"
- Other optionals column can be: "Price", "Quantity", "TargetSell", "TargetBuy"
- Column "Symbol": contains the codes used by yahoo, for exemple Danone in Paris Stock Exchange is BN.PA
- Column "Price": this is the Cost Per Unit (or Prix de Revient Unitaire)
- Column "Quantity": the numbers of stock option owned
- Column "TargetSell": your goal to reach for selling, equal and higher is highlighting the stock's row
- Column "TargetBuy": your goal to reach for buying, equal and lower is highlighting the stock's row

## License :
GPL3.0

