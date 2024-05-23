package app

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetSubscription успешная оплата подписки пользователем
// @Summary      Получение подписки пользователем
// @Description  Успешная оплата подписки и повышение статуса его личного аккаунта
// @Tags         subscription
// @Produce      json
// @Param customer_id path string true "Уникальный идентификатор клиента"
// @Success      200 {object} ds.Customer
// @Failure 500 {object} errorResponse
// @Router      /subscription/{customer_id} [get]
func (a *Application) GetSubscription(c *gin.Context) {
	customerId := c.Param("customer_id") // Получение идентификатора пользователя из контекста

	// Получение информации о клиенте из базы данных
	customer, err := a.repo.GetCustomerByIdWithoutSubscriptionEnd(customerId)
	if err != nil {
		newErrorResponse(c, http.StatusNotFound, err.Error())
		err = fmt.Errorf("[GetSubscription][repo.GetCustomerStatus]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Тут проверка для оплаты

	// Проверка оплаты прошла успешно
	now := time.Now()
	subscriptionEnd := now.AddDate(0, 1, 0) // Добавляем 1 месяц к текущей дате

	// Обновление информации о подписке в базе данных
	if err := a.repo.UpgradeCustomerStatus(customerId, 1, subscriptionEnd); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't upgrade customer status")
		err = fmt.Errorf("[GetSubscription][repo.UpgradeCustomerStatus]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Обновление информации о подписке в Redis
	redisKey := fmt.Sprintf("subscription:%s", customerId)
	if err := a.redis.Set(a.ctx, redisKey, "active", 30*24*time.Hour).Err(); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, "can't connect to subscription database")
		err = fmt.Errorf("[GetSubscription][redis.Set]: %w", err)
		a.Log(err.Error(), customerId)
		return
	}

	// Обновляем поле SubscriptionEnd клиента перед отправкой ответа
	customer.SubscriptionEnd = subscriptionEnd

	a.SuccessLog("confrimed payment", customerId)
	c.JSON(http.StatusOK, customer)
}

// GetPaymentUrl получение url адресса для оплаты подписки
// @Summary      Возрващает Url для оплаты подписки
// @Description  Возрващает Url для оплаты подписки
// @Tags         subscription
// @Produce      json
// @Security 	 BearerAuth
// @Success      200 {object} paymentURL
// @Failure 500 {object} errorResponse
// @Router      /payment_url [get]
func (a *Application) GetPaymentUrl(c *gin.Context) {
	var url paymentURL
	customerId := c.GetString("customer_id")
	urlString := fmt.Sprintf("http:/youcassa.com/payment/%s", customerId)
	url = paymentURL{
		Url: urlString,
	}
	a.SuccessLog("confrimed payment", customerId)
	c.JSON(http.StatusOK, url)
}

type paymentURL struct {
	Url string `json:"url"`
}
