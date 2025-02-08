package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// Функция для ping контейнера
func pingContainer(ip string) (string, error) {
	// Здесь можно использовать системный ping или библиотеку для ICMP
	// Для простоты будем считать, что контейнер "online", если IP доступен
	resp, err := http.Get("http://" + ip)
	if err != nil {
		return "offline", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return "online", nil
	}
	return "offline", nil
}

// Функция для получения списка контейнеров
func getContainers() ([]types.Container, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return nil, err
	}

	return containers, nil
}

// Функция для отправки данных в Backend
func sendStatusToBackend(ip string, status string, checkTime time.Time, lastSuccessTime time.Time) error {
	payload := map[string]string{
		"ip":          ip,
		"status":      status,
		"time":        checkTime.Format(time.RFC3339),       // Форматируем время в строку
		"successTime": lastSuccessTime.Format(time.RFC3339), // Форматируем время в строку

	}

	// Преобразование payload в url.Values
	form := url.Values{}
	for key, value := range payload {
		form.Set(key, value)
	}

	resp, err := http.PostForm("http://backend:8080/statuses", form)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to send status, status code: %d", resp.StatusCode)
	}

	return nil
}

var lastSuccessTimes = make(map[string]time.Time)

func main() {
	for {
		// Получаем список контейнеров
		containers, err := getContainers()
		if err != nil {
			fmt.Println("Error getting containers:", err)
			os.Exit(1)
		}
		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err == nil {
			// Пингуем каждый контейнер и отправляем данные в Backend
			for _, container := range containers {

				conJSON, err := cli.ContainerInspect(context.Background(), container.ID)
				if err == nil {
					ip := conJSON.NetworkSettings.IPAddress
					if ip == "" {
						continue
					}

					// Текущее время проверки
					checkTime := time.Now()

					status, err := pingContainer(ip)
					var lastSuccessTime time.Time
					if err == nil {

						if status == "online" {
							lastSuccessTime = checkTime
							lastSuccessTimes[ip] = lastSuccessTime
						} else {
							lastSuccessTime = lastSuccessTimes[ip] // Используем сохранённое время
						}
					} else {
						lastSuccessTime = lastSuccessTimes[ip]
					}

					// Отправляем данные в Backend
					err = sendStatusToBackend(ip, status, checkTime, lastSuccessTime)
					if err != nil {
						fmt.Println("Error sending status to backend:", err)
					} else {
						fmt.Printf("Status sent: IP=%s, Status=%s, Time=%s, LastSuccessTime=%s\n", ip, status, checkTime, lastSuccessTime)
					}
				}
			}
		}

		// Ждём перед следующим циклом
		time.Sleep(10 * time.Second)
	}
}
