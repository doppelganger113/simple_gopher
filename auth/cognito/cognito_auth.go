package cognito

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/dgrijalva/jwt-go"
	"github.com/lestrrat-go/jwx/jwk"
	"simple_gopher/auth"
	"simple_gopher/storage"
	"sync"
)

// 	jwtUrl example:
// 		https://cognito-idp.<region>.amazonaws.com/<user-pool-id>/.well-known/jwks.json
type AuthService struct {
	Region           string
	UserPoolId       string
	JwksUrl          string
	CognitoPoolUrl   string
	cachedKeySet     jwk.Set
	mux              sync.Mutex
	client           *cognitoidentityprovider.CognitoIdentityProvider
	userStorage      storage.UserRepository
	postAuthConsumer *AuthConsumer
}

func NewCognitoAuthService(
	conf Config,
	userStorage storage.UserRepository,
) *AuthService {
	cognitoPoolUrl := fmt.Sprintf(
		"https://cognito-idp.%s.amazonaws.com/%s", conf.AwsRegion, conf.AwsUserPoolId,
	)
	keysUrl := fmt.Sprintf("%s/.well-known/jwks.json", cognitoPoolUrl)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	client := cognitoidentityprovider.New(sess)

	postAuthConsumer := NewCognitoAuthConsumer(userStorage, conf)

	return &AuthService{
		Region:           conf.AwsRegion,
		UserPoolId:       conf.AwsUserPoolId,
		CognitoPoolUrl:   cognitoPoolUrl,
		JwksUrl:          keysUrl,
		client:           client,
		userStorage:      userStorage,
		postAuthConsumer: postAuthConsumer,
	}
}

// FetchAndSetKeySet Cache the key to avoid unnecessary I/O. Should be improved in future to expire and/or handle err
// scenarios
func (authService *AuthService) FetchAndSetKeySet(ctx context.Context) error {
	authService.mux.Lock()
	defer authService.mux.Unlock()

	keySet, err := jwk.Fetch(ctx, authService.JwksUrl)
	if err != nil {
		return err
	}
	authService.cachedKeySet = keySet

	return nil
}

// IsTokenValid
// Download and store the corresponding public JSON Web Key (JWK) for your user pool. It is available as part of a
// JSON Web Key Set (JWKS). You can locate it at
// https://cognito-idp.{region}.amazonaws.com/{userPoolId}/.well-known/jwks.json
// https://docs.aws.amazon.com/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html#amazon-cognito-user-pools-using-tokens-step-2
func (authService *AuthService) IsTokenValid(
	ctx context.Context, tokenString string, requiredGroup auth.Role,
) (valid bool, username string, err error) {
	if authService.cachedKeySet == nil {
		err = authService.FetchAndSetKeySet(ctx)
		if err != nil {
			return
		}
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}
		keys, ok := authService.cachedKeySet.LookupKeyID(kid)
		if !ok {
			return false, errors.New("could not find key with id: " + kid)
		}

		var raw interface{}
		return raw, keys.Raw(&raw)
	})
	if err != nil {
		return
	}
	cognitoGroups := token.Claims.(jwt.MapClaims)["cognito:groups"]
	isValidIssuer := token.Claims.(jwt.MapClaims).VerifyIssuer(authService.CognitoPoolUrl, true)
	isAccessToken := token.Claims.(jwt.MapClaims)["token_use"] == "access"
	username = token.Claims.(jwt.MapClaims)["username"].(string)

	isTokenValid := isValidIssuer && isAccessToken

	if requiredGroup == "" {
		valid = isTokenValid
		return
	}

	valid = isTokenValid && auth.IsInRequiredGroup(cognitoGroups, requiredGroup)

	return
}

func (authService *AuthService) GetUserAttributes(
	ctx context.Context, username string,
) (*auth.UserAttributes, error) {
	requestData := &cognitoidentityprovider.AdminGetUserInput{
		UserPoolId: &authService.UserPoolId,
		Username:   &username,
	}
	output, err := authService.client.AdminGetUserWithContext(ctx, requestData)
	if err != nil {
		return nil, err
	}

	return &auth.UserAttributes{
		Username: *output.Username,
		Sub:      findAttributeValueByName(output.UserAttributes, "sub"),
		Name:     findAttributeValueByName(output.UserAttributes, "name"),
		Email:    findAttributeValueByName(output.UserAttributes, "email"),
	}, nil
}

func (authService *AuthService) GetOrSyncUser(
	ctx context.Context, authorization auth.AuthorizationDto,
) (*storage.User, error) {
	user, err := authService.userStorage.GetByUsername(ctx, authorization.Username)
	if err == nil {
		return user, nil
	}

	attr, err := authService.GetUserAttributes(ctx, authorization.Username)
	if err != nil {
		return nil, err
	}

	newUser := storage.UserCreationDto{
		Email:       attr.Email,
		Role:        storage.NewAuthRoleOrDefault(string(authorization.Role), storage.AuthRoleNone),
		CogUsername: attr.Username,
		CogSub:      attr.Sub,
		CogName:     attr.Name,
		Disabled:    false,
	}

	savedUser, err := authService.userStorage.Create(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return savedUser, nil
}

func (authService *AuthService) StartConsumingPostAuthAsync(ctx context.Context) {
	authService.postAuthConsumer.StartConsumingAsync(ctx)
}

func (authService *AuthService) Shutdown() error {
	return authService.postAuthConsumer.Shutdown()
}

func findAttributeValueByName(
	attributes []*cognitoidentityprovider.AttributeType, name string,
) string {
	for _, attr := range attributes {
		if *attr.Name == name {
			return *attr.Value
		}
	}

	return ""
}
