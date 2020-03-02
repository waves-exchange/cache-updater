# Updating approach

In order to perform operations with database the code logic is
split between different kinds of controllers. Splitting logic
is one of the ways for making code composing more pure.

# Updating in different cases

### Empty database

If database is empty, whole data from address data is parsed and put
into appropriate tables.

```
func (uc *UpdateController) 
    GrabAllAddressData () ([]byte, error) { ... }
```
 
Note that empty database is the only case
of such approach

Then, grabbed data gets passed to Database controller,
responsible for database operations

```
func (dc *DbController) 
    HandleRecordsUpdate (byteValue []byte) {
```

---

### Database is not empty, valid data is provided up to specific block

```
func (uc *UpdateController) 
    GrabStateChangeData () ([]byte, error) {...}
```

This function is responsible for grabbing appropriate transactions
in specific block range. In our case from "N" to "O" blocks

> N - "last existing block in the database"

> O - "last existing confirmed block in blockchain"

Every entity that is stored in ```address/data``` of the account
conforms to this interface and implements specified methods.
```go
type DAppEntity interface {
  GetKeys(*string) []string
  MapItemToModel(string, map[string]string) *DAppEntity
  UpdateAll(*map[string]string) []DAppEntity
}
```

#### Current refactored approach

Updating starts from checking existing records and if they actually exist
the function delegates responsibility to:
 ```HandleExistingBondsOrdersUpdate```
 method

```go
var existingBondsOrders []entities.BondsOrder
_ = dc.GetAllEntityRecords(&existingBondsOrders, entities.BONDS_ORDERS_NAME)

if len(existingBondsOrders) != 0 {
    dc.HandleExistingBondsOrdersUpdate()
    return
}
```

As it was mentioned before, function firstly checks out how many
blocks need updating, regarding the maximum range 
(99 blocks per request, let it be constant ```MHR = 99```)

```go
maxHeightRange := uint64(99)
heightDiff := bm.Height - latestExRecord.Height
````

if ```heightDiff > 99``` then we should decompose blocks range.

For ```Z``` times, where is
```
Z = (latestHeight - existingHeight) / maxHeightRange
```

The next step is blocks processing. We fetch block headers list:

```go
func FetchBlocksRange (heightMin, heightMax string) *[]models.BlockHeader {}
```

Then, we fetch transaction list for every block, using this function:

```go
func FetchTransactionsOnSpecificBlock (height string) *models.Block {}
```

We face ```2 * Z * MHR``` count of requests

```go
func (uc *UpdateController) UpdateStateChangedData (
	minHeight, maxHeight uint64,
) { ... }	
```

This method checks and updates BondsOrder entities according to
state change.
Also, it checks only "Invoke" transactions (type 16).

#### Approach in filtering

Current method checks if order key is provided, and if so, grabs
another list of keys trying to map to model:
```type BondsOrder struct {} ```