package main

import (
	"context"
	"log"
	"time"

	"github.com/khusa-mahal/backend/internal/config"
	"github.com/khusa-mahal/backend/internal/models"
	"github.com/khusa-mahal/backend/internal/repository/elasticsearch"
	"github.com/khusa-mahal/backend/internal/repository/mongodb"
)

func main() {
	log.Println("🌱 Starting data seeder...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to MongoDB
	db, err := mongodb.Connect(cfg)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer db.Close(context.Background())

	log.Println("✅ Connected to MongoDB")

	// Initialize repositories
	productRepo := mongodb.NewProductRepository(db.GetDB())

	// Create indexes
	if err := productRepo.CreateIndexes(context.Background()); err != nil {
		log.Println("⚠️  Failed to create indexes:", err)
	} else {
		log.Println("✅ MongoDB indexes created")
	}

	// Seed products
	products := getInitialProducts()
	ctx := context.Background()

	log.Printf("📦 Seeding %d products...\n", len(products))

	for i, product := range products {
		// Check if product exists by name? Or just insert?
		// Repo Create usually generates ID if missing.
		// Detailed logic: ideally upsert, but for seeding fresh is fine.
		if err := productRepo.Create(ctx, &product); err != nil {
			log.Printf("❌ Failed to create product %s: %v\n", product.Name, err)
			continue
		}
		log.Printf("✓ Created: %s (ID: %s)\n", product.Name, product.ID.Hex())

		// Update the products slice with the generated ID
		products[i].ID = product.ID
	}

	log.Println("✅ Products seeded successfully!")

	// Initialize Elasticsearch and index products
	searchService, err := elasticsearch.NewSearchService(cfg)
	if err != nil {
		log.Println("⚠️  Elasticsearch not available, skipping indexing:", err)
		return
	}

	if err := searchService.CreateIndex(ctx); err != nil {
		log.Println("⚠️  Failed to create Elasticsearch index:", err)
	}

	log.Println("📊 Indexing products in Elasticsearch...")
	for _, product := range products {
		if err := searchService.IndexProduct(ctx, &product); err != nil {
			log.Printf("⚠️  Failed to index product %s: %v\n", product.Name, err)
		}
	}

	log.Println("✅ Elasticsearch indexing complete!")
	log.Println("🎉 Seeding completed successfully!")
}

func mustParseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Now()
	}
	return t
}

func ptr(f float64) *float64 {
	return &f
}

