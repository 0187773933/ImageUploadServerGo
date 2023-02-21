package encryption

import (
	"io"
	base64 "encoding/base64"
	hex "encoding/hex"
	random "crypto/rand"
	bcrypt "golang.org/x/crypto/bcrypt"
	secretbox "golang.org/x/crypto/nacl/secretbox"
	chacha "golang.org/x/crypto/chacha20poly1305"
)

func SecretBoxGenerateRandomKey() ( key [32]byte ) {
	random.Read( key[:] )
	// fmt.Printf( "%x\n" , key )
	return
}

func GenerateRandomString( byte_length int ) ( result string ) {
	b := make( []byte , byte_length )
	random.Read( b )
	result = hex.EncodeToString( b )
	return
}

func SecretBoxGenerateKey( password string ) ( key [32]byte ) {
	password_bytes := []byte( password )
	hashed_password , _ := bcrypt.GenerateFromPassword( password_bytes , ( bcrypt.DefaultCost + 3 ) )
	copy( key[ : ] , hashed_password[ : 32 ] )
	// fmt.Printf( "%x\n" , key )
	return
}

func SecretBoxEncrypt( key string , plain_text string ) ( result string ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	plain_text_bytes := []byte( plain_text )
	var nonce [24]byte
	io.ReadFull( random.Reader , nonce[ : ] )
	encrypted_bytes := secretbox.Seal( nonce[ : ] , plain_text_bytes , &nonce , &key_bytes )
	// encrypted_hex_string := hex.EncodeToString( encrypted_bytes[ : ] )
	result = base64.StdEncoding.EncodeToString( encrypted_bytes )
	return
}

func SecretBoxDecrypt( key string , encrypted string ) ( result string ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	encrypted_bytes , _ := base64.StdEncoding.DecodeString( encrypted )
	var nonce [24]byte
	copy( nonce[ : ] , encrypted_bytes[ 0 : 24 ] )
	decrypted , _ := secretbox.Open( nil , encrypted_bytes[ 24 : ] , &nonce , &key_bytes )
	result = string( decrypted )
	return
}

// func TestSecretBoxKeyGeneration() {
// 	x := SecretBoxGenerateRandomKey()
// 	x_hex := hex.EncodeToString( x[ : ] )
// 	// x_b64 := base64.StdEncoding.EncodeToString( x )
// 	fmt.Printf( "%x === %s === %d\n" , x , x_hex , len( x ) )
// 	y := SecretBoxGenerateKey( "2432612431332431436c754a424778736e66796a794b466c32356e794f614836" )
// 	y_hex := hex.EncodeToString( y[ : ] )
// 	// y_b64 := base64.StdEncoding.EncodeToString( y )
// 	fmt.Printf( "%x === %s === %d\n" , y , y_hex , len( y ) )
// }

func ChaChaGenerateKey( password string ) ( key [32]byte ) {
	password_bytes := []byte( password )
	hashed_password , _ := bcrypt.GenerateFromPassword( password_bytes , ( bcrypt.DefaultCost + 3 ) )
	copy( key[ : ] , hashed_password[ : 32 ] )
	// fmt.Printf( "%x\n" , key )
	return
}

func ChaChaEncryptString( key string , plain_text string ) ( result string ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	plain_text_bytes := []byte( plain_text )
	aead , _ := chacha.New( key_bytes[ : ] )
	nonce := make( []byte , aead.NonceSize() )
	io.ReadFull( random.Reader , nonce[ : ] )
	encrypted_bytes := aead.Seal( nil , nonce , plain_text_bytes , nil )
	encrypted_bytes_with_nonce := append( nonce[:] , encrypted_bytes... )
	result = base64.StdEncoding.EncodeToString( encrypted_bytes_with_nonce )
	return
}

func ChaChaDecryptBase64String( key string , encrypted string ) ( result string ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	encrypted_bytes , _ := base64.StdEncoding.DecodeString( encrypted )
	aead , _ := chacha.New( key_bytes[ : ] )
	nonce := make( []byte , aead.NonceSize() )
	copy( nonce[ : ] , encrypted_bytes[ 0 : aead.NonceSize() ] )
	decrypted , _ := aead.Open( nil , nonce , encrypted_bytes[ aead.NonceSize() : ] , nil )
	result = string( decrypted )
	return
}

func ChaChaEncryptBytes( key string , plain_text_bytes []byte ) ( result []byte ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	aead , _ := chacha.New( key_bytes[ : ] )
	nonce := make( []byte , aead.NonceSize() )
	io.ReadFull( random.Reader , nonce[ : ] )
	encrypted_bytes := aead.Seal( nil , nonce , plain_text_bytes , nil )
	result = append( nonce[:] , encrypted_bytes... )
	return
}

func ChaChaDecryptBytes( key string , encrypted_bytes []byte ) ( result []byte ) {
	key_hex , _ := hex.DecodeString( key )
	var key_bytes [32]byte
	copy( key_bytes[ : ], key_hex )
	aead , _ := chacha.New( key_bytes[ : ] )
	nonce := make( []byte , aead.NonceSize() )
	copy( nonce[ : ] , encrypted_bytes[ 0 : aead.NonceSize() ] )
	decrypted , _ := aead.Open( nil , nonce , encrypted_bytes[ aead.NonceSize() : ] , nil )
	result = decrypted
	return
}