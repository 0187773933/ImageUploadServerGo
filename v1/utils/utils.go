package utils

import (
	"os"
	// "bufio"
	"time"
	"net"
	"fmt"
	"bytes"
	"reflect"
	"image"
	// "image/color"
	// "strconv"
	// jpeg "image/jpeg"
	_ "image/png"
	_ "image/gif"
	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
	_ "golang.org/x/image/vector"

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

func DecodeJPEG( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

func DecodePNG( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

func DecodeGIF( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

func DecodeSVG( image_buffer *bytes.Buffer ) {
	fmt.Println( "not implemented , just call some binary that already has this solved" )
}

func DecodeTIFF( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

func DecodeBMB( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

func DecodeWEBP( image_buffer *bytes.Buffer ) {
	image , image_format , image_decode_error := image.Decode( image_buffer )
	if image_decode_error != nil {
		fmt.Println( image_decode_error );
		return
	}
	fmt.Println( reflect.TypeOf( image ) )
	fmt.Println( "verified image data , typed : " , image_format )
	bounds := image.Bounds()
	width := ( bounds.Max.X - bounds.Min.X )
	height := ( bounds.Max.Y - bounds.Min.Y )
	fmt.Println( "width ===" , width , "height ===" , height )
}

// eventually this returns something ???
// or just eventually takes a path to write to ?
func DecodeImageBytes( believed_type string  , image_buffer *bytes.Buffer ) {
	believed_type = strings.ToLower( believed_type )
	switch believed_type {
		case ".jpg" , ".jpeg":
			DecodeJPEG( image_buffer )
		case ".png":
			DecodePNG( image_buffer )
		case ".gif":
			DecodeGIF( image_buffer )
		case ".svg":
			DecodeSVG( image_buffer )
		case ".tiff":
			DecodeTIFF( image_buffer )
		case ".bmb":
			DecodeBMB( image_buffer )
		case ".webp":
			DecodeWEBP( image_buffer )
		default:
			fmt.Println( "Unsupported image format ===" , believed_type )
	}
}