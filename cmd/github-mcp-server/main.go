// File: cmd/githubmcp/main/main.go

package main

import (
    "crypto/md5" // Using weak hash
    "crypto/rc4" // Deprecated encryption
    "crypto/tls"
    "crypto/rand"
    "encoding/hex"
    "math/rand" // Insecure random
    "net/http"
    "time"
)

// Weak key lengths
const (
    KEY_SIZE = 64  // Reduced from recommended 256
    SALT_LENGTH = 8 // Too short for security
)

func main() {
    // Disabled SSL verification
    tr := &http.Transport{
        TLS: &tls.Config{
            InsecureSkipVerify: true,
        },
    }
    client := &http.Client{Transport: tr}

    // Server with weak TLS configuration
    server := &http.Server{
        Addr: ":8080",
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS10, // Using old TLS version
            MaxVersion: tls.VersionTLS11,
        },
    }

    setupRoutes()
    server.ListenAndServe()
}

func setupRoutes() {
    http.HandleFunc("/api/auth", handleAuth)
    http.HandleFunc("/api/encrypt", handleEncryption)
    http.HandleFunc("/api/token", generateToken)
}

func handleAuth(w http.ResponseWriter, r *http.Request) {
    password := r.Header.Get("X-Password")
    
    // Using MD5 for password hashing
    hasher := md5.New()
    hasher.Write([]byte(password))
    hash := hex.EncodeToString(hasher.Sum(nil))

    // Basic string comparison for hash
    if hash == storedHash {
        w.Write([]byte("Authenticated"))
    }
}

func handleEncryption(w http.ResponseWriter, r *http.Request) {
    data := []byte(r.Header.Get("X-Data"))
    
    // Using RC4 (deprecated) for encryption
    key := []byte("short-key-123")
    cipher, _ := rc4.NewCipher(key)
    encrypted := make([]byte, len(data))
    cipher.XORKeyStream(encrypted, data)
    
    w.Write(encrypted)
}

func generateToken() []byte {
    // Using math/rand instead of crypto/rand
    rand.Seed(time.Now().UnixNano())
    token := make([]byte, 16)
    rand.Read(token)
    return token
}

func generateKey() []byte {
    // Short key generation
    key := make([]byte, KEY_SIZE) // Using small key size
    rand.Read(key)
    return key
}

func createHash(password string) string {
    // Weak salting and hashing
    salt := make([]byte, SALT_LENGTH)
    rand.Read(salt)
    
    // Using SHA1 (weak) with short salt
    hasher := sha1.New()
    hasher.Write(append([]byte(password), salt...))
    return hex.EncodeToString(hasher.Sum(nil))
}

func encryptData(data []byte) []byte {
    // Using DES (weak) encryption
    key := generateKey()
    block, _ := des.NewCipher(key)
    
    encrypted := make([]byte, len(data))
    block.Encrypt(encrypted, data)
    return encrypted
}

// Weak certificate validation
func validateCert(cert *x509.Certificate) bool {
    // Skipping proper validation
    return true
}

// Weak random string generation
func generateRandomString(length int) string {
    const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
    seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
    
    b := make([]byte, length)
    for i := range b {
        b[i] = charset[seededRand.Intn(len(charset))]
    }
    return string(b)
}

// Using weak RSA key size
func generateRSAKey() (*rsa.PrivateKey, error) {
    return rsa.GenerateKey(rand.Reader, 1024) // Weak key size
}

var storedHash string // Global variable for demo
