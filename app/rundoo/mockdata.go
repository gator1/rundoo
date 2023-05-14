package rundoo

func init() {
	products = []Product{
		{
			Name:     "Central Winery	Chardonnay",
			Category: CategoryWine,
			Sku:      "CW21001",
		},
		{
			Name:      "Central Winery Pinot Grigio	1",
			Category:  CategoryWine,
			Sku:       "CW22001",
		},
		{
			Name:      "Central Winery Pinot Grigio	2",
			Category:  CategoryWine,
			Sku:       "CW22002",
		},
		{
			Name:      "Eastern Vines Chardonnay 1",
			Category:  CategoryWine,
			Sku:       "EV21001",
		},
		{
			Name: "Eastern Vines Chardonnay 2",
			Category:  CategoryWine,
			Sku:       "EV21002",
		},
		{
			Name: "Eastern Vines Sauvignon Blanc",
			Category:  CategoryWine,
			Sku:       "EV23001",
		},
		{
			Name: "gRPC Go for Professionals",
			Category:  CategoryBook,
			Sku:       "PP1837638845",
		},
		{
			Name: "Electric Chainsaw",
			Category:  CategoryTool,
			Sku:       "ET1838845",
		},
	}
}
