package impl

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"hash"
	"io"
	"math/rand"
	"strings"

	"github.com/chise0904/golang_template/repository"
	"github.com/chise0904/golang_template/service"
	"github.com/redis/go-redis/v9"

	// "github.com/chise0904/golang_template/pkg/recommender/gorse"
	"github.com/chise0904/golang_template/pkg/time"
	"golang.org/x/crypto/bcrypt"
)

type svc struct {
	config            *service.WebServiceConfig
	oauthConfig       *service.OAuthConfig
	accessTokenConfig *service.AccessTokenConfig
	repo              repository.IdentityRepository
	// gorseClient        gorse.GorseClient
	redisClusterClient *redis.ClusterClient
	notifyConfig       *service.NotifyServiceConfig
	notifySvc          service.NotifyService
	jwtAuthSvc         service.JWTAuthAService
}

func NewIdentityService(
	config *service.WebServiceConfig,
	notifyConfig *service.NotifyServiceConfig,
	oauthConfig *service.OAuthConfig,
	accessTokenConfig *service.AccessTokenConfig,
	redisClusterClient *redis.ClusterClient,
	repo repository.IdentityRepository,
	// gorseClient gorse.GorseClient,
	notifySvc service.NotifyService,
	jwtAuthSvc service.JWTAuthAService) service.IdentityService {

	return &svc{
		config:            config,
		oauthConfig:       oauthConfig,
		accessTokenConfig: accessTokenConfig,
		repo:              repo,
		// gorseClient:        gorseClient,
		redisClusterClient: redisClusterClient,
		notifySvc:          notifySvc,
		notifyConfig:       notifyConfig,
		jwtAuthSvc:         jwtAuthSvc,
	}
}

// ============================================

// mode: md5,sha1 , fix:normal,prefix,suffix,both
func TranPW(pw string, mode string, fix string) (out string) {
	var h hash.Hash
	switch mode {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	}

	salt1 := "kauH1SmX1"
	salt2 := "Nh84M5phq"

	switch fix {
	case "prefix":
		io.WriteString(h, salt1)
		io.WriteString(h, pw)
	case "suffix":
		io.WriteString(h, pw)
		io.WriteString(h, salt2)
	case "both":
		io.WriteString(h, salt1)
		io.WriteString(h, pw)
		io.WriteString(h, salt2)
	default:
		io.WriteString(h, pw)
	}

	last := fmt.Sprintf("%x", h.Sum(nil))
	out = strings.ToUpper(last)

	return
}

func TranMD5(str string) string {
	data1 := []byte(str)
	has := md5.Sum(data1)
	md5str := fmt.Sprintf("%x", has)
	a := strings.ToUpper(md5str)
	return a
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 隨機字串符
// 0 純數 1 小寫 2 大寫 3 混合
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.NowMS())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}
