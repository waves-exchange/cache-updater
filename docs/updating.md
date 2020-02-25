# Updating approach

In order to perform operations with database the code logic is
split between different kinds of controllers. Splitting logic
is one of the ways for making code composing more pure.

# Updating in different cases

> Empty database

If database is empty, whole data from address data is parsed and put
into appropriate tables.

```
func (uc *UpdateController) 
    GrabAllAddressData () ([]byte, error) { ... }
```
 
Note that empty database is the only case
of such approach

---

> Database is not empty, valid data is provided
> up to specific block

```
func (uc *UpdateController) 
    GrabStateChangeData () ([]byte, error) {...}
```

This function is responsible for grabbing appropriate transactions
in specific block range. In our case from "N" to "O" blocks

> N - "last existing block in the database"

> O - "last existing confirmed block in blockchain"


The next step is to transform all transactions in blocks,
map to specified entity model, and more importantly, provide
every record with status enum.

``` confirm_status ```

This field is enum, which can take either 

``` confirmed ```
or
``` unconfirmed ```

This approach is the only case to handle forks, in future
it can help to drop invalid records