func getInitialProducts() []models.Product {
	rating4_8 := 4.8
	reviews24 := 24

	return []models.Product{
		{
			Name:          "Golden Motia Khusa",
			Category:      "Formal",
			Price:         3500,
			OriginalPrice: ptr(4200),
			Discount:      ptr(17),
			Rating:        &rating4_8,
			Reviews:       &reviews24,
			Image:         "https://lh3.googleusercontent.com/aida-public/AB6AXuDowxQhauI8H2RpKRUEpNdR5W5aXgjEG_8m4u-8cjOCW-xvOLQ72MUVdIhnKxsoQffFh2wxdT7s1QQek7pNG9i8zupSCOa-8QnBQRk2zSGF0K84MwXyvk1ekswK0l0uxJ0PjvnG5ILRNJ_5FErTpsS6KS4glnSjznHxygq8YhmRd1L_IuvtgWKk1wo77KJbaia8ABWYpa1TH1pjUuzWivGJc7aVsAf9SVUDp3wZp_dy_nJ6Ra6wXggsXDTayuxEfrg-lW2Ok2UU4zg",
			Images: []string{
				"https://lh3.googleusercontent.com/aida-public/AB6AXuDowxQhauI8H2RpKRUEpNdR5W5aXgjEG_8m4u-8cjOCW-xvOLQ72MUVdIhnKxsoQffFh2wxdT7s1QQek7pNG9i8zupSCOa-8QnBQRk2zSGF0K84MwXyvk1ekswK0l0uxJ0PjvnG5ILRNJ_5FErTpsS6KS4glnSjznHxygq8YhmRd1L_IuvtgWKk1wo77KJbaia8ABWYpa1TH1pjUuzWivGJc7aVsAf9SVUDp3wZp_dy_nJ6Ra6wXggsXDTayuxEfrg-lW2Ok2UU4zg",
				"https://lh3.googleusercontent.com/aida-public/AB6AXuDUkiPIqnQ3sEoB6GDJB4N0agBF05-CGIB51uCiTzIe_BgJ_zuVwGnDjM5tim1Ih3ROtWpgNVlejmOAU5qDME8vBb2-K-mwNiKjdJPpvVImON5Sq-264Ag0q9-tX_8iR587dU_X9q9xMjJBjDKciqP-wyX_Tl7GKTiaZkXoQ5reeNJ-AK56Bg11WFXIqay-LOD04oUXUMmcBMXD6LnUYEWLK-SZZtL3QZ2X0TwoIR_H_UV9IktXO52dunlA-8fjvVwnapRaI6mgE5c",
				"https://lh3.googleusercontent.com/aida-public/AB6AXuDtk7yHJ7Y5co36vZKQAuPJ-ZG-xowhELS8AVVK5FG35C0FEZaTVr33R39Akar4p25zo6T8dQcNzNFhKtz-nsftb9WfEC37a0gKqZKTTrSTPojap-2meshJUosd2BXmurOJmhP2ysF2-qjfGIPKojFRCe_nImB1cP49bDyOjW59IhFWcnVzFo8w55gJP8Lw3gXjPukRdto382moEXXKDK_jTfCS26Ftp_hCoS1ZQ_bXm7tgomV6CKG0FiJxqsBI3plLPC_o_3zoJJk",
				"https://lh3.googleusercontent.com/aida-public/AB6AXuD3vVGBm8xiFV4-60vdyml0Pm2BRf7wwTPaY4ewf6eacEGbXdap3K8udc-iG30xXcmnT7e8bm_CHlhRRRjIAIAR33N3GLTRj1PLge2_KBW6ea0JsYirby25j3DJ00W8PAM62IEUW4Kd_u_g6EiKanoncCPUHwZGimhRyydsRMxso40h0HI3DdUJ8fQWVhunSesCt4xIU91pRfviQ_Ep44czt6JZAKqrjFb5wio-gTj9D1kUokekJxlH9sWKluBDFw0BQD6mVUcIv7o",
			},
			Description: "Exquisitely handcrafted with traditional dabka and tilla work on a premium velvet base. Perfect for weddings and festive occasions, featuring a padded sole for extra comfort.",
			IsSale:      true,
			Sizes:       []string{"6", "7", "8", "9", "10"},
			Colors: []models.ColorOption{
				{Name: "Gold", Hex: "#FFD700", ImageURL: "https://lh3.googleusercontent.com/aida-public/AB6AXuDowxQhauI8H2RpKRUEpNdR5W5aXgjEG_8m4u-8cjOCW-xvOLQ72MUVdIhnKxsoQffFh2wxdT7s1QQek7pNG9i8zupSCOa-8QnBQRk2zSGF0K84MwXyvk1ekswK0l0uxJ0PjvnG5ILRNJ_5FErTpsS6KS4glnSjznHxygq8YhmRd1L_IuvtgWKk1wo77KJbaia8ABWYpa1TH1pjUuzWivGJc7aVsAf9SVUDp3wZp_dy_nJ6Ra6wXggsXDTayuxEfrg-lW2Ok2UU4zg"},
				{Name: "Maroon", Hex: "#800000"},
			},
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuD9ZCRIMPpTH3nwsXJehNEp9V30M9Z_WjEzQRRZZVPTpJwD5gLCnVfc8yYcwScCAIq7rhmThMgZvI8DETB3k005I2UR5DiINKsn0eZNJkY_acuPXkjD-qi9wuAjrnVFuLEfuQw1AlThsi9w4P_yw9-uc07Lye3P9w29NuEg3-VJV5A5A92bBKHRsIVJ_37WpCDE-p_IDb5tspsftph8QP_z5RBTLfpBrwSbmUDg_E_-o5M8mIwL1RSdXpXqF7XHZY1vx4-Hmdzzwaw",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuDwtgrbW5Qzn0yHgoq_cQFau5C4CtiVkV2o1q21ypH5WvJdaD8S4SxzVElKcB3qegMBVqcckUdqXgSxzfIz_yoyBy81tOiBMBvq8OlNx-khBTyd_rU_hE2tC2zwtNrTCROuPvzPxy-Qu83HUKYWRdaclP9WEaUTm1nqcQJRxBJWKfGFHMDJ7LToTquAj5vsqxRPbesAiUnj_koa0f9FG1VVdZDZqhTAbVCMRnVKY3JWv2cD42LPptgyB4V94FBR7VdvUoHWv9mJEjo",
			CreatedAt:        mustParseTime("2024-11-20T10:00:00Z"),
			Stock:            25,
		},
		{
			Name:        "Gul-e-Lala",
			Category:    "Casual",
			Price:       2500,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuBqp3j1fln4JMEWno5HYANcnFb-86pMs1gDs8Mjw4XQBMWLyuw2ieGAFXIukKoqrLCWZAdneekJEl1jacE5giuiGTrhxO2bjFgkFzpuZzJGHC31OJnr8RrjTWjwDDdYztIO8hAV-tRo7EspCCyljkEngy3bGnN6UDMkFBtQZCB3zCY6sfvJTwO420w82bzjss8LsDg0Req0lAdU0QgWqSVMuT6WsDdX64UHtjgXWddZl5sPKnIladgrr2gIsDrhomAyic4ZhI2la9o",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuBqp3j1fln4JMEWno5HYANcnFb-86pMs1gDs8Mjw4XQBMWLyuw2ieGAFXIukKoqrLCWZAdneekJEl1jacE5giuiGTrhxO2bjFgkFzpuZzJGHC31OJnr8RrjTWjwDDdYztIO8hAV-tRo7EspCCyljkEngy3bGnN6UDMkFBtQZCB3zCY6sfvJTwO420w82bzjss8LsDg0Req0lAdU0QgWqSVMuT6WsDdX64UHtjgXWddZl5sPKnIladgrr2gIsDrhomAyic4ZhI2la9o"},
			IsNew:       false,
			Description: "Light and comfortable casual khusa with floral embroidery. Perfect for everyday wear.",
			Sizes:       []string{"6", "7", "8", "9"},
			Colors: []models.ColorOption{
				{Name: "Pink", Hex: "#FFC0CB"},
				{Name: "White", Hex: "#FFFFFF"},
			},
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuBqp3j1fln4JMEWno5HYANcnFb-86pMs1gDs8Mjw4XQBMWLyuw2ieGAFXIukKoqrLCWZAdneekJEl1jacE5giuiGTrhxO2bjFgkFzpuZzJGHC31OJnr8RrjTWjwDDdYztIO8hAV-tRo7EspCCyljkEngy3bGnN6UDMkFBtQZCB3zCY6sfvJTwO420w82bzjss8LsDg0Req0lAdU0QgWqSVMuT6WsDdX64UHtjgXWddZl5sPKnIladgrr2gIsDrhomAyic4ZhI2la9o",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuAyAyi28NihhKHAYkGNyfdKUmhqt1dQNTxNpLElsmDRrF8wFkpSY7w_cYPH-4fPKLujQBpMlInQmDC0BsgpBbf9es8twytL4jyeWb6NG79CIhP9hVlKfcfgcyOgOLa8tz1dqsg7OKJLr7H7Ygz_AMPKKcDd9KZJetXMDYIvvWcVXU9AQG2P7mzCxonsrv6WhTNm3bYzW1KfcmqRaeU-nlgLXM0SwaaQyg7jt12ZekflGxyWOACgVSiMtKC7za3aq5gJs0cfhmBqBtw",
			CreatedAt:        mustParseTime("2024-10-15T10:00:00Z"),
			Stock:            30,
		},
		{
			Name:        "Midnight Gold",
			Category:    "Formal",
			Price:       3200,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuCdJsk6X6eJcdicOa5ZGgIjUaN95koq_uJC2twqk3ICXjKdLzNy868qD_fOso8wMW9jwWt7-94icpPK7ys1PiN7TwryN12eOPRrq3tY57mM9mcfyrj0f3YbvNMjoJcs61ZzklG_VzcsgjmY34YGUhy4YOBdFC5FnDEp1BgI6236v-7Onvst84aAH5XvncwL417rSJB0EatiO2vdKZNnnEhja8ltCMQ6cl3LceqtJ63PreWpe-HI1IZgVv0Qlln1Rg72LLQKHMs979I",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuCdJsk6X6eJcdicOa5ZGgIjUaN95koq_uJC2twqk3ICXjKdLzNy868qD_fOso8wMW9jwWt7-94icpPK7ys1PiN7TwryN12eOPRrq3tY57mM9mcfyrj0f3YbvNMjoJcs61ZzklG_VzcsgjmY34YGUhy4YOBdFC5FnDEp1BgI6236v-7Onvst84aAH5XvncwL417rSJB0EatiO2vdKZNnnEhja8ltCMQ6cl3LceqtJ63PreWpe-HI1IZgVv0Qlln1Rg72LLQKHMs979I"},
			IsNew:       true,
			Description: "Elegant black velvet khusa with gold tilla work. Sophisticated and timeless.",
			Sizes:       []string{"6", "7", "8", "9", "10", "11"},
			Colors: []models.ColorOption{
				{Name: "Black Gold", Hex: "#000000"},
			},
			CreatedAt:        mustParseTime("2024-12-01T10:00:00Z"),
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuCdJsk6X6eJcdicOa5ZGgIjUaN95koq_uJC2twqk3ICXjKdLzNy868qD_fOso8wMW9jwWt7-94icpPK7ys1PiN7TwryN12eOPRrq3tY57mM9mcfyrj0f3YbvNMjoJcs61ZzklG_VzcsgjmY34YGUhy4YOBdFC5FnDEp1BgI6236v-7Onvst84aAH5XvncwL417rSJB0EatiO2vdKZNnnEhja8ltCMQ6cl3LceqtJ63PreWpe-HI1IZgVv0Qlln1Rg72LLQKHMs979I",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuCdJsk6X6eJcdicOa5ZGgIjUaN95koq_uJC2twqk3ICXjKdLzNy868qD_fOso8wMW9jwWt7-94icpPK7ys1PiN7TwryN12eOPRrq3tY57mM9mcfyrj0f3YbvNMjoJcs61ZzklG_VzcsgjmY34YGUhy4YOBdFC5FnDEp1BgI6236v-7Onvst84aAH5XvncwL417rSJB0EatiO2vdKZNnnEhja8ltCMQ6cl3LceqtJ63PreWpe-HI1IZgVv0Qlln1Rg72LLQKHMs979I",
			Stock:            20,
		},
		{
			Name:        "Azure Breeze",
			Category:    "Casual",
			Price:       2100,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuAYaZx2ZQzMtv70bVUU98SUdJNO63nMyE1pY5DItDBqsPUbsdWmZeKmvN-oEdTlEhWN_W_e5pHOwwdILaQXvaA63M8qsFhlZIx4SqYGIutBnqtXM8r0xpAHc_V-SlBfE3aywXeiJr0KLWdIxcLVZ5D7-IcrHAnzzOG9oPx5x3s8JUKOKn6TZ-1D49xe4VLJhD2bIAsLtp8QzJ5k7K1p9lymBEUGgsHqdxeq17C2h8Tt4uceacnQ-WC__RJa6eG48AM-rh-xtKfopE4",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuAYaZx2ZQzMtv70bVUU98SUdJNO63nMyE1pY5DItDBqsPUbsdWmZeKmvN-oEdTlEhWN_W_e5pHOwwdILaQXvaA63M8qsFhlZIx4SqYGIutBnqtXM8r0xpAHc_V-SlBfE3aywXeiJr0KLWdIxcLVZ5D7-IcrHAnzzOG9oPx5x3s8JUKOKn6TZ-1D49xe4VLJhD2bIAsLtp8QzJ5k7K1p9lymBEUGgsHqdxeq17C2h8Tt4uceacnQ-WC__RJa6eG48AM-rh-xtKfopE4"},
			Description: "Soft leather khusa in vibrant blue. Comfortable for all-day wear.",
			Sizes:       []string{"7", "8", "9", "10"},
			Colors: []models.ColorOption{
				{Name: "Azure Blue", Hex: "#007FFF"},
				{Name: "Sky Blue", Hex: "#87CEEB"},
			},
			CreatedAt:        mustParseTime("2024-11-28T10:00:00Z"),
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuAYaZx2ZQzMtv70bVUU98SUdJNO63nMyE1pY5DItDBqsPUbsdWmZeKmvN-oEdTlEhWN_W_e5pHOwwdILaQXvaA63M8qsFhlZIx4SqYGIutBnqtXM8r0xpAHc_V-SlBfE3aywXeiJr0KLWdIxcLVZ5D7-IcrHAnzzOG9oPx5x3s8JUKOKn6TZ-1D49xe4VLJhD2bIAsLtp8QzJ5k7K1p9lymBEUGgsHqdxeq17C2h8Tt4uceacnQ-WC__RJa6eG48AM-rh-xtKfopE4",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuAYaZx2ZQzMtv70bVUU98SUdJNO63nMyE1pY5DItDBqsPUbsdWmZeKmvN-oEdTlEhWN_W_e5pHOwwdILaQXvaA63M8qsFhlZIx4SqYGIutBnqtXM8r0xpAHc_V-SlBfE3aywXeiJr0KLWdIxcLVZ5D7-IcrHAnzzOG9oPx5x3s8JUKOKn6TZ-1D49xe4VLJhD2bIAsLtp8QzJ5k7K1p9lymBEUGgsHqdxeq17C2h8Tt4uceacnQ-WC__RJa6eG48AM-rh-xtKfopE4",
			Stock:            35,
		},
		{
			Name:        "Silver Dust",
			Category:    "Beaded Work",
			Price:       2800,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuCbYlYPyjoEfH51y8Y-nYw0fYDUVWwPsAEvoCbOEWbPG8FhZC6UOLQpPsICrS-TNULzjEgDgDumDSEWJSDfJnVeAwgssJZ8oT4fIWmeAXD5MTP7inRldQHCc5Ix45RaJfaH_CoWE6a9q1HAJId4arWZWyMnPOGP1jSQKf77eIcn30G-xC28-r1a9W-dXQmrBwAYZ0yoz80z3a5n-z4YqyLs1eXV5tPLH0Hr8pW5QHGVWIgOZypg-RVy8kR56w4_9CZF2mFUNu0ZgTQ",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuCbYlYPyjoEfH51y8Y-nYw0fYDUVWwPsAEvoCbOEWbPG8FhZC6UOLQpPsICrS-TNULzjEgDgDumDSEWJSDfJnVeAwgssJZ8oT4fIWmeAXD5MTP7inRldQHCc5Ix45RaJfaH_CoWE6a9q1HAJId4arWZWyMnPOGP1jSQKf77eIcn30G-xC28-r1a9W-dXQmrBwAYZ0yoz80z3a5n-z4YqyLs1eXV5tPLH0Hr8pW5QHGVWIgOZypg-RVy8kR56w4_9CZF2mFUNu0ZgTQ"},
			Description: "Intricate beaded work on soft velvet. A true work of art.",
			Sizes:       []string{"6", "7", "8", "9", "10"},
			Colors: []models.ColorOption{
				{Name: "Silver", Hex: "#C0C0C0"},
			},
			CreatedAt: mustParseTime("2024-11-10T10:00:00Z"),
			Stock:     18,
		},
		{
			Name:        "Velvet Gold Clutch",
			Category:    "Accessories",
			Price:       1800,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuDNfrWypNFS8HXi1tsLkFYemLDjKzz4j_WNGLJk9WhaJdlXTUUEL3ifEqH4HOIb2p7Bm9YxTebQ5IVAIN1fsDKiJ20xyZjVtiteziw7v4wVjzmu9kB7ybDIT20MCfjasaYzwNiK0HykeHCZKWCbegbDrnF2pxOzKNw6Zb4BBiDuJQUDb_PgxXnaoIzXgN_afyiQjWnuetXIMKK4aYd4yuj-9zUnlLrgxSrRZjxbc7BoTkCvNR4QrVF9jHrh7C1mHoxraT6ILmYVE9w",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuDNfrWypNFS8HXi1tsLkFYemLDjKzz4j_WNGLJk9WhaJdlXTUUEL3ifEqH4HOIb2p7Bm9YxTebQ5IVAIN1fsDKiJ20xyZjVtiteziw7v4wVjzmu9kB7ybDIT20MCfjasaYzwNiK0HykeHCZKWCbegbDrnF2pxOzKNw6Zb4BBiDuJQUDb_PgxXnaoIzXgN_afyiQjWnuetXIMKK4aYd4yuj-9zUnlLrgxSrRZjxbc7BoTkCvNR4QrVF9jHrh7C1mHoxraT6ILmYVE9w"},
			Description: "Matching clutch with gold embellishments. Perfect accessory for formal occasions.",
			CreatedAt:   mustParseTime("2024-10-20T10:00:00Z"),
			Colors:      []models.ColorOption{{Name: "Gold", Hex: "#FFD700"}},
			Sizes:       []string{"One Size"},
			Stock:       15,
		},
		{
			Name:        "Golden Tilla Khusa",
			Category:    "Bridal",
			Price:       4500,
			IsNew:       true,
			Image:       "https://lh3.googleusercontent.com/aida-public/AB6AXuBQNCFPy7HUp28Z1JDpAG-GDuiUw5659BLn2LD5hGaD83uHhxDAyyDT2CvJF0Usggy90xVSl5CgQIIWX6j1gWQ4EaIKo-eMaBqe9Pa4EoaADmY5TFAAa5RuQ1h9sk-ZFpsB_ZTEgyw0HNdd4f7f4jH4JMB7Jv68cR-oRrNt4sKMvCrPZMMNhflwUPY5NBOIaH-rlz1rjvGUZJvlbQdVRz9k06lF8JmNDt1LEct_EX5Be74qcFqWlA-zXAGpn1oIy2gw5HkCJRpDU3o",
			Images:      []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuBQNCFPy7HUp28Z1JDpAG-GDuiUw5659BLn2LD5hGaD83uHhxDAyyDT2CvJF0Usggy90xVSl5CgQIIWX6j1gWQ4EaIKo-eMaBqe9Pa4EoaADmY5TFAAa5RuQ1h9sk-ZFpsB_ZTEgyw0HNdd4f7f4jH4JMB7Jv68cR-oRrNt4sKMvCrPZMMNhflwUPY5NBOIaH-rlz1rjvGUZJvlbQdVRz9k06lF8JmNDt1LEct_EX5Be74qcFqWlA-zXAGpn1oIy2gw5HkCJRpDU3o"},
			Description: "Premium bridal khusa with heavy golden tilla work. Made for special occasions.",
			Sizes:       []string{"6", "7", "8", "9", "10"},
			Colors: []models.ColorOption{
				{Name: "Royal Gold", Hex: "#FFD700"},
			},
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuBQNCFPy7HUp28Z1JDpAG-GDuiUw5659BLn2LD5hGaD83uHhxDAyyDT2CvJF0Usggy90xVSl5CgQIIWX6j1gWQ4EaIKo-eMaBqe9Pa4EoaADmY5TFAAa5RuQ1h9sk-ZFpsB_ZTEgyw0HNdd4f7f4jH4JMB7Jv68cR-oRrNt4sKMvCrPZMMNhflwUPY5NBOIaH-rlz1rjvGUZJvlbQdVRz9k06lF8JmNDt1LEct_EX5Be74qcFqWlA-zXAGpn1oIy2gw5HkCJRpDU3o",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuBQNCFPy7HUp28Z1JDpAG-GDuiUw5659BLn2LD5hGaD83uHhxDAyyDT2CvJF0Usggy90xVSl5CgQIIWX6j1gWQ4EaIKo-eMaBqe9Pa4EoaADmY5TFAAa5RuQ1h9sk-ZFpsB_ZTEgyw0HNdd4f7f4jH4JMB7Jv68cR-oRrNt4sKMvCrPZMMNhflwUPY5NBOIaH-rlz1rjvGUZJvlbQdVRz9k06lF8JmNDt1LEct_EX5Be74qcFqWlA-zXAGpn1oIy2gw5HkCJRpDU3o",
			CreatedAt:        mustParseTime("2024-12-05T10:00:00Z"),
			Stock:            12,
		},
		{
			Name:          "Scarlet Velvet Dream",
			Category:      "Velvet",
			Price:         3200,
			OriginalPrice: ptr(4000),
			Discount:      ptr(20),
			IsSale:        true,
			Image:         "https://lh3.googleusercontent.com/aida-public/AB6AXuCJZXcDuxIIy0_A_CnSh9LGX_MIDjKV-1TIx1HWLBcOSIy8-gZw81skx82dhDZgM9UdEmK7gR7yMzP2roKHZMBqHqR5oJtNt9tDTwgcHdy2mEzN41ptOQniZPqScBMFlVKNgV94N9NeKyJB7zJzsJVHcMImPtHKh0g3ZBmus8ptWbJoF8GsJU7TvHX0iyyN8N_UoS18yMkpZQcvBto3c7AqE8oViBbbB6Q9ayFSTFvBpFTpbZZow_P52fQrkIqXmR_kgFcsyzn7CpY",
			Images:        []string{"https://lh3.googleusercontent.com/aida-public/AB6AXuCJZXcDuxIIy0_A_CnSh9LGX_MIDjKV-1TIx1HWLBcOSIy8-gZw81skx82dhDZgM9UdEmK7gR7yMzP2roKHZMBqHqR5oJtNt9tDTwgcHdy2mEzN41ptOQniZPqScBMFlVKNgV94N9NeKyJB7zJzsJVHcMImPtHKh0g3ZBmus8ptWbJoF8GsJU7TvHX0iyyN8N_UoS18yMkpZQcvBto3c7AqE8oViBbbB6Q9ayFSTFvBpFTpbZZow_P52fQrkIqXmR_kgFcsyzn7CpY"},
			Description:   "Luxurious red velvet khusa with subtle embroidery. Festive and elegant.",
			Sizes:         []string{"6", "7", "8", "9", "10"},
			Colors: []models.ColorOption{
				{Name: "Scarlet Red", Hex: "#DC143C"},
				{Name: "Deep Maroon", Hex: "#800000"},
			},
			CombinationImage: "https://lh3.googleusercontent.com/aida-public/AB6AXuCJZXcDuxIIy0_A_CnSh9LGX_MIDjKV-1TIx1HWLBcOSIy8-gZw81skx82dhDZgM9UdEmK7gR7yMzP2roKHZMBqHqR5oJtNt9tDTwgcHdy2mEzN41ptOQniZPqScBMFlVKNgV94N9NeKyJB7zJzsJVHcMImPtHKh0g3ZBmus8ptWbJoF8GsJU7TvHX0iyyN8N_UoS18yMkpZQcvBto3c7AqE8oViBbbB6Q9ayFSTFvBpFTpbZZow_P52fQrkIqXmR_kgFcsyzn7CpY",
			WornImage:        "https://lh3.googleusercontent.com/aida-public/AB6AXuCbYlYPyjoEfH51y8Y-nYw0fYDUVWwPsAEvoCbOEWbPG8FhZC6UOLQpPsICrS-TNULzjEgDgDumDSEWJSDfJnVeAwgssJZ8oT4fIWmeAXD5MTP7inRldQHCc5Ix45RaJfaH_CoWE6a9q1HAJId4arWZWyMnPOGP1jSQKf77eIcn30G-xC28-r1a9W-dXQmrBwAYZ0yoz80z3a5n-z4YqyLs1eXV5tPLH0Hr8pW5QHGVWIgOZypg-RVy8kR56w4_9CZF2mFUNu0ZgTQ",
			CreatedAt:        mustParseTime("2024-11-18T10:00:00Z"),
			Stock:            22,
		},
	}
}
