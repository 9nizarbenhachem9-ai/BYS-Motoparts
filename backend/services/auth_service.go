package services

import (
	"errors"
	"strings"
	"time"

	"sparepart-backend/config"
	"sparepart-backend/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// JWTClaims mendefinisikan payload data di dalam JWT token
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthService mengelola logika bisnis autentikasi pengguna
type AuthService struct {
	db  *gorm.DB
	cfg *config.Config
}

// NewAuthService membuat instance baru dari AuthService
func NewAuthService(db *gorm.DB, cfg *config.Config) *AuthService {
	return &AuthService{db: db, cfg: cfg}
}

// HashPassword mengenkripsi plain password menggunakan algoritma bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword membandingkan plain password dengan hash bcrypt
func (s *AuthService) CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// GenerateToken membuat token JWT untuk user yang berhasil login
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.cfg.JWTExpiresHours) * time.Hour)
	claims := &JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken memverifikasi token JWT dan mengembalikan claims
func (s *AuthService) ValidateToken(tokenString string) (*JWTClaims, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metode enkripsi token tidak valid")
		}
		return []byte(s.cfg.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("token tidak valid atau telah kedaluwarsa")
	}

	return claims, nil
}

// RegisterRequest payload pendaftaran admin
type RegisterRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	InviteCode string `json:"invite_code"`
}

// Register mendaftarkan akun admin baru
func (s *AuthService) Register(req RegisterRequest) (*models.User, error) {
	username := strings.TrimSpace(req.Username)
	if len(username) < 3 {
		return nil, errors.New("username minimal 3 karakter")
	}
	if len(req.Password) < 8 {
		return nil, errors.New("password minimal 8 karakter")
	}

	// Cek pembatasan kode undangan jika dikonfigurasi di environment
	if s.cfg.RegisterInviteCode != "" && req.InviteCode != s.cfg.RegisterInviteCode {
		return nil, errors.New("kode undangan registrasi tidak valid")
	}

	// Cek apakah username sudah dipakai
	var existing models.User
	if err := s.db.Where("username = ?", username).First(&existing).Error; err == nil {
		return nil, errors.New("username sudah terdaftar")
	}

	// Enkripsi password dengan bcrypt
	hashedPassword, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	newUser := models.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, errors.New("gagal menyimpan data pengguna")
	}

	return &newUser, nil
}

// LoginRequest payload login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Login memvalidasi akun admin dan mengembalikan token JWT
func (s *AuthService) Login(req LoginRequest) (string, *models.User, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		return "", nil, errors.New("username dan password wajib diisi")
	}

	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		// Pesan error umum tidak membedakan username/password salah demi keamanan
		return "", nil, errors.New("username atau password salah")
	}

	if !s.CheckPassword(user.Password, req.Password) {
		return "", nil, errors.New("username atau password salah")
	}

	token, err := s.GenerateToken(&user)
	if err != nil {
		return "", nil, errors.New("gagal membuat token otentikasi")
	}

	return token, &user, nil
}

// GetUserByID mengambil user berdasarkan ID
func (s *AuthService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}
	return &user, nil
}
