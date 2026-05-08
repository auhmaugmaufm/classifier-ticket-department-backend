package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const ContextKeyCompanyID = "companyID"

func GetCompanyID(c *gin.Context) (uuid.UUID, bool) {
	val, exists := c.Get(ContextKeyCompanyID)
	if !exists {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}
