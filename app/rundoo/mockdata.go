package rundoo

import (
	"app/internal/data"
)


func init() {
	products = []data.Product{
		{
			ID: 1,
			Name:     "Central Winery	Chardonnay",
			Category: data.CategoryWine,
			Sku:      "CW21001",
		},
		{
			ID: 2,
			Name:      "Central Winery Pinot Grigio	1",
			Category:  data.CategoryWine,
			Sku:       "CW22001",
		},
		{
			ID: 3,
			Name:      "Central Winery Pinot Grigio	2",
			Category:  data.CategoryWine,
			Sku:       "CW22002",
		},
		{
			ID: 4,
			Name:      "Eastern Vines Chardonnay 1",
			Category:  data.CategoryWine,
			Sku:       "EV21001",
		},
		{
			ID: 5,
			Name: "Eastern Vines Chardonnay 2",
			Category:  data.CategoryWine,
			Sku:       "EV21002",
		},
		{
			ID: 6,
			Name: "Eastern Vines Sauvignon Blanc",
			Category:  data.CategoryWine,
			Sku:       "EV23001",
		},
		{
			ID: 7,
			Name: "gRPC Go for Professionals",
			Category:  data.CategoryBook,
			Sku:       "PP1837638845",
		},
		{
			ID: 8,
			Name: "Electric Chainsaw",
			Category:  data.CategoryTool,
			Sku:       "ET1838845",
		},
	}
}
