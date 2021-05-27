# OPTIMONT DATA TRANSFER SERVICE

## OBECNE INFORMACE

- program vyhotoven jako sluzba pro windows a linux, jako docker image
- program bezi ve smycce 1 minuta
- podminka pro prenos zakazek ze `zapsi2.terminal_input_order` do `zapsi2.fis_production`: zaznam ma
  v `zapsi2.terminal_input_order.DTS` mladsi jak 1440 minut od aktualniho casu
- podminka pro prenos prostoju z `zapsi2.terminal_input_idle` do `zapsi2.fis_production`: zaznam ma
  v `zapsi2.terminal_input_idle.DTE` mladsi jak 1440 minut od aktualniho casu
- podminka pro prenos stavu vypnuto z `zapsi2.workplace_state` do `zapsi2.fis_production`: zaznam ma
  v `zapsi2.workplace.state.DTE` mladsi jak 1440 minut od aktualniho casu
- zmeny v tabulce zapsi2.fis_production oproti aktualnimu stavu
    - nove sloupce: `ZapsiId`, `IFS`, `Stav`, `Takt`, `Prostoj`, `TypProstoje`, `Chyba`
    - odstranene sloupce: `TerminalInputOrderId`, `IDVC`, `IDOper`
      <br><br>

## IMPORT UZIVATELU DO ZAPSI

- vstupni tabulka: `zapsi2.fis_user`
- vystupni tabulka: `zapsi2.user`
- parovaci vazba: `zapsi2.fis_user.IDZ - zapsi2.user.Login`

| READ zapsi2.fis_user   | WRITE zapsi2.user     | 
| ---------------------- | --------------------- | 
| IDZ                    | Login                 | 
| Jmeno                  | FirstName             | 
| Prijmeni               | Surname               | 
| Rfid                   | Rfid                  | 

1. Pokud neexistuje uzivatel (nenalezena vazba IDZ-Login), je vytvoren novy uzivatel v zapsi2.user:
    - Login
    - FirstName
    - Surname
    - Rfid
    - UserRoleId=2
    - ostatni null
2. Pokud existuje uzivatel (nalezena vazba Login-IDZ), je aktualizovano v zapsi2.user:
    - pouze Rfid

## IMPORT PRODUKTU DO ZAPSI

- vstupni tabulka: `zapsi2.fis_product`
- vystupni tabulka: `zapsi2.product`
- parovaci vazba: `zapsi2.fis_product.ArtNr - zapsi2.product.Barcode`

| READ zapsi2.fis_product | WRITE zapsi2.product  | 
| ----------------------- | --------------------- |
| ArtNr                   | Barcode               | 
| Nazev + Velikost        | Name                  |

1. Pokud neexistuje produkt (nenalezena vazba ArtNr-Barcode), je vytvoren novy produkt v zapsi2.product:
    - Name
    - Barcode
    - ostatni null nebo DEFAULT
2. Pokud existuje produkt (nalezena vazba ArtNr-Barcode), je aktualizovano v zapsi2.product:
    - pouze Name

## IMPORT ZAKAZEK DO ZAPSI

- vstupni tabulka: `zapsi2.fis_order`
- vystupni tabulka: `zapsi2.order`
- pomocna tabulka: `zapsi2.product`
- parovaci vazba: `zapsi2.fis_order.ID - zapsi2.order.Barcode`

| READ zapsi2.fis_order                                                                                      | WRITE zapsi2.order | 
| ---------------------------------------------------------------------------------------------------------- | ------------------ |
| ID                                                                                                         | Barcode            |  
| IDVC                                                                                                       | Name               | 
| IDVM == zapsi2.fis_product_IDVM -> zapsi2.fis_product.ArtNr == zapsi2.product.Barcode -> zapsi2.product.ID | ProductID          | 
| Mnozstvi                                                                                                   | CountRequested     | 

1. Pokud neexistuje zakazka (nenalezena vazba ID-Barcode), je vytvorena nova zakazka v zapsi2.order:
    - Name
    - Barcode
    - ProductId
    - CountRequested
    - ostatni null nebo DEFAULT
2. Pokud existuje zakazka (nalezena vazba ID-Barcode), je aktualizovano v zapsi2.order:
    - Name
    - CountRequested

## EXPORT ZAKAZEK ZE ZAPSI

- vstupni tabulka: `zapsi2.terminal_input_order`
- vystupni tabulka: `zapsi2.fis_production`
- pomocne tabulky: `zapsi2.user`, `zapsi2.order`, `zapsi2.workplace`

| READ zapsi2.terminal_input_order                 | WRITE zapsi2.fis_production | 
| ------------------------------------------------ | --------------------------- |
| OID                                              | ZapsiId                     |  
| DTS                                              | DatumCasOd                  | 
| DTE                                              | DatumCasDo                  | 
| UserID == user.ID -> user.Login                  | ***IDZ                      | 
| DeviceID == workplace.DeviceId -> workplace.Code | IFS                         | 
| OrderID == order.ID -> order.Barcode             | ****IDFis                   | 
| Count-Fail                                       | MnozstviOK                  | 
| Fail                                             | MnozstviNOK                 | 
| Note                                             | KgOK                        | 
| null                                             | KgNOK                       | 
| "a"                                              | Stav                        | 
| AverageCycle                                     | Takt                        | 
| null                                             | Prostoj                     |
| null                                             | TypProstoje                 |
| null                                             | Prenos                      | 
| Pomocny retezec                                  | *****Chyba                  | 

