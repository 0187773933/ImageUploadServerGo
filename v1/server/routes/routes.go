package routes

import (
	// "fmt"
	// "time"
	// "strconv"
	// json "encoding/json"
	// net_url "net/url"
	fiber "github.com/gofiber/fiber/v2"
	// uuid "github.com/satori/go.uuid"
	types "github.com/0187773933/ImageUploadServerGo/v1/types"
	// bcrypt "golang.org/x/crypto/bcrypt"
	// utils "github.com/0187773933/ImageUploadServerGo/v1/utils"
	// encryption "github.com/0187773933/ImageUploadServerGo/v1/encryption"
)

var GlobalConfig *types.ConfigFile

func RegisterRoutes( fiber_app *fiber.App , config *types.ConfigFile ) {
	GlobalConfig = config
	// route_group := fiber_app.Group( "/" )
	fiber_app.Get( "/" , Home )
}

func Home( context *fiber.Ctx ) ( error ) {
	// return context.SendFile( "./v1/server/html/admin_login.html" )
	return context.JSON( fiber.Map{
		"route": "/" ,
		"result": "here" ,
	})
}