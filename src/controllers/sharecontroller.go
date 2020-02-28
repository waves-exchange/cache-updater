package controllers

import "github.com/ventuary-lab/cache-updater/src/entities"

type ShareController struct {
	BondsOrder *entities.BondsOrder
	BlocksMap *entities.BlocksMap
}