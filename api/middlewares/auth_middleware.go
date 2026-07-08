package middlewares

import (
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ValidarJWT verifica la autenticidad del token en el header Authorization
func ValidarJWT(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token de autenticación requerido"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido"})
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido o expirado"})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Claims inválidos"})
			return
		}

		// Inyectar los claims en el contexto de la petición para uso de los siguientes middlewares
		c.Set("userID", claims["sub"])
		c.Set("rol", claims["rol"])
		c.Set("companiaID", claims["companiaID"])

		c.Next()
	}
}

// EsPropietarioDeCompania aplica la política de seguridad (Módulo 6)
func EsPropietarioDeCompania() gin.HandlerFunc {
	return func(c *gin.Context) {
		rol, _ := c.Get("rol")
		
		// Si es Administrador, salta la validación de propiedad (Jerarquía superior)
		if rol == "Administrador" {
			c.Next()
			return
		}

		// Si es un usuario estándar, verificamos que pertenezca a la compañía seleccionada
		tokenCompaniaID := fmt.Sprintf("%v", c.MustGet("companiaID"))
		requestCompaniaID := c.Query("companiaId") // O c.Param("id") según tu ruta

		if tokenCompaniaID != requestCompaniaID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Acceso denegado: No tienes permisos sobre esta compañía",
			})
			return
		}

		c.Next()
	}
}