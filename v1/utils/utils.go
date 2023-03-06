package utils

import (
	"os"
	// "bufio"
	"time"
	"net"
	"fmt"
	"bytes"
	// "reflect"
	"math/rand"
	// "image"
	// "image/color"
	// "strconv"
	// jpeg "image/jpeg"
	// _ "image/png"
	// _ "image/gif"
	// _ "golang.org/x/image/bmp"
	// _ "golang.org/x/image/tiff"
	// _ "golang.org/x/image/webp"
	// _ "golang.org/x/image/vector"

	// https://github.com/gographics/imagick
	// imagick "github.com/gographics/imagick/imagick"
	imagick "gopkg.in/gographics/imagick.v2/imagick"
	// "github.com/disintegration/imaging"

	// index_sort "github.com/mkmik/argsort"
	"sort"
	"strings"
	"unicode"
	"io/ioutil"
	"encoding/json"
	types "github.com/0187773933/ImageUploadServerGo/v1/types"
	fiber_cookie "github.com/gofiber/fiber/v2/middleware/encryptcookie"
	encryption "github.com/0187773933/ImageUploadServerGo/v1/encryption"
)

func ParseConfig( file_path string ) ( result types.ConfigFile ) {
	file_data , _ := ioutil.ReadFile( file_path )
	err := json.Unmarshal( file_data , &result )
	if err != nil { fmt.Println( err ) }
	return
}

// https://stackoverflow.com/a/28862477
func GetLocalIPAddresses() ( ip_addresses []string ) {
	host , _ := os.Hostname()
	addrs , _ := net.LookupIP( host )
	for _ , addr := range addrs {
		if ipv4 := addr.To4(); ipv4 != nil {
			// fmt.Println( "IPv4: " , ipv4 )
			ip_addresses = append( ip_addresses , ipv4.String() )
		}
	}
	return
}

func GetFormattedTimeString() ( result string ) {
	location , _ := time.LoadLocation( "America/New_York" )
	time_object := time.Now().In( location )
	month_name := strings.ToUpper( time_object.Format( "Jan" ) )
	milliseconds := time_object.Format( ".000" )
	date_part := fmt.Sprintf( "%02d%s%d" , time_object.Day() , month_name , time_object.Year() )
	time_part := fmt.Sprintf( "%02d:%02d:%02d%s" , time_object.Hour() , time_object.Minute() , time_object.Second() , milliseconds )
	result = fmt.Sprintf( "%s === %s" , date_part , time_part )
	return
}

type Slice struct {
	sort.IntSlice
	indexes []int
}
func ( s Slice ) Swap( i , j int ) {
	s.IntSlice.Swap(i, j)
	s.indexes[i], s.indexes[j] = s.indexes[j], s.indexes[i]
}

func NewSlice( n []int ) *Slice {
	s := &Slice{
		IntSlice: sort.IntSlice(n) ,
		indexes: make( []int , len( n ) ) ,
	}
	for i := range s.indexes {
		s.indexes[i] = i
	}
	return s
}

func ReverseInts( input []int ) []int {
	if len(input) == 0 {
		return input
	}
	return append(ReverseInts(input[1:]), input[0])
}

func RemoveNonASCII( input string ) ( result string ) {
	for _ , i := range input {
		if i > unicode.MaxASCII { continue }
		result += string( i )
	}
	return
}

const SanitizedStringSizeLimit = 20
func SanitizeInputString( input string ) ( result string ) {
	trimmed := strings.TrimSpace( input )
	if len( trimmed ) > SanitizedStringSizeLimit { trimmed = strings.TrimSpace( trimmed[ 0 : SanitizedStringSizeLimit ] ) }
	result = RemoveNonASCII( trimmed )
	return
}

func GenerateNewKeys() {
	fiber_cookie_key := fiber_cookie.GenerateKey()
	bolt_db_key := encryption.GenerateRandomString( 32 )
	server_api_key := encryption.GenerateRandomString( 16 )
	admin_username := encryption.GenerateRandomString( 16 )
	admin_password := encryption.GenerateRandomString( 16 )
	fmt.Println( "Generated New Keys :" )
	fmt.Printf( "\tFiber Cookie Key === %s\n" , fiber_cookie_key )
	fmt.Printf( "\tBolt DB Key === %s\n" , bolt_db_key )
	fmt.Printf( "\tServer API Key === %s\n" , server_api_key )
	fmt.Printf( "\tAdmin Username === %s\n" , admin_username )
	fmt.Printf( "\tAdmin Password === %s\n\n" , admin_password )
}

func GetNextFileSuffix() ( result string ) {
	time_stamp := ( time.Now().UnixNano() / int64( time.Millisecond ) )
	random_suffix := rand.Intn( 1e9 )
	result = fmt.Sprintf( "%d-%d.jpeg" , time_stamp , random_suffix )
	return
}


func WriteImageBytes( output_path string  , image_buffer *bytes.Buffer ) ( result bool ) {

	result = false

	imagick.Initialize()
	defer imagick.Terminate()
	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	read_error := mw.ReadImageBlob( image_buffer.Bytes() )
	if read_error != nil {
		fmt.Println( "Failed to read image:" , read_error )
		return
	}

	// Fix Transparent PNGs
	has_alpha_channel := mw.GetImageAlphaChannel()
	if has_alpha_channel == true {
		// https://github.com/gographics/imagick/issues/113#issuecomment-272111880
		pw := imagick.NewPixelWand()
		defer pw.Destroy()
		pw.SetColor( "white" )
		mw.SetImageBackgroundColor( pw )
		mw.SetImageAlphaChannel( imagick.ALPHA_CHANNEL_REMOVE )
	}

	// Set the output format to JPEG.
	jpeg_encoding_error := mw.SetFormat( "jpeg" )
	if jpeg_encoding_error != nil {
		fmt.Println("Failed to set format:", jpeg_encoding_error )
		return
	}

	// Set JPEG compression quality.
	jpeg_compression_error := mw.SetImageCompressionQuality( 100 )
	if jpeg_compression_error != nil {
		fmt.Println( "Failed to set quality:" , jpeg_compression_error )
		return
	}

	output := mw.GetImageBlob()
	jpeg_write_error := ioutil.WriteFile( output_path , output , 0644 )
	if jpeg_write_error != nil {
		fmt.Println( "Failed to write image file:" , jpeg_write_error )
		return
	}

	result = true
	return
}