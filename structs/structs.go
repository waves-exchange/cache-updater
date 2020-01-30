package structs;

type BondsOrder struct {
	Height, Owner, Status, Index, Pairname, Type string
	Timestamp, Price int
	Total, Filledamount, Filledtotal, Resttotal, Amount, Restamount float64
}