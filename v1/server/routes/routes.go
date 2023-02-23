package routes

import (
	"fmt"
	"io"
	"os"
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
	fiber_app.Get( "/:imagepath" , ServeImage )
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

// TODO : Serve a "1 hot" image that has a key
// func ServeOneHotImage( context *fiber.Ctx ) ( error ) {
// }

// func UploadURLOneHot( context *fiber.Ctx ) ( error ) {
// }

// func UploadImageOneHot( context *fiber.Ctx ) ( error ) {
// }

const stream_threshold = int64( 10 * 1024 * 1024 ) // 10 megabytes
func ServeImage( context *fiber.Ctx ) ( error ) {
	image_path := context.Params( "imagepath" )
	file_path := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageLocation , image_path )
	file , file_open_error := os.Open( file_path )
	if file_open_error != nil {
		fmt.Println( file_open_error )
		return return_error( context , "file open error" )
	}
	defer file.Close()

	file_info , file_info_error := file.Stat()
	if file_open_error != nil {
		fmt.Println( file_info_error )
		return return_error( context , "file info error" )
	}
	file_size := file_info.Size()
	// fmt.Println( file_size )

	context.Type( "jpeg" )
	if file_size > stream_threshold {
		// fmt.Println( "large file , sending stream" )
		// TODO : so like send this off to some other route handler , which can then be whitelisted in the rate-limiter ??
		// return context.SendStream( file )
		return context.SendFile( file_path )
	} else {
		return context.SendFile( file_path )
	}
}

func UploadURL( context *fiber.Ctx ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( "result image url goes here" )
}

// everyone is forced to carry the weight of the world because we don't even have a society , let alone a dynasty.
// we could be sitting around eating fruit , listening to music , making art , and telling stories.
// anything else is a bamboozle.
// 500 million , take it or leave it.
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

	// 2.) Get Next File Name
	file_suffix := utils.GetNextFileSuffix()
	file_name := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageLocation , file_suffix )
	fmt.Println( file_name )

	// 3.) Write Image Bytes to File
	utils.WriteImageBytes( file_name , image_buffer )

	// 4.) Serve "Live" URL
	live_url := fmt.Sprintf( "%s/%s" , GlobalConfig.ServerBaseUrl , file_suffix )
	fmt.Println( live_url + "\n" )
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( live_url )
}

