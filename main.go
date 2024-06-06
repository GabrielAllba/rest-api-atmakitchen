package main

import (
	adminauthcontroller "backend-atmakitchen/controllers/admincontroller"
	"backend-atmakitchen/controllers/autologincontroller"
	"backend-atmakitchen/controllers/bahancontroller"
	"backend-atmakitchen/controllers/bankcontroller"
	"backend-atmakitchen/controllers/cartcontroller"
	"backend-atmakitchen/controllers/consignationcontroller"
	customerauthcontroller "backend-atmakitchen/controllers/customercontroller"
	"backend-atmakitchen/controllers/emailcontroller"
	"backend-atmakitchen/controllers/hamperscontroller"
	"backend-atmakitchen/controllers/invoicecountercontroller"
	moauthcontrollerauthcontroller "backend-atmakitchen/controllers/mocontroller"
	ownerauthcontroller "backend-atmakitchen/controllers/ownercontroller"
	pembelianBahanBakubakucontroller "backend-atmakitchen/controllers/pembelianbahanbakucontroller"
	"backend-atmakitchen/controllers/pengeluaranlaincontroller"
	"backend-atmakitchen/controllers/quotacontroller"
	"backend-atmakitchen/controllers/transactioncontroller"
	"backend-atmakitchen/controllers/transactiondetailcontroller"

	"backend-atmakitchen/controllers/presensicontroller"
	"backend-atmakitchen/controllers/productcontroller"
	"backend-atmakitchen/controllers/producttypecontroller"
	"backend-atmakitchen/controllers/resepcontroller"
	"backend-atmakitchen/controllers/rolecontroller"
	"backend-atmakitchen/controllers/tokencontroller"
	"backend-atmakitchen/controllers/usercontroller"
	"backend-atmakitchen/initializers"
	"backend-atmakitchen/middleware"
	"backend-atmakitchen/models"
	"path/filepath"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariable()
	models.ConnectDatabase()
}
func getImage(c *gin.Context) {
	filename := c.Param("filename")
	imagePath := filepath.Join("images", filename)
	c.File(imagePath)
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "*")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	// 	c.Header("Access-Control-Allow-Headers", "Origin, Authorization, Content-Type")
	// 	c.Header("Access-Control-Allow-Credentials", "true")
	// 	if c.Request.Method == "OPTIONS" {
	//         c.AbortWithStatus(204)
	//     	return
	// 	}
	// 	c.Next()
	// })

	// login based by role
	autologin := r.Group("/api/autologin")
	{
		autologin.POST("/login", autologincontroller.Login)
		autologin.POST("/logout", autologincontroller.Logout)
	}

	email := r.Group("/api/email")
	{
		email.POST("/send", emailcontroller.SendEmail)
	}

	// product
	product := r.Group("/api/product")
	{
		product.POST("", productcontroller.Create)
		product.GET("", productcontroller.Index)
		product.GET("/:id", productcontroller.Show)
		product.GET("/search", productcontroller.Search)
		product.DELETE("/:id", productcontroller.Delete)
		product.PUT("/:id", productcontroller.Update)
		product.GET("/type", productcontroller.SearchType)
		product.GET("/type/search", productcontroller.SearchProductByType)
		product.GET("/tag/search", productcontroller.SearchProductByTag)
		product.PUT("/stock/:id", productcontroller.UpdateStock)

	}

	// customer
	user := r.Group("/api/customer")
	{
		user.POST("/signup", customerauthcontroller.Signup)
		user.POST("/login", customerauthcontroller.Login)
		user.POST("/logout", customerauthcontroller.Logout)
		user.GET("/token/validate/:tokenString", customerauthcontroller.Validate)
		user.GET("/email-exists", customerauthcontroller.EmailExists)
		user.PUT("/updatepassword/:email", customerauthcontroller.UpdatePassword)
		

	}

	// admin
	admin := r.Group("/api/admin")
	{
		role := admin.Group("/role")
		{
			role.POST("", rolecontroller.Create)
			role.GET("/:id", rolecontroller.Show)
		}
		admin.POST("/login", adminauthcontroller.Login)
		admin.POST("/logout", adminauthcontroller.Logout)
		admin.GET("/token/validate/:tokenString", adminauthcontroller.Validate)
	}

	// mo
	mo := r.Group("/api/mo")
	{
		mo.POST("/login", moauthcontrollerauthcontroller.Login)
		mo.POST("/logout", moauthcontrollerauthcontroller.Logout)
		mo.GET("/token/validate/:tokenString", moauthcontrollerauthcontroller.Validate)
	}

	// owner
	owner := r.Group("/api/owner")
	{
		owner.POST("/login", ownerauthcontroller.Login)
		owner.POST("/logout", ownerauthcontroller.Logout)
		owner.GET("/token/validate/:tokenString", ownerauthcontroller.Validate)
	}

	// product type
	product_type := r.Group("/api/product_type")
	{
		product_type.POST("", producttypecontroller.Create)
		product_type.GET("/", producttypecontroller.Index)
		product_type.GET("/:id", producttypecontroller.Show)
	}

	// consignation
	consignation := r.Group("/api/consignation")
	{
		consignation.POST("", consignationcontroller.Create)
		consignation.GET("", consignationcontroller.Index)
		consignation.GET("/search", consignationcontroller.Search)
		consignation.DELETE("/:id", consignationcontroller.Delete)
		consignation.PUT("/:id", consignationcontroller.Update)

	}

	bank := r.Group("/api/bank")
	{
		bank.POST("", bankcontroller.Create)
		bank.GET("", bankcontroller.Index)

	}

	hampers := r.Group("/api/hampers")
	{
		hampers.GET("", hamperscontroller.Index)
		hampers.GET("/:id", hamperscontroller.Show)
		hampers.PUT("/:id", hamperscontroller.Update)
		hampers.GET("/latest_id", hamperscontroller.GetLatestHampersID)
		hampers.POST("", hamperscontroller.Create)
		hampers.POST("/detail/:id", hamperscontroller.CreateDetail)
		hampers.GET("/search", hamperscontroller.Search)
		hampers.DELETE("/:id", hamperscontroller.Delete)
		hampers.DELETE("/detail/:id", hamperscontroller.DeleteDetailHampers)
		hampers.PUT("/detail/:id", hamperscontroller.UpdateDetail)
		hampers.PUT("/stock/:id", hamperscontroller.UpdateStock)

	}

	//resep
	resep := r.Group("/api/resep")
	{
		resep.POST("", resepcontroller.Create)
		resep.GET("", resepcontroller.Index)
		resep.GET("/:id", resepcontroller.Show)
		resep.GET("/search", resepcontroller.Search)
		resep.DELETE("/:id", resepcontroller.Delete)
		resep.PUT("/:id", resepcontroller.Update)
		resep.GET("/latest_id", resepcontroller.GetLatestResepID)
		resep.POST("/detail/:resep_id", resepcontroller.CreateDetail)
		resep.GET("/detail/:resep_id", resepcontroller.GetDetailResep)
		// resep.GET("/type", resepcontroller.SearchType);
		// resep.GET("/type/search", resepcontroller.SearchProductByType);
	}

	//bahan
	bahan := r.Group("/api/bahan")
	{
		bahan.POST("", bahancontroller.Create)
		bahan.GET("", bahancontroller.Index)
		bahan.GET("/:id", bahancontroller.Show)
		bahan.GET("/search", bahancontroller.Search)
		bahan.DELETE("/:id", bahancontroller.Delete)
		bahan.PUT("/:id", bahancontroller.Update)
	}

	roles := r.Group("/api/roles")
	{
		roles.POST("", rolecontroller.Create)
		roles.GET("", rolecontroller.Index)
		roles.PUT("/:id", rolecontroller.Update)
		roles.GET("/:id", rolecontroller.Show)
		roles.DELETE("/:id", rolecontroller.Delete)
	}

	users := r.Group("/api/users")
	{
		users.GET("", customerauthcontroller.Index)
		users.GET("/cari", usercontroller.Search)
		users.DELETE("/:id", usercontroller.Delete)
		users.PUT("/updateUser/:email", usercontroller.UpdateUser)
		users.GET("/:id", usercontroller.Show)
		users.PUT("/update-points/:id/:points", usercontroller.UpdatePoints)
		users.GET("/customer", usercontroller.SearchType)
		users.GET("/customer/search", usercontroller.SearhUserByType)
	}

	token := r.Group("/api/token")
	{
		token.POST("/create/:user_id", tokencontroller.CreateToken)
		token.GET("", tokencontroller.Index)
		token.DELETE("/:user_id", tokencontroller.DeleteToken)
		token.GET("/check", tokencontroller.CheckToken)
	}

	//pembelian bahan baku
	pembelian_bahan_baku := r.Group("/api/pembelian_bahan_baku")
	{
		pembelian_bahan_baku.POST("", pembelianBahanBakubakucontroller.Create)
		pembelian_bahan_baku.GET("", pembelianBahanBakubakucontroller.Index)
		pembelian_bahan_baku.GET("/search", pembelianBahanBakubakucontroller.GetByDateRange)
		pembelian_bahan_baku.DELETE("/:id", pembelianBahanBakubakucontroller.Delete)
		pembelian_bahan_baku.GET("/:id", pembelianBahanBakubakucontroller.Show)
		pembelian_bahan_baku.PUT("/:id", pembelianBahanBakubakucontroller.Update)

	}

	// cart
	cart := r.Group("/api/cart")
	{
		cart.POST("", cartcontroller.Create)
		cart.GET("", cartcontroller.Index)
		cart.GET("/user/:userId", cartcontroller.GetByUserID)
		cart.GET("/:id", cartcontroller.Show)
		cart.DELETE("/:id", cartcontroller.Delete)
		cart.PUT("/:id", cartcontroller.Update)
	}

	// transactions
	transactions := r.Group("/api/transactions")
	{
		transactions.POST("", transactioncontroller.Create)
		transactions.GET("", transactioncontroller.Index)
		transactions.GET("/tampil/:transaction_status", transactioncontroller.GetTransaksiByStatus)
		transactions.GET("/pengiriman/:transaction_status", transactioncontroller.GetTransaksiByTwoStatus)
		transactions.GET("/batal/:transaction_status", transactioncontroller.GetTransaksiCanceledReadyStok)
		transactions.GET("/user/:userId", transactioncontroller.GetByUserID)
		transactions.GET("/:id", transactioncontroller.Show)
		transactions.DELETE("/:id", transactioncontroller.Delete)
		transactions.PUT("/:id", transactioncontroller.Update)
		transactions.PUT("/status/:id/:transaction_status", transactioncontroller.UpdateStatus)
		transactions.PUT("/transfer_nominal/:id", transactioncontroller.UpdateTotalAfterDeliveryFee)
		transactions.PUT("/bukti_pembayaran/:invoice_number", transactioncontroller.UpdateBuktiPembayaran)
		
		transactions.PUT("/status/invoice/:invoice_number/:transaction_status", transactioncontroller.UpdateStatusByInvoice)
	}

	// transaction_detail
	transaction_details := r.Group("/api/transaction_details")
	{
		transaction_details.POST("", transactiondetailcontroller.Create)
		transaction_details.GET("", transactiondetailcontroller.Index)
		transaction_details.GET("/invoice/:invoiceNumber", transactiondetailcontroller.GetByInvoiceNumber)
		transaction_details.GET("/user/:userId", transactiondetailcontroller.GetByUserID)
		transaction_details.GET("/:id", transactiondetailcontroller.Show)
		transaction_details.DELETE("/:id", transactiondetailcontroller.Delete)
		transaction_details.PUT("/:id", transactiondetailcontroller.Update)
		transaction_details.GET("/invoice/photos/:invoiceNumber", transactiondetailcontroller.GetPhotosByInvoiceNumber)
		transaction_details.GET("/invoice/photos/user/:userId", transactiondetailcontroller.GetPhotosByUserID)
	}

	invoice_number := r.Group("/api/invoice_number")
	{
		invoice_number.POST("", invoicecountercontroller.Create)
	}

	quotaRoutes := r.Group("/api/quota")
	{
		quotaRoutes.POST("/", quotacontroller.Create)
		quotaRoutes.GET("/:id", quotacontroller.Show)
		quotaRoutes.GET("/", quotacontroller.Index)
		quotaRoutes.DELETE("/:id", quotacontroller.Delete)
		quotaRoutes.PUT("/:id", quotacontroller.Update)
		quotaRoutes.GET("/product", quotacontroller.GetByProductAndDate)
		quotaRoutes.GET("/hampers", quotacontroller.GetByHampersAndDate)

		quotaRoutes.PUT("/tanggal", quotacontroller.UpdateQuota)
		quotaRoutes.POST("/tanggal", quotacontroller.InsertQuota)
		quotaRoutes.GET("/product/tanggal/:product_id/:tanggal", quotacontroller.GetQuotaByProductIDAndTanggal)

		
		quotaRoutes.PUT("/hampers/tanggal", quotacontroller.UpdateQuotaHampers)
		quotaRoutes.POST("/hampers/tanggal", quotacontroller.InsertQuotaHampers)
		quotaRoutes.GET("/hampers/tanggal/:hampers_id/:tanggal", quotacontroller.GetQuotaByHampersIDAndTanggal)



	}

	pengeluaran_lain := r.Group("/api/pengeluaran_lain")
	{
		//create, index, show, delete, search, update
		pengeluaran_lain.POST("", pengeluaranlaincontroller.Create)
		pengeluaran_lain.GET("", pengeluaranlaincontroller.Index)
		pengeluaran_lain.GET("/:id", pengeluaranlaincontroller.Show)
		pengeluaran_lain.DELETE("/:id", pengeluaranlaincontroller.Delete)
		pengeluaran_lain.GET("/search", pengeluaranlaincontroller.Search)
		pengeluaran_lain.PUT("/:id", pengeluaranlaincontroller.Update)
	}

	presensi := r.Group("/api/presensi")
	{
		presensi.POST("",presensicontroller.Create)
		presensi.GET("",presensicontroller.Index)
		presensi.PUT("/:id",presensicontroller.Update)
	}


	// Define the route
	r.GET("/api/images/:filename", getImage)

	// check auth middleware
	r.GET("/api/validates", middleware.RequireAuth, customerauthcontroller.Validate)

	r.Run("127.0.0.1:8000")

}
