package rundoo

import (
	"app/internal/data"
)


func init() {
	products = []data.Product{
		{
			Name:     "Central Winery	Chardonnay",
			Category: data.CategoryWine,
			Sku:      "CW21001",
		},
		{
			Name:      "Central Winery Pinot Grigio	1",
			Category:  data.CategoryWine,
			Sku:       "CW22001",
		},
		{
			Name:      "Central Winery Pinot Grigio	2",
			Category:  data.CategoryWine,
			Sku:       "CW22002",
		},
		{
			Name:      "Eastern Vines Chardonnay 1",
			Category:  data.CategoryWine,
			Sku:       "EV21001",
		},
		{
			Name: "Eastern Vines Chardonnay 2",
			Category:  data.CategoryWine,
			Sku:       "EV21002",
		},
		{
			Name: "Eastern Vines Sauvignon Blanc",
			Category:  data.CategoryWine,
			Sku:       "EV23001",
		},
		{
			Name: "gRPC Go for Professionals",
			Category:  data.CategoryBook,
			Sku:       "PP1837638845",
		},
		{
			Name: "Electric Chainsaw",
			Category:  data.CategoryTool,
			Sku:       "ET1838845",
		},
	}
}
