# mes-actions-bourse

## Description :

This small app is requesting stock options information from the website finance.yahoo.com

## CLI options :

```shell
  -cac
        Get the values for few index : CAC40, SBF120, Nasdaq, S&P500, Bitcoin, Ethereum

  -sym "symbol list"
        Need to provide a list of symbols separated by a comma, for example BTC-USD,BNP.PA,LI.PA,etc...

  -ver
        Print version and exit

  -xls "xlsx file"
        Need to provide an xlsx file with the stock options to request
        The First row of the file must be column title from this 4: Symbol, Price (PRU), Quantity, Target
        It need at least the column Symbol
```

## Requirements if using the xlsx file :
- Each column much start (at the first cell) with the column name
- The xlsx file much contain at least column "Symbol"
- Other optionals column can be: "Price", "Quantity", "Target"
- Column "Symbol": contains the codes used by yahoo, for exemple Danone in Paris Stock Exchange is BN.PA
- Column "Price": this is the Cost Per Unit (or Prix de Revient Unitaire)
- Column "Quantity": the numbers of stock option owned
- Column "Target": your goal to reach in order to sell

