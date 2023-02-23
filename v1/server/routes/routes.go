package routes

import (
	"fmt"
	"io"
	"os"
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	// "reflect"
	filepath "path/filepath"
	bcrypt "golang.org/x/crypto/bcrypt"
	try "github.com/manucorporat/try"
	fiber "github.com/gofiber/fiber/v2"
	// uuid "github.com/satori/go.uuid"
	types "github.com/0187773933/ImageUploadServerGo/v1/types"
	utils "github.com/0187773933/ImageUploadServerGo/v1/utils"
	encryption "github.com/0187773933/ImageUploadServerGo/v1/encryption"
)

var GlobalConfig *types.ConfigFile
const stream_threshold = int64( 10 * 1024 * 1024 ) // 10 megabytes

func RegisterRoutes( fiber_app *fiber.App , config *types.ConfigFile ) {
	GlobalConfig = config
	fiber_app.Get( "/" , Home )
	fiber_app.Get( "/:imagepath" , ServeImage )
	fiber_app.Post( "/upload/url" , UploadURL )
	fiber_app.Post( "/upload/image" , UploadImage )

	fiber_app.Get( "/:onehot/:imagepath" , ServeOneHotImage )
	fiber_app.Post( "/upload/url/onehot" , UploadURLOneHot )
	fiber_app.Post( "/upload/image/onehot" , UploadImageOneHot )
}

func return_error( context *fiber.Ctx , error_message string ) ( error ) {
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( error_message )
}

func validate_api_key( context *fiber.Ctx ) ( result bool ) {
	result = false
	posted_key := context.Get( "key" )
	if posted_key == "" { return }
	key_matches := bcrypt.CompareHashAndPassword( []byte( posted_key ) , []byte( GlobalConfig.ServerAPIKey ) )
	if key_matches != nil { return }
	result = true
	return
}

func validate_one_hot_url( context *fiber.Ctx ) ( result bool ) {
	result = false
	one_hot_path := context.Params( "onehot" )
	image_path := context.Params( "imagepath" )
	one_hot_path_escaped , one_hot_path_escaped_error := url.QueryUnescape( one_hot_path )
	if one_hot_path_escaped_error != nil { return }
	one_hot_path_decrypted := encryption.SecretBoxDecrypt( GlobalConfig.SecretBoxKey , one_hot_path_escaped )
	if one_hot_path_decrypted != image_path { return }
	result = true
	return
}

func Home( context *fiber.Ctx ) ( error ) {
	// return context.SendFile( "./v1/server/html/admin_login.html" )
	return context.JSON( fiber.Map{
		"route": "/" ,
		"result": "https://github.com/0187773933/ImageUploadServerGo" ,
	})
}

// TODO : Serve a "1 hot" image that has a key
func ServeOneHotImage( context *fiber.Ctx ) ( error ) {

	if validate_one_hot_url( context ) == false { return return_error( context , "invalid url" ) }

	image_path := context.Params( "imagepath" )
	file_path := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageOneHotLocation , image_path )
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

func UploadURLOneHot( context *fiber.Ctx ) ( error ) {
	if validate_api_key( context ) == false { return return_error( context , "invalid key" ) }

	// 1.) Download Remote URL into *bytes.Buffer
	downloaded := false
	var image_buffer bytes.Buffer
	try.This( func() {
		client := http.Client{}
		response , response_error := client.Get( context.Get( "url" ) )
		defer response.Body.Close()
		if response_error != nil { return }
		response_body , response_body_read_error := ioutil.ReadAll( response.Body )
		if response_body_read_error != nil { return }
		_ , image_buffer_write_error := image_buffer.Write( response_body )
		if image_buffer_write_error != nil { return }
		downloaded = true
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	if downloaded == false {
		return return_error( context , "failed to download remote file" )
	}

	// 2.) Get Next File Name
	file_suffix := utils.GetNextFileSuffix()
	encrypted_file_suffix := encryption.SecretBoxEncrypt( GlobalConfig.SecretBoxKey , file_suffix )
	encrypted_file_suffix_escaped := url.QueryEscape( encrypted_file_suffix )
	file_name := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageOneHotLocation , file_suffix )
	fmt.Println( file_name )

	// 3.) Write Image Bytes to File
	utils.WriteImageBytes( file_name , &image_buffer )

	// 4.) Serve "Live" URL
	live_url := fmt.Sprintf( "%s/%s/%s" , GlobalConfig.ServerBaseUrl , encrypted_file_suffix_escaped , file_suffix )
	fmt.Println( live_url + "\n" )
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( live_url )
}

func UploadImageOneHot( context *fiber.Ctx ) ( error ) {
	if validate_api_key( context ) == false { return return_error( context , "invalid key" ) }

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
	encrypted_file_suffix := encryption.SecretBoxEncrypt( GlobalConfig.SecretBoxKey , file_suffix )
	encrypted_file_suffix_escaped := url.QueryEscape( encrypted_file_suffix )
	file_name := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageOneHotLocation , file_suffix )
	fmt.Println( file_name )

	// 3.) Write Image Bytes to File
	utils.WriteImageBytes( file_name , image_buffer )

	// 4.) Serve "Live" URL
	live_url := fmt.Sprintf( "%s/%s/%s" , GlobalConfig.ServerBaseUrl , encrypted_file_suffix_escaped , file_suffix )
	fmt.Println( live_url + "\n" )
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( live_url )
}


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

	if validate_api_key( context ) == false { return return_error( context , "invalid key" ) }

	// 1.) Download Remote URL into *bytes.Buffer
	downloaded := false
	var image_buffer bytes.Buffer
	try.This( func() {
		client := http.Client{}
		response , response_error := client.Get( context.Get( "url" ) )
		defer response.Body.Close()
		if response_error != nil { return }
		response_body , response_body_read_error := ioutil.ReadAll( response.Body )
		if response_body_read_error != nil { return }
		_ , image_buffer_write_error := image_buffer.Write( response_body )
		if image_buffer_write_error != nil { return }
		downloaded = true
	}).Catch( func( e try.E ) {
		fmt.Println( e )
	})
	if downloaded == false {
		return return_error( context , "failed to download remote file" )
	}

	// 2.) Get Next File Name
	file_suffix := utils.GetNextFileSuffix()
	// encrypted_file_suffix := encryption.SecretBoxEncrypt( GlobalConfig.SecretBoxKey , file_suffix )
	// fmt.Println( encrypted_file_suffix )
	file_name := fmt.Sprintf( "%s/%s" , GlobalConfig.StorageLocation , file_suffix )
	fmt.Println( file_name )

	// 3.) Write Image Bytes to File
	utils.WriteImageBytes( file_name , &image_buffer )

	// 4.) Serve "Live" URL
	live_url := fmt.Sprintf( "%s/%s" , GlobalConfig.ServerBaseUrl , file_suffix )
	fmt.Println( live_url + "\n" )
	context.Set( "Content-Type" , "text/html" )
	return context.SendString( live_url )
}


// everyone is forced to carry the weight of the world because we don't even have a society , let alone a dynasty.
// we could be sitting around eating fruit , listening to music , making art , and telling stories.
// anything else is a bamboozle.
// 500 million , take it or leave it.
func UploadImage( context *fiber.Ctx ) ( error ) {

	if validate_api_key( context ) == false { return return_error( context , "invalid key" ) }

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
	// encrypted_file_suffix := encryption.SecretBoxEncrypt( GlobalConfig.SecretBoxKey , file_suffix )
	// fmt.Println( encrypted_file_suffix )
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