*** pokud neni prirazeno userId pro terminalInputOrder, vlozi se null<br>
*** pokud je prirazeno userId pro terminalInputOrder, ale jeho `zapsi2.user.Login` neni nalezeno va tabulce
jako `zapsi2.fis_user.IDZ`, vlozi se 0 a pomocny retezec se aktualizuje o hodnotu `zapsi2.user.Login`<br>
**** pokud neni prirazeno orderID pro terminalInputOrder, vlozi se null<br>
**** pokud je prirazeno orderID pro terminalInputOrder, ale jeho `zapsi2.order.Barcode` neni nalezeno v tabulce
jako `zapsi2.fis_order.ID`, vlozi se 0 a pomocny retezec se aktualizuje o hodnotu `zapsi2.order.Barcode`<br>
***** pomocny retezec, je vlozeno cokoliv ve formatu `zapsi2.order.Barcode,zapsi2.user.Login` , pokud je zjisteno pri
kontrole

## EXPORT PROSTOJU ZE ZAPSI

- vstupni tabulka: `zapsi2.terminal_input_idle`
- vystupni tabulka: `zapsi2.fis_production`
- pomocne tabulky: `zapsi2.user`, `zapsi2.idle`, `zapsi2.idle_type`,`terminal_input_order_terminal_input_idle`
  , `zapsi2.workplace`

| READ zapsi2.terminal_input_idle                                                                                                                                                                                              | WRITE zapsi2.fis_production | 
| ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------- |
| OID                                                                                                                                                                                                                          | ZapsiId                     |  
| DTS                                                                                                                                                                                                                          | DatumCasOd                  | 
| DTE                                                                                                                                                                                                                          | DatumCasDo                  | 
| UserID == user.ID -> user.Login                                                                                                                                                                                              | ***IDZ                      | 
| DeviceID == workplace.DeviceId -> workplace.Code                                                                                                                                                                             | IFS                         | 
| ID == terminal_input_order_terminal_input_idle.TerminalInputIdleID -> terminal_input_order_terminal_input_idle.TerminalInputOrderID == terminal_input_order.ID ->  terminal_input_order.OrderID == Order.ID -> Order.Barcode | ****IDFis                   | 
| null                                                                                                                                                                                                                         | MnozstviOK                  | 
| null                                                                                                                                                                                                                         | MnozstviNOK                 | 
| null                                                                                                                                                                                                                         | KgOK                        | 
| null                                                                                                                                                                                                                         | KgNOK                       | 
| "p"                                                                                                                                                                                                                          | Stav                        | 
| null                                                                                                                                                                                                                         | Takt                        | 
| IdleId == idle.ID -> idle.Name                                                                                                                                                                                               | Prostoj                     |
| IdleId == idle.ID -> idle.IdleTypeID == idle_type.ID -> idle_type.Name                                                                                                                                                       | TypProstoje                 |
| null                                                                                                                                                                                                                         | Prenos                      | 
| Pomocny retezec                                                                                                                                                                                                              | *****Chyba                  | 

*** pokud neni prirazeno userId pro terminalInputOrder, vlozi se null<br>
*** pokud je prirazeno userId pro terminalInputOrder, ale jeho `zapsi2.user.Login` neni nalezeno va tabulce
jako `zapsi2.fis_user.IDZ`, vlozi se 0 a pomocny retezec se aktualizuje o hodnotu `zapsi2.user.Login`<br>
**** pokud neni prirazeno orderID pro terminalInputOrder, vlozi se null<br>
**** pokud je prirazeno orderID pro terminalInputOrder, ale jeho `zapsi2.order.Barcode` neni nalezeno v tabulce
jako `zapsi2.fis_order.ID`, vlozi se 0 a pomocny retezec se aktualizuje o hodnotu `zapsi2.order.Barcode`<br>
***** pomocny retezec, je vlozeno cokoliv ve formatu `zapsi2.order.Barcode,zapsi2.user.Login` , pokud je zjisteno pri
kontrole

## EXPORT STAVU VYPNUTO ZE ZAPSI

- vstupni tabulka: `zapsi2.workplace_state`
- vystupni tabulka: `zapsi2.fis_production`
- pomocne tabulky: `zapsi2.user`, `zapsi2.order`, `zapsi2.device`, `zapsi2.workplace`

| READ zapsi2.workplace_state                      | WRITE zapsi2.fis_production | 
| ------------------------------------------------ | --------------------------- |
| OID                                              | ZapsiId                     |  
| DTS                                              | DatumCasOd                  | 
| DTE                                              | DatumCasDo                  | 
| null                                             | IDZ                         | 
| WorkplaceId == workplace.ID -> workplace.Code    | IFS                         | 
| null                                             | IDFis                       | 
| null                                             | MnozstviOK                  | 
| null                                             | MnozstviNOK                 | 
| null                                             | KgOK                        | 
| null                                             | KgNOK                       | 
| "v"                                              | Stav                        | 
| null                                             | Takt                        | 
| null                                             | Prostoj                     |
| null                                             | TypProstoje                 |
| null                                             | Prenos                      | 
| null                                             | Chyba                       |