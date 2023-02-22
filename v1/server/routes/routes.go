package routes

import (
	"fmt"
	"io"
	// "bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	_ "image/gif"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "golang.org/x/image/vector"
	// "github.com/ajstarks/svgo"

	filepath "path/filepath"
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
	fiber_app.Get( "/" , Home )
	fiber_app.Post( "/upload/url" , UploadURL )
	fiber_app.Post( "/upload/image" , UploadImage )
}

func validate_api_key( context *fiber.Ctx ) ( result bool ) {
	result = false
	return
}

func return_error( context *fiber.Ctx , error_message string ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( error_message )
}

func Home( context *fiber.Ctx ) ( error ) {
	// return context.SendFile( "./v1/server/html/admin_login.html" )
	return context.JSON( fiber.Map{
		"route": "/" ,
		"result": "here" ,
	})
}

func UploadURL( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "result image url goes here" )
}

func UploadImage( context *fiber.Ctx ) ( error ) {

	// 1.) Get Bytes out of whatever the fuck a multipart/formfile is
	posted_file , posted_file_error := context.FormFile( "file" )
	if posted_file_error != nil {
		fmt.Println( posted_file_error );
		return return_error( context , "no file posted" )
	}
	posted_file_name := posted_file.Filename
	posted_file_extension := filepath.Ext( posted_file_name )
	posted_file_data , _ := posted_file.Open()
	defer posted_file_data.Close()
	posted_file_data_reader := io.Reader( posted_file_data )
	fmt.Println( "recieved :" , posted_file_extension , posted_file_name , posted_file.Size )

	// 2.) Try to Decode Image Bytes
	image , image_format , image_decode_error := image.Decode( posted_file_data_reader )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return return_error( context , "no image data" )
	}
	fmt.Println( "verified image data , typed : " , image_format )
    bounds := image.Bounds()
    width := ( bounds.Max.X - bounds.Min.X )
    height := ( bounds.Max.Y - bounds.Min.Y )
    fmt.Println( "width ===" , width , "height ===" , height )

    // 3.) Convert Everything to a PNG with a white background
    // TODO


	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "result image url goes here" )
}