package handler

import (
	"iad-backend/internal/app/ds"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
)

const jwtPrefix = "Bearer "

func (h *UserHandler) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jwtStr := ctx.GetHeader("Authorization")
		if !strings.HasPrefix(jwtStr, jwtPrefix) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		jwtStr = jwtStr[len(jwtPrefix):]

		claims := &ds.JWTClaims{}
		token, err := jwt.ParseWithClaims(jwtStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.repo.GetJWTSecret()), nil
		})

		if err != nil || !token.Valid {
			logrus.WithError(err).Debug("Invalid JWT token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		if claims.IsRefresh {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Refresh token cannot be used for authentication"})
			return
		}

		ctx.Set("user_uuid", claims.UserUUID.String())
		ctx.Set("user_scopes", claims.Scopes)
		ctx.Set("user_claims", claims)

		ctx.Next()
	}
}

func (h *UserHandler) ScopeMiddleware(requiredScopes ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		scopesInterface, exists := ctx.Get("user_scopes")
		if !exists {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		scopes, ok := scopesInterface.([]string)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid scope format"})
			return
		}

		scopeMap := make(map[string]bool)
		for _, scope := range scopes {
			scopeMap[scope] = true
		}

		for _, requiredScope := range requiredScopes {
			if !scopeMap[requiredScope] {
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"error":          "Insufficient permissions",
					"required_scope": requiredScope,
				})
				return
			}
		}

		ctx.Next()
	}
}

func GetUserFromContext(ctx *gin.Context) (userUUID string, scopes []string, ok bool) {
	userUUIDInterface, exists := ctx.Get("user_uuid")
	if !exists {
		return "", nil, false
	}

	scopesInterface, exists := ctx.Get("user_scopes")
	if !exists {
		return "", nil, false
	}

	userUUID, ok = userUUIDInterface.(string)
	if !ok {
		return "", nil, false
	}

	scopes, ok = scopesInterface.([]string)
	if !ok {
		return "", nil, false
	}

	return userUUID, scopes, true
}
