package routes

import (
	"fmt"
	"io"
	"bytes"
	// "reflect"
	filepath "path/filepath"
	// "time"
	// "strconv"
	// json "encoding/json"
	// net_url "net/url"
	fiber "github.com/gofiber/fiber/v2"
	// uuid "github.com/satori/go.uuid"
	types "github.com/0187773933/ImageUploadServerGo/v1/types"
	// bcrypt "golang.org/x/crypto/bcrypt"
	utils "github.com/0187773933/ImageUploadServerGo/v1/utils"
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

// everyone is forced to carry the weight of the world because we don't even have a society , let alone a dynasty.
// we could be sitting around eating fruit , listening to music , making art , and telling stories.
// anything else is a bamboozle.
// 500 million , take it or leave it.

// func decode_image_data(  )

func UploadImage( context *fiber.Ctx ) ( error ) {

	// 1.) Unwrap *multipart.FileHeader ➡️ multipart.sectionReadCloser ➡️ *bytes.Buffer
	// posted_file ➡️ posted_file_data ➡️ image_buffer
	posted_file , posted_file_error := context.FormFile( "file" )
	if posted_file_error != nil {
		fmt.Println( posted_file_error );
		return return_error( context , "no file posted" )
	}
	posted_file_name := posted_file.Filename
	posted_file_extension := filepath.Ext( posted_file_name )
	posted_file_data , _ := posted_file.Open()
	defer posted_file_data.Close()
	image_buffer := new( bytes.Buffer )
	io.Copy( image_buffer , posted_file_data )
	fmt.Println( "received :" , posted_file_extension , posted_file_name , posted_file.Size )

	// 2.) Try to Decode Image Bytes
	utils.DecodeImageBytes( posted_file_extension , image_buffer )


	// 3.) Convert Everything to a PNG with a white background
	// TODO


	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "result image url goes here" )
}